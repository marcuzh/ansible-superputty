package ini_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/marcuzh/ansible-superputty/pkg/ansible/inventory/ini"

	"github.com/relex/aini"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFile_Parse_ReturnExpectedHosts(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		inputData     []string
		expectedHosts []ini.Host
	}{
		"empty data": {
			inputData:     []string{},
			expectedHosts: []ini.Host{},
		},
		"single empty group": {
			inputData:     []string{"[my-empty-group]"},
			expectedHosts: []ini.Host{},
		},
		"single host entry": {
			inputData: []string{"example.com ansible_host=10.10.10.10"},
			expectedHosts: []ini.Host{
				{Name: "example.com", Hostname: "10.10.10.10"},
			},
		},
		"multiple host entry": {
			inputData: []string{
				"example.com ansible_host=10.10.10.10",
				"foo.com ansible_host=127.0.0.1",
			},
			expectedHosts: []ini.Host{
				{Name: "example.com", Hostname: "10.10.10.10"},
				{Name: "foo.com", Hostname: "127.0.0.1"},
			},
		},
		"single host entry with no IP/DNS": {
			inputData: []string{"example.com"},
			expectedHosts: []ini.Host{
				{Name: "example.com"},
			},
		},
		"single host in a group": {
			inputData: []string{"[my-group-name]", "example.com ansible_host=10.10.10.10"},
			expectedHosts: []ini.Host{
				{Name: "example.com", Hostname: "10.10.10.10"},
			},
		},
	}

	for name, tc := range tests {
		test := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			invFile := ini.File{}
			err := invFile.Parse(aini.Parse, strings.NewReader(strings.Join(test.inputData, "\n")))

			require.NoError(t, err)
			assert.ElementsMatchf(t, test.expectedHosts, invFile.Hosts, "expected hosts to be equal")
		})
	}
}

func TestFile_Parse_ReturnError(t *testing.T) {
	alwaysErrorsParser := func(r io.Reader) (*aini.InventoryData, error) {
		return nil, errors.New("always errors")
	}

	invFile := ini.File{}
	err := invFile.Parse(alwaysErrorsParser, strings.NewReader(""))

	require.Error(t, err)
}
