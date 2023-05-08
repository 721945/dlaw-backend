package utils

import "github.com/721945/dlaw-backend/models"

func ContainsTag(tags []models.Tag, tag models.Tag) bool {
	for _, t := range tags {
		if t.Name == tag.Name {
			return true
		}
	}

	return false
}
