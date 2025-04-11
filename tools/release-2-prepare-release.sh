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

DIRECTORY=$(pwd)

########################################################################################################################
# 0. Check if there are uncommitted changes as these would automatically be committed (local)
########################################################################################################################

if [[ $(git status --porcelain) ]]; then
  # Changes
  echo "❌ There are untracked files or changed files, aborting."
  exit 1
fi

PROJECT_VERSION=$(../mvnw -f ../pom.xml -q -Dexec.executable=echo -Dexec.args="\${project.version}" --non-recursive exec:exec)
RELEASE_VERSION=${PROJECT_VERSION%"-SNAPSHOT"}
IFS='.' read -ra VERSION_SEGMENTS <<< "$RELEASE_VERSION"
NEW_VERSION="${VERSION_SEGMENTS[0]}.${VERSION_SEGMENTS[1]}.$((VERSION_SEGMENTS[2] + 1))-SNAPSHOT"

########################################################################################################################
# 1. Do a simple release-prepare command
########################################################################################################################

if ! docker compose run releaser bash /ws/mvnw -e -P with-c,with-dotnet,with-go,with-java,with-python,enable-all-checks,update-generated-code -Dmaven.repo.local=/ws/out/.repository release:prepare -DautoVersionSubmodules=true -DreleaseVersion="$RELEASE_VERSION" -DdevelopmentVersion="$NEW_VERSION" -Dtag="v$RELEASE_VERSION"; then
    echo "❌ Got non-0 exit code from docker compose, aborting."
    exit 1
fi

########################################################################################################################
# 2. Push the changes (local)
########################################################################################################################

if ! git push; then
    echo "❌ Got non-0 exit code from pushing changes to git, aborting."
    exit 1
fi

########################################################################################################################
# 3. Do a simple release-perform command skip signing of artifacts and deploy to local directory (inside the Docker container)
########################################################################################################################

echo "Performing Release:"
docker compose build
if ! docker compose run releaser bash /ws/mvnw -e -Dmaven.repo.local=/ws/out/.repository -DaltDeploymentRepository=snapshot-repo::default::file:/ws/out/.local-artifacts-dir release:perform; then
    echo "❌ Got non-0 exit code from docker compose, aborting."
    exit 1
fi

########################################################################################################################
# 4. Sign all artifacts
########################################################################################################################

echo "Signing artifacts:"
find ../out/.local-artifacts-dir -print | grep -E '^((.*\.pom)|(.*\.jar)|(.*\.kar)|(.*\.nar)|(.*-features\.xml)|(.*-cyclonedx\.json)|(.*-cyclonedx\.xml)|(.*-site\.xml)|(.*\.zip))$' | while read -r line ; do
    echo "Processing $line"
    if ! gpg -ab "$line"; then
        echo "❌ Got non-0 exit code from signing artifact, aborting."
        exit 1
    fi
done

########################################################################################################################
# 5. Deploy the artifacts to Nexus
########################################################################################################################

echo "Deploying artifacts:"
if ! ../mvnw -f ../jenkins.pom -X -P deploy-releases wagon:upload; then
    echo "❌ Got non-0 exit code from staging artifacts, aborting."
    exit 1
fi

########################################################################################################################
# 6. Prepare a directory for the release candidate
########################################################################################################################

echo "Staging release candidate:"
read -r -p 'Release-Candidate number: ' rcNumber
RELEASE_CANDIDATE="rc$rcNumber"
RELEASE_VERSION=$(find ../out/.local-artifacts-dir/org/apache/plc4x/plc4x-parent/ -maxdepth 1 -type d | grep -vE 'plc4x-parent/$' | xargs -n 1 basename)
mkdir -p "../out/stage/${RELEASE_VERSION}/${RELEASE_CANDIDATE}"
cp ../README.md "../out/stage/${RELEASE_VERSION}/${RELEASE_CANDIDATE}/README"
cp ../RELEASE_NOTES "../out/stage/${RELEASE_VERSION}/${RELEASE_CANDIDATE}"
cp "../out/.local-artifacts-dir/org/apache/plc4x/plc4x-parent/${RELEASE_VERSION}/plc4x-parent-${RELEASE_VERSION}-source-release.zip" "../out/stage/${RELEASE_VERSION}/${RELEASE_CANDIDATE}/apache-plc4x-${RELEASE_VERSION}-source-release.zip"
cp "../out/.local-artifacts-dir/org/apache/plc4x/plc4x-parent/${RELEASE_VERSION}/plc4x-parent-${RELEASE_VERSION}-source-release.zip.asc" "../out/stage/${RELEASE_VERSION}/${RELEASE_CANDIDATE}/apache-plc4x-${RELEASE_VERSION}-source-release.zip.asc"
cp "../out/.local-artifacts-dir/org/apache/plc4x/plc4x-parent/${RELEASE_VERSION}/plc4x-parent-${RELEASE_VERSION}-source-release.zip.sha512" "../out/stage/${RELEASE_VERSION}/${RELEASE_CANDIDATE}/apache-plc4x-${RELEASE_VERSION}-source-release.zip.sha512"
cp "../out/.local-artifacts-dir/org/apache/plc4x/plc4x-parent/${RELEASE_VERSION}/plc4x-parent-${RELEASE_VERSION}-cyclonedx.json" "../out/stage/${RELEASE_VERSION}/${RELEASE_CANDIDATE}/apache-plc4x-${RELEASE_VERSION}-cyclonedx.json"
cp "../out/.local-artifacts-dir/org/apache/plc4x/plc4x-parent/${RELEASE_VERSION}/plc4x-parent-${RELEASE_VERSION}-cyclonedx.json.asc" "../out/stage/${RELEASE_VERSION}/${RELEASE_CANDIDATE}/apache-plc4x-${RELEASE_VERSION}-cyclonedx.json.asc"
cp "../out/.local-artifacts-dir/org/apache/plc4x/plc4x-parent/${RELEASE_VERSION}/plc4x-parent-${RELEASE_VERSION}-cyclonedx.xml" "../out/stage/${RELEASE_VERSION}/${RELEASE_CANDIDATE}/apache-plc4x-${RELEASE_VERSION}-cyclonedx.xml"
cp "../out/.local-artifacts-dir/org/apache/plc4x/plc4x-parent/${RELEASE_VERSION}/plc4x-parent-${RELEASE_VERSION}-cyclonedx.xml.asc" "../out/stage/${RELEASE_VERSION}/${RELEASE_CANDIDATE}/apache-plc4x-${RELEASE_VERSION}-cyclonedx.xml.asc"

########################################################################################################################
# 7. Upload the release candidate artifacts to SVN
########################################################################################################################

cd "../out/stage/${RELEASE_VERSION}" || exit
svn import "${RELEASE_CANDIDATE}" "https://dist.apache.org/repos/dist/dev/plc4x/${RELEASE_VERSION}/${RELEASE_CANDIDATE}" -m"Staging of ${RELEASE_CANDIDATE} of PLC4X ${RELEASE_VERSION}"

########################################################################################################################
# 8. TODO: Make sure the currently used GPG key is available in the KEYS file
########################################################################################################################

########################################################################################################################
# 9. TODO: Close the Nexus staging repository
########################################################################################################################

########################################################################################################################
# 10. TODO: Send out the [VOTE] and [DISCUSS] emails
########################################################################################################################
