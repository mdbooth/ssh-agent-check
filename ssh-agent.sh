#!/bin/sh

env_file=~/.ssh-agent.env

if [ -f "$env_file" ]; then
    . "$env_file"
    if ! ssh-agent-check; then
        kill $SSH_AGENT_PID 2>/dev/null
        rm "$env_file"
    fi
fi

if [ ! -f "$env_file" ]; then
    ssh-agent > "$env_file"
    . "$env_file"
fi

unset env_file
