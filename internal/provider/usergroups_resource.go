package provider

import (
	"context"
	"fmt"
	jumpcloud "github.com/Spotnana-Tech/sec-jumpcloud-client-go"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"slices"
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
	client *jumpcloud.Client
}

// UserGroupResourceModel is the local model for this resource type.
type UserGroupResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	Type             types.String `tfsdk:"type"`
	Email            types.String `tfsdk:"email"`
	MembershipMethod types.String `tfsdk:"membership_method"`
	Members          types.Set    `tfsdk:"members"`
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
				Computed:            true,
				Description:         "User Group ID",
				MarkdownDescription: "User Group ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "User Group Name",
				MarkdownDescription: "User Group Name",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "User Group Description",
				MarkdownDescription: "User Group Description",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Computed:            true,
				Description:         "ex. user_group or device_group type",
				MarkdownDescription: "ex. user_group or device_group type",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				Computed:            true,
				Description:         "User group email address",
				MarkdownDescription: "User group email address",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"membership_method": schema.StringAttribute{
				Computed:            true,
				Description:         "Can be STATIC or DYNAMIC_AUTOMATED or DYNAMIC_REVIEW_REQUIRED",
				MarkdownDescription: "Can be STATIC or DYNAMIC_AUTOMATED or DYNAMIC_REVIEW_REQUIRED",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"members": schema.SetAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "User emails associated with this group",
				MarkdownDescription: "This is a set of user emails associated with this group.",
				ElementType:         types.StringType,
			},
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

	// Get the members emails from the plan
	var planMemberEmails []string
	memberUserEmailSet, _ := plan.Members.ToSetValue(ctx)
	diags = memberUserEmailSet.ElementsAs(ctx, &planMemberEmails, false) //nolint:all

	// Get the user ids from the emails
	var memberUserIds []string
	for _, member := range planMemberEmails {
		userId, _ := r.client.GetUserIDFromEmail(member)
		memberUserIds = append(memberUserIds, userId)
	}

	// Cast local model to client model
	group := jumpcloud.UserGroup{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	// Create new group, check for errors
	g, err := r.client.CreateUserGroup(group)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating group",
			"Could not create group, unexpected error: "+err.Error(),
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Created Jumpcloud User Group: %s", g.Name))

	// Add members
	for _, userID := range memberUserIds {
		ok, _ := r.client.AddUserToGroup(g.ID, userID)
		if !ok {
			resp.Diagnostics.AddError(
				"Error adding user to group",
				"Could not add user to group, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Get the newly created group
	newGroup, _ := r.client.GetUserGroup(g.ID)

	// Get the members
	var memberEmails []attr.Value // This is the terraform structure requirement
	members, _ := r.client.GetGroupMembers(newGroup.ID)
	for _, member := range members {
		email, _ := r.client.GetUserEmailFromID(member.To.ID)
		memberEmails = append(memberEmails, types.StringValue(email))
	}
	returnedMembers, _ := types.SetValue(types.StringType, memberEmails)
	// Map response body to schema and populate Computed attribute values
	plan = UserGroupResourceModel{
		ID:               types.StringValue(newGroup.ID),
		Description:      types.StringValue(newGroup.Description),
		Name:             types.StringValue(newGroup.Name),
		Email:            types.StringValue(newGroup.Email),
		Type:             types.StringValue(newGroup.Type),
		MembershipMethod: types.StringValue(newGroup.MembershipMethod),
		Members:          returnedMembers,
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
	// Retrieve values from state
	var state UserGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed group value from the jumpcloud client
	tflog.Info(ctx, fmt.Sprintf("Looking Up Group ID: %s", state.ID.ValueString()))
	group, err := r.client.GetUserGroup(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Jumpcloud Group",
			"Could not read Jumpcloud Group ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
	// Get the members
	var memberEmails []attr.Value // This is the terraform structure requirement
	members, _ := r.client.GetGroupMembers(state.ID.ValueString())
	for _, member := range members {
		email, _ := r.client.GetUserEmailFromID(member.To.ID)
		memberEmails = append(memberEmails, types.StringValue(email))
	}
	returnedMembers, _ := types.SetValue(types.StringType, memberEmails)

	// Overwrite items with refreshed state
	state = UserGroupResourceModel{
		ID:               types.StringValue(group.ID),
		Name:             types.StringValue(group.Name),
		Description:      types.StringValue(group.Description),
		Type:             types.StringValue(group.Type),
		Email:            types.StringValue(group.Email),
		MembershipMethod: types.StringValue(group.MembershipMethod),
		Members:          returnedMembers,
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
	diags := req.Plan.Get(ctx, &plan)  //nolint:all
	diags = req.State.Get(ctx, &state) // existing resource state as defined in the terraform state file
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Isolate the members from the plan and state
	stateMembers, _ := state.Members.ToSetValue(ctx)
	plannedMembers, _ := plan.Members.ToSetValue(ctx)

	// Turn them in to []string
	var newMembers []string
	plannedMembers.ElementsAs(ctx, &newMembers, false)
	var oldMembers []string
	stateMembers.ElementsAs(ctx, &oldMembers, false)

	// Cast local model to client model
	groupModification := jumpcloud.UserGroup{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	// Update group, reference the state's group Id
	_, err := r.client.UpdateUserGroup(state.ID.ValueString(), groupModification)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Modifying Group",
			"Could not modify Group ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Get current group membership
	var currentMemberEmails []string
	currentMembers, _ := r.client.GetGroupMembers(state.ID.ValueString())
	for _, member := range currentMembers {
		email, _ := r.client.GetUserEmailFromID(member.To.ID)
		currentMemberEmails = append(currentMemberEmails, email)
	}

	// TODO if plan member not in current members, add to group
	for _, member := range newMembers {
		if !slices.Contains(currentMemberEmails, member) {
			uid, _ := r.client.GetUserIDFromEmail(member)
			ok, _ := r.client.AddUserToGroup(state.ID.ValueString(), uid)
			if !ok {
				resp.Diagnostics.AddError(
					"Error Adding User to Group",
					"Could not add user to group, unexpected error: "+err.Error(),
				)
				return
			}
		}
	}

	// TODO if current member not in plan members, remove from group
	for _, member := range oldMembers {
		if !slices.Contains(newMembers, member) {
			// Remove member from group
			uid, _ := r.client.GetUserIDFromEmail(member)
			ok, _ := r.client.RemoveUserFromGroup(state.ID.ValueString(), uid)
			if !ok {
				resp.Diagnostics.AddError(
					"Error Removing User from Group",
					"Could not remove user from group, unexpected error: "+err.Error(),
				)
				return
			}
		}
	}

	// Get the updated group
	groupState, err := r.client.GetUserGroup(state.ID.ValueString()) //nolint:all
	tflog.Info(ctx, fmt.Sprintf("Group Name: %s Group ID: %s", groupState.Name, groupState.ID))

	// Get the members
	var updatedMemberEmails []attr.Value // This is the terraform structure requirement
	updatedMembers, _ := r.client.GetGroupMembers(state.ID.ValueString())
	// Iterate through the members and get their emails
	for _, member := range updatedMembers {
		email, _ := r.client.GetUserEmailFromID(member.To.ID)
		updatedMemberEmails = append(updatedMemberEmails, types.StringValue(email))
	}
	finalMembers, _ := types.SetValue(types.StringType, updatedMemberEmails)
	// Map response body to schema and populate Computed attribute values
	plan = UserGroupResourceModel{
		ID:               types.StringValue(groupState.ID),
		Name:             types.StringValue(groupState.Name),
		Description:      types.StringValue(groupState.Description),
		Type:             types.StringValue(groupState.Type),
		Email:            types.StringValue(groupState.Email),
		MembershipMethod: types.StringValue(groupState.MembershipMethod),
		Members:          finalMembers,
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *jcUserGroupsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state UserGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing group. This object will be purged from the state file so there is no need to return values
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

	// This is where we import our client for this type of resource
	client, ok := req.ProviderData.(*jumpcloud.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *jumpcloud.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// ImportState imports the resource state from live resources via their ID attribute.
func (r *jcUserGroupsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
