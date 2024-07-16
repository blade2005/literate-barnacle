package thecatapi

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSendRequest(t *testing.T) {
	c := NewClient(os.Getenv("MOCK_URL_API_KEY"))
	res, err := c.GetImagesSearch(nil)
	assert.Nil(t, err, "expecting nil error")
	assert.NotNil(t, res, "expecting non-nil result")
}
