package kubespray

import (
	"cli/tools/virtualenv"
	"cli/ui"
)

type VirtualEnvironments struct {
	MAIN      *virtualenv.VirtualEnv
	KUBESPRAY *virtualenv.VirtualEnv
}

type KubesprayExecutor struct {
	ClusterName string
	ClusterPath string
	K8sVersion  string
	SshUser     string
	SshPKey     string

	Venvs VirtualEnvironments

	Ui *ui.Ui
}

func (e *KubesprayExecutor) Init() error {
	if err := e.KubitectInit(TAG_INIT); err != nil {
		return err
	}

	return e.KubitectHostsSetup()
}
