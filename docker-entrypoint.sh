#!/bin/sh

set -e

if [ ! -f config.yml ] && [ "no$STORAGE_TYPE" != "no" ]; then
    # try to do ENV overriding
    echo "Initial the configuration ..."
    cat > config.yml <<EOT
storage:
    type: $STORAGE_TYPE
    address: $STORAGE_ADDRESS
    username: $STORAGE_USERNAME
    password: $STORAGE_PASSWORD
    namespace: $STORAGE_NAMESPACE
EOT
    chown goschedule:goschedule config.yml
fi

if [ "no$@" = "no" ]; then
    set $@ "goschedule-console"
fi

exec gosu goschedule $@