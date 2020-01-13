package playlists

import (
	"log"

	"github.com/zmb3/spotify"
)

// Duplicates will remove duplicate tracks from the given playlist.
func Duplicates(playlistID spotify.ID, client *spotify.Client) error {
	// Retrieve the tracks currently on the playlist.
	playlistTracks, err := client.GetPlaylistTracks(playlistID)
	if err != nil {
		return err
	}
	total := playlistTracks.Total
	log.Printf("Playlist contains %d total tracks", total)

	// Page through the playlist, caching all unique tracks.
	// Note: Spotify's ID system has been shown to not be a reliable
	// uniqueness identifier. To compare tracks start with track name,
	// then artist name.
	currentTracks := map[string]spotify.SimpleTrack{}
	for page := 1; ; page++ {
		if playlistTracks.Tracks != nil {
			for _, playlistTrack := range playlistTracks.Tracks {
				temp, cached := currentTracks[playlistTrack.Track.Name]
				if !cached {
					currentTracks[playlistTrack.Track.Name] = playlistTrack.Track.SimpleTrack
				} else if temp.Artists[0].Name != playlistTrack.Track.SimpleTrack.Artists[0].Name {
					currentTracks[playlistTrack.Track.Name] = playlistTrack.Track.SimpleTrack
				}
			}
		}

		err = client.NextPage(playlistTracks)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			return err
		}
	}

	// Guard clause will exit if there are no duplicate tracks.
	unique := len(currentTracks)
	log.Printf("Playlist contains %d unique tracks\n", unique)
	if unique >= total {
		log.Printf("Playlist does not contain any duplicates.\n")
		return nil
	}

	// Convert the map of tracks to a slice.
	var tracks []spotify.ID
	for _, simpleTrack := range currentTracks {
		tracks = append(tracks, simpleTrack.ID)
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
