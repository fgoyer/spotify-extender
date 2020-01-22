package spotext

import (
	"log"

	"github.com/zmb3/spotify"
)

// CopyTracks copies unique tracks from playlistA to playlistB
func CopyTracks(playlistA spotify.ID, playlistB spotify.ID, client *spotify.Client) error {
	from, err := getPlaylistTracks(playlistA, client)
	if err != nil {
		return err
	}

	to, err := getPlaylistTracks(playlistB, client)
	if err != nil {
		return err
	}

	tracks := difference(from, to)

	// Push track results to a playlist by ID.
	resultCount := len(tracks)
	log.Printf("Copying %v new tracks\n", resultCount)
	if resultCount > 0 {
		var temp []spotify.ID
		for i := len(tracks) - 1; i >= 0; i-- {
			temp = append(temp, tracks[i].ID)
			if len(temp) == 100 || i == 0 {
				snapshot, err := client.AddTracksToPlaylist(playlistB, temp...)
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
