package main

import (
	"./routers"
)

func main() {
	r := routers.Init()
	r.Run(":31204")
}
