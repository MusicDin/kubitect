#======================================================================================
# Libvirt connection
#======================================================================================

variable "action" {
  type        = string
  description = "Which action has to be done on the cluster (create, upgrade, add_worker, or remove_worker)"
  default     = "create"

  validation {
    condition     = contains(["create", "upgrade", "add_worker", "remove_worker"], var.action)
    error_message = "Variable 'action' is invalid.\nDefault value is \"create\".\nPossible values are: [\"create\", \"upgrade\", \"add_worker\", \"remove_worker\"]."
  }
}

variable "libvirt_provider_uri" {
  type        = string
  description = "Libvirt provider's URI"
  default     = "qemu:///system"
}

variable "libvirt_resource_pool_name" {
  type        = string
  description = "The libvirt resource pool name"
}

variable "libvirt_resource_pool_location" {
  type        = string
  description = "Location where resource pool will be initialized"
  default     = "/var/lib/libvirt/pools/"
}

#======================================================================================
# Network
#======================================================================================

variable "network_name" {
  type        = string
  description = "Network name"
  default     = "k8s-network"
}

variable "network_interface" {
  type        = string
  description = "Network interface used for VMs (cloud-init) and Keepalived"
  default     = "ens3"
}

variable "network_forward_mode" {
  type        = string
  description = "Network forward mode"
  default     = "nat"

  validation {
    condition     = contains(["nat", "route"], var.network_forward_mode)
    error_message = "Variable 'network_forward_mode' is invalid.\nPossible values are: [\"nat\", \"route\"]."
  }
}

variable "network_virtual_bridge" {
  type        = string
  description = "Network virtual bridge"
  default     = "virbr1"
}

variable "network_mac" {
  type        = string
  description = "Network MAC address"
  default     = "52:54:00:4f:e3:88"
}

variable "network_gateway" {
  type        = string
  description = "Network gateway IP address"
  default     = "192.168.113.1"
}

variable "network_mask_bits" {
  type        = number
  description = "Bits used for network"
  default     = 24

  validation {
    condition     = var.network_mask_bits > 0 && var.network_mask_bits <= 32
    error_message = "Valid value for bits used for network is between 1 and 32. (Default value is 24)."
  }
}

variable "network_dhcp_ip_start" {
  type        = string
  description = "DHCP IP range start"
  default     = "192.168.113.2"
}

variable "network_dhcp_ip_end" {
  type        = string
  description = "DHCP IP range end"
  default     = "192.168.113.254"
}

#======================================================================================
# Kubernetes infrastructure
#======================================================================================

#============================#
# General variables          #
#============================#

variable "vm_user" {
  type        = string
  description = "SSH user for VMs"
  default     = "user"
}

variable "vm_ssh_private_key" {
  type        = string
  description = "Location of private ssh key for VMs"
}

variable "vm_ssh_known_hosts" {
  type        = string
  description = "Add virtual machines to SSH known hosts"
  default     = "true"

  validation {
    condition     = contains(["true", "false"], var.vm_ssh_known_hosts)
    error_message = "Variable 'vm_ssh_known_hosts' is invalid.\nPossible values are: [\"true\", \"false\"]."
  }
}

variable "vm_distro" {
  type        = string
  description = "Linux distribution used on VMs. Possible values: [ubuntu, centos, debian]"
  default     = "N/A"
}

variable "vm_image_source" {
  type        = string
  description = "Image source, which can be path on host's filesystem or URL."
}

variable "vm_name_prefix" {
  type        = string
  description = "Prefix added to names of VMs"
  default     = "vm"
}

#============================#
# Load balancer variables    #
#============================#

variable "vm_lb_cpu" {
  type        = number
  description = "The number of vCPU allocated to the HAProxy load balancer"
  default     = 1
}

variable "vm_lb_ram" {
  type        = number
  description = "The amount of RAM allocated to the HAProxy load balancer"
  default     = 4096
}

variable "vm_lb_storage" {
  type        = number
  description = "The amount of disk (in Bytes) allocated to the HAProxy load balancer. Default: 15GB"
  default     = 16106127360
}

variable "vm_lb_macs_ips" {
  type        = map(string)
  description = "MAC (key) and IP (value) addresses of HAProxy load balancer nodes"
}

variable "vm_lb_vip" {
  type        = string
  description = "The IP address of HAProxy load balancer floating VIP"
}

#============================#
# Master nodes variables     #
#============================#

variable "vm_master_cpu" {
  type        = number
  description = "The number of vCPU allocated to the master node"
  default     = 1
}

variable "vm_master_ram" {
  type        = number
  description = "The amount of RAM allocated to the master node"
  default     = 4096
}

variable "vm_master_storage" {
  type        = number
  description = "The amount of disk (in Bytes) allocated to the master node. Default: 15GB"
  default     = 16106127360
}

variable "vm_master_macs_ips" {
  type        = map(string)
  description = "MAC and IP addresses of master nodes"

  validation {
    condition     = length(var.vm_master_macs_ips) > 0
    error_message = "Variable 'vm_master_macs_ips' is invalid.\nAt least one master node should be defined."
  }
}

#============================#
# Worker nodes variables     #
#============================#

variable "vm_worker_cpu" {
  type        = number
  description = "The number of vCPU allocated to the worker node"
  default     = 2
}

variable "vm_worker_ram" {
  type        = number
  description = "The amount of RAM allocated to the worker node"
  default     = 8192
}
variable "vm_worker_storage" {
  type        = number
  description = "The amount of disk (in Bytes) allocated to the worker node. Default: 30GB"
  default     = 32212254720
}
variable "vm_worker_macs_ips" {
  type        = map(string)
  description = "MAC and IP addresses of worker nodes"
}

variable "vm_worker_node_label" {
  type        = string
  description = "Worker node role label"
  default     = ""
}

#======================================================================================
# General kubernetes (k8s) variables
#======================================================================================

variable "k8s_kubespray_url" {
  type        = string
  description = "The Git repository URL to clone Kubespray from"
}

variable "k8s_kubespray_version" {
  type        = string
  description = "The version of Kubespray that will be used to deploy Kubernetes"
}

variable "k8s_version" {
  type        = string
  description = "The version of Kuberenetes that will be deployed"
}

variable "k8s_network_plugin" {
  type        = string
  description = "The overlay network plugin used by Kubernetes cluster"

  validation {
    condition     = contains(["flannel", "weave", "calico", "cilium", "canal", "kube-router"], var.k8s_network_plugin)
    error_message = "Variable 'k8s_network_plugin' is invalid.\nPossible values are: [\"flannel\", \"weave\", \"calico\", \"cilium\", \"canal\", \"kube-router\"]."
  }
}

variable "k8s_dns_mode" {
  type        = string
  description = "The DNS service used by Kubernetes cluster (coredns/kubedns)"

  validation {
    condition     = contains(["coredns", "kubedns"], var.k8s_dns_mode)
    error_message = "Variable 'k8s_dns_mode' is invalid.\nPossible values are: [\"coredns\", \"kubedns\"]."
  }
}

#======================================================================================
# Kubespray addons
#======================================================================================

variable "kubespray_custom_addons_enabled" {
  type        = string
  description = "If enabled, custom addons.yml will be used"
  default     = "false"

  validation {
    condition     = contains(["true", "false"], var.kubespray_custom_addons_enabled)
    error_message = "Variable 'kubespray_custom_addons_enabled' is invalid.\nPossible values are: [\"true\", \"false\"]."
  }
}

variable "kubespray_custom_addons_path" {
  type        = string
  description = "If enabled, custom addons.yml will be used"
  default     = ""
}

variable "k8s_dashboard_enabled" {
  type        = string
  description = "Sets up Kubernetes dashboard if enabled"
  default     = "false"

  validation {
    condition     = contains(["true", "false"], var.k8s_dashboard_enabled)
    error_message = "Variable 'k8s_dashboard_enabled' is invalid.\nPossible values are: [\"true\", \"false\"]."
  }
}

variable "helm_enabled" {
  type        = string
  description = "Sets up Helm if enabled"
  default     = "false"

  validation {
    condition     = contains(["true", "false"], var.helm_enabled)
    error_message = "Variable 'helm_enabled' is invalid.\nPossible values are: [\"true\", \"false\"]."
  }
}

variable "metallb_enabled" {
  type        = string
  description = "Sets up MetalLB if enabled"
  default     = "false"

  validation {
    condition     = contains(["true", "false"], var.metallb_enabled)
    error_message = "Variable 'metallb_enabled' is invalid.\nPossible values are: [\"true\", \"false\"]."
  }
}

variable "metallb_version" {
  type        = string
  description = "MetalLB version"
  default     = ""
}

variable "metallb_port" {
  type        = number
  description = "Kubernetes MetalLB port"
  default     = 7472
}

variable "metallb_cpu_limit" {
  type        = string
  description = "MetalLB pod CPU limit"
  default     = "100m"
}

variable "metallb_mem_limit" {
  type        = string
  description = "MetalLB pod memory (RAM) limit"
  default     = "100Mi"
}

variable "metallb_protocol" {
  type        = string
  description = "MetalLB protocol (layer2/bgp)"
  default     = "layer2"

  validation {
    condition     = contains(["layer2", "bgp"], var.metallb_protocol)
    error_message = "Variable 'metallb_protocol' is invalid.\nPossible values are: [\"layer2\", \"bgp\"]."
  }
}

variable "metallb_ip_range" {
  type        = string
  description = "IP range that MetalLB will use for services of type LoadBalancer"
  default     = ""
}

variable "metallb_peers" {
  type        = list(object({
    peer_ip  = string
    peer_asn = number
    my_asn   = number
  }))
  description = "List of MetalLB peers"
}
