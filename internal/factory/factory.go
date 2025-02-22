package factory

import (
	"context"
	"errors"
	"github.com/timwehrle/spot-on/internal/auth"
	"github.com/zmb3/spotify/v2"
)

type Factory struct {
	authManager *auth.Manager
}

func New(authManager *auth.Manager) *Factory {
	return &Factory{
		authManager: authManager,
	}
}

func (f *Factory) Client() (*spotify.Client, error) {
	ctx := context.Background()
	client, err := f.authManager.GetClient(ctx)
	if err != nil {
		return nil, errors.New("could not create Spotify client, authentication required")
	}
	return client, nil
}
