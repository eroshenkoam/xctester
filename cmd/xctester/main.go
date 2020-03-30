package main

import (
	"github.com/eroshenkoam/xctester/commands"
	"log"
)

func main() {
	log.SetPrefix("xctester: ")
	commands.Execute()
}

