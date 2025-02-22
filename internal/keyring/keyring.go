package keyring

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
	"time"
)

var (
	service = "spot-on"
	user    = "spotify-token"
)

func Set(secret *oauth2.Token) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	errCh := make(chan error, 1)

	data, err := json.Marshal(secret)
	if err != nil {
		return err
	}

	go func() {
		errCh <- keyring.Set(service, user, string(data))
		close(errCh)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("failed to set secret: %w", err)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("timeout while trying to set secret in keyring")
	}
}

func Get() (*oauth2.Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resultCh := make(chan *oauth2.Token, 1)
	errCh := make(chan error, 1)

	var secret string
	var err error
	var token *oauth2.Token

	go func() {
		defer close(resultCh)
		defer close(errCh)
		secret, err = keyring.Get(service, user)
		err = json.Unmarshal([]byte(secret), &token)
		if err != nil {
			errCh <- err
		} else {
			resultCh <- token
		}
	}()

	select {
	case token := <-resultCh:
		return token, nil
	case err := <-errCh:
		return nil, fmt.Errorf("failed to get secret: %w", err)
	case <-ctx.Done():
		return nil, fmt.Errorf("timeout while trying to get secret in keyring")
	}
}

func Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		errCh <- keyring.Delete(service, user)
		close(errCh)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("failed to delete secret: %w", err)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("timeout while trying to delete secret in keyring")
	}
}

func MockInit() {
	keyring.MockInit()
}

func MockInitWithError(err error) {
	keyring.MockInitWithError(err)
}
