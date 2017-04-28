// Copyright 2017 Intel Corporation.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "github.com/intel-go/yanff/flow"
import "github.com/intel-go/yanff/packet"

import "flag"

var (
	load  uint
	mode  uint
	cores uint

	RECV_PORT1 uint
	RECV_PORT2 uint
	SEND_PORT1 uint
	SEND_PORT2 uint
)

func main() {
	flag.UintVar(&load, "load", 1000, "Use this for regulating 'load intensity', number of iterations")
	flag.UintVar(&mode, "mode", 2, "Benching mode: 2, 12 - two handles; 3, 13 - tree handles; 4, 14 - four handles. 2,3,4 - one flow; 12,13,14 - two flows")
	flag.UintVar(&cores, "cores", 35, "Number of cores to use by system")

	flag.UintVar(&SEND_PORT1, "SEND_PORT1", 1, "port for 1st sender")
	flag.UintVar(&SEND_PORT2, "SEND_PORT2", 1, "port for 2nd sender")
	flag.UintVar(&RECV_PORT1, "RECV_PORT1", 0, "port for 1st receiver")
	flag.UintVar(&RECV_PORT2, "RECV_PORT2", 0, "port for 2nd receiver")

	// Call flag.Parse() here is needed to parse -cores option before 'core' value is passed to SystemInit.
	// We use this only in benchmarks. This call prevents other flags to be added (dpdk related flags inside low package).
	flag.Parse()

	// Initialize YANFF library at requested number of cores
	flow.SystemInit(cores)

	var tempFlow *flow.Flow
	var afterFlow *flow.Flow

	// Receive packets from zero port. One queue will be added automatically.
	firstFlow0 := flow.SetReceiver(uint8(RECV_PORT1))
	firstFlow1 := flow.SetReceiver(uint8(RECV_PORT2))

	firstFlow := flow.SetMerger(firstFlow0, firstFlow1)
	if mode > 10 {
		tempFlow = flow.SetPartitioner(firstFlow, 150, 150)
	}

	// Handle second flow via some heavy function
	flow.SetHandler(firstFlow, heavyFunc, nil)
	flow.SetHandler(firstFlow, heavyFunc, nil)
	if mode%10 > 2 {
		flow.SetHandler(firstFlow, heavyFunc, nil)
	}
	if mode%10 > 3 {
		flow.SetHandler(firstFlow, heavyFunc, nil)
	}
	if mode > 10 {
		flow.SetHandler(tempFlow, heavyFunc, nil)
		flow.SetHandler(tempFlow, heavyFunc, nil)
		if mode%10 > 2 {
			flow.SetHandler(tempFlow, heavyFunc, nil)
		}
		if mode%10 > 3 {
			flow.SetHandler(tempFlow, heavyFunc, nil)
		}
		afterFlow = flow.SetMerger(firstFlow, tempFlow)
	} else {
		afterFlow = firstFlow
	}
	secondFlow := flow.SetPartitioner(afterFlow, 150, 150)

	// Send both flows each one to one port. Queues will be added automatically.
	flow.SetSender(firstFlow, uint8(SEND_PORT1))
	flow.SetSender(secondFlow, uint8(SEND_PORT2))

	flow.SystemStart()
}

func heavyFunc(currentPacket *packet.Packet, context flow.UserContext) {
	currentPacket.ParseEtherIPv4()
	T := (currentPacket.IPv4.DstAddr)
	for j := uint(0); j < load; j++ {
		T += uint32(j)
	}
	currentPacket.IPv4.SrcAddr = 263 + (T)
}
