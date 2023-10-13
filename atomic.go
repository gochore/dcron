package dcron

import "context"

//go:generate go get go.uber.org/mock/mockgen
//go:generate go run go.uber.org/mock/mockgen -source=atomic.go -destination mock_dcron/atomic.go
//go:generate go mod tidy

// Atomic provides distributed atomic operation for dcron,
// it can be implemented easily via Redis/SQL and so on.
type Atomic interface {
	// SetIfNotExists stores the key/value and return true if the key is not existed,
	// or does nothing and return false.
	// Note that the key/value should be kept for at least one minute.
	// For example, `SetNX(key, value, time.Minute)` via redis.
	SetIfNotExists(ctx context.Context, key, value string) bool
}
