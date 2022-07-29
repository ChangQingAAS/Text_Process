package main

import (
	"bufio"
	"fmt"
	"github.com/ChangQingAAS/GoDemo/TextBisDataProcess/src"
	"io"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const CACHE_LEN = 2

var (
	kfMap map[string]src.TimeData
	// chanSema chan int // 信号量管道
	chanData chan *src.KfPerson
	db       *sqlx.DB
)

func HandleError(err error, why string) {
	if err != nil {
		fmt.Println("ERROR OCCUR!!! ", err, why)
	}
}

func init() {
	_, err1 := os.Stat("./resources/kaifang_good_dbok.mark")
	if err1 == nil {
		fmt.Println("database init finish!")
		return
	}

	var err error
	db, err := sqlx.Open("mysql", "root:root@tcp(127.0.0.1:3306)/bigText")
	HandleError(err, "sqlx.Open")
	defer db.Close()
	fmt.Println("数据库已经打开")

	// create table if not exists
	_, err = db.Exec("create table if not exists kfperson(id int primary key auto_increment,name,varchar(20),idcard char(18),sex char(1));")
	HandleError(err, "db.Exec create table")
	fmt.Println("table has created")

	// 初始化信号量管道(控制并发数
	chanData = make(chan *src.KfPerson, 1000000) // 假定当前通道的缓存能力是无限的
	fmt.Println("管道已经初始化")

	// 开辟携程，源源不断地从数据管道获取信息，插入数据库
	for i := 0; i < 500; i++ {
		go insertKfPerson()
		fmt.Printf("协程%d已经开辟\n", i)
	}

	// open source file
	file, e := os.Open("./resources/kaifang_good.txt")
	HandleError(e, "os.Open")
	defer file.Close()
	reader := bufio.NewReader(file)
	fmt.Println("大数据文本已经打开")

	// 分批次读入大数据文本
	// 还是要基于缓存的读取
	for {
		lineBytes, _, err := reader.ReadLine()

		if err == io.EOF {
			// close the channel of data
			close(chanData)
			break
		}

		HandleError(err, "reader.ReadLine")

		lineStr := string(lineBytes)
		fields := strings.Split(lineStr, "，")
		name, idcard := fields[0], fields[1]

		// 抛弃过长的名字
		name = strings.TrimSpace(name)
		if len(strings.Split(name, "")) > 20 {
			fmt.Printf("%s 名字过长，已经抛弃了")
			continue
		}

		// 开有限条协程，从管道中读取数据
		kfPerson := src.KfPerson{Name: name, IdCard: idcard}
		chanData <- &kfPerson
	}
	fmt.Println("init success!!!")

	_, err = os.Create("./resources/kaifang_good_dbok.mark")
	if err == nil {
		fmt.Println("the signal of database initiation finish!")
	}
}

func insertKfPerson() {
	for kfPerson := range chanData {
		// 循环插入直至成功
		for {
			result, err := db.Exec("insert into kfperson(name,idcard) value(?,?);", kfPerson.Name, kfPerson.IdCard)
			HandleError(err, "db.Exec insert")
			if err != nil {
				// db性能被消耗尽了，休息五秒
				<-time.After(5 * time.Second)
			} else {
				if n, e := result.RowsAffected(); e == nil && n > 0 {
					fmt.Printf("insert %s success!\n", kfPerson.Name)
					break
				}
			}
		}
	}
}

func main() {
	// open database
	db, err := sqlx.Open("mysql", "root:root@tcp(127.0.0.1:3306)/bigText")
	HandleError(err, "sqlx.Open")
	defer db.Close()

	// init cache
	kfMap = make(map[string]src.TimeData, 0)

	// for loop to query
	var name string
	for {
		fmt.Print("please input the name of the kaifang： ")
		fmt.Scanf("%s", &name)

		// exit the program
		if name == "exit" {
			break
		}

		// check all cache
		if name == "cache" {
			fmt.Printf("共缓存了%d条结果: \n", len(kfMap))
			for key := range kfMap {
				fmt.Println(key)
			}
			continue
		}

		// 先查看内存中是否有结果
		if td, ok := kfMap[name]; ok {
			qr := td.(*src.QueryResult)
			fmt.Println(qr)
			qr.Count += 1
			fmt.Println(qr.Value)
			fmt.Println("共查询到%d条结果", len(qr.Value))
			continue
		}

		// cache中没有，goto database
		kfpeople := make([]src.KfPerson, 0)
		query := "select id,name,idcard from kfperson where name like '" + name + "' ;"
		e := db.Select(&kfpeople, query)
		HandleError(e, "db.Select")
		fmt.Printf("共查询到%d条结果： \n", len(kfpeople))
		fmt.Println(kfpeople)

		// add query result into cache
		queryResult := src.QueryResult{Value: kfpeople}
		queryResult.CacheTime = time.Now().UnixNano()
		queryResult.Count = 1
		kfMap[name] = &queryResult

		// 淘汰存放时间过长的缓存
		if len(kfMap) > CACHE_LEN {
			delKey := src.UpdateCache(&kfMap)
			fmt.Printf("%s已经被淘汰出缓存!\n", delKey)
		}
	}
	fmt.Println("ALL OVER!")
}
