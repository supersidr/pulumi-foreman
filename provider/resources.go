// Copyright 2016-2024, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package foreman

import (
	"path"

	// Allow embedding bridge-metadata.json in the provider.
	_ "embed"

	foreman "github.com/masikrus/terraform-provider-foreman/foreman" // Import the upstream provider

	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge"
	tks "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge/tokens"
	shimv2 "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfshim/sdk-v2"
	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"

	"github.com/masikrus/pulumi-foreman/provider/pkg/version"
)

// all of the token components used below.
const (
	// This variable controls the default name of the package in the package
	// registries for nodejs and python:
	mainPkg = "foreman"
	// modules:
	mainMod = "index" // the foreman module
)

//go:embed cmd/pulumi-resource-foreman/bridge-metadata.json
var metadata []byte

// Provider returns additional overlaid schema and metadata associated with the provider.
func Provider() tfbridge.ProviderInfo {
	// Create a Pulumi provider mapping
	prov := tfbridge.ProviderInfo{
		// Instantiate the Terraform provider
		//
		// The [pulumi-terraform-bridge](https://github.com/pulumi/pulumi-terraform-bridge) supports 3
		// types of Terraform providers:
		//
		// 1. Providers written with the terraform-plugin-sdk/v1:
		//
		//    If the provider you are bridging is written with the terraform-plugin-sdk/v1, then you
		//    will need to adapt the boilerplate:
		//
		//    - Change the import "shimv2" to "shimv1" and change the associated import to
		//      "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfshim/sdk-v1".
		//
		//    You can then proceed as normal.
		//
		// 2. Providers written with terraform-plugin-sdk/v2:
		//
		//    This boilerplate is already geared towards providers written with the
		//    terraform-plugin-sdk/v2, since it is the most common provider framework used. No
		//    adaptions are needed.
		//
		// 3. Providers written with terraform-plugin-framework:
		//
		//    If the provider you are bridging is written with the terraform-plugin-framework, then
		//    you will need to adapt the boilerplate:
		//
		//    - Remove the `shimv2` import and add:
		//
		//      	pfbridge "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfbridge"
		//
		//    - Replace `shimv2.NewProvider` with `pfbridge.ShimProvider`.
		//
		//    - In provider/cmd/pulumi-tfgen-foreman/main.go, replace the
		//      "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfgen" import with
		//      "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfgen". Remove the `version.Version`
		//      argument to `tfgen.Main`.
		//
		//    - In provider/cmd/pulumi-resource-foreman/main.go, replace the
		//      "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge" import with
		//      "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfbridge". Replace the arguments to the
		//      `tfbridge.Main` so it looks like this:
		//
		//      	tfbridge.Main(context.Background(), "foreman", foreman.Provider(),
		//			tfbridge.ProviderMetadata{PulumiSchema: pulumiSchema})
		//
		//   Detailed instructions can be found at
		//   https://pulumi-developer-docs.readthedocs.io/projects/pulumi-terraform-bridge/en/latest/docs/guides/new-pf-provider.html
		//   After that, you can proceed as normal.
		//
		// This is where you give the bridge a handle to the upstream terraform provider. SDKv2
		// convention is to have a function at "github.com/masikrus/terraform-provider-foreman/provider".New
		// which takes a version and produces a factory function. The provider you are bridging may
		// not do that. You will need to find the function (generally called in upstream's main.go)
		// that produces a:
		//
		// - *"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema".Provider (for SDKv2)
		// - *"github.com/hashicorp/terraform-plugin-sdk/v1/helper/schema".Provider (for SDKv1)
		// - "github.com/hashicorp/terraform-plugin-framework/provider".Provider (for plugin-framework)
		//
		//nolint:lll
		P: shimv2.NewProvider(foreman.Provider()),

		Name:    "foreman",
		Version: version.Version,
		// DisplayName is a way to be able to change the casing of the provider name when being
		// displayed on the Pulumi registry
		DisplayName: "Foreman",
		// Change this to your personal name (or a company name) that you would like to be shown in
		// the Pulumi Registry if this package is published there.
		Publisher: "masikrus",
		// LogoURL is optional but useful to help identify your package in the Pulumi Registry
		// if this package is published there.
		//
		// You may host a logo on a domain you control or add an PNG logo (100x100) for your package
		// in your repository and use the raw content URL for that file as your logo URL.
		LogoURL: "",
		// PluginDownloadURL is an optional URL used to download the Provider
		// for use in Pulumi programs
		// e.g. https://github.com/org/pulumi-provider-name/releases/download/v${VERSION}/
		PluginDownloadURL: "",
		Description:       "A Pulumi package for creating and managing foreman cloud resources.",
		// category/cloud tag helps with categorizing the package in the Pulumi Registry.
		// For all available categories, see `Keywords` in
		// https://www.pulumi.com/docs/guides/pulumi-packages/schema/#package.
		Keywords:   []string{"foreman", "category/cloud"},
		License:    "Apache-2.0",
		Homepage:   "https://www.pulumi.com",
		Repository: "https://github.com/masikrus/pulumi-foreman",
		// The GitHub Org for the provider - defaults to `terraform-providers`. Note that this should
		// match the TF provider module's require directive, not any replace directives.
		GitHubOrg: "",
		DataSources: map[string]*tfbridge.DataSourceInfo{
			"foreman_architecture": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_computeprofile": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_computeresource": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_defaulttemplate": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_domain": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_environment": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_location": {
				Fields: map[string]*tfbridge.SchemaInfo{
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_global_parameter": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_hostgroup": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_httpproxy": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_image": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_jobtemplate": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_katello_content_credential": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_katello_content_view": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_katello_lifecycle_environment": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_katello_product": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_katello_repository": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_katello_sync_plan": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_media": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_model": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_operatingsystem": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_organization": {
				Fields: map[string]*tfbridge.SchemaInfo{
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_parameter": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_partitiontable": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_provisioningtemplate": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_puppetclass": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_setting": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_smartclassparameter": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_smartproxy": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_subnet": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_templateinput": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_templatekind": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_user": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
			"foreman_usergroup": {
				// Tok: ForemanDataSource(mainMod, "getEnvironment"),
				Fields: map[string]*tfbridge.SchemaInfo{
					// HACK: remove this field for now as it breaks dotnet codegen due to our current type naming strategy.
					// https://github.com/pulumi/pulumi-terraform-bridge/issues/1118
					"__meta_":  {Omit: true},
					"__meta__": {Omit: true},
				},
			},
		},
		MetadataInfo: tfbridge.NewProviderMetadata(metadata),
		Config: map[string]*tfbridge.SchemaInfo{
			"server_hostname": {
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"FOREMAN_SERVER_HOSTNAME"},
				},
			},
			"server_protocol": {
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"FOREMAN_PROTOCOL"},
				},
			},
			"client_username": {
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"FOREMAN_CLIENT_USERNAME"},
				},
			},
			"client_password": {
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"FOREMAN_CLIENT_PASSWORD"},
				},
			},
			"organization_id": {
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"FOREMAN_ORGANIZATION_ID"},
				},
			},
			"location_id": {
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"FOREMAN_LOCATION_ID"},
				},
			},
		},
		// If extra types are needed for configuration, they can be added here.
		ExtraTypes: map[string]schema.ComplexTypeSpec{},
		JavaScript: &tfbridge.JavaScriptInfo{
			// RespectSchemaVersion ensures the SDK is generated linking to the correct version of the provider.
			PackageName:          "pulumi-foreman",
			RespectSchemaVersion: true,
		},
		Python: &tfbridge.PythonInfo{
			// RespectSchemaVersion ensures the SDK is generated linking to the correct version of the provider.
			RespectSchemaVersion: true,
			// Enable modern PyProject support in the generated Python SDK.
			PyProject: struct{ Enabled bool }{true},
		},
		Golang: &tfbridge.GolangInfo{
			// Set where the SDK is going to be published to.
			ImportBasePath: path.Join(
				"github.com/masikrus/pulumi-foreman/sdk/",
				tfbridge.GetModuleMajorVersion(version.Version),
				"go",
				mainPkg,
			),
			// Opt in to all available code generation features.
			GenerateResourceContainerTypes: true,
			GenerateExtraInputTypes:        true,
			// RespectSchemaVersion ensures the SDK is generated linking to the correct version of the provider.
			RespectSchemaVersion: true,
		},
		CSharp: &tfbridge.CSharpInfo{
			// RespectSchemaVersion ensures the SDK is generated linking to the correct version of the provider.
			RespectSchemaVersion: true,
			// Use a wildcard import so NuGet will prefer the latest possible version.
			PackageReferences: map[string]string{
				"Pulumi": "3.*",
			},
		},
	}

	// MustComputeTokens maps all resources and datasources from the upstream provider into Pulumi.
	//
	// tokens.SingleModule puts every upstream item into your provider's main module.
	//
	// You shouldn't need to override anything, but if you do, use the [tfbridge.ProviderInfo.Resources]
	// and [tfbridge.ProviderInfo.DataSources].
	prov.MustComputeTokens(tks.SingleModule("foreman_", mainMod,
		tks.MakeStandard(mainPkg)))

	prov.MustApplyAutoAliases()
	prov.SetAutonaming(255, "-")

	return prov
}
