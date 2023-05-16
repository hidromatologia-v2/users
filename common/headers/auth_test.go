package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthorization(t *testing.T) {
	assert.Greater(t, len(Authorization("INVALID")), 0)
}
