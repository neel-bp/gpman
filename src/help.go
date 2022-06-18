package src

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
	gitauth	    connect remote git repository for syncing credentials
	push        push all local credentials to remote git repository
`

const STOREHELP string = `
Usage:	gpman store [OPTIONS] [ARGUMENTS...]

store credentials (username/email, password) against a given site/service
Note: when not using the wizard the order of argument is:
[SITE/SERVICE] [USERNAME] [PASSWORD]

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
