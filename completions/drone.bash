# Base options.
_drone_short=('-t' '--token' '-s' '--server'
    '--autoscaler' '--help' '-h' '--version'
    '-v')

_drone_long=('build:s' 'cron:s' 'log:s' 'encrypt:s'
             'exec:s' 'info:s' 'repo:s' 'user:s'
             'secret:s' 'server:s' 'queue:s' 'orgsecret:s'
             'autoscale:s' 'convert:s' 'lint:s' 'sign:s'
             'jsonnet:s' 'starlark:s' 'plugins:s' 'template:s'
             'help' 'h')


# Commands: build
_drone_build_short=('-h')
_drone_build_long=('ls:s' 'last:s' 'info:s' 'create:s' 'stop:s' 'restart:s' 'approve:s' 'decline:s' 'promote:s' 'rollback:s' 'queue:s' '--help:s')

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
_drone_cron_long=('ls:s' 'info:s' 'add:s' 'rm:s' 'disable:s' 'enable:s' 'exec:s' '--help')

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
_drone_log_long=('purge:s' 'view:s' '--help')

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
_drone_repo_long=('ls:s' 'info:s' 'enable:s' 'update:s' 'disable:s' 'repair:s' 'chown:s' 'sync:s' '--help')

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
_drone_user_long=('ls:s' 'info:s' 'add:s' 'update:s' 'rm:s' 'block:s' '--help')

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

_drone_secret_short=('-h')
_drone_secret_long=('add:s' 'rm:s' 'update:s' 'info:s' 'ls:s' '--help')

_drone_secret_add_short=()
_drone_secret_add_long=('--name:u' '--data:u' '--allow-pull-request' '--allow-push-on-pull-request')

_drone_secret_rm_short=()
_drone_secret_rm_long=('--repository:u' '--name:u')

_drone_secret_update_short=()
_drone_secret_update_long=('--name:u' '--data:u' '--allow-pull-request' '--alow-push-on-pull-request')

_drone_secret_info_short=()
_drone_secret_info_long=('--repository:u' '--name:u' '--format:u')

_drone_secret_ls_short=()
_drone_secret_ls_long=('--repository:u' '--format:u')

_drone_server_short=('-h')
_drone_server_long=('ls' 'info' 'create' 'destroy' 'env' '--help')

_drone_server_ls_short=('-s:u' '-a' '-l' '-H')
_drone_server_ls_long=('--state:u' '--all' '--long' '--headers' '--format:u')

_drone_server_info_short=()
_drone_server_info_long=('--format:u')

_drone_server_create_short=()
_drone_server_create_long=('--format:u')

_drone_server_destroy_short=()
_drone_server_destroy_long=('--format:u' '--force')

_drone_server_env_short=()
_drone_server_env_long=('--shell:u' '--no-proxy' '--clear')

_drone_queue_short=('-h')
_drone_queue_long=('ls:s' 'pause:s' 'resume:s' '--help')

_drone_queue_ls_short=()
_drone_queue_ls_long=('--format:u')

_drone_queue_pause_short=()
_drone_queue_pause_long=()

_drone_queue_resume_short=()
_drone_queue_resume_long=()

_drone_orgsecret_short=('-h')
_drone_orgsecret_long=('add:s' 'rm:s' 'update:s' 'info:s' 'ls:s')

_drone_orgsecret_add_short=()
_drone_orgsecret_add_long=('--allow-pull-request' '--allow-push-on-pull-request')

_drone_orgsecret_rm_short=()
_drone_orgsecret_rm_long=()

_drone_orgsecret_update_short=()
_drone_orgsecret_update_long=('--allow-pull-request' '--allow-push-on-pull-request')

_drone_orgsecret_info_short=()
_drone_orgsecret_info_long=('--format:u')

_drone_orgsecret_ls_short=()
_drone_orgsecret_ls_long=('--filter:u' '--format:u')

_drone_autoscale_short=('-h')
_drone_autoscale_long=('pause:s' 'resume:s' 'version:s' '--help')

_drone_autoscale_pause_short=()
_drone_autoscale_pause_long=()

_drone_autoscale_resume_short=()
_drone_autoscale_resume_long=()

_drone_autoscale_version_short=()
_drone_autoscale_version_long=('--format:u')

_drone_convert_short=()
_drone_convert_long=('--save')

_drone_lint_short=()
_drone_lint_long=('--trusted' ':f')

_drone_sign_short=()
_drone_sign_long=('--save' ':f')

_drone_jsonnet_short=('-V:u')
_drone_jsonnet_long=('--source:f' '--target:f' '--stream' '--format' '--stdout' '--string' '--extVar:u')

_drone_starlark_short=()
_drone_starlark_long=('--source:f' '--target:f' '--format' '--stdout' '--max-execution-steps:u' '--repo.name:u' '--repo.namespace:u' '--repo.slug:u' '--build.event:u' '--build.branch:u' '--build.source:u' '--build.source_repo:u' '--build.target:u' '--build.ref:u' '--build.commit:u' '--build.message:u' '--build.title:u' '--build.link:u' '--build.environment:u' '--build.debug:u')

_drone_plugins_short=('-h')
_drone_plugins_long=('admit:s' 'config:s' 'convert:s' 'env:s' 'registry:s' 'secret:s' '--help')

_drone_plugins_admit_short=()
_drone_plugins_admit_long=('--user:u' '--endpoint:u' '--secret:u' '--ssl-skip-verify:u')

_drone_plugins_config_short=('-h')
_drone_plugins_config_long=('get:s' '--help')

_drone_plugins_config_get_short=()
_drone_plugins_config_get_long=('--ref:u' '--source:u' '--target:u' '--before:u' '--after:u' '--path:u' '--endpoint:u' '--secret:u' '--ssl-skip-verify:u')

_drone_plugins_convert_short=()
_drone_plugins_convert_long=('--path:u' '--ref:u' '--source:u' '--target:u' '--before:u' '--after:u' '--repository:u' '--endpoint:u' '--secret:u' '--ssl-skip-verify:u')

_drone_plugins_env_short=()
_drone_plugins_env_long=('--ref:u' '--source:u' '--target:u' '--before:u' '--after:u' '--repository:u' '--endpoint:u' '--secret:u' '--skip-skip-verify:u')

_drone_plugins_registry_short=('-h')
_drone_plugins_registry_long=('list:s' '--help')

_drone_plugins_registry_list_short=()
_drone_plugins_registry_list_long=('--ref:u' '--source:u' '--target:u' '--before:u' '--after:u' '--event:u' '--repo:u' '--endpoint:u' '--secret:u' '--ssl-skip-verify:u' '--format:u')

_drone_plugins_secret_short=('-h')
_drone_plugins_secret_long=('get:s' '--help')

_drone_plugins_secret_get_short=()
_drone_plugins_secret_get_long=('--ref:u' '--source:u' '--target:u' '--before:u' '--after:u' '--event:u' '--repo:u' '--endpoint:u' '--secret:u' '--ssl-skip-verify:u')

_drone_template_short=('-h')
_drone_template_long=('add:s' 'info:s' 'ls:s' 'update:s' 'rm:s' '--help')

_drone_template_add_short=()
_drone_template_add_long=('--name:u' '--namespace:u' '--data:u')

_drone_template_info_short=()
_drone_template_info_long=('--namespace:u' '--name:u' '--format:u')

_drone_template_ls_short=()
_drone_template_ls_long=('--namespace:u' '--format:u')

_drone_template_update_short=()
_drone_template_update_long=('--name:u' '--namespace:u' '--data:u')

_drone_template_rm_short=()
_drone_template_rm_long=('--name:u' '--namespace:u')

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
        previous_arg="${current_arg}"
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
    
    for i in "${accepts_files[@]}"; do
        if [[ "${previous_arg}" == "" ]] && (echo "${previous_arg}" | grep -q '^[^-]'); then
            mapfile -t COMPREPLY < <(find ./ -maxdepth 1 -type f -not -path './' | sed 's|^\./||' | grep "^${current_arg}" 2> /dev/null)
            return
        fi
    done

    mapfile -t COMPREPLY < <(printf '%s\n' "${unknown_argument_options[@]}" "${file_argument_options[@]}" "${no_argument_options[@]}" | grep "^${current_arg}" 2> /dev/null)
}

complete -F _drone drone
