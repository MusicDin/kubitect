package cluster

import (
	"cli/config/modelconfig"
	"cli/env"
	"cli/tools/git"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToApplyAction(t *testing.T) {
	a, err := ToApplyActionType("create")
	assert.Equal(t, CREATE, a)
	assert.NoError(t, err)

	a, _ = ToApplyActionType(UPGRADE.String())
	assert.Equal(t, UPGRADE, a)

	a, _ = ToApplyActionType(string(SCALE))
	assert.Equal(t, SCALE, a)
}

func TestToApplyAction_Invalid(t *testing.T) {
	_, err := ToApplyActionType("invalid")
	assert.EqualError(t, err, "unknown cluster action: invalid")
}

func TestVerifyClusterDir_ClusterNotExists(t *testing.T) {
	c := MockCluster(t)
	assert.ErrorContains(t, verifyClusterDir(c.Path), "cluster does not exist")
}

func TestVerifyClusterDir_MissingFiles(t *testing.T) {
	c := MockCluster(t)
	assert.NoError(t, os.MkdirAll(c.Path, os.ModePerm))
	assert.Equal(t, NewInvalidClusterDirError(env.ProjectRequiredFiles), verifyClusterDir(c.Path))
}

func TestCloneAndCopyReqFiles_MissingFiles(t *testing.T) {
	tmpDir := t.TempDir()
	proj := git.NewGitProject(env.ConstProjectUrl, "v1.0.0")

	err := cloneAndCopyReqFiles(proj, path.Join(tmpDir, "tmp"), tmpDir)
	assert.ErrorContains(t, err, "Missing files")
}

func TestCloneAndCopyReqFiles(t *testing.T) {
	tmpDir := t.TempDir()
	proj := git.NewGitProject(env.ConstProjectUrl, "master")

	err := cloneAndCopyReqFiles(proj, path.Join(tmpDir, "tmp"), tmpDir)
	assert.NoError(t, err)
}

func TestCloneAndCopyReqFiles_InvalidURL(t *testing.T) {
	tmpDir := t.TempDir()
	proj := git.NewGitProject("invalid", "master")

	err := cloneAndCopyReqFiles(proj, path.Join(tmpDir, "tmp"), tmpDir)
	assert.ErrorContains(t, err, "git clone: failed to clone project")
}

func TestGenerateMissingKeys(t *testing.T) {
	c := MockCluster(t)

	// Unset PrivateKeyPath to force generating SSH keys.
	c.NewConfig.Cluster.NodeTemplate.SSH.PrivateKeyPath = ""
	assert.NoError(t, c.generateSshKeys())

	// Keys should not be regenerated since files exist
	timeout := time.After(10 * time.Second)
	done := make(chan bool)
	go func() {
		assert.NoError(t, c.generateSshKeys())
		done <- true
	}()

	select {
	case <-timeout:
		assert.Fail(t, "Keys should not be recreated after being generated")
	case <-done:
	}
}

func TestGenerateMissingKeys_PKPathProvided(t *testing.T) {
	c := MockCluster(t)
	assert.NoError(t, c.generateSshKeys())
}

func TestPrepare_MissingFiles(t *testing.T) {
	c := MockCluster(t)

	// Remove 1 required file
	assert.NoError(t, os.RemoveAll(path.Join(c.Path, env.ProjectRequiredFiles[0])))

	assert.ErrorContains(t, c.prepare(), "is missing some required files")
}

func TestPrepare_MissingFiles_LocalCluster(t *testing.T) {
	c := MockLocalCluster(t)

	// Remove 1 required file
	assert.NoError(t, os.RemoveAll(path.Join(c.Path, env.ProjectRequiredFiles[0])))

	assert.ErrorContains(t, c.prepare(), "is missing some required files")
}

func TestPlan(t *testing.T) {
	c := MockCluster(t)

	assert.NoError(t, c.ApplyNewConfig())
	assert.NoError(t, c.Sync())

	// Make "blocking" change
	c.NewConfig.Kubernetes.Version = modelconfig.Version("v1.2.3")

	_, err := c.plan(SCALE)
	assert.EqualError(t, err, "Aborted. Configuration file contains errors.")
}

func TestApply_Create(t *testing.T) {
	c := MockCluster(t)

	// Skip required files check
	tmp := env.ProjectRequiredFiles
	env.ProjectRequiredFiles = []string{}
	defer func() { env.ProjectRequiredFiles = tmp }()

	assert.NoError(t, c.Apply(CREATE.String()))
}

func TestApply_Upgrade_AskToCreate(t *testing.T) {
	c := MockCluster(t)

	// Skip required files check
	tmp := env.ProjectRequiredFiles
	env.ProjectRequiredFiles = []string{}
	defer func() { env.ProjectRequiredFiles = tmp }()

	assert.NoError(t, c.Apply(UPGRADE.String()))
	assert.Contains(t, c.Ui().ReadStdout(t), "Cannot upgrade cluster 'cluster-mock'. It has not been created yet.")
}

func TestApply_Upgrade_NoChanges(t *testing.T) {
	c := MockCluster(t)

	assert.NoError(t, c.ApplyNewConfig())
	assert.NoError(t, c.Sync())

	// Skip required files check
	tmp := env.ProjectRequiredFiles
	env.ProjectRequiredFiles = []string{}
	defer func() { env.ProjectRequiredFiles = tmp }()

	assert.NoError(t, c.Apply(UPGRADE.String()))
	assert.Contains(t, c.Ui().ReadStdout(t), "No changes detected.")
}

func TestApply_Upgrade(t *testing.T) {
	c := MockCluster(t)

	assert.NoError(t, c.ApplyNewConfig())
	assert.NoError(t, c.Sync())

	// Make some changes to the new config
	c.NewConfig.Kubernetes.Version = modelconfig.Version("v1.2.3")

	// Skip required files check
	tmp := env.ProjectRequiredFiles
	env.ProjectRequiredFiles = []string{}
	defer func() { env.ProjectRequiredFiles = tmp }()

	assert.NoError(t, c.Apply(UPGRADE.String()))
}

func TestApply_Scale(t *testing.T) {
	c := MockCluster(t)

	assert.NoError(t, c.ApplyNewConfig())
	assert.NoError(t, c.Sync())

	// Append worker node
	c.NewConfig.Cluster.Nodes.Worker.Instances = append(
		c.NewConfig.Cluster.Nodes.Worker.Instances,
		modelconfig.WorkerInstance{Id: "worker"},
	)

	// Skip required files check
	tmp := env.ProjectRequiredFiles
	env.ProjectRequiredFiles = []string{}
	defer func() { env.ProjectRequiredFiles = tmp }()

	assert.NoError(t, c.Apply(SCALE.String()))
}
