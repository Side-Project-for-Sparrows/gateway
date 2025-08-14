package main

import (
	"github.com/Side-Project-for-Sparrows/gateway/config"
	"github.com/Side-Project-for-Sparrows/gateway/internal/router"
)

func main() {
	initialize()
	route()
}

func initialize() {
	config.InitAll()
	config.ConstructAll() // 병관아 또 까먹고 지우면 안돼. jwt util 초기화 여기서한다?
}

func route() {
	router.InitRoute()
}
