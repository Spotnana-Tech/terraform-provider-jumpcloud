package provider

import (
	"context"
	"fmt"
	jcclient "github.com/Spotnana-Tech/sec-jumpcloud-client-go"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &jcUserGroupsResource{}
	_ resource.ResourceWithConfigure   = &jcUserGroupsResource{}
	_ resource.ResourceWithImportState = &jcUserGroupsResource{}
)

// NewUserGroupsResource is a helper function to simplify the provider implementation.
func NewUserGroupsResource() resource.Resource {
	return &jcUserGroupsResource{}
}

// jcUserGroupsResource is the resource implementation.
type jcUserGroupsResource struct {
	client *jcclient.Client
}

// UserGroupResourceModel is the local model for this resource type.
type UserGroupResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	Email       types.String `tfsdk:"email"`
	//MemberQuery types.Map    `tfsdk:"member_query"`
}

// Metadata returns the resource type name.
func (r *jcUserGroupsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// Resources that end with this string will be routed to this resource implementation.
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
				MarkdownDescription: "User Group Description",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("UserGroup Managed by Terraform Provider snjumpcloud"),
			},
			"type": schema.StringAttribute{
				Computed: true,
			},
			"email": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
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
	var plan UserGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Cast local model to client model
	group := jcclient.UserGroup{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Type:        plan.Type.ValueString(),
		Email:       plan.Email.ValueString(),
	}
	// Create new group
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
	plan = UserGroupResourceModel{
		ID:          types.StringValue(newGroup.ID),
		Description: types.StringValue(newGroup.Description),
		Name:        types.StringValue(newGroup.Name),
		Email:       types.StringValue(newGroup.Email),
		Type:        types.StringValue(newGroup.Type),
	}
	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *jcUserGroupsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Get refreshed group value from jcclient
	tflog.Info(ctx, fmt.Sprintf("Looking Up Group ID: %s", state.ID.ValueString()))
	group, err := r.client.GetUserGroup(state.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Jumpcloud Group",
			"Could not read Jumpcloud Group ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state = UserGroupResourceModel{
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
	// Retrieve values from plan and state
	var plan, state UserGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("STATE ID: %s", state.ID.ValueString()))

	// Cast local model to client model
	groupModification := jcclient.UserGroup{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	// Update group, reference the state's group Id
	group, err := r.client.UpdateUserGroup(state.ID.ValueString(), groupModification)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Modifying Group",
			"Could not modify Group ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Group Name: %s", group.Name))
	tflog.Info(ctx, fmt.Sprintf("Group Desc: %s", group.Description))
	tflog.Info(ctx, fmt.Sprintf("Group ID: %s", group.ID))

	plan = UserGroupResourceModel{
		Description: types.StringValue(group.Description),
		ID:          types.StringValue(group.ID),
		Name:        types.StringValue(group.Name),
		Email:       types.StringValue(state.Email.ValueString()),
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *jcUserGroupsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserGroupResourceModel
	diags := req.State.Get(ctx, &state)
	tflog.Info(ctx, fmt.Sprintf("STATE: %s", state.ID.ValueString()))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing group
	err := r.client.DeleteUserGroup(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting UserGroup",
			"Could not delete user group, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *jcUserGroupsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	// This is where we import our client for this type of resource!
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

func (r *jcUserGroupsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
