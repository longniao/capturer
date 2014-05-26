package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/msbranco/goconfig"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// 抓取成功后更新状态
func update_status(db *sql.DB, user_id int) {
	stmt, err := db.Prepare("update table set status = 1 where record_id = ? ")
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return
	}

	stmt.Exec(user_id)
}

// 抓取文件到本地
func download_handler(webroot, remote_url string, call_date string, call_sheet_id string) error {

	//log.Println(remote_url)
	resp, err := http.Get(remote_url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 判断并创建文件夹
	call_date = webroot + call_date
	createDir(call_date)

	file_name := call_date + "/" + call_sheet_id + ".mp3"
	err = ioutil.WriteFile(file_name, body, 0777)
	if err != nil {
		panic(err)
	}

	return err
}

// 判断并创建文件夹
func createDir(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {

		e := os.Mkdir(path, 0777)
		if e != nil {
			panic(err)
		}

		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}

}

// 执行
func run(db *sql.DB, webroot string) {

	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// 查询
	rows, err := db.Query("select record_id, call_sheet_id, remote_url, call_date from table where status = ? limit 1", 0)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var record_id int
	var call_sheet_id string
	var remote_url string
	var call_date string
	for rows.Next() {
		err := rows.Scan(&record_id, &call_sheet_id, &remote_url, &call_date)
		if err != nil {
			log.Fatal(err)
		}

		err = download_handler(webroot, remote_url, call_date, call_sheet_id)
		if err != nil {
			log.Fatal(err)
		} else {
			update_status(db, record_id)
			log.Println(record_id, call_date, remote_url)
		}

	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	// 配置
	c, err := goconfig.ReadConfigFile("config.ini")
	if err != nil {
		print(err)
	}

	config, err := c.GetString("serverconfig", "config")

	dbhost, err := c.GetString(config, "dbhost")
	dbport, err := c.GetString(config, "dbport")
	dbuser, err := c.GetString(config, "dbuser")
	dbpassword, err := c.GetString(config, "dbpassword")
	webroot, err := c.GetString(config, "webroot")

	// 连接数据库
	db, err := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/flounder?autocommit=true")

	if err != nil {
		log.Fatalf("Open database error: %s\n", err)
	}
	defer db.Close()

	for {
		run(db, webroot)
		log.Println(time.Now())
		time.Sleep(1000000000)
	}
}
