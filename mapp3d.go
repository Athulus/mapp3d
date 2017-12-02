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
	digest "github.com/xinsnake/go-http-digest-auth-client"
)

var client http.Client
var baseURL string
var uname string
var id string
var key string
var printer = "um3"
var transport digest.DigestTransport

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
			Name: "lights",
			Flags: []cli.Flag{
				cli.UintFlag{Name: "hue", Value: 0},
				cli.UintFlag{Name: "saturation", Value: 0},
				cli.UintFlag{Name: "brightness, value", Value: 100},
			},
			Usage:  "mess with the printers lights",
			Action: lights,
		},
		{
			Name: "print",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "model, m", Usage: "the path to the 3d model you want to print (.stl)"},
				cli.StringFlag{Name: "slicer", Usage: "the path to the slicing configuration file you want to use"},
			},
			Usage:  "input a 3d model, and it will be sliced and the gcode sent to your printer",
			Action: print,
		},
		{
			Name:   "slice",
			Usage:  "configure the slicing properties for your prints",
			Action: slice,
		},
		{
			Name:   "init",
			Usage:  "register the app with the printer. get an access token",
			Action: startup,
		},
	}

	app.Run(os.Args)
}

//checks for config file and setup
func init() {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}
	settings, err := ioutil.ReadFile(usr.HomeDir + "/.mapp3d")
	if err != nil {
		fmt.Println(err)
		fmt.Println("we need to create a configuaration file bforeto continue")
		fmt.Println("the file will be saved to '~/.mapp3d'")
		fmt.Println("creating config file ...")
		settings = makeConfigFile(usr.HomeDir + "/.mapp3d")
	}
	var cfg config
	err = json.Unmarshal(settings, &cfg)
	baseURL = cfg.BaseURL
	uname = cfg.User
	key = cfg.Key
	id = cfg.ID
	transport = digest.NewTransport(key, id)

}

func startup(c *cli.Context) {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}
	settings, err := ioutil.ReadFile(usr.HomeDir + "/.mapp3d")
	if err != nil {
		fmt.Println(err)
		fmt.Println("creating config file ...")
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
	fmt.Println("enter the baseurl for the printer API")
	baseURL, err = bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		panic(err)
	}
	baseURL = strings.TrimSpace(baseURL)

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
	type lightsRequest struct {
		Hue        uint `json:"hue"`
		Saturation uint `json:"saturation"`
		Brightness uint `json:"brightness"`
	}
	body, err := json.Marshal(lightsRequest{c.Uint("hue"), c.Uint("saturation"), c.Uint("brightness")})
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("PUT", baseURL+"printer/led", bytes.NewBuffer(body))
	res, err := transport.RoundTrip(req)
	resBody, err := ioutil.ReadAll(res.Body)
	for k, v := range res.Header {
		fmt.Println(k, v)
	}
	fmt.Println(res.Status)
	fmt.Println(string(resBody))
}

func status(c *cli.Context) {
	request, err := http.NewRequest("GET", baseURL+"auth/check/"+id, nil)
	if err != nil {
		fmt.Println(err)
	}
	response, err := client.Do(request)
	body, err := ioutil.ReadAll(response.Body)
	fmt.Println("the printer status is", string(body))
}

func slice(c *cli.Context) {}
func print(c *cli.Context) {}
