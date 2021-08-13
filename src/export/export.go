package main

import (
	"database/sql"
	_ "database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)
import "os"

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	fmt.Println("my Go")
	f, err := os.Open(dir + "\\dbconfig.txt")
	if err != nil {
		fmt.Println("read file fail", err)
	}
	defer f.Close()

	fd, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("read to fd fail", err)
	}
	//fmt.Println(string(fd))
	var ss = string(fd)
	var arr = strings.Split(ss, ",")
	//fmt.Println(arr[0])
	var url = arr[0] + ":" + arr[1] + "@(" + arr[2] + ")/" + arr[3]
	//fmt.Println(url)
	db, _ := sql.Open("mysql", url) // 设置连接数据库的参数
	defer db.Close()                //关闭数据库
	err1 := db.Ping()               //连接数据库
	if err1 != nil {
		fmt.Println("数据库连接失败")
		return
	}
	sql, err := os.Open(dir + "\\sql.txt")
	if err != nil {
		fmt.Println("read file fail", err)
	}
	defer sql.Close()
	ssql, err := ioutil.ReadAll(sql)
	var bomstart = "\uFEFF"
	var realsql string
	if strings.HasPrefix(string(ssql), bomstart) {
		realsql = strings.ReplaceAll(string(ssql), bomstart, "")
	} else {
		realsql = string(ssql)
	}

	if !strings.HasPrefix(realsql, "select") {
		fmt.Println("非查询sql,不允许执行")
		fmt.Println(string(ssql))
		return
	}
	rows, _ := db.Query(realsql) //获取所有数据

	columns, _ := rows.Columns()
	columnLength := len(columns)
	cache := make([]interface{}, columnLength) //临时存储每行数据
	for index, _ := range cache {              //为每一列初始化一个指针
		var a interface{}
		cache[index] = &a
	}
	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("Sheet1")

	for rows.Next() {
		_ = rows.Scan(cache...)
		row := sheet.AddRow()
		row.SetHeightCM(1) //设置每行的高度
		for _, data := range cache {
			cell := row.AddCell()
			cell.Value = Strval(*data.(*interface{}))
		}

	}
	_ = rows.Close()

	err = file.Save(dir + "\\file.xlsx")
	if err != nil {
		panic(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("导出文件成功,文件存放位置:" + dir)
}

func Strval(value interface{}) string {
	// interface 转 string
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}
