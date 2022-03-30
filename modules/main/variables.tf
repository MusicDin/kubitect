#======================================================================================
# General configuration
#======================================================================================

variable "action" {
  type        = string
  description = "Action (create, upgrade)."
  default     = "create"
  nullable    = false

  validation {
    condition     = contains(["create", "upgrade"], var.action)
    error_message = "Variable 'action' is invalid. Possible values are: ['create', 'upgrade']."
  }
}

variable "libvirt_provider_uri" {
  type        = string
  description = "Libvirt provider's URI."
  default     = "qemu:///system"
  nullable    = false
}

variable "hosts_mainResourcePoolPath" {
  type        = string
  description = "Path where main resource pool will be initialized."
  default     = "/var/lib/libvirt/pools/"
  nullable    = false

  validation {
    condition     = length(var.hosts_mainResourcePoolPath) != 0
    error_message = "Main resource pool path cannot be empty."
  }
}

variable "hosts_dataResourcePools" {
  type        = list(object({
    name: string
    path: string
  }))
  description = "Location where main resource pool will be initialized."
  default = []
  nullable    = false

  validation {
    condition = length(var.hosts_dataResourcePools.*.name) == length(distinct(var.hosts_dataResourcePools.*.name))
    error_message = "Duplicate data resource pool name found!\nMake sure that the data resource pool names on the same host are unique."
  }
}

#======================================================================================
# Cluster infrastructure configuration
#======================================================================================

variable "cluster_name" {
  type        = string
  description = "Cluster name used as a prefix for various cluster component names."
  default     = "vm"
  nullable    = false
}

#================================
# Node template
#================================

variable "cluster_nodeTemplate_user" {
  type        = string
  description = "Username used to SSH to the virtual machines."
  default     = "user"
  nullable    = false
}

variable "cluster_nodeTemplate_ssh_privateKeyPath" {
  type        = string
  description = "Location of private SSH key that will be used for virtual machines."

  #validation {
  #  condition     = fileexists(var.cluster_nodeTemplate_ssh_privateKeyPath) && fileexists("${var.cluster_nodeTemplate_ssh_privateKeyPath}.pub")
  #  error_message = "Invalid path to private and/or public SSH key. \n\nPrivate and public key must both exist. Public key should be on the same path as the private key, but with '.pub' suffix."
  #}
}

variable "cluster_nodeTemplate_ssh_addToKnownHosts" {
  type        = bool
  description = "Add virtual machines to SSH known hosts."
  default     = true
  nullable    = false
}

variable "cluster_nodeTemplate_image_source" {
  type        = string
  description = "Image source, which can be path on host's filesystem or URL."

  validation {
    condition     = length(var.cluster_nodeTemplate_image_source) != 0
    error_message = "Virtual machine (VM) image source is missing. Please specify local path or URL to the image."
  }
}

variable "cluster_nodeTemplate_updateOnBoot" {
  type        = bool
  description = "Update system on boot."
  default     = true
  nullable    = false
}

#================================
# Cluster network
#================================

variable "cluster_network_mode" {
  type        = string
  description = "Network mode."
  default     = "nat"
  nullable    = false

  validation {
    condition     = contains(["nat", "route", "bridge"], var.cluster_network_mode)
    error_message = "Variable 'network_mode' is invalid.\nPossible values are: [\"nat\", \"route\", \"bridge\"]."
  }
}

variable "cluster_network_cidr" {
  type        = string
  description = "Network CIDR."

  validation {
    condition     = can(regex("^([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])(.([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])){3}/([1-9]|[1-2][0-9]|3[0-2])$", var.cluster_network_cidr))
    error_message = "Invalid network CIDR."
  }
}

variable "cluster_network_gateway" {
  type        = string
  description = "Network gateway."
  default     = null
}

variable "cluster_network_bridge" {
  type        = string
  description = "Network (virtual) bridge."
  default     = "virbr0"
  nullable    = false
}


variable "cluster_network_dns" {
  type        = list(string)
  description = "List of DNS servers used by virtual machines."
  default     = []
  nullable    = false
}

#======================================================================================
# HAProxy load balancer VMs parameters
#======================================================================================

variable "cluster_nodes_loadBalancer_vip" {
  type        = string
  description = "HAProxy load balancer virtual IP address (VIP)."

  /*
  validation {
    condition = (
      cidrhost(var.network_cidr, 0) == cidrhost("${var.lb_vip}/${split("/", var.network_cidr)[1]}", 0)
    )
    error_message = "HAProxy load balancer virtual IP address (VIP) has to be within network CIDR."
  }
  */
}

variable "cluster_nodes_loadBalancer_default_cpu" {
  type        = number
  description = "The default number of vCPU allocated to the HAProxy load balancer."
  default     = 1
  nullable    = false
}

variable "cluster_nodes_loadBalancer_default_ram" {
  type        = number
  description = "The default amount of RAM (in GiB) allocated to the HAProxy load balancer."
  default     = 4
  nullable    = false
}

variable "cluster_nodes_loadBalancer_default_mainDiskSize" {
  type        = number
  description = "Size of the main disk (in GiB) that is attached to the HAProxy load balancer."
  default     = 16
  nullable    = false
}

variable "cluster_nodes_loadBalancer_instances" {
  type = list(object({
    id      = number
    host    = optional(string)
    mac     = optional(string)
    ip      = optional(string)
    cpu     = optional(number)
    ram     = optional(number)
    mainDiskSize = optional(number)
  }))
  description = "HAProxy load balancer node instances."

  validation {
    condition = (
      alltrue([for node in var.cluster_nodes_loadBalancer_instances : (node.id >= 0 && node.id <= 200)])
      && compact(tolist([for node in var.cluster_nodes_loadBalancer_instances : node.id])) == distinct(compact(tolist([for node in var.cluster_nodes_loadBalancer_instances : node.id])))
      && compact(tolist([for node in var.cluster_nodes_loadBalancer_instances : node.mac])) == distinct(compact(tolist([for node in var.cluster_nodes_loadBalancer_instances : node.mac])))
      && compact(tolist([for node in var.cluster_nodes_loadBalancer_instances : node.ip])) == distinct(compact(tolist([for node in var.cluster_nodes_loadBalancer_instances : node.ip])))
    )
    error_message = "HAProxy load balancer nodes configuration is incorrect. Make sure that:\n - every ID is unique and that it's value is between 0 and 200,\n - every MAC and IP address is unique or null."
  }
}

#======================================================================================
# Master node VMs parameters
#======================================================================================

variable "cluster_nodes_master_default_cpu" {
  type        = number
  description = "The default number of vCPU allocated to the master node."
  default     = 1
  nullable    = false
}

variable "cluster_nodes_master_default_ram" {
  type        = number
  description = "The default amount of RAM (in GiB) allocated to the master node."
  default     = 4
  nullable    = false
}

variable "cluster_nodes_master_default_mainDiskSize" {
  type        = number
  description = "Size of the main disk (in GiB) that is attached to the master node."
  default     = 16
  nullable    = false
}

variable "cluster_nodes_master_default_dataDisks" {
  type        = list(object({
    name: string
    pool: string
    size: number
  }))
  description = "List of additional data disks that are attached to the master node."
  default     = []
  nullable    = false
}

variable "cluster_nodes_master_instances" {
  type = list(object({
    id      = number
    host    = optional(string)
    mac     = optional(string)
    ip      = optional(string)
    cpu     = optional(number)
    ram     = optional(number)
    mainDiskSize = optional(number)
    dataDisks = optional(list(object({ 
      name: string
      pool: string
      size: number
    })))
  }))
  description = "Master node instances (control plane)"

  validation {
    condition = (
      compact(tolist([for node in var.cluster_nodes_master_instances : node.id])) == distinct(compact(tolist([for node in var.cluster_nodes_master_instances : node.id])))
      && compact(tolist([for node in var.cluster_nodes_master_instances : node.mac])) == distinct(compact(tolist([for node in var.cluster_nodes_master_instances : node.mac])))
      && compact(tolist([for node in var.cluster_nodes_master_instances : node.ip])) == distinct(compact(tolist([for node in var.cluster_nodes_master_instances : node.ip])))
      # && length(var.cluster_nodes_master_instances) % 2 != 0
    )
    error_message = "Master nodes configuration is incorrect. Make sure that: \n - number of master nodes is odd (not divisible by 2),\n - every ID is unique,\n - every MAC and IP address is unique or null."
  }
}

#======================================================================================
# Worker node VMs parameters
#======================================================================================

variable "cluster_nodes_worker_default_cpu" {
  type        = number
  description = "The default number of vCPU allocated to the worker node."
  default     = 2
  nullable    = false
}

variable "cluster_nodes_worker_default_ram" {
  type        = number
  description = "The default amount of RAM (in GiB) allocated to the worker node."
  default     = 8
  nullable    = false
}

variable "cluster_nodes_worker_default_mainDiskSize" {
  type        = number
  description = "Size of the main disk (in GiB) that is attached to the worker node."
  default     = 32
  nullable    = false
}

variable "cluster_nodes_worker_default_dataDisks" {
  type        = list(object({
    name: string
    pool: string
    size: number
  }))
  description = "List of additional data disks that are attached to the worker node."
  default     = []
  nullable    = false
}

variable "cluster_nodes_worker_default_label" {
  type        = string
  description = "Worker node role label."
  default     = ""
  nullable    = false
}

variable "cluster_nodes_worker_instances" {
  type = list(object({
    id      = number
    host    = optional(string)
    mac     = optional(string)
    ip      = optional(string)
    cpu     = optional(number)
    ram     = optional(number)
    mainDiskSize = optional(number)
    dataDisks = optional(list(object({ 
      name: string
      pool: string
      size: number
    })))
    #label   = optional(string)
  }))
  description = "Worker node instances."

  validation {
    condition = (
      compact(tolist([for node in var.cluster_nodes_worker_instances : node.id])) == distinct(compact(tolist([for node in var.cluster_nodes_worker_instances : node.id])))
      && compact(tolist([for node in var.cluster_nodes_worker_instances : node.mac])) == distinct(compact(tolist([for node in var.cluster_nodes_worker_instances : node.mac])))
      && compact(tolist([for node in var.cluster_nodes_worker_instances : node.ip])) == distinct(compact(tolist([for node in var.cluster_nodes_worker_instances : node.ip])))
    )
    error_message = "Worker nodes configuration is incorrect. Make sure that:\n - every ID is unique,\n - every MAC and IP address is unique or null."
  }
}

#======================================================================================
# General Kubernetes configuration
#======================================================================================

variable "kubernetes_version" {
  type        = string
  description = "The version of Kuberenetes that will be deployed."
}

variable "kubernetes_networkPlugin" {
  type        = string
  description = "The overlay network plugin used by Kubernetes cluster."
  default     = "calico"
  nullable    = false

  validation {
    condition     = contains(["flannel", "weave", "calico", "cilium", "canal", "kube-router"], var.kubernetes_networkPlugin)
    error_message = "Variable 'k8s_network_plugin' is invalid.\nPossible values are: [\"flannel\", \"weave\", \"calico\", \"cilium\", \"canal\", \"kube-router\"]."
  }
}

variable "kubernetes_dnsMode" {
  type        = string
  description = "The DNS service used by Kubernetes cluster (coredns/kubedns)."
  default     = "coredns"
  nullable    = false

  validation {
    condition     = contains(["coredns", "kubedns"], var.kubernetes_dnsMode)
    error_message = "Variable 'k8s_dns_mode' is invalid.\nPossible values are: [\"coredns\", \"kubedns\"]."
  }
}

variable "kubernetes_kubespray_url" {
  type        = string
  description = "The Git repository URL to clone Kubespray from."
  default     = "https://github.com/kubernetes-sigs/kubespray.git"
  nullable    = false
}

variable "kubernetes_kubespray_version" {
  type        = string
  description = "The version of Kubespray that will be used to deploy Kubernetes."
}

variable "kubernetes_kubespray_addons_enabled" {
  type        = bool
  description = "If enabled, configured Kubespray addons will be applied."
  default     = false
  nullable    = false
}

variable "kubernetes_kubespray_addons_configPath" {
  type        = string
  description = "If Kubespray addons are enabled, addons configuration file on this path will be used."
  default     = ""
}

variable "kubernetes_other_copyKubeconfig" {
  type        = bool
  description = "If enabled, kubeconfig (config/admin.conf) will be copied to '~/.kube/' directory."
  default     = false
  nullable    = false
}


#
# Further work required on these variables:
#

#variable "k8s_dashboard_enabled" {
#  type        = bool
#  description = "Sets up Kubernetes dashboard if enabled"
#  default     = false
#  nullable    = false
#}

#variable "k8s_dashboard_rbac_enabled" {
#  type        = bool
#  description = "If enabled, Kubernetes dashboard service account will be created"
#  default     = false
#  nullable    = false
#}

#variable "k8s_dashboard_rbac_user" {
#  type        = string
#  description = "Kubernetes dashboard service account user"
#  default     = "admin"
#  nullable    = false
#}

#======================================================================================
# Other internal variables
#======================================================================================

variable "internal" {
  type = object({
    is_bridge = string
    vm_types = object({
      load_balancer = string
      master        = string
      worker        = string
    })
  })
  description = "Internal variables passed from parent module."
}