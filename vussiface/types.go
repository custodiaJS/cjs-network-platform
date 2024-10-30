package vussiface

import cjsnetworkplatform "github.com/custodiaJs/cjs-network-platform"

type VirtualDevice struct {
	name    string
	inChan  chan cjsnetworkplatform.Packet
	outChan chan cjsnetworkplatform.Packet
}
