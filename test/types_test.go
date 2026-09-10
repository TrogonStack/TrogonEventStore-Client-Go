package test

import (
	"testing"
	"time"

	"github.com/TrogonStack/TrogonEventStore-Client-Go/trogoneventstore"
	"github.com/stretchr/testify/assert"
)

func TestTypes(t *testing.T) {
	t.Run("TestConsistentMetadataSerializationStreamAcl", func(t *testing.T) {
		acl := trogoneventstore.Acl{}
		acl.AddReadRoles("admin")
		acl.AddWriteRoles("admin")
		acl.AddDeleteRoles("admin")
		acl.AddMetaReadRoles("admin")
		acl.AddMetaWriteRoles("admin")

		expected := trogoneventstore.StreamMetadata{}
		expected.SetMaxAge(2 * time.Second)
		expected.SetCacheControl(15 * time.Second)
		expected.SetTruncateBefore(1)
		expected.SetMaxCount(12)
		expected.SetAcl(acl)
		expected.AddCustomProperty("foo", "bar")

		bytes, err := expected.ToJson()
		assert.NoError(t, err, "failed to serialize in JSON")

		meta, err := trogoneventstore.StreamMetadataFromJson(bytes)
		assert.NoError(t, err, "failed to parse Metadata from props")
		assert.Equal(t, expected, *meta, "consistency serialization failure")
	})

	t.Run("TestConsistentMetadataSerializationUserStreamAcl", func(t *testing.T) {
		expected := trogoneventstore.StreamMetadata{}
		expected.SetMaxAge(2 * time.Second)
		expected.SetCacheControl(15 * time.Second)
		expected.SetTruncateBefore(1)
		expected.SetMaxCount(12)
		expected.SetAcl(trogoneventstore.UserStreamAcl)
		expected.AddCustomProperty("foo", "bar")

		bytes, err := expected.ToJson()
		assert.NoError(t, err, "failed to serialize in JSON")

		meta, err := trogoneventstore.StreamMetadataFromJson(bytes)
		assert.NoError(t, err, "failed to parse Metadata from props")
		assert.Equal(t, expected, *meta, "consistency serialization failure")
	})

	t.Run("TestConsistentMetadataSerializationSystemStreamAcl", func(t *testing.T) {
		expected := trogoneventstore.StreamMetadata{}
		expected.SetMaxAge(2 * time.Second)
		expected.SetCacheControl(15 * time.Second)
		expected.SetTruncateBefore(1)
		expected.SetMaxCount(12)
		expected.SetAcl(trogoneventstore.SystemStreamAcl)
		expected.AddCustomProperty("foo", "bar")

		bytes, err := expected.ToJson()
		assert.NoError(t, err, "failed to serialize in JSON")

		meta, err := trogoneventstore.StreamMetadataFromJson(bytes)
		assert.NoError(t, err, "failed to parse Metadata from props")
		assert.Equal(t, expected, *meta, "consistency serialization failure")
	})

	t.Run("TestCustomPropertyRetrievalFromStreamMetadata", func(t *testing.T) {
		expected := trogoneventstore.StreamMetadata{}
		expected.AddCustomProperty("foo", "bar")

		foo := expected.CustomProperty("foo")
		assert.Equal(t, "bar", foo, "custom property value mismatch")
	})

	t.Run("TestUnknownCustomPropertyRetrievalFromStreamMetadata", func(t *testing.T) {
		expected := trogoneventstore.StreamMetadata{}
		expected.AddCustomProperty("foo", "bar")

		foo := expected.CustomProperty("foes")
		assert.Empty(t, foo, "custom property value mismatch")
	})

	t.Run("TestGetAllCustomPropertiesFromStreamMetadata", func(t *testing.T) {
		expected := trogoneventstore.StreamMetadata{}
		expected.AddCustomProperty("foo", 123)
		expected.AddCustomProperty("foes", "baz")

		props := expected.CustomProperties()
		assert.Equal(t, map[string]interface{}{"foo": 123, "foes": "baz"}, props, "custom properties mismatch")
	})
}
