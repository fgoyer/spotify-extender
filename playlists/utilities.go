package playlists

import (
	"github.com/zmb3/spotify"
)

// difference returns the elements in `a` that aren't in `b`.
func difference(a, b []spotify.FullTrack) []spotify.FullTrack {
	mb := make(map[spotify.ID]struct{}, len(b))
	for _, x := range b {
		mb[x.ID] = struct{}{}
	}
	var diff []spotify.FullTrack
	for _, x := range a {
		if _, found := mb[x.ID]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func getPlaylistTracks(playlistID spotify.ID, client *spotify.Client) ([]spotify.FullTrack, error) {
	// Retrieve the tracks currently on the playlist.
	playlistTracks, err := client.GetPlaylistTracks(playlistID)
	if err != nil {
		return nil, err
	}

	currentTracks := make([]spotify.FullTrack, 0)
	for page := 1; ; page++ {
		if playlistTracks.Tracks != nil {
			for _, track := range playlistTracks.Tracks {
				currentTracks = append(currentTracks, track.Track)
			}
		}

		err = client.NextPage(playlistTracks)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			return nil, err
		}
	}

	return currentTracks, nil
}

// Note: Spotify's ID system has been shown to not be a reliable
// uniqueness identifier. To compare tracks start with track name,
// then artist name.
func uniqueTracks(fullTracks []spotify.FullTrack) (map[string]spotify.FullTrack, error) {
	currentTracks := map[string]spotify.FullTrack{}
	if fullTracks != nil {
		for _, fullTrack := range fullTracks {
			temp, cached := currentTracks[fullTrack.Name]
			if !cached {
				currentTracks[fullTrack.Name] = fullTrack
			} else if temp.Artists[0].Name != fullTrack.SimpleTrack.Artists[0].Name {
				currentTracks[fullTrack.Name] = fullTrack
			}
		}
	}

	return currentTracks, nil
}
