package main

import "fmt"

func main() {
	var listaDeNomes []string

	listaDeNomes = append(listaDeNomes, "DESEC")
	listaDeNomes = append(listaDeNomes, "Golang")
	listaDeNomes = append(listaDeNomes, "Course")

	fmt.Println(listaDeNomes[3])
}
