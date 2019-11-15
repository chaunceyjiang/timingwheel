package timingwheel

import (
	"errors"
	"time"
)

type TimingWheel struct {
	tick int64 // 每个格的大小 毫秒

	wheelSize int64 // 时间轮的格数

	interval int64

	currentTime int64 // 当前时间

	overFlowWheel *TimingWheel // 溢出的时间轮

	exitC chan struct{}

	buckets []*bucket
}

func NewTimingWheel(tick time.Duration, wheelSize int64) *TimingWheel {
	tickMS := int64(tick / time.Millisecond)
	if tickMS <= 0 {
		panic(errors.New("tick must be greater than or equal to 1ms"))
	}
	if wheelSize <= 0 {
		panic(errors.New("wheelSize must be greater than or equal to 1"))
	}
	currentTime := timeToMS(time.Now().UTC())

	return newTimingWheel(tickMS, currentTime, wheelSize)
}

func newTimingWheel(tick int64, startMs int64, wheelSize int64) *TimingWheel {
	buckets := make([]*bucket, wheelSize)
	for i := range buckets {
		buckets[i] = newBucket()
	}
	tw := &TimingWheel{
		tick:      tick,
		wheelSize: wheelSize,
		interval:  tick * wheelSize,
		buckets:   buckets,
		exitC:     make(chan struct{}),
	}

	return tw
}

//func (tw *TimingWheel) Start() {
//	ticker := time.NewTimer(tw.tick)
//}
func (tw *TimingWheel) Stop() {

}
