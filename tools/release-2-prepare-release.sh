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

DIRECTORY="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# Values shared with the other release scripts (Nexus staging profile, dist.apache.org URLs).
if [[ ! -f "$DIRECTORY/tools/release-common.sh" ]]; then
    echo "❌ '$DIRECTORY/tools/release-common.sh' not found, aborting."
    exit 1
fi
# shellcheck source=release-common.sh
source "$DIRECTORY/tools/release-common.sh"


########################################################################################################################
# 0. Check if there are uncommitted changes as these would automatically be committed (local)
########################################################################################################################

if [[ $(git -C "$DIRECTORY" status --porcelain) ]]; then
  # Changes
  echo "❌ There are untracked files or changed files, aborting."
  exit 1
fi

# Maven 4 prefixes even quiet output with "[INFO] [stdout] ", so take the last token of the
# last line rather than the whole output.
PROJECT_VERSION=$("$DIRECTORY"/mvnw -f "$DIRECTORY"/pom.xml -q -Dexec.executable=echo -Dexec.args="\${project.version}" --non-recursive exec:exec | tail -n 1 | awk '{print $NF}')
if [[ ! "$PROJECT_VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+(-SNAPSHOT)?$ ]]; then
    echo "❌ Could not read a usable project version, got '$PROJECT_VERSION'."
    echo "   Everything below derives the branch and tag names from it, so aborting."
    exit 1
fi
RELEASE_VERSION=${PROJECT_VERSION%"-SNAPSHOT"}
TAG_NAME="v$RELEASE_VERSION"
IFS='.' read -ra VERSION_SEGMENTS <<< "$RELEASE_VERSION"
NEW_VERSION="${VERSION_SEGMENTS[0]}.${VERSION_SEGMENTS[1]}.$((VERSION_SEGMENTS[2] + 1))-SNAPSHOT"

# Check if a local tag already exists (This can happen if a first release attempt failed)
if git -C "$DIRECTORY" rev-parse "$TAG_NAME" >/dev/null 2>&1; then
  echo "❌ Tag '$TAG_NAME' exists locally. Please delete with 'git tag -d $TAG_NAME'"
  exit 1
else
  echo "✅ Tag '$TAG_NAME' does not exist locally."
fi

# Check if a remote tag already exists (This can happen if a first release attempt failed)
if git -C "$DIRECTORY" ls-remote --tags origin | grep -q "refs/tags/$TAG_NAME$"; then
  echo "❌ Tag '$TAG_NAME' exists on remote 'origin'. Please delete with 'git push origin --delete $TAG_NAME'"
  exit 1
else
  echo "✅ Tag '$TAG_NAME' does not exist on remote 'origin'."
fi

# The Antora documentation version of this branch was already set by 'release-1-create-branch.sh'.

########################################################################################################################
# 1. Do a simple release-prepare command
########################################################################################################################

if ! docker compose -f "$DIRECTORY/tools/docker-compose.yaml" run releaser \
        bash -c "/ws/mvnw -e -P with-c,with-dotnet,with-go,with-java,with-python,enable-all-checks,update-generated-code -Dmaven.repo.local=/ws/out/.repository release:prepare -DautoVersionSubmodules=true -DreleaseVersion='$RELEASE_VERSION' -DdevelopmentVersion='$NEW_VERSION' -Dtag='$TAG_NAME'"; then
    echo "❌ Got non-0 exit code from docker compose, aborting."
    exit 1
fi

########################################################################################################################
# 2. Push the changes (local)
########################################################################################################################

if ! git -C "$DIRECTORY" push; then
    echo "❌ Got non-0 exit code from pushing changes to git, aborting."
    exit 1
fi

TAG_COMMIT_HASH=$(git -C "$DIRECTORY" rev-list -n 1 "$TAG_NAME")
echo "✅ Tag '$TAG_NAME' has hash '$TAG_COMMIT_HASH'"

########################################################################################################################
# 3. Do a simple release-perform command skip signing of artifacts and deploy to local directory (inside the Docker container)
########################################################################################################################

echo "Performing Release:"
if ! docker compose -f "$DIRECTORY/tools/docker-compose.yaml" run releaser \
        bash -c "/ws/mvnw -e -Dmaven.repo.local=/ws/out/.repository -DaltDeploymentRepository=snapshot-repo::file:/ws/out/.local-artifacts-dir release:perform"; then
    echo "❌ Got non-0 exit code from docker compose, aborting."
    exit 1
fi

########################################################################################################################
# 4. Sign all artifacts
########################################################################################################################

echo "Signing artifacts:"
SIGNED_ARTIFACTS=0
# The loop reads from a process substitution rather than from a pipe on purpose: the last stage
# of a pipeline runs in a subshell, where an "exit" would only end the loop and let the release
# carry on to the deployment below with unsigned artifacts.
while read -r line ; do
    echo "Processing $line"
    if ! gpg -ab "$line"; then
        echo "❌ Got non-0 exit code from signing artifact, aborting."
        exit 1
    fi
    SIGNED_ARTIFACTS=$((SIGNED_ARTIFACTS + 1))
done < <(find "$DIRECTORY/out/.local-artifacts-dir" -print | grep -E '^((.*\.pom)|(.*\.jar)|(.*\.kar)|(.*\.nar)|(.*-features\.xml)|(.*-cyclonedx\.json)|(.*-cyclonedx\.xml)|(.*-site\.xml)|(.*\.zip))$')

if [[ "$SIGNED_ARTIFACTS" -eq 0 ]]; then
    echo "❌ Found no artifacts to sign in '$DIRECTORY/out/.local-artifacts-dir', aborting."
    exit 1
fi
echo "✅ Signed $SIGNED_ARTIFACTS artifacts."

########################################################################################################################
# 5. Deploy the artifacts to Nexus and close the staging repo
########################################################################################################################

echo "Deploying artifacts:"
# Clean up any pre-existing properties file, as otherwise we'll also deploy that,
# and that will cause errors when closing.
rm "$DIRECTORY/out/.local-artifacts-dir/$STAGING_PROFILE_ID.properties"
if ! "$DIRECTORY/mvnw" -f "$DIRECTORY/tools/stage.pom" nexus-staging:deploy-staged-repository -DstagingProfileId=$STAGING_PROFILE_ID; then
    echo "❌ Got non-0 exit code from staging artifacts, aborting."
    exit 1
fi

# Get the url of the closed repo to add that to the VOTE email
DEPLOY_PROPS="$DIRECTORY/out/.local-artifacts-dir/$STAGING_PROFILE_ID.properties"
STAGING_REPO_ID=$(grep stagingRepository.id "$DEPLOY_PROPS" | cut -d= -f2)
STAGING_REPO_URL="$NEXUS_URL/content/repositories/$STAGING_REPO_ID"
echo "✅ Staging repository closed: $STAGING_REPO_URL"

########################################################################################################################
# 7. Prepare a directory for the release candidate
########################################################################################################################

echo "Staging release candidate:"
read -r -p 'Release-Candidate number: ' rcNumber
RELEASE_CANDIDATE="rc$rcNumber"

# RELEASE_VERSION is the one derived from the project version at the top of this script. It used
# to be re-derived here by listing the built artifacts, which silently produced several words -
# and with it a set of broken paths - whenever "out" still held a previous attempt. Rather than
# guess the version from whatever is lying around, check that the build we just did is there.
ARTIFACTS_DIR="$DIRECTORY/out/.local-artifacts-dir/org/apache/plc4x/plc4x-parent/$RELEASE_VERSION"
if [[ ! -d "$ARTIFACTS_DIR" ]]; then
    echo "❌ Found no built artifacts for $RELEASE_VERSION in '$ARTIFACTS_DIR', aborting."
    echo "   These versions are present:"
    find "$DIRECTORY/out/.local-artifacts-dir/org/apache/plc4x/plc4x-parent/" -maxdepth 1 -mindepth 1 -type d -exec basename {} \; 2>/dev/null | sed 's|^|     |'
    exit 1
fi

STAGE_DIR="$DIRECTORY/out/stage/$RELEASE_VERSION/$RELEASE_CANDIDATE"
mkdir -p "$STAGE_DIR"
cp "$DIRECTORY/README.md" "$STAGE_DIR/README"
cp "$DIRECTORY/RELEASE_NOTES" "$STAGE_DIR"
cp "$ARTIFACTS_DIR/plc4x-parent-$RELEASE_VERSION-source-release.zip" "$STAGE_DIR/apache-plc4x-$RELEASE_VERSION-source-release.zip"
cp "$ARTIFACTS_DIR/plc4x-parent-$RELEASE_VERSION-source-release.zip.asc" "$STAGE_DIR/apache-plc4x-$RELEASE_VERSION-source-release.zip.asc"
cp "$ARTIFACTS_DIR/plc4x-parent-$RELEASE_VERSION-source-release.zip.sha512" "$STAGE_DIR/apache-plc4x-$RELEASE_VERSION-source-release.zip.sha512"
cp "$ARTIFACTS_DIR/plc4x-parent-$RELEASE_VERSION-cyclonedx.json" "$STAGE_DIR/apache-plc4x-$RELEASE_VERSION-cyclonedx.json"
cp "$ARTIFACTS_DIR/plc4x-parent-$RELEASE_VERSION-cyclonedx.json.asc" "$STAGE_DIR/apache-plc4x-$RELEASE_VERSION-cyclonedx.json.asc"
cp "$ARTIFACTS_DIR/plc4x-parent-$RELEASE_VERSION-cyclonedx.xml" "$STAGE_DIR/apache-plc4x-$RELEASE_VERSION-cyclonedx.xml"
cp "$ARTIFACTS_DIR/plc4x-parent-$RELEASE_VERSION-cyclonedx.xml.asc" "$STAGE_DIR/apache-plc4x-$RELEASE_VERSION-cyclonedx.xml.asc"

########################################################################################################################
# 8. Make sure the currently used GPG key is available in the KEYS file
########################################################################################################################

ORIGINAL_FILE="$STAGE_DIR/apache-plc4x-$RELEASE_VERSION-source-release.zip"
SIGNATURE_FILE="$STAGE_DIR/apache-plc4x-$RELEASE_VERSION-source-release.zip.asc"

# KEYS_URL comes from tools/release-common.sh
TEMP_DIR=$(mktemp -d)
KEYS_FILE="$TEMP_DIR/KEYS"
KEYRING="$TEMP_DIR/pubring.kbx"

# Fetch KEYS file
echo "🔽 Downloading KEYS file from $KEYS_URL"
if ! curl -fsSL "$KEYS_URL" -o "$KEYS_FILE"; then
    echo "❌ Could not download the KEYS file from $KEYS_URL, aborting."
    rm -rf "$TEMP_DIR"
    exit 1
fi
if [[ ! -s "$KEYS_FILE" ]]; then
    echo "❌ The KEYS file downloaded from $KEYS_URL is empty, aborting."
    rm -rf "$TEMP_DIR"
    exit 1
fi

# Import keys into temporary keyring
echo "🔑 Importing KEYS into temporary GPG keyring"
if ! gpg --no-default-keyring --keyring "$KEYRING" --import "$KEYS_FILE" > /dev/null 2>&1; then
    echo "❌ Could not import the KEYS file into a temporary keyring, aborting."
    rm -rf "$TEMP_DIR"
    exit 1
fi
# Without this an empty keyring would make the verification below fail, and the message would
# blame the release manager's key instead of the KEYS file that never arrived.
IMPORTED_KEYS=$(gpg --no-default-keyring --keyring "$KEYRING" --with-colons --list-keys 2>/dev/null | grep -c '^pub:')
if [[ "$IMPORTED_KEYS" -eq 0 ]]; then
    echo "❌ The temporary keyring contains no keys after importing $KEYS_URL, aborting."
    rm -rf "$TEMP_DIR"
    exit 1
fi
echo "🔑 Imported $IMPORTED_KEYS keys from the KEYS file"

# Verify the signature
echo "🧾 Verifying signature on $ORIGINAL_FILE with $SIGNATURE_FILE"
if gpg --no-default-keyring --keyring "$KEYRING" --verify "$SIGNATURE_FILE" "$ORIGINAL_FILE" > /dev/null 2>&1; then
    echo "✅ Signature is valid and signed by a key in the Apache PLC4X KEYS file"
else
    echo "❌ Signature is invalid or the key is not in the Apache PLC4X KEYS file"
    exit 1
fi

########################################################################################################################
# 8b. Make sure that key belongs to an Apache identity
########################################################################################################################

# Being in the KEYS file is not quite enough: the key used to sign an Apache release has to carry
# the release manager's {apache-id}@apache.org address. A key can have several user ids, and only
# one of them has to be the Apache one.

SIGNATURE_STATUS=$(gpg --no-default-keyring --keyring "$KEYRING" --status-fd 1 \
    --verify "$SIGNATURE_FILE" "$ORIGINAL_FILE" 2>/dev/null)
# Field 3 of VALIDSIG is the key that made the signature, which is the signing subkey for anyone
# who signs with one. The last field is the fingerprint of the primary key, which is what carries
# the user ids, so prefer that and fall back to field 3 for older gpg output.
SIGNING_KEY=$(echo "$SIGNATURE_STATUS" | awk '/^\[GNUPG:\] VALIDSIG /{print ($12 != "" ? $12 : $3); exit}')
if [[ -z "$SIGNING_KEY" ]]; then
    echo "❌ Could not work out which key signed the release, aborting."
    rm -rf "$TEMP_DIR"
    exit 1
fi

# Field 2 is the validity of the user id: skip the ones that are revoked ("r") or expired ("e"),
# so that an @apache.org address the key no longer stands behind does not satisfy the check.
SIGNING_UIDS=$(gpg --no-default-keyring --keyring "$KEYRING" --with-colons --list-keys "$SIGNING_KEY" 2>/dev/null \
    | awk -F: '$1 == "uid" && $2 !~ /^[re]$/ {print $10}')
if echo "$SIGNING_UIDS" | grep -qiE '<[^>]+@apache\.org>'; then
    echo "✅ The signing key is registered to an apache.org address"
else
    echo "❌ The signing key is not registered to an apache.org address, aborting."
    echo "   Releases have to be signed with a key carrying your {apache-id}@apache.org address."
    echo "   Key $SIGNING_KEY currently carries:"
    echo "$SIGNING_UIDS" | sed 's|^|     |'
    rm -rf "$TEMP_DIR"
    exit 1
fi

# Cleanup
rm -rf "$TEMP_DIR"

########################################################################################################################
# 9. Validate the sha512 hashes
########################################################################################################################

ORIGINAL_FILE="$STAGE_DIR/apache-plc4x-$RELEASE_VERSION-source-release.zip"
HASHES_FILE="$STAGE_DIR/apache-plc4x-$RELEASE_VERSION-source-release.zip.sha512"

ACTUAL_HASH=$(shasum -a 512 "$ORIGINAL_FILE" | awk '{ print $1 }')
EXPECTED_HASH=$(tr -d ' \n\r' < "$HASHES_FILE")

if [ "$ACTUAL_HASH" = "$EXPECTED_HASH" ]; then
    echo "✅ SHA-512 hash matches"
else
    echo "❌ SHA-512 hash does not match"
    echo "Expected: $EXPECTED_HASH"
    echo "Actual:   $ACTUAL_HASH"
    exit 1
fi

########################################################################################################################
# 10. Upload the release candidate artifacts to SVN
########################################################################################################################

cd "$DIRECTORY/out/stage/$RELEASE_VERSION" || exit
svn import "$RELEASE_CANDIDATE" "$DIST_DEV/$RELEASE_VERSION/$RELEASE_CANDIDATE" -m"Staging of $RELEASE_CANDIDATE of PLC4X $RELEASE_VERSION"

########################################################################################################################
# 11. Prepare the [VOTE] and [DISCUSS] emails
########################################################################################################################

cat > "$DIRECTORY/out/stage/vote-email.eml" <<EOF
To: dev@plc4x.apache.org
Subject: [VOTE] Apache PLC4X $RELEASE_VERSION $RELEASE_CANDIDATE
Content-Type: text/plain; charset=UTF-8

Apache PLC4X $RELEASE_VERSION has been staged under [2] and it’s time to vote
on accepting it for release. All Maven artifacts are available under [1].
Voting will be open for 72hr.

A minimum of 3 binding +1 votes and more binding +1 than binding -1
are required to pass.

Release tag: v$RELEASE_VERSION
Hash for the release tag: $TAG_COMMIT_HASH

Per [3] "Before voting +1 PMC members are required to download
the signed source code package, compile it as provided, and test
the resulting executable on their own platform, along with also
verifying that the package meets the requirements of the ASF policy
on releases."

You can achieve the above by following [4].

[ ]  +1 accept (indicate what you validated - e.g. performed the non-RM items in [4])
[ ]  -1 reject (explanation required)


[1] $STAGING_REPO_URL
[2] $DIST_DEV/$RELEASE_VERSION/$RELEASE_CANDIDATE
[3] https://www.apache.org/dev/release.html#approving-a-release
[4] https://plc4x.apache.org/plc4x/latest/developers/release/validation.html
EOF
echo "✅ Vote email generated to $DIRECTORY/out/stage/vote-email.eml"

cat > "$DIRECTORY/out/stage/discuss-email.eml" <<EOF
To: dev@plc4x.apache.org
Subject: [DISCUSS] Apache PLC4X $RELEASE_VERSION $RELEASE_CANDIDATE
Content-Type: text/plain; charset=UTF-8

This is the discussion thread for the corresponding VOTE thread.

Please keep discussions in this thread to simplify the counting of votes.

If you have to vote -1 please mention a brief description on why and then take the details to this thread.
EOF
echo "✅ Discuss email generated to $DIRECTORY/out/stage/discuss-email.eml"
