package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTemplateResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
				resource "listmonk_template" "test" {
					body = "<p>Hello world</p>"
					name = "tf-test"
					subject = "test1"
					type = "tx"
				}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("listmonk_template.test", "body", "<p>Hello world</p>"),
					resource.TestCheckResourceAttr("listmonk_template.test", "name", "tf-test"),
					resource.TestCheckResourceAttr("listmonk_template.test", "subject", "test1"),
					resource.TestCheckResourceAttr("listmonk_template.test", "type", "tx"),
				),
			},
			// ImportState testing
			// TODO: fix this
			// {
			// 	Config:            providerConfig,
			// 	ResourceName:      "listmonk_template.test",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	ImportStateId:     "1",
			// },
			// Update and Read testing
			{
				Config: providerConfig + `
			resource "listmonk_template" "test" {
				body = "<p>Hello there</p>"
				name = "tf-test"
				subject = "test1"
				type = "tx"
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("listmonk_template.test", "body", "<p>Hello there</p>"),
					resource.TestCheckResourceAttr("listmonk_template.test", "name", "tf-test"),
					resource.TestCheckResourceAttr("listmonk_template.test", "subject", "test1"),
					resource.TestCheckResourceAttr("listmonk_template.test", "type", "tx"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
