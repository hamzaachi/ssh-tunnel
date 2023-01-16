package main

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const dbfile string = "channels.db"

type Recod struct {
	ID   int
	Name string
	Type string
	Port int
}

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

func (con *Service) RetrieveByID(id int) (*Recod, error) {
	r := Recod{}

	stm, err := con.db.Prepare("SELECT * FROM channels WHERE id = ?")
	if err != nil {
		return nil, err
	}

	defer stm.Close()

	if err = stm.QueryRow(id).Scan(&r.ID, &r.Name, &r.Type, &r.Port); err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}

	return &r, nil
}

func (con *Service) RetrieveByName(name string) ([]Recod, error) {

	stm, err := con.db.Prepare("SELECT * FROM channels WHERE name = ?")
	if err != nil {
		return nil, err
	}

	defer stm.Close()

	rows, err := stm.Query(name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	l := []Recod{}
	for rows.Next() {
		r := Recod{}
		err := rows.Scan(&r.ID, &r.Name, &r.Type, &r.Port)
		if err != nil {
			return nil, err
		}
		l = append(l, r)
	}

	return l, nil
}

func (con *Service) List() ([]Recod, error) {

	stm, err := con.db.Prepare("SELECT * FROM channels;")
	if err != nil {
		return nil, err
	}

	defer stm.Close()

	rows, err := stm.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	l := []Recod{}
	for rows.Next() {
		r := Recod{}
		err := rows.Scan(&r.ID, &r.Name, &r.Type, &r.Port)
		if err != nil {
			return nil, err
		}
		l = append(l, r)
	}

	return l, nil
}

func (con *Service) GetPorts() ([]int, error) {

	stm, err := con.db.Prepare("SELECT localport FROM channels;")
	if err != nil {
		return nil, err
	}

	defer stm.Close()

	rows, err := stm.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var l []int
	for rows.Next() {
		r := Recod{}
		err := rows.Scan(&r.Port)
		if err != nil {
			return nil, err
		}
		l = append(l, r.Port)
	}

	return l, nil
}
