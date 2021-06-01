package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"log"
)

func errorHandleDemo() {
	Select()
}

type dbObj struct {
	db *sql.DB
}

type MySqlConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}

func (d *dbObj) Open() *sql.DB {
	var err error
	mysqlConfig := MySqlConfig{
		User:     "root",
		Password: "root",
		Host:     "127.0.0.1",
		Port:     3306,
		Database: "test",
	}
	d.db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		mysqlConfig.User,
		mysqlConfig.Password,
		mysqlConfig.Host,
		mysqlConfig.Port,
		mysqlConfig.Database,
	))
	if err != nil {
		panic(err)
	}
	return d.db
}
func (d *dbObj) Close() {
	d.Close()
}

type userInfo struct {
	id      int
	orgcode string
	name    string
	version int
}

func SelectAll() {
	dbc := &dbObj{}
	db := dbc.Open()
	defer db.Close()

	stmt, _ := db.Prepare("SELECT orgcode,`name` FROM  userinfo WHERE id > ?")
	rows, _ := stmt.Query(0) //query为多行
	defer rows.Close()
	user := &userInfo{}

	for rows.Next() {
		err := rows.Scan(&user.orgcode, &user.name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(user.orgcode, ":", user.name)
	}
}

func Select() {
	dbc := &dbObj{}
	db := dbc.Open()
	defer db.Close()

	stmt, err := db.Prepare("SELECT orgcode,`name` FROM  userinfo WHERE ID= ?")
	if err != nil {
		log.Fatal("select failed!", err)
	}
	rows := stmt.QueryRow(1008) //QueryRow为单行
	user := &userInfo{}

	err = rows.Scan(&user.orgcode, &user.name)
	if err != nil {
		//if err == sql.ErrNoRows {
		//
		//}
		//%+v能把堆栈信息打印出来
		log.Fatalf("sql error:%+v", errors.Wrap(err, "No chaxundaojilu"))
	}
	fmt.Println(user.orgcode, ":", user.name)

}

func Insert() {
	dbc := &dbObj{}
	db := dbc.Open()
	defer db.Close()

	result, err := db.Exec("INSERT  userinfo (orgcode,imei,`name`) VALUE(?,?,?)", "cccc", 1009, "cccc")
	if err != nil {
		log.Fatal(err)
	}
	rowsaffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("获取受影响行数失败,err:%v", err)
		return
	}
	fmt.Println("受影响行数:", rowsaffected)
}

func Delete() {
	dbc := &dbObj{}
	db := dbc.Open()
	defer db.Close()
	result, err := db.Exec("DELETE FROM userinfo WHERE id=?", 1009)
	if err != nil {
		log.Fatal(err)
	}
	rowsaffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("获取受影响行数失败,err:%v", err)
		return
	}
	fmt.Println("受影响行数:", rowsaffected)
}

func Update() {
	dbc := &dbObj{}
	db := dbc.Open()
	defer db.Close()
	result, err := db.Exec("UPDATE userinfo SET `name`= ? WHERE id=?", "lcbbb", 1008)
	if err != nil {
		log.Fatal(err)
	}
	rowsaffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("获取受影响行数失败,err:%v", err)
		return
	}
	fmt.Println("受影响行数:", rowsaffected)
}

func Transaction() {
	dbc := &dbObj{}
	db := dbc.Open()
	defer db.Close()
	tx, _ := db.Begin()
	tx.Exec("UPDATE userinfo SET `name`= ? WHERE id=?", "lcaaa", 1007)
	result, err := tx.Exec("UPDATE userinfo SET `name`= ? WHERE id=?", "lcbbb", 1008)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit() //提交事务
	rowsaffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("获取受影响行数失败,err:%v", err)
		return
	}
	fmt.Println("受影响行数:", rowsaffected)
}

func ConcurrenceUpdate() {
	dbc := &dbObj{}
	db := dbc.Open()
	defer db.Close()
	getone := func(db *sql.DB, id int) *userInfo {
		stmt, _ := db.Prepare("SELECT orgcode,`name`,version FROM  userinfo WHERE ID= ?")
		rows := stmt.QueryRow(id) //
		user := &userInfo{}
		err := rows.Scan(&user.orgcode, &user.name, &user.version)
		if err != nil {
			log.Fatal(err)
		}
		return user
	}
	udateone := func(db *sql.DB, name string, id int, version int) {
		result, err := db.Exec("UPDATE userinfo SET `name`= ?, version=version+1 WHERE id=? AND version=?", name, id, version)
		if err != nil {
			log.Fatal(err)
		}
		rowsaffected, err := result.RowsAffected()
		if err != nil {
			fmt.Printf("并发更新获取受影响行数失败,err:%v", err)
			return
		}
		fmt.Println("并发更新受影响行数:", rowsaffected)
	}
	num := 10
	for i := 0; i < num; i++ {
		go func() {
			u := getone(db, 1008)
			fmt.Printf("获取数据:%v\r\n", u)
			udateone(db, "lc并发更新测试", 1008, u.version)
		}()
	}
	select {}
}
