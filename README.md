# Deployment-Tool
A [Go](https://golang.org/dl/) based deployment tool that allows the users to deploy the web application on the server using SSH information and pem file. 
This application is intended for non tecnhincal users they can just open the GUI and choose the server just deploy.
This application expect the pem files and `config.json` file on the root folder of the project. In the `config.json` file have all the username IP and pem file path of the server.

## Dependencies
This application uses the [Fyne](https://developer.fyne.io/started/) package for build the GUI , make sure all the requirement for fyne are installed based on your platform.

## Run and Build
```
go run main.go

go build
```

## Windows build from Linux
```
 GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc go build
```
## Usage
In the `main.go` file you can just add you shell script inside the ` Executeshell ` function. Your server `.pem` file is expected in the root of the project folder.
The file name must match the `myserver.pem` or you can modifiy the `.pem` file name in the code. after taking the build also the file should be present in the same folder of executables.
```
script := `
	cd /home/web/public_html/
	git pull origin master
	php artisan migrate
	ls
	cat index.html
	`
```
## License
GNU General Public License v3.0


