package kubespray

func (e *KubesprayExecutor) Upgrade() error {
	if err := e.KubitectInit(TAG_INIT, TAG_KUBESPRAY, TAG_GEN_NODES); err != nil {
		return err
	}

	if err := e.KubitectHostsSetup(); err != nil {
		return err
	}

	if err := e.KubesprayUpgrade(); err != nil {
		return err
	}

	return e.KubitectFinalize()
}
