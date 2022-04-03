package job

type Job struct {
	ID int
	F  func() error
}

type JobExecutable func() error
