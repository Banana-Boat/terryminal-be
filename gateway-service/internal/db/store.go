package db

import "database/sql"

type Store struct {
	db       *sql.DB // 外部包不可见
	*Queries         // 匿名字段，默认包含Queries中的所有字段
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}
