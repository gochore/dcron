package dcron

type Job interface {
	Key() string
	Spec() string

	Before(ctx JobContext) (skip bool)
	Run() error
	After(ctx JobContext)
}
