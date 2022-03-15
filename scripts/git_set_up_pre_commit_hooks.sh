#!/usr/bin/env bash

set -e

SCRIPT_DIR=$(dirname $0)
cd $SCRIPT_DIR/..

if grep -sq "user-service linter" .git/hooks/pre-commit ; then
    echo "Pre-commit hook already set up."
    exit 0
fi

cat <<EOF > .git/hooks/pre-commit
# Set up user-service linter
cd user-service && make lint
EOF

chmod +x .git/hooks/pre-commit

echo "Pre-commit hook successfully set up."
