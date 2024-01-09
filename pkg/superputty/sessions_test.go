package superputty_test

import (
	"errors"
	"io"
	"os"
	"testing"

	xml2 "github.com/marcuzh/ansible-superputty/pkg/xml"

	"github.com/marcuzh/ansible-superputty/pkg/superputty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type AlwaysErrorsMarshaller struct{}

func (d AlwaysErrorsMarshaller) MarshalIndent(_ interface{}) ([]byte, error) {
	return nil, errors.New("deliberate marshall error")
}

func TestSessionData_ToXML(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		sessionData *superputty.SessionData
		fixture     string
	}{
		"empty session data": {
			fixture:     "../../testdata/empty-sessions.xml",
			sessionData: &superputty.SessionData{},
		},
	}

	for name, tc := range tests {
		test := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			expected, err := readFixtureFile(test.fixture)
			require.NoError(t, err)

			xml, err := test.sessionData.ToXML(xml2.DefaultXMLMarshaller{})
			require.NoError(t, err)

			assert.Equal(t, string(expected), xml)
		})
	}
}

func TestSessionData_ToXML_ReturnError(t *testing.T) {
	sessionData := superputty.SessionData{}

	xml, err := sessionData.ToXML(AlwaysErrorsMarshaller{})

	require.Error(t, err)
	assert.Empty(t, xml)
}

func readFixtureFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	return io.ReadAll(file)
}
