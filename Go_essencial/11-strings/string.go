package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	var resultado string
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Digite algo: ")
	resultado, _ = reader.ReadString('\n')

	palavras := strings.Split(resultado, " ")

	fmt.Println("Quantidade de palavras digitadas: ", len(palavras))
}
