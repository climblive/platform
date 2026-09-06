package domain

import "uuid"

type UUIDGenerator interface {
	Generate() uuid.UUID
}
