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
	Attributes struct {
		Sudo struct {
			Enabled         types.Bool `tfsdk:"enabled"`
			WithoutPassword types.Bool `tfsdk:"withoutPassword"`
		} `tfsdk:"sudo"`
		LdapGroups []struct {
			Name types.String `tfsdk:"name"`
		} `tfsdk:"ldapGroups"`
		PosixGroups []struct {
			ID   types.Int64  `tfsdk:"id"`
			Name types.String `tfsdk:"name"`
		} `tfsdk:"posixGroups"`
		Radius struct {
			Reply []struct {
				Name  types.String `tfsdk:"name"`
				Value types.String `tfsdk:"value"`
			} `tfsdk:"reply"`
		} `tfsdk:"radius"`
		SambaEnabled types.Bool `tfsdk:"sambaEnabled"`
	} `tfsdk:"attributes"`
	Description types.String `tfsdk:"description"`
	Email       types.String `tfsdk:"email"`
	ID          types.String `tfsdk:"id"`
	MemberQuery struct {
		QueryType types.String `tfsdk:"queryType"`
		Filters   []struct {
			Field    types.String `tfsdk:"field"`
			Operator types.String `tfsdk:"operator"`
			Value    types.String `tfsdk:"value"`
		} `tfsdk:"filters"`
	} `tfsdk:"memberQuery"`
	MemberQueryExemptions []struct {
		Attributes struct {
		} `tfsdk:"attributes"`
		ID   types.String `tfsdk:"id"`
		Type types.String `tfsdk:"type"`
	} `tfsdk:"memberQueryExemptions"`
	MemberSuggestionsNotify types.Bool   `tfsdk:"memberSuggestionsNotify"`
	MembershipMethod        types.String `tfsdk:"membership_method"`
	Name                    types.String `tfsdk:"name"`
	SuggestionCounts        struct {
		Add    int `tfsdk:"add"`
		Remove int `tfsdk:"remove"`
		Total  int `tfsdk:"total"`
	} `tfsdk:"suggestionCounts"`
	Type types.String `tfsdk:"type"`
}

// jcUserGroupsDataSourceModel maps the data source schema data.
type jcUserGroupsDataSourceModel struct {
	UserGroups []jcUserGroupsModel `tfsdk:"usergroups"`
}

// NewjcUserGroupDataSource is a helper function to simplify the provider implementation.
func NewjcUserGroupDataSource() datasource.DataSource {
	return &jcUserGroupDataSource{}
}

// jcUserGroupDataSource is the data source implementation.
// This struct accepts a client pointer to the JumpCloud API client so terraform can make its changes to the system
type jcUserGroupDataSource struct {
	client *jcclient.Client
}

// Metadata returns the data source type name.
func (d *jcUserGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_usergroup"
}

// Schema defines the schema for the data source.
func (d *jcUserGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"usergroups": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"membership_method": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *jcUserGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state jcUserGroupsDataSourceModel
	groups, err := d.client.GetAllUserGroups()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Jumpcloud User Groups",
			err.Error(),
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Hey there, I'm a log message"))
	tflog.Info(ctx, fmt.Sprintf("Read Jumpcloud User Groups: %v", len(groups)))
	// Map response body to model
	//for _, group := range groups {
	//	jcUserGroupState := jcUserGroupsModel{
	//		ID:          types.Int64Value(int64(coffee.ID)),
	//		Name:        types.StringValue(coffee.Name),
	//		Teaser:      types.StringValue(coffee.Teaser),
	//		Description: types.StringValue(coffee.Description),
	//		Price:       types.Float64Value(coffee.Price),
	//		Image:       types.StringValue(coffee.Image),
	//	}
	// Map response to state
	for _, group := range groups {
		jcUserGroupState := jcUserGroupsModel{
			ID:               types.StringValue(group.ID),
			Name:             types.StringValue(group.Name),
			Description:      types.StringValue(group.Description),
			MembershipMethod: types.StringValue(group.MembershipMethod),
			Type:             types.StringValue(group.Type),
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
