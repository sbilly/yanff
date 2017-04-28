// Copyright 2017 Intel Corporation.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"github.com/intel-go/yanff/flow"
)

var (
	RECV_PORT  uint
	SEND_PORT1 uint
	SEND_PORT2 uint
)

// Main function for constructing packet processing graph.
func main() {
	flag.UintVar(&RECV_PORT, "RECV_PORT", 0, "port for receiver")
	flag.UintVar(&SEND_PORT1, "SEND_PORT1", 0, "port for 1st sender")
	flag.UintVar(&SEND_PORT2, "SEND_PORT2", 1, "port for 2nd sender")

	// Init YANFF system at 16 available cores.
	flow.SystemInit(16)

	// Receive packets from 0 port
	flow1 := flow.SetReceiver(uint8(RECV_PORT))

	flow2 := flow.SetPartitioner(flow1, 1000, 100)

	flow.SetSender(flow1, uint8(SEND_PORT1))
	flow.SetSender(flow2, uint8(SEND_PORT2))

	// Begin to process packets.
	flow.SystemStart()
}
