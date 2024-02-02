package provider

import (
	"context"
	"fmt"
	jcclient "github.com/Spotnana-Tech/sec-jumpcloud-client-go"
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
	_ resource.Resource                = &jcAppAssociationResource{}
	_ resource.ResourceWithConfigure   = &jcAppAssociationResource{}
	_ resource.ResourceWithImportState = &jcAppAssociationResource{}
)

// NewAppAssociationResource is a helper function to simplify the provider implementation.
func NewAppAssociationResource() resource.Resource {
	return &jcAppAssociationResource{}
}

// jcAppAssociationResource is the resource implementation.
type jcAppAssociationResource struct {
	client *jcclient.Client
}
type AppAssociationSchemaModel struct {
	ID               types.String `tfsdk:"app_id"`
	Name             types.String `tfsdk:"name"`
	DisplayName      types.String `tfsdk:"display_name"`
	DisplayLabel     types.String `tfsdk:"display_label"`
	AssociatedGroups types.Set    `tfsdk:"associated_groups"`
}
type Association struct {
	GroupID   types.String `tfsdk:"group_id"`
	GroupName types.String `tfsdk:"group_name"`
}

// AppAssociationResourceModel is the local model for this resource type.
// This may be wrong? Using API schema here but need to use TF provider schema above
type AppAssociationResourceModel []struct {
	ID                 types.String        `tfsdk:"app_id"`
	Type               types.String        `tfsdk:"type"`
	Name               types.String        `tfsdk:"name"`
	CompiledAttributes *CompiledAttributes `tfsdk:"compiledAttributes"`
	Paths              *PathAttributes     `tfsdk:"paths"`
}
type CompiledAttributes struct {
	LdapGroups *LdapGroups `tfsdk:"ldapGroups"`
}
type LdapGroups []struct {
	Name types.String `tfsdk:"name"`
}
type PathAttributes [][]struct {
	Attributes types.Map `tfsdk:"attributes"`
	To         *To       `tfsdk:"to"`
}
type To struct {
	Attributes *CompiledAttributes `tfsdk:"attributes"`
	ID         types.String        `tfsdk:"id"`
	Type       types.String        `tfsdk:"type"`
}

// Configure adds the provider configuration to the resource.
func (r *jcAppAssociationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Metadata returns the resource type name.
func (r *jcAppAssociationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_association"
}

// Schema defines the schema for the resource.
func (r *jcAppAssociationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"app_id": schema.StringAttribute{
				Optional: true,
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
				MarkdownDescription: "",
				DeprecationMessage:  "",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *jcAppAssociationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// We should not be Creating or Deleting Apps via TF provider... yet
}

// Read refreshes the Terraform state with the latest data.
func (r *jcAppAssociationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AppAssociationSchemaModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("State: %v", state))
	tflog.Info(ctx, fmt.Sprintf("Looking Up App ID: %s %s", state.ID.ValueString(), state.Name.ValueString()))
	app, err := r.client.GetApplication(state.ID.ValueString())
	tflog.Info(ctx, fmt.Sprintf("Looked Up App ID: %s %s", app.ID, app.DisplayName))
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
	// Iterate through our app associations
	for _, a := range associations {
		_id := a.To.ID // App ID
		idAssociations = append(idAssociations, types.StringValue(_id))
	}
	appAssociations, _ := types.SetValue(types.StringType, idAssociations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Overwrite items with refreshed state
	state = AppAssociationSchemaModel{
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
func (r *jcAppAssociationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state AppAssociationSchemaModel
	diags := req.State.Get(ctx, &state)
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	oldstate, _ := state.AssociatedGroups.ToSetValue(ctx)
	newstate, _ := plan.AssociatedGroups.ToSetValue(ctx)
	tflog.Info(ctx, fmt.Sprintf("OLD: %v", oldstate))
	tflog.Info(ctx, fmt.Sprintf("NEW: %v", newstate))

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
		tflog.Info(ctx, fmt.Sprintf("Checking %s now...\n", group))
		// if group is not in the old state, associate it
		if !slices.Contains(oldElements, group) {
			tflog.Info(ctx, fmt.Sprintf("Group %s is not currently associated with this app\n", group))
			// Associate the group with the app
			err = r.client.AssociateGroupWithApp(state.ID.ValueString(), group)
		}
	}
	// For each group in our old configuration
	for _, group := range oldElements {
		if !slices.Contains(newElements, group) {
			tflog.Info(ctx, fmt.Sprintf("Group %s is no longer associated with this app\n", group))
			// Disassociate the group with the app
			err = r.client.RemoveGroupFromApp(state.ID.ValueString(), group)
		}
	}

	associations, err := r.client.GetAppAssociations(state.ID.ValueString(), "user_group")
	tflog.Info(ctx, fmt.Sprintf("Associations: %s", associations))
	var idAssociations []attr.Value
	// Iterate through our app associations
	for _, a := range associations {
		idAssociations = append(idAssociations, types.StringValue(a.To.ID))
	}
	appAssociations, _ := types.SetValue(types.StringType, idAssociations)
	tflog.Info(ctx, fmt.Sprintf("CurrentAssociated Groups: %s\n\n", state.AssociatedGroups.Elements()))
	tflog.Info(ctx, fmt.Sprintf("Associations: %s\n\n", appAssociations.Elements()))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Overwrite items with refreshed state
	state = AppAssociationSchemaModel{
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
func (r *jcAppAssociationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// We should not be Creating or Deleting Apps via TF provider... yet
	return
}

func (r *jcAppAssociationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("app_id"), req, resp)
}
