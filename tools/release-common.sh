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

# Values that more than one of the release scripts has to agree on. Sourced by
# "release-2-prepare-release.sh", "release-3-finish-release.sh" and "validate-release.sh" - it is
# not meant to be executed.

# The Nexus staging profile of "org.apache.plc4x". If deploying to Nexus starts failing with
# "404 not found", this is the first thing to check: log in to $NEXUS_URL, open "Staging Profiles",
# select "org.apache.plc4x" and take the id out of the browser URL after "#stagingProfiles;".
STAGING_PROFILE_ID=15cd9d785359f8

NEXUS_URL="https://repository.apache.org"

# Where release candidates are staged and where releases end up. Everything below
# https://dist.apache.org/repos/dist/ needs an Apache committer account to write to.
DIST_BASE="https://dist.apache.org/repos/dist"
DIST_DEV="$DIST_BASE/dev/plc4x"
DIST_RELEASE="$DIST_BASE/release/plc4x"
KEYS_URL="$DIST_RELEASE/KEYS"

# At the ASF the same LDAP account is behind the Nexus staging repository and the dist.apache.org
# SVN, so the credentials Maven already has in "settings.xml" are the ones svn keeps asking for.
# "read_asf_credentials" digs them out, "svn_authenticated" uses them if they are there.
#
# Nothing is lost when they cannot be read: the functions fall back to letting svn prompt, which
# is what the release scripts did before. A password that Maven has encrypted ("{...}") is such a
# case - decrypting that here would mean reimplementing Maven's key handling.
ASF_USERNAME=""
ASF_PASSWORD=""

read_asf_credentials() {
    local settings="${MAVEN_SETTINGS:-$HOME/.m2/settings.xml}"
    local server_id="apache.releases.https"

    if [[ ! -f "$settings" ]] || ! command -v xmllint > /dev/null 2>&1; then
        return 0
    fi

    local username password
    username=$(xmllint --xpath "string(//*[local-name()='server'][*[local-name()='id']='$server_id']/*[local-name()='username'])" "$settings" 2>/dev/null)
    password=$(xmllint --xpath "string(//*[local-name()='server'][*[local-name()='id']='$server_id']/*[local-name()='password'])" "$settings" 2>/dev/null)

    if [[ -z "$username" || -z "$password" || "$password" == \{* ]]; then
        return 0
    fi

    ASF_USERNAME="$username"
    ASF_PASSWORD="$password"
}

# Runs an svn command with those credentials, passing the password through stdin so that it never
# shows up in the process list. If they turn out not to work - the Nexus credentials can be a user
# token, which the SVN does not know about - the same command is run again interactively.
svn_authenticated() {
    if [[ -n "$ASF_USERNAME" && -n "$ASF_PASSWORD" ]]; then
        if printf '%s' "$ASF_PASSWORD" \
                | svn --username "$ASF_USERNAME" --password-from-stdin --non-interactive "$@"; then
            return 0
        fi
        echo "⚠️  The credentials from settings.xml were not accepted by svn, asking instead."
    fi
    svn "$@"
}
