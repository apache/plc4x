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

# BSD and GNU sed disagree about "-i": BSD wants a backup suffix as a separate argument, GNU
# treats that argument as the expression. Editing through a temporary file works on both, and
# leaves nothing behind if the edit fails - a stray backup file would otherwise trip the check
# for a clean working tree the next time one of these scripts runs.
#   $1 the file to edit, the rest are passed to sed unchanged
sed_in_place() {
  local file="$1"; shift
  if ! sed "$@" "$file" > "$file.sed.tmp"; then
      rm -f "$file.sed.tmp"
      return 1
  fi
  mv "$file.sed.tmp" "$file"
}

########################################################################################################################
# 0. Check if there are uncommitted changes as these would automatically be committed (local)
########################################################################################################################

if [[ $(git -C "$DIRECTORY" status --porcelain) ]]; then
  # Changes
  echo "❌ There are untracked files or changed files, aborting."
  exit 1
fi

########################################################################################################################
# 1. Get and calculate the current version (local)
########################################################################################################################

# Maven 4 prefixes even quiet output with "[INFO] [stdout] ", so take the last token of the
# last line rather than the whole output.
PROJECT_VERSION=$("$DIRECTORY/mvnw" -f "$DIRECTORY/pom.xml" -q -Dexec.executable=echo -Dexec.args="\${project.version}" --non-recursive exec:exec | tail -n 1 | awk '{print $NF}')
if [[ ! "$PROJECT_VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+(-SNAPSHOT)?$ ]]; then
    echo "❌ Could not read a usable project version, got '$PROJECT_VERSION'."
    echo "   Everything below derives the branch and tag names from it, so aborting."
    exit 1
fi
RELEASE_VERSION=${PROJECT_VERSION%"-SNAPSHOT"}
RELEASE_SHORT_VERSION=${RELEASE_VERSION%".0"}
BRANCH_NAME="rel/$RELEASE_SHORT_VERSION"
IFS='.' read -ra VERSION_SEGMENTS <<< "$RELEASE_VERSION"
NEW_VERSION="${VERSION_SEGMENTS[0]}.$((VERSION_SEGMENTS[1] + 1)).0-SNAPSHOT"
NEW_DOCS_VERSION=${NEW_VERSION%"-SNAPSHOT"}
ANTORA_DESCRIPTOR="$DIRECTORY/website/asciidoc/antora.yml"
echo "Current Version: '$PROJECT_VERSION'"
echo "Release Version: '$RELEASE_VERSION'"
echo "Release Branch Name: '$BRANCH_NAME'"
echo "New develop Version: '$NEW_VERSION'"
echo "New develop Docs Version: '$NEW_DOCS_VERSION'"

########################################################################################################################
# 3. Ask if the RELEASE_NOTES have been filled out at all (local)
########################################################################################################################

read -r -p "Have the RELEASE_NOTES been updated for this version? (yes/no) " yn
case $yn in
	yes ) echo continuing with the process;;
	no ) echo Please update the RELEASE_NOTES first;
		exit 1;;
	* ) echo invalid response;
		exit 1;;
esac

########################################################################################################################
# 4 Remove the "(Unreleased)" prefix from the current version of the RELEASE_NOTES file (local)
########################################################################################################################

if ! sed_in_place "$DIRECTORY/RELEASE_NOTES" "s/(Unreleased) Apache PLC4X $PROJECT_VERSION*/Apache PLC4X $RELEASE_VERSION/"; then
    echo "❌ Got non-0 exit code from updating RELEASE_NOTES, aborting."
    exit 1
fi

# Commit this change to git.
git -C "$DIRECTORY" add --all
git -C "$DIRECTORY" commit -m "chore: updated generated code"

########################################################################################################################
# 5. Do a simple maven branch command with pushChanges=false
########################################################################################################################

# Attempt to read user.name and user.email (local first, then global)
GIT_USER_NAME=$(git config user.name || git config --global user.name)
GIT_USER_EMAIL=$(git config user.email || git config --global user.email)

# Check if either is still unset
if [[ -z "$GIT_USER_NAME" || -z "$GIT_USER_EMAIL" ]]; then
  echo "❌ Git user.name and/or user.email not configured."
  echo
  echo "Please run one of the following commands:"
  echo "  git config --global user.name \"Your Name\""
  echo "  git config --global user.email \"you@example.com\""
  echo
  echo "Or configure them just for this repository:"
  echo "  git config user.name \"Your Name\""
  echo "  git config user.email \"you@example.com\""
  exit 1
fi

if ! docker compose -f "$DIRECTORY/tools/docker-compose.yaml" run releaser \
        bash -c "git config user.name \"$GIT_USER_NAME\" && \
           git config user.email \"$GIT_USER_EMAIL\" && \
           /ws/mvnw -e -P with-c,with-dotnet,with-go,with-java,with-python,enable-all-checks,update-generated-code -Dmaven.repo.local=/ws/out/.repository release:branch -DautoVersionSubmodules=true -DpushChanges=false -DdevelopmentVersion='$NEW_VERSION' -DbranchName='$BRANCH_NAME'"; then
    echo "❌ Got non-0 exit code from docker compose, aborting."
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
==============================================================\
"
echo NEW_VERSION
if ! sed_in_place "$DIRECTORY/RELEASE_NOTES" "1s/.*/$NEW_HEADER/"; then
    echo "❌ Got non-0 exit code from adding a new header to RELEASE_NOTES, aborting."
    exit 1
fi

########################################################################################################################
# 7. Point the documentation of "develop" at the next version (local)
########################################################################################################################

# Every branch names the concrete version it documents in its own Antora descriptor. The
# "urls.latest_version_segment" / "urls.latest_prerelease_version_segment" keys in
# website/antora-playbook.yml then map the newest of them onto the "latest" and
# "pre-release" URL segments, so no branch ever has to know whether it is the current
# release. That only works if every branch keeps its version unique and up to date.
#
# Sets "version:" (and optionally "prerelease:") in an antora.yml. Arguments:
#   $1 the antora.yml to edit, $2 the version to set, $3 the prerelease value (optional)
update_antora_version() {
  local descriptor="$1" version="$2" prerelease="$3"

  if [[ ! -f "$descriptor" ]]; then
    echo "❌ Antora descriptor '$descriptor' not found, aborting."
    exit 1
  fi
  if ! grep -qE "^version:" "$descriptor"; then
    echo "❌ No 'version:' key found in '$descriptor', aborting."
    exit 1
  fi
  if ! sed_in_place "$descriptor" -E "s|^version:.*|version: '$version'|"; then
    echo "❌ Got non-0 exit code from updating the version in '$descriptor', aborting."
    exit 1
  fi
  if [[ -n "$prerelease" ]]; then
    if ! grep -qE "^prerelease:" "$descriptor"; then
      echo "❌ No 'prerelease:' key found in '$descriptor', aborting."
      exit 1
    fi
    if ! sed_in_place "$descriptor" -E "s|^prerelease:.*|prerelease: $prerelease|"; then
      echo "❌ Got non-0 exit code from updating the prerelease flag in '$descriptor', aborting."
      exit 1
    fi
  fi
  echo "✅ '$descriptor' now documents version '$version'."
}

# "develop" stays a prerelease, it just moves on to the next version.
update_antora_version "$ANTORA_DESCRIPTOR" "$NEW_DOCS_VERSION" "True"

########################################################################################################################
# 8. Commit the change (local)
########################################################################################################################

if ! git -C "$DIRECTORY" add --all; then
    echo "❌ Got non-0 exit code from adding all changes files, aborting."
    exit 1
fi
if ! git -C "$DIRECTORY" commit -m "chore: prepared the RELEASE_NOTES and the documentation version for the next version."; then
    echo "❌ Got non-0 exit code from committing changes files, aborting."
    exit 1
fi

########################################################################################################################
# 9. Push the changes (local)
########################################################################################################################

if ! git -C "$DIRECTORY" push; then
    echo "❌ Got non-0 exit code from pushing changes, aborting."
    exit 1
fi

########################################################################################################################
# 10. Switch to the release branch (local)
########################################################################################################################

if ! git -C "$DIRECTORY" checkout "$BRANCH_NAME"; then
    echo "❌ Got non-0 exit code from switching branches to the release branch, aborting."
    exit 1
fi

# Make sure the release branch is also pushed to the remote.
if ! git -C "$DIRECTORY" push --set-upstream origin "$BRANCH_NAME"; then
    echo "❌ Got non-0 exit code from pushing changes, aborting."
    exit 1
fi

########################################################################################################################
# 11. Point the documentation of the release branch at the version it is going to release (local)
########################################################################################################################

# The release branch documents the version it is going to release, but stays flagged as a
# prerelease: nothing has been voted on yet. Nothing of it is published while the vote runs -
# Antora only reads the branches listed in website/antora-playbook.yml, and
# 'release-3-finish-release.sh' adds it there after a successful vote. The flag matters at that
# moment: Antora skips prereleases when it picks the latest version of a component, so the branch
# would appear under its own version number rather than as ".../plc4x/latest/...".
# 'release-3-finish-release.sh' clears the flag in the same run, and that is what makes this
# branch the published release.
update_antora_version "$ANTORA_DESCRIPTOR" "$RELEASE_VERSION" "True"

# Usually there is nothing to commit here: the branch was cut from "develop" before step 6b moved
# it on, so it already documents the version being released. Only a "develop" that was out of sync
# leaves an actual change behind.
if [[ $(git -C "$DIRECTORY" status --porcelain) ]]; then
    if ! git -C "$DIRECTORY" add "$ANTORA_DESCRIPTOR"; then
        echo "❌ Got non-0 exit code from adding the Antora descriptor, aborting."
        exit 1
    fi
    if ! git -C "$DIRECTORY" commit -m "chore: set the documentation version of the release branch to $RELEASE_VERSION (still a prerelease)."; then
        echo "❌ Got non-0 exit code from committing the Antora descriptor, aborting."
        exit 1
    fi
    if ! git -C "$DIRECTORY" push; then
        echo "❌ Got non-0 exit code from pushing the Antora descriptor, aborting."
        exit 1
    fi
else
    echo "✅ The release branch already documents $RELEASE_VERSION."
fi

echo "✅ Release branch creation complete. We have switched the local branch to the release branch. Please continue with 'release-2-prepare-release.sh' as soon as the release branch is ready for being released."