package db

import "errors"

var (
	ErrParameterNotFound = errors.New("kurohelperdb: parameter not found")
	ErrNoRowsAffected    = errors.New("kurohelperdb: No rows affected")
	ErrUniqueViolation   = errors.New("kurohelperdb: 有重複主鍵的資料行存在")
	ErrInvalidRating     = errors.New("kurohelperdb: 評分需為 0 至 100，且最多一位小數")
)
