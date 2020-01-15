package playlists

import (
	"log"

	"github.com/zmb3/spotify"
)

// RemoveDuplicates will remove duplicate tracks from the given playlist.
func RemoveDuplicates(playlistID spotify.ID, client *spotify.Client) error {
	// Retrieve the tracks currently on the playlist.
	currentTracks, err := getPlaylistTracks(playlistID, client)
	// playlistTracks, err := client.GetPlaylistTracks(playlistID)

	// total := playlistTracks.Total
	log.Printf("Playlist contains %d total tracks", len(currentTracks))

	unique, err := uniqueTracks(currentTracks)
	if err != nil {
		return err
	}

	// Guard clause will exit if there are no duplicate tracks.
	count := len(unique)
	log.Printf("Playlist contains %d unique tracks\n", count)
	if count >= len(currentTracks) {
		log.Printf("Playlist does not contain any duplicates.\n")
		return nil
	}

	// Convert the map of full tracks to a slice of IDs.
	var tracks []spotify.ID
	for _, fullTrack := range unique {
		tracks = append(tracks, fullTrack.ID)
	}

	// Clear the playlist before adding the remaining tracks.
	err = client.ReplacePlaylistTracks(playlistID, tracks[0])
	if err != nil {
		return err
	}

	var temp []spotify.ID
	for i := len(tracks) - 1; i >= 1; i-- {
		temp = append(temp, tracks[i])
		if len(temp) == 100 || i == 1 {
			snapshot, err := client.AddTracksToPlaylist(playlistID, temp...)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Added %d tracks to playlist. Snapshot: %s\n", len(temp), snapshot)
			temp = nil
		}
	}

	log.Println("Duplicate tracks removed!")
	return nil
}
