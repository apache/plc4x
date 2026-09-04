#!/usr/bin/env bash

# ----------------------------------------------------------------------------
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#    https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
# ----------------------------------------------------------------------------

# Checks everything a release needs before any of it is done for real.
#
# A release runs for the better part of an hour, and almost everything that has gone wrong with it
# so far was a precondition that could have been checked in seconds: a Nexus password that had
# expired, a signing key that is not in the KEYS file, a gpg agent that cannot sign inside the
# release container. This script gathers those checks in one place, so a release fails in a minute
# rather than after the third full build of the day.
#
# It only reads: nothing here creates a commit, a tag, a staging repository or an SVN directory.
#
#   tools/release-preflight.sh              run every check
#   tools/release-preflight.sh --offline    skip the checks that need the network
#
# It is meant to be run at the start of the release scripts, and can be run on its own at any time.

DIRECTORY="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
if [[ ! -f "$DIRECTORY/tools/release-common.sh" ]]; then
    echo "❌ '$DIRECTORY/tools/release-common.sh' not found, aborting."
    exit 1
fi
# shellcheck source=release-common.sh
source "$DIRECTORY/tools/release-common.sh"

OFFLINE=false
if [[ "$1" == "--offline" ]]; then
    OFFLINE=true
fi

# Every check adds to one of these, and the summary at the end decides the exit code. Checks do
# not abort the script on purpose: a release manager should see everything that needs fixing in
# one go, instead of running into the problems one release attempt at a time.
FAILURES=()
WARNINGS=()

fail() {
    echo "❌ $1"
    FAILURES+=("$1")
}

warn() {
    echo "⚠️  $1"
    WARNINGS+=("$1")
}

ok() {
    echo "✅ $1"
}

########################################################################################################################
# 1. The tools the release scripts call
########################################################################################################################

echo
echo "── Tools ─────────────────────────────────────────────────────────────────────"

for tool in git docker gpg svn curl java shasum; do
    if command -v "$tool" > /dev/null 2>&1; then
        ok "'$tool' is installed"
    else
        fail "'$tool' is not installed, the release scripts need it."
    fi
done

# The whole build runs in a container, so a docker binary without a running daemon gets us nowhere.
if command -v docker > /dev/null 2>&1; then
    if docker info > /dev/null 2>&1; then
        ok "The docker daemon is running"
    else
        fail "The docker daemon is not running, start Docker Desktop before releasing."
    fi
fi

########################################################################################################################
# 2. Git
########################################################################################################################

echo
echo "── Git ───────────────────────────────────────────────────────────────────────"

GIT_USER_NAME=$(git -C "$DIRECTORY" config user.name || git -C "$DIRECTORY" config --global user.name)
GIT_USER_EMAIL=$(git -C "$DIRECTORY" config user.email || git -C "$DIRECTORY" config --global user.email)
if [[ -n "$GIT_USER_NAME" && -n "$GIT_USER_EMAIL" ]]; then
    ok "Git identity is set to '$GIT_USER_NAME <$GIT_USER_EMAIL>'"
else
    fail "Git 'user.name' and/or 'user.email' are not configured."
fi

if [[ $(git -C "$DIRECTORY" status --porcelain) ]]; then
    warn "The working tree has uncommitted changes, the release scripts refuse to start like that."
else
    ok "The working tree is clean"
fi

if [[ "$OFFLINE" == false ]]; then
    # This authenticates against the remote without changing anything. It does not prove the push
    # will be accepted: the ".asf.yaml" rulesets are evaluated on the real push, and a bypass only
    # applies to the identity that ends up pushing.
    PUSH_OUTPUT=$(git -C "$DIRECTORY" push --dry-run origin HEAD 2>&1)
    PUSH_STATUS=$?
    if [[ $PUSH_STATUS -eq 0 ]]; then
        ok "Can authenticate against 'origin' for pushing"
    elif echo "$PUSH_OUTPUT" | grep -qiE "permission denied|authentication|403|could not read from remote"; then
        fail "Cannot push to 'origin' - check the remote, your ssh key and your access rights."
    else
        # A rejected push is a different matter than a rejected login: the branch may simply have
        # moved on, which is not something this check has an opinion about.
        warn "'git push --dry-run' did not succeed, but not for lack of access: $(echo "$PUSH_OUTPUT" | tail -n 1)"
    fi
fi

########################################################################################################################
# 3. GPG
########################################################################################################################

echo
echo "── Signing ───────────────────────────────────────────────────────────────────"

# The release is signed with whatever key git is configured to use, or the default secret key.
SIGNING_KEY=$(git -C "$DIRECTORY" config user.signingkey)
if [[ -z "$SIGNING_KEY" ]]; then
    SIGNING_KEY=$(gpg --list-secret-keys --with-colons 2>/dev/null | awk -F: '/^sec:/ {print $5; exit}')
fi

if [[ -z "$SIGNING_KEY" ]]; then
    fail "No gpg secret key found, a release has to be signed."
else
    ok "Signing with key '$SIGNING_KEY'"

    # An Apache release has to be signed with a key that carries the release manager's
    # "{apache-id}@apache.org" address, so check that before the signing takes an hour to fail.
    SIGNING_UIDS=$(gpg --with-colons --list-keys "$SIGNING_KEY" 2>/dev/null \
        | awk -F: '$1 == "uid" && $2 !~ /^[re]$/ {print $10}')
    if echo "$SIGNING_UIDS" | grep -qiE '<[^>]+@apache\.org>'; then
        ok "The signing key carries an apache.org address"
    else
        fail "The signing key '$SIGNING_KEY' carries no (valid) apache.org address."
    fi

    # Signing hundreds of artifacts only to find out that the agent cannot sign is a long way to
    # go for that news, so sign one throwaway file here.
    SIGN_TEST_DIR=$(mktemp -d)
    echo "plc4x release preflight" > "$SIGN_TEST_DIR/test.txt"
    if gpg --batch --yes --local-user "$SIGNING_KEY" -ab "$SIGN_TEST_DIR/test.txt" > /dev/null 2>&1; then
        ok "The signing key can sign without interaction"
    else
        fail "Cannot sign with '$SIGNING_KEY' - is the gpg agent running and the passphrase cached?"
    fi
    rm -rf "$SIGN_TEST_DIR"

    if [[ "$OFFLINE" == false ]]; then
        # Being able to sign is not enough, the key also has to be published in the KEYS file, or
        # nobody can verify the release. The release scripts check this after signing everything.
        KEYS_DIR=$(mktemp -d)
        if curl -fsSL "$KEYS_URL" -o "$KEYS_DIR/KEYS" && [[ -s "$KEYS_DIR/KEYS" ]]; then
            SIGNING_FINGERPRINT=$(gpg --with-colons --fingerprint "$SIGNING_KEY" 2>/dev/null \
                | awk -F: '/^fpr:/ {print $10; exit}')
            if gpg --no-default-keyring --keyring "$KEYS_DIR/pubring.kbx" --import "$KEYS_DIR/KEYS" > /dev/null 2>&1 \
                && gpg --no-default-keyring --keyring "$KEYS_DIR/pubring.kbx" --with-colons --fingerprint 2>/dev/null \
                    | grep -q "$SIGNING_FINGERPRINT"; then
                ok "The signing key is published in the KEYS file"
            else
                fail "The signing key is not in $KEYS_URL - add it before releasing."
            fi
        else
            warn "Could not download the KEYS file from $KEYS_URL, skipping that check."
        fi
        rm -rf "$KEYS_DIR"
    fi
fi

########################################################################################################################
# 4. Nexus
########################################################################################################################

if [[ "$OFFLINE" == false ]]; then
    echo
    echo "── Nexus ─────────────────────────────────────────────────────────────────────"

    # "rc-list" is the read-only counterpart of the staging deployment: same server id, same
    # credentials, same "stage.pom", so it fails in exactly the same way a stale password would
    # fail an hour into the release. Letting Maven do the call also means the credentials are
    # read and decrypted the way Maven does it, rather than parsed out of settings.xml here.
    NEXUS_OUTPUT=$("$DIRECTORY/mvnw" -q -f "$DIRECTORY/tools/stage.pom" \
        nexus-staging:rc-list -DstagingProfileId="$STAGING_PROFILE_ID" 2>&1)
    NEXUS_STATUS=$?
    if [[ $NEXUS_STATUS -eq 0 ]]; then
        ok "Can authenticate against $NEXUS_URL with the 'apache.releases.https' credentials"
    elif echo "$NEXUS_OUTPUT" | grep -q "401"; then
        fail "$NEXUS_URL rejects the 'apache.releases.https' credentials from your settings.xml (401)."
    else
        fail "Could not talk to $NEXUS_URL, run 'mvnw -f tools/stage.pom nexus-staging:rc-list -DstagingProfileId=$STAGING_PROFILE_ID' to see why."
    fi
fi

########################################################################################################################
# 5. SVN
########################################################################################################################

if [[ "$OFFLINE" == false ]]; then
    echo
    echo "── SVN ───────────────────────────────────────────────────────────────────────"

    # The release candidate is imported into "dist/dev", and the finished release is moved to
    # "dist/release", so both have to be reachable. This only proves read access - a commit is the
    # only thing that proves write access, and that is not something to do in a check.
    if svn --non-interactive ls "$DIST_DEV" > /dev/null 2>&1; then
        ok "'$DIST_DEV' is reachable"
    else
        fail "Cannot read '$DIST_DEV' - are your svn credentials cached?"
    fi
fi

########################################################################################################################
# 6. Disk space
########################################################################################################################

echo
echo "── Resources ─────────────────────────────────────────────────────────────────"

# A full release build with all languages, the local repository and the staged artifacts needs
# a good deal of room, and running out halfway through is a particularly annoying way to fail.
# "df -k" is the portable spelling, "df -g" only exists on macOS.
AVAILABLE_GB=$(df -k "$DIRECTORY" 2>/dev/null | awk 'NR==2 {print int($4 / 1024 / 1024)}')
if [[ -n "$AVAILABLE_GB" ]]; then
    if [[ "$AVAILABLE_GB" -ge 20 ]]; then
        ok "${AVAILABLE_GB}GB free on the volume holding the checkout"
    else
        warn "Only ${AVAILABLE_GB}GB free on the volume holding the checkout, a release needs about 20GB."
    fi
fi

########################################################################################################################
# 7. Summary
########################################################################################################################

echo
echo "──────────────────────────────────────────────────────────────────────────────"
if [[ ${#WARNINGS[@]} -gt 0 ]]; then
    echo "⚠️  ${#WARNINGS[@]} warning(s):"
    for warning in "${WARNINGS[@]}"; do
        echo "     - $warning"
    done
fi
if [[ ${#FAILURES[@]} -gt 0 ]]; then
    echo "❌ ${#FAILURES[@]} problem(s) that will stop the release:"
    for failure in "${FAILURES[@]}"; do
        echo "     - $failure"
    done
    exit 1
fi
echo "✅ Everything the release needs is in place."
