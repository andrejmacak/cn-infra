// Copyright (c) 2017 Cisco and/or its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package core

import (
	"os"
	"os/signal"
)

// ExampleAgent struct with public channel used to close it
type ExampleAgent struct {
	CloseChannel chan *struct{}
}

// EventLoopWithInterrupt init Agent with plugins. Agent can be interrupted from outside using public CloseChannel
func (exampleAgent *ExampleAgent) EventLoopWithInterrupt(agent *Agent) {
	err := agent.Start()
	if err != nil {
		agent.log.Error("Error loading core", err)
		os.Exit(1)
	}
	defer func() {
		err := agent.Stop()
		if err != nil {
			agent.log.Errorf("Agent stop error '%+v'", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt)
	select {
	case <-sigChan:
		agent.log.Info("Interrupt received, returning.")
		return
	case <-exampleAgent.CloseChannel:
		err := agent.Stop()
		if err != nil {
			agent.log.Errorf("Agent stop error '%v'", err)
			os.Exit(1)
		}
		os.Exit(0)
	}
}
