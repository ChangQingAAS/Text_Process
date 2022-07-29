package src

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// 一次性将全部数据读入内存
//func test() {
//	contentBytes, err := ioutil.ReadFile("./resources/kaifangX.txt")
//	if err != nil {
//		fmt.Println("读入失败", err)
//	}
//	contentStr := string(contentBytes)
//
//	// 逐行打印
//	lineStrs := strings.Split(contentStr, "\n\r")
//	for _, lineStr := range lineStrs {
//		fmt.Println(lineStr)
//	}
//}

// 基于磁盘和缓存的读取
func ReadData() {

	file, _ := os.Open("./resources/kaifangX.txt")
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		lineBytes, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		fmt.Println(string(lineBytes))
	}
}

/**
通过 ioutil.ReadFile(filename) API可以傻瓜式地一次性读入全部数据

然而考虑到数据总量在1G以上，一次性读入内存，不但增加内存压力，产生爆内存的风险，还会长时间阻塞，效率极低；

正确的做法是，使用缓冲区进行逐行读取
*/
