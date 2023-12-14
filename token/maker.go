package token

import "time"

// an interface for managing tokens
type Maker interface {
	// CreateToken create a new token for specific username and duration
	CreateToken(username string, role string, duration time.Duration) (string, *Payload, error)

	// verify token checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
