#=====================================================================================
# Evaluate configuration variable
#=====================================================================================

locals {

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

  # Evaluate SSH configuration
  ssh = {
    privateKeyPath = (
      try(var.config.cluster.nodeTemplate.ssh.privateKeyPath, null) != null
      ? pathexpand(var.config.cluster.nodeTemplate.ssh.privateKeyPath)
      : pathexpand(var.defaults_config.default.ssh.privateKeyPath)
    )
  }

}


#=====================================================================================
# SSH Keys
#=====================================================================================

# Generates SSH keys if path to the private key is not provided. #
resource "null_resource" "generate_ssh_keys" {

  count = try(var.config.cluster.nodeTemplate.ssh.privateKeyPath, null) == null ? 1 : 0

  provisioner "local-exec" {

    command = <<-EOF
      dirname $SSH_PK_PATH | xargs mkdir -p
      if [ ! -e $SSH_PK_PATH ] && [ ! -e $\{SSH_PK_PATH\}.pub ]; then \
        ssh-keygen -f $SSH_PK_PATH -q -N ""; \
      fi
    EOF

    environment = {
      SSH_PK_PATH = local.ssh.privateKeyPath
    }
  }

}
