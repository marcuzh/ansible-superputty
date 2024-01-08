package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

const inventoryFileFlag = "inventory"

type cmdConfig struct {
	inventoryFile string
}

var (
	config = cmdConfig{} //nolint: gochecknoglobals // deliberately allowing
)

type RootCmd struct {
	rootCmd *cobra.Command
	config  *cmdConfig
}

func NewCmd() *RootCmd {
	return &RootCmd{
		config: &config,
		rootCmd: &cobra.Command{
			Use:   "ansible-superputty",
			Short: "A CLI tool to generate SuperPuTTY configuration",
			Long:  `A CLI tool to generate SuperPuTTY configuration`,
		},
	}
}

func (c *RootCmd) SetArgs(args []string) {
	c.rootCmd.SetArgs(args)
}

func (c *RootCmd) SetOut(writer io.Writer) {
	c.rootCmd.SetOut(writer)
}

func (c *RootCmd) SetErr(writer io.Writer) {
	c.rootCmd.SetErr(writer)
}

func (c *RootCmd) Setup() error {
	c.rootCmd.Flags().StringVarP(&c.config.inventoryFile, inventoryFileFlag, "i", "", "path to ansible inventory file")
	err := c.rootCmd.MarkFlagRequired(inventoryFileFlag)
	if err != nil {
		return err
	}
	c.rootCmd.RunE = c.run
	return nil
}

func (c *RootCmd) Execute() error {
	return c.rootCmd.Execute()
}

func (c *RootCmd) run(_ *cobra.Command, _ []string) error {
	return nil
}
