package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	jcclient "github.com/Spotnana-Tech/sec-jumpcloud-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &jcUserGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &jcUserGroupDataSource{}
)

// jcUserGroupsModel maps the provider schema data to a Go type.
type jcUserGroupsModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
}

// jcUserGroupsDataSourceModel maps the data source schema data.
// TODO: Update this struct value to types.ListType or types.SetType
type jcUserGroupsDataSourceModel struct {
	UserGroups []jcUserGroupsModel `tfsdk:"usergroups"`
}

// NewjcUserGroupDataSource is a helper function to simplify the provider implementation.
func NewjcUserGroupDataSource() datasource.DataSource {
	return &jcUserGroupDataSource{}
}

// jcUserGroupDataSource is the data source implementation.
// This struct accepts a client pointer to the JumpCloud Go client so terraform can make its changes to the system
type jcUserGroupDataSource struct {
	client *jcclient.Client
}

// Metadata returns the data source type name.
func (d *jcUserGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_usergroups"
}

// Schema defines the schema for the data source.
func (d *jcUserGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"usergroups": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "A list of Jumpcloud User Groups",
				MarkdownDescription: "A list of Jumpcloud User Groups",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "The ID of the User Group",
							MarkdownDescription: "The ID of the User Group",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							Description:         "The Name of the User Group",
							MarkdownDescription: "The Name of the User Group",
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed:            true,
							Description:         "Types can be user_group or system_group",
							MarkdownDescription: "Types can be user_group or system_group",
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *jcUserGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get all user groups
	var state jcUserGroupsDataSourceModel
	groups, err := d.client.GetAllUserGroups()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Jumpcloud User Groups",
			err.Error(),
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Read Jumpcloud User Groups: %v", len(groups)))

	// Map response to state
	for _, group := range groups {
		jcUserGroupState := jcUserGroupsModel{
			ID:          types.StringValue(group.ID),
			Name:        types.StringValue(group.Name),
			Description: types.StringValue(group.Description),
			Type:        types.StringValue(group.Type),
		}
		state.UserGroups = append(state.UserGroups, jcUserGroupState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return

	}
}

// Configure adds the provider configured client to the data source.
func (d *jcUserGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	// This is where we import our client for this type of data source
	client, ok := req.ProviderData.(*jcclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}
