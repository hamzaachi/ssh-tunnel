package main

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const dbfile string = "channels.db"

func connect(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbfile)

	if err != nil {
		return nil, err
	}

	return db, nil

}

func (con *Service) Insert() error {

	stm, err := con.db.Prepare("INSERT INTO channels(name, type, localport) VALUES(?,?,?)")
	if err != nil {
		return err
	}

	defer stm.Close()

	_, err = stm.Exec(con.Name, con.Category, con.LocalPort)

	if err != nil {
		return err
	}
	return nil
}
