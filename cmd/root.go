package cmd

import (
	"io"

	"github.com/marcuzh/ansible-superputty/pkg/cmd"

	"github.com/spf13/cobra"
)

const inventoryFileFlag = "inventory"

type cmdFlags struct {
	inventoryFile string
}

type RootCmd struct {
	cobraCmd *cobra.Command
	config   *cmdFlags
	executor cmd.CommandExecutor
}

func NewRootCmd() (*RootCmd, error) {
	flags := cmdFlags{
		inventoryFile: "",
	}

	cobraCmd := newCobraCommand()
	err := configureFlags(&cobraCmd, &flags)
	if err != nil {
		return nil, err
	}

	return &RootCmd{
		executor: &cmd.DefaultCommandExecutor{},
		config:   &flags,
		cobraCmd: &cobraCmd,
	}, nil
}

func (c *RootCmd) WithExecutor(executor cmd.CommandExecutor) *RootCmd {
	c.executor = executor
	return c
}

func (c *RootCmd) SetArgs(args []string) {
	c.cobraCmd.SetArgs(args)
}

func (c *RootCmd) SetOut(writer io.Writer) {
	c.cobraCmd.SetOut(writer)
}

func (c *RootCmd) SetErr(writer io.Writer) {
	c.cobraCmd.SetErr(writer)
}

func (c *RootCmd) Setup() {
	c.cobraCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return c.executor.Execute(c.config.inventoryFile)
	}
}

func (c *RootCmd) Execute() error {
	return c.cobraCmd.Execute()
}

func configureFlags(command *cobra.Command, flags *cmdFlags) error {
	command.Flags().StringVarP(&flags.inventoryFile, inventoryFileFlag, "i", "", "path to ansible inventory file")
	return command.MarkFlagRequired(inventoryFileFlag)
}

func newCobraCommand() cobra.Command {
	command := cobra.Command{
		Use:   "ansible-superputty",
		Short: "A CLI tool to generate SuperPuTTY configuration",
		Long:  `A CLI tool to generate SuperPuTTY configuration`,
	}
	return command
}
