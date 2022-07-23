package echo

import (
	stdContext "context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Hello returns a greeting for the named person.
func Hello(name string) (string, error) {
	// If no name was given, return an error with a message.
	if name == "" {
		return name, errors.New("empty name")
	}
	// Create a message using a random format.
	// message := fmt.Sprintf(randomFormat(), name)
	message := fmt.Sprint(name)
	return message, nil
}

func TestHelloEmpty(t *testing.T) {
	msg, err := Hello("")
	if msg != "" || err == nil {
		t.Fatalf(`Hello("") = %q, %v, want "", error`, msg, err)
	}
}

func waitForServerStart(e *Echo, errChan <-chan error, isTLS bool) error {
	ctx, cancel := stdContext.WithTimeout(stdContext.Background(), 200*time.Millisecond)
	defer cancel()

	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			var addr net.Addr
			addr = e.ListenerAddr()
			fmt.Println(addr)
			if addr != nil && strings.Contains(addr.String(), ":") {
				return nil // was started
			}
		case err := <-errChan:
			if err == http.ErrServerClosed {
				return nil
			}
			return err
		}
	}
}

func TestEchoStart(t *testing.T) {
	e := New()
	errChan := make(chan error)

	go func() {
		err := e.Start()
		if err != nil {
			errChan <- err
		}
	}()

	err := waitForServerStart(e, errChan, false)
	assert.NoError(t, err)
}
