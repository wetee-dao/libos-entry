package util

// func SetKey(sfs Fs, appKey subkey.KeyPair, filename string) error {
// 	_, err := sfs.Stat(filename)
// 	if err != nil {
// 		_, err = sfs.Create(filename)
// 		if err != nil {
// 			return errors.New("SetKey: " + err.Error())
// 		}
// 	}

// 	if err := sfs.WriteFile(filename, appKey.Seed(), 0o600); err != nil {
// 		return errors.New("SetKey to file: " + err.Error())
// 	}
// 	return nil
// }

// func GetKey(sfs Fs, filename string) (subkey.KeyPair, error) {
// 	keyBytes, err := sfs.ReadFile(filename)
// 	// 没有key文件，创建一个新的key
// 	if err != nil || len(keyBytes) == 0 {
// 		util.LogWithRed("GetKey", "Key not found. Generating and storing a new Key "+err.Error())
// 		newKey, err := sr25519.Scheme{}.Generate()
// 		if err != nil {
// 			return nil, err
// 		}
// 		if err := SetKey(sfs, newKey, filename); err != nil {
// 			return nil, err
// 		}
// 		return newKey, nil
// 	}

// 	appKey, err := sr25519.Scheme{}.FromSeed(keyBytes)
// 	if err != nil {
// 		return nil, errors.New("GetKey: " + err.Error())
// 	}

// 	return appKey, nil
// }
