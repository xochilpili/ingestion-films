package utils

import (
	"slices"
	"strings"
)

func ExcludeGenre(genres []string, unwanted []string) bool {
	idx := slices.IndexFunc(genres, func(g string) bool {
		for _, t := range unwanted {
			if strings.EqualFold(g, t) {
				return true
			}
		}
		return false
	})

	return idx != -1
}
