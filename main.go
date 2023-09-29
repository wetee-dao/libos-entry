package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/afero"
	"golang.org/x/sys/unix"
	"wetee.app/libos-entrypoint/libos"
	"wetee.app/libos-entrypoint/utils"
)

func main() {
	log.SetPrefix("[LibOS Entrypoint] ")

	// Get libOS based on uname
	// 获取 libOS 类型
	libOS, err := utils.GetLibOS()
	if err != nil {
		utils.ExitWithMsg("Failed to get libOS: %s", err)
	}

	// Use filesystem from libOS
	// 获取libOS的文件系统
	hostfs := afero.NewOsFs()
	fmt.Println("OS hostfs: ", hostfs)

	var service string

	switch libOS {
	case "Gramine":
		log.Println("Geted libOS: Gramine")

		service, err = libos.InitGramineEntry(hostfs)
		if err != nil {
			utils.ExitWithMsg("activating Gramine entry failed: %s", err)
		}

		// case "Occlum":
		// 	log.Println("Geted libOS: Occlum")

		// 	service, err =  initOcclumEntry(hostfs)
		// 	if err != nil {
		// 		exit("activating Occlum entry failed: %s", err)
		// 	}
	}

	// Start service
	// 开启服务
	if err := unix.Exec(service, os.Args, os.Environ()); err != nil {
		utils.ExitWithMsg("%s", err)
	}
}
