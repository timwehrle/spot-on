package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/timwehrle/spot-on/internal/keyring"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"log"
	"net/http"
	"sync"
)

const redirectURI = "http://localhost:8080/callback"

var (
	auth  = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopePlaylistReadPrivate))
	state = "abc123"
	ch    = make(chan *spotify.Client)
	mu    sync.Mutex
)

// Manager handles Spotify authentication and token retrieval
type Manager struct{}

// NewAuthManager creates an instance of AuthManager
func NewAuthManager() *Manager {
	return &Manager{}
}

// StartAuthServer runs the local OAuth callback server
func (m *Manager) StartAuthServer() {
	http.HandleFunc("/callback", m.completeAuth)
	go func() {
		log.Println("Starting Auth Server on :8080")
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

// GetAuthURL returns the authentication URL
func (m *Manager) GetAuthURL() string {
	return auth.AuthURL(state)
}

// completeAuth handles the OAuth callback
func (m *Manager) completeAuth(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}

	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// Store token
	err = keyring.Set(tok)
	if err != nil {
		log.Fatalf("%v", err)
	}

	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}

// GetClient retrieves a new Spotify client with stored token
func (m *Manager) GetClient(ctx context.Context) (*spotify.Client, error) {
	token, err := keyring.Get()
	if err != nil {
		return nil, errors.New("no token found, please authenticate")
	}

	return spotify.New(auth.Client(ctx, token)), nil
}
