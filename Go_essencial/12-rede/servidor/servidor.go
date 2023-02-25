package main

import (
	"bufio"
	"log"
	"net"
)

func main() {
	listener, erro := net.Listen("tcp", "0.0.0.0:9090")
	if erro != nil {
		log.Fatal("Erro ao configurar o Listener.", erro.Error())
	}

	conexao, erro := listener.Accept()
	if erro != nil {
		log.Fatal("Erro ao iniciar uma nova conexao.", erro.Error())
	}

	for {
		mensagem, _ := bufio.NewReader(conexao).ReadString('\n')
		conexao.Write([]byte(mensagem))
	}

}
