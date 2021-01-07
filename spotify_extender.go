package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fgoyer/spotify-extender/spotext"
	"github.com/zmb3/spotify"
)

const redirectURI = "http://localhost:8080/callback"

var (
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopePlaylistModifyPrivate, spotify.ScopePlaylistModifyPublic)
	ch    = make(chan *spotify.Client)
	state = "extender"

	chillhopLofi2020 = spotext.GenreSearch{PlaylistID: "6SIX5vPxNfxkZW0rehrHP4", Query: "year:2020 AND genre:\"chillhop\" AND genre:\"lo-fi beats\""}
	chillhopLofi2019 = spotext.GenreSearch{PlaylistID: "1e8Bk00Ah6mrX40giHTpKK", Query: "year:2019 AND genre:\"chillhop\" AND genre:\"lo-fi beats\""}
	dwTest           = spotext.GenreSearch{PlaylistID: "7FcV5iqwkqiYSfK3FxNMOe", Query: ""}
	duplicates       = spotify.ID("7zs8x17jLEWdowTXsPKojx")
)

func main() {
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	url := auth.AuthURL(state)
	log.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// wait for auth to complete
	client := <-ch
	client.AutoRetry = true

	// use the client to make calls that require authorization
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("You are logged in as:", user.DisplayName)

	err = spotext.RemoveDuplicates(chillhopLofi2020.PlaylistID, client)
	// err = spotext.Compile(chillhopLofi2020, client)
	if err != nil {
		log.Fatal(err)
	}
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	token, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}

	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	client := auth.NewClient(token)
	fmt.Fprintf(w, "Login Completed!")
	ch <- &client
}
