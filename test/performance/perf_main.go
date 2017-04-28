// Copyright 2017 Intel Corporation.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "github.com/intel-go/yanff/flow"
import "github.com/intel-go/yanff/packet"

import "flag"

var (
	load   uint
	loadRW uint
	cores  uint

	RECV_PORT1 uint
	RECV_PORT2 uint
	SEND_PORT1 uint
	SEND_PORT2 uint
)

func main() {
	flag.UintVar(&load, "load", 1000, "Use this for regulating 'load intensity', number of iterations")
	flag.UintVar(&loadRW, "loadRW", 50, "Use this for regulating 'load read/write intensity', number of iterations")
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

	// Receive packets from zero port. One queue per receive will be added automatically.
	firstFlow0 := flow.SetReceiver(uint8(RECV_PORT1))
	firstFlow1 := flow.SetReceiver(uint8(RECV_PORT2))

	firstFlow := flow.SetMerger(firstFlow0, firstFlow1)

	// Handle second flow via some heavy function
	flow.SetHandler(firstFlow, heavyFunc, nil)

	// Split for two senders and send
	secondFlow := flow.SetPartitioner(firstFlow, 150, 150)

	flow.SetSender(firstFlow, uint8(SEND_PORT1))
	flow.SetSender(secondFlow, uint8(SEND_PORT2))

	flow.SystemStart()
}

func heavyFunc(currentPacket *packet.Packet, context flow.UserContext) {
	currentPacket.ParseEtherIPv4()
	T := currentPacket.IPv4.DstAddr
	for j := uint32(0); j < uint32(load); j++ {
		T += j
	}
	for i := uint32(0); i < uint32(loadRW); i++ {
		currentPacket.IPv4.DstAddr = currentPacket.IPv4.SrcAddr + i
	}
	currentPacket.IPv4.SrcAddr = 263 + (T)
}
