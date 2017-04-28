// Copyright 2017 Intel Corporation.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"

	"github.com/intel-go/yanff/flow"
	"github.com/intel-go/yanff/packet"

	"github.com/intel-go/yanff/test/stability/test1/common"
)

var (
	cores     uint
	SEND_PORT uint
	RECV_PORT uint
)

// Main function for constructing packet processing graph.
func main() {
	flag.UintVar(&RECV_PORT, "RECV_PORT", 0, "port for receiver")
	flag.UintVar(&SEND_PORT, "SEND_PORT", 1, "port for sender")

	// Init YANFF system at 16 available cores.
	flow.SystemInit(16)

	inputFlow := flow.SetReceiver(uint8(RECV_PORT))
	flow.SetHandler(inputFlow, fixPacket, nil)
	flow.SetSender(inputFlow, uint8(SEND_PORT))

	// Begin to process packets.
	flow.SystemStart()
}

func fixPacket(pkt *packet.Packet, context flow.UserContext) {
	offset := pkt.ParseL4Data()
	if offset < 0 {
		println("ParseL4 returned negative value", offset)
		println("TEST FAILED")
		return
	}

	ptr := (*common.Packetdata)(pkt.Data)
	if ptr.F2 != 0 {
		fmt.Printf("Bad data found in the packet: %x\n", ptr.F2)
		println("TEST FAILED")
		return
	}

	ptr.F2 = ptr.F1
}
