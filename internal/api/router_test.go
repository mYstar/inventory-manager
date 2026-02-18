package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test404(t *testing.T) {
	router := initTestRouter()
	response := sendTestRequest(router, "GET", "/not-existing", "")
	assert.Equal(t, 404, response.Code)
}
