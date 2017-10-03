/*
 * Copyright (C) 2017 Sylvain Afchain
 *
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *
 */

package insamelock

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

var activated bool

type RWMutex struct {
	mutex  sync.RWMutex
	frames string
}

func (i *RWMutex) saveFrames() {
	i.frames = ""

	pc := make([]uintptr, 15)
	n := runtime.Callers(3, pc)

	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, more := frames.Next()
		i.frames += fmt.Sprintf("%s:%d\n", frame.File, frame.Line)

		if !more {
			break
		}
	}
}

func (i *RWMutex) Lock() {
	if !activated {
		i.mutex.Lock()
		return
	}

	got := make(chan bool)
	go func() {
		select {
		case <-got:
		case <-time.After(30 * time.Second):
			panic(fmt.Sprintf("\n-- POTENTIAL DEADLOCK --\n%s-- POTENTIAL DEADLOCK --\n", i.frames))
		}
	}()

	i.mutex.Lock()

	// stop the timer
	got <- true

	// save the current stack
	i.saveFrames()
}

func (i *RWMutex) Unlock() {
	i.mutex.Unlock()
}

func (i *RWMutex) RLock() {
	i.mutex.RLock()
}

func (i *RWMutex) RUnlock() {
	i.mutex.RUnlock()
}

func init() {
	if os.Getenv("INSANELOCK") == "true" {
		activated = true
	}
}
