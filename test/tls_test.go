package test

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/TrogonStack/TrogonEventStore-Client-Go/trogoneventstore"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTLS(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		config, err := trogoneventstore.ParseConnectionString(fmt.Sprintf("esdb://admin:changeit@%s", "localhost:2111,localhost:2112,localhost:2113"))
		if err != nil {
			t.Fatalf("Unexpected configuration error: %s", err.Error())
		}

		c, err := trogoneventstore.NewClient(config)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err.Error())
		}
		defer c.Close()

		numberOfEventsToRead := 1
		numberOfEvents := uint64(numberOfEventsToRead)
		opts := trogoneventstore.ReadAllOptions{
			From:           trogoneventstore.Start{},
			Direction:      trogoneventstore.Backwards,
			ResolveLinkTos: true,
		}

		_, err = c.ReadAll(context.Background(), opts, numberOfEvents)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to verify certificate")
	})

	t.Run("DefaultsWithCertificate", func(t *testing.T) {
		config, err := trogoneventstore.ParseConnectionString(fmt.Sprintf("esdb://admin:changeit@%s", "localhost:2111,localhost:2112,localhost:2113"))
		if err != nil {
			t.Fatalf("Unexpected configuration error: %s", err.Error())
		}

		b, err := os.ReadFile("../certs/ca/ca.crt")
		if err != nil {
			t.Fatalf("failed to read node certificate ../certs/ca/ca.crt: %s", err.Error())
		}
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(b) {
			t.Fatalf("failed to append node certificates: %s", err.Error())
		}
		config.RootCAs = cp

		c, err := trogoneventstore.NewClient(config)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err.Error())
		}

		defer c.Close()

		numberOfEventsToRead := 1
		numberOfEvents := uint64(numberOfEventsToRead)
		opts := trogoneventstore.ReadAllOptions{
			From:           trogoneventstore.Start{},
			Direction:      trogoneventstore.Backwards,
			ResolveLinkTos: true,
		}
		stream, err := c.ReadAll(context.Background(), opts, numberOfEvents)
		require.NoError(t, err)
		defer stream.Close()
		evt, err := stream.Recv()
		require.Nil(t, evt)
		require.True(t, errors.Is(err, io.EOF))
	})

	t.Run("WithoutCertificateAndVerify", func(t *testing.T) {
		config, err := trogoneventstore.ParseConnectionString(fmt.Sprintf("esdb://admin:changeit@%s?tls=true&tlsverifycert=false", "localhost:2111,localhost:2112,localhost:2113"))
		if err != nil {
			t.Fatalf("Unexpected configuration error: %s", err.Error())
		}

		c, err := trogoneventstore.NewClient(config)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err.Error())
		}

		defer c.Close()

		numberOfEventsToRead := 1
		numberOfEvents := uint64(numberOfEventsToRead)
		opts := trogoneventstore.ReadAllOptions{
			From:           trogoneventstore.Start{},
			Direction:      trogoneventstore.Backwards,
			ResolveLinkTos: true,
		}
		stream, err := c.ReadAll(context.Background(), opts, numberOfEvents)
		require.NoError(t, err)
		defer stream.Close()

		evt, err := stream.Recv()
		require.Nil(t, evt)
		require.True(t, errors.Is(err, io.EOF))
	})

	t.Run("testTLSWithoutCertificate", func(t *testing.T) {
		config, err := trogoneventstore.ParseConnectionString(fmt.Sprintf("esdb://admin:changeit@%s?tls=true&tlsverifycert=true", "localhost:2111,localhost:2112,localhost:2113"))
		if err != nil {
			t.Fatalf("Unexpected configuration error: %s", err.Error())
		}

		c, err := trogoneventstore.NewClient(config)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err.Error())
		}

		defer c.Close()

		numberOfEventsToRead := 1
		numberOfEvents := uint64(numberOfEventsToRead)
		opts := trogoneventstore.ReadAllOptions{
			From:           trogoneventstore.Start{},
			Direction:      trogoneventstore.Backwards,
			ResolveLinkTos: true,
		}
		_, err = c.ReadAll(context.Background(), opts, numberOfEvents)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to verify certificate")
	})

	t.Run("WithCertificate", func(t *testing.T) {
		config, err := trogoneventstore.ParseConnectionString(fmt.Sprintf("esdb://admin:changeit@%s?tls=true&tlsverifycert=true", "localhost:2111,localhost:2112,localhost:2113"))
		if err != nil {
			t.Fatalf("Unexpected configuration error: %s", err.Error())
		}

		b, err := os.ReadFile("../certs/ca/ca.crt")
		if err != nil {
			t.Fatalf("failed to read node certificate ../certs/ca/ca.crt: %s", err.Error())
		}
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(b) {
			t.Fatalf("failed to append node certificates: %s", err.Error())
		}
		config.RootCAs = cp

		c, err := trogoneventstore.NewClient(config)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err.Error())
		}

		defer c.Close()

		numberOfEventsToRead := 1
		numberOfEvents := uint64(numberOfEventsToRead)
		opts := trogoneventstore.ReadAllOptions{
			From:           trogoneventstore.Start{},
			Direction:      trogoneventstore.Backwards,
			ResolveLinkTos: true,
		}
		stream, err := c.ReadAll(context.Background(), opts, numberOfEvents)
		require.NoError(t, err)
		defer stream.Close()
		evt, err := stream.Recv()
		require.Nil(t, evt)
		require.True(t, errors.Is(err, io.EOF))
	})

	t.Run("WithCertificateFromAbsoluteFile", func(t *testing.T) {
		absPath, err := filepath.Abs("../certs/ca/ca.crt")
		if err != nil {
			t.Fatalf("Unexpected error: %s", err.Error())
		}

		s := fmt.Sprintf("esdb://admin:changeit@%s?tls=true&tlsverifycert=true&tlsCAFile=%s", "localhost:2111,localhost:2112,localhost:2113", absPath)
		config, err := trogoneventstore.ParseConnectionString(s)
		if err != nil {
			t.Fatalf("Unexpected configuration error: %s", err.Error())
		}

		c, err := trogoneventstore.NewClient(config)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err.Error())
		}

		defer c.Close()

		numberOfEventsToRead := 1
		numberOfEvents := uint64(numberOfEventsToRead)
		opts := trogoneventstore.ReadAllOptions{
			From:           trogoneventstore.Start{},
			Direction:      trogoneventstore.Backwards,
			ResolveLinkTos: true,
		}
		stream, err := c.ReadAll(context.Background(), opts, numberOfEvents)
		require.NoError(t, err)
		defer stream.Close()

		evt, err := stream.Recv()
		require.Nil(t, evt)
		require.True(t, errors.Is(err, io.EOF))
	})

	t.Run("WithCertificateFromRelativeFile", func(t *testing.T) {
		config, err := trogoneventstore.ParseConnectionString(fmt.Sprintf("esdb://admin:changeit@%s?tls=true&tlsverifycert=true&tlsCAFile=../certs/ca/ca.crt", "localhost:2111,localhost:2112,localhost:2113"))
		if err != nil {
			t.Fatalf("Unexpected configuration error: %s", err.Error())
		}

		c, err := trogoneventstore.NewClient(config)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err.Error())
		}

		defer c.Close()

		WaitForAdminToBeAvailable(t, c)
		numberOfEventsToRead := 1
		numberOfEvents := uint64(numberOfEventsToRead)
		opts := trogoneventstore.ReadAllOptions{
			From:           trogoneventstore.Start{},
			Direction:      trogoneventstore.Backwards,
			ResolveLinkTos: true,
		}
		stream, err := c.ReadAll(context.Background(), opts, numberOfEvents)
		require.NoError(t, err)
		defer stream.Close()
		evt, err := stream.Recv()
		require.Nil(t, evt)
		require.True(t, errors.Is(err, io.EOF))
	})

	t.Run("WithInvalidCertificate", func(t *testing.T) {
		config, err := trogoneventstore.ParseConnectionString(fmt.Sprintf("esdb://admin:changeit@%s?tls=true&tlsverifycert=true", "localhost:2111,localhost:2112,localhost:2113"))
		if err != nil {
			t.Fatalf("Unexpected configuration error: %s", err.Error())
		}

		b, err := os.ReadFile("../certs/untrusted-ca/ca.crt")
		if err != nil {
			t.Fatalf("failed to read node certificate ../certs/untrusted-ca/ca.crt: %s", err.Error())
		}
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(b) {
			t.Fatalf("failed to append node certificates: %s", err.Error())
		}
		config.RootCAs = cp

		c, err := trogoneventstore.NewClient(config)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err.Error())
		}

		defer c.Close()

		numberOfEventsToRead := 1
		numberOfEvents := uint64(numberOfEventsToRead)
		opts := trogoneventstore.ReadAllOptions{
			From:           trogoneventstore.Start{},
			Direction:      trogoneventstore.Backwards,
			ResolveLinkTos: true,
		}
		_, err = c.ReadAll(context.Background(), opts, numberOfEvents)
		esdbErr, ok := trogoneventstore.FromError(err)
		require.False(t, ok)
		require.NotNil(t, esdbErr)
		assert.Contains(t, esdbErr.Error(), "certificate signed by unknown authority")
	})
}

func WaitForAdminToBeAvailable(t *testing.T, db *trogoneventstore.Client) {
	for count := 0; count < 50; count++ {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		t.Logf("[debug] checking if admin user is available...%v/50", count)

		stream, err := db.ReadStream(ctx, "$users", trogoneventstore.ReadStreamOptions{}, 1)

		if ctx.Err() != nil {
			t.Log("[debug] request timed out, retrying...")
			cancel()
			time.Sleep(1 * time.Second)
			continue
		}

		if stream != nil {
			_, err = stream.Recv()
			if err == nil {
				t.Log("[debug] admin is available!")
				cancel()
				stream.Close()
				return
			}
		}

		if err != nil {
			if esdbError, ok := trogoneventstore.FromError(err); !ok {
				if esdbError.Code() == trogoneventstore.ErrorCodeResourceNotFound ||
					esdbError.Code() == trogoneventstore.ErrorCodeUnauthenticated ||
					esdbError.Code() == trogoneventstore.ErrorCodeDeadlineExceeded ||
					esdbError.Code() == trogoneventstore.ErrorUnavailable {
					time.Sleep(1 * time.Second)
					t.Logf("[debug] not available retrying...")
					cancel()
					continue
				}

				t.Fatalf("unexpected error when waiting the admin account to be available: %+v", esdbError)
			}

			t.Fatalf("unexpected error when waiting the admin account to be available: %+v", err)
		}
	}

	t.Fatalf("failed to access admin account in a timely manner")
}
