package main

import (
	fte "../lib"
	cfg "../lib/common"
	ctr "../lib/controller"
)

func init() {
	cfg.InitConfig()
	fte.InitMail()
	fte.InitView()
	ctr.InitSession()
}

func main() {
	fte.InitRouter()
}
