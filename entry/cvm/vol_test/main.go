package main

import (
	"fmt"

	volume "github.com/wetee-dao/libos-entry/libos/tee_volume"
)

var key = []byte{208, 35, 86, 92, 106, 206, 205, 212, 85, 182, 48, 244, 30, 163, 14, 59, 81, 204, 83, 127, 13, 184, 187, 146, 125, 93, 90, 16, 135, 65, 23, 233}

func main() {
	v, err := volume.NewSecretMount(8, 0, key, "/mnt/x")
	if err != nil {
		fmt.Println("Error mounting device:", err)
		return
	}

	err = v.Mount()
	if err != nil {
		fmt.Println("Error mounting device:", err)
		return
	}

	// time.Sleep(20 * time.Second)
	// err = v.Unmount()
	// if err != nil {
	// 	fmt.Println("Error unmounting device:", err)
	// 	return
	// }
}
