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
__modify_main_tf() {
	cd $MAIN_TF_MODIFIER_PATH
	__activate_virtual_env
	ansible-playbook $MAIN_TF_MODIFIER_PLAYBOOK_FILE \
		--inventory $MAIN_TF_MODIFIER_INVENTORY_FILE \
		--extra-vars "config_path=$CONFIG_PATH" \
		|| __err "An error has occured during main.tf modification."
}

#
# Reset main.tf to default (localhost only) configuration.
#
__reset() {
	cd $MAIN_TF_MODIFIER_PATH
	__activate_virtual_env
	ansible-playbook $MAIN_TF_MODIFIER_PLAYBOOK_FILE \
		--inventory $MAIN_TF_MODIFIER_INVENTORY_FILE \
		--extra-vars "action_type=reset" \
		|| __err "An error has occured during the reset of the main.tf file."
}

#
# Modify and apply the configuration.
#
__apply() {
	__modify_main_tf
	terraform -chdir=$ROOTDIR init -upgrade
	terraform -chdir=$ROOTDIR apply \
		-var "config_type=yaml" \
		-var "config_path=$CONFIG_PATH" \
		$@
}

#
# Modify and plan the configuration.
#
__plan() {
	__modify_main_tf
	terraform -chdir=$ROOTDIR init -upgrade
	terraform -chdir=$ROOTDIR plan \
		-var "config_type=yaml" \
		-var "config_path=$CONFIG_PATH" \
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
		  apply    - Modify main.tf and apply new configuration.
		  plan     - Modify main.tf and plan new configuration.
		  generate - Only generate main.tf.
		  reset    - Resets main.tf to default (localhost only) 
		             configuration.

		> Other options:
		  -c, --config  - Path to cluster configuration.
		  -h, --help    - Shows help.
		  -v, --version - Shows script version.
	EOF
}


cmd="$1"
flag=""

#
# Shift first argument (cmd) if it exists.
#
if [ "$#" -gt 0 ]; then
	shift
fi

#
# Check whether custom path (--config) is set.
#
for arg in "$@"; do

	shift

	if [ "$flag" = "--config" ]; then
		__set_config_path $arg
		flag=""
		continue
	fi

	case $arg in
		"-c"|"--config")
			flag="--config"
			;;

		"-c="*|"--config="*)
			__set_config_path $(echo "$arg" | cut -d'=' -f 2)
			;;

		*)
			set -- "$@" $arg
	esac
done

#
# Throw an error if the --config flag is present,
# but the path has not been provided.
#
if [ ! -z "$flag" ]; then
	__err "Option '$flag' requires an argument."
fi

#
# Commands.
#
case $cmd in
	"-h"|"--help")
		__help
		;;

	"-v"|"--version")
		__version
		;;

	"apply")
		__apply $@
		;;

	"plan")
		__plan $@
		;;

	"generate")
		__modify_main_tf
		;;

	"reset")
		__reset
		;;

	*)
		__help
esac

exit 0
