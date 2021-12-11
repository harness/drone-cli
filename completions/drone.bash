declare valid_commands=('build' 'cron' 'log' 'encrypt'
                        'exec' 'info' 'repo' 'user'
                        'secret' 'server' 'queue' 'orgsecret'
                        'autoscale' 'convert' 'lint' 'sign'
                        'jsonnet' 'starlark' 'plugins' 'template'
                        'help' 'h')

declare valid_args=('-t' '--token' '-s' '--server'
                    '--autoscaler' '--help' '-h' '--version'
                    '-v')

declare contains_sub_commands=('build' 'cron' 'log' 'repo'
                               'user' 'secret' 'server' 'queue'
                               'orgsecret' 'autoscale' 'plugins' 'template')

# Commands: build
_drone_build_short=('-h')
_drone_build_long=('ls' 'last' 'info' 'create' 'stop' 'restart' 'approve' 'decline' 'promote' 'rollback' 'queue' '--help')

_drone_build_ls_short=()
_drone_build_ls_long=('--format:u' '--branch:u' '--event:u' '--status:u' '--limit:u' '--page:u')

_drone_build_last_short=()
_drone_build_last_long=('--format:u' '--branch:u')

_drone_build_info_short=()
_drone_build_info_long=('--format:u')

_drone_build_create_short=('-p:u')
_drone_build_create_long=('--commit:u' '--branch:u' '--param:u' '--format:u')

_drone_build_stop_short=()
_drone_build_stop_long=()

_drone_build_restart_short=('-p:u')
_drone_build_restart_long=('--param:u' '--format:u')

_drone_build_approve_short=()
_drone_build_approve_long=()

_drone_build_decline_short=()
_drone_build_decline_long=()

_drone_build_promote_short=('-p:u')
_drone_build_promote_long=('--param:u' '--format:u')

_drone_build_rollback_short=('-p:u')
_drone_build_rollback_long=('--param:u')

_drone_build_queue_short=()
_drone_build_queue_long=('--format:u' '--repo:u' '--branch:u' '--event:u' '--status:u')

# Commands: cron
_drone_cron_short=('-h')
_drone_cron_long=('ls' 'info' 'add' 'rm' 'disable' 'enable' 'exec' '--help')

_drone_cron_ls_short=()
_drone_cron_ls_long=()

_drone_cron_info_short=()
_drone_cron_info_long=('--format:u')

_drone_cron_add_short=()
_drone_cron_add_long=('--branch:u')

_drone_cron_rm_short=()
_drone_cron_rm_long=()

_drone_cron_disable_short=()
_drone_cron_disable_long=()

_drone_cron_enable_short=()
_drone_cron_enable_long=()

_drone_cron_exec_short=()
_drone_cron_exec_long=()

# Commands: log
_drone_log_short=('-h')
_drone_log_long=('purge' 'view' '--help')

_drone_log_purge_short=()
_drone_log_purge_long=()

_drone_log_view_short=()
_drone_log_view_long=()

# Commands: encrypt
_drone_encrypt_short=()
_drone_encrypt_long=('--allow-pull-requests' '--allow-push-on-pull-request')

# Commands: exec
_drone_exec_short=()
_drone_exec_long=('--pipeline:u' '--include:u' '--exclude:u' '--resume-at:u' '--clone' '--trusted' '--timeout:u' '--volume:u' '--network:u' '--registry:u' '--secret-file:f' '--env-file:f' '--privileged:u' '--netrc-username:u' '--netrc-password:u' '--netrc-machine:u' '--branch:u' '--event:u' '--instance:u' '--ref:u' '--sha:u' '--repo:u' '--deploy-to:u')

# Commands: info
_drone_info_short=()
_drone_info_long=('--format:u')

# Commands: repo
_drone_repo_short=('-h')
_drone_repo_long=('ls' 'info' 'enable' 'update' 'disable' 'repair' 'chown' 'sync' '--help')

_drone_repo_ls_short=()
_drone_repo_ls_long=('--format:u' '--org:u' '--active')

_drone_repo_info_short=()
_drone_repo_info_long=('--format:u')

_drone_repo_enable_short=()
_drone_repo_enable_long=()

_drone_repo_update_short=()
_drone_repo_update_long=('--trusted' '--protected' '--throttle' '--timeout' '--visibility' '--ignore-forks' '--ignore-pull-requests' '--auto-cancel-pull-requests' '--auto-cancel-pushes' '--config:u' '--build-counter:u' '--unsafe')

_drone_repo_disable_short=()
_drone_repo_disable_long=()

_drone_repo_repair_short=()
_drone_repo_repair_long=()

_drone_repo_chown_short=()
_drone_repo_chown_long=()

_drone_repo_sync_short=()
_drone_repo_sync_long=('--format:u')

# Commands: user
_drone_user_short=('-h')
_drone_user_long=('ls' 'info' 'add' 'update' 'rm' 'block' '--help')

_drone_user_ls_short=()
_drone_user_ls_long=('--format:u')

_drone_user_info_short=()
_drone_user_info_long=('--format:u')

_drone_user_add_short=()
_drone_user_add_long=('--admin' '--machine' '--token:u')

_drone_user_update_short=()
_drone_user_update_long=('--admin' '--active')

_drone_user_rm_short=()
_drone_user_rm_long=()

_drone_user_block_short=()
_drone_user_block_long=()

_drone() {
    declare command="${COMP_WORDS[1]}"
    declare sub_command="${COMP_WORDS[2]}"
    declare num_of_args="${#COMP_WORDS[@]}"
    declare current_arg="${COMP_WORDS[COMP_CWORD]}"
    declare previous_arg="${COMP_WORDS[COMP_CWORD-1]}"

    # If we're still typing the first command or we haven't typed any, return valid commands for that.
    if [[ "${num_of_args}" == "2" && "${current_arg:+x}" != "x" ]]; then
        COMPREPLY=("${valid_commands[@]}")
        return
    elif [[ "${num_of_args}" == "2" ]]; then
        mapfile -t COMPREPLY < <(printf "%s\n" "${valid_commands[@]}" "${valid_args[@]}" | grep -- "^${current_arg}" 2> /dev/null)
        return
    fi

    # Otherwise start processing subcommands.
    for i in "${contains_sub_commands[@]}"; do
        if [[ "${command}" == "${i}" ]]; then
            has_sub_command="1"
            break
        fi
    done

    unknown_argument_options=()
    file_argument_options=()
    no_argument_options=()

    if [[ "${sub_command:+x}" == "x" && "${num_of_args}" -gt "3" ]] && (( "${has_sub_command}" )); then
        short_opts="_drone_${command}_${sub_command}_short[@]"
        long_opts="_drone_${command}_${sub_command}_long[@]"
        declare command="${sub_command}"
    else
        short_opts="_drone_${command}_short[@]"
        long_opts="_drone_${command}_long[@]"
    fi

    short_opts=("${!short_opts}")
    long_opts=("${!long_opts}")
    
    # If we couldn't find any arguments for the specified command, abort.
    if [[ "${#short_opts[@]}" == "0" && "${#long_opts[@]}" == "0" ]]; then
        COMPREPLY=()
        return
    fi
    
    # Process argument definitions.
    for i in "${short_opts[@]}" "${long_opts[@]}"; do
        option_extension="$(echo "${i}" | grep -o ':[uf]$' | sed 's|^:||')"
        option="$(echo "${i}" | sed 's|:[uf]$||')"

        if [[ "${option_extension}" == "u" ]]; then
            unknown_argument_options+=("${option}")
        elif [[ "${option_extension}" == "f" ]]; then
            file_argument_options+=("${option}")
        else
            no_argument_options+=("${option}")
        fi
    done

    # Actually generate the completions.
    for i in "${unknown_argument_options[@]}"; do
        if [[ "${previous_arg}" == "${i}" ]]; then
            COMPREPLY=()
            return
        fi
    done

    for i in "${file_argument_options[@]}"; do
        if [[ "${previous_arg}" == "${i}" ]]; then
            mapfile -t COMPREPLY < <(find ./ -maxdepth 1 -type f -not -path './' | sed 's|^\./||' | grep "^${current_arg}" 2> /dev/null)
            return
        fi
    done

    mapfile -t COMPREPLY < <(printf '%s\n' "${unknown_argument_options[@]}" "${file_argument_options[@]}" "${no_argument_options[@]}" | grep "^${current_arg}" 2> /dev/null)
}

complete -F _drone drone
