package data

import "database/sql"

type Models struct {
}

func NewModels(conn *sql.DB) Models {
	return Models{}
}
