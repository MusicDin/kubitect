TKK_SCRIPT_NAME=tkk
TKK_APPLY_ACTIONS="\
        create\
        upgrade\
        add_worker\
        remove_worker"
 
# Enable programmable completion facilities are enabled.
shopt -s progcomp


#================================================
# Complete function
#================================================

#
# Autocomplete function
#
_tkk_completion() {
    local cur prev firstword lastword complete_words complete_options
 
    # Don't break words at ':' and '='
    COMP_WORDBREAKS=${COMP_WORDBREAKS//[:=]}
 
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    firstword=$(__tkk_get_firstword)
    lastword=$(__tkk_get_lastword)
 
    GLOBAL_COMMANDS="\
        apply\
        create"
 
    GLOBAL_OPTIONS="\
        -h --help\
        -v --version"
 
    APPLY_OPTIONS="\
        -c --config\
        -a --action\
           --auto-approve"

    CREATE_COMMANDS="\
        config"

    CREATE_CONFIG_OPTIONS="\
        --tfvars"

    case "${firstword}" in 

        apply)
            case "${prev}" in
                -c|--config)
                    ;;
                    
                -a|--action)
                    complete_words="$TKK_APPLY_ACTIONS"
                    ;;

                *)
                    complete_options="$APPLY_OPTIONS"
                    ;;
            esac
            ;;

        create)
            case "${lastword}" in
                config)
                    case "${prev}" in	
                        --tfvars)
                            return 0
                            ;;

                        *)
                            complete_options="$CREATE_CONFIG_OPTIONS"
                            ;;
                    esac
                    ;;

                *)
                    complete_words="$CREATE_COMMANDS"
                    ;;
            esac
            ;;
        
        # GLOBAL
        *)
            case "${prev}" in
                *)
                    complete_words="$GLOBAL_COMMANDS"
                    complete_options="$GLOBAL_OPTIONS"
                    ;;
            esac
            ;;
    esac
 
    
    if [[ -z $complete_options ]] && [[ -z $complete_words ]]; then
        # Print filenames if option and word lists are empty
        compopt -o default
        COMPREPLY=()
    
    elif [[ -z $complete_words ]] && [[ $cur == "-"* ]]; then
        # Print options if word list is empty and current word starts with '-'.
        COMPREPLY=( $( compgen -W "$complete_options" -- $cur ))
 
    else
        # Print words
        COMPREPLY=( $( compgen -W "$complete_words" -- $cur ))

    fi
 
    return 0
}
 

#================================================
# Helper functions
#================================================

#
# Path completion
#
__tkk_path_completion() {
    local files=("/some/path/$2"*)
    [[ -e ${files[0]} ]] && COMPREPLY=( "${files[@]##*/}" )
}

#
# Determines the first non-option word of the command line.
#
__tkk_get_firstword() {
    local firstword i
 
    firstword=
    for ((i = 1; i < ${#COMP_WORDS[@]}; ++i)); do
        if [[ ${COMP_WORDS[i]} != -* ]]; then
            firstword=${COMP_WORDS[i]}
            break
        fi
    done
 
    echo $firstword
}

#
# Determines the last non-option word of the command line. 
#
__tkk_get_lastword() {
    local lastword i
 
    lastword=
    for ((i = 1; i < ${#COMP_WORDS[@]} && i < 3; ++i)); do
        if [[ ${COMP_WORDS[i]} != -* ]] && [[ -n ${COMP_WORDS[i]} ]] && [[ ${COMP_WORDS[i]} != $cur ]]; then
            lastword=${COMP_WORDS[i]}
        fi
    done
 
    echo $lastword
}
 

#================================================
# Complete command
#================================================

complete -F _tkk_completion tkk