package logger

import (
	"testing"
	"time"
)

func TestInfof(t *testing.T)  {
	Infof("hello world!")
	time.Sleep(time.Second * 1)
}
