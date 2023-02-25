package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conexao, erro := net.Dial("tcp", "10.0.0.211:9090")
	if erro != nil {
		log.Fatal("Erro ao conectar com o servidor: ", erro.Error())
	}

	for {
		mensagem, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		conexao.Write([]byte(mensagem))

		resposta, _ := bufio.NewReader(conexao).ReadString('\n')
		fmt.Println("SERVIDOR: ", resposta)
	}

}
