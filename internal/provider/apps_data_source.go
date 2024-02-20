package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Spotnana-Tech/sec-jumpcloud-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &jcAppsDataSource{}
	_ datasource.DataSourceWithConfigure = &jcAppsDataSource{}
)

// jcAppsModel maps the provider schema data to a Go type.
type jcAppsModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	DisplayName  types.String `tfsdk:"display_name"`
	DisplayLabel types.String `tfsdk:"display_label"`
	Description  types.String `tfsdk:"description"`
	SsoType      types.String `tfsdk:"sso_type"`
	Url          types.String `tfsdk:"url"`
}

// jcAppsDataSourceModel maps the data source schema data.
// TODO: Update this struct value to types.ListType or types.SetType
type jcAppsDataSourceModel struct {
	Apps []jcAppsModel `tfsdk:"apps"`
}

// NewjcAppsDataSource is a helper function to simplify the provider implementation.
func NewjcAppsDataSource() datasource.DataSource {
	return &jcAppsDataSource{}
}

// jcAppsDataSource is the data source implementation.
// This struct accepts a client pointer to the JumpCloud Go client so terraform can make its changes to the system
type jcAppsDataSource struct {
	client *jumpcloud.Client
}

// Metadata returns the data source type name.
func (d *jcAppsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_apps"
}

// Schema defines the schema for the data source.
func (d *jcAppsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"apps": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "A list of Jumpcloud Applications",
				MarkdownDescription: "A list of Jumpcloud Applications",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "The unique identifier for the application",
							MarkdownDescription: "The unique identifier for the application",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							Description:         "The name of the application (may be the template or integration name if generated via console)",
							MarkdownDescription: "The name of the application, may be the template or integration name if generated via console",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							Description:         "The display name for the application",
							MarkdownDescription: "The display name for the application",
						},
						"display_label": schema.StringAttribute{
							Computed:            true,
							Description:         "The display label for the application",
							MarkdownDescription: "The display label for the application",
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"sso_type": schema.StringAttribute{
							Computed:            true,
							Description:         "The type of SSO for the application, some apps are merely bookmarks to a webpage",
							MarkdownDescription: "The type of SSO for the application, some apps are merely bookmarks to a webpage",
						},
						"url": schema.StringAttribute{
							Computed:            true,
							Description:         "The URL for the application",
							MarkdownDescription: "The URL for the application",
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *jcAppsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get all user groups
	var state jcAppsDataSourceModel
	apps, err := d.client.GetAllApplications()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Jumpcloud User Groups",
			err.Error(),
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Read Jumpcloud User Groups: %v", len(apps)))

	// Map response to state
	for _, app := range apps {
		appstate := jcAppsModel{
			ID:           types.StringValue(app.ID),
			Name:         types.StringValue(app.Name),
			DisplayName:  types.StringValue(app.DisplayName),
			DisplayLabel: types.StringValue(app.DisplayLabel),
			Description:  types.StringValue(app.Description),
			Url:          types.StringValue(app.Sso.URL),
			SsoType:      types.StringValue(app.Sso.Type),
		}
		state.Apps = append(state.Apps, appstate)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return

	}
}

// Configure adds the provider configured client to the data source.
func (d *jcAppsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	// This is where we import our client for this type of data source
	client, ok := req.ProviderData.(*jumpcloud.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}
