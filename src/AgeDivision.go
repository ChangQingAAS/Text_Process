package src

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

/*
主协程负责逐行读取文本大数据

另外开设34条子协程，负责对不同省份文件进行写入，从34个不同的管道中扫描数据并写出文件；

主协程根据身份证号反映的不同省份，将读入的信息丢入不同的管道，由对应的子协程进行文件写出；

当文件读取完毕时，关闭所有的数据管道（通知子协程停止数据扫描）；

主协程通过等待组等待所有子协程完成任务；
*/
type Ager struct {
	decade   string // 190X, ... 200x,201x
	file     *os.File
	chanData chan string
}

var Awg sync.WaitGroup

func AwriteFile(ager *Ager) {
	// 死循环读取管道，管道关闭时循环结束
	for lineStr := range ager.chanData {
		ager.file.WriteString(lineStr)
		fmt.Print(ager.decade, "x 写入", lineStr)
	}
	// 标记协程结束
	Awg.Done()

}
func AgeDivision() {
	// create many Ager object
	agersMap := make(map[string]*Ager)
	for i := 190; i < 202; i++ {
		ager := Ager{decade: strconv.Itoa(i)}
		file, _ := os.OpenFile("./resources/age/"+ager.decade+"x.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		ager.file = file
		defer ager.file.Close()

		// 创建每个省的数据管道
		ager.chanData = make(chan string, 0)
		agersMap[ager.decade] = &ager
		fmt.Println("管道已经创建")
	}

	// set a routine for every decade
	for _, ager := range agersMap {
		Awg.Add(1)
		go AwriteFile(ager)
	}

	// read good data
	file, _ := os.Open("./resources/kaifang_good.txt")
	defer file.Close()
	reader := bufio.NewReader(file)

	// determine the first two digits of the ID number line by line
	for {
		lineStr, err := reader.ReadString('\n')
		if err == io.EOF {
			for _, ager := range agersMap {
				close(ager.chanData)
				fmt.Println("管道已关闭")
			}
			break
		}

		// 拿出decade
		decade := strings.Split(lineStr, "，")[1][6:9]

		// 对号入座，写入相应的管道
		if agersMap[decade] != nil {
			agersMap[decade].chanData <- lineStr + "\n"
		} else {
			fmt.Println("\n\n\n\n", lineStr, "\n\n\n\n")
		}
	}

	// blocking if the wait routine stop
	Awg.Wait()
}
