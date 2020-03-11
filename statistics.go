package dcron

// Statistics records statistics info for a cron or a job.
type Statistics struct {
	TotalTask   int64
	PassedTask  int64
	FailedTask  int64
	SkippedTask int64
	MissedTask  int64

	TotalRun   int64
	PassedRun  int64
	FailedRun  int64
	RetriedRun int64
}

// Add return a new Statistics with two added.
func (s Statistics) Add(delta Statistics) Statistics {
	s.TotalTask += delta.TotalTask
	s.PassedTask += delta.PassedTask
	s.FailedTask += delta.FailedTask
	s.SkippedTask += delta.SkippedTask
	s.MissedTask += delta.MissedTask
	s.TotalRun += delta.TotalRun
	s.PassedRun += delta.PassedRun
	s.FailedRun += delta.FailedRun
	s.RetriedRun += delta.RetriedRun
	return s
}
