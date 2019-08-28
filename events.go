package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"

	"1c-connect-events/amocrm"

	"github.com/ros-tel/1c-connect-pipe"
	"gopkg.in/yaml.v2"
)

type (
	TConfig struct {
		AgentLogin string `yaml:"agent_login"`
		CrmType    string `yaml:"crm_type"`

		CRMs struct {
			AmoCRM amocrm.AmoCRM `yaml:"amocrm"`
		} `yaml:"crms"`
	}
)

var (
	config *TConfig

	config_file = flag.String("config", "", "Usage: -config=<config_file>")
	debug       = flag.Bool("debug", false, "Print debug information on stderr")
)

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	flag.Parse()

	// Load the configuration file
	if *config_file == "" {
		*config_file = "config" + string(os.PathSeparator) + "config.yml"
	}

	getConfig(*config_file)

	switch config.CrmType {
	case "amocrm":
		amocrm.Start(config.CRMs.AmoCRM, debug)
	default:
		log.Fatal("Unsupported CRM: " + config.CrmType)
	}

	// Инициируем клиент
	pipe_client := pipe.InitClient(config.AgentLogin, *debug)

	/*
	 *  \/ Подписываемся на события софтфона: \/
	 */
	// ... входящие вызовы
	pipe_client.SendCommand(pipe.Command{
		Action:    "EventSubscribe",
		Mode:      "Softphone",
		Object:    "Call",
		Initiator: "Incoming",
	})
	// ... исходящие вызовы
	pipe_client.SendCommand(pipe.Command{
		Action:    "EventSubscribe",
		Mode:      "Softphone",
		Object:    "Call",
		Initiator: "Self",
	})
	/*
	 *  /\ Подписались на события софтфона /\
	 */

	// В отдельной рутине читаем события из каналов
	go func() {
		for {
			select {
			case e := <-pipe_client.Event:
				if *debug {
					log.Printf("Event received %+v", e)
				}
				if e.Mode == "Softphone" && e.Object == "Call" {
					switch config.CrmType {
					case "amocrm":
						amocrm.SendEvent(e)
					}
				}
			case r := <-pipe_client.Result:
				if *debug {
					log.Printf("CommandResult received %+v", r)
				}
			}
		}
	}()

	// Запускаем клиент
	pipe_client.Start()
}

// Load the YAML config file
func getConfig(configFile string) {
	var err error
	var input = io.ReadCloser(os.Stdin)
	if input, err = os.Open(configFile); err != nil {
		log.Fatalln(err)
	}
	defer input.Close()

	// Read the config file
	yamlBytes, err := ioutil.ReadAll(input)

	if err != nil {
		log.Fatalln(err)
	}

	// Parse the config
	if err := yaml.Unmarshal(yamlBytes, &config); err != nil {
		//log.Fatalf("Content: %v", yamlBytes)
		log.Fatalf("Could not parse %q: %v", configFile, err)
	}
}
