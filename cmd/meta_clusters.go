package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MusicDin/kubitect/pkg/app"
	"github.com/MusicDin/kubitect/pkg/cluster"
)

type MetaClusters []cluster.ClusterMeta

func (cs MetaClusters) Names() []string {
	var names []string
	for _, c := range cs {
		names = append(names, c.Name)
	}
	return names
}

func (cs MetaClusters) FindByName(name string) *cluster.ClusterMeta {
	for _, c := range cs {
		if c.Name == name {
			return &c
		}
	}

	return nil
}

func (cs MetaClusters) CountByName(name string) int {
	var i = 0

	for _, c := range cs {
		if c.Name == name {
			i++
		}
	}

	return i
}

// AllClusters returns list of clusters meta from both global (project)
// and local clusters directory (if working directory is a Kubitect project).
func AllClusters(ctx app.AppContext) (MetaClusters, error) {
	cs, err := clusters(ctx, false)

	if err != nil {
		return nil, err
	}

	lcs, err := clusters(ctx, true)

	if err == nil {
		cs = append(cs, lcs...)
	}

	return cs, nil
}

// clusters returns list of clusters meta located in either global (project)
// or local clusters directory.
func clusters(ctx app.AppContext, local bool) (MetaClusters, error) {
	var path string

	if local {
		// Ignore local clusters if the local and global cluster directory
		// paths are the same.
		if ctx.LocalClustersDir() == ctx.ClustersDir() {
			return nil, nil
		}

		path = ctx.LocalClustersDir()
	} else {
		path = ctx.ClustersDir()
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read clusters directory: %v", err)
	}

	var cs MetaClusters

	for _, f := range files {
		if f.IsDir() {
			name := f.Name()

			cs = append(cs, cluster.ClusterMeta{
				AppContext: ctx,
				Name:       name,
				Path:       filepath.Join(path, name),
				Local:      local,
			})
		}
	}

	return cs, nil
}
