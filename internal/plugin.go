package internal

import (
	"fmt"

	"github.com/creasty/defaults"
	"github.com/mach-composer/mach-composer-plugin-helpers/helpers"
	"github.com/mach-composer/mach-composer-plugin-sdk/plugin"
	"github.com/mach-composer/mach-composer-plugin-sdk/schema"
	"github.com/mitchellh/mapstructure"
)

type Plugin struct {
	environment string
	provider    string
	siteConfigs map[string]SiteConfig
}

func NewCommercelayerPlugin() schema.MachComposerPlugin {
	state := &Plugin{
		provider:    "0.0.8",
		siteConfigs: map[string]SiteConfig{},
	}

	return plugin.NewPlugin(&schema.PluginSchema{
		Identifier: "commercelayer",

		Configure: state.Configure,
		IsEnabled: state.IsEnabled,

		GetValidationSchema: state.GetValidationSchema,

		// Config
		SetSiteConfig: state.SetSiteConfig,

		// Renders
		RenderTerraformProviders: state.TerraformRenderProviders,
		RenderTerraformResources: state.TerraformRenderResources,
		RenderTerraformComponent: state.RenderTerraformComponent,
	})
}

func (p *Plugin) Configure(environment string, provider string) error {
	p.environment = environment
	if provider != "" {
		p.provider = provider
	}
	return nil
}

func (p *Plugin) IsEnabled() bool {
	return true
}

func (p *Plugin) GetValidationSchema() (*schema.ValidationSchema, error) {
	result := getSchema()
	return result, nil
}

func (p *Plugin) SetSiteConfig(site string, data map[string]any) error {
	if len(data) == 0 {
		return nil
	}

	cfg := SiteConfig{}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}

	if err := defaults.Set(&cfg); err != nil {
		return err
	}
	p.siteConfigs[site] = cfg
	return nil
}

func (p *Plugin) TerraformRenderProviders(site string) (string, error) {
	result := fmt.Sprintf(`
		commercelayer = {
			source  = "incentro-dc/commercelayer"
			version = "%s"
		}`, helpers.VersionConstraint(p.provider))
	return result, nil
}

func (p *Plugin) TerraformRenderResources(site string) (string, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return "", nil
	}

	apiEndpoint, err := cfg.ApiEndpoint()
	if err != nil {
		return "", err
	}

	authEndpoint, err := cfg.AuthEndpoint()
	if err != nil {
		return "", err
	}

	templateContext := struct {
		Settings     *SiteConfig
		ApiEndpoint  string
		AuthEndpoint string
	}{
		Settings:     cfg,
		ApiEndpoint:  apiEndpoint,
		AuthEndpoint: authEndpoint,
	}

	template := `
		provider "commercelayer" {
			{{ renderProperty "client_id" .Settings.ClientID }}
			{{ renderProperty "client_secret" .Settings.ClientSecret }}
			{{ renderProperty "api_endpoint" .ApiEndpoint }}
			{{ renderProperty "auth_endpoint" .AuthEndpoint }}
		}
	`
	return helpers.RenderGoTemplate(template, templateContext)
}

func (p *Plugin) RenderTerraformComponent(site string, component string) (*schema.ComponentSchema, error) {
	result := &schema.ComponentSchema{
		Providers: []string{
			"commercelayer = commercelayer",
		},
	}
	return result, nil
}

func (p *Plugin) getSiteConfig(site string) *SiteConfig {
	cfg, ok := p.siteConfigs[site]
	if !ok {
		return nil
	}
	return &cfg
}
