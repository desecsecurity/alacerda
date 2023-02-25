package main

import (
	"github.com/desec/estrutura"
	"github.com/desec/helpers"
	"github.com/sirupsen/logrus"
)

func main() {
	helpers.Urgente()
	println(estrutura.Teste)

	logrus.SetLevel(logrus.DebugLevel)

	logrus.Debug("Mensagem de Debug do Logrus")
	logrus.Warnln("Mensagem de Warn do Logrus")
}
