package main

import (
	"github.com/henryhsue/oura-app/oura"
	"github.com/henryhsue/oura-app/sheets"
)

func main() {
	oura.Run()
	sheets.WriteToSheet()
}
