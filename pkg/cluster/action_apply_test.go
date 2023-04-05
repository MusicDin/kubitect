package cluster

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/MusicDin/kubitect/pkg/env"
	"github.com/MusicDin/kubitect/pkg/models/config"

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

func TestPrepare(t *testing.T) {
	c := MockLocalCluster(t)
	assert.NoError(t, c.prepare())
}

func TestPlan(t *testing.T) {
	c := MockCluster(t)

	assert.NoError(t, c.ApplyNewConfig())
	assert.NoError(t, c.Sync())

	// Make "blocking" change
	ver := fmt.Sprintf("%s.%s", env.ProjectK8sVersions[0], "99")
	c.NewConfig.Kubernetes.Version = config.KubernetesVersion(ver)

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
	ver := fmt.Sprintf("%s.%s", env.ProjectK8sVersions[0], "99")
	c.NewConfig.Kubernetes.Version = config.KubernetesVersion(ver)

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
		config.WorkerInstance{Id: "worker"},
	)

	// Skip required files check
	tmp := env.ProjectRequiredFiles
	env.ProjectRequiredFiles = []string{}
	defer func() { env.ProjectRequiredFiles = tmp }()

	assert.NoError(t, c.Apply(SCALE.String()))
}
