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
    local help_args=('--help')
    local drone_events=('cron' 'custom' 'push' 'pull_request' 'tag' 'promote' 'rollback')

    local cur prev words cword
    _init_completion || return
    
    cmd="${words[1]}"
    subcmd="${words[2]}"
    subsubcmd="${words[3]}"
    
    case "${cmd}" in
        build)
            _lvl_check 3 || return
            options=('ls' 'last' 'info' 'create' 'stop' 'restart' 'approve' 'decline' 'promote' 'rollback' 'queue')

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
                    local options=('--commit' '--branch' '--param' '--format' "${help_args[@]}")

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
                    local options=('--param' '--format' "${help_args[@]}")
                    
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
                    local options=('--param' '--format' "${help_args[@]}")
                    
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
                    local options=('--param' "${help_args[@]}")
                    
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
            _lvl_check 3 || return
            options=('ls' 'info' 'add' 'rm' 'disable' 'enable' 'exec')

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
            _lvl_check 3 || return
            local options=('purge' 'view')

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
            _lvl_check 3 || return
            local options=('--allow-pull-request' '--allow-push-on-pull-request' "${help_args[@]}")

            _gen_compreply "${cur}"
            return
            ;;

        exec)
            _lvl_check 3 || return
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
            _lvl_check 3 || return
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

        repo)
            _lvl_check 3 || return
            local options=('ls' 'info' 'enable' 'update' 'disable' 'repair' 'chown' 'sync')

            case "${subcmd}" in
                ls)
                    _lvl_check 4 || return
                    local options=('--format' '--org' '--active')

                    case "${prev}" in
                        --format|--org)
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

                enable)
                    _lvl_check 4 || return
                    local options=("${help_args[@]}")

                    _gen_compreply "${cur}"
                    return
                    ;;

                update)
                    _lvl_check 4 || return
                    local options=('--trusted' '--protected' '--throttle' '--timeout' '--visibility' '--ignore-forks' '--ignore-pull-requests' '--auto-cancel-pull-requests' '--auto-cancel-pushes' '--auto-cancel-running' '--config' '--build-counter' '--unsafe')

                    case "${prev}" in
                        --throttle|--timeout|--visibility|--config|--build-counter)
                            COMPREPLY=()
                            return
                            ;;
                    esac

                    _gen_compreply "${cur}"
                    return
                    ;;

                disable)
                    _lvl_check 4 || return
                    local options=("${help_args[@]}")

                    _gen_compreply "${cur}"
                    return
                    ;;

                repair)
                    _lvl_check 4 || return
                    local options=("${help_args[@]}")

                    _gen_compreply "${cur}"
                    return
                    ;;

                chown)
                    _lvl_check 4 || return
                    local options=("${help_args[@]}")

                    _gen_compreply "${cur}"
                    return
                    ;;

                sync)
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

                *)
                    _gen_compreply "${cur}"
                    return
                    ;;
            esac
            ;;

        user)
            _lvl_check 3 || return
            local options=('ls' 'info' 'add' 'update' 'rm' 'block')

            case "${subcmd}" in
                ls)
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
                    local options=('--admin' '--machine' '--token')

                    case "${prev}" in
                        --token)
                            COMPREPLY=()
                            return
                            ;;
                    esac

                    _gen_compreply "${cur}"
                    return
                    ;;

                update)
                    _lvl_check 4 || return
                    local options=('--active' '--admin')

                    _gen_compreply "${cur}"
                    return
                    ;;

                rm)
                    _lvl_check 4 || return
                    local options=("${help_args[@]}")

                    _gen_compreply "${cur}"
                    return
                    ;;

                block)
                    _lvl_check 4 || return
                    local options=("${help_args[@]}")

                    _gen_compreply "${cur}"
                    return
                    ;;

                *)
                    _gen_compreply "${cur}"
                    return
                    ;;
            esac
            ;;

        secret|orgsecret)
            _lvl_check 3 || return
            local options=('add' 'rm' 'update' 'info' 'ls')

            case "${subcmd}" in
                add)
                    _lvl_check 4 || return
                    local options=('--repository' '--name' '--data' '--allow-pull-request' '--allow-push-on-pull-request' "${help_args[@]}")

                    case "${prev}" in
                        --repository|--name|--data)
                            COMPREPLY=()
                            return
                            ;;
                    esac

                    _gen_compreply "${cur}"
                    return
                    ;;

                rm)
                    _lvl_check 4 || return
                    local options=('--repository' '--name' "${help_args[@]}")

                    case "${prev}" in
                        --repository|--name)
                            COMPREPLY=()
                            return
                            ;;
                    esac

                    _gen_compreply "${cur}"
                    return
                    ;;

                update)
                    _lvl_check 4 || return
                    local options=('--repository' '--name' '--data' '--allow-pull-request' '--allow-push-on-pull-request' "${help_args[@]}")

                    case "${prev}" in
                        --repository|--name|--data)
                            COMPREPLY=()
                            return
                            ;;
                    esac

                    _gen_compreply "${cur}"
                    return
                    ;;

                info)
                    _lvl_check 4 || return
                    local options=('--repository' '--name' '--format' "${help_args[@]}")

                    case "${prev}" in
                        --repository|--name|--format)
                            COMREPLY=()
                            return
                            ;;
                    esac

                    _gen_compreply "${cur}"
                    return
                    ;;

                ls)
                    _lvl_check 4 || return
                    local options=('--repository' '--format')

                    case "${prev}" in
                        --repository|--format)
                            COMPREPLY=()
                            return
                            ;;
                    esac

                    _gen_compreply "${cur}"
                    return
                    ;;

                *)
                    _gen_compreply "${cur}"
                    return
                    ;;
            esac
            ;;

        server)
            _lvl_check 3 || return
            local options=('ls' 'info' 'create' 'destroy' 'env')

            case "${subcmd}" in
                ls)
                    _lvl_check 4 || return
                    local options=('--state' '--all' '--long' '--headers' '--format' "${help_args[@]}")

                    case "${prev}" in
                        --state|--format|-s)
                            COMPREPLY=()
                            return
                            ;;
                    esac

                    _gen_compreply "${cur}"
                    return
                    ;;

                info|create)
                    _lvl_check 4 || return
                    local options=('--format' "${help_args[@]}")

                    case "${prev}" in
                        --format)
                            COMPREPLY=()
                            return
                            ;;
                    esac

                    _gen_compreply "${cur}"
                    ;;

                destroy)
                    _lvl_check 4 || return
                    local options=('--format' '--force' "${help_args[@]}")

                    case "${prev}" in
                        --format)
                            COMPREPLY=()
                            return
                            ;;
                    esac

                    _gen_compreply "${cur}"
                    return
                    ;;

                env)
                    _lvl_check 4 || return
                    local options=('--shell' '--no-proxy' '--clear' "${help_args[@]}")

                    case "${prev}" in
                        --shell)
                            local options=('bash' 'fish' 'powershell')
                            _gen_compreply "${cur}"
                            return
                            ;;
                    esac

                    _gen_compreply "${cur}"
                    ;;

                *)
                    _gen_compreply "${cur}"
                    return
                    ;;
            esac
            ;;

        queue)
            _lvl_check 3 || return
            local options=('ls' 'pause' 'resume')

            case "${subcmd}" in
                ls)
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

                pause|resume)
                    _lvl_check 4 || return
                    local options=("${help_args[@]}")

                    _gen_compreply "${cur}"
                    return
                    ;;
                *)
                    _gen_compreply "${cur}"
                    return
                    ;;
            esac
            ;;

        autoscale)
            _lvl_check 3 || return
            local options=('pause' 'resume' 'version')

            case "${subcmd}" in
                pause|resume)
                    _lvl_check 4 || return
                    local options=("${help_args[@]}")

                    _gen_compreply "${cur}"
                    return
                    ;;

                version)
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

                *)
                    _gen_compreply "${cur}"
                    return
                    ;;
            esac
            ;;

        convert|save)
            _lvl_check 3 || return
            local options=('--save' "${help_args[@]}")

            _gen_compreply "${cur}"
            return
            ;;

        lint)
            _lvl_check 3 || return
            local options=('--trusted' "${help_args[@]}")

            _gen_compreply "${cur}"
            return
            ;;

        jsonnet)
            _lvl_check 3 || return
            local options=('--source' '--target' '--stream' '--format' '--stdout' '--string' '--extVar' "${help_args[@]}")

            case "${prev}" in
                --source|--target)
                    _filedir
                    return
                    ;;

                --extVar|-V)
                    COMPREPLY=()
                    return
                    ;;
            esac

            _gen_compreply "${cur}"
            ;;

        starlark)
            _lvl_check 3 || return
            local options=('--source' '--target' '--format' '--stdout' '--max-execution-steps' '--repo.name' '--repo.namespace' '--repo.slug' '--build.event' '--build.branch' '--build.source' '--build.source_repo' '--build.target' '--build.ref' '--build.commit' '--build.message' '--build.title' '--build.link' '--build.environment' '--build.debug' "${help_args[@]}")

            case "${prev}" in
                --source|--target)
                    _filedir
                    return
                    ;;

                --max-execution-steps|--repo.name|--repo.namespace|--repo.slug|--build.event|--build.branch|--build.source|--build.source_repo|--build.target|--build.ref|--build.commit|--build.message|--build.title|--build.link|--build.environment)
                    COMPREPLY=()
                    return
                    ;;
            esac

            _gen_compreply "${cur}"
            return
            ;;

        plugins)
            _lvl_check 3 || return
            local options=('admit' 'config' 'convert' 'env' 'registry' 'secret')

            case "${subcmd}" in
                admit)
                    _lvl_check 4 || return
                    local options=('--user' '--endpoint' '--secret' '--ssl-skip-verify' "${help_args[@]}")

                    case "${prev}" in
                        --user|--endpoint|--secret|--ssl-skip-verify)
                            COMPREPLY=()
                            return
                            ;;
                    esac

                    _gen_compreply "${cur}"
                    return
                    ;;

                config)
                    _lvl_check 4 || return
                    local options=('get')

                    case "${subsubcmd}" in
                        get)
                            _lvl_check 5 || return
                            local options=('--ref' '--source' '--target' '--before' '--after' '--path' '--endpoint' '--secret' '--ssl-skip-verify')

                            case "${prev}" in
                                --ref|--source|--target|--before|--after|--path|--endpoint|--secret|--ssl-skip-verify)
                                    COMPREPLY=()
                                    return
                                    ;;
                            esac

                            _gen_compreply "${cur}"
                            return
                            ;;
                        *)
                            _gen_compreply "${cur}"
                            return
                            ;;
                    esac
                    ;;

                convert)
                    _lvl_check 4 || return
                    local options=('--path' '--ref' '--source' '--target' '--before' '--after' '--repository' '--endpoint' '--secret' '--ssl-skip-verify')

                    case "${prev}" in
                        --path|--ref|--source|--target|--before|--after|--repository|--endpoint|--secret|--ssl-skip-verify)
                            COMPREPLY=()
                            return
                            ;;
                    esac

                    _gen_compreply "${cur}"
                    return
                    ;;

                env)
                    _lvl_check 4 || return
                    local options=('--ref' '--source' '--target' '--before' '--after' '--repository' '--endpoint' '--secret' '--ssl-skip-verify')

                    case "${prev}" in
                        --ref|--source|--target|--before|--after|--repository|--endpoint|--secret|--ssl-skip-verify)
                            COMPREPLY=()
                            return
                            ;;
                    esac

                    _gen_compreply "${cur}"
                    ;;

                registry)
                    _lvl_check 4 || return
                    local options=('list')

                    case "${subsubcmd}" in
                        list)
                            _lvl_check 5 || return
                            local options=('--ref' '--source' '--target' '--before' '--after' '--event' '--repo' '--endpoint' '--secret' '--ssl-skip-verify' '--format')

                            case "${prev}" in
                                --ref|--source|--target|--before|--after|--event|--repo|--endpoint|--secret|--ssl-skip-verify|--format)
                                    COMPREPLY=()
                                    return
                                    ;;
                            esac

                            _gen_compreply "${cur}"
                            return
                            ;;
                        *)
                            _gen_compreply "${cur}"
                            return
                            ;;
                    esac
                    ;;

                secret)
                    _lvl_check 4 || return
                    local options=('get')

                    case "${subsubcmd}" in
                        get)
                            _lvl_check 5 || return
                            local options=('--ref' '--source' '--target' '--before' '--after' '--event' '--repo' '--endpoint' '--secret' '--ssl-skip-verify')

                            case "${prev}" in
                                --ref|--source|--target|--before|--after|--event|--repo|--endpoint|--secret|--ssl-skip-verify)
                                    COMPREPLY=()
                                    return
                                    ;;
                            esac

                            _gen_compreply "${cur}"
                            return
                            ;;

                        *)
                            _gen_compreply "${cur}"
                            return
                            ;;
                    esac
                    ;;
                *)
                    _gen_compreply "${cur}"
                    return
                    ;;
            esac
            ;;

        template)
            _lvl_check 3 || return
            local options=('add' 'info' 'ls' 'add' 'update' 'rm')

            case "${subcmd}" in
                add|update)
                    _lvl_check 4 || return
                    local options=('--name' '--namespace' '--data' "${help_args[@]}")

                    case "${prev}" in
                        --name|--namespace|--data)
                            COMPREPLY=()
                            return
                            ;;
                    esac

                    _gen_compreply "${cur}"
                    return
                    ;;

                info)
                    _lvl_check 4 || return
                    local options=('--namespace' '--name' '--format')

                    case "${prev}" in
                        --namespace|--name|--format)
                            COMREPLY=()
                            return
                            ;;
                    esac

                    _gen_compreply "${cur}"
                    return
                    ;;

                ls)
                    _lvl_check 4 || return
                    local options=('--namespace' '--format')

                    case "${prev}" in
                        --namespace|--format)
                            COMPREPLY=()
                            return
                            ;;
                    esac

                    _gen_compreply "${cur}"
                    return
                    ;;

                rm)
                    _lvl_check 4 || return
                    local options=('--name' '--namespace')

                    case "${prev}" in
                        --name|--namespace)
                            COMPREPLY=()
                            return
                            ;;
                    esac

                    _gen_compreply "${cur}"
                    return
                    ;;
            esac
            ;;

        help|h)
            _gen_compreply "${cur}"
            ;;

        *)
            _gen_compreply "${cur}"
            return
            ;;
    esac
}

complete -F _drone drone
