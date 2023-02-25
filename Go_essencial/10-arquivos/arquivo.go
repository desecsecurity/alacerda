package main

import (
	"log"
	"os"
	"path/filepath"
)

func main() {
	// Criando um arquivo
	//os.Create("meu_arquivo.txt")

	filepath.FromSlash(filepath.Join(os.TempDir(), "mnha_Pasta", "arquivo.txt"))

	arquivo, erro := os.OpenFile("arquivo_teste.txt", os.O_APPEND|os.O_CREATE, 660)

	if erro != nil {
		log.Panic(erro.Error())
	}
	defer arquivo.Close()

	arquivo.WriteString("DESEC Security\n")

}
