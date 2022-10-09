package mydbs

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Mydata struct {
	Myip    string
	Mytime  int
	Myvalue string
}

var (
	create_table = `create table if not exists mydata (myip text,mytime integer,myvalue text)`
	dbname       = "data/data.db"
	dbtype       = "sqlite3"
	sqldb        *sql.DB
)

func initDir(filepath string) {
	_, err := os.Stat(filepath)
	if err == nil {
		//文件存在
		return
	}
	if os.IsNotExist(err) {
		//文件夹不存在
		err := os.Mkdir(filepath, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func init() {
	//初始化文件夹
	initDir("data")
	//初始化数据库
	var err error
	sqldb, err = sql.Open(dbtype, dbname)
	checkErr(err)
	TableCreate()
}

func checkErr(err error) {
	if err != nil {
		panic("sql error")
	}
}

func TableCreate() {
	sqldb.Exec(create_table)
}

func DataInsert(data *Mydata) {
	insert_into_sql := "insert into mydata (myip,mytime,myvalue) values(?,?,?)"
	sqldb.Exec(insert_into_sql, data.Myip, data.Mytime, data.Myvalue)
}

func DataUpdate(data *Mydata) {
	update_sql := "update mydata set mytime=?,myvalue=? where myip=?"
	sqldb.Exec(update_sql, data.Mytime, data.Myvalue, data.Myip)
}

func DataSelect(myip string) (Mydata, error) {
	var m Mydata
	select_ip_sql := "select * from mydata where myip=?"
	rows := sqldb.QueryRow(select_ip_sql, myip)
	err := rows.Scan(&m.Myip, &m.Mytime, &m.Myvalue)
	return m, err
}
