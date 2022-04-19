/*
创建时间: 2020/3/28
作者: zjy
功能介绍:

*/

package snowflake

import (
	"sync"
	"testing"
	"time"
)

func TestGenID(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(10000)                        //using 10000 goroutine to generate 10000 ids
	fastgid,erro := NewNodeGID(1000,10)
	if erro != nil {
		t.Error(erro)
		return
	}
	results := make(chan int64, 10000) //store result
	for i := 0; i < 10000; i++ {
		go func() {
			id := fastgid.NextId()
			t.Logf("id: %b \t %x \t %d", id, id, id)
			t.Logf("time: %d  \t node:%d \t seq:%d", fastgid.Time(id),fastgid.Node(), fastgid.Step(id))
			results <- id
			defer wg.Done()
		}()
	}

	wg.Wait()
	m := make(map[int64]bool)
	for i := 0; i < 10000; i++ {
		select {
		case id := <-results:
			if _, ok := m[id]; ok {
				t.Errorf("id 重复id: %x", id)
				//return
			} else {
				m[id] = true
			}
		case <-time.After(2 * time.Second):
			t.Errorf("Expect 10000 ids in results, but got %d", i)
			return
		}
	}



}

func BenchmarkGenID(b *testing.B) {
	fastgid,erro := NewNodeGID(1,10)
	if erro != nil {
		b.Error(erro)
		return
	}
	for i := 0; i < b.N; i++ {
		fastgid.NextId()
	}
}

func BenchmarkGenIDP(b *testing.B) {
	fastgid,erro := NewNodeGID(2,10)
	if erro != nil {
		b.Error(erro)
		return
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			fastgid.NextId()
		}
	})
}
