package dcron

import (
	"sync"
	"time"
)

const (
	minCountKeep = 16
)

type Group interface {
	inc(platAt time.Time, fn func() bool) bool
}

func NewGroup(limit int) Group {
	if limit <= 0 {
		return nil
	}
	return &innerGroup{
		limit: limit,
	}
}

type innerGroup struct {
	sync.Mutex

	limit  int
	counts []*groupCount
}

func (g *innerGroup) inc(platAt time.Time, fn func() bool) bool {
	g.Lock()
	defer g.Unlock()
	defer g.tidy()

	var gc *groupCount
	for i := len(g.counts) - 1; i > 0; i-- {
		v := g.counts[i]
		if v.platAt.Equal(platAt) {
			gc = v
		}
	}
	if gc == nil {
		gc = &groupCount{
			platAt: platAt,
			count:  0,
		}
		g.counts = append(g.counts, gc)
	}

	if gc.count < g.limit && fn() {
		gc.count++
		return true
	}
	return false
}

func (g *innerGroup) tidy() {
	if len(g.counts) > 2*minCountKeep {
		g.counts = g.counts[len(g.counts)-minCountKeep:]
	}
}

type groupCount struct {
	platAt time.Time
	count  int
}
