package playlists

import (
	"log"

	"github.com/zmb3/spotify"
)

// GenreSearch contains the query to search for and playlist ID for where to put the results.
type GenreSearch struct {
	PlaylistID spotify.ID
	Query      string
}

// Compile initiates the given search and places the results in the given playlist.
func Compile(s GenreSearch, client *spotify.Client) error {

	// search
	log.Println("Searching...")
	results, err := client.Search(s.Query, spotify.SearchTypeTrack)
	if err != nil {
		log.Fatal(err)
	}

	tracks := make([]spotify.ID, 0)
	// Traverse the pages of the search results. Spotify returns 20 results per page,
	// let's enforce a strict 5,000 track (250 page) limit per search.
	for page := 1; page <= 250; page++ {
		// handle track results
		if results.Tracks != nil {
			for _, item := range results.Tracks.Tracks {
				tracks = append(tracks, item.ID)
			}
		}

		err = client.NextTrackResults(results)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	// Retrieve the tracks currently on the playlist.
	playlistTracks, err := client.GetPlaylistTracks(s.PlaylistID)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Playlist has %d total tracks", playlistTracks.Total)
	currentTracks := make([]spotify.ID, 0)
	for page := 1; ; page++ {
		if playlistTracks.Tracks != nil {
			for _, track := range playlistTracks.Tracks {
				currentTracks = append(currentTracks, track.Track.ID)
			}
		}

		err = client.NextPage(playlistTracks)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	// Remove any tracks that are already in the playlist
	tracks = difference(tracks, currentTracks)

	// Push track results to a playlist.
	resultCount := len(tracks)
	log.Printf("Found %v new tracks\n", resultCount)
	if resultCount > 0 {
		var temp []spotify.ID
		for i := len(tracks) - 1; i >= 0; i-- {
			temp = append(temp, tracks[i])
			if len(temp) == 100 || i == 0 {
				snapshot, err := client.AddTracksToPlaylist(s.PlaylistID, temp...)
				if err != nil {
					log.Fatal(err)
				}
				log.Printf("Added %v tracks to playlist. Snapshot: %v\n", len(temp), snapshot)
				temp = nil
			}
		}
	}

	return nil
}
