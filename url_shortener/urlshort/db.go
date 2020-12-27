package urlshort

import (
	"database/sql"
	"fmt"
)

// Conn DB connection struct
type Conn struct {
	DB *sql.DB
}

// NewDB returns DB connection struct
func NewDB() (Conn, error) {
	conn, err := newDB()
	if err != nil {
		return Conn{}, err
	}
	if err := conn.init(); err != nil {
		return Conn{}, err
	}
	return conn, nil
}

func newDB() (Conn, error) {
	db, err := sql.Open("sqlite3", "./path.db")
	if err != nil {
		return Conn{}, err
	}
	return Conn{DB: db}, nil
}

func (conn *Conn) init() error {
	createTable := `CREATE TABLE IF NOT EXISTS paths (
		path varchar(20) not null,
		url varchar(100) not null,
		PRIMARY KEY(path)
	);`
	_, err := conn.DB.Exec(createTable)
	if err != nil {
		return fmt.Errorf("db create Table error: %v", err)
	}

	insertData, err := conn.DB.Prepare("INSERT INTO paths (path, url) VALUES (?, ?)")
	if err != nil {
		return err
	}

	insertData.Exec("/urlshort-db", "https://github.com/gophercises/urlshort")
	insertData.Exec("/urlshort-final-db", "https://github.com/gophercises/urlshort/tree/solution")
	return nil
}

// GetPathsToURLsByDB create and initilize DB, then return
// Paths and URLs mapping info.
func (conn *Conn) GetPathsToURLsByDB() ([]PathToURL, error) {
	rows, err := conn.DB.Query("SELECT * FROM paths")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var PathsToURLs []PathToURL
	for rows.Next() {
		p := PathToURL{}
		if err := rows.Scan(&p.Path, &p.URL); err != nil {
			return nil, err
		}
		PathsToURLs = append(PathsToURLs, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return PathsToURLs, nil
}
