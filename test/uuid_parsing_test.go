package test

import (
	"testing"

	"github.com/TrogonStack/TrogonEventStore-Client-Go/trogoneventstore"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUUIDParsingSuite(t *testing.T) {
	t.Run("UUIDParsingTests", func(t *testing.T) {
		expected := uuid.New()
		most, least := trogoneventstore.UUIDAsInt64(expected)
		actual, err := trogoneventstore.ParseUUIDFromInt64(most, least)

		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
