package utils

import (
	"fmt"
	"log"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/spf13/afero"
)

func SaveKey(appFs afero.Fs, appKey *signature.KeyringPair, filename string) error {
	if err := afero.WriteFile(appFs, filename, []byte(appKey.URI), 0o600); err != nil {
		return fmt.Errorf("failed to store Key to file: %v", err)
	}
	return nil
}

func LoadKey(appFs afero.Fs, filename string) (*signature.KeyringPair, error) {
	keyBytes, err := afero.ReadFile(appFs, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read Key from file: %v", err)
	}
	appKey, err := signature.KeyringPairFromSecret(string(keyBytes), 42)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal Key: %v", err)
	}

	return &appKey, nil
}

func GetKey(appFs afero.Fs, KeyFile string) (*signature.KeyringPair, error) {
	log.Println("geting Key")
	existingKey, err := LoadKey(appFs, KeyFile)
	if err != nil {
		return nil, err
	}

	// generate new Key if not present and store it
	if existingKey == nil {
		log.Println("Key not found. Generating and storing a new Key")
		newKey, err := GetSignerKey()
		if err != nil {
			return nil, err
		}
		if err := SaveKey(appFs, newKey, KeyFile); err != nil {
			return nil, err
		}
		return newKey, nil
	}

	log.Println("found Key:", existingKey.Address)
	return existingKey, nil
}
