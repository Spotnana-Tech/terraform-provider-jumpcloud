package provider

import (
	"context"
	"fmt"
	jcclient "github.com/Spotnana-Tech/sec-jumpcloud-client-go"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &jcUserGroupsResource{}
	_ resource.ResourceWithConfigure = &jcUserGroupsResource{}
)

// NewUserGroupsResource is a helper function to simplify the provider implementation.
func NewUserGroupsResource() resource.Resource {
	return &jcUserGroupsResource{}
}

// jcUserGroupsResource is the resource implementation.
type jcUserGroupsResource struct {
	client *jcclient.Client
}

// orderResourceModel maps the resource schema data.
type userGroupsResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Items       []UserGroup  `tfsdk:"user_groups"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

type UserGroups []UserGroup

type UserGroup struct {
	//Attributes struct {
	//	Sudo struct {
	//		Enabled         bool `tfsdk:"enabled"`
	//		WithoutPassword bool `tfsdk:"withoutPassword"`
	//	} `tfsdk:"sudo"`
	//	LdapGroups []struct {
	//		Name string `tfsdk:"name"`
	//	} `tfsdk:"ldap_groups"`
	//	PosixGroups []struct {
	//		ID   int    `tfsdk:"id"`
	//		Name string `tfsdk:"name"`
	//	} `tfsdk:"posix_groups"`
	//	Radius struct {
	//		Reply []struct {
	//			Name  string `tfsdk:"name"`
	//			Value string `tfsdk:"value"`
	//		} `tfsdk:"reply"`
	//	} `tfsdk:"radius"`
	//	SambaEnabled bool `tfsdk:"samba_enabled"`
	//} `tfsdk:"attributes"`
	Description string `tfsdk:"description"`
	Email       string `tfsdk:"email"`
	ID          string `tfsdk:"id"`
	MemberQuery struct {
		QueryType string `tfsdk:"query_type"`
		Filters   []struct {
			Field    string `tfsdk:"field"`
			Operator string `tfsdk:"operator"`
			Value    string `tfsdk:"value"`
		} `tfsdk:"filters"`
	} `tfsdk:"member_query"`
	//MemberQueryExemptions []struct {
	//	Attributes struct {
	//	} `tfsdk:"attributes"`
	//	ID   string `tfsdk:"id"`
	//	Type string `tfsdk:"type"`
	//} `tfsdk:"member_query_exemptions"`
	MemberSuggestionsNotify bool   `tfsdk:"member_suggestions_notify"`
	MembershipMethod        string `tfsdk:"membership_method"`
	Name                    string `tfsdk:"name"`
	//SuggestionCounts        struct {
	//	Add    int `tfsdk:"add"`
	//	Remove int `tfsdk:"remove"`
	//	Total  int `tfsdk:"total"`
	//} `tfsdk:"suggestion_counts"`
	Type string `tfsdk:"type"`
}
type CreateUserGroup struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	Email       types.String `tfsdk:"email"`
	//MemberQuery types.Map    `tfsdk:"member_query"`
}

type NewUserGroup struct {
	//Attributes struct {
	//	Sudo struct {
	//		Enabled         bool `tfsdk:"enabled"`
	//		WithoutPassword bool `tfsdk:"without_password"`
	//	} `tfsdk:"sudo"`
	//	LdapGroups []struct {
	//		Name string `tfsdk:"name"`
	//	} `tfsdk:"ldap_groups"`
	//	PosixGroups []struct {
	//		ID   int    `tfsdk:"id"`
	//		Name string `tfsdk:"name"`
	//	} `tfsdk:"posix_groups"`
	//	Radius struct {
	//		Reply []struct {
	//			Name  string `tfsdk:"name"`
	//			Value string `tfsdk:"value"`
	//		} `tfsdk:"reply"`
	//	} `tfsdk:"radius"`
	//	SambaEnabled bool `tfsdk:"samba_enabled"`
	//} `tfsdk:"attributes"`
	Description string `tfsdk:"description"`
	Email       string `tfsdk:"email"`
	ID          string `tfsdk:"id"`
	MemberQuery struct {
		QueryType string `tfsdk:"query_type"`
		Filters   []struct {
			Field    string `tfsdk:"field"`
			Operator string `tfsdk:"operator"`
			Value    string `tfsdk:"value"`
		} `tfsdk:"filters"`
	} `tfsdk:"member_query"`
	//MemberQueryExemptions []struct {
	//	Attributes struct {
	//	} `tfsdk:"attributes"`
	//	ID   string `tfsdk:"id"`
	//	Type string `tfsdk:"type"`
	//} `tfsdk:"member_query_exemptions"`
	MemberSuggestionsNotify bool   `tfsdk:"member_suggestions_notify"`
	MembershipMethod        string `tfsdk:"membership_method"`
	Name                    string `tfsdk:"name"`
	//SuggestionCounts        struct {
	//	Add    int `tfsdk:"add"`
	//	Remove int `tfsdk:"remove"`
	//	Total  int `tfsdk:"total"`
	//} `tfsdk:"suggestion_counts"`
	Type string `tfsdk:"type"`
}

// Metadata returns the resource type name.
func (r *jcUserGroupsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_usergroup"
}

// Schema defines the schema for the resource.
func (r *jcUserGroupsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Required: true,
			},
			"type": schema.StringAttribute{
				Computed: true,
			},
			"email": schema.StringAttribute{
				Computed: true,
			},
			//"membership_method": schema.StringAttribute{
			//	Computed: true,
			//	Optional: true,
			//},
			//"member_suggestions_notify": schema.BoolAttribute{
			//	Computed: true,
			//},
			// This is the greatest nested schema!
			//"member_query": schema.MapNestedAttribute{
			//	Computed: true,
			//	NestedObject: schema.NestedAttributeObject{
			//		Attributes: map[string]schema.Attribute{
			//			"query_type": schema.StringAttribute{
			//				Computed: true,
			//			},
			//			"filters": schema.ListNestedAttribute{
			//				Computed: true,
			//				Optional: true,
			//				NestedObject: schema.NestedAttributeObject{
			//					Attributes: map[string]schema.Attribute{
			//						"field": schema.StringAttribute{
			//							Computed: true,
			//						},
			//						"operator": schema.StringAttribute{
			//							Computed: true,
			//						},
			//						"value": schema.StringAttribute{
			//							Computed: true,
			//						},
			//					},
			//				},
			//			},
			//		},
			//	},
			//},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *jcUserGroupsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan CreateUserGroup
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("PLAN: %v", plan))
	// Cast local model to client model
	group := jcclient.UserGroup{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Type:        plan.Type.ValueString(),
		Email:       plan.Email.ValueString(),
	}
	//for _, item := range plan {
	//	groups = append(groups, jcclient.UserGroup(item))
	//}

	// Create new order
	newGroup, err := r.client.CreateUserGroup(group)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating group",
			"Could not create group, unexpected error: "+err.Error(),
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Created Jumpcloud User Group: %s", newGroup.Name))

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(newGroup.ID)
	plan = CreateUserGroup{
		Description: types.StringValue(newGroup.Description),
		Name:        types.StringValue(newGroup.Name),
		Email:       types.StringValue(newGroup.Email),
		Type:        types.StringValue(newGroup.Type),
	}
	tflog.Info(ctx, fmt.Sprintf("Created Jumpcloud User Group: %v", newGroup.Name))
	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *jcUserGroupsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CreateUserGroup
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Get refreshed group value from jcclient
	tflog.Info(ctx, fmt.Sprintf("Looking Up GroupId: %v", state.ID.ValueString()))
	group, err := r.client.GetUserGroup(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Jumpcloud Group",
			"Could not read Jumpcloud Group ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state = CreateUserGroup{
		Description: types.StringValue(group.Description),
		ID:          types.StringValue(group.ID),
		Name:        types.StringValue(group.Name),
		Email:       types.StringValue(group.Email),
		Type:        types.StringValue(group.Type),
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *jcUserGroupsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *jcUserGroupsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

// Configure adds the provider configured client to the resource.
func (r *jcUserGroupsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*jcclient.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *jcclient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}
