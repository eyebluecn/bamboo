package main

import (
	"github.com/eyebluecn/bamboo/code/core"
	"github.com/eyebluecn/bamboo/code/support"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	core.APPLICATION = &support.BambooApplication{}
	core.APPLICATION.Start()

}
