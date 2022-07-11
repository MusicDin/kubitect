#=====================================================================================
# Evaluate configuration variables
#=====================================================================================

locals {

  # VM user
  user = (
    try(var.config.cluster.nodeTemplate.user, null) != null
    ? var.config.cluster.nodeTemplate.user
    : var.defaults_config.default.user
  )

  # Precedence: kubitect.yaml > infrastructure.yaml > defaults.yaml
  # Note: Cannot be part of locals.distro, as distro section is then "self-referencing"
  os_distro = (
    try(var.config.cluster.nodeTemplate.os.distro, null) != null
    ? var.config.cluster.nodeTemplate.os.distro
    : (
      try(var.infra_config.cluster.nodeTemplate.os.distro, null) != null
      ? var.infra_config.cluster.nodeTemplate.os.distro
      : var.defaults_config.default.distro
    )
  )

  # If user has changed os distro in cluster configuration then values in
  # infra_config must be ignored.
  os_distro_changed = (
    try(var.infra_config.cluster.nodeTemplate.os.distro, null) != null
    ? var.infra_config.cluster.nodeTemplate.os.distro != local.os_distro
    : true
  )

  # Evaluate os configuration
  # Precedence: kubitect.yaml > infrastructure.yaml > defaults.yaml
  os = {
    distro = local.os_distro

    source = (
      try(var.config.cluster.nodeTemplate.os.source, null) != null
      ? var.config.cluster.nodeTemplate.os.source
      : (
        try(var.infra_config.cluster.nodeTemplate.os.source, null) != null && local.os_distro_changed == false
        ? var.infra_config.cluster.nodeTemplate.os.source
        : var.defaults_config.distros[index(var.defaults_config.distros.*.name, local.os_distro)].url
      )
    )

    networkInterface = (
      try(var.config.cluster.nodeTemplate.networkInterface, null) != null
      ? var.config.cluster.nodeTemplate.networkInterface
      : (
        try(var.infra_config.cluster.nodeTemplate.networkInterface, null) != null && local.os_distro_changed == false
        ? var.infra_config.cluster.nodeTemplate.networkInterface
        : var.defaults_config.distros[index(var.defaults_config.distros.*.name, local.os_distro)].networkInterface
      )
    )
  }

  # Evaluate SSH key path
  ssh = {
    privateKeyPath = (
      try(var.config.cluster.nodeTemplate.ssh.privateKeyPath, null) != null
      ? pathexpand(var.config.cluster.nodeTemplate.ssh.privateKeyPath)
      : pathexpand(var.defaults_config.default.ssh.privateKeyPath)
    )
  }
}
