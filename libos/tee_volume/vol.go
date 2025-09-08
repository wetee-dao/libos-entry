package volume

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetDeviceNameByMajMin 根据主设备号(MAJ)和次设备号(MIN)获取设备路径（如 /dev/sda）
func GetDeviceNameByMajMin(maj, min int) (string, error) {
	// /sys/dev/block/<MAJ:MIN> 是指向设备实际路径的符号链接
	sysPath := fmt.Sprintf("/sys/dev/block/%d:%d", maj, min)

	// 读取符号链接指向的目标
	target, err := os.Readlink(sysPath)
	if err != nil {
		return "", fmt.Errorf("无法读取符号链接 %s: %w", sysPath, err)
	}

	// 目标路径格式通常为 "../../devices/.../block/sda"
	// 提取最后一个路径段（如 "sda"）
	parts := strings.Split(filepath.Clean(target), "/")
	deviceName := parts[len(parts)-1]
	if deviceName == "" {
		return "", fmt.Errorf("无法从路径 %s 提取设备名", target)
	}

	// 组合为完整设备路径（如 /dev/sda）
	return filepath.Join("/dev", deviceName), nil
}
