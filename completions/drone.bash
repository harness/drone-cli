_gen_compreply() {
    mapfile -t COMPREPLY < <(compgen -W '${options[@]}' -- "${1}")
}

_lvl_check() {
    if [[ "${#words[@]}" -lt "${1}" ]]; then
        _gen_compreply "${cur}"
        return 1
    fi
}

_drone() {
    local options=('build' 'cron' 'log' 'encrypt' 'exec' 'info' 'repo' 'user' 'secret' 'server' 'queue' 'orgsecret' 'autoscale' 'convert' 'lint' 'sign' 'jsonnet' 'starlark' 'plugins' 'template' 'help' 'h')
    local help_args=('--help' '-h')
    local drone_events=('cron' 'custom' 'push' 'pull_request' 'tag' 'promote' 'rollback')

    local cur prev words cword
    _init_completion || return
    
    cmd="${words[1]}"
    subcmd="${words[2]}"
    subsubcmd="${words[3]}"
    
    _lvl_check 3 || return

    if ! compgen -W '${options[@]}' -- "${cmd}" 1> /dev/null; then
        COMPREPLY=()
        return
    fi

    case "${cmd}" in
        build)
            options=('ls' 'last' 'info' 'create' 'stop' 'restart' 'approve' 'decline' 'promote' 'rollback' 'queue' "${help_args[@]}")

            case "${subcmd}" in
                ls)
                    _lvl_check 4 || return
                    local options=('--format' '--branch' '--event' '--status' '--limit' '--page' "${help_args[@]}")

                    case "${prev}" in
                        --format|--branch|--event|--status|--limit|--page)
                            COMPREPLY=()
                            return
                            ;;
                    esac
                    
                    _gen_compreply "${cur}"
                    return
                    ;;

                last)
                    _lvl_check 4 || return
                    local options=('--format' '--branch' "${help_args[@]}")

                    case "${prev}" in
                        --format|--branch)
                            COMPREPLY=()
                            return
                            ;;
                    esac

                    _gen_compreply "${cur}"
                    return
                    ;;

                info)
                    _lvl_check 4 || return
                    local options=('--format' "${help_args[@]}")

                    case "${prev}" in
                        --format)
                            COMPREPLY=()
                            return
                            ;;
                    esac

                    _gen_compreply "${cur}"
                    return
                    ;;

                create)
                    _lvl_check 4 || return
                    local options=('--commit' '--branch' '--param' '--format' '-p' "${help_args[@]}")

                    case "${prev}" in
                        --commit|--branch|--param|--format|-p)
                            COMPREPLY=()
                            return
                            ;;
                    esac

                    _gen_compreply "${cur}"
                    return
                    ;;
                stop)
                    _lvl_check 4 || return
                    local options=("${help_args[@]}")

                    _gen_compreply "${cur}"
                    return
                    ;;

                restart)
                    _lvl_check 4 || return
                    local options=('--param' '--format' '-p' "${help_args[@]}")
                    
                    case "${prev}" in
                        --param|--format|-p)
                            COMPREPLY=()
                            return
                            ;;
                    esac

                    _gen_compreply "${cur}"
                    return
                    ;;

                approve)
                    _lvl_check 4 || return
                    local options=("${help_args[@]}")

                    _gen_compreply "${cur}"
                    return
                    ;;

                decline)
                    _lvl_check 4 || return
                    local options=("${help_args[@]}")

                    _gen_compreply "${cur}"
                    return
                    ;;

                promote)
                    _lvl_check 4 || return
                    local options=('--param' '--format' '-p' "${help_args[@]}")
                    
                    case "${prev}" in
                        --param|--format|-p)
                            COMPREPLY=()
                            return
                            ;;
                    esac

                    _gen_compreply "${cur}"
                    return
                    ;;

                rollback)
                    _lvl_check 4 || return
                    local options=('--param' '-p' "${help_args[@]}")
                    
                    case "${prev}" in
                        --param|-p)
                            COMPREPLY=()
                            return
                            ;;
                    esac

                    _gen_compreply "${cur}"
                    return
                    ;;

                queue)
                    _lvl_check 4 || return
                    local options=('--repo' '--branch' '--event' '--status' "${help_args[@]}")

                    case "${prev}" in
                        --repo|--branch|--event|--status)
                            COMPREPLY=()
                            return
                            ;;
                    esac

                    _gen_compreply
                    return
                    ;;
                *)
                    _gen_compreply "${subcmd}"
                    return
                    ;;
            esac
            ;;

        cron)
            options=('ls' 'info' 'add' 'rm' 'disable' 'enable' 'exec' "${help_args[@]}")

            case "${subcmd}" in
                ls)
                    _lvl_check 4 || return
                    local options=("${help_args[@]}")

                    _gen_compreply "${cur}"
                    return
                    ;;

                info)
                    _lvl_check 4 || return
                    local options=('--format' "${help_args[@]}")

                    case "${prev}" in
                        --format)
                            COMPREPLY=()
                            return
                            ;;
                    esac

                    _gen_compreply "${cur}"
                    return
                    ;;

                add)
                    _lvl_check 4 || return
                    local options=('--branch' "${help_args[@]}")

                    case "${prev}" in
                        --branch)
                            COMPREPLY=()
                            return
                            ;;
                    esac

                    _gen_compreply "${cur}"
                    return
                    ;;

                rm)
                    _lvl_check 4 || return
                    local options=("${help_args[@]}")

                    _gen_compreply "${cur}"
                    return
                    ;;

                disable)
                    _lvl_check 4 || return
                    local options=("${help_args[@]}")

                    _gen_compreply "${cur}"
                    return
                    ;;

                enable)
                    _lvl_check 4 || return
                    local options=("${help_args[@]}")

                    _gen_compreply "${cur}"
                    return
                    ;;

                exec)
                    _lvl_check 4 || return
                    local options=("${help_args[@]}")

                    _gen_compreply "${cur}"
                    return
                    ;;

                *)
                    _gen_compreply "${subcmd}"
                    return
                    ;;
            esac
            ;;

        log)
            local options=('purge' 'view' "${help_args[@]}")

            case "${subcmd}" in
                purge)
                    _lvl_check 4 || return
                    local options=("${help_args[@]}")

                    _gen_compreply "${cur}"
                    return
                    ;;

                view)
                    _lvl_check 4 || return
                    local options=("${help_args[@]}")

                    _gen_compreply "${cur}"
                    return
                    ;;

                *)
                    _gen_compreply "${subcmd}"
                    return
                    ;;
            esac
            ;;

        encrypt)
            local options=('--allow-pull-request' '--allow-push-on-pull-request' "${help_args[@]}")

            _gen_compreply "${cur}"
            return
            ;;

        exec)
            local options=('--pipeline' '--include' '--exclude' '--resume-at' '--clone' '--trusted' '--timeout' '--volume' '--network' '--registry' '--secret-file' '--env-file' '--privileged' '--netrc-username' '--netrc-password' '--netrc-machine' '--branch' '--event' '--instance' '--ref' '--sha' '--repo' '--deploy-to' "${help_args[@]}")

            case "${prev}" in
                --pipeline|--include|--exclude|--resume-at|--timeout|--volume|--network|--registry|--privileged|--netrc-username|--netrc-password|--netrc-machine|--branch|--instance|--ref|--sha|--repo|--deploy-to)
                    COMPREPLY=()
                    return
                    ;;

                --secret-file|--env-file)
                    _filedir
                    return
                    ;;

                --event)
                    local options=("${drone_events[@]}")
                    _gen_compreply "${cur}"
                    return
                    ;;
            esac

            _gen_compreply "${cur}"
            return
            ;;

        info)
            local options=('--format')

            case "${prev}" in
                --format)
                    COMPREPLY=()
                    return
                    ;;
            esac

            _gen_compreply "${cur}"
            return
            ;;
    esac
}

complete -F _drone drone
