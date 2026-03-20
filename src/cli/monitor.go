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
	"log"
	"runtime"
	"time"
)

type Monitor struct{}

var _ Command = (*Monitor)(nil) // enforcement of an interface

func (h *Monitor) Name() string {
	return "monitor"
}

func (h *Monitor) Description() string {
	return "display monitoring information"
}

func (h *Monitor) Execute(args []string) error {

	const (
		ByteToMiB = 1024 * 1024
		Interval  = 30 * time.Second
	)
	go func() {
		for {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			log.Printf("[monitor] alloc: %d MB, sys: %d MB, goroutines: %d",
				m.Alloc/ByteToMiB,
				m.Sys/ByteToMiB,
				runtime.NumGoroutine(),
			)

			time.Sleep(Interval)
		}
	}()

	return nil
}
