package embedded

import "embed"

var (
	//go:embed static/*
	Static embed.FS
)
