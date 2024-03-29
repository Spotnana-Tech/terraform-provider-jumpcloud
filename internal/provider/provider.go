package provider

import (
	"context"
	"github.com/Spotnana-Tech/sec-jumpcloud-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &jumpcloudProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &jumpcloudProvider{
			version: version,
		}
	}
}

// jumpcloudProvider is the provider implementation.
type jumpcloudProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *jumpcloudProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "jumpcloud"
	resp.Version = p.version
}

// jumpcloudProviderModel maps provider schema data to a Go type.
type jumpcloudProviderModel struct {
	ApiKey types.String `tfsdk:"api_key"`
}

// Schema defines the provider-level schema for configuration data.
func (p *jumpcloudProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				Description:         "The JumpCloud API key. This is a sensitive value and should be stored in environment variables, never in code.",
				MarkdownDescription: "The JumpCloud API key. This is a sensitive value and should be stored in environment variables, never in code.",
			},
		},
	}
}

// Configure prepares a JumpCloud API client for data sources and resources.
func (p *jumpcloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Jumpcloud client")

	// Retrieve provider data from configuration
	var config jumpcloudProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the API key from the configuration
	var apiKey string
	if !config.ApiKey.IsNull() {
		apiKey = config.ApiKey.ValueString()
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing JumpCloud API Key",
			"The provider cannot create the JumpCloud API client as the required api_key has not been supplied. ",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing JumpCloud API Key",
			"The provider cannot create the JumpCloud API client as there is an unknown configuration value for the JumpCloud API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the JC_API_KEY environment variable."+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Set provider-level log fields
	ctx = tflog.SetField(ctx, "jumpcloud_host", jumpcloud.HostURL)
	ctx = tflog.SetField(ctx, "api_key", apiKey)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "api_key") // Mask the API key in the logs
	tflog.Debug(ctx, "Creating Jumpcloud client")

	// Create a new jumpcloudProvider client using the configuration values
	client, err := jumpcloud.NewClient(apiKey)

	// If the client is not created, or the host is not the expected value, return an error
	if err != nil || !strings.Contains(client.HostURL.String(), "console.jumpcloud.com") {
		resp.Diagnostics.AddError(
			"Unable to Create JumpCloud API Client",
			"An unexpected error occurred when creating the JumpCloud API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				// TODO implement better error handling in JC client
				"JumpCloud Client Error: "+err.Error(),
		)
		return
	}

	// Make the JumpCloud client available during DataSource and Resource
	//type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
	tflog.Info(ctx, "Configured Jumpcloud client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *jumpcloudProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewjcUserGroupDataSource,
		NewjcGroupLookupDataSource,
		NewjcAppsDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *jumpcloudProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewUserGroupsResource,
		NewAppResource,
	}
}
