package statement

import (
	"database/sql"
	"log"
)

func InsertTable(db *sql.DB) {
	stmIns, err := db.Prepare("INSERT INTO pet VALUES (?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmIns.Exec("Claws", "Gwen", "cat", "m", "1994-03-17", nil)
	if err != nil {
		panic(err.Error())
	}
	_, err = stmIns.Exec("Buffy", "Harold", "dog", "f", "1989-05-13", nil)
	if err != nil {
		panic(err.Error())
	}
	_, err = stmIns.Exec("Fang", "Benny", "dog", "m", "1990-08-27", nil)
	if err != nil {
		panic(err.Error())
	}
	_, err = stmIns.Exec("Bowser", "Diane", "dog", "m", "1979-08-31", "1995-07-29")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmIns.Exec("Chirpy", "Gwen", "bird", "f", "1998-09-11", nil)
	if err != nil {
		panic(err.Error())
	}
	_, err = stmIns.Exec("Whistler", "Gwen", "bird", nil, "1997-12-09", nil)
	if err != nil {
		panic(err.Error())
	}
	_, err = stmIns.Exec("Slim", "Benny", "snake", "m", "1996-04-29", nil)
	if err != nil {
		panic(err.Error())
	}
}

func QueryUser(db *sql.DB) {
	var (
		id   string
		name string
	)
	rows, err := db.Query("select id,name from user where id = ?", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
