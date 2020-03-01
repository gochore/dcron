package dcron

type Job interface {
	Key() string
	Spec() string

	Before(ctx Context) (skip bool)
	Run() error
	After(ctx Context)
}
