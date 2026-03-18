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
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
)

type Identity struct {
	PrivateKey *ecdh.PrivateKey
	PublicKey  *ecdh.PublicKey
}

func (i *Identity) String() string {
	return hex.EncodeToString(i.PublicKey.Bytes())
}

func identityPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(home, ".cloaq")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}

	return filepath.Join(dir, "identity.key"), nil
}

func saveIdentity(path string, key []byte) error {
	return os.WriteFile(path, key, 0600)
}

func loadIdentity(path string) (*ecdh.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ecdh.X25519().NewPrivateKey(data)
}

func GenerateIdentity() (*Identity, error) {
	path, err := identityPath()
	if err != nil {
		return nil, err
	}

	var priv *ecdh.PrivateKey

	if _, err := os.Stat(path); err == nil {
		priv, err = loadIdentity(path)
		if err != nil {
			return nil, err
		}
	} else {

		priv, err = ecdh.X25519().GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}

		if err := saveIdentity(path, priv.Bytes()); err != nil {
			return nil, err
		}
	}

	return &Identity{
		PrivateKey: priv,
		PublicKey:  priv.Public().(*ecdh.PublicKey),
	}, nil
}
func (i *Identity) DeriveSharedKey(peerPub *ecdh.PublicKey) ([]byte, error) {
	// Perform ECDH key exchange
	secret, err := i.PrivateKey.ECDH(peerPub)
	if err != nil {
		return nil, err
	}
	// Hash the shared secret to derive symmetric key
	hash := sha256.Sum256(secret)
	return hash[:], nil
}
