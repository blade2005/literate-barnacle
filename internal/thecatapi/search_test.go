package thecatapi

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendRequest(t *testing.T) {
	c := NewClient(os.Getenv("MOCK_URL_API_KEY"))
	res, err := c.GetImagesSearch(nil)
	assert.NoError(t, err, "expecting nil error")
	assert.NotNil(t, res, "expecting non-nil result")
}
