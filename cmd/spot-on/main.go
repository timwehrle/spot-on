package main

import (
	"context"
	"fmt"
	"github.com/timwehrle/spot-on/internal/auth"
	"github.com/timwehrle/spot-on/internal/factory"
	"log"
)

func main() {
	authManager := auth.NewAuthManager()

	authManager.StartAuthServer()

	url := authManager.GetAuthURL()
	fmt.Println("Please log in to Spotify by visiting the following page:", url)

	f := factory.New(authManager)

	fmt.Println("Waiting for authentication...")
	client, err := f.Client()
	if err != nil {
		log.Fatal(err)
	}

	playlists, err := client.CurrentUsersPlaylists(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for i, playlist := range playlists.Playlists {
		fmt.Printf("%d. %s\n", i+1, playlist.Name)
	}
}
