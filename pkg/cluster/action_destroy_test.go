package cluster

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDestroy(t *testing.T) {
	c := MockCluster(t)

	// Create terraform state file
	err := os.MkdirAll(path.Dir(c.TfStatePath()), os.ModePerm)
	require.NoError(t, err)

	err = ioutil.WriteFile(c.TfStatePath(), []byte(""), os.ModePerm)
	require.NoError(t, err)

	assert.NoError(t, c.Destroy())
}

func TestDestroy_NoTfStateFile(t *testing.T) {
	c := MockCluster(t)
	assert.EqualError(t, c.Destroy(), "cluster 'cluster-mock' is already destroyed (or not yet initialized).")
}
