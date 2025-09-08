package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	volume "github.com/wetee-dao/libos-entry/libos/tee_volume"
)

func main() {
	ctx := context.Background()
	devicePath := "/dev/sda"
	mappingId := "tee_vol"

	luksDev, err := volume.NewCryptLuks(devicePath, "/run/tee-vol/luks_key.bin", mappingId)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := luksDev.Detach(ctx); err != nil {
			fmt.Println("Error detaching LUKS device:", err)
		}
	}()

	isLuks, err := luksDev.CheckFormat(ctx)
	if err != nil {
		panic(err)
	}
	if !isLuks {
		fmt.Println("Device is not a LUKS device, formatting it")
		if err := luksDev.Format(ctx); err != nil {
			panic(fmt.Errorf("formatting device %s as LUKS: %w", devicePath, err))
		}
		fmt.Println("Device formatted successfully")
	} else {
		fmt.Println("Device is already a LUKS device")
	}

	fmt.Println("Opening LUKS device", "mappingName", mappingId)
	if err := luksDev.Attach(ctx); err != nil {
		panic(fmt.Errorf("opening LUKS device %s: %w", devicePath, err))
	}
	fmt.Println("LUKS device opened successfully", "mappingName", mappingId)

	fmt.Println("Check if ext4 filesystem is present on device")
	isExt4, err := luksDev.CheckExt4Format(ctx)
	if err != nil {
		panic(fmt.Errorf("checking if device is ext4: %w", err))
	}

	if !isExt4 {
		fmt.Println("No ext4 filesystem identified, creating new ext4 filesystem")
		if err := luksDev.FormatExt4(ctx); err != nil {
			panic(fmt.Errorf("formatting device %s to ext4: %w", "/dev/mapper/"+mappingId, err))
		}
		fmt.Println("Created ext4 filesystem on device")
	} else {
		fmt.Println("ext4 filesystem present on device")
	}

	mountTo := "/mnt/x"
	fmt.Printf("Mounting device %s to %s", "/dev/mapper/"+mappingId, mountTo)
	if err := mount(ctx, "/dev/mapper/"+mappingId, mountTo); err != nil {
		panic(err)
	}

	time.Sleep(20 * time.Second)
}

// mount wraps the mount command and attaches the device to the specified mountPoint.
func mount(ctx context.Context, devPath, mountPoint string) error {
	if err := os.MkdirAll(mountPoint, 0o755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	if out, err := exec.CommandContext(ctx, "mount", "-o", "sync,data=journal", devPath, mountPoint).CombinedOutput(); err != nil {
		return fmt.Errorf("mount: %w, output: %q", err, out)
	}
	return nil
}
