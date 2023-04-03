package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DbHdr *sql.DB

func InitDB() (err error) {
	fmt.Println("==========start to initial DB handler======")
	//組合sql連線字串
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", UserName, Password, Addr, Port, Database)
	dsn := "root:iloveaoe@tcp(127.0.0.1:3306)/user"
	DbHdr, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	err = DbHdr.Ping()
	if err != nil {
		return err
	}
	return nil
}
