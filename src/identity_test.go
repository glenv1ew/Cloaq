// NOTICE

// Project Name: Cloaq
// Copyright © 2026 Neil Talap and/or its designated Affiliates.

// This software is licensed under the Dragonfly Public License (DPL) 1.0.

// All rights reserved. The names "Neil Talap" and any associated logos or branding
// are trademarks of the Licensor and may not be used without express written permission,
// except as provided in Section 7 of the License.

// For commercial licensing inquiries or permissions beyond the scope of this
// license, please create an issue in github.

package network

import (
	"bytes"
	"testing"
)

func TestSharedKeyDerivation(t *testing.T) {

	alice, err := GenerateIdentity()
	if err != nil {
		t.Fatal(err)
	}

	bob, err := GenerateIdentity()
	if err != nil {
		t.Fatal(err)
	}

	keyA, err := alice.DeriveSharedKey(bob.PublicKey)
	if err != nil {
		t.Fatal(err)
	}

	keyB, err := bob.DeriveSharedKey(alice.PublicKey)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(keyA, keyB) {
		t.Fatal("shared keys do not match")
	}
}
