package main

import (
	"desec_estruturas/modelos"
	"fmt"
)

var pessoa = modelos.Pessoa{}

func main() {
	pessoa.Nome = "Alan Lacerda"
	pessoa.Endereco = "Canada"
	pessoa.Idade = 100
	pessoa.Fuma = false

	imprimeDetalhePessoa()
}

func imprimeDetalhePessoa() {
	fmt.Println("Nome: ", pessoa.Nome)
	fmt.Println("Endereco: ", pessoa.Endereco)
	fmt.Println("Idade: ", pessoa.Idade)
	fmt.Println("Fumante: ", pessoa.Fuma)
}
