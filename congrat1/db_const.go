package congrat1

import (
	"database/sql"
	"fmt"
)

const (
	DriverName     = "mysql"
	DataSourceName = "root:liuyilun134@tcp(127.0.0.1:3306)/nts?charset=utf8&parseTime=True&loc=Local"
)

func UseDBConnection(f func(db *sql.DB) error) {
	db, err := sql.Open(DriverName, DataSourceName)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = db.Close()
	}()
	err = f(db)
	if err != nil {
		fmt.Println(err)
	}
}
