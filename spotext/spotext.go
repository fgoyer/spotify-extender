package spotext

import (
	"github.com/zmb3/spotify"
)

// SpotifyExtender is a client for working with the Spotify Web API.
type SpotifyExtender struct {
	client  *spotify.Client
	baseURL string

	AutoRetry bool
}

// New creates a new instance of the SpotifyExtender
func (s *SpotifyExtender) New() {

}
