package example

import "time"

type Case1 struct {
	Array []string
	Name  int
}

type Config struct {
	TestArray []string
	Case1     *Case1
	Case2     *Case2
	Case3     []Case1
}

type Case2 struct {
	WrongDuration time.Duration
}
