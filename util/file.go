package util

import (
	"fmt"
	"log"

	"github.com/spf13/afero"
	"github.com/vedhavyas/go-subkey/v2"
	"github.com/vedhavyas/go-subkey/v2/sr25519"
)

func SaveKey(appFs afero.Fs, appKey subkey.KeyPair, filename string, sf SecretFunction) error {
	_, err := appFs.Stat(filename)
	if err != nil {
		_, err = appFs.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create Key file: %v", err)
		}
	}
	// 加密数据
	val, err := sf.Encrypt(appKey.Seed())
	if err != nil {
		return fmt.Errorf("failed to encrypt Key: %v", err)
	}

	if err := afero.WriteFile(appFs, filename, val, 0o600); err != nil {
		return fmt.Errorf("failed to store Key to file: %v", err)
	}
	return nil
}

func LoadKey(appFs afero.Fs, filename string, sf SecretFunction) (subkey.KeyPair, error) {
	keyBytes, err := afero.ReadFile(appFs, filename)
	if err != nil {
		return nil, nil
	}

	// 没有key文件，返回nil
	if keyBytes == nil && len(keyBytes) == 0 {
		return nil, nil
	}

	// 解密数据
	keyBytes, err = sf.Decrypt(keyBytes)
	if err != nil {
		return nil, nil
	}

	appKey, err := sr25519.Scheme{}.FromSeed(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal Key: %v", err)
	}

	return appKey, nil
}

func GetKey(appFs afero.Fs, KeyFile string, sf SecretFunction) (subkey.KeyPair, error) {
	log.Println("geting Key")
	existingKey, err := LoadKey(appFs, KeyFile, sf)
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
		if err := SaveKey(appFs, newKey, KeyFile, sf); err != nil {
			return nil, err
		}
		return newKey, nil
	}

	log.Println("found Key:", existingKey.SS58Address(42))
	return existingKey, nil
}
