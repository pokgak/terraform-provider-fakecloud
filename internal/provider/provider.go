// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fakecloud "github.com/pokgak/fakecloud/sdk"
)

// Ensure the implementation satisfies various provider interfaces.
var _ provider.Provider = &FakecloudProvider{}

// FakecloudProvider defines the provider implementation.
type FakecloudProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ScaffoldingProviderModel describes the provider data model.
type FakecloudProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *FakecloudProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "fakecloud"
	resp.Version = p.version
}

func (p *FakecloudProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *FakecloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config FakecloudProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Fakecloud API Host",
			"The provider cannot create the Fakecloud API client as there is an unknown configuration value for the Fakecloud API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the FAKECLOUD_HOST environment variable.",
		)
	}

	// if config.Username.IsUnknown() {
	// 	resp.Diagnostics.AddAttributeError(
	// 		path.Root("username"),
	// 		"Unknown Fakecloud API Username",
	// 		"The provider cannot create the Fakecloud API client as there is an unknown configuration value for the Fakecloud API username. "+
	// 			"Either target apply the source of the value first, set the value statically in the configuration, or use the FAKECLOUD_USERNAME environment variable.",
	// 	)
	// }

	// if config.Password.IsUnknown() {
	// 	resp.Diagnostics.AddAttributeError(
	// 		path.Root("password"),
	// 		"Unknown Fakecloud API Password",
	// 		"The provider cannot create the Fakecloud API client as there is an unknown configuration value for the Fakecloud API password. "+
	// 			"Either target apply the source of the value first, set the value statically in the configuration, or use the FAKECLOUD_PASSWORD environment variable.",
	// 	)
	// }

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("FAKECLOUD_HOST")
	username := os.Getenv("FAKECLOUD_USERNAME")
	password := os.Getenv("FAKECLOUD_PASSWORD")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Fakecloud API Host",
			"The provider cannot create the Fakecloud API client as there is a missing or empty value for the Fakecloud API host. "+
				"Set the host value in the configuration or use the FAKECLOUD_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	// if username == "" {
	// 	resp.Diagnostics.AddAttributeError(
	// 		path.Root("username"),
	// 		"Missing Fakecloud API Username",
	// 		"The provider cannot create the Fakecloud API client as there is a missing or empty value for the Fakecloud API username. "+
	// 			"Set the username value in the configuration or use the FAKECLOUD_USERNAME environment variable. "+
	// 			"If either is already set, ensure the value is not empty.",
	// 	)
	// }

	// if password == "" {
	// 	resp.Diagnostics.AddAttributeError(
	// 		path.Root("password"),
	// 		"Missing Fakecloud API Password",
	// 		"The provider cannot create the Fakecloud API client as there is a missing or empty value for the Fakecloud API password. "+
	// 			"Set the password value in the configuration or use the FAKECLOUD_PASSWORD environment variable. "+
	// 			"If either is already set, ensure the value is not empty.",
	// 	)
	// }

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new Fakecloud client using the configuration values
	client, err := fakecloud.NewClient(host, username, password)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Fakecloud API Client",
			"An unexpected error occurred when creating the Fakecloud API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Fakecloud Client Error: "+err.Error(),
		)
		return
	}

	// Make the Fakecloud client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *FakecloudProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewVirtualMachineResource,
	}
}

func (p *FakecloudProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewVirtualMachinesDataSource,
		NewVirtualMachineDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &FakecloudProvider{
			version: version,
		}
	}
}
