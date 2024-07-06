package rest

import (
	"strconv"

	"github.com/climblive/platform/backend/internal/domain"
)

func parseResourceID(id string) domain.ResourceID {
	number, _ := strconv.Atoi(id)
	return domain.ResourceID(number)
}
