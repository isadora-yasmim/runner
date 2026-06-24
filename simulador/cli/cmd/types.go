package cmd

// ResponseOutput é a estrutura padrão de resposta do assinador.jar,
// compartilhada pelos comandos sign, verify e stop.
type ResponseOutput struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
