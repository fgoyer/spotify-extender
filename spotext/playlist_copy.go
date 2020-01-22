package spotext

import (
	"log"

	"github.com/zmb3/spotify"
)

// CopyTracks copies unique tracks from playlistA to playlistB
func CopyTracks(playlistA string, playlistB string, client *spotify.Client) error {
	fromID := spotify.ID(playlistA)
	toID := spotify.ID(playlistB)

	from, err := getPlaylistTracks(fromID, client)
	if err != nil {
		return err
	}

	to, err := getPlaylistTracks(toID, client)
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
				snapshot, err := client.AddTracksToPlaylist(toID, temp...)
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
