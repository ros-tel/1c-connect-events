package asteriskami

import (
	"log"

	"github.com/heltonmarx/goami/ami"
	"github.com/ros-tel/1c-connect-pipe"
	"github.com/satori/go.uuid"
)

var (
	config Ami
	debug  *bool

	commands = make(chan string, 100)
)

type (
	Ami struct {
		Username string            `yaml:"user"`
		Password string            `yaml:"password"`
		Addr     string            `yaml:"addr"`
		Events   map[string]string `yaml:"events"`
	}
)

func Start(c Ami, d *bool) {
	config = c
	debug = d

	socket, err := ami.NewSocket(c.Addr)
	if err != nil {
		log.Fatalf("socket error: %v\n", err)
	}
	if _, err := ami.Connect(socket); err != nil {
		log.Fatalf("connect error: %v\n", err)
	}

	//Login
	uuid := uuid.NewV4() // ami.GetUUID() в версии v1.0.1-0.20190919213858-fa0badef58ed не может под Windows работать
	if err := ami.Login(socket, c.Username, c.Password, "Off", uuid.String()); err != nil {
		log.Fatalf("login error: %v\n", err)
	}
	log.Printf("login ok!\n")

	go processCalls(socket)
}

// Получение события и постановка в очереди
func SendEvent(e *pipe.Event) {
	if e.Mode == "ServicesClients" && e.Object == "AgentOnlineStatus" && e.Initiator == "Self" {
		if *debug {
			log.Println("Key:", e.Mode+"."+e.Object+"."+e.Initiator+"."+e.Status)
		}
		if command, ok := config.Events[e.Mode+"."+e.Object+"."+e.Initiator+"."+e.Status]; ok {
			commands <- command
		}
	}
}

// Собираем логи о звонках
func processCalls(api *ami.Socket) {
	var err error
	for {
		select {
		case command := <-commands:
			if *debug {
				log.Println("Send Command:", command)
			}
			err = api.Send(command)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
