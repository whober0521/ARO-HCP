# copy from https://github.com/Azure/ARO-HCP/blob/952e3ce091e7a53fdacdeef2c44bbed43e176008/maestro/Makefile#L14
deploy-server:
    TENANT_ID="72f988bf-86f1-41af-91ab-2d7cd011db47"
    REGION_RG="aro-hcp-$(region)-$(user)"
    EVENTGRID_NS="maestro-eventgrid-$(region)"
    MAESTRO_KV="maestro-kv-$(region)-$(user)"
    SERVICE_RG="aro-hcp-$(region)-svc-$(user)"
    AKS="aro-hcp-aks"
    MAESTRO_MI="maestro-server"
    HELM_CHART="../maestro/deploy/helm/server"
    TEST="dev"

    EVENTGRID_HOSTNAME=$(az event namespace show -g "${REGION_RG}" -n "${EVENTGRID_NS}" --query "properties.topicSpacesConfiguration.hostname")
    MAESTRO_MI_CLIENT_ID=$(az identity show -g "${SERVICE_RG}" -n "${MAESTRO_MI}" --query "clientId")
    ISTO_VERSION=$(az aks show -g "${SERVICE_RG}" -n "${AKS}" --query "serviceMeshProfile.istio.revisions[-1]")

    kubectl create namespace maestro --dry-run=client -o json | kubectl apply -f - && \
    kubectl label namespace maestro "istio.io/rev=$${ISTO_VERSION}" --overwrite=true && \
    helm upgrade --install maestro-server "${HELM_CHART}" \
        --namespace maestro \
        --set broker.host=$${EVENTGRID_HOSTNAME} \
        --set credsKeyVault.name=$${MAESTRO_KV} \
        --set azure.clientId=$${MAESTRO_MI_CLIENT_ID} \
        --set azure.tenantId=$${TENANT_ID} \
        --set image.base='quay.io/redhat-user-workloads/maestro-rhtap-tenant/maestro/maestro'\
        --set database.containerizedDb=true \
        --set database.ssl=disable