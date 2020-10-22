package utils

import (
	"digger/models"
	"math/rand"
	"sync"
	"time"
)

const (
	decreaseFactorStep = 2 // 0.5 -> 0.25 -> 0.125 ...
)

var (
	rander *rand.Rand
)

func init() {
	rander = rand.New(rand.NewSource(time.Now().UnixNano()))
}

type proxyState struct {
	proxyId      int
	lock         *sync.Mutex
	weightsRange [2]int // private weights range
	// choose factor, range (0,1]
	factor float32
}

func (p *proxyState) Fail() {
	p.lock.Lock()
	p.lock.Unlock()

	if p.factor < 0.1 {
		return
	}
	p.factor = p.factor / decreaseFactorStep
}

func (p *proxyState) Success() {
	p.lock.Lock()
	p.lock.Unlock()

	if p.factor >= 1 {
		return
	}
	p.factor = p.factor * decreaseFactorStep
}

type proxyStageManager struct {
	lock       *sync.Mutex
	proxiesMap map[int][]*proxyState // taskId -> []proxyStage
}

// select proxy by weights
func (s *proxyStageManager) selectProxy(taskId int, selectRange []*models.Proxy) (*models.Proxy, *proxyState) {
	s.lock.Lock()
	defer s.lock.Unlock()

	states := s.proxiesMap[taskId]
	if states == nil {
		states = make([]*proxyState, len(selectRange))
		for i := range selectRange {
			states[i] = &proxyState{
				proxyId:      selectRange[i].Id,
				lock:         new(sync.Mutex),
				weightsRange: [2]int{},
				factor:       1,
			}
		}
		s.proxiesMap[taskId] = states
	}
	var lastInt = 0
	for _, s := range states {
		rangeMax := int(s.factor * 100)
		s.weightsRange[0] = lastInt
		s.weightsRange[1] = lastInt + rangeMax
		lastInt = s.weightsRange[1] + 1
	}
	loc := rander.Intn(lastInt)
	proxyId := 0
	var state *proxyState
	for _, s := range states {
		if loc >= s.weightsRange[0] && loc <= s.weightsRange[1] {
			proxyId = s.proxyId
			state = s
		}
	}
	for _, p := range selectRange {
		if p.Id == proxyId {
			return p, state
		}
	}
	return nil, nil
}
