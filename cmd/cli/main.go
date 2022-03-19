package main

import (
	"errors"
	"log"
	"os"

	"github.com/ald3v/celeritas"
	"github.com/fatih/color"
)

const version = "1.0.0"
var cel celeritas.Celeritas

func main(){
	var message string
	arg1, arg2, arg3, err := validateInput()
	if err != nil {
		exitGracefully(err)
	}

	setup(arg1,arg2)

	switch arg1 {
	case "help":
		showHelp()

	case "new":
		if arg2 == "" {
			exitGracefully(errors.New("new requires an applicaiton name"))
		}
		doNew(arg2)

	case "version":
		color.Yellow("Application version: " + version)
	case "migrate":
		if arg2 == "" {
			arg2 = "up"
		}
		err = doMigrate(arg2, arg3)
		if err != nil {
			exitGracefully(err)
		}
		message = "Migrations complete!"

	case "make":
		if arg2 == "" {
			exitGracefully(errors.New("make requires a subcommand: (migration|model|handler)"))
		}
		err = doMake(arg2, arg3)
		if err != nil {
			exitGracefully(err)
		}
	default:
		log.Println(arg2, arg3)
	}

	exitGracefully(nil, message)

}

func validateInput()(string,string,string, error){
	var arg1, arg2, arg3 string

	if len(os.Args) > 1 {
		arg1 = os.Args[1]

		if len(os.Args) >= 3 {
			arg2 = os.Args[2]
		}

		if len(os.Args) >= 4 {
			arg3 = os.Args[3]
		}
	}else {
		color.Red("Error: command required")
		showHelp()
		return "","","", errors.New("command required")
	}

	return arg1,arg2,arg3, nil
}

func exitGracefully(err error, msg ... string){
		message := ""
		if len(message) > 0 {
			message = msg[0]
		}
		if err != nil {
			color.Red("Error: %v\n", err)
		}

		if len(message) > 0 {
			color.Yellow(message)
		}else{
			color.Green(message)
		}

		os.Exit(1)
}