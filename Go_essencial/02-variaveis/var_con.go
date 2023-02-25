package main

import "fmt"

var minhaVariavel string

const minhaConstante string = "DESEC"

func main() {
	var minhaVariavel string
	const minhaConstante string = "DESEC"

	minhaVariavel = "meu codigo em golang"

	fmt.Println(minhaVariavel)
	fmt.Println(minhaConstante)

	imprime()
}
func imprime() {
	minhaVariavel = "variavel em imprime"

	fmt.Println(minhaVariavel)
	fmt.Println(minhaConstante)
}
