#!/bin/sh

#======================================================================================
# terraform-kubespray-kvm (tkk) helper script
#======================================================================================
#
# See 'tkk.sh --help' for help
#

# Script version
VERSION="0.0.1"

# Project root directory
ROOTDIR="$(cd $(dirname $0)/.. && pwd)"

# Ansible main.tf modifier
MAIN_TF_MODIFIER_PATH="$ROOTDIR/ansible/main-tf-modifier"
MAIN_TF_MODIFIER_PLAYBOOK_FILE="modify-main-tf.yml"
MAIN_TF_MODIFIER_INVENTORY_FILE="hosts.ini"

# Other paths
VENV_PATH="$ROOTDIR/venv"
REQUIREMENTS_PATH="$ROOTDIR/requirements.txt"
CONFIG_PATH="$ROOTDIR/cluster.yml"

# Colors..
COLOR_RED='\033[0;31m'
COLOR_GREEN='\033[0;32m'
COLOR_CLEAR='\033[0m'

# Options
TKK_OPTIONS_SHORT=c:,h,v
TKK_OPTIONS_LONG=config:,help,version

#
# Prints green ok status message.
#
__print_ok() {
	echo "[ ${COLOR_GREEN}OK${COLOR_CLEAR} ] $1"
}

#
# Prints red error message.
#
__print_err() {
	echo "[ ${COLOR_RED}ERROR${COLOR_CLEAR} ] $1"
}

#
# Print an error and exit the script.
#
__err() {
	__print_err "$1\n"
	exit 1
}

#
# Set custom config path.
#
__set_config_path() {
	CONFIG_PATH="$(cd $(dirname $1) && pwd)/$(basename $1)"
	__print_ok "--config=$CONFIG_PATH\n"
}

#
# Install Ansible with other dependencies within virtualenv
#
__activate_virtual_env() {
	virtualenv -p python3 $VENV_PATH \
		&& . $VENV_PATH/bin/activate \
		&& pip3 install -r $REQUIREMENTS_PATH
}

#
# Trigger main.tf file modification.
#
__generate() {
	cd $MAIN_TF_MODIFIER_PATH
	__activate_virtual_env
	ansible-playbook $MAIN_TF_MODIFIER_PLAYBOOK_FILE \
		--inventory $MAIN_TF_MODIFIER_INVENTORY_FILE \
		--extra-vars "config_path=$CONFIG_PATH" \
		|| __err "An error has occured during main.tf modification."
	terraform -chdir=$ROOTDIR init -upgrade
}

#
# Generate main.tf file that read Terraform variables (terraform.tfvars)
# as an input instead of YAML configuration.
#
__generate_tf() {
	cd $MAIN_TF_MODIFIER_PATH
	__activate_virtual_env
	ansible-playbook $MAIN_TF_MODIFIER_PLAYBOOK_FILE \
		--inventory $MAIN_TF_MODIFIER_INVENTORY_FILE \
		--extra-vars "action_type=generate-tf" \
		|| __err "An error has occured during the reset of the main.tf file."
	terraform -chdir=$ROOTDIR init -upgrade
}

#
# Modify and apply the configuration.
#
__apply() {
	__generate
	terraform -chdir=$ROOTDIR apply \
		-var "config_path=$CONFIG_PATH" \
		-compact-warnings \
		$@
}

#
# Modify and plan the configuration.
#
__plan() {
	__generate
	terraform -chdir=$ROOTDIR plan \
		-var "config_path=$CONFIG_PATH" \
		-compact-warnings \
		$@
}

#
# Print script version.
#
__version() {
	cat <<-EOF
		tkk.sh - $VERSION
	EOF
}

#
# Prints help.
#
__help() {
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

		> Commands:
		  apply       - Modify main.tf and apply new configuration.
		  plan        - Modify main.tf and plan new configuration.
		  generate    - Only generate main.tf.
		  generate-tf - Generate main.tf file that uses Terraform
		                variables (terraform.tfvars) as an input
		                instead of YAML configuration.

		> Other options:
		  -c, --config  - Path to cluster configuration.
		  -h, --help    - Shows help.
		  -v, --version - Shows script version.
	EOF
}


#
# Read options.
#
OPTS=$(getopt \
	--unquoted \
	--options $TKK_OPTIONS_SHORT \
	--longoptions $TKK_OPTIONS_LONG \
	-- "$@") \
	|| __err "Error reading options."

eval set -- "$OPTS"

#
# Set global options.
#
while :; do
	case "$1" in
		-v | --version )
			__version
			exit 0
			shift
			;;

		-h | --help )
			__help
			exit 0
			shift
			;;

		-c | --config)
			__set_config_path $arg
			shift 2
			;;

		--)
			shift
			break
			;;

		*)
			__err "Unexpected option: $1"
			exit 41
			;;
	esac
done

#
# Commands.
#
case $1 in

	apply)
		__apply
		;;

	plan)
		__plan
		;;

	generate)
		__generate
		;;

	generate-tf)
		__generate_tf
		;;

	*)
		__help
		;;
esac
