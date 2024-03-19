package provider

import (
	"context"
	"fmt"
	"github.com/Spotnana-Tech/sec-jumpcloud-client-go"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &jcUserGroupDataLookupSource{}
	_ datasource.DataSourceWithConfigure = &jcUserGroupDataLookupSource{}
)

// jcUserGroupsLookupModel maps the provider schema data to a Go type.
type jcUserGroupsLookupModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	Members     types.Set    `tfsdk:"members"`
}

// jcUserGroupsDataSourceModel maps the data source schema data.
type jcUserGroupsLookupDataSourceModel struct {
	UserGroups []jcUserGroupsLookupModel `tfsdk:"groups"`
	Name       types.String              `tfsdk:"name"`
	Limit      types.Int64               `tfsdk:"limit"`
}

// NewjcUserGroupLookupDataSource is a helper function to simplify the provider implementation.
func NewjcUserGroupLookupDataSource() datasource.DataSource {
	return &jcUserGroupDataLookupSource{}
}

// jcUserGroupDataLookupSource is the data source implementation.
// This struct accepts a client pointer to the JumpCloud Go client so terraform can make its changes to the system.
type jcUserGroupDataLookupSource struct {
	client *jumpcloud.Client
}

// Metadata returns the data source type name.
func (d *jcUserGroupDataLookupSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_lookup"
}

// Schema defines the schema for the data source.
func (d *jcUserGroupDataLookupSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"groups": schema.ListNestedAttribute{
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
						"members": schema.SetAttribute{
							Computed:            true,
							Description:         "The members of the User Group",
							MarkdownDescription: "The members of the User Group",
							ElementType:         types.StringType,
						},
					},
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				Computed:            false,
				Description:         "The name to filter on",
				MarkdownDescription: "The name to filter on",
			},
			"limit": schema.Int64Attribute{
				Optional:            true,
				Computed:            false,
				Description:         "The limit of results to return",
				MarkdownDescription: "The limit of results to return",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *jcUserGroupDataLookupSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get all user groups
	var state jcUserGroupsLookupDataSourceModel
	var name string
	var limit int
	//diags := req.State.Get(ctx, &state)
	diags := req.Config.Get(ctx, &state)
	diags = req.Config.GetAttribute(ctx, path.Root("name"), &name)
	diags = req.Config.GetAttribute(ctx, path.Root("limit"), &limit)
	tflog.Info(ctx, fmt.Sprintf("Request: name %s, limit %d", name, limit))
	if limit == 0 {
		limit = 5
	}
	tflog.Info(ctx, fmt.Sprintf("Filters: %s : %v",
		name,
		limit,
	),
	)

	groups, err := d.client.SearchUserGroups(
		"name",
		name,
		limit,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Jumpcloud User Groups",
			err.Error(),
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Read Jumpcloud User Groups: %v", len(groups)))
	if len(groups) == 0 {
		resp.Diagnostics.AddError(
			"Unable to Read Jumpcloud User Groups",
			"No user groups found",
		)
		return
	}
	// Map response to state
	for _, group := range groups {
		// Get the members
		var memberEmails []attr.Value // This is the terraform structure requirement
		members, _ := d.client.GetGroupMembers(group.ID)
		for _, member := range members {
			email, _ := d.client.GetUserEmailFromID(member.To.ID)
			memberEmails = append(memberEmails, types.StringValue(email))
		}
		returnedMembers, _ := types.SetValue(types.StringType, memberEmails)
		jcUserGroupsLookupState := jcUserGroupsLookupModel{
			ID:          types.StringValue(group.ID),
			Name:        types.StringValue(group.Name),
			Description: types.StringValue(group.Description),
			Type:        types.StringValue(group.Type),
			Members:     returnedMembers,
		}
		state.UserGroups = append(state.UserGroups, jcUserGroupsLookupState)
	}
	state.Limit = types.Int64Value(int64(limit))
	state.Name = types.StringValue(name)
	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return

	}
}

// Configure adds the provider configured client to the data source.
func (d *jcUserGroupDataLookupSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
