package cmd_test

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/marcuzh/ansible-superputty/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockCmdExecutor struct {
	mock.Mock
}

func (ce *MockCmdExecutor) Execute(inventoryFile string) error {
	args := ce.Called(inventoryFile)
	return args.Error(0)
}

func Test_ExecuteCommand(t *testing.T) {
	tests := map[string]struct {
		args             string
		expectedCmdFlags any
		returnArgs       any
	}{
		"valid short hand args": {
			args:             "-i inventories/test",
			expectedCmdFlags: "inventories/test",
		},
		"valid long hand args": {
			args:             "--inventory inventories/test",
			expectedCmdFlags: "inventories/test",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			bStdOut, bStdErr := bytes.NewBufferString(""), bytes.NewBufferString("")

			rootCmd := setupRootCmd(strings.Split(tc.args, " "), bStdOut, bStdErr)

			mockCmdExecutor := new(MockCmdExecutor)
			mockCmdExecutor.On("Execute", tc.expectedCmdFlags).Return(tc.returnArgs)
			defer mockCmdExecutor.AssertExpectations(t)

			rootCmd.WithExecutor(mockCmdExecutor)

			err := rootCmd.Execute()
			require.NoError(t, err)

			stdOut, err := io.ReadAll(bStdOut)
			require.NoError(t, err)

			stdErr, err := io.ReadAll(bStdErr)
			require.NoError(t, err)

			assert.Empty(t, string(stdOut))
			assert.Empty(t, string(stdErr))

			mockCmdExecutor.AssertExpectations(t)
		})
	}
}

func Test_ExecuteCommand_Errors(t *testing.T) {
	defaultMockBehaviourSetup := func(m *MockCmdExecutor) {}
	executeNotCalledAssertion := func(m *MockCmdExecutor) {
		m.AssertNotCalled(t, "Execute")
	}

	tests := map[string]struct {
		args                   string
		expectedStdErr         string
		mockBehaviourSetup     func(m *MockCmdExecutor)
		mockBehaviourAssertion func(m *MockCmdExecutor)
	}{
		"no flags supplied": {
			args:                   "",
			expectedStdErr:         `Error: required flag(s) "inventory" not set`,
			mockBehaviourSetup:     defaultMockBehaviourSetup,
			mockBehaviourAssertion: executeNotCalledAssertion,
		},
		"unrecognised flag": {
			args:                   "-i inventories/test --foo",
			expectedStdErr:         `Error: unknown flag: --foo`,
			mockBehaviourSetup:     defaultMockBehaviourSetup,
			mockBehaviourAssertion: executeNotCalledAssertion,
		},
		"inventory flag supplied but with no value": {
			args:                   "-i",
			expectedStdErr:         `Error: flag needs an argument: 'i' in -i`,
			mockBehaviourSetup:     defaultMockBehaviourSetup,
			mockBehaviourAssertion: executeNotCalledAssertion,
		},
		"valid args - execution error": {
			args:           "-i inventories/test",
			expectedStdErr: `Error: command execution error`,
			mockBehaviourSetup: func(m *MockCmdExecutor) {
				m.On("Execute", "inventories/test").Return(errors.New("command execution error"))
			},
			mockBehaviourAssertion: func(m *MockCmdExecutor) {
				m.AssertCalled(t, "Execute", "inventories/test")
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			bStdOut, bStdErr := bytes.NewBufferString(""), bytes.NewBufferString("")

			rootCmd := setupRootCmd(strings.Split(tc.args, " "), bStdOut, bStdErr)

			mockCmdExecutor := new(MockCmdExecutor)
			tc.mockBehaviourSetup(mockCmdExecutor)
			defer tc.mockBehaviourAssertion(mockCmdExecutor)

			rootCmd.WithExecutor(mockCmdExecutor)

			err := rootCmd.Execute()
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

func setupRootCmd(args []string, stdOut io.Writer, stdErr io.Writer) *cmd.RootCmd {
	rootCmd := cmd.NewRootCmd()

	rootCmd.SetOut(stdOut)
	rootCmd.SetErr(stdErr)
	rootCmd.SetArgs(args)

	return rootCmd
}
