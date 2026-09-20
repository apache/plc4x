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
# shellcheck shell=bash
# shellcheck disable=SC2034  # these are consumed by the scripts that source this file

# Values that more than one of the release scripts has to agree on, and the ownership settings
# every script that starts the releaser container needs. Sourced by all of the "release-*"
# scripts and by "validate-release.sh" - it is not meant to be executed.

# The Nexus staging profile of "org.apache.plc4x". If deploying to Nexus starts failing with
# "404 not found", this is the first thing to check: log in to $NEXUS_URL, open "Staging Profiles",
# select "org.apache.plc4x" and take the id out of the browser URL after "#stagingProfiles;".
STAGING_PROFILE_ID=15cd9d785359f8

NEXUS_URL="https://repository.apache.org"

# The nexus-staging-maven-plugin serializes its requests with XStream, which reflects into
# "java.util" collection classes. Since JDK 16 that module is not open to the unnamed module, so
# every goal of the plugin that sends a list - "rc-release", "rc-drop", ... - dies client-side with
# "No converter available ... java.util.Arrays$ArrayList" before anything reaches Nexus. Opening
# the package for the Maven JVM is the only workaround; the plugin has not been released since.
# Kept separate from any MAVEN_OPTS already set, which is prepended so it still wins.
NEXUS_MAVEN_OPTS="${MAVEN_OPTS:+$MAVEN_OPTS }--add-opens java.base/java.util=ALL-UNNAMED"

# Where release candidates are staged and where releases end up. Everything below
# https://dist.apache.org/repos/dist/ needs an Apache committer account to write to.
DIST_BASE="https://dist.apache.org/repos/dist"
DIST_DEV="$DIST_BASE/dev/plc4x"
DIST_RELEASE="$DIST_BASE/release/plc4x"
KEYS_URL="$DIST_RELEASE/KEYS"

# The docker-compose file runs the releaser container as "$RELEASE_UID:$RELEASE_GID", defaulting to
# "0:0". Only Linux needs the override: there the bind-mounted checkout keeps the uid the container
# used, so a build as root leaves an "out" directory, "target" directories and release commits that
# the user running these scripts can no longer delete. Docker Desktop on macOS and Windows remaps
# the ownership by itself, so the default is left alone there.
if [[ "$(uname)" == "Linux" ]]; then
  RELEASE_UID="$(id -u)"
  RELEASE_GID="$(id -g)"
  export RELEASE_UID RELEASE_GID
fi

# HOME for the container, see the docker-compose file. Created here because the container may not be
# able to create it itself: the mount root belongs to the user, but a missing directory would have to
# be created by whoever runs the build.
mkdir -p "$DIRECTORY/out/home"

# ----------------------------------------------------------------------------------------------
# Finding the Apache remote
# ----------------------------------------------------------------------------------------------
# Not everybody clones the Apache repository as "origin". The contributor workflow the website
# documents has you fork on GitHub, clone the fork as "origin" and add the Apache repository under
# a second name - usually "upstream". Every "git push origin" in the release scripts then pushed
# the release branch and the release tag into that fork instead. So rather than trusting a name,
# the remotes are matched by URL against the canonical repository, on gitbox or on the GitHub
# mirror, over ssh or over https.

# The repository this checkout is a clone of, as it is named under "apache/" and "repos/asf/".
APACHE_REPO_NAME="plc4x"

# Filled in by "resolve_apache_remote": the name and URL of the remote that points at the Apache
# repository, or an empty name and the reason in APACHE_REMOTE_PROBLEM.
APACHE_REMOTE=""
APACHE_REMOTE_URL=""
APACHE_REMOTE_PROBLEM=""
# How many remotes matched, so that the help below does not tell somebody who already has the
# repository as a remote to add it again.
APACHE_REMOTE_MATCHES=0

resolve_apache_remote() {
    APACHE_REMOTE=""
    APACHE_REMOTE_URL=""
    APACHE_REMOTE_PROBLEM=""
    APACHE_REMOTE_MATCHES=0

    # An explicit override wins over the URL match, for everything the match cannot know about:
    # a mirror, a private clone a release is rehearsed in, a host that is reached through a proxy.
    if [[ -n "${PLC4X_REMOTE:-}" ]]; then
        if APACHE_REMOTE_URL=$(git -C "$DIRECTORY" remote get-url "$PLC4X_REMOTE" 2>/dev/null); then
            APACHE_REMOTE="$PLC4X_REMOTE"
        else
            APACHE_REMOTE_URL=""
            APACHE_REMOTE_PROBLEM="PLC4X_REMOTE is set to '$PLC4X_REMOTE', but this checkout has no remote of that name."
        fi
        return
    fi

    # Anchored at the end so that, in the "plc4x" checkout, "apache/plc4x-extras" is not taken for
    # "apache/plc4x". The optional ".git" and trailing slash cover the ways the URL gets written.
    local pattern="(gitbox\.apache\.org[:/]repos/asf/${APACHE_REPO_NAME}|github\.com[:/]apache/${APACHE_REPO_NAME})(\.git)?/?$"
    local matches=()
    local name url
    while read -r name url; do
        [[ -z "$name" ]] && continue
        if [[ "$url" =~ $pattern ]]; then
            matches+=("$name")
        fi
    done < <(git -C "$DIRECTORY" remote -v 2>/dev/null | awk '$3 == "(fetch)" { print $1, $2 }')

    APACHE_REMOTE_MATCHES=${#matches[@]}
    case ${#matches[@]} in
        0)
            APACHE_REMOTE_PROBLEM="No remote of this checkout points at the Apache '$APACHE_REPO_NAME' repository."
            ;;
        1)
            APACHE_REMOTE="${matches[0]}"
            APACHE_REMOTE_URL=$(git -C "$DIRECTORY" remote get-url "$APACHE_REMOTE" 2>/dev/null)
            ;;
        *)
            # Deliberately not guessed at: which of them the release ends up in is not a decision
            # a release script should make on its own.
            APACHE_REMOTE_PROBLEM="Several remotes point at the Apache '$APACHE_REPO_NAME' repository: ${matches[*]}."
            ;;
    esac
}

# Prints what the remotes of this checkout are and the two ways out of it. Only useful once
# "resolve_apache_remote" has failed.
apache_remote_help() {
    local remotes
    remotes=$(git -C "$DIRECTORY" remote -v 2>/dev/null | awk '$3 == "(fetch)" { printf "     %-12s %s\n", $1, $2 }')
    echo "   Remotes of this checkout:"
    if [[ -n "$remotes" ]]; then
        echo "$remotes"
    else
        echo "     (none)"
    fi
    if [[ "$APACHE_REMOTE_MATCHES" -eq 0 ]]; then
        echo "   Add the Apache repository as a remote:"
        echo "     git remote add apache https://github.com/apache/$APACHE_REPO_NAME.git"
        echo "   Or name the remote the release should use explicitly:"
    else
        echo "   Name the remote the release should use explicitly:"
    fi
    echo "     export PLC4X_REMOTE=<remote-name>"
}

# For the scripts that push, fetch or query the remote: resolve it, or stop before anything is
# done in the wrong repository.
require_apache_remote() {
    resolve_apache_remote
    if [[ -z "$APACHE_REMOTE" ]]; then
        echo "❌ $APACHE_REMOTE_PROBLEM"
        apache_remote_help
        exit 1
    fi
}

# Resolved once here so that every script that sources this file has APACHE_REMOTE available.
# Scripts that actually talk to the remote call "require_apache_remote" to insist on it.
resolve_apache_remote
