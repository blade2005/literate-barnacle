package thecatapi

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendRequest(t *testing.T) {
	c := NewClient(os.Getenv("MOCK_URL_API_KEY"))
	res, err := c.GetImagesSearch(nil)
	fmt.Printf("res res: %#v\n", res)
	assert.NoError(t, err, "expecting nil error")
	assert.NotNil(t, res, "expecting non-nil result")
}
