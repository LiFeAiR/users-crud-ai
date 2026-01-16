package utils

import "database/sql"

func NewNullInt32(i int32) sql.NullInt32 {
	return sql.NullInt32{Int32: i, Valid: true}
}

func NewNullInt32Nil() sql.NullInt32 {
	return sql.NullInt32{Valid: false}
}
