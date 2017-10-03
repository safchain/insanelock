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
	buffer := make([]byte, 10000)
	runtime.Stack(buffer, true)

	i.frames = string(buffer)
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
			err := fmt.Sprintf("\n-- POTENTIAL DEADLOCK --\n")
			err += fmt.Sprintf("--   HOLDING THE LOCK --\n")
			err += fmt.Sprintf("%s\n", i.frames)
			err += fmt.Sprintf("--   TRYING TO LOCK --\n")

			buffer := make([]byte, 10000)
			runtime.Stack(buffer, true)

			err += fmt.Sprintf("%s\n", string(buffer))
			panic(err)
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
