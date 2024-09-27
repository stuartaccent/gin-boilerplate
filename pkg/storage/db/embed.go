package db

import "embed"

var (
	//go:embed migrations/*.sql
	Migrations embed.FS
)
