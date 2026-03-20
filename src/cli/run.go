// NOTICE

// Project Name: Cloaq
// Copyright © 2026 Neil Talap and/or its designated Affiliates.

// This software is licensed under the Dragonfly Public License (DPL) 1.0.

// All rights reserved. The names "Neil Talap" and any associated logos or branding
// are trademarks of the Licensor and may not be used without express written permission,
// except as provided in Section 7 of the License.

// For commercial licensing inquiries or permissions beyond the scope of this
// license, please create an issue in github.

package cli

import (
	network "cloaq/src"
	"cloaq/src/monitor"
	"encoding/hex"

	"cloaq/src/tun"
	"log"
	"runtime"
	"sync/atomic"
)

type Run struct {
	port  int      `yaml:"port"`
	peers []string `yaml:"peers"`
}

var _ Command = (*Run)(nil) // enforcement of an interface

func (s *Run) Name() string {
	return "run"
}

func (s *Run) Description() string {
	return "display configuration run"
}
func (s *Run) Execute(args []string) error {
	log.Println("starting cloaq...")
	log.Println("os:", runtime.GOOS, "arch:", runtime.GOARCH)

	identity, err := network.CreateOrLoadIdentity()
	if err != nil {
		log.Fatal("identity creation failed: ", err)
	}

	nodeID := hex.EncodeToString(identity.PublicKey.Bytes())
	log.Println("node identity loaded:", nodeID[:16]+"...")

	tr, err := network.NewTransport(":9000")
	if err != nil {
		log.Fatal("transport init error:", err)
	}

	dev, err := tun.InitDevice("cloaq0")
	if err != nil {
		log.Fatal("tunnel init error:", err)
	}
	defer func(dev *tun.LinuxDevice) {
		err := dev.Close()
		if err != nil {

		}
	}(dev)

	if err := dev.Start(); err != nil {
		log.Fatal("vnic start error:", err)
	}
	log.Println("vnic initialized:", dev.Name())

	packetChan := make(chan network.Packet, 100)
	m := &monitor.Monitor{}

	network.SafeRuntime("Monitor", func() {
		if err := m.Execute(nil); err != nil {
			log.Printf("monitor error: %v", err)
		}
	})

	network.SafeRuntime("ReadLoop", func() {
		if err := network.ReadLoop(dev, packetChan); err != nil {
			log.Printf("readloop error: %v", err)
		}
	})

	log.Println("ipv6 tun gateway created. processing traffic...")

	for pkt := range packetChan {

		if len(s.peers) > 0 {
			target := s.peers[0]

			onionedData := network.Encapsulate(pkt.Data)

			err := tr.SendTo(target, onionedData)
			if err == nil {

				atomic.AddUint64(&monitor.BytesSent, uint64(len(onionedData)))
				log.Printf("[sent] %d bytes to %s from %s", len(onionedData), target, nodeID[:8])
			} else {
				log.Printf("[error] send failed to %s: %v", target, err)
			}
		}
	}

	return nil
}
