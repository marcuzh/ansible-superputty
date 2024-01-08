package cmd_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/marcuzh/ansible-superputty/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ExecuteCommand(t *testing.T) {
	tests := map[string]struct {
		args string
	}{
		"valid short hand args": {
			args: "-i inventories/test",
		},
		"valid long hand args": {
			args: "--inventory inventories/test",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			bStdOut, bStdErr := bytes.NewBufferString(""), bytes.NewBufferString("")

			cmd, err := setupCmd(strings.Split(tc.args, " "), bStdOut, bStdErr)
			require.NoError(t, err)

			err = cmd.Execute()
			require.NoError(t, err)

			stdOut, err := io.ReadAll(bStdOut)
			require.NoError(t, err)

			stdErr, err := io.ReadAll(bStdErr)
			require.NoError(t, err)

			assert.Empty(t, string(stdOut))
			assert.Empty(t, string(stdErr))
		})
	}
}

func Test_ExecuteCommand_Errors(t *testing.T) {
	tests := map[string]struct {
		args           string
		expectedStdErr string
	}{
		"no flags supplied": {
			args:           "",
			expectedStdErr: `Error: required flag(s) "inventory" not set`,
		},
		"unrecognised flag": {
			args:           "-i inventories/test --foo",
			expectedStdErr: `Error: unknown flag: --foo`,
		},
		"inventory flag supplied but with no value": {
			args:           "-i",
			expectedStdErr: `Error: flag needs an argument: 'i' in -i`,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			bStdOut, bStdErr := bytes.NewBufferString(""), bytes.NewBufferString("")

			cmd, err := setupCmd(strings.Split(tc.args, " "), bStdOut, bStdErr)
			require.NoError(t, err)

			err = cmd.Execute()
			require.Error(t, err)

			stdOut, err := io.ReadAll(bStdOut)
			require.NoError(t, err)

			stdErr, err := io.ReadAll(bStdErr)
			require.NoError(t, err)

			assert.True(t, strings.HasPrefix(string(stdOut), "Usage:"))
			assert.Equal(t, tc.expectedStdErr, strings.ReplaceAll(string(stdErr), "\n", ""))
		})
	}
}

func setupCmd(args []string, stdOut io.Writer, stdErr io.Writer) (*cmd.RootCmd, error) {
	command := cmd.NewCmd()
	err := command.Setup()
	if err != nil {
		return nil, err
	}

	command.SetOut(stdOut)
	command.SetErr(stdErr)
	command.SetArgs(args)

	return command, nil
}
