package main

import (
	"net/http"
	"fmt"
)

func main()  {
	
}

func run() {
	fmt.Println("Старт сервиса")
	err := http.ListenAndServe()
}