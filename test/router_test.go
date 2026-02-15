package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test404(t *testing.T) {
	response := sendTestRequest("GET", "/not-existing", "")
	assert.Equal(t, 404, response.Code)
}
