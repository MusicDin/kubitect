package kubespray

import (
	"path"

	"github.com/MusicDin/kubitect/embed"
	"github.com/MusicDin/kubitect/pkg/models/config"
	"github.com/MusicDin/kubitect/pkg/utils/template"
)

const groupVarsDir = "group_vars"

// fetchTemplate fetches an embedded template with a given name
// and returns it as a string.
//
// It panics if the resource is not found.
func fetchTemplate(name string) string {
	tpl, err := embed.GetTemplate(name + ".tpl")
	if err != nil {
		panic(err)
	}

	return template.TrimTemplate(string(tpl.Content))
}

type KubesprayAllTemplate struct {
	InfraNodes config.Nodes
	configDir  string
}

func NewKubesprayAllTemplate(configDir string, infraNodes config.Nodes) KubesprayAllTemplate {
	return KubesprayAllTemplate{
		configDir:  configDir,
		InfraNodes: infraNodes,
	}
}

func (t KubesprayAllTemplate) Name() string {
	return "all.yaml"
}

func (t KubesprayAllTemplate) Write() error {
	dstPath := path.Join(t.configDir, groupVarsDir, "all", t.Name())
	return template.Write(t, dstPath)
}

func (t KubesprayAllTemplate) Template() string {
	return fetchTemplate(t.Name())
}

type KubesprayK8sClusterTemplate struct {
	Config    config.Config
	configDir string
}

func NewKubesprayK8sClusterTemplate(configDir string, config config.Config) KubesprayK8sClusterTemplate {
	return KubesprayK8sClusterTemplate{
		configDir: configDir,
		Config:    config,
	}
}

func (t KubesprayK8sClusterTemplate) Name() string {
	return "k8s-cluster.yaml"
}

func (t KubesprayK8sClusterTemplate) Write() error {
	dstPath := path.Join(t.configDir, groupVarsDir, "k8s_cluster", t.Name())
	return template.Write(t, dstPath)
}

func (t KubesprayK8sClusterTemplate) Template() string {
	return fetchTemplate(t.Name())
}

type KubesprayAddonsTemplate struct {
	configDir string
	Addons    string
}

func NewKubesprayAddonsTemplate(configDir string, addons string) KubesprayAddonsTemplate {
	return KubesprayAddonsTemplate{
		configDir: configDir,
		Addons:    addons,
	}
}

func (t KubesprayAddonsTemplate) Name() string {
	return "addons.yaml"
}

func (t KubesprayAddonsTemplate) Write() error {
	dstPath := path.Join(t.configDir, groupVarsDir, "k8s_cluster", t.Name())
	return template.Write(t, dstPath)
}

func (t KubesprayAddonsTemplate) Template() string {
	return "{{ .Addons }}"
}

type KubesprayEtcdTemplate struct {
	configDir string
}

func NewKubesprayEtcdTemplate(configDir string) KubesprayEtcdTemplate {
	return KubesprayEtcdTemplate{configDir}
}

func (t KubesprayEtcdTemplate) Name() string {
	return "etcd.yaml"
}

func (t KubesprayEtcdTemplate) Write() error {
	dstPath := path.Join(t.configDir, groupVarsDir, t.Name())
	return template.Write(t, dstPath)
}

func (t KubesprayEtcdTemplate) Template() string {
	return fetchTemplate(t.Name())
}

type HostsTemplate struct {
	configDir string
	Hosts     []config.Host
}

func NewHostsTemplate(configDir string, hosts []config.Host) HostsTemplate {
	return HostsTemplate{
		configDir: configDir,
		Hosts:     hosts,
	}
}

func (t HostsTemplate) Name() string {
	return "hosts.yaml"
}

func (t HostsTemplate) Write() error {
	dstPath := path.Join(t.configDir, t.Name())
	return template.Write(t, dstPath)
}

func (t HostsTemplate) Functions() map[string]interface{} {
	return map[string]interface{}{
		"isRemoteHost": isRemoteHost,
	}
}

// isRemoteHost returns true id host's connection type equals REMOTE.
func isRemoteHost(host config.Host) bool {
	return host.Connection.Type == config.REMOTE
}

func (t HostsTemplate) Template() string {
	return fetchTemplate(t.Name())
}

type NodesTemplate struct {
	configDir   string
	ConfigNodes config.Nodes
	InfraNodes  config.Nodes
}

func NewNodesTemplate(configDir string, configNodes, infraNodes config.Nodes) NodesTemplate {
	return NodesTemplate{
		configDir:   configDir,
		ConfigNodes: configNodes,
		InfraNodes:  infraNodes,
	}
}

func (t NodesTemplate) Name() string {
	return "nodes.yaml"
}

func (t NodesTemplate) Write() error {
	dstPath := path.Join(t.configDir, t.Name())
	return template.Write(t, dstPath)
}

func (t NodesTemplate) Template() string {
	return fetchTemplate(t.Name())
}
