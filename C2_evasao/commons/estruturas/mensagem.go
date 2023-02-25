package estruturas

type Mensagem struct {
	AgentID       string
	AgentHostname string
	AgentCWD      string
	Comandos      []Commando
}
