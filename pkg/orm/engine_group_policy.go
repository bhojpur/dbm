package orm

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"math/rand"
	"sync"
	"time"
)

// GroupPolicy is be used by chosing the current slave from slaves
type GroupPolicy interface {
	Slave(*EngineGroup) *Engine
}

// GroupPolicyHandler should be used when a function is a GroupPolicy
type GroupPolicyHandler func(*EngineGroup) *Engine

// Slave implements the chosen of slaves
func (h GroupPolicyHandler) Slave(eg *EngineGroup) *Engine {
	return h(eg)
}

// RandomPolicy implmentes randomly chose the slave of slaves
func RandomPolicy() GroupPolicyHandler {
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	return func(g *EngineGroup) *Engine {
		return g.Slaves()[r.Intn(len(g.Slaves()))]
	}
}

// WeightRandomPolicy implmentes randomly chose the slave of slaves
func WeightRandomPolicy(weights []int) GroupPolicyHandler {
	var rands = make([]int, 0, len(weights))
	for i := 0; i < len(weights); i++ {
		for n := 0; n < weights[i]; n++ {
			rands = append(rands, i)
		}
	}
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	return func(g *EngineGroup) *Engine {
		var slaves = g.Slaves()
		idx := rands[r.Intn(len(rands))]
		if idx >= len(slaves) {
			idx = len(slaves) - 1
		}
		return slaves[idx]
	}
}

// RoundRobinPolicy returns a group policy handler
func RoundRobinPolicy() GroupPolicyHandler {
	var pos = -1
	var lock sync.Mutex
	return func(g *EngineGroup) *Engine {
		var slaves = g.Slaves()
		lock.Lock()
		defer lock.Unlock()
		pos++
		if pos >= len(slaves) {
			pos = 0
		}
		return slaves[pos]
	}
}

// WeightRoundRobinPolicy returns a group policy handler
func WeightRoundRobinPolicy(weights []int) GroupPolicyHandler {
	var rands = make([]int, 0, len(weights))
	for i := 0; i < len(weights); i++ {
		for n := 0; n < weights[i]; n++ {
			rands = append(rands, i)
		}
	}
	var pos = -1
	var lock sync.Mutex
	return func(g *EngineGroup) *Engine {
		var slaves = g.Slaves()
		lock.Lock()
		defer lock.Unlock()
		pos++
		if pos >= len(rands) {
			pos = 0
		}
		idx := rands[pos]
		if idx >= len(slaves) {
			idx = len(slaves) - 1
		}
		return slaves[idx]
	}
}

// LeastConnPolicy implements GroupPolicy, every time will get the least connections slave
func LeastConnPolicy() GroupPolicyHandler {
	return func(g *EngineGroup) *Engine {
		var slaves = g.Slaves()
		connections := 0
		idx := 0
		for i := 0; i < len(slaves); i++ {
			openConnections := slaves[i].DB().Stats().OpenConnections
			if i == 0 {
				connections = openConnections
				idx = i
			} else if openConnections <= connections {
				connections = openConnections
				idx = i
			}
		}
		return slaves[idx]
	}
}
