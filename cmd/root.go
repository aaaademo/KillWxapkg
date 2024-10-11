package cmd

import (
	"log"
	"sync"

	. "github.com/Ackites/KillWxapkg/internal/cmd"
	. "github.com/Ackites/KillWxapkg/internal/config"
	"github.com/Ackites/KillWxapkg/internal/restore"
)

func Execute(appID, input, outputDir, fileExt string, restoreDir bool, pretty bool, noClean bool, save bool, sensitive bool) {
	log.Println("开始处理AppID：", appID)
	// 存储配置
	configManager := NewSharedConfigManager()
	configManager.Set("appID", appID)
	configManager.Set("input", input)
	configManager.Set("outputDir", outputDir)
	configManager.Set("fileExt", fileExt)
	configManager.Set("restoreDir", restoreDir)
	configManager.Set("pretty", pretty)
	configManager.Set("noClean", noClean)
	configManager.Set("save", save)
	configManager.Set("sensitive", sensitive)

	if input == "" {
		// input 等于 默认值 加上 appID
		input = SetDefaultInput(appID)
		// 提示 设置 默认路径发生错误，并结束程序
		if input == "" {
			log.Println("未设置输入路径，请设置后重试")
			return
		}
	}


	inputFiles := ParseInput(input, fileExt)

	if len(inputFiles) == 0 {
		log.Println("未找到任何文件")
		return
	}

	// 确定输出目录
	if outputDir == "" {
		// outputDir = DetermineOutputDir(input, appID)
		outputDir = "outputs/"+ appID
	}

	var wg sync.WaitGroup
	for _, inputFile := range inputFiles {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			err := ProcessFile(file, outputDir, appID, save)
			if err != nil {
				log.Printf("处理文件 %s 时出错: %v\n", file, err)
			} else {
				log.Printf("成功处理文件: %s\n", file)
			}
		}(inputFile)
	}
	wg.Wait()

	// 还原工程目录结构
	restore.ProjectStructure(outputDir, restoreDir)
}
