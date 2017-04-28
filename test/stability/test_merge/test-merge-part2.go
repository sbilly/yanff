// Copyright 2017 Intel Corporation.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"github.com/intel-go/yanff/flow"
)

var (
	RECV_PORT1 uint
	RECV_PORT2 uint
	SEND_PORT  uint
)

// Main function for constructing packet processing graph.
func main() {
	flag.UintVar(&RECV_PORT1, "RECV_PORT1", 0, "port for 1st receiver")
	flag.UintVar(&RECV_PORT2, "RECV_PORT2", 1, "port for 2nd receiver")
	flag.UintVar(&SEND_PORT, "SEND_PORT", 0, "port for sender")

	// Init YANFF system at requested number of cores.
	flow.SystemInit(16)

	// Receive packets from 0 and 1 ports
	inputFlow1 := flow.SetReceiver(uint8(RECV_PORT1))
	inputFlow2 := flow.SetReceiver(uint8(RECV_PORT2))

	outputFlow := flow.SetMerger(inputFlow1, inputFlow2)

	flow.SetSender(outputFlow, uint8(SEND_PORT))

	// Begin to process packets.
	flow.SystemStart()
}
