#!/bin/sh

#======================================================================================
# terraform-kubespray-kvm (tkk) helper script
#======================================================================================

# See 'tkk.sh --help' for help

VERSION="0.0.1"

ROOTDIR="$(cd $(dirname $0)/.. && pwd)"

MAIN_TF_MODIFIER_PATH="$ROOTDIR/ansible/main-tf-modifier"
MAIN_TF_MODIFIER_PLAYBOOK_FILE="modify-main-tf.yml"
MAIN_TF_MODIFIER_INVENTORY_FILE="hosts.ini"


err() {
	echo "Error: $1"
	exit 1
}

# Trigger main.tf file modification
modifyMainTf() {
	cd $MAIN_TF_MODIFIER_PATH
	ansible-playbook $MAIN_TF_MODIFIER_PLAYBOOK_FILE \
		--inventory $MAIN_TF_MODIFIER_INVENTORY_FILE \
		|| err "An error has occured during main.tf modification."
}

# Sets main.tf to default (localhost only) configuration
reset() {
	cd $MAIN_TF_MODIFIER_PATH
	ansible-playbook $MAIN_TF_MODIFIER_PLAYBOOK_FILE \
		--inventory $MAIN_TF_MODIFIER_INVENTORY_FILE \
		--extra-vars 'action_type=reset' \
		|| err "An error has occured during the reset of main.tf file."
}

# Modify and apply configuration
apply() {
	shift
	modifyMainTf
	terraform -chdir=$ROOTDIR init -upgrade
	terraform -chdir=$ROOTDIR apply $@
}

# Modify and plan configuration
plan() {
	shift
	modifyMainTf
	terraform -chdir=$ROOTDIR init -upgrade
	terraform -chdir=$ROOTDIR plan $@
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
		  1.) Modify servers section in cluster.yml file.
		  2.) Run 'sh tkk.sh apply' or 'sh tkk.sh plan'.

		> Main commands:
		  apply    - Modify main.tf and apply new configuration.
		  plan     - Modify main.tf and plan new configuration.
		  generate - Only generate main.tf.
		  reset    - Resets main.tf to default (localhost only) 
		             configuration.

		> Other commands:
		  -h, --help      - Shows help.
		  -v, --version   - Show script version
	EOF
}


if [ "$1" = "apply" ]; then
	apply $@
elif [ "$1" = "plan" ]; then
	plan $@
elif [ "$1" = "generate" ]; then
	modifyMainTf
elif [ "$1" = "reset" ]; then
	reset
elif [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
	help
elif [ "$1" = "-v" ] || [ "$1" = "--version" ]; then
	version
else
	help
fi
