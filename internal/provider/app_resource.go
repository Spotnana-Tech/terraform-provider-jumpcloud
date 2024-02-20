package provider

import (
	"context"
	"fmt"
	"github.com/Spotnana-Tech/sec-jumpcloud-client-go"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"slices"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &jcAppResource{}
	_ resource.ResourceWithConfigure   = &jcAppResource{}
	_ resource.ResourceWithImportState = &jcAppResource{}
)

// NewAppResource is a helper function to simplify the provider implementation.
func NewAppResource() resource.Resource {
	return &jcAppResource{}
}

// jcAppResource is the resource implementation.
type jcAppResource struct {
	client *jumpcloud.Client
}

// AppSchemaModel is the local model for this resource type.
type AppSchemaModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	DisplayName      types.String `tfsdk:"display_name"`
	DisplayLabel     types.String `tfsdk:"display_label"`
	AssociatedGroups types.Set    `tfsdk:"associated_groups"`
}

// Metadata returns the resource type name.
func (r *jcAppResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app"
}

// Schema defines the schema for the resource.
func (r *jcAppResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"display_name": schema.StringAttribute{
				Computed: true,
			},
			"display_label": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			// Set attribute does not care about order
			"associated_groups": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Description:         "Group IDs associated with this app",
				MarkdownDescription: "This is a set of group IDs associated with this app.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *jcAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// We should not be Creating or Deleting Apps via TF provider... yet
}

// Read refreshes the Terraform state with the latest data.
func (r *jcAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get the current state
	var state AppSchemaModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the app by ID
	tflog.Info(ctx, fmt.Sprintf("Looking Up App ID: %s %s", state.ID.ValueString(), state.Name.ValueString()))
	app, err := r.client.GetApplication(state.ID.ValueString())
	tflog.Info(ctx, fmt.Sprintf("Look Up Results: %s %s", app.ID, app.DisplayName))
	// Get the app associations
	associations, err := r.client.GetAppAssociations(state.ID.ValueString(), "user_group")
	tflog.Info(ctx, fmt.Sprintf("Associations: %s", associations))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Jumpcloud Group",
			"Could not read Jumpcloud Group ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// A temp holder for associations
	var idAssociations []attr.Value
	// Iterate through our app associations, add them to idAssociations
	for _, a := range associations {
		_id := a.To.ID // App ID
		idAssociations = append(idAssociations, types.StringValue(_id))
	}
	appAssociations, diags := types.SetValue(types.StringType, idAssociations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Overwrite items with refreshed state
	state = AppSchemaModel{
		ID:               types.StringValue(app.ID),
		Name:             types.StringValue(app.Name),
		DisplayName:      types.StringValue(app.DisplayName),
		DisplayLabel:     types.StringValue(app.DisplayLabel),
		AssociatedGroups: appAssociations,
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *jcAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan and state
	var plan, state AppSchemaModel
	diags := req.State.Get(ctx, &state)
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Turning stated and planned associated groups into sets for use in comparison
	oldstate, _ := state.AssociatedGroups.ToSetValue(ctx)
	newstate, _ := plan.AssociatedGroups.ToSetValue(ctx)

	// Get the current app associations
	CurrentAssociations, err := r.client.GetAppAssociations(state.ID.ValueString(), "user_group")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Jumpcloud Group",
			"Could not read Jumpcloud Group ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Turn CurrentAssociations to a []string of group IDs
	var currentGroups []string
	for _, association := range CurrentAssociations {
		currentGroups = append(currentGroups, association.To.ID)
	}
	tflog.Info(ctx, fmt.Sprintf("Looking Up App ID: %s %s\n", state.ID.ValueString(), state.DisplayName.ValueString()))
	tflog.Info(ctx, fmt.Sprintf("Currently has %v Groups associated\n", len(currentGroups)))

	// Turn oldstate and newstate into []string of group IDs - these are the associations of the app to usergroups
	var oldElements []string
	var newElements []string
	diags = oldstate.ElementsAs(ctx, &oldElements, false)
	diags = newstate.ElementsAs(ctx, &newElements, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// For each group in our new configuration
	for _, group := range newElements {

		// if group is not in the old state, associate it
		if !slices.Contains(oldElements, group) {
			tflog.Info(ctx, fmt.Sprintf("ADDING GROUPID %s TO %s \n", group, state.DisplayLabel.ValueString()))

			// Associate the group with the app
			err = r.client.AssociateGroupWithApp(state.ID.ValueString(), group)
		}
	}

	// For each group in our old configuration
	for _, group := range oldElements {

		// if group is not in the new state, disassociate it
		if !slices.Contains(newElements, group) {
			tflog.Info(ctx, fmt.Sprintf("REMOVING GROUPID %s FROM %s \n", group, state.DisplayLabel.ValueString()))

			// Disassociate the group with the app
			err = r.client.RemoveGroupFromApp(state.ID.ValueString(), group)
		}
	}

	// Get the app associations
	associations, err := r.client.GetAppAssociations(state.ID.ValueString(), "user_group")

	// Temp holder for associations to be added to state
	var idAssociations []attr.Value
	// Iterate through our app associations, add their terraform approved types to idAssociations
	for _, a := range associations {
		idAssociations = append(idAssociations, types.StringValue(a.To.ID))
	}
	// Turn this slice in to a set for terraform
	appAssociations, _ := types.SetValue(types.StringType, idAssociations)
	tflog.Info(ctx, fmt.Sprintf("CurrentAssociated Groups: %s\n\n", state.AssociatedGroups.Elements()))
	tflog.Info(ctx, fmt.Sprintf("Associations: %s\n\n", appAssociations.Elements()))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Overwrite items with refreshed state
	state = AppSchemaModel{
		ID:               types.StringValue(state.ID.ValueString()),
		Name:             types.StringValue(state.Name.ValueString()),
		DisplayName:      types.StringValue(state.DisplayName.ValueString()),
		DisplayLabel:     types.StringValue(state.DisplayLabel.ValueString()),
		AssociatedGroups: appAssociations,
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *jcAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// We should not be Creating or Deleting Apps via TF provider... yet
	return
}

// Configure adds the provider configuration to the resource.
func (r *jcAppResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	// This is where we import our client for this type of resource!
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

// ImportState imports the resource state from an existing resource.
func (r *jcAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
