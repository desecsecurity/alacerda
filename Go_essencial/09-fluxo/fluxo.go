package main

import "fmt"

func main() {
	var listaDeNomes []string

	listaDeNomes = append(listaDeNomes, "DESEC")
	listaDeNomes = append(listaDeNomes, "Golang")
	listaDeNomes = append(listaDeNomes, "Course")

	for indice, nome := range listaDeNomes {
		if nome == "Course" {
			fmt.Println("Indice: ", indice)
		}
	}
}
