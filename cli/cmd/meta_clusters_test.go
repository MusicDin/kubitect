package cmd

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/MusicDin/kubitect/cli/app"
	"github.com/MusicDin/kubitect/cli/cluster"

	"github.com/stretchr/testify/assert"
)

// MockMetaClusters returns a slice of mocked MetaClusters based on the
// provided slice of cluster names. Clusters whose name contains a "local"
// prefix are marked as local.
func MockMetaClusters(t *testing.T, names []string) MetaClusters {
	t.Helper()

	ctxOptions := app.AppContextOptions{Local: false}
	ctx := app.MockAppContext(t, ctxOptions)

	var clusters MetaClusters
	for _, clsName := range names {
		clsIsLocal := strings.HasPrefix(clsName, "local")
		clsPath := path.Join(ctx.ClustersDir(), clsName)

		if clsIsLocal {
			clsPath = path.Join(ctx.LocalClustersDir(), clsName)
		}

		// Create cluster files
		assert.NoError(t, os.MkdirAll(clsPath, os.ModePerm), "Failed to create a mock cluster directory!")

		c := cluster.ClusterMeta{
			Name:       clsName,
			Path:       clsPath,
			Local:      clsIsLocal,
			AppContext: ctx,
		}

		clusters = append(clusters, c)
	}

	return clusters
}

func TestNames(t *testing.T) {
	names := []string{
		"mock-cluster-1",
		"mock-cluster-1",
		"local-cluster-3",
	}

	mcs := MockMetaClusters(t, names)
	assert.ElementsMatch(t, names, mcs.Names())
}

func TestCountByName(t *testing.T) {
	names := []string{
		"mock-cluster-1",
		"mock-cluster-1",
		"local-cluster-3",
	}

	mcs := MockMetaClusters(t, names)
	assert.Equal(t, 0, mcs.CountByName("mock-cluster-0"))
	assert.Equal(t, 2, mcs.CountByName("mock-cluster-1"))
}

func TestFindByName(t *testing.T) {
	mcs := MockMetaClusters(t, []string{"mock-cluster"})

	c := mcs.FindByName("mock-cluster")
	assert.NotNil(t, c, "FindByName failed to return an existing cluster!")
	assert.Equal(t, "mock-cluster", c.Name)

	assert.Nil(t, mcs.FindByName("invalid"))
}

func TestAllClusters(t *testing.T) {
	names := []string{
		"mock-cluster",
		"local-cluster",
	}

	mcs := MockMetaClusters(t, names)

	cs, err := AllClusters(mcs[0].AppContext)
	assert.NoError(t, err)
	assert.Len(t, cs, 2)
}

func TestAllClusters_InvalidClustersDir(t *testing.T) {
	cs, err := AllClusters(app.MockAppContext(t))
	assert.ErrorContains(t, err, "failed to read clusters directory")
	assert.Empty(t, cs)
}
