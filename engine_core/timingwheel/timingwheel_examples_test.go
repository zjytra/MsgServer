package timingwheel_test

import (
	"bytes"
	"fmt"
	"github.com/zjytra/MsgServer/devlop/xutil/osutil"
	"github.com/zjytra/MsgServer/devlop/xutil/timeutil"
	"testing"
	"time"

	"github.com/zjytra/MsgServer/engine_core/timingwheel"
)

func Example_startTimer() {
	tw := timingwheel.NewTimingWheel(time.Millisecond, 20)
	tw.Start()
	defer tw.Stop()

	exitC := make(chan time.Time, 1)
	tw.AfterFunc(time.Second, func() {
		fmt.Println("The timer fires")
		exitC <- time.Now().UTC()
	})

	<-exitC

	// Output:
	// The timer fires
}

func Example_stopTimer() {
	tw := timingwheel.NewTimingWheel(time.Millisecond, 20)
	tw.Start()
	defer tw.Stop()

	t := tw.AfterFunc(time.Second, func() {
		fmt.Println("The timer fires")
	})

	<-time.After(900 * time.Millisecond)
	// Stop the timer before it fires
	t.Stop()

	// Output:
	//
}

func Test_AfterFunc(t *testing.T) {
	tw := timingwheel.NewTimingWheel(time.Millisecond, 20)
	tw.Start()
	defer tw.Stop()
	t.Logf("now %v",timeutil.GetTimeNow())
	times := []struct{
		timed time.Duration
		te timingwheel.TimeEvent
	}{
		{time.Millisecond, func() {
			t.Logf("AfterFunc time.Second %v ,协程 %d",timeutil.GetTimeNow(),osutil.GetGID())
		}},
		{time.Millisecond * 20, func() {
			t.Logf("AfterFunc time.Second* 20 %v 协程 %d",timeutil.GetTimeNow(),osutil.GetGID())
		}},
		{time.Millisecond * 30, func() {
			t.Logf("AfterFunc time.Second* 30 %v 协程 %d",timeutil.GetTimeNow(),osutil.GetGID())
		}},
	}
	for _, s := range times {
		tw.AfterFunc(s.timed,s.te)
	}
	<- time.After(time.Minute)

}

func TestScheduleFunc(t *testing.T) {
	tw := timingwheel.NewTimingWheel(time.Millisecond, 20)
	tw.Start()
	defer tw.Stop()
	fmt.Println("now",timeutil.GetTimeNow())
	t1 := tw.ScheduleFunc(&EveryScheduler{time.Millisecond * 20}, func() {
		fmt.Println("The timer fires",timeutil.GetCurrentTimeMs())

	})

	<- time.After(time.Minute)
	t1.Stop()
	<- time.After(time.Second * 100)
	fmt.Println("TestScheduleFunc end")
}

func BenchmarkAfterFunc(b *testing.B) {
	tw := timingwheel.NewTimingWheel(time.Millisecond, 20)
	tw.Start()

	b.ResetTimer()
	for i := 0 ;i <b.N;i++ {
		b.RunParallel(func(pb *testing.PB) {
			var buf bytes.Buffer
			for pb.Next() {
				// The loop body is executed b.N times total across all goroutines.
				buf.Reset()
				tw.ScheduleFunc(&EveryScheduler{time.Millisecond * 20}, func() {
					fmt.Println("The timer fires",timeutil.GetCurrentTimeMs())
				})
			}
		})
	}
	<- time.After(time.Second * 30)
	tw.Stop()
	
}