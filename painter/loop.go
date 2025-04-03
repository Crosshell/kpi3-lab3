package painter

import (
	"image"
	"image/color"
	"sync"

	"golang.org/x/exp/shiny/screen"
)

type Receiver interface {
	Update(t screen.Texture)
}

type Loop struct {
	Receiver Receiver

	next    screen.Texture
	prev    screen.Texture
	state   TextureState
	stateMu sync.RWMutex

	mq       messageQueue
	stopChan chan struct{}
	stopReq  bool
}

var size = image.Pt(800, 800)

type TextureState struct {
	Background color.Color
	BgRect     *BgRect
	Figures    []Figure
}

func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)
	l.state = TextureState{Background: color.Black}
	l.stopChan = make(chan struct{})

	go func() {
		for {
			select {
			case <-l.stopChan:
				return
			default:
				op := l.mq.pull()
				if op != nil {
					l.stateMu.Lock()
					if update := op.Do(l.next); update {
						l.Receiver.Update(l.next)
						l.next, l.prev = l.prev, l.next
					}
					l.stateMu.Unlock()
				}
			}
		}
	}()
}

func (l *Loop) Post(op Operation) {
	l.mq.push(op)
}

func (l *Loop) StopAndWait() {
	l.stopReq = true
	l.stopChan <- struct{}{}
}

type messageQueue struct {
	ops  []Operation
	mu   sync.Mutex
	cond *sync.Cond
}

func (mq *messageQueue) push(op Operation) {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	mq.ops = append(mq.ops, op)
	if mq.cond != nil {
		mq.cond.Signal()
	}
}

func (mq *messageQueue) pull() Operation {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	for len(mq.ops) == 0 {
		if mq.cond == nil {
			mq.cond = sync.NewCond(&mq.mu)
		}
		mq.cond.Wait()
	}

	op := mq.ops[0]
	mq.ops[0] = nil
	mq.ops = mq.ops[1:]
	return op
}

func (mq *messageQueue) empty() bool {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	return len(mq.ops) == 0
}