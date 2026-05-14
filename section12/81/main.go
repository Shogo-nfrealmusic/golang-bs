package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DbConnection *sql.DB

type Person struct {
	Name string
	Age int
}

func main() {
	DbConnection, _ = sql.Open("sqlite3", "./example.sql")
	defer DbConnection.Close()
	cmd := `CREATE TABLE IF NOT EXISTS person(
		name STRING,
		age INT)`
	_, err := DbConnection.Exec(cmd)
	if err != nil {
		log.Fatal(err)
	}

	cmd = "INSERT INTO person (name, age) VALUES (?, ?)"
	_, err = DbConnection.Exec(cmd, "taro", 20)
	if err != nil {
		log.Fatal(err)
	}

	// cmd = "UPDATE person SET age = ? WHERE name = ?"
	// _, err = DbConnection.Exec(cmd, 25, "taro")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// cmd = "SELECT * FROM person"
	// rows, _ := DbConnection.Query(cmd)
	// defer rows.Close()
	// var pp []Person
	// for rows.Next() {
	// 	var p Person
	// 	err := rows.Scan(&p.Name, &p.Age)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	pp = append(pp, p)
	// }
	// for _, p := range pp {
	// 	fmt.Println(p.Name, p.Age)
	// }

	cmd = "DELETE FROM person WHERE name = ?"
	_, err = DbConnection.Exec(cmd, "taro")
	if err != nil {
		log.Fatal(err)
	}
}