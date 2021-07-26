package model

import (
	"time"
)

type Article struct {
	title   string
	time    time.Time
	content string
}
