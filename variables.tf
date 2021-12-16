#======================================================================================
# General configuration
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

  validation {
    condition     = length(var.libvirt_resource_pool_location) != 0
    error_message = "Libvirt resource pool location cannot be empty."
  }
}

#======================================================================================
# Global VM configuration
#======================================================================================

variable "vm_user" {
  type        = string
  description = "Username used to SSH to the VM"
  default     = "user"
}

variable "vm_ssh_private_key" {
  type        = string
  description = "Location of private ssh key for VMs"

  validation {
    condition     = fileexists(var.vm_ssh_private_key) && fileexists("${var.vm_ssh_private_key}.pub")
    error_message = "Invalid path to private and/or public SSH key. \nPrivate key should be on path 'var.vm_ssh_private_key' and public key should be on the same path with suffix '.pub'."
  }
}

variable "vm_ssh_known_hosts" {
  type        = bool
  description = "Add virtual machines to SSH known hosts"
  default     = true
}

variable "vm_distro" {
  type        = string
  description = "Linux distribution used on VMs. (ubuntu, centos, debian, n/a)"
  default     = "N/A"
}

variable "vm_image_source" {
  type        = string
  description = "Image source, which can be path on host's filesystem or URL."

  validation {
    condition     = length(var.vm_image_source) != 0
    error_message = "Virtual machine (VM) image source is missing. Please specify local path or URL to the image."
  }
}

variable "vm_name_prefix" {
  type        = string
  description = "Prefix added to names of VMs"
  default     = "vm"

  validation {
    condition     = length(var.vm_name_prefix) != 0
    error_message = "Virtual machine (VM) name prefix cannot be empty."
  }
}

variable "vm_network_interface" {
  type        = string
  description = "Network interface used by VMs to connect to the network"
  default     = "ens3"
}

#======================================================================================
# Network configuration
#======================================================================================

variable "network_name" {
  type        = string
  description = "Network name"
  default     = "k8s-network"
}

variable "network_mode" {
  type        = string
  description = "Network mode"
  default     = "nat"

  validation {
    condition     = contains(["nat", "route", "bridge"], var.network_mode)
    error_message = "Variable 'network_mode' is invalid.\nPossible values are: [\"nat\", \"route\", \"bridge\"]."
  }
}

variable "network_bridge" {
  type        = string
  description = "Network (virtual) bridge"
  default     = null
}

variable "network_cidr" {
  type        = string
  description = "Network CIDR"
  validation {
    condition     = can(regex("^([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])(.([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])){3}/([1-9]|[1-2][0-9]|3[0-2])$", var.network_cidr))
    error_message = "Invalid network CIDR."
  }
}

variable "network_gateway" {
  type        = string
  description = "Network gateway"
  default     = null
}

variable "network_dns_list" {
  type        = list(string)
  description = "List of DNS servers used by VMs"
  default     = []
}

#======================================================================================
# HAProxy load balancer VMs parameters
#======================================================================================

variable "lb_default_cpu" {
  type        = number
  description = "The default number of vCPU allocated to the HAProxy load balancer"
  default     = 1
}

variable "lb_default_ram" {
  type        = number
  description = "The default amount of RAM allocated to the HAProxy load balancer"
  default     = 4096
}

variable "lb_default_storage" {
  type        = number
  description = "The default amount of disk (in Bytes) allocated to the HAProxy load balancer. Default: 15GB"
  default     = 16106127360
}

variable "lb_vip" {
  type        = string
  description = "HAProxy load balancer virtual IP address (VIP)"

  /*validation {
    condition = (
      cidrhost(var.network_cidr, 0) == cidrhost("${var.lb_vip}/${split("/", var.network_cidr)[1]}", 0)
    )
    error_message = "HAProxy load balancer virtual IP address (VIP) has to be within network CIDR."
  }*/
}

variable "lb_nodes" {
  type = list(object({
    id  = number
    mac = string
    ip  = string
    # Waiting non-experimental release of optional function #
    #mac     = optional(string)
    #ip      = optional(string)
    #cpu     = optional(number)
    #ram     = optional(number)
    #storage = optional(number)
  }))
  description = "HAProxy load balancer nodes configuration"

  validation {
    condition = (
      alltrue([for node in var.lb_nodes : (node.id >= 0 && node.id <= 200)])
      && compact(tolist([for node in var.lb_nodes : node.id])) == distinct(compact(tolist([for node in var.lb_nodes : node.id])))
      && compact(tolist([for node in var.lb_nodes : node.mac])) == distinct(compact(tolist([for node in var.lb_nodes : node.mac])))
      && compact(tolist([for node in var.lb_nodes : node.ip])) == distinct(compact(tolist([for node in var.lb_nodes : node.ip])))
    )
    error_message = "HAProxy load balancer nodes configuration is incorrect. Make sure that:\n - every ID is unique and that it's value is between 0 and 200,\n - every MAC and IP address is unique or null."
  }
}

#======================================================================================
# Master node VMs parameters
#======================================================================================

variable "master_default_cpu" {
  type        = number
  description = "The default number of vCPU allocated to the master node"
  default     = 1
}

variable "master_default_ram" {
  type        = number
  description = "The default amount of RAM allocated to the master node"
  default     = 4096
}

variable "master_default_storage" {
  type        = number
  description = "The default amount of disk (in Bytes) allocated to the master node. Default: 15GB"
  default     = 16106127360
}

variable "master_nodes" {
  type = list(object({
    id  = number
    mac = string
    ip  = string
    # Waiting non-experimental release of optional function #
    #mac     = optional(string)
    #ip      = optional(string)
    #cpu     = optional(number)
    #ram     = optional(number)
    #storage = optional(number)
  }))
  description = "Master nodes configuration"

  validation {
    condition = (
      length(var.master_nodes) % 2 != 0
      && compact(tolist([for node in var.master_nodes : node.id])) == distinct(compact(tolist([for node in var.master_nodes : node.id])))
      && compact(tolist([for node in var.master_nodes : node.mac])) == distinct(compact(tolist([for node in var.master_nodes : node.mac])))
      && compact(tolist([for node in var.master_nodes : node.ip])) == distinct(compact(tolist([for node in var.master_nodes : node.ip])))
    )
    error_message = "Master nodes configuration is incorrect. Make sure that: \n - number of master nodes is odd (not divisible by 2),\n - every ID is unique,\n - every MAC and IP address is unique or null."
  }
}

#======================================================================================
# Worker node VMs parameters
#======================================================================================

variable "worker_default_cpu" {
  type        = number
  description = "The default number of vCPU allocated to the worker node"
  default     = 2
}

variable "worker_default_ram" {
  type        = number
  description = "The default amount of RAM allocated to the worker node"
  default     = 8192
}

variable "worker_default_storage" {
  type        = number
  description = "The default amount of disk (in Bytes) allocated to the worker node. Default: 30GB"
  default     = 32212254720
}

variable "worker_node_label" {
  type        = string
  description = "Worker node role label"
  default     = ""
}

variable "worker_nodes" {
  type = list(object({
    id  = number
    mac = string
    ip  = string
    # Waiting non-experimental release of optional function #
    #mac     = optional(string)
    #ip      = optional(string)
    #cpu     = optional(number)
    #ram     = optional(number)
    #storage = optional(number)
  }))
  description = "Worker nodes configuration"

  validation {
    condition = (
      compact(tolist([for node in var.worker_nodes : node.id])) == distinct(compact(tolist([for node in var.worker_nodes : node.id])))
      && compact(tolist([for node in var.worker_nodes : node.mac])) == distinct(compact(tolist([for node in var.worker_nodes : node.mac])))
      && compact(tolist([for node in var.worker_nodes : node.ip])) == distinct(compact(tolist([for node in var.worker_nodes : node.ip])))
    )
    error_message = "Worker nodes configuration is incorrect. Make sure that:\n - every ID is unique,\n - every MAC and IP address is unique or null."
  }
}

#======================================================================================
# General Kubernetes configuration
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

variable "k8s_copy_kubeconfig" {
  type        = bool
  description = "If enabled, kubeconfig (config/admin.conf) will be copied to ~/.kube directory"
  default     = false
}

#======================================================================================
# Kubespray addons
#======================================================================================

variable "kubespray_custom_addons_enabled" {
  type        = bool
  description = "If enabled, custom addons.yml will be used"
  default     = false
}

variable "kubespray_custom_addons_path" {
  type        = string
  description = "If custom addons are enabled, addons YAML file from this path will be used"
  default     = ""
}

variable "k8s_dashboard_enabled" {
  type        = bool
  description = "Sets up Kubernetes dashboard if enabled"
  default     = false
}

variable "k8s_dashboard_rbac_enabled" {
  type        = bool
  description = "If enabled, Kubernetes dashboard service account will be created"
  default     = false
}

variable "k8s_dashboard_rbac_user" {
  type        = string
  description = "Kubernetes dashboard service account user"
  default     = ""
}

variable "helm_enabled" {
  type        = bool
  description = "Sets up Helm if enabled"
  default     = false
}

variable "local_path_provisioner_enabled" {
  type        = bool
  description = "Sets up Rancher's local path provisioner if enabled"
  default     = false
}

variable "local_path_provisioner_version" {
  type        = string
  description = "Local path provisioner version"
  default     = ""
}

variable "local_path_provisioner_namespace" {
  type        = string
  description = "Namespace in which local path provisioner will be installed"
  default     = "local-path-provisioner"
}

variable "local_path_provisioner_storage_class" {
  type        = string
  description = "Local path provisioner storage class"
  default     = "local-storage"
}

variable "local_path_provisioner_reclaim_policy" {
  type        = string
  description = "Local path provisioner reclaim policy"
  default     = "Delete"

  validation {
    condition     = contains(["Delete", "Retain"], var.local_path_provisioner_reclaim_policy)
    error_message = "Variable 'local_path_provisioner_reclaim_policy' is invalid.\nPossible values are: [\"Delete\", \"Retain\"]."
  }
}

variable "local_path_provisioner_claim_root" {
  type        = string
  description = "Local path provisioner claim root"
  default     = "/opt/local-path-provisioner/"
}

variable "metallb_enabled" {
  type        = bool
  description = "Sets up MetalLB if enabled"
  default     = false
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
  type = list(object({
    peer_ip  = string
    peer_asn = number
    my_asn   = number
  }))
  description = "List of MetalLB peers"
}
