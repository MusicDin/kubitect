#!/bin/sh

#======================================================================================
# terraform-kubespray-kvm (tkk) helper script
#======================================================================================
#
# See 'tkk --help' for help
#

# Script version
TKK_SCRIPT_NAME="tkk"
TKK_SCRIPT_VERSION="0.0.1"

# Other variables
TKK_HOME=${TKK_HOME:-"$HOME/.tkk"}
TKK_CLUSTER_NAME=${TKK_CLUSTER_NAME:-"default"}
TKK_ACTION=""

# Other paths
TKK_CONFIG_PATH=""
TKK_REQUIREMENTS_PATH="requirements.txt"
TKK_CLUSTER_PATH=""

# Text colors..
TKK_TEXT_RED="\033[0;31m"
TKK_TEXT_GREEN="\033[0;32m"
TKK_TEXT_BOLD="\033[1m"
TKK_TEXT_ITALIC="\033[3m"
TKK_TEXT_UNDERLINE="\033[4m"
TKK_TEXT_CLEAR="\033[0m"

# Options
TKK_OPTIONS_SHORT=a:,c:,l,h,v
TKK_OPTIONS_LONG=action:,config:,local,help,version,cluster:,auto-approve
TKK_OPTION_ACTION=""
TKK_OPTION_LOCAL=""
TKK_OPTION_CLUSTER_NAME=""
TKK_OPTION_AUTO_APPROVE=""

#
# Prints bold message.
#
__print_bold() {
    echo "${TKK_TEXT_BOLD}${1}${TKK_TEXT_CLEAR}"
}

#
# Prints italic message.
#
__print_italic() {
    echo "${TKK_TEXT_ITALIC}${1}${TKK_TEXT_CLEAR}"
}

#
# Prints underline message.
#
__print_underline() {
    echo "${TKK_TEXT_UNDERLINE}${1}${TKK_TEXT_CLEAR}"
}


#
# Prints green ok status message.
#
__print_ok() {
    stamp=$(__print_bold "OK")
    echo "[ ${TKK_TEXT_GREEN}${stamp}${TKK_TEXT_CLEAR} ] $1"
}

#
# Prints red error message.
#
__print_err() {
    stamp=$(__print_bold "ERROR")
    echo "[ ${TKK_TEXT_RED}${stamp}${TKK_TEXT_CLEAR} ] $1"
}

#
# Print an error and exit the script.
#
__err() {
    __print_err "$1\n"
    exit 1
}

#
# Set cluster name and path.
#
__set_cluster(){

    if [ -n "$TKK_OPTION_CLUSTER_NAME" ]; then
        TKK_CLUSTER_NAME="$TKK_OPTION_CLUSTER_NAME"
    fi

    # Local dir
    if [ -n "$TKK_OPTION_LOCAL" ]; then
        TKK_CLUSTER_PATH="."
        __print_ok "--cluster=."
        return
    fi
    
    TKK_CLUSTER_PATH="$TKK_HOME/clusters/$TKK_CLUSTER_NAME"
    __print_ok "--cluster=$TKK_CLUSTER_NAME"
}

#
# Set custom config path.
#
__set_config_path() {
    if [ -n $1 ]; then
        TKK_CONFIG_PATH="$(cd $(dirname $1) && pwd)/$(basename $1)"
        __print_ok "--config=$TKK_CONFIG_PATH"
    fi
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
# List clusters in "$TKK_HOME/clusters" directory.
#
__list_clusters() {

    clusters_path="$TKK_HOME/clusters"

    if [ -z "$(ls -d "$clusters_path/"* 2>/dev/null)" ]; then
        echo "There is no initialized clusters."
        return
    fi

    echo "Clusters: "
    for dir in "$clusters_path/"*; do
        if [ -d "$dir" ]; then
            echo "- $(basename "$dir")"
        fi
    done
}

#
# Initialize a cluster.
# Prepare a cluster directory.
#
__init_cluster() {

    # Skip init phase for local actions
    if [ -n $TKK_OPTION_LOCAL ]; then
        return
    fi

    local tkk_url="https://github.com/MusicDin/terraform-kvm-kubespray"
    local tkk_version="feature/multiple-servers"

    # Fail if python3 is not installed.
    command -v python3 >/dev/null 2>&1 \
        || __err "Python3 needs to be installed."

    # Fail if virtualenv is not installed.
    command -v virtualenv >/dev/null 2>&1 \
        || __err "Virtualenv (pip3 install virtualenv) needs to be installed."

    mkdir -p "$TKK_CLUSTER_PATH"
    cd "$TKK_CLUSTER_PATH"

    # Clone git project
    git init . --quiet
    git fetch $tkk_url $tkk_version --depth 1 --quiet
    git checkout FETCH_HEAD --quiet
    git reset --hard --quiet

    __print_ok "Successfully initialized cluster '$TKK_CLUSTER_NAME'."
}

#
# Create virtual environment and install Ansible with other
# required dependencies.
#
__activate_virtual_env() {

    venv_path="$TKK_CLUSTER_PATH/venv"

    virtualenv -p python3 "$venv_path" \
        && . "$venv_path/bin/activate" \
        && pip3 install -r $TKK_REQUIREMENTS_PATH
}

#
# Generates main.tf file for the given cluster config file.
#
__generate_config() {

    __activate_virtual_env

    local extra_args

    if [ -n $TKK_OPTION_LOCAL ]; then
        extra_args="--extra-vars tkk_cluster_path=$(pwd)"
    fi

    cd "$TKK_CLUSTER_PATH/ansible/tkk"
    ansible-playbook "tkk.yaml" \
        --inventory "hosts.ini" \
        --extra-vars "config_path=$TKK_CONFIG_PATH" \
        $extra_args \
        --tags "apply" \
        || __err "An error has occured during the main.tf generation."

    cd "../.."
    terraform -chdir="$TKK_CLUSTER_PATH" init -upgrade
}

#
# Starts the cluster cretion process.
#
__apply() {
    __init_cluster
    __generate_config
    terraform -chdir="$TKK_CLUSTER_PATH" apply \
       -var action="$TKK_ACTION" \
       -input=false \
       -compact-warnings \
       $TKK_OPTION_AUTO_APPROVE
       #-var config_path="$TKK_CONFIG_PATH" \
}

#
# Modify and plan the configuration.
#
__plan() {
    __generate_config
    # terraform -chdir="$TKK_HOME/clusters/$TKK_CLUSTER_NAME" plan \
    #     -var action="$TKK_ACTION" \
    #     -compact-warnings
        # -var config_path="$TKK_CONFIG_PATH" \
}

#
# Destroy the cluster.
#
__destroy() {
    terraform -chdir="$TKK_HOME/clusters/$TKK_CLUSTER_NAME" destroy \
        -compact-warnings \
        $TKK_OPTION_AUTO_APPROVE
}

#
# Remove cluster directory and it's content.
# Before purging, trigger cluster destruction. 
#
__purge() {
    TKK_OPTION_AUTO_APPROVE="--auto-approve"
    __destroy
    rm -rf "$TKK_HOME/clusters/$TKK_CLUSTER_NAME"
    __print_ok "Successfully purged '$TKK_CLUSTER_NAME' cluster"
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

	$(__print_bold "> $TKK_SCRIPT_NAME - $TKK_SCRIPT_VERSION")

	$(__print_bold "> $(__print_underline "Quick start"):")
	    Run '$TKK_SCRIPT_NAME apply' command to create the cluster.
	    Optionally prepare custom cluster config and provide it 
	    using --config option.

	$(__print_bold "> $(__print_underline "Commands"):")

	    apply           - Modify main.tf and apply new configuration.
	      -a, --action  - Action to be executed on the cluster.
	                      (default: create)
	      -c, --config  - Custom path to the configuration file.
	          --cluster - Specify the cluster to be used.
	                      (default: default)

	    plan            - Creates terraform cluster plan.
	      -a, --action  - Action to be executed on the cluster.
	                      (default: create)
	      -c, --config  - Custom path to the configuration file.

	    destroy         - Destroys the cluster.

	    list
	      clusters      - List initialized clusters.

	$(__print_bold "> $(__print_underline "Global options"):")
	    -h, --help         - Shows help.
	    -v, --version      - Shows script version.
	        --auto-approve - Automatically approves any user
	                         confirmation request.
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
        -v | --version)
            __version
            exit 0
            ;;

        -h | --help)
            __help
            exit 0
            ;;

        -a | --action)
            TKK_OPTION_ACTION=$2
            shift 2
            ;;

        -c | --config)
            __set_config_path $2
            shift 2
            ;;

        --cluster)
            TKK_OPTION_CLUSTER_NAME=$2
            shift 2
            ;;

        --local)
            TKK_OPTION_LOCAL=$1
            shift
            ;;

        --auto-approve)
            TKK_OPTION_AUTO_APPROVE=$1
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
# Commands.
#
case $1 in

    apply)
        __set_cluster
        __set_action
        __apply
        ;;

    plan)
        __set_cluster
        __set_action
    	__plan
    	;;

    destroy)
        __set_cluster
        __destroy
        ;;

    purge)
        __set_cluster
        __purge
        ;;

    ls|list)
        case $2 in
            clusters)
                __list_clusters
                ;;

            *)
                # Temporary..
                __list_clusters
                ;;
        esac
        ;;

    *)
        __help
        ;;
esac
