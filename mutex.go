package dcron

//go:generate mockgen -source=mutex.go -destination mock_dcron/mutex.go
type Mutex interface {
	SetIfNotExists(key, value string) bool
}
