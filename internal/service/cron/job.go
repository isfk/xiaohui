package cron

type Job interface {
	Run()
	String() string
}
