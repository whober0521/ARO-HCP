// copy from https://github.com/Azure/ARO-HCP/blob/main/dev-infrastructure/configurations/region.bicepparam
using '../templates/region.bicep'

// dns
param baseDNSZoneName = 'hcp.osadev.cloud'
param baseDNSZoneResourceGroup = 'global'

// maestro
param maestroKeyVaultName = 'maestro-kv-$(region)-$(user)'
param maestroEventGridNamespacesName = 'maestro-eventgrid-$(region)'
param maestroEventGridMaxClientSessionsPerAuthName = 4

// These parameters are always overriden in the Makefile
param currentUserId = ''
