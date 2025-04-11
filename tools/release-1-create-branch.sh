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

########################################################################################################################
# 1. Get and calculate the current version (local)
########################################################################################################################

PROJECT_VERSION=$(../mvnw -f ../pom.xml -q -Dexec.executable=echo -Dexec.args="\${project.version}" --non-recursive exec:exec)
RELEASE_VERSION=${PROJECT_VERSION%"-SNAPSHOT"}
RELEASE_SHORT_VERSION=${RELEASE_VERSION%".0"}
BRANCH_NAME="rel/$RELEASE_SHORT_VERSION"
IFS='.' read -ra VERSION_SEGMENTS <<< "$RELEASE_VERSION"
NEW_VERSION="${VERSION_SEGMENTS[0]}.$((VERSION_SEGMENTS[1] + 1)).0-SNAPSHOT"
echo "Current Version: '$PROJECT_VERSION'"
echo "Release Version: '$RELEASE_VERSION'"
echo "Release Branch Name: '$BRANCH_NAME'"
echo "New develop Version: '$NEW_VERSION'"

########################################################################################################################
# 3. Ask if the RELEASE_NOTES have been filled out at all (local)
########################################################################################################################

read -pr "Have the RELEASE_NOTES been updated for this version? (yes/no) " yn
case $yn in
	yes ) echo continuing with the process;;
	no ) echo Please update the RELEASE_NOTES first;
		exit 1;;
	* ) echo invalid response;
		exit 1;;
esac

########################################################################################################################
# 4. Do a simple maven branch command with pushChanges=false
########################################################################################################################

if ! docker compose run releaser bash /ws/mvnw -e -P with-c,with-dotnet,with-go,with-java,with-python,enable-all-checks,update-generated-code -Dmaven.repo.local=/ws/out/.repository release:branch -DautoVersionSubmodules=true -DpushChanges=false -DdevelopmentVersion="$NEW_VERSION" -DbranchName="$BRANCH_NAME"; then
    echo "❌ Got non-0 exit code from docker compose, aborting."
    exit 1
fi

########################################################################################################################
# 5 Remove the "(Unreleased)" prefix from the current version of the RELEASE_NOTES file (local)
########################################################################################################################

if ! sed -i '' "s/(Unreleased) Apache PLC4X $PROJECT_VERSION*/Apache PLC4X $RELEASE_VERSION/" ../RELEASE_NOTES; then
    echo "❌ Got non-0 exit code from updating RELEASE_NOTES, aborting."
    exit 1
fi

########################################################################################################################
# 6. Add a new section for the new version to the RELEASE_NOTES file (local)
########################################################################################################################

NEW_HEADER="==============================================================\n\
(Unreleased) Apache PLC4X $NEW_VERSION\n\
==============================================================\n\
\n\
New Features\n\
------------\n\
\n\
Incompatible changes\n\
--------------------\n\
\n\
Bug Fixes\n\
---------\n\
\n\
==============================================================\n\
"
echo NEW_VERSION
if ! sed -i '' "1s/.*/$NEW_HEADER/" ../RELEASE_NOTES; then
    echo "❌ Got non-0 exit code from adding a new header to RELEASE_NOTES, aborting."
    exit 1
fi

########################################################################################################################
# 7. Commit the change (local)
########################################################################################################################

if ! git add --all; then
    echo "❌ Got non-0 exit code from adding all changes files, aborting."
    exit 1
fi
if ! git commit -m "chore: prepared the RELEASE_NOTES for the next version."; then
    echo "❌ Got non-0 exit code from committing changes files, aborting."
    exit 1
fi

########################################################################################################################
# 8. Push the changes (local)
########################################################################################################################

if ! git push; then
    echo "❌ Got non-0 exit code from pushing changes, aborting."
    exit 1
fi

########################################################################################################################
# 9. Switch to the release branch (local)
########################################################################################################################

if ! git checkout "$BRANCH_NAME"; then
    echo "❌ Got non-0 exit code from switching branches to the release branch, aborting."
    exit 1
fi

echo "✅ Release branch creation complete. We have switched the local branch to the release branch. Please continue with 'release-2-prepare-release.sh' as soon as the release branch is ready for being released."