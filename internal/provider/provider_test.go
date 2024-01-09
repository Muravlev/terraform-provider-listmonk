package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var (
	providerConfig = fmt.Sprintf(`
		provider "listmonk" {
			host = "http://localhost:%s"
			username = "listmonk"
			password = "listmonk"
		}
`, dockerClient.ContainerPort)
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"listmonk": providerserver.NewProtocol6WithError(New("test")()),
	}
)
