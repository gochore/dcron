package dcron

type Mutex interface {
	SetIfNotExists(key, value string) bool
}
