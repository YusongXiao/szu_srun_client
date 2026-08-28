#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=srun.sh
. "$SCRIPT_DIR/srun.sh"

TEST_ENV_FILE="$(mktemp)"
trap 'rm -f "$TEST_ENV_FILE"' EXIT
ENV_FILE="$TEST_ENV_FILE"

assert_credentials() {
	local expected_username="$1"
	local expected_password="$2"

	if [[ "$USERNAME" != "$expected_username" || "$PASSWORD" != "$expected_password" ]]; then
		printf 'expected username %q, got %q\n' "$expected_username" "$USERNAME" >&2
		exit 1
	fi
}

printf '%s\n' 'username=123456' 'password=old-password' > "$ENV_FILE"
SRUN_USERNAME=456789 SRUN_PASSWORD=new-password resolve_credentials
assert_credentials 456789 new-password

printf '%s\n' 'SRUN_USERNAME=123456' 'SRUN_PASSWORD=old-password' > "$ENV_FILE"
SRUN_USERNAME=456789 SRUN_PASSWORD=new-password resolve_credentials
assert_credentials 456789 new-password

printf '%s\n' 'username=123456' 'password=old-password' > "$ENV_FILE"
resolve_credentials
assert_credentials 123456 old-password

printf 'credential resolution tests passed\n'
