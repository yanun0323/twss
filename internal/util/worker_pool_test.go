package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkerPool(t *testing.T) {
	wp := NewWorkerPool("test", 10)
	wp.Run()
	x := 0
	for i := 0; i < 10000; i++ {
		func(num int) {
			wp.Push(func() {
				x++
			})
		}(i)
	}
	assert.Nil(t, wp.Shutdown(3*time.Second))
}
