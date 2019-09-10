package spider

// 任务管理

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// 默认爬取数量
const maxNum = 120

// 获取任务列表
func GetJobs() map[string]int {
	// 任务文件目录
	pwd, _ := os.Getwd()
	fileDir := pwd + "/static/job"

	// 解析目录下的txt文件，生成任务列表
	jobs := map[string]int{}
	files, _ := ioutil.ReadDir(fileDir)
	for _, f := range files {
		// 判断文件名是否合法
		fileName := f.Name()
		if match, _ := regexp.MatchString("-job.txt", fileName); !match {
			continue
		}
		// open file
		file, err := os.Open(fileDir + "/" + fileName)
		if err != nil {
			fmt.Println(fileName, "打开失败", err)
			continue
		}
		defer file.Close()
		// 按行获取文件内容
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if line != "" {
				// 根据逗号分隔，判断是否有默认爬取数量，如无设置为maxNum
				splitArr := strings.Split(line, ",")
				if len(splitArr) >= 2 {
					num, _ := strconv.Atoi(splitArr[1])
					jobs[splitArr[0]] = num
				} else {
					jobs[splitArr[0]] = maxNum
				}
			}
		}
	}
	return jobs
}
