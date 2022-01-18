#!/bin/sh

#======================================================================================
# terraform-kubespray-kvm (tkk) helper script
#======================================================================================

# See 'tkk.sh --help' for help

VERSION="0.0.1"

BASEDIR="$(dirname $0)/.."

#CONFIG_PATH="./cluster.yml"

MAIN_TF_MODIFIER_PATH="$BASEDIR/ansible/main-tf-modifier"
MAIN_TF_MODIFIER_INVENTORY_PATH="$MAIN_TF_MODIFIER_PATH/hosts.ini"
MAIN_TF_MODIFIER_PLAYBOOK_PATH="$MAIN_TF_MODIFIER_PATH/modify-main-tf.yml"

# Trigger main.tf file modification
modifyMainTf() {
    ansible-playbook $MAIN_TF_MODIFIER_PLAYBOOK_PATH \
        --inventory $MAIN_TF_MODIFIER_INVENTORY_PATH
}

# Modify and apply configuration
apply() {
    shift
    modifyMainTf
    terraform -chdir=$BASEDIR apply $@
}

# Modify and plan configuration
plan() {
    shift
    modifyMainTf
    terraform -chdir=$BASEDIR plan $@
}

version() {
	cat <<-EOF
		tkk.sh - $VERSION
	EOF
}

help() {
	cat <<-EOF

		> tkk.sh - $VERSION
        
		  Script is useful when deploying Kubernetes cluster on 
		  multiple physical servers.

		  It triggers Ansible playbook that modifies main.tf file
		  based on cluster configuration and runs terraform apply
		  or plan.

		  Enjoy.

		> How to use:
		  1.) Modify servers in cluster.yml 
		  2.) Run 'tkk.sh apply' or 'tkk.sh plan'

		> Main commands:
		  apply    - Modify main.tf and apply configuration
		  plan     - Modify main.tf and plan configuration.

		> Other commands:
		  -h, --help      - Shows help.
          -v, --version   - Show script version
	EOF
}

# Apply 
if [ "$1" = "apply" ]; then
    apply $@
elif [ "$1" = "plan" ]; then
    plan $@
elif [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    help
elif [ "$1" = "-v" ] || [ "$1" = "--version" ]; then
    version
else
	help
fi