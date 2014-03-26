package commands

import (
	"cf/app"
	"cf/commands"
	"flag"
	"fmt"
	"github.com/codegangsta/cli"
	"strings"
	testreq "testhelpers/requirements"
)

func NewContext(cmdName string, args []string) *cli.Context {
	targetCommand := findCommand(cmdName)

	flagSet := new(flag.FlagSet)
	for i, _ := range targetCommand.Flags {
		targetCommand.Flags[i].Apply(flagSet)
	}

	// move all flag args to the beginning of the list, go requires them all upfront
	firstFlagIndex := -1
	for index, arg := range args {
		if strings.HasPrefix(arg, "-") {
			firstFlagIndex = index
			break
		}
	}
	if firstFlagIndex > 0 {
		args := args[0:firstFlagIndex]
		flags := args[firstFlagIndex:]
		flagSet.Parse(append(flags, args...))
	} else {
		flagSet.Parse(args[0:])
	}

	globalSet := new(flag.FlagSet)

	return cli.NewContext(cli.NewApp(), flagSet, globalSet)
}

func findCommand(cmdName string) (cmd cli.Command) {
	cmdFactory := commands.ConcreteFactory{}
	reqFactory := &testreq.FakeReqFactory{}
	cmdRunner := commands.NewRunner(cmdFactory, reqFactory)
	myApp, _ := app.NewApp(cmdRunner)

	for _, cmd := range myApp.Commands {
		if cmd.Name == cmdName {
			return cmd
		}
	}
	panic(fmt.Sprintf("command %s does not exist", cmdName))
	return
}

//func NewContext(cmdName string, args []string) *cli.Context {
//	commandFactory := commands.ConcreteFactory{}
//	cmdRunner := commands.NewRunner(commandFactory, &testreq.FakeReqFactory{})
//	testApp, _ := app.NewApp(cmdRunner)
//
//
//	args = append([]string{cmdName}, args...)
//
//	set := flagSet(os.Args[0], testApp.Flags)
//	set.SetOutput(ioutil.Discard)
//	err := set.Parse(args)
//	if err != nil {
//		panic(err)
//	}
//
//	nerr := normalizeFlags(testApp.Flags, set)
//	if nerr != nil {
//		panic(nerr)
//	}
//
//	return cli.NewContext(testApp, set, set)
//}
//
//// Copy-pasta private functions from codegangsta
//func flagSet(name string, flags []cli.Flag) *flag.FlagSet {
//	set := flag.NewFlagSet(name, flag.ContinueOnError)
//
//	for _, f := range flags {
//		f.Apply(set)
//	}
//	return set
//}
//
//func normalizeFlags(flags []cli.Flag, set *flag.FlagSet) error {
//	visited := make(map[string]bool)
//	set.Visit(func(f *flag.Flag) {
//		visited[f.Name] = true
//	})
//
//	for _, f := range flags {
//		var name string
//		switch typedF := f.(type) {
//		case cli.StringFlag:
//			name = typedF.Name
//		case cli.IntFlag:
//			name = typedF.Name
//		case cli.BoolFlag:
//			name = typedF.Name
//		case cli.StringSliceFlag:
//			name = typedF.Name
//		case cli.Float64Flag:
//			name = typedF.Name
//		}
//
//		parts := strings.Split(name, ",")
//		if len(parts) == 1 {
//			continue
//		}
//		var ff *flag.Flag
//		for _, name := range parts {
//			name = strings.Trim(name, " ")
//			if visited[name] {
//				if ff != nil {
//					return errors.New("Cannot use two forms of the same flag: " + name + " " + ff.Name)
//				}
//				ff = set.Lookup(name)
//			}
//		}
//		if ff == nil {
//			continue
//		}
//		for _, name := range parts {
//			name = strings.Trim(name, " ")
//			set.Set(name, ff.Value.String())
//		}
//	}
//	return nil
//}
