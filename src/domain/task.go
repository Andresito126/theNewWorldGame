package domain

import "time"

type Task struct {
	Type string
	Resource string
	Duration  time.Duration 
	TargetX int
	TargetY int
	TargetID int
}


