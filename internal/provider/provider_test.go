package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	providerConfig = `
		terraform {
			required_providers {
				listmonk = {
					source = "github.com/Muravlev/listmonk"
				}
			}
		}
		provider "listmonk" {
			host = "http://localhost:9000"
		}
`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"listmonk": providerserver.NewProtocol6WithError(New("test")()),
	}
)
