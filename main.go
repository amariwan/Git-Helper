package main

import (
	"flag"
	"log"

	"github.com/artyom/autoflags"

	"aland-mariwan.de/githelper/directoryTreeCommands"
	"aland-mariwan.de/githelper/githelperJsonCommands"
)

// Importing packages

// Main function
func main() {
	var application = Application{}
	application.Init()
	application.Run()
}

type commandlineArguments struct {
	Command  string
	LogLevel int `flag:"l,log level."`
}

type Application struct {
	args         commandlineArguments
	cloneAndPull githelperJsonCommands.CloneAndPull
	status       githelperJsonCommands.Status
	dPull        directoryTreeCommands.Pull
	dStatus      directoryTreeCommands.Status
}

func (application *Application) Init() {
	application.parseArgs()
	if flag.NArg() > 0 {
		application.args.Command = flag.Args()[0]
	}

	if len(application.args.Command) <= 0 {
		log.Fatal("You need to provide an command. Exiting...")
	}

	application.cloneAndPull = githelperJsonCommands.CloneAndPull{}
	application.status = githelperJsonCommands.Status{}
	application.dPull = directoryTreeCommands.Pull{}
	application.dStatus = directoryTreeCommands.Status{}
}

func (application *Application) Run() {
	switch application.args.Command {
	case "clone":
		application.cloneAndPull.DoClone = true
		application.cloneAndPull.DoPull = false
		application.cloneAndPull.CloneAndPullCmd()
	case "pull":
		application.cloneAndPull.DoClone = false
		application.cloneAndPull.DoPull = true
		application.cloneAndPull.CloneAndPullCmd()
	case "sync":
		application.cloneAndPull.DoClone = true
		application.cloneAndPull.DoPull = true
		application.cloneAndPull.CloneAndPullCmd()
	case "status":
		application.status.StatusCmd()
	case "run":
		if flag.NArg() > 1 {
			var args []string
			for k, v := range flag.Args() {
				if k > 0 {
					if k == 1 && v == "--" {
						continue
					}
					args = append(args, v)

				}
			}
			log.Print(args)
		}
	case "dpull":
		application.dPull.PullCmd()
	case "dstatus":
		application.dStatus.StatusCmd()
	}
}

func (application *Application) setCommandlineDefaultValues() {
	application.args = commandlineArguments{
		Command:  "",
		LogLevel: 1,
	}
}

func (application *Application) parseArgs() {
	application.setCommandlineDefaultValues()
	autoflags.Define(&application.args)
	flag.Parse()
}
