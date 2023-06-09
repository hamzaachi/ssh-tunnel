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

func (con *Tunnel) Insert(ctx context.Context) error {

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

func (con *Tunnel) Delete(ctx context.Context, Name, Type string) error {

	stm, err := con.db.Prepare("DELETE FROM channels WHERE name = ? AND type = ?")
	if err != nil {
		return err
	}

	defer stm.Close()

	_, err = stm.Exec(Name, Type)

	if err != nil {
		return err
	}

	return nil
}

func (con *Tunnel) RetrieveByID(ctx context.Context, Name, Type string) ([]Recod, error) {
	r := Recod{}

	stm, err := con.db.Prepare("SELECT * FROM channels WHERE name = ? AND type = ?")
	if err != nil {
		return nil, err
	}

	defer stm.Close()

	l := []Recod{}
	if err = stm.QueryRow(Name, Type).Scan(&r.ID, &r.Name, &r.Type, &r.Port); err == sql.ErrNoRows {
		return l, nil
	}

	l = append(l, r)
	return l, nil
}

func (con *Tunnel) RetrieveByName(ctx context.Context, name string) ([]Recod, error) {

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

func (con *Tunnel) List(ctx context.Context) ([]Recod, error) {

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

func (con *Tunnel) GetPorts(ctx context.Context) ([]int, error) {

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
