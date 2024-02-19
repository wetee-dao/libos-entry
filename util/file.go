package util

import (
	"fmt"
	"log"

	"github.com/vedhavyas/go-subkey/v2"
	"github.com/vedhavyas/go-subkey/v2/sr25519"
)

func SaveKey(appFs Fs, appKey subkey.KeyPair, filename string) error {
	_, err := appFs.Stat(filename)
	if err != nil {
		_, err = appFs.Create(filename)
		if err != nil {
			return fmt.Errorf("SaveKey failed to create Key file: %v", err)
		}
	}

	if err := appFs.WriteFile(filename, appKey.Seed(), 0o600); err != nil {
		return fmt.Errorf("SaveKey failed to store Key to file: %v", err)
	}
	return nil
}

func LoadKey(appFs Fs, filename string) (subkey.KeyPair, error) {
	keyBytes, err := appFs.ReadFile(filename)
	if err != nil {
		return nil, nil
	}

	// 没有key文件，返回nil
	if keyBytes == nil && len(keyBytes) == 0 {
		return nil, nil
	}

	appKey, err := sr25519.Scheme{}.FromSeed(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal Key: %v", err)
	}

	return appKey, nil
}

func GetKey(appFs Fs, KeyFile string) (subkey.KeyPair, error) {
	existingKey, err := LoadKey(appFs, KeyFile)
	if err != nil {
		return nil, err
	}

	// generate new Key if not present and store it
	if existingKey == nil {
		log.Println("Key not found. Generating and storing a new Key")
		newKey, err := sr25519.Scheme{}.Generate()
		if err != nil {
			return nil, err
		}
		if err := SaveKey(appFs, newKey, KeyFile); err != nil {
			return nil, err
		}
		return newKey, nil
	}

	log.Println("found Key:", existingKey.SS58Address(42))
	return existingKey, nil
}
