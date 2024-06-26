package session

import (
	"database/sql"
	"fmt"
	"geeorm/dialect"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var testDB *sql.DB
var TestDial, _ = dialect.GetDialect("sqlite3")

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open("sqlite3", "../gee.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		os.Exit(1)
	}
	if testDB == nil {
		fmt.Println("Failed to open database connection.")
		os.Exit(1)
	}

	code := m.Run()

	testDB.Close()
	os.Exit(code)
}

func TestSessionExec(t *testing.T) {
	if testDB == nil {
		t.Fatal("Database connection is nil")
	}
	dialect, ok := dialect.GetDialect("sqlite3")
	if !ok {
		t.Fatal("get dialect failed")
	}
	session := New(testDB, dialect)
	session.Raw("DROP TABLE IF EXISTS User;").Exec()
	session.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := session.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	affected, _ := result.RowsAffected()
	if affected != 2 {
		t.Error("wrong rows: ", affected)
	}
}

func TestQueryRows(t *testing.T) {
	dialect, ok := dialect.GetDialect("sqlite3")
	if !ok {
		t.Fatal("get dialect failed")
	}
	session := New(testDB, dialect)
	session.Raw("DROP TABLE IF EXISTS User;").Exec()
	session.Raw("CREATE TABLE User(Name text);").Exec()
	session.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	rows, _ := session.Raw("SELECT Name FROM User").QueryRows()
	var names []string
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			t.Fatal("Failed to scan row:", err)
		}
		names = append(names, name)
	}
	if len(names) != 2 {
		t.Error("Query rows wrong: ", len(names))
	}
}

func TestTable(t *testing.T) {
	session := New(testDB, TestDial)
	session.Model(&User{})
	session.DropTable()
	//session.CreateTable()
	if !session.HasTable() {
		t.Fatal("Failed to create table User")
	}
}

func TestSession_Model(t *testing.T) {
	session := New(testDB, TestDial)
	session.Model(&User{})
	session.DropTable()
	session.CreateTable()
	if !(session.RefTable().FieldNames[0] == "Name") || !(session.RefTable().FieldNames[1] == "Age") {
		t.Fatal("session table name not fit")
	}
}
