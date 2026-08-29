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

# Resolve the project directory from the location of this script, so that it does not matter which
# directory it is started from. This used to be "$(pwd)" with a "find .." further down, which
# deleted files from the PARENT of the checkout - and that is where the plc4x-build-tools and
# plc4x-extras checkouts usually live - whenever the script was not started from "tools".
DIRECTORY="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

########################################################################################################################
# 1. Work out what to roll the versions back to
########################################################################################################################

# Maven 4 prefixes even quiet output with "[INFO] [stdout] ", so take the last token of the
# last line rather than the whole output.
CURRENT_VERSION=$("$DIRECTORY/mvnw" -f "$DIRECTORY/pom.xml" -q -Dexec.executable=echo -Dexec.args="\${project.version}" --non-recursive exec:exec | tail -n 1 | awk '{print $NF}')

# The version to go back to cannot be derived reliably: after "release:branch" the poms hold the
# next minor version, after "release:prepare" the next bugfix version, and in between the release
# version itself. So take it as an argument, or ask.
PRE_RELEASE_VERSION="$1"
if [[ -z "$PRE_RELEASE_VERSION" ]]; then
    echo "The poms currently say: $CURRENT_VERSION"
    read -r -p "Which version should they be set back to? " PRE_RELEASE_VERSION
fi
if [[ -z "$PRE_RELEASE_VERSION" ]]; then
    echo "❌ No version given, aborting."
    exit 1
fi
if [[ ! "$PRE_RELEASE_VERSION" =~ -SNAPSHOT$ ]]; then
    echo "❌ '$PRE_RELEASE_VERSION' is not a SNAPSHOT version, aborting."
    echo "   The development version is what this rolls back to, for example 1.0.0-SNAPSHOT."
    exit 1
fi

RELEASE_VERSION=${PRE_RELEASE_VERSION%"-SNAPSHOT"}
IFS='.' read -ra VERSION_SEGMENTS <<< "$RELEASE_VERSION"
BRANCH_NAME="rel/${VERSION_SEGMENTS[0]}.${VERSION_SEGMENTS[1]}"
TAG_NAME="v$RELEASE_VERSION"

echo "Rolling back to:      $PRE_RELEASE_VERSION"
echo "Release branch:       $BRANCH_NAME"
echo "Release tag:          $TAG_NAME"

########################################################################################################################
# 2. Set the versions back (Docker container)
########################################################################################################################

if ! docker compose -f "$DIRECTORY/tools/docker-compose.yaml" run releaser \
        bash /ws/mvnw -e -P with-c,with-dotnet,with-go,with-java,with-python,update-generated-code \
        -Dmaven.repo.local=/ws/out/.repository versions:set -DprocessAllModules=true -DnewVersion="$PRE_RELEASE_VERSION"; then
    echo "❌ Got non-0 exit code from setting the versions back, aborting."
    exit 1
fi

########################################################################################################################
# 3. Delete the files the release plugin left behind (local)
########################################################################################################################

find "$DIRECTORY" -type f -name 'release.properties' -delete
find "$DIRECTORY" -type f -name 'pom.xml.versionsBackup' -delete
find "$DIRECTORY" -type f -name 'pom.xml.releaseBackup' -delete
echo "✅ Removed the left-over release.properties and pom backup files."

########################################################################################################################
# 4. Delete the release branch and the tag
########################################################################################################################

# These are separate questions because an attempt can fail before either of them was created, and
# because deleting them on the remote is not something to do by accident.
delete_ref() {
  local kind="$1" name="$2" local_cmd="$3" remote_cmd="$4"
  local exists_local=false exists_remote=false
  if [[ "$kind" == "branch" ]]; then
    git -C "$DIRECTORY" show-ref --verify --quiet "refs/heads/$name" && exists_local=true
  else
    git -C "$DIRECTORY" show-ref --verify --quiet "refs/tags/$name" && exists_local=true
  fi
  git -C "$DIRECTORY" ls-remote --exit-code origin "$name" > /dev/null 2>&1 && exists_remote=true

  if [[ "$exists_local" == false && "$exists_remote" == false ]]; then
    echo "✅ No $kind '$name' to delete."
    return 0
  fi
  echo
  echo "$kind '$name' exists: local=$exists_local remote=$exists_remote"
  read -r -p "Delete it? (yes/no) " yn
  if [[ "$yn" != "yes" ]]; then
    echo "✅ Keeping $kind '$name'."
    return 0
  fi
  if [[ "$exists_local" == true ]]; then
    eval "$local_cmd" || echo "⚠️  Could not delete the local $kind '$name'."
  fi
  if [[ "$exists_remote" == true ]]; then
    eval "$remote_cmd" || echo "⚠️  Could not delete the remote $kind '$name'."
  fi
  echo "✅ Deleted $kind '$name'."
}

delete_ref tag "$TAG_NAME" \
    "git -C \"$DIRECTORY\" tag -d \"$TAG_NAME\"" \
    "git -C \"$DIRECTORY\" push origin --delete \"$TAG_NAME\""
delete_ref branch "$BRANCH_NAME" \
    "git -C \"$DIRECTORY\" branch -D \"$BRANCH_NAME\"" \
    "git -C \"$DIRECTORY\" push origin --delete \"$BRANCH_NAME\""

########################################################################################################################
# 5. What this does not undo
########################################################################################################################

cat <<EOF

Not undone by this script - check these by hand:
  - Commits already pushed to 'develop' by 'release-1-create-branch.sh': the finalized
    RELEASE_NOTES, the new section for the next version and the version in
    website/asciidoc/antora.yml.
  - A release candidate staged in SVN under https://dist.apache.org/repos/dist/dev/plc4x/
    ("svn rm" it).
  - A Nexus staging repository opened by 'release-2-prepare-release.sh'
    (drop it at https://repository.apache.org).
  - The drafted emails in out/stage.
EOF
