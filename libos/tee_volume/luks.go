package volume

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

// CryptLuks 表示一个加密设备，用于LUKS加密管理
type CryptLuks struct {
	devPath   string // 物理设备路径
	metaPath  string // 头部信息备份文件路径
	keyFile   string // 密钥文件路径
	mappingID string // 设备映射名称
}

// NewCryptLuks 创建CryptLuks实例（不与实际设备交互）
func NewCryptLuks(devPath, keyFile, mappingID string) (*CryptLuks, error) {
	var randSuffix [4]byte
	if _, err := rand.Read(randSuffix[:]); err != nil {
		return nil, errors.Wrap(err, "生成头部文件随机后缀失败")
	}
	return &CryptLuks{
		devPath:   devPath,
		keyFile:   keyFile,
		mappingID: mappingID,
	}, nil
}

// CheckFormat 检查设备是否已格式化为LUKS
func (c *CryptLuks) CheckFormat(ctx context.Context) (bool, error) {
	var exitErr *exec.ExitError
	output, err := execCryptCommand(ctx, "isLuks", "--verbose", c.devPath)
	if errors.As(err, &exitErr) {
		if exitErr.ExitCode() == 1 {
			// 退出码1表示不是LUKS设备
			return false, nil
		}
		return false, fmt.Errorf("检查设备 %s 是否为LUKS格式失败: %w, 输出: %s, 错误: %s",
			c.devPath, err, output, exitErr.Stderr)
	} else if err != nil {
		return false, fmt.Errorf("检查LUKS格式时出错: %w", err)
	}
	return true, nil
}

// Format 格式化设备为LUKS2格式
func (c *CryptLuks) Format(ctx context.Context) error {
	// 创建头部文件所在目录
	if err := os.MkdirAll(filepath.Dir(c.metaPath), 0o700); err != nil {
		return fmt.Errorf("create meta dir %s 失败: %w", filepath.Dir(c.metaPath), err)
	}

	// 构建cryptsetup命令参数
	args := []string{
		"luksFormat",
		"--type=luks2",                  // 使用LUKS2格式
		"--cipher=aes-xts-plain64",      // 加密算法
		"--pbkdf=argon2id",              // 密钥派生函数
		"--pbkdf-memory=10240",          // 内存使用限制(10MiB)
		"--integrity=hmac-sha256",       // 完整性校验算法
		"--integrity-no-wipe",           // 不擦除设备（未写入区块视为无效）
		"--sector-size=4096",            // 扇区大小4KiB
		"--batch-mode",                  // 非交互模式
		fmt.Sprintf("-d=%s", c.keyFile), // 密钥文件路径
		c.devPath,
	}

	_, err := execCryptCommand(ctx, args...)
	return err
}

// Attach 激活LUKS设备（打开映射）
func (c *CryptLuks) Attach(ctx context.Context) error {
	// 执行luksOpen命令
	args := []string{
		"luksOpen",
		fmt.Sprintf("-d=%s", c.keyFile),
		c.devPath,
		c.mappingID,
	}

	if _, err := execCryptCommand(ctx, args...); err != nil {
		return fmt.Errorf("激活设备映射 %s 失败: %w", c.mappingID, err)
	}
	return nil
}

// Detach 关闭LUKS设备映射
func (c *CryptLuks) Detach(ctx context.Context) error {
	_, err := execCryptCommand(ctx, "luksClose", c.mappingID)
	return err
}

// CheckExt4Format 检查设备是否为ext4格式
func (c *CryptLuks) CheckExt4Format(_ context.Context) (bool, error) {
	const (
		ext4Magic   = uint16(0xef53) // ext4魔术字
		magicOffset = 1080           // 魔术字在超级块中的偏移量
	)

	mappingPath := filepath.Join("/dev/mapper", c.mappingID)
	file, err := os.Open(mappingPath)
	if err != nil {
		return false, fmt.Errorf("打开映射设备失败: %w", err)
	}
	defer file.Close()

	// 定位到魔术字位置
	if _, err := file.Seek(magicOffset, 0); err != nil {
		return false, fmt.Errorf("定位魔术字偏移失败: %w", err)
	}

	// 读取并验证魔术字
	var magic uint16
	if err := binary.Read(file, binary.LittleEndian, &magic); err != nil {
		// 未格式化设备可能返回I/O错误
		if errors.Is(err, syscall.EIO) {
			return false, nil
		}
		return false, fmt.Errorf("读取魔术字失败: %w", err)
	}

	return magic == ext4Magic, nil
}

// CreateExt4 格式化设备为ext4文件系统
func (c *CryptLuks) FormatExt4(ctx context.Context) error {
	mappingPath := filepath.Join("/dev/mapper", c.mappingID)

	// 擦除 ext4 超级块备份区域
	if err := clearExt4Blocks(ctx, mappingPath); err != nil {
		return fmt.Errorf("clear ext4 error: %w", err)
	}

	// 执行格式化
	return formatExt4(ctx, mappingPath)
}

var numberRegex = regexp.MustCompile(`\d+`)

// 执行 cryptsetup 命令的通用函数
func execCryptCommand(ctx context.Context, args ...string) (string, error) {
	// 创建 cryptsetup 运行所需目录
	if err := os.MkdirAll("/run/cryptsetup", 0o755); err != nil {
		return "", errors.Wrap(err, "create cryptsetup run dir error")
	}

	fmt.Println("cryptsetup", args)
	cmd := exec.CommandContext(ctx, "cryptsetup", args...)
	output, err := cmd.Output()
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return "", errors.Wrap(err, "cryptsetup command error:"+string(exitErr.Stderr))
	} else if err != nil {
		return "", fmt.Errorf("cryptsetup command error: %w", err)
	}

	return string(output), nil
}

// 清理 ext4 文件系统相关区块
func clearExt4Blocks(ctx context.Context, devPath string) error {
	// 干运行获取超级块备份位置
	cmd := exec.CommandContext(ctx, "mkfs.ext4", "-F", "-n", devPath)
	output, err := cmd.Output()
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return fmt.Errorf("获取 ext4 区块信息失败: %w, 错误输出: %q", err, exitErr.Stderr)
	} else if err != nil {
		return fmt.Errorf("执行 mkfs.ext4 失败: %w, 输出: %q", err, output)
	}

	// 解析超级块备份位置
	delimiter := "Superblock backups stored on blocks:"
	_, blockListStr, ok := strings.Cut(string(output), delimiter)
	if !ok {
		_, blockListStr, ok = strings.Cut(string(output), "超级块的备份存储于下列块：")
		if !ok {
			return fmt.Errorf("parsing mkfs.ext4 output: delimiter %q not found in output %q", delimiter, output)
		}
	}

	// 提取区块编号
	blockStrs := numberRegex.FindAllString(blockListStr, -1)
	if len(blockStrs) == 0 {
		return fmt.Errorf("parsing mkfs.ext4 output: no block numbers found in output %q", output)
	}

	// 转换为整数并添加 0号区块
	blockNums := make([]int64, 0, len(blockStrs)+1)
	for _, s := range blockStrs {
		num, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("parsing mkfs.ext4 output: parsing block number %q: %w", s, err)
		}
		if num < 0 {
			return fmt.Errorf("parsing mkfs.ext4 output: invalid block number %d", num)
		}
		blockNums = append(blockNums, int64(num))
	}
	blockNums = append(blockNums, 0)

	// 清零指定区块
	if err := zeroDeviceBlocks(devPath, blockNums); err != nil {
		return fmt.Errorf("clearing ext4 blocks on device %s: %w", devPath, err)
	}
	return nil
}

// 直接写入零值到指定区块（绕过页缓存）
func zeroDeviceBlocks(devPath string, indices []int64) error {
	// 以直接IO模式打开设备
	fd, err := unix.Open(devPath, unix.O_WRONLY|unix.O_DIRECT, 0)
	if err != nil {
		return fmt.Errorf("打开设备失败: %w", err)
	}
	defer unix.Close(fd)

	const blockSize = 4096

	// 分配页对齐的零缓冲区（直接IO要求）
	buf, err := unix.Mmap(-1, 0, blockSize, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_ANONYMOUS|unix.MAP_PRIVATE)
	if err != nil {
		return fmt.Errorf("分配零缓冲区失败: %w", err)
	}
	defer func() { _ = unix.Munmap(buf) }()

	// 写入每个区块
	for _, index := range indices {
		offset := index * blockSize
		for written := 0; written < blockSize; {
			n, err := unix.Pwrite(fd, buf[written:], offset+int64(written))
			if err != nil {
				return fmt.Errorf("写入区块 %d 失败: %w", index, err)
			}
			written += n
		}
	}

	return nil
}

// 格式化设备为ext4文件系统
func formatExt4(ctx context.Context, devPath string) error {
	cmd := exec.CommandContext(ctx, "mkfs.ext4", devPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("执行 ext4 格式化失败: % w, 输出: % q", err, output)
	}
	return nil
}
