package data

import (
	"embed"
	_ "embed"
)

//go:embed apartment.json
var ApartmentStr string

//go:embed *
var ConfigJsonsFile embed.FS
