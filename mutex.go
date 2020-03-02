package dcron

//go:generate mockgen -source=mutex.go -destination mock_dcron/mutex.go

// Mutex provides distributed mutex for dcron,
// it can be implemented easily via Redis/SQL and so on.
type Mutex interface {
	// SetIfNotExists stores the key/value and return true if the key is not existed,
	// or does nothing and return false.
	// Note that the key/value should be kept for at least one minute.
	// For example, `SetNX(key, value, time.Minute)` via redis.
	SetIfNotExists(key, value string) bool
}
