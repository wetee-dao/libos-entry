package main

import (
	"fmt"

	"github.com/wetee-dao/libos-entry/libos"
	volume "github.com/wetee-dao/libos-entry/libos/tee_volume"
)

// var key = []byte{208, 35, 86, 92, 106, 206, 205, 212, 85, 182, 48, 244, 30, 163, 14, 59, 81, 204, 83, 127, 13, 184, 187, 146, 125, 93, 90, 16, 135, 65, 23, 233}

// secret_mount implements TEEServer.
func (c CvmServer) secret_mount(major *int64, minor *int64, mount_path *string, container_index *uint64) CrossResponse {
	mountKey := libos.DiskKeys[int(*container_index)][*mount_path].Key
	v, err := volume.NewSecretMount(*major, *minor, mountKey, *mount_path)
	if err != nil {
		fmt.Println("WeTEELOG init device error:", err)
		return CrossResponse{code: 1, data: []byte(err.Error())}
	}

	err = v.Mount()
	if err != nil {
		fmt.Println("Error mounting error:", err)
		return CrossResponse{code: 1, data: []byte(err.Error())}
	}

	// time.Sleep(20 * time.Second)
	// err = v.Unmount()
	// if err != nil {
	// 	fmt.Println("Error unmounting device:", err)
	// 	return CrossResponse{code: 1, data: []byte(err.Error())}
	// }

	return CrossResponse{
		code: 0,
	}
}
