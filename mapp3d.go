package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"strings"

	"github.com/urfave/cli"
)

var client http.Client
var baseURL = "http://ultimaker.hackrva.org/api/api/v1/"

type requestBody struct {
	Application string `json:"application"`
	User        string `json:"user"`
}

type config struct {
	User    string
	ID      string `json:"id"`
	Key     string `json:"key"`
	BaseURL string
	Printer string
}

func main() {
	app := cli.NewApp()

	app.Name = "mapp3d"
	app.Usage = "CLI tool to talk to the ultimaker 3 API"

	app.Commands = []cli.Command{
		{
			Name:   "status",
			Usage:  "prints out the printers status",
			Action: status,
		},
		{
			Name:   "lights",
			Usage:  "mess with the printers lights",
			Action: lights,
		},
		{
			Name:   "init",
			Usage:  "register the app with the printer. get an access token",
			Action: startup,
		},
	}

	app.Run(os.Args)
}

func startup(c *cli.Context) {
	fmt.Println("init test")
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}
	settings, err := ioutil.ReadFile(usr.HomeDir + "/.mapp3d")
	if err != nil {
		fmt.Println(err)
		fmt.Println("creating settings file ...")
		settings = makeConfigFile(usr.HomeDir + "/.mapp3d")
	}
	fmt.Println(string(settings))
}

// MakeConfigFile will generate the apps configuartion and save it to a file
func makeConfigFile(filePath string) []byte {

	fmt.Println("enter a username")
	uname, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		panic(err)
	}
	uname = strings.TrimSpace(uname)
	reqBody, err := json.Marshal(requestBody{"max's application", uname})
	fmt.Println(string(reqBody))

	res, err := http.Post(baseURL+"auth/request", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

	var cfg config
	json.Unmarshal(body, &cfg)
	cfg.BaseURL = baseURL
	cfg.Printer = "um3"
	cfg.User = uname

	fmt.Println(cfg)

	cfgJSON, err := json.Marshal(cfg)
	ioutil.WriteFile(filePath, cfgJSON, 0666)
	return cfgJSON
}

func lights(c *cli.Context) {
	fmt.Println("lights test")
}

func status(c *cli.Context) {
	fmt.Println("status test")
	request, err := http.NewRequest("GET", baseURL+"printer/status", nil)
	if err != nil {
		fmt.Println(err)
	}
	response, err := client.Do(request)
	body, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
}
