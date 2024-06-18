#!/bin/bash

if [ -n "$1" ]; then
    case "$1" in
    "next-js")\
        npx --yes create-next-app@latest . --js --use-npm --tailwind --eslint --app --no-src-dir --no-import-alias 
        ;;
    "next-ts")\
        npx --yes create-next-app@latest . --ts --use-npm --tailwind --eslint --app --no-src-dir --no-import-alias
        ;;
    "nest")\
        npm install -g @nestjs/cli
        nest new . -p npm
        ;;
    *)
        echo "Unknown template: $1"
        ;;
    esac
fi

# Start code-server
exec code-server --bind-addr 0.0.0.0:8080 . --auth none
