package model

import (
	"time"
)

type Article struct {
	Title   string
	Time    time.Time
	Content string
}
