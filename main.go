package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type Config struct {
	Config []Apps `json:"apps"`
}
type Apps struct {
	Name      string `json:"name"`
	Appcode   string `json:"appcode"`
	Username  string `json:"username"`
	Ipaddress string `json:"ipaddress"`
	Pemfile   string `json:"pemfile"`
	Passcode  string `json:"passcode"`
}

var app_selected string
var appcreds = make(map[string]string)

func main() {

	os.Setenv("FYNE_THEME", "light")

	a := app.New()
	win := a.NewWindow("TiaTech Application Deployment Tool [Version 1.2.0]")
	win.Resize(fyne.NewSize(800, 400))
	win.CenterOnScreen()
	apps := parseConfigFile()
	applications := []string{}

	for _, items := range apps.Config {

		applications = append(applications, items.Name)
		appcreds[items.Name] = items.Appcode
		content, _ := json.Marshal(items)
		appcreds[items.Appcode] = string(content)
	}

	var split_index int
	if len(applications) > 6 {
		split_index = 5
	} else if len(applications) > 3 && len(applications) <= 6 {
		split_index = 3
	} else {
		split_index = 2
	}

	firstset := append(applications[:split_index])
	secondset := append(applications[split_index:])

	var secondradios *widget.RadioGroup
	var firstradios *widget.RadioGroup
	firstradios = widget.NewRadioGroup(firstset, func(vals string) {
		secondradios.SetSelected("")
		firstradios.SetSelected(vals)
		app_selected = appcreds[vals]

	})

	secondradios = widget.NewRadioGroup(secondset, func(vals string) {
		firstradios.SetSelected("")
		secondradios.SetSelected(vals)
		app_selected = appcreds[vals]
	})

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Security Code")
	textArea := widget.NewMultiLineEntry()
	textArea.Resize(fyne.NewSize(100, 200))
	cancelbutton := widget.NewButton("CLOSE", func() {

		os.Exit(0)
	})

	deploybutton := widget.NewButton("DEPLOY", func() {
		passcode := password.Text
		if app_selected != "" {
			textArea.SetText("The Deployment is in Progress ....")
			if password.Text != "" {
				result := Executeshell(app_selected, passcode)
				textArea.SetText(result)
			} else {
				textArea.SetText("Please Provide your Security Code")
			}
		} else {
			textArea.SetText("Please choose an option!")
		}
	})

	vlayout := container.New(layout.NewVBoxLayout(),
		container.New(layout.NewGridLayout(2), firstradios, secondradios),
		container.New(layout.NewGridLayout(1), password),
		container.New(layout.NewGridLayout(2), deploybutton, cancelbutton),
	)
	win.SetContent(container.New(layout.NewGridLayout(2), vlayout, textArea))

	win.ShowAndRun()
}

func parseConfigFile() Config {

	jsonFile, err := os.Open("config.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var config Config

	json.Unmarshal(byteValue, &config)

	return config

}

func getPHPWebShellScript() string {

	return `
	cd /home/admin/web/public_html
	ls
	`
}

func getDotNetShellScript() string {

	return `
	cd /home/admin/web/public_html
	ls
	`
}

func getTiaImageShellScript() string {

	return `
	cd /home/admin/web/public_html
	ls
	`
}

func getAssetShellScript() string {

	return `
	cd /home/admin/web/public_html
	ls
	`
}

func getTeleShellScript() string {

	return `
	cd /home/admin/web/public_html
	ls
	`
}

func getFinAppShellScript() string {

	return `
	cd /home/admin/web/public_html
	ls
	`
}

func Executeshell(app_selected string, passcode string) string {

	var script string

	if app_selected == "numrweb" {
		script = getPHPWebShellScript()
	}

	if app_selected == "numrdotnet" {
		script = getDotNetShellScript()
	}

	if app_selected == "numrtiaimage" {
		script = getTiaImageShellScript()
	}

	if app_selected == "aaset" {
		script = getAssetShellScript()
	}

	if app_selected == "tele" {
		script = getTeleShellScript()
	}
	if app_selected == "finapp" {
		script = getFinAppShellScript()
	}
	var config Apps
	json.Unmarshal([]byte(appcreds[app_selected]), &config)

	usrname := config.Username
	ipaddr := config.Ipaddress
	pemfile := config.Pemfile
	pass := config.Passcode
	plain_text := decyptMyPassCode(pass)

	if plain_text != passcode {
		return "Invalid Security Code!"
	}

	if usrname != "" && ipaddr != "" && pemfile != "" {

		c := exec.Command("ssh", usrname+"@"+ipaddr, "-i", "myserver.pem")
		var buf = new(bytes.Buffer)
		buf.WriteString(script)
		c.Stdin = buf

		b, e := c.Output()
		if e != nil {
			fmt.Println(e)
		}
		return (string(b))
	} else {
		return "Config file is not valid!"
	}
}

func encryptMyPassCode(passcode string) string {

	password := []byte(passcode)
	key := []byte("adbb347fd8f1260b7796fcc17bda48df")

	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
	}

	gcm, err := cipher.NewGCM(c)

	if err != nil {
		fmt.Println(err)
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}

	encrypted_text := gcm.Seal(nonce, nonce, password, nil)

	return hex.EncodeToString(encrypted_text)
}

func decyptMyPassCode(hexastring string) string {

	key := []byte("adbb347fd8f1260b7796fcc17bda48df")

	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
	}

	cipher, _ := hex.DecodeString(hexastring)

	nonceSize := gcm.NonceSize()
	if len(cipher) < nonceSize {
		fmt.Println(err)
	}

	nonce, cipher := cipher[:nonceSize], cipher[nonceSize:]
	plaintext, _ := gcm.Open(nil, nonce, cipher, nil)

	return string(plaintext)

}
