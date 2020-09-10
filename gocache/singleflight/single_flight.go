package singleflight

import "sync"

type Call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu sync.Mutex
	m  map[string]*Call
}

func (g *Group) Do(key string, fn func(key string) (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*Call)
	}

	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	} else {
		c := new(Call)
		g.m[key] = c
		g.mu.Unlock()

		c.wg.Add(1)
		c.val, c.err = fn(key)
		c.wg.Done()

		g.mu.Lock()
		delete(g.m, key)
		g.mu.Unlock()

		return c.val, c.err
	}
}
