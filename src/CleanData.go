package src

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

/*
一些数据条目的第二个字段并不为身份证号，而是一些邮政编码或者座机号，对我们用处不大，所以我们先来做一些清洗。

把好的数据放在一块，坏的数据也不要轻易扔掉，也总结在一起，以后万一用其他的挖掘方式还可能用到
*/

// 将有带身份证号的数据整理到一起
func CleanData() {
	// 基于磁盘和缓存的读取
	file, _ := os.Open("./resources/kaifangX.txt")
	defer file.Close()

	// 准备一个优质文件
	goodFile, _ := os.OpenFile("./resources/kaifang_good.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	defer goodFile.Close()

	//准备一个劣质文件
	badFile, _ := os.OpenFile("./resources/kaifang_bad.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	defer badFile.Close()

	reader := bufio.NewReader(file)
	for {
		lineBytes, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		lineStr := string(lineBytes)

		fmt.Println(string(lineStr))

		fields := strings.Split(lineStr, "，")
		if len(fields) > 1 && len(fields[1]) == 18 {
			// 防止有空白行
			//身份证号18位，这里偷懒了，仔细的话应该用正则判断
			goodFile.WriteString(lineStr + "\n")
			fmt.Println("Good: ", lineStr)
		} else {
			badFile.WriteString(lineStr + "\n")
			fmt.Println("Bad: ", lineStr)
		}
	}
}
