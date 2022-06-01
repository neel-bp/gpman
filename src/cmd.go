package src

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"syscall"

	"github.com/atotto/clipboard"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/term"
)

const HELP string = `
gpman is a simple cli password manager

Usage:

	gpman <command> [OPTIONS] [ARGUMENTS]

The commands are:

	store       store credentials (username/email, password) against a given site/service
	get         retrieve stored credentials
	list        list all stored credentials
	delete      delete credentials stored against given site/service
	help        print help message or command specific help for given command
`

const STOREHELP string = `
Usage:	gpman store [OPTIONS] [ARGUMENTS...]

store credentials (username/email, password) against a given site/service

Options:

-w [BOOL]	starts store wizard which asks for site, username, password 
		and passphrase (passphrase is used to encrypt given values)

-gp [INT]	generates a password of given length and stores it as 
		password against given site/service, when given password prompt is skipped in wizard
		if wizard flag (-w) is not present password argument is no longer needed 
		generated password is considered instead.

-sc [BOOL]	by default generate password flag (-gp) includes special characters 
		to opt out of them special character flag is used (-sc)
`

const GETHELP string = `
Usage: gpman get [OPTIONS] [SITE]:REQUIRED

retrieve stored credentials

Options:

-clip [BOOL]	by default get command prints password to stdout but it is possible to
		write password directly to clipboard using -clip flag requires xclip
		on linux
`

const LISTHELP string = `
Usage: gpman list [OPTIONS]

list all stored credentials
NOTE: gpman assumes you are saving all credentials with same passphrase but does
not restricts you from using multiple passphrase however if you choose to decrypt
password or username while using list command the command will fail because it only
prompts you for one passphrase and tries decrypting all credentials using that passphrase

Options:

-p [BOOL]	list command by default only print name of site and does not decrypt username and password
		-p flag decrypts passwords of all sites and prints to stdout along with everything else
	
-u [BOOL]	list command by default only print name of site and does not decrypt username and password
		-u flag decrypts usernames of all sites and prints to stdout along with everything else
`

const DELHELP string = `
Usage: gpman delete [SITE]:REQUIRED

delete credentials stored against given site/service
`

var SaveCmd = flag.NewFlagSet("store", flag.ExitOnError)
var GetCmd = flag.NewFlagSet("get", flag.ExitOnError)
var ListCmd = flag.NewFlagSet("list", flag.ExitOnError)
var DelCmd = flag.NewFlagSet("delete", flag.ExitOnError)
var HelpCmd = flag.NewFlagSet("help", flag.ExitOnError)

var HELPMAP = map[string]string{
	SaveCmd.Name(): STOREHELP,
	GetCmd.Name():  GETHELP,
	ListCmd.Name(): LISTHELP,
	DelCmd.Name():  DELHELP,
}

// Errors =====
var ErrInsuficientArgs = errors.New("wrong Number of arguments\nRun 'gpman help' for Usage")

func HandleSaveCommand(args []string) error {

	wizard := SaveCmd.Bool("w", false, "")
	passlen := SaveCmd.Int("gp", 0, "")
	specialChars := SaveCmd.Bool("sc", false, "")
	SaveCmd.Usage = func() {
		fmt.Fprintln(SaveCmd.Output(), STOREHELP)
	}
	SaveCmd.Parse(args)
	if *wizard {

		var username, site string
		var password, passphrase []byte

		fmt.Print("enter site/service: ")
		_, err := fmt.Scan(&site)
		if err != nil {
			return err
		}

		fmt.Print("enter username: ")
		_, err = fmt.Scan(&username)
		if err != nil {
			return err
		}

		if *passlen == 0 {
			fmt.Print("enter password: ")
			password, err = term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return err
			}
		} else {
			passwordb, err := GenerateRandomPswd(*passlen, !*specialChars)
			if err != nil {
				return err
			}
			password = []byte(passwordb)
		}

		fmt.Print("\nenter passphrase: ")
		passphrase, err = term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}

		return JsonWriter(string(passphrase), site, username, string(password))
	}

	if len(SaveCmd.Args()) < 3 && *passlen == 0 {
		return ErrInsuficientArgs
	}
	if *passlen > 0 && len(SaveCmd.Args()) < 2 {
		return ErrInsuficientArgs
	}
	fmt.Print("enter passphrase: ")
	passphrase, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}
	var passwd string
	if *passlen > 0 {
		passwd, err = GenerateRandomPswd(*passlen, !*specialChars)
		if err != nil {
			return err
		}
	} else {
		passwd = SaveCmd.Args()[2]
	}
	return JsonWriter(string(passphrase), SaveCmd.Args()[0], SaveCmd.Args()[1], passwd)
}

func HandleGetCommand(args []string) error {
	if len(args) == 0 {
		return ErrInsuficientArgs
	}
	clip := GetCmd.Bool("clip", false, "")
	GetCmd.Usage = func() {
		fmt.Fprintln(GetCmd.Output(), GETHELP)
	}
	GetCmd.Parse(args)
	site := GetCmd.Args()[0]
	fmt.Print("enter passphrase: ")
	passphrase, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}
	username, passwd, err := JsonReader(string(passphrase), site)
	if err != nil {
		return err
	}
	if *clip {
		fmt.Printf("\nyour credentials for %s:\nusername: %s\nyour password has been copied to clipboard", site, username)
		err := clipboard.WriteAll(passwd)
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("\nyour credentials for %s:\nusername: %s\npassword: %s", site, username, passwd)
	}
	return nil

}

func HandleDeleteCommand(args []string) error {
	if len(args) == 0 {
		return ErrInsuficientArgs
	}
	DelCmd.Usage = func() {
		fmt.Fprintln(DelCmd.Output(), DELHELP)
	}
	DelCmd.Parse(args)
	site := DelCmd.Args()[0]
	err := JsonDelete(site)
	if err != nil {
		return err
	}
	return nil
}

func HandleListCommand(args []string) error {
	user_flag := ListCmd.Bool("u", false, "")
	pass_flag := ListCmd.Bool("p", false, "")
	HEADERS := []string{"site/service", "username", "password"}
	var passphrase []byte
	var err error
	ListCmd.Usage = func() {
		fmt.Fprintln(ListCmd.Output(), LISTHELP)
	}
	ListCmd.Parse(args)
	if *user_flag || *pass_flag {
		fmt.Print("enter passphrase: ")
		passphrase, err = term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
	}
	data, err := ListPasses(string(passphrase), *user_flag, *pass_flag)
	if err != nil {
		return err
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(HEADERS)
	for _, v := range data {
		table.Append(v)
	}
	fmt.Println()
	table.Render()
	return nil
}

func HandleHelpCommand(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(os.Stdout, HELP)
		return nil
	}
	h, ok := HELPMAP[args[0]]
	if !ok {
		return fmt.Errorf("gpman help %s: unknown help topic. Run 'gpman help'", args[0])
	}
	fmt.Fprintln(os.Stdout, h)
	return nil

}

func CommandHandler(args []string) error {
	if len(args) < 2 {
		HandleHelpCommand(nil)
		return nil
	}
	switch args[1] {
	case SaveCmd.Name():
		return HandleSaveCommand(args[2:]) // os.Args[0] is program name os.Args[1] is subcommand name rest of the commands are passed to appropriate handler
	case GetCmd.Name():
		return HandleGetCommand(args[2:])
	case DelCmd.Name():
		return HandleDeleteCommand(args[2:])
	case ListCmd.Name():
		return HandleListCommand(args[2:])
	case HelpCmd.Name():
		return HandleHelpCommand(args[2:])
	default:
		return fmt.Errorf("gpman %s: unknown command\nRun 'gpman help' for usage", args[1])
	}
}
