package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcess(t *testing.T) {
	token, err := process("scope=&token_type=bearer&access_token=c20478d263c282beb19451be75dbb342a35a968c")
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, "c20478d263c282beb19451be75dbb342a35a968c", token)
	_, err = process("scope=&token_type=bearer&token=c20478d263c282beb19451be75dbb342a35a968c")
	require.NotNil(t, err)
}
