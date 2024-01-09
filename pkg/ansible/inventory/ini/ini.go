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

type AINIParser interface {
	Parse(r io.Reader) (*aini.InventoryData, error)
}

type DefaultAINIParser struct{}

func (p DefaultAINIParser) Parse(r io.Reader) (*aini.InventoryData, error) {
	return aini.Parse(r)
}

func (f *File) Parse(ainiParser AINIParser, reader io.Reader) error {
	data, err := ainiParser.Parse(reader)
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
