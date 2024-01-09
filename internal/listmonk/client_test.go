package listmonk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	client := &Client{
		Host: "http://localhost:9000/",
	}

	t.Run("GetTemplates", func(t *testing.T) {
		_, err := client.GetTemplates()
		assert.NoError(t, err)
	})

	t.Run("GetTemplate", func(t *testing.T) {
		_, err := client.GetTemplate(1)
		assert.NoError(t, err)
	})

	var templateID int
	t.Run("CreateTemplate", func(t *testing.T) {
		template := Template{
			Name:    "Test Template",
			Body:    "<html><body><h1>Test Template</h1></body></html>",
			Subject: "Test Template",
			Type:    "tx",
		}
		t1, err := client.CreateTemplate(&template)
		assert.NoError(t, err)
		templateID = t1.ID
	})

	t.Run("UpdateTemplate", func(t *testing.T) {
		template := Template{
			ID:      templateID,
			Name:    "Test Template Updated",
			Body:    "<html><body><h1>Test Template</h1></body></html>",
			Subject: "Test Template Upd",
			Type:    "tx",
		}
		_, err := client.UpdateTemplate(&template)
		assert.NoError(t, err)
	})

	t.Run("DeleteTemplate", func(t *testing.T) {
		err := client.DeleteTemplate(templateID)
		assert.NoError(t, err)
	})
}
