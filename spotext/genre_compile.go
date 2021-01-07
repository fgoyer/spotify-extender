package spotext

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
		return err
	}
	log.Println("Searching complete, compiling tracks...")
	tracks := make([]spotify.FullTrack, 0)
	// Traverse the pages of the search results. Spotify returns 20 results per page,
	// let's enforce a strict 2,000 track (99 page) limit per search.
	for page := 1; page < 100; page++ {
		// handle track results
		if results.Tracks != nil {
			for _, item := range results.Tracks.Tracks {
				tracks = append(tracks, item)
			}
		}

		err = client.NextTrackResults(results)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			return err
		}
	}
	log.Println("Tracks compiled, retrieving playlist...")
	// Retrieve the tracks currently on the playlist.
	currentTracks, err := getPlaylistTracks(s.PlaylistID, client)
	if err != nil {
		return err
	}
	log.Printf("Playlist has %d total tracks", len(currentTracks))

	// Remove any tracks that are already in the playlist
	tracks = difference(tracks, currentTracks)

	// Push track results to a playlist by ID.
	resultCount := len(tracks)
	log.Printf("Found %v new tracks\n", resultCount)
	if resultCount > 0 {
		var temp []spotify.ID
		for i := len(tracks) - 1; i >= 0; i-- {
			temp = append(temp, tracks[i].ID)
			if len(temp) == 100 || i == 0 {
				snapshot, err := client.AddTracksToPlaylist(s.PlaylistID, temp...)
				if err != nil {
					return err
				}
				log.Printf("Added %v tracks to playlist. Snapshot: %v\n", len(temp), snapshot)
				temp = nil
			}
		}
	}

	return nil
}
