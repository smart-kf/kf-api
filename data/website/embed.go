package website

import "embed"

//go:embed views
var FS embed.FS

//go:embed static/*
var StaticFS embed.FS
