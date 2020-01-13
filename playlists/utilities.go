package playlists

import (

	"github.com/zmb3/spotify"
)

// difference returns the elements in `a` that aren't in `b`.
func difference(a, b []spotify.ID) []spotify.ID {
	mb := make(map[spotify.ID]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []spotify.ID
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
