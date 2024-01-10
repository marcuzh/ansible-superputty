package cmd

type CommandExecutor interface {
	Execute(inventoryFile string) error
}

type DefaultCommandExecutor struct{}

func (ce *DefaultCommandExecutor) Execute(_ string) error {
	return nil
}
