package poc

import (
	"context"
	"os"
	"text/template"

	"gopkg.in/yaml.v3"
)

type configProviderImpl struct {
	baseVariableOverrides *VariableOverrides
	config                string
	region                string
	user                  string
}

func Print(region, user string) error {
	config := NewConfigProvider("poc/config.yaml", region, user)

	// ===================  print config values  ===================
	// println("Clouds:")
	// cloud, err := config.GetAllClouds()
	// if err != nil {
	// 	return err
	// }

	// for _, c := range cloud {
	// 	println(c)
	// 	cloudv, err := config.GetCloudVariables(context.Background(), c)

	// 	if err != nil {
	// 		return err
	// 	}

	// 	for k, v := range cloudv {
	// 		println(k, v)
	// 	}
	// }

	// println()
	// println("DeployEnvs:")
	// envs, err := config.GetAllDeployEnvs()
	// if err != nil {
	// 	return err
	// }

	// for _, e := range envs {
	// 	println(e.Cloud, e.DeployEnv)
	// 	envv, err := config.GetDeployenvVariables(context.Background(), e.Cloud, e.DeployEnv)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	for k, v := range envv {
	// 		println(k, v)
	// 	}
	// }

	// println()
	// println("Regions:")
	// regions, err := config.GetAllRegions()
	// if err != nil {
	// 	return err
	// }

	// for _, r := range regions {
	// 	println(r.Cloud, r.DeployEnv, r.Region)

	// 	envv, err := config.GetVariables(context.Background(), r.Cloud, r.DeployEnv, r.Region)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	for k, v := range envv {
	// 		println(k, v)
	// 	}
	// }
	// ===================  print config values  ===================

	// Create output file
	output, err := os.Create("poc/output/helm.sh")
	if err != nil {
		return err
	}
	defer output.Close()

	// Parse and execute template with config values
	tmpl, err := template.New("helm.sh").ParseFiles("poc/input/helm.sh")
	if err != nil {
		return err
	}
	variables, err := config.GetDeployenvVariables(context.Background(), "public", "dev")
	if err != nil {
		return err
	}

	err = tmpl.ExecuteTemplate(output, "helm.sh", variables)
	if err != nil {
		return err
	}

	// Create output file
	output, err = os.Create("poc/output/test.bicepparam")
	if err != nil {
		return err
	}
	defer output.Close()

	// Parse and execute template with config values
	tmpl, err = template.New("test.bicepparam").ParseFiles("poc/input/test.bicepparam")
	if err != nil {
		return err
	}
	variables, err = config.GetDeployenvVariables(context.Background(), "public", "dev")
	if err != nil {
		return err
	}

	err = tmpl.ExecuteTemplate(output, "test.bicepparam", variables)
	if err != nil {
		return err
	}

	return nil
}

func NewConfigProvider(config, region, user string) *configProviderImpl {
	return &configProviderImpl{
		config: config,
		region: region,
		user:   user,
	}
}

func (cp *configProviderImpl) GetAllClouds() ([]string, error) {
	_, err := cp.loadConfig()
	cloudInfo := []string{}
	if err == nil {
		for cloud, _ := range cp.baseVariableOverrides.Overrides {
			cloudInfo = append(cloudInfo, cloud)
		}
	}
	return cloudInfo, err
}

func (cp *configProviderImpl) GetCloudVariables(ctx context.Context, cloud string) (Variables, error) {
	variableOverrides, err := cp.loadConfig()
	variables := Variables{}
	if err == nil {
		for k, v := range variableOverrides.Defaults {
			variables[k] = v
		}
		if cloudOverride, ok := variableOverrides.Overrides[cloud]; ok {
			for k, v := range cloudOverride.Defaults {
				variables[k] = v
			}
		}
	}

	if err == nil {
		variables, err = NewDefaultVarHandler(cp.region, cp.user).ReplaceVariables(variables)
	}

	return variables, nil
}

func (cp *configProviderImpl) GetAllDeployEnvs() ([]*DeployEnvInfo, error) {
	_, err := cp.loadConfig()
	deployEnvInfo := []*DeployEnvInfo{}
	if err == nil {
		for cloud, cloudVariableOverrides := range cp.baseVariableOverrides.Overrides {
			for deployEnv := range cloudVariableOverrides.Overrides {
				deployEnvInfo = append(deployEnvInfo, &DeployEnvInfo{
					Cloud:     cloud,
					DeployEnv: deployEnv,
				})
			}
		}
	}

	return deployEnvInfo, err
}

func (cp *configProviderImpl) GetDeployenvVariables(ctx context.Context, cloud, deployenv string) (Variables, error) {
	variableOverrides, err := cp.loadConfig()
	variables := Variables{}
	if err == nil {
		for k, v := range variableOverrides.Defaults {
			variables[k] = v
		}
		if cloudOverride, ok := variableOverrides.Overrides[cloud]; ok {
			for k, v := range cloudOverride.Defaults {
				variables[k] = v
			}

			if deployEnvOverride, ok := cloudOverride.Overrides[deployenv]; ok {
				for k, v := range deployEnvOverride.Defaults {
					variables[k] = v
				}
			}
		}
	}

	if err == nil {
		variables, err = NewDefaultVarHandler(cp.region, cp.user).ReplaceVariables(variables)
	}

	return variables, nil
}

func (cp *configProviderImpl) GetAllRegions() ([]*RegionInfo, error) {
	_, err := cp.loadConfig()
	regionInfo := []*RegionInfo{}
	if err == nil {
		for cloud, cloudVariableOverrides := range cp.baseVariableOverrides.Overrides {
			for deployEnv, deployEnvVariableOverride := range cloudVariableOverrides.Overrides {
				for region, _ := range deployEnvVariableOverride.Overrides {
					regionInfo = append(regionInfo, &RegionInfo{
						Cloud:     cloud,
						DeployEnv: deployEnv,
						Region:    region,
					})
				}
			}
		}
	}

	return regionInfo, err
}

func (cp *configProviderImpl) GetVariables(ctx context.Context, cloud, deployEnv, region string) (Variables, error) {
	variableOverrides, err := cp.loadConfig()
	variables := Variables{}
	if err == nil {
		for k, v := range variableOverrides.Defaults {
			variables[k] = v
		}
		if cloudOverride, ok := variableOverrides.Overrides[cloud]; ok {
			for k, v := range cloudOverride.Defaults {
				variables[k] = v
			}

			if deployEnvOverride, ok := cloudOverride.Overrides[deployEnv]; ok {
				for k, v := range deployEnvOverride.Defaults {
					variables[k] = v
				}

				if regionOverride, ok := deployEnvOverride.Overrides[region]; ok {
					for k, v := range regionOverride {
						variables[k] = v
					}
				}
			}
		}
	}
	if err == nil {
		variables, err = NewDefaultVarHandler(cp.region, cp.user).ReplaceVariables(variables)
	}
	return variables, err
}

func (cp *configProviderImpl) loadConfig() (*VariableOverrides, error) {
	builtInConfigContent, err := os.ReadFile(cp.config)
	if err == nil {
		err = cp.updateVariables(builtInConfigContent)
	}

	return cp.baseVariableOverrides, err
}

func (cp *configProviderImpl) updateVariables(configContent []byte) error {
	currentVariableOverrides := &VariableOverrides{}
	err := yaml.Unmarshal(configContent, currentVariableOverrides)

	if err == nil {
		cp.baseVariableOverrides = currentVariableOverrides
	}
	return err
}
