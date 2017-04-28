// Copyright 2017 Intel Corporation.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "github.com/intel-go/yanff/flow"
import "github.com/intel-go/yanff/packet"
import "github.com/intel-go/yanff/rules"

import "flag"
import "time"

var (
	L2Rules *rules.L2Rules
	L3Rules *rules.L3Rules
	load    uint
	cores   uint

	RECV_PORT  uint
	SEND_PORT1 uint
	SEND_PORT2 uint
)

func main() {
	flag.UintVar(&load, "load", 1000, "Use this for regulating 'load intensity', number of iterations")
	flag.UintVar(&cores, "cores", 16, "Number of cores to use by system")
	flag.UintVar(&RECV_PORT, "RECV_PORT", 0, "port for receiver")
	flag.UintVar(&SEND_PORT1, "SEND_PORT1", 1, "port for 1st sender")
	flag.UintVar(&SEND_PORT2, "SEND_PORT2", 2, "port for 2nd sender")

	// Initialize YANFF library at requested number of cores
	flow.SystemInit(cores)

	// Start regular updating forwarding rules
	L2Rules = rules.GetL2RulesFromJSON("demoL2_ACL.json")
	L3Rules = rules.GetL3RulesFromJSON("demoL3_ACL.json")
	go updateSeparateRules()

	// Receive packets from zero port. One queue will be added automatically.
	firstFlow := flow.SetReceiver(uint8(RECV_PORT))

	// Separate packets for additional flow due to some rules
	secondFlow := flow.SetSeparator(firstFlow, L3Separator, nil)

	// Handle second flow via some heavy function
	flow.SetHandler(firstFlow, heavyFunc, nil)

	// Send both flows each one to one port. Queues will be added automatically.
	flow.SetSender(firstFlow, uint8(SEND_PORT1))
	flow.SetSender(secondFlow, uint8(SEND_PORT2))

	flow.SystemStart()
}

func L3Separator(currentPacket *packet.Packet, context flow.UserContext) bool {
	currentPacket.ParseEtherIPv4()
	localL2Rules := L2Rules
	localL3Rules := L3Rules
	return rules.L2_ACL_permit(currentPacket, localL2Rules) &&
		rules.L3_ACL_permit(currentPacket, localL3Rules)
}

func heavyFunc(currentPacket *packet.Packet, context flow.UserContext) {
	for i := uint(0); i < load; i++ {
	}
}

func updateSeparateRules() {
	for true {
		time.Sleep(time.Second * 5)
		L2Rules = rules.GetL2RulesFromJSON("demoL2_ACL.json")
		L3Rules = rules.GetL3RulesFromJSON("demoL3_ACL.json")
	}
}
