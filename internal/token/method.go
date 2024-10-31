package token

import "github.com/google/uuid"

type JWT interface {
	CreateToken(userID uuid.UUID, isAccess bool) (string, error)
}
