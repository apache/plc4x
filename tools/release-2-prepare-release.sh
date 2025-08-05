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

########################################################################################################################
# 0. Check if there are uncommitted changes as these would automatically be committed (local)
########################################################################################################################

if [[ $(git -C "$DIRECTORY" status --porcelain) ]]; then
  # Changes
  echo "❌ There are untracked files or changed files, aborting."
  exit 1
fi

PROJECT_VERSION=$("$DIRECTORY"/mvnw -f "$DIRECTORY"/pom.xml -q -Dexec.executable=echo -Dexec.args="\${project.version}" --non-recursive exec:exec)
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
        bash -c "/ws/mvnw -e -Dmaven.repo.local=/ws/out/.repository -DaltDeploymentRepository=snapshot-repo::default::file:/ws/out/.local-artifacts-dir release:perform"; then
    echo "❌ Got non-0 exit code from docker compose, aborting."
    exit 1
fi

########################################################################################################################
# 4. Sign all artifacts
########################################################################################################################

echo "Signing artifacts:"
find $DIRECTORY/out/.local-artifacts-dir -print | grep -E '^((.*\.pom)|(.*\.jar)|(.*\.kar)|(.*\.nar)|(.*-features\.xml)|(.*-cyclonedx\.json)|(.*-cyclonedx\.xml)|(.*-site\.xml)|(.*\.zip))$' | while read -r line ; do
    echo "Processing $line"
    if ! gpg -ab "$line"; then
        echo "❌ Got non-0 exit code from signing artifact, aborting."
        exit 1
    fi
done

########################################################################################################################
# 5. Deploy the artifacts to Nexus and close the staging repo
########################################################################################################################

# If this doesn't work and results in errors complaining about "404 not found", the stating profile id
# might have changed. To update that, go to https://repository.apache.org and log in. Then select "Staging Profiles"
# and click on the one named "org.apache.plc4x". The browser URL will be updated to something like:
# https://repository.apache.org/#stagingProfiles;15cd9d785359f8 ... the id after the "#stagingProfiles;" is the
# staging profile id.
echo "Deploying artifacts:"
STAGING_PROFILE_ID=15cd9d785359f8
# Clean up any pre-existing properties file, as otherwise we'll also deploy that,
# and that will cause errors when closing.
rm "$DIRECTORY/out/.local-artifacs-dir/$STAGING_PROFILE_ID.properties"
if ! "$DIRECTORY/mvnw" -f "$DIRECTORY/tools/stage.pom" nexus-staging:deploy-staged-repository -DstagingProfileId=$STAGING_PROFILE_ID; then
    echo "❌ Got non-0 exit code from staging artifacts, aborting."
    exit 1
fi

# Get the url of the closed repo to add that to the VOTE email
DEPLOY_PROPS="$DIRECTORY/out/.local-artifacts-dir/$STAGING_PROFILE_ID.properties"
STAGING_REPO_ID=$(grep stagingRepository.id "$DEPLOY_PROPS" | cut -d= -f2)
NEXUS_URL="https://repository.apache.org"
STAGING_REPO_URL="$NEXUS_URL/content/repositories/$STAGING_REPO_ID"
echo "✅ Staging repository closed: $STAGING_REPO_URL"

########################################################################################################################
# 7. Prepare a directory for the release candidate
########################################################################################################################

echo "Staging release candidate:"
read -r -p 'Release-Candidate number: ' rcNumber
RELEASE_CANDIDATE="rc$rcNumber"
RELEASE_VERSION=$(find "$DIRECTORY/out/.local-artifacts-dir/org/apache/plc4x/plc4x-parent/" -maxdepth 1 -type d | grep -vE 'plc4x-parent/$' | xargs -n 1 basename)
mkdir -p "$DIRECTORY/out/stage/$RELEASE_VERSION/$RELEASE_CANDIDATE"
cp "$DIRECTORY/README.md" "$DIRECTORY/out/stage/$RELEASE_VERSION/$RELEASE_CANDIDATE/README"
cp "$DIRECTORY/RELEASE_NOTES" "$DIRECTORY/out/stage/$RELEASE_VERSION/$RELEASE_CANDIDATE"
cp "$DIRECTORY/out/.local-artifacts-dir/org/apache/plc4x/plc4x-parent/$RELEASE_VERSION/plc4x-parent-$RELEASE_VERSION-source-release.zip" "$DIRECTORY/out/stage/$RELEASE_VERSION/$RELEASE_CANDIDATE/apache-plc4x-$RELEASE_VERSION-source-release.zip"
cp "$DIRECTORY/out/.local-artifacts-dir/org/apache/plc4x/plc4x-parent/$RELEASE_VERSION/plc4x-parent-$RELEASE_VERSION-source-release.zip.asc" "$DIRECTORY/out/stage/$RELEASE_VERSION/$RELEASE_CANDIDATE/apache-plc4x-$RELEASE_VERSION-source-release.zip.asc"
cp "$DIRECTORY/out/.local-artifacts-dir/org/apache/plc4x/plc4x-parent/$RELEASE_VERSION/plc4x-parent-$RELEASE_VERSION-source-release.zip.sha512" "$DIRECTORY/out/stage/$RELEASE_VERSION/$RELEASE_CANDIDATE/apache-plc4x-$RELEASE_VERSION-source-release.zip.sha512"
cp "$DIRECTORY/out/.local-artifacts-dir/org/apache/plc4x/plc4x-parent/$RELEASE_VERSION/plc4x-parent-$RELEASE_VERSION-cyclonedx.json" "$DIRECTORY/out/stage/$RELEASE_VERSION/$RELEASE_CANDIDATE/apache-plc4x-$RELEASE_VERSION-cyclonedx.json"
cp "$DIRECTORY/out/.local-artifacts-dir/org/apache/plc4x/plc4x-parent/$RELEASE_VERSION/plc4x-parent-$RELEASE_VERSION-cyclonedx.json.asc" "$DIRECTORY/out/stage/$RELEASE_VERSION/$RELEASE_CANDIDATE/apache-plc4x-$RELEASE_VERSION-cyclonedx.json.asc"
cp "$DIRECTORY/out/.local-artifacts-dir/org/apache/plc4x/plc4x-parent/$RELEASE_VERSION/plc4x-parent-$RELEASE_VERSION-cyclonedx.xml" "$DIRECTORY/out/stage/$RELEASE_VERSION/$RELEASE_CANDIDATE/apache-plc4x-$RELEASE_VERSION-cyclonedx.xml"
cp "$DIRECTORY/out/.local-artifacts-dir/org/apache/plc4x/plc4x-parent/$RELEASE_VERSION/plc4x-parent-$RELEASE_VERSION-cyclonedx.xml.asc" "$DIRECTORY/out/stage/$RELEASE_VERSION/$RELEASE_CANDIDATE/apache-plc4x-$RELEASE_VERSION-cyclonedx.xml.asc"

########################################################################################################################
# 8. Make sure the currently used GPG key is available in the KEYS file
########################################################################################################################

ORIGINAL_FILE="$DIRECTORY/out/stage/$RELEASE_VERSION/$RELEASE_CANDIDATE/apache-plc4x-$RELEASE_VERSION-source-release.zip"
SIGNATURE_FILE="$DIRECTORY/out/stage/$RELEASE_VERSION/$RELEASE_CANDIDATE/apache-plc4x-$RELEASE_VERSION-source-release.zip.asc"

KEYS_URL="https://dist.apache.org/repos/dist/release/plc4x/KEYS"
TEMP_DIR=$(mktemp -d)
KEYS_FILE="$TEMP_DIR/KEYS"
KEYRING="$TEMP_DIR/pubring.kbx"

# Fetch KEYS file
echo "🔽 Downloading KEYS file from $KEYS_URL"
curl -fsSL "$KEYS_URL" -o "$KEYS_FILE"

# Import keys into temporary keyring
echo "🔑 Importing KEYS into temporary GPG keyring"
gpg --no-default-keyring --keyring "$KEYRING" --import "$KEYS_FILE" > /dev/null 2>&1

# Verify the signature
echo "🧾 Verifying signature on $ORIGINAL_FILE with $SIGNATURE_FILE"
if gpg --no-default-keyring --keyring "$KEYRING" --verify "$SIGNATURE_FILE" "$ORIGINAL_FILE" > /dev/null 2>&1; then
    echo "✅ Signature is valid and signed by a key in the Apache PLC4X KEYS file"
else
    echo "❌ Signature is invalid or the key is not in the Apache PLC4X KEYS file"
    exit 1
fi

# TODO: Check that the signature references an apache email address...

# Cleanup
rm -rf "$TEMP_DIR"

########################################################################################################################
# 9. Validate the sha512 hashes
########################################################################################################################

ORIGINAL_FILE="$DIRECTORY/out/stage/$RELEASE_VERSION/$RELEASE_CANDIDATE/apache-plc4x-$RELEASE_VERSION-source-release.zip"
HASHES_FILE="$DIRECTORY/out/stage/$RELEASE_VERSION/$RELEASE_CANDIDATE/apache-plc4x-$RELEASE_VERSION-source-release.zip.sha512"

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
svn import "$RELEASE_CANDIDATE" "https://dist.apache.org/repos/dist/dev/plc4x/$RELEASE_VERSION/$RELEASE_CANDIDATE" -m"Staging of $RELEASE_CANDIDATE of PLC4X $RELEASE_VERSION"

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
[2] https://dist.apache.org/repos/dist/dev/plc4x/$RELEASE_VERSION/$RELEASE_CANDIDATE
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
