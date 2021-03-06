/*
 * Copyright 2014 Jason Woods.
 *
 * This file is a modification of code from Logstash Forwarder.
 * Copyright 2012-2013 Jordan Sissel and contributors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"github.com/driskell/log-courier/src/lc-lib/core"
	"github.com/driskell/log-courier/src/lc-lib/registrar"
	"testing"
	"time"
)

func newTestStdinRegistrar() (*core.Pipeline, *StdinRegistrar) {
	pipeline := core.NewPipeline()
	return pipeline, newStdinRegistrar(pipeline)
}

func newEventSpool(offset int64) []*core.EventDescriptor {
	// Prepare an event spool with single event of specified offset
	return []*core.EventDescriptor{
		&core.EventDescriptor{
			Stream: nil,
			Offset: offset,
			Event: []byte{},
		},
	}
}

func TestStdinRegistrarWait(t *testing.T) {
	p, r := newTestStdinRegistrar()

	// Start the stdin registrar
	go func() {
		r.Run()
	}()

	c := r.Connect()
	c.Add(registrar.NewAckEvent(newEventSpool(13)))
	c.Send()

	r.Wait(13)

	wait := make(chan int)
	go func() {
		p.Wait()
		wait <- 1
	}()

	select {
	case <-wait:
		break
	case <-time.After(5 * time.Second):
		t.Error("Timeout waiting for stdin registrar shutdown")
		return
	}

	if r.last_offset != 13 {
		t.Error("Last offset was incorrect: ", r.last_offset)
	} else if r.wait_offset == nil || *r.wait_offset != 13 {
		t.Error("Wait offset was incorrect: ", r.wait_offset)
	}
}
