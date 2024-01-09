package ini

import (
	"fmt"
	"io"

	"github.com/relex/aini"
)

const AnsibleHostVar = "ansible_host"

type Host struct {
	Name     string
	Hostname string
}

type File struct {
	Hosts []Host
}

type AINIParser func(r io.Reader) (*aini.InventoryData, error)

func (f *File) Parse(ainiParser AINIParser, reader io.Reader) error {
	data, err := ainiParser(reader)
	if err != nil {
		return fmt.Errorf("ansible inventory file parse error: %w", err)
	}

	for _, iniHost := range data.Hosts {
		f.Hosts = append(f.Hosts, Host{
			Name:     iniHost.Name,
			Hostname: iniHost.InventoryVars[AnsibleHostVar],
		})
	}

	return nil
}
