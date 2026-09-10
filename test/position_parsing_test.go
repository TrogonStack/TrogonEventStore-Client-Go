package test

import (
	"testing"

	"github.com/TrogonStack/TrogonEventStore-Client-Go/trogoneventstore"
	"github.com/stretchr/testify/assert"
)

func TestPositionParsingSuite(t *testing.T) {
	t.Run("StreamPositionTests", func(t *testing.T) {
		pos, err := trogoneventstore.ParseStreamPosition("C:123/P:456")
		assert.NoError(t, err)
		assert.NotNil(t, pos)

		obj, err := trogoneventstore.ParseStreamPosition("C:-1/P:-1")
		assert.NoError(t, err)

		_, ok := obj.(trogoneventstore.End)
		assert.True(t, ok)

		obj, err = trogoneventstore.ParseStreamPosition("C:0/P:0")
		assert.NoError(t, err)

		_, ok = obj.(trogoneventstore.Start)
		assert.True(t, ok)

		obj, err = trogoneventstore.ParseStreamPosition("-1")
		assert.NoError(t, err)

		_, ok = obj.(trogoneventstore.End)
		assert.True(t, ok)

		obj, err = trogoneventstore.ParseStreamPosition("0")
		assert.NoError(t, err)

		_, ok = obj.(trogoneventstore.Start)
		assert.True(t, ok)

		obj, err = trogoneventstore.ParseStreamPosition("42")
		assert.NoError(t, err)

		value, ok := obj.(trogoneventstore.StreamRevision)
		assert.True(t, ok)
		assert.Equal(t, uint64(42), value.Value)
	})
}
