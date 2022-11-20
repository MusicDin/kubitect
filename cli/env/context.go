package env

import (
	"os"
	"path"
	"path/filepath"
)

const (
	EnvHome            = "KUBITECT_HOME"
	DefaultHomeDir     = ".kubitect"
	DefaultShareDir    = "share"
	DefaultClustersDir = "clusters"
)

type ContextOptions struct {
	// Local deployment. Use working dir as project home dir.
	Local bool
}

type Context struct {
	local      bool
	workingDir string
	homeDir    string
}

func (c *Context) Local() bool {
	return c.local
}

func (c *Context) ShowTerraformPlan() bool {
	return false
}

func (c *Context) WorkingDir() string {
	return c.workingDir
}

func (c *Context) HomeDir() string {
	return c.homeDir
}

func (c *Context) ShareDir() string {
	return path.Join(c.homeDir, DefaultShareDir)
}

func (c *Context) ClustersDir() string {
	return filepath.Join(c.homeDir, DefaultClustersDir)
}

func (c *Context) LocalClustersDir() string {
	return filepath.Join(c.workingDir, DefaultHomeDir, DefaultClustersDir)
}

func (o *ContextOptions) Context() *Context {
	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	hd := filepath.Join(wd, DefaultHomeDir)

	if !o.Local {
		userHomeDir, err := os.UserHomeDir()

		if err != nil {
			panic(err)
		}

		def := filepath.Join(userHomeDir, DefaultHomeDir)
		hd = envVar(EnvHome, def)
	}

	c := Context{
		homeDir:    hd,
		workingDir: wd,
		local:      o.Local,
	}

	return &c
}

func envVar(key, def string) string {
	v, ok := os.LookupEnv(key)

	if ok {
		return v
	}

	return def
}
