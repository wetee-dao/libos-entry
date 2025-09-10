package volume

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type SecretVolume struct {
	device     *CryptLuks
	mountPoint string
}

func NewSecretMount(maj, min int64, key []byte, mountPoint string) (*SecretVolume, error) {
	ctx := context.Background()
	devicePath, err := GetDeviceNameByMajMin(maj, min)
	if err != nil {
		return nil, err
	}
	secretId := fmt.Sprintf("vol_%d_%d", maj, min)
	keyPath := "/run/secret__" + fmt.Sprintf("%d_%d", maj, min) + ".key"

	os.WriteFile(keyPath, key, 0600)
	luksDev, err := NewCryptLuks(devicePath, keyPath, secretId)
	if err != nil {
		return nil, err
	}

	isLuks, err := luksDev.CheckFormat(ctx)
	if err != nil {
		return nil, err
	}
	if !isLuks {
		fmt.Println("WeTEELOG secretMount device is not a LUKS device, formatting it")
		if err := luksDev.Format(ctx); err != nil {
			return nil, fmt.Errorf("formatting device %s as LUKS: %w", devicePath, err)
		}
		fmt.Println("WeTEELOG secretMount device formatted successfully")
	} else {
		fmt.Println("WeTEELOG secretMount device is already a LUKS device")
	}

	fmt.Println("WeTEELOG secretMount opening LUKS device", "mappingName", secretId)
	if err := luksDev.Attach(ctx); err != nil {
		return nil, fmt.Errorf("opening LUKS device %s: %w", devicePath, err)
	}
	fmt.Println("WeTEELOG secretMount LUKS device opened successfully", "mappingName", secretId)

	isExt4, err := luksDev.CheckExt4Format(ctx)
	if err != nil {
		luksDev.Detach(ctx)
		return nil, fmt.Errorf("checking if device is ext4: %w", err)
	}

	if !isExt4 {
		fmt.Println("No ext4 filesystem identified, creating new ext4 filesystem")
		if err := luksDev.FormatExt4(ctx); err != nil {
			luksDev.Detach(ctx)
			return nil, fmt.Errorf("formatting device %s to ext4: %w", "/dev/mapper/"+secretId, err)
		}
		fmt.Println("WeTEELOG secretMount created ext4 filesystem on device")
	} else {
		fmt.Println("WeTEELOG secretMount ext4 filesystem present on device")
	}

	isMounted, alreadyMount, err := isMounted("/dev/mapper/" + secretId)
	if err != nil {
		return nil, err
	}

	if isMounted {
		err := umount(ctx, alreadyMount)
		if err != nil {
			return nil, err
		}
	}

	return &SecretVolume{device: luksDev, mountPoint: mountPoint}, nil
}

func (s *SecretVolume) Mount() error {
	fmt.Printf("WeTEELOG secretMount mounting device %s to %s \n", "/dev/mapper/"+s.device.secretId, s.mountPoint)
	if err := mount(context.Background(), "/dev/mapper/"+s.device.secretId, s.mountPoint); err != nil {
		return err
	}

	return nil
}

func (s *SecretVolume) Unmount() error {
	ctx := context.Background()
	defer func() {
		if err := s.device.Detach(ctx); err != nil {
			fmt.Println("Error detaching LUKS device:", err)
		}
	}()
	return umount(ctx, s.mountPoint)
}

// GetDeviceNameByMajMin 根据主设备号(MAJ)和次设备号(MIN)获取设备路径（如 /dev/xxxx）
func GetDeviceNameByMajMin(maj, min int64) (string, error) {
	// /sys/dev/block/<MAJ:MIN> 是指向设备实际路径的符号链接
	sysPath := fmt.Sprintf("/sys/dev/block/%d:%d", maj, min)

	// 读取符号链接指向的目标
	target, err := os.Readlink(sysPath)
	if err != nil {
		return "", fmt.Errorf("无法读取符号链接 %s: %w", sysPath, err)
	}

	// 目标路径格式通常为 "../../devices/.../block/xxxx"
	parts := strings.Split(filepath.Clean(target), "/")
	deviceName := parts[len(parts)-1]
	if deviceName == "" {
		return "", fmt.Errorf("无法从路径 %s 提取设备名", target)
	}

	// 组合为完整设备路径（如 /dev/xxxx）
	return filepath.Join("/dev", deviceName), nil
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

// umount wraps the umount command and detaches the device from the specified mountPoint.
func umount(ctx context.Context, mountPoint string) error {
	if out, err := exec.CommandContext(ctx, "umount", mountPoint).CombinedOutput(); err != nil {
		return fmt.Errorf("umount: %w, output: %q", err, out)
	}
	return nil
}

// IsMounted 检查设备或目录是否已挂载
// device: 设备路径（如 /dev/sdb1）或挂载点目录（如 /mnt/data）
func isMounted(device string) (bool, string, error) {
	// 打开/proc/mounts文件，该文件记录了当前所有挂载信息
	file, err := os.Open("/proc/mounts")
	if err != nil {
		return false, "", fmt.Errorf("cant open /proc/mounts: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// /proc/mounts每行格式: 设备 挂载点 文件系统类型 挂载选项 dump  fsck顺序
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue // 无效行，跳过
		}

		mountDevice := parts[0]
		mountPoint := parts[1]

		// 规范化路径（处理相对路径和符号链接）
		normalizedDevice, err := filepath.EvalSymlinks(device)
		if err != nil {
			normalizedDevice = device // 规范化失败时使用原始路径
		}

		// 检查是否匹配设备路径或挂载点
		if mountDevice == device || mountDevice == normalizedDevice ||
			mountPoint == device || mountPoint == normalizedDevice {
			return true, mountPoint, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, "", fmt.Errorf("scanner /proc/mounts failed: %w", err)
	}

	return false, "", nil
}
