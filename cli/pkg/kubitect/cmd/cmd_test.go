package cmd

import (
	app2 "github.com/MusicDin/kubitect/cli/pkg/app"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Execute(t *testing.T, cmdFunc any, ctx ...app2.AppContextMock) (string, error) {
	return ExecuteWithArgs(t, cmdFunc, nil, ctx...)
}

func ExecuteWithArgs(t *testing.T, cmdFunc interface{}, args []string, opts ...app2.AppContextMock) (string, error) {
	t.Helper()

	var ctx app2.AppContextMock
	var cmd *cobra.Command

	if len(opts) > 0 {
		ctx = opts[0]
	} else {
		ctx = app2.MockAppContext(t)
	}

	if f, ok := cmdFunc.(func() *cobra.Command); ok {
		cmd = f()
	} else if f, ok := cmdFunc.(func(...app2.AppContextOptions) *cobra.Command); ok {
		cmd = f(ctx.Options())
	} else {
		assert.FailNow(t, "Provided command either does not return *cobra.Command or does not accept ...app.AppContextOptions!")
	}

	stdout := ctx.Ui().Streams().Out().File()
	stderr := ctx.Ui().Streams().Err().File()

	cmd.SetOut(stdout)
	cmd.SetErr(stderr)

	err := cmd.Execute()
	out := string(ctx.Ui().ReadStdout(t))

	return out, err
}

func TestRootCmd_Help(t *testing.T) {
	out, err := Execute(t, NewRootCmd)
	assert.NoError(t, err)
	assert.Contains(t, out, rootLong)
}

func TestExportCmd_Help(t *testing.T) {
	out, err := Execute(t, NewExportCmd)
	assert.NoError(t, err)
	assert.Contains(t, out, exportLong)
}

// func TestListCmd_NoClusters(t *testing.T) {
// 	out, err := Execute(t, NewListCmd)
// 	assert.NoError(t, err)
// 	assert.Contains(t, out, "No clusters initialized yet.")
// }
