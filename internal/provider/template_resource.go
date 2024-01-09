package provider

import (
	"context"
	"fmt"
	"strconv"
	"terraform-provider-listmonk/internal/listmonk"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &templateResource{}
	_ resource.ResourceWithConfigure   = &templateResource{}
	_ resource.ResourceWithImportState = &templateResource{}
)

// NewtemplateResource is a helper function to simplify the provider implementation.
func NewTemplateResource() resource.Resource {
	return &templateResource{}
}

// Configure adds the provider configured client to the resource.
func (r *templateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*listmonk.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *listmonk.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// templateResource is the resource implementation.
type templateResource struct {
	client *listmonk.Client
}

// templateResourceModel describes the resource data model.
type templateResourceModel struct {
	ID        types.String `tfsdk:"id"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
	Name      types.String `tfsdk:"name"`
	Body      types.String `tfsdk:"body"`
	Type      types.String `tfsdk:"type"`
	IsDefault types.Bool   `tfsdk:"is_default"`
	Subject   types.String `tfsdk:"subject"`
}

// Metadata returns the resource type name.
func (t *templateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template"
}

// Schema defines the schema for the resource.
func (t *templateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Template data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Template identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Template created at",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "Template updated at",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Template name",
				Required:            true,
			},
			"body": schema.StringAttribute{
				MarkdownDescription: "Template body",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Template type",
				Required:            true,
			},
			"is_default": schema.BoolAttribute{
				MarkdownDescription: "Template is default",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"subject": schema.StringAttribute{
				MarkdownDescription: "Template subject",
				Required:            true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (t *templateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan templateResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	template := listmonk.Template{
		Body:    plan.Body.ValueString(),
		Name:    plan.Name.ValueString(),
		Subject: plan.Subject.ValueString(),
		Type:    plan.Type.ValueString(),
	}

	// Create the resource
	r, err := t.client.CreateTemplate(&template)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create template",
			fmt.Sprintf("Failed to create template: %s", err),
		)

		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(strconv.Itoa(r.ID))
	plan.CreatedAt = types.StringValue(r.CreatedAt)
	plan.UpdatedAt = types.StringValue(r.UpdatedAt)
	plan.IsDefault = types.BoolValue(r.IsDefault)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (t *templateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state templateResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx = tflog.SetField(ctx, "id", state.ID)
	tflog.Info(ctx, "Reading template")

	// Get refreshed template value from Listmonk
	templateId, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to parse template ID (reading)",
			fmt.Sprintf("Unable to parse template ID: %s", err),
		)
	}
	template, err := t.client.GetTemplate(templateId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read template",
			fmt.Sprintf("Failed to read template: %s", err),
		)

		return
	}

	// Overwrite items with refreshed state
	state.CreatedAt = types.StringValue(template.CreatedAt)
	state.UpdatedAt = types.StringValue(template.UpdatedAt)
	state.Name = types.StringValue(template.Name)
	state.Body = types.StringValue(template.Body)
	state.Type = types.StringValue(template.Type)
	state.IsDefault = types.BoolValue(template.IsDefault)
	state.Subject = types.StringValue(template.Subject)

	// Set refreshed state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (t *templateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan templateResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	templateId, err := strconv.ParseInt(plan.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to parse template ID (updating)",
			fmt.Sprintf("Unable to parse template ID: %s", err),
		)
		return
	}
	// Generate API request body from plan
	template := listmonk.Template{
		ID:      int(templateId),
		Body:    plan.Body.ValueString(),
		Name:    plan.Name.ValueString(),
		Subject: plan.Subject.ValueString(),
		Type:    plan.Type.ValueString(),
	}

	// Update existing template
	r, err := t.client.UpdateTemplate(&template)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update template",
			fmt.Sprintf("Failed to update template: %s", err),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	// plan.ID = types.Int64Value(int64(r.ID))
	// plan.CreatedAt = types.StringValue(r.CreatedAt)
	plan.UpdatedAt = types.StringValue(r.UpdatedAt)
	plan.IsDefault = types.BoolValue(r.IsDefault)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (t *templateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state templateResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	templateId, err := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to parse template ID (deleting)",
			fmt.Sprintf("Unable to parse template ID: %s", err),
		)
		return
	}
	// Delete existing template
	err = t.client.DeleteTemplate(int(templateId))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Template",
			"Could not delete template, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *templateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	id := path.Root("id")
	resource.ImportStatePassthroughID(ctx, id, req, resp)
}
