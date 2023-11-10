package dcron

// Statistics records statistics info for a cron or a job.
type Statistics struct {
	TotalTask   int64 // Total count of tasks processed
	PassedTask  int64 // Number of tasks successfully executed
	FailedTask  int64 // Number of tasks that failed during execution due to errors
	SkippedTask int64 // Number of tasks skipped due to BeforeFunc returning true
	MissedTask  int64 // Number of tasks executed by other instances

	TotalRun   int64 // Total count of execution runs
	PassedRun  int64 // Number of successfully executed runs
	FailedRun  int64 // Number of runs that have failed due to errors
	RetriedRun int64 // Number of runs that encountered errors and were subsequently retried
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
