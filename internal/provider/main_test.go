package provider

import (
	"terraform-provider-listmonk/internal/tests"
	"testing"
)

var (
	dockerClient = tests.NewTestDocker()
)

func TestMain(m *testing.M) {
	defer dockerClient.Cleanup()
	m.Run()
}
