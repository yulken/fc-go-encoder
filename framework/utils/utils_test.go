package utils_test

import (
	"encoder/framework/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsJson(t *testing.T) {
	json := `{
		"id": 1,
		"status": "Active"
	}`

	err := utils.IsJson(json)
	require.Nil(t, err)

	invalidJson := `yikes`

	err = utils.IsJson(invalidJson)
	require.Error(t, err)
}
