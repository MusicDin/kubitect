#!/bin/sh

#======================================================================================
# terraform-kubespray-kvm (tkk) helper script
#======================================================================================
#
# See 'tkk --help' for help
#

# Script version
TKK_PATH="$(cd $(dirname $0)/.. && pwd)"
TKK_SCRIPT_NAME="tkk"
TKK_SCRIPT_VERSION="0.0.1"

# Ansible main.tf modifier
TKK_MAIN_TF_MODIFIER_PATH="$TKK_PATH/ansible/main-tf-modifier"
TKK_MAIN_TF_MODIFIER_PLAYBOOK_FILE="modify-main-tf.yml"
TKK_MAIN_TF_MODIFIER_INVENTORY_FILE="hosts.ini"

# Other paths
TKK_VENV_PATH="$TKK_PATH/venv"
TKK_REQUIREMENTS_PATH="$TKK_PATH/requirements.txt"
TKK_CONFIG_PATH="$TKK_PATH/cluster.yaml"

# Other variables
TKK_ACTION=""

# Colors..
TKK_COLOR_RED='\033[0;31m'
TKK_COLOR_GREEN='\033[0;32m'
TKK_COLOR_CLEAR='\033[0m'

# Options
TKK_OPTIONS_SHORT=a:,c:,h,v
TKK_OPTIONS_LONG=action:,config:,help,version
TKK_OPTION_ACTION=""
TKK_OPTION_TFVARS="false"

#
# Prints green ok status message.
#
__print_ok() {
    echo "[ ${TKK_COLOR_GREEN}OK${TKK_COLOR_CLEAR} ] $1"
}

#
# Prints red error message.
#
__print_err() {
    echo "[ ${TKK_COLOR_RED}ERROR${TKK_COLOR_CLEAR} ] $1"
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
    TKK_CONFIG_PATH="$(cd $(dirname $1) && pwd)/$(basename $1)"
    __print_ok "--config=$TKK_CONFIG_PATH"
}

#
# Set project path. If project path environment variable 
# is not, then set project path to current directory (".").
#
__set_project_path() {
    if [ -z $TKK_PATH]; then
        TKK_PATH="$(cd $(dirname .) && pwd)"
    fi
    __print_ok "Project path set to: $TKK_PATH"
}

#
# Set action to be executed on the cluster.
#
__set_action() {
    case $TKK_OPTION_ACTION in
        upgrade)
            TKK_ACTION="upgrade"
            ;;

        add_worker)
            TKK_ACTION="add_worker"
            ;;

        remove_worker)
            TKK_ACTION="remove_worker"
            ;;

        # Default action.
        ""|create)
            TKK_ACTION="create"
            ;;
        
        *)
            __err "Unsupported action '$TKK_OPTION_ACTION'"
    esac
    __print_ok "--action=$TKK_ACTION"
}

#
# Install Ansible with other dependencies within virtualenv.
#
__activate_virtual_env() {
    virtualenv -p python3 $TKK_VENV_PATH \
        && . $TKK_VENV_PATH/bin/activate \
        && pip3 install -r $TKK_REQUIREMENTS_PATH
}

#
# Trigger main.tf file modification.
#
__create_config() {
    cd $TKK_MAIN_TF_MODIFIER_PATH
    __activate_virtual_env
    ansible-playbook $TKK_MAIN_TF_MODIFIER_PLAYBOOK_FILE \
        --inventory $TKK_MAIN_TF_MODIFIER_INVENTORY_FILE \
        --extra-vars "config_path=$TKK_CONFIG_PATH" \
        || __err "An error has occured during main.tf modification."
    terraform -chdir=$TKK_PATH init -upgrade
}

#
# Generate main.tf file that read Terraform variables (terraform.tfvars)
# as an input instead of YAML configuration.
#
__create_config_tfvars() {
    cd $TKK_MAIN_TF_MODIFIER_PATH
    __activate_virtual_env
    ansible-playbook $TKK_MAIN_TF_MODIFIER_PLAYBOOK_FILE \
        --inventory $TKK_MAIN_TF_MODIFIER_INVENTORY_FILE \
        --extra-vars "action_type=generate-tf" \
        || __err "An error has occured during the reset of the main.tf file."
    terraform -chdir=$TKK_PATH init \
        -upgrade
}

#
# Modify and apply the configuration.
#
__apply() {
    __create_config
    terraform -chdir=$TKK_PATH apply \
        -var action="$TKK_ACTION" \
        -var config_path="$TKK_CONFIG_PATH" \
        -compact-warnings
}

#
# Modify and plan the configuration.
#
__plan() {
    __create_config
    terraform -chdir=$TKK_PATH plan \
        -var action="$TKK_ACTION" \
        -var config_path="$TKK_CONFIG_PATH" \
        -compact-warnings
}

#
# Print script version.
#
__version() {
    cat <<-EOF
		$TKK_SCRIPT_NAME - $TKK_SCRIPT_VERSION
	EOF
}

#
# Prints help.
#
__help() {
	cat <<-EOF

	> $TKK_SCRIPT_NAME - $TKK_SCRIPT_VERSION

	    Script is useful when deploying Kubernetes cluster on 
	    multiple hosts.

	    It triggers Ansible playbook that modifies 'main.tf' 
	    file based on the cluster configuration and then runs 
	    appropriate terraform command.

	    Enjoy.

	> Quick start:
	    1.) Modify hosts section in cluster.yml file.
	    2.) Run '$TKK_SCRIPT_NAME apply' command to create the cluster.

	> Commands:
	    apply         - Modify main.tf and apply new configuration.
	    -a, --action  - Action to be executed on the cluster.
	                    (default: create)
	    -c, --config  - Custom path to the configuration file.

	    create        - Only generate main.tf.
	        --tfvars  - Generate main.tf file that uses Terraform
	                    variables (terraform.tfvars) as an input
	                    instead of YAML configuration.

	> Global options:
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

        -a | --action)
            TKK_OPTION_ACTION=$2
            shift 2
            ;;

        -c | --config)
            __set_config_path $2
            shift 2
            ;;

        --tfvars)
            TKK_OPTION_TFVARS="true"
            shift
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
# Prepare environment
#
#__set_project_path

#
# Commands.
#
case $1 in

    apply)
        __set_action
        __apply
        ;;

    plan)
        __set_action
    	__plan
    	;;

    create)
        if [ "$TKK_OPTION_TFVARS" = "true" ]; then
            __create_config_tfvars
        else
            __create_config
        fi
        ;;

    *)
        __help
        ;;
esac
