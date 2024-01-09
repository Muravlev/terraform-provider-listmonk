// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"terraform-provider-listmonk/internal/listmonk"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ListmonkProvider satisfies various provider interfaces.
var _ provider.Provider = &ListmonkProvider{}

// ListmonkProvider defines the provider implementation.
type ListmonkProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ListmonkProviderModel describes the provider data model.
type ListmonkProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Headers  types.Map    `tfsdk:"headers"`
}

func (p *ListmonkProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "listmonk"
	resp.Version = p.version
}

func (p *ListmonkProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Required:    true,
				Description: "URL of the listmonk instance. Example: `https://listmonk.example.com`",
			},
			"username": schema.StringAttribute{
				Optional:    true,
				Description: "Username of the listmonk instance. Example: `username`",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Password of the listmonk instance. Example: `password`",
			},
			"headers": schema.MapAttribute{
				Optional:    true,
				Description: "Headers to be sent with each request. Example: `{ \"X-Listmonk-Header\": \"value\" }`",
				ElementType: types.StringType,
				Sensitive:   true,
			},
		},
	}
}

func (p *ListmonkProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config ListmonkProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown listmonk host",
			"Please set the listmonk host",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}
	headers := map[string]string{}
	config.Headers.ElementsAs(ctx, &headers, false)
	for k, v := range headers {
		headers[k] = v
	}
	// Example client configuration for data sources and resources
	client := listmonk.NewClient(
		config.Host.ValueString(),
		config.Username.ValueString(),
		config.Password.ValueString(),
		headers,
	)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *ListmonkProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewTemplateResource,
	}
}

func (p *ListmonkProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewTemplateDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ListmonkProvider{
			version: version,
		}
	}
}
