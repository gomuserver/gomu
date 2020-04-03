package main

import (
	"os"
	"strings"
	"testing"

	flag "github.com/hatchify/parg"
	"github.com/hatchify/simply"
)

func TestParse_Empty(context *testing.T) {
	input := "gomu"
	os.Args = strings.Split(input, " ")

	command, err := getCommand()

	errTest := simply.Target(err, context, "Error should be nil")
	errTest.Assert().Equals(nil)
	errTest.Validate(errTest)

	cmdTest := simply.Target(command, context, "Cmd should not be nil")
	cmdTest.Assert().DoesNotEqual(nil)
	cmdTest.Validate(cmdTest)

	actTest := simply.Target(command.Action, context, "Action should be empty")
	actResult := actTest.Equals("")
	actTest.Validate(actResult)

	argTest := simply.Target(command.Arguments, context, "Args should be empty")
	argResult := argTest.Equals([]*flag.Argument{})
	argTest.Validate(argResult)

	flagTest := simply.Target(command.Flags, context, "Flags should be empty")
	flagResult := flagTest.Equals(map[string]*flag.Flag{})
	flagTest.Validate(flagResult)
}

// This test cannot pass with default parse rules
//   1) bool flags immediately preceding command names
//   2_ array flags before command is set
// Both justify custom config for flag parsing
func TestConfig_1Flag_1FlagMatch_1BoolFlag_Cmd_1Flag_2Arg_2FlagArrayMatch(context *testing.T) {
	input := "gomu -include test1 -include test2 -name-only sync -b JIRA-Ticket mod-common simply -i hatchify vroomy"
	os.Args = strings.Split(input, " ")

	command, err := getCommand()

	test := simply.Target(err, context, "Error should not exist")
	result := test.Assert().Equals(nil)
	test.Validate(result)

	test = simply.Target(command, context, "Command should exist")
	result = test.Assert().DoesNotEqual(nil)
	test.Validate(result)

	test = simply.Target(command.Action, context, "Action should be <sync>")
	result = test.Equals("sync")
	test.Validate(result)

	test = simply.Target(len(command.Arguments), context, "Arguments should have 2 elements")
	result = test.Equals(2)
	test.Validate(result)

	test = simply.Target(command.Arguments[0].Name, context, "Argument[0] should be mod-common")
	result = test.Equals("mod-common")
	test.Validate(result)

	test = simply.Target(command.Arguments[1].Name, context, "Argument[1] should be simply")
	result = test.Equals("simply")
	test.Validate(result)

	test = simply.Target(len(command.Flags), context, "Flags should have 3 elements")
	result = test.Equals(3)
	test.Validate(result)

	includes, _ := command.Flags["-include"].Value.([]string)
	test = simply.Target(includes[0], context, "-i[0] should be test1")
	result = test.Equals("test1")
	test.Validate(result)
	test = simply.Target(includes[1], context, "-i[1] should be test2")
	result = test.Equals("test2")
	test.Validate(result)
	test = simply.Target(includes[2], context, "-i[2] should be hatchify")
	result = test.Equals("hatchify")
	test.Validate(result)
	test = simply.Target(includes[3], context, "-i[3] should be vroomy")
	result = test.Equals("vroomy")
	test.Validate(result)

	branch, _ := command.Flags["-branch"].Value.(string)
	test = simply.Target(branch, context, "-b should be JIRA-Ticket")
	result = test.Equals("JIRA-Ticket")
	test.Validate(result)

	nameOnly, _ := command.Flags["-name-only"].Value.(bool)
	test = simply.Target(nameOnly, context, "-name-only should be true")
	result = test.Equals(true)
	test.Validate(result)
}

func TestConfig_1BoolFlag_2FlagArray_Cmd_1Flag_1Arg_1BoolFlag_1Arg_2FlagArrayMatch_1BoolFlag(context *testing.T) {
	input := "gomu -name -include test1 -include test2 sync -b JIRA-Ticket mod-common -c simply -i hatchify vroomy -pr"
	os.Args = strings.Split(input, " ")

	command, err := getCommand()

	test := simply.Target(err, context, "Error should not exist")
	result := test.Assert().Equals(nil)
	test.Validate(result)

	test = simply.Target(command, context, "Command should exist")
	result = test.Assert().DoesNotEqual(nil)
	test.Validate(result)

	test = simply.Target(command.Action, context, "Action should be <sync>")
	result = test.Equals("sync")
	test.Validate(result)

	test = simply.Target(len(command.Arguments), context, "Arguments should have 2 elements")
	result = test.Equals(2)
	test.Validate(result)

	test = simply.Target(command.Arguments[0].Name, context, "Argument[0] should be mod-common")
	result = test.Equals("mod-common")
	test.Validate(result)

	test = simply.Target(command.Arguments[1].Name, context, "Argument[1] should be simply")
	result = test.Equals("simply")
	test.Validate(result)

	test = simply.Target(len(command.Flags), context, "Flags should have 3 elements")
	result = test.Equals(5)
	test.Validate(result)

	includes, _ := command.Flags["-include"].Value.([]string)
	test = simply.Target(includes[0], context, "-i[0] should be test1")
	result = test.Equals("test1")
	test.Validate(result)
	test = simply.Target(includes[1], context, "-i[1] should be test2")
	result = test.Equals("test2")
	test.Validate(result)
	test = simply.Target(includes[2], context, "-i[2] should be hatchify")
	result = test.Equals("hatchify")
	test.Validate(result)
	test = simply.Target(includes[3], context, "-i[3] should be vroomy")
	result = test.Equals("vroomy")
	test.Validate(result)

	branch, _ := command.Flags["-branch"].Value.(string)
	test = simply.Target(branch, context, "-b should be JIRA-Ticket")
	result = test.Equals("JIRA-Ticket")
	test.Validate(result)

	nameOnly, _ := command.Flags["-name-only"].Value.(bool)
	test = simply.Target(nameOnly, context, "-name should be true")
	result = test.Equals(true)
	test.Validate(result)

	pr, _ := command.Flags["-pull-request"].Value.(bool)
	test = simply.Target(pr, context, "-pr should be true")
	result = test.Equals(true)
	test.Validate(result)

	commit, _ := command.Flags["-commit"].Value.(bool)
	test = simply.Target(commit, context, "-c should be true")
	result = test.Equals(true)
	test.Validate(result)
}

func TestConfig_1BoolFlag_2FlagArray_1BoolFlag_Cmd_1Flag_1Arg_1BoolFlag_1Arg_2FlagArrayMatchDiffId(context *testing.T) {
	input := "gomu -name -include test1 test2 -pr sync -b JIRA-Ticket mod-common -c simply -i hatchify vroomy"
	os.Args = strings.Split(input, " ")

	command, err := getCommand()

	test := simply.Target(err, context, "Error should not exist")
	result := test.Assert().Equals(nil)
	test.Validate(result)

	test = simply.Target(command, context, "Command should exist")
	result = test.Assert().DoesNotEqual(nil)
	test.Validate(result)

	test = simply.Target(command.Action, context, "Action should be <sync>")
	result = test.Equals("sync")
	test.Validate(result)

	test = simply.Target(len(command.Arguments), context, "Arguments should have 2 elements")
	result = test.Equals(2)
	test.Validate(result)

	test = simply.Target(command.Arguments[0].Name, context, "Argument[0] should be mod-common")
	result = test.Equals("mod-common")
	test.Validate(result)

	test = simply.Target(command.Arguments[1].Name, context, "Argument[1] should be simply")
	result = test.Equals("simply")
	test.Validate(result)

	test = simply.Target(len(command.Flags), context, "Flags should have 3 elements")
	result = test.Equals(5)
	test.Validate(result)

	includes, _ := command.Flags["-include"].Value.([]string)
	test = simply.Target(includes[0], context, "-i[0] should be test1")
	result = test.Equals("test1")
	test.Validate(result)
	test = simply.Target(includes[1], context, "-i[1] should be test2")
	result = test.Equals("test2")
	test.Validate(result)
	test = simply.Target(includes[2], context, "-i[2] should be hatchify")
	result = test.Equals("hatchify")
	test.Validate(result)
	test = simply.Target(includes[3], context, "-i[3] should be vroomy")
	result = test.Equals("vroomy")
	test.Validate(result)

	branch, _ := command.Flags["-branch"].Value.(string)
	test = simply.Target(branch, context, "-b should be JIRA-Ticket")
	result = test.Equals("JIRA-Ticket")
	test.Validate(result)

	nameOnly, _ := command.Flags["-name-only"].Value.(bool)
	test = simply.Target(nameOnly, context, "-name should be true")
	result = test.Equals(true)
	test.Validate(result)

	pr, _ := command.Flags["-pull-request"].Value.(bool)
	test = simply.Target(pr, context, "-pr should be true")
	result = test.Equals(true)
	test.Validate(result)

	commit, _ := command.Flags["-commit"].Value.(bool)
	test = simply.Target(commit, context, "-c should be true")
	result = test.Equals(true)
	test.Validate(result)
}

func TestConfig_1BoolFlag_3FlagArrayCmdMatch_Cmd_1Flag_1Arg_1BoolFlag_1Arg_2FlagArrayMatch_1BoolFlag(context *testing.T) {
	input := "gomu -name -include test1 test2 sync list -b JIRA-Ticket mod-common -c simply -i hatchify vroomy -pr"
	os.Args = strings.Split(input, " ")

	command, err := getCommand()

	test := simply.Target(err, context, "Error should not exist")
	result := test.Assert().Equals(nil)
	test.Validate(result)

	test = simply.Target(command, context, "Command should exist")
	result = test.Assert().DoesNotEqual(nil)
	test.Validate(result)

	test = simply.Target(command.Action, context, "Action should be <list>")
	result = test.Equals("list")
	test.Validate(result)

	test = simply.Target(len(command.Arguments), context, "Arguments should have 2 elements")
	result = test.Equals(2)
	test.Validate(result)

	test = simply.Target(command.Arguments[0].Name, context, "Argument[0] should be mod-common")
	result = test.Equals("mod-common")
	test.Validate(result)

	test = simply.Target(command.Arguments[1].Name, context, "Argument[1] should be simply")
	result = test.Equals("simply")
	test.Validate(result)

	test = simply.Target(len(command.Flags), context, "Flags should have 5 elements")
	result = test.Equals(5)
	test.Validate(result)

	includes, _ := command.Flags["-include"].Value.([]string)
	test = simply.Target(includes[0], context, "-i[0] should be test1")
	result = test.Equals("test1")
	test.Validate(result)
	test = simply.Target(includes[1], context, "-i[1] should be test2")
	result = test.Equals("test2")
	test.Validate(result)
	test = simply.Target(includes[2], context, "-i[2] should be sync")
	result = test.Equals("sync")
	test.Validate(result)
	test = simply.Target(includes[3], context, "-i[3] should be hatchify")
	result = test.Equals("hatchify")
	test.Validate(result)
	test = simply.Target(includes[4], context, "-i[4] should be vroomy")
	result = test.Equals("vroomy")
	test.Validate(result)

	branch, _ := command.Flags["-branch"].Value.(string)
	test = simply.Target(branch, context, "-b should be JIRA-Ticket")
	result = test.Equals("JIRA-Ticket")
	test.Validate(result)

	nameOnly, _ := command.Flags["-name-only"].Value.(bool)
	test = simply.Target(nameOnly, context, "-name should be true")
	result = test.Equals(true)
	test.Validate(result)

	pr, _ := command.Flags["-pull-request"].Value.(bool)
	test = simply.Target(pr, context, "-pr should be true")
	result = test.Equals(true)
	test.Validate(result)

	commit, _ := command.Flags["-commit"].Value.(bool)
	test = simply.Target(commit, context, "-c should be true")
	result = test.Equals(true)
	test.Validate(result)
}
