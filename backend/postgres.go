package backend

import (
	"database/sql"

	"github.com/chadac/vfs"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	connStr string
	tableName string
}

func (s *PostgresStore) Connect() (*sql.DB, error) {
}

func (s *PostgresStore) Get(id Checksum) (*Object, error) {
	db, err := s.Connect()
	rows, err := db.Query("SELECT object FROM %1 WHERE id = %2", s.tableName, id)
	if err != nil {
		return nil, err
	}
	for rows.Next()
}
