// Copyright 2017 Intel Corporation.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"github.com/intel-go/yanff/flow"
)

var (
	SEND_PORT uint
	RECV_PORT uint
)

// Main function for constructing packet processing graph.
func main() {
	flag.UintVar(&SEND_PORT, "SEND_PORT", 1, "port for sender")
	flag.UintVar(&RECV_PORT, "RECV_PORT", 0, "port for receiver")

	// Initialize YANFF library at requested number of cores.
	flow.SystemInit(16)

	// Receive packets from 0 port and send to 1 port.
	flow1 := flow.SetReceiver(uint8(RECV_PORT))
	flow.SetSender(flow1, uint8(SEND_PORT))

	// Begin to process packets.
	flow.SystemStart()
}
