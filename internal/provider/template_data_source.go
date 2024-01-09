// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"
	"terraform-provider-listmonk/internal/listmonk"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	// _ datasource.DataSource              = &TemplateDataSource{}
	_ datasource.DataSourceWithConfigure = &TemplateDataSource{}
)

func NewTemplateDataSource() datasource.DataSource {
	return &TemplateDataSource{}
}

// TemplateDataSource defines the data source implementation.
type TemplateDataSource struct {
	client *listmonk.Client
}

// TemplateDataSourceModel describes the data source data model.
type TemplateDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
	Name      types.String `tfsdk:"name"`
	Body      types.String `tfsdk:"body"`
	Type      types.String `tfsdk:"type"`
	IsDefault types.Bool   `tfsdk:"is_default"`
	Subject   types.String `tfsdk:"subject"`
}

func (d *TemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template"
}

func (d *TemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Template data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Template identifier",
				Required:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Template created at",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "Template updated at",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Template name",
				Computed:            true,
			},
			"body": schema.StringAttribute{
				MarkdownDescription: "Template body",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Template type",
				Computed:            true,
			},
			"is_default": schema.BoolAttribute{
				MarkdownDescription: "Template is default",
				Computed:            true,
			},
			"subject": schema.StringAttribute{
				MarkdownDescription: "Template subject",
				Computed:            true,
			},
		},
	}
}

func (d *TemplateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

	d.client = client
}

func (d *TemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TemplateDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the template id from the configuration.
	var templateIdStr string
	req.Config.GetAttribute(ctx, path.Root("id"), &templateIdStr)
	templateId, err := strconv.Atoi(templateIdStr)
	if err != nil {
		resp.Diagnostics.AddError("Import error", fmt.Sprintf("Unable to to parse template id to int, got error: %s", err))
		return
	}

	// Get the template from the client.
	template, err := d.client.GetTemplate(templateId)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Template, got error: %s", err))
		return
	}

	// Set the data source state from the client response.
	data.ID = types.StringValue(strconv.Itoa(template.ID))
	data.CreatedAt = types.StringValue(template.CreatedAt)
	data.UpdatedAt = types.StringValue(template.UpdatedAt)
	data.Name = types.StringValue(template.Name)
	data.Body = types.StringValue(template.Body)
	data.Type = types.StringValue(template.Type)
	data.IsDefault = types.BoolValue(template.IsDefault)
	data.Subject = types.StringValue(template.Subject)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
