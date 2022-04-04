package tagmd

import (
	"fmt"

	COMM "github.com/GGroups/rttm_login/comm"
	"github.com/jmoiron/sqlx"
	_ "github.com/logoove/sqlite"
)

type TagM struct {
	Id     int    `json:"tid" db:"tid"`
	Name   string `json:"name" db:"name"`
	Type   string `json:"type" db:"type"`
	Erasev int    `json:"erasev" db:"erasev"`
	Timev  int    `json:"timev" db:"timev"`
	Descv  int    `json:"descv" db:"descv"`
	Desc   string `json:"desc" db:"Desc"`
}

const (
	SQLCRE_TAG = `CREATE TABLE TagM ("tid" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	"name" TEXT NOT NULL, 
	"type" INTEGER NOT NULL, 
	"erasev" INTEGER NOT NULL,
	"timev" INTEGER NOT NULL,
	"descv" INTEGER NOT NULL, 
	"Desc");`

	SQLSEL_TAG = `SELECT tid, name, type, erasev, timev, descv,Desc FROM TagM; `

	SQLINS_TAG = `insert into TagM (
		"name", "type", "erasev", "timev", "descv", "Descv") 
	 VALUES (?,?,?,?,?,?)`

	SQLUPD_TAG = `update TagM set name=?, type=?, erasev=?, timev=?, descv=?, Desc=? where tid=?`
)

func InitTag() error {
	db, err := sqlx.Open(COMM.LITE3, COMM.DB_FILE)
	if err != nil {
		return err
	}
	_, err = db.Exec(SQLCRE_TAG)
	if err != nil {
		return err
	}
	db.Close()
	return nil
}

func CreatTag(tg *TagM) error {
	db, err := sqlx.Open(COMM.LITE3, COMM.DB_FILE)
	if err != nil {
		return err
	}
	u := *tg
	_, err = db.Exec(SQLINS_TAG, u.Name, u.Type, u.Erasev, u.Timev, u.Descv, u.Desc)
	if err != nil {
		return err
	}
	db.Close()
	return nil
}

func SetTag(tg *TagM) error {
	db, err := sqlx.Open(COMM.LITE3, COMM.DB_FILE)
	if err != nil {
		return err
	}
	u := *tg
	stmt, err := db.Prepare(SQLUPD_TAG)
	if err != nil {
		fmt.Println("##db.Prepare", err)
		return err
	}
	_, err = stmt.Exec(u.Name, u.Type, u.Erasev, u.Timev, u.Descv, u.Desc, u.Id)
	if err != nil {
		fmt.Println("##db.Prepare", err)
		return err
	}
	db.Close()
	return nil
}

func GetAllTags(tgs *[]TagM) error {
	db, err := sqlx.Open(COMM.LITE3, COMM.DB_FILE)
	if err != nil {
		return err
	}
	err = db.Select(tgs, SQLSEL_TAG)
	if err != nil {
		return err
	}
	db.Close()
	return nil
}
