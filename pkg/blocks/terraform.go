package blocks

import (
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type terraformBlock struct {
	requiredVersion   string
	requiredProviders []provider
	remoteBackend     remoteBackendConfig
	experiments       []string
}

type provider struct {
	name    string
	source  string
	version string
	cfg     map[string]string
}

type remoteBackendConfig struct {
	backendType   string
	backendConfig map[string]string
}

func NewTerraformBlock(requiredVersion string, experimental ...string) *terraformBlock {
	return &terraformBlock{requiredVersion: requiredVersion, experiments: experimental}
}

func (x *terraformBlock) AddProvider(name, source, version string, config map[string]string) *terraformBlock {
	x.requiredProviders = append(x.requiredProviders, provider{
		name:    name,
		source:  source,
		version: version,
		cfg:     config,
	})
	return x
}

func (x *terraformBlock) SetRemoteBackend(backendType string, backendConfig map[string]string) *terraformBlock {
	x.remoteBackend.backendType = backendType
	x.remoteBackend.backendConfig = backendConfig
	return x
}

func (x *terraformBlock) Type() string {
	return "terraform"
}

func (x *terraformBlock) String() string {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	tfBlock := rootBody.AppendNewBlock(x.Type(), nil)
	tfBlockBody := tfBlock.Body()
	tfBlockBody.SetAttributeValue("required_version", cty.StringVal(">= 1.9.0, < 2.0.0"))

	if len(x.requiredProviders) > 0 {
		tfBlockBody.AppendNewline()
		providerBlock := tfBlockBody.AppendNewBlock("required_providers", nil)
		providerBody := providerBlock.Body()
		for _, provider := range x.requiredProviders {
			providerBody.SetAttributeValue(provider.name, cty.MapVal(map[string]cty.Value{
				"version": cty.StringVal(provider.version),
				"source":  cty.StringVal(provider.source)}))

			rootBody.AppendNewline()
			providerCofigBody := rootBody.AppendNewBlock("provider", []string{provider.name}).Body()
			for k, v := range provider.cfg {
				providerCofigBody.SetAttributeValue(k, cty.StringVal(v))
			}
		}
	}

	if x.remoteBackend.backendType != "" {
		tfBlockBody.AppendNewline()
		backendBlock := tfBlockBody.AppendNewBlock("backend", []string{x.remoteBackend.backendType})
		backendBody := backendBlock.Body()
		for k, v := range x.remoteBackend.backendConfig {
			backendBody.SetAttributeValue(k, cty.StringVal(v))
		}
	}

	if len(x.experiments) > 0 {
		tfBlockBody.AppendNewline()
		vals := make([]cty.Value, 0, len(x.experiments))
		for _, e := range x.experiments {
			vals = append(vals, cty.StringVal(e))
		}
		tfBlockBody.SetAttributeValue("experiments", cty.ListVal(vals))
	}

	return string(f.Bytes())
}
