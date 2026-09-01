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

# Everything below writes to the release branch, so make sure that is where we are before
# anything is changed - running this on "develop" by accident would rewrite its documentation
# version.
RELEASE_BRANCH=$(git -C "$DIRECTORY" rev-parse --abbrev-ref HEAD)
if [[ ! "$RELEASE_BRANCH" =~ ^rel/ ]]; then
    echo "❌ Expected to be on a 'rel/...' branch, but found '$RELEASE_BRANCH', aborting."
    echo "   Check out the release branch of the release that just passed its vote."
    exit 1
fi

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
# Helpers for maintaining the download page
########################################################################################################################

# The download page exists on every branch that publishes documentation, and each copy has to list
# the release, because ".../plc4x/latest/users/download.html" is served from the release branch
# while ".../plc4x/pre-release/..." is served from "develop".
#
# Prints the section of the RELEASE_NOTES belonging to one version, converted to AsciiDoc.
#   $1 the RELEASE_NOTES file, $2 the version
release_notes_to_asciidoc() {
  awk -v ver="$2" '
    /^=+$/ {
      if (state == 2) { exit }        # banner of the next version, we are done
      if (state == 1) { state = 2 }   # closing banner of our own title
      else { pending = 1 }
      next
    }
    {
      if (pending) {
        pending = 0
        title = $0
        sub(/^\(Unreleased\) /, "", title)
        if (title == "Apache PLC4X " ver) { state = 1 }
        next
      }
      if (state == 2) { print }
    }
  ' "$1" | awk '
    # Each bullet ends up on a single line. The RELEASE_NOTES wrap them and indent the
    # continuations, and neither keeping nor stripping that indentation is safe: a continuation
    # that happens to start with something like "18446744073709551615." would be read as an
    # ordered list item by AsciiDoc. Joining the lines avoids the question and renders the same.
    function flush() { if (buffer != "") { print buffer; buffer = "" } }
    { lines[NR] = $0 }
    END {
      last = NR
      while (last > 0 && lines[last] ~ /^[ \t]*$/) { last-- }   # trim trailing blank lines
      buffer = ""
      for (i = 1; i <= last; i++) {
        l = lines[i]
        # "New Features" underlined with dashes becomes a level 4 heading
        if (l ~ /^[A-Za-z]/ && lines[i+1] ~ /^-+$/) { flush(); print "==== " l; i++; started = 1; continue }
        if (l ~ /^[ \t]*$/) { if (!started) { continue }   # skip the blank lines before the first heading
                               flush(); print ""; continue }
        started = 1
        sub(/^[ \t]+/, "", l)
        if (l ~ /^[-*] /) {                       # a new bullet
          flush()
          sub(/^- PLC4X-[0-9]+ /, "- ", l)        # older notes carried the issue id in front
          buffer = l
        } else if (buffer != "") {                # a wrapped continuation of the current bullet
          buffer = buffer " " l
        } else {                                  # ordinary prose
          print l
        }
      }
      flush()
    }
  '
}

# Adds a release to a download page: the entry that was current is moved down to "Previous
# Releases" with its source-release link pointed at the archive, and the new one takes its place.
#   $1 the download.adoc to edit, $2 the version, $3 the RELEASE_NOTES to take the notes from
update_download_page() {
  local page="$1" version="$2" notes="$3"
  local anchor="release-${version//./_}"

  if [[ ! -f "$page" ]]; then
      echo "❌ Download page '$page' not found, aborting."
      exit 1
  fi
  if grep -q "^\[#$anchor\]$" "$page"; then
      echo "✅ $version is already listed on '$page'."
      return 0
  fi

  local current_line previous_line
  current_line=$(grep -n '^== Current Releases$' "$page" | head -1 | cut -d: -f1)
  previous_line=$(grep -n '^== Previous Releases$' "$page" | head -1 | cut -d: -f1)
  if [[ -z "$current_line" || -z "$previous_line" || "$previous_line" -le "$current_line" ]]; then
      echo "❌ '$page' has no usable 'Current Releases'/'Previous Releases' structure, aborting."
      exit 1
  fi

  # This function always makes $version the current release and pushes what was there into
  # "Previous Releases". That is only right when $version really is the newest release: a bugfix
  # off an old branch (0.13.2 while 1.0.0 is current) would otherwise archive 1.0.0. Where that
  # happens the page needs a human decision about how to present two maintained lines.
  local current_version
  current_version=$(sed -n "$((current_line + 1)),$((previous_line - 1))p" "$page" \
      | grep -oE '^\[#release-[0-9_]+\]$' | head -1 | sed -E 's|^\[#release-||; s|\]$||; s|_|.|g')
  if [[ -n "$current_version" && "$current_version" != "$version" ]] \
          && [[ "$(printf '%s\n%s\n' "$version" "$current_version" | sort -V | tail -1)" != "$version" ]]; then
      echo "❌ '$page' lists $current_version as the current release, which is newer than $version, aborting."
      echo "   This looks like a bugfix release off an older branch. Add the entry by hand so that"
      echo "   $current_version stays the current release."
      exit 1
  fi

  local notes_body
  notes_body=$(release_notes_to_asciidoc "$notes" "$version")
  if [[ -z "$notes_body" ]]; then
      echo "❌ Found no RELEASE_NOTES section for $version in '$notes', aborting."
      exit 1
  fi

  {
    head -n "$current_line" "$page"
    echo
    echo "[#$anchor]"
    echo "=== $version Official https://www.apache.org/dyn/closer.lua/plc4x/$version/apache-plc4x-$version-source-release.zip[source release] [ https://downloads.apache.org/plc4x/$version/apache-plc4x-$version-source-release.zip.sha512[SHA512] ] [ https://downloads.apache.org/plc4x/$version/apache-plc4x-$version-source-release.zip.asc[ASC] ]"
    echo
    echo "$notes_body"
    echo
    echo "== Previous Releases"
    # What used to be current, with its source-release link moved to the archive.
    sed -n "$((current_line + 1)),$((previous_line - 1))p" "$page" \
        | sed 's|https://www\.apache\.org/dyn/closer\.lua/plc4x/|https://archive.apache.org/dist/plc4x/|g'
    tail -n +"$((previous_line + 1))" "$page"
  } > "$page.tmp"

  if ! grep -q "^\[#$anchor\]$" "$page.tmp"; then
      echo "❌ Could not add $version to '$page', aborting."
      rm -f "$page.tmp"
      exit 1
  fi
  mv "$page.tmp" "$page"
  echo "✅ $version added to '$page'."
}

########################################################################################################################
# 1. Publish the documentation of this branch as the current release (local)
########################################################################################################################

# Up to this point the release branch documented the version it was going to release, but was
# still flagged as a prerelease, so ".../plc4x/latest/..." kept serving the previous release
# while the vote was running (see 'release-1-create-branch.sh'). The vote has passed, so the
# flag can be cleared. Antora then picks this branch as the latest version of the component
# and publishes it as ".../plc4x/latest/...", which at the same time demotes the previous
# release branch to its own version number. Nothing has to be changed on that older branch.
#
# The version is taken from the tag 'release-2-prepare-release.sh' created, so that this also
# does the right thing for a bugfix release, where 'release-1-create-branch.sh' never ran.

ANTORA_DESCRIPTOR="$DIRECTORY/website/asciidoc/antora.yml"

RELEASE_TAG=$(git -C "$DIRECTORY" describe --tags --abbrev=0 --match "v*" 2>/dev/null)
if [[ -z "$RELEASE_TAG" ]]; then
    echo "❌ Could not determine the release tag of this branch, aborting."
    echo "   Are you on the release branch, and did 'release-2-prepare-release.sh' run?"
    exit 1
fi
RELEASED_VERSION=${RELEASE_TAG#v}

# "git describe" finds the nearest tag reachable from HEAD, which on a maintenance branch that has
# not been prepared yet is the PREVIOUS release. Check the tag actually belongs to this branch.
TAG_LINE=$(echo "$RELEASED_VERSION" | cut -d. -f1,2)
if [[ "$RELEASE_BRANCH" != "rel/$TAG_LINE" ]]; then
    echo "❌ The most recent tag is '$RELEASE_TAG', which does not belong to branch '$RELEASE_BRANCH', aborting."
    echo "   Expected a tag on the rel/$TAG_LINE line. Has 'release-2-prepare-release.sh' run on this branch?"
    exit 1
fi

read -r -p "Did the vote for $RELEASED_VERSION pass and is it being released? (yes/no) " yn
case $yn in
	yes ) echo "Publishing the documentation of $RELEASED_VERSION as the current release.";;
	no ) echo "Not publishing anything. Re-run this script once the vote has passed.";
		exit 1;;
	* ) echo "invalid response";
		exit 1;;
esac

if [[ ! -f "$ANTORA_DESCRIPTOR" ]]; then
    echo "❌ Antora descriptor '$ANTORA_DESCRIPTOR' not found, aborting."
    exit 1
fi
# Check the RELEASE_NOTES before writing anything: the download page needs them further down, and
# aborting halfway would leave the descriptor modified and the next run blocked by the check for a
# clean working tree at the top of this script.
if [[ -z "$(release_notes_to_asciidoc "$DIRECTORY/RELEASE_NOTES" "$RELEASED_VERSION")" ]]; then
    echo "❌ Found no RELEASE_NOTES section for $RELEASED_VERSION, aborting."
    echo "   Expected a block headed 'Apache PLC4X $RELEASED_VERSION' - if it still says"
    echo "   '(Unreleased) ... -SNAPSHOT', 'release-1-create-branch.sh' never finalized it."
    exit 1
fi
if ! grep -qE "^version:" "$ANTORA_DESCRIPTOR" || ! grep -qE "^prerelease:" "$ANTORA_DESCRIPTOR"; then
    echo "❌ '$ANTORA_DESCRIPTOR' has no 'version:' and/or 'prerelease:' key, aborting."
    exit 1
fi
if ! sed_in_place "$ANTORA_DESCRIPTOR" -E "s|^version:.*|version: '$RELEASED_VERSION'|"; then
    echo "❌ Got non-0 exit code from updating the version in the Antora descriptor, aborting."
    exit 1
fi
if ! sed_in_place "$ANTORA_DESCRIPTOR" -E "s|^prerelease:.*|prerelease: False|"; then
    echo "❌ Got non-0 exit code from clearing the prerelease flag, aborting."
    exit 1
fi
echo "✅ '$ANTORA_DESCRIPTOR' now publishes $RELEASED_VERSION as the current release."

# The copy of the download page on this branch is what ".../plc4x/latest/users/download.html"
# will be served from, so it has to list the release too.
update_download_page "$DIRECTORY/website/asciidoc/modules/users/pages/download.adoc" \
    "$RELEASED_VERSION" "$DIRECTORY/RELEASE_NOTES"

if [[ $(git -C "$DIRECTORY" status --porcelain) ]]; then
    if ! git -C "$DIRECTORY" add "$ANTORA_DESCRIPTOR" \
            "$DIRECTORY/website/asciidoc/modules/users/pages/download.adoc"; then
        echo "❌ Got non-0 exit code from adding the Antora descriptor, aborting."
        exit 1
    fi
    if ! git -C "$DIRECTORY" commit -m "chore: published the documentation of $RELEASED_VERSION as the current release."; then
        echo "❌ Got non-0 exit code from committing the Antora descriptor, aborting."
        exit 1
    fi
    if ! git -C "$DIRECTORY" push; then
        echo "❌ Got non-0 exit code from pushing the Antora descriptor, aborting."
        exit 1
    fi
else
    echo "✅ The Antora descriptor was already up to date."
fi

########################################################################################################################
# 2. Update the website configuration on "develop"
########################################################################################################################

# Two things have to happen on "develop" for this release to show up on the website:
#
#   - the release branch has to be listed in the "content.sources" list of
#     website/antora-playbook.yml, because Antora only reads the branches listed there - without
#     it the descriptor changed above has no effect at all
#   - the released version has to be added to website/resources/plc4x-doap.rdf, which is what
#     Apache's tooling reads to keep track of the project's release activity
#
# Both files live on "develop", which is not the branch we are on, so the edits are made in a
# throw-away worktree and pushed straight to "develop" - the release branch stays checked out in
# the working copy the release manager is using.

if ! git -C "$DIRECTORY" fetch origin develop; then
    echo "❌ Got non-0 exit code from fetching 'develop', aborting."
    exit 1
fi

DEVELOP_WORKTREE=$(mktemp -d)
cleanup_develop_worktree() {
  git -C "$DIRECTORY" worktree remove --force "$DEVELOP_WORKTREE" >/dev/null 2>&1
  rm -rf "$DEVELOP_WORKTREE"
  # If the remove failed the directory is gone but the .git/worktrees entry is not, so prune it.
  git -C "$DIRECTORY" worktree prune >/dev/null 2>&1
}
trap cleanup_develop_worktree EXIT

# A detached worktree, so this works even if "develop" happens to be checked out somewhere else.
rmdir "$DEVELOP_WORKTREE"
if ! git -C "$DIRECTORY" worktree add --detach "$DEVELOP_WORKTREE" origin/develop; then
    echo "❌ Got non-0 exit code from creating a worktree for 'develop', aborting."
    exit 1
fi

PLAYBOOK="$DEVELOP_WORKTREE/website/antora-playbook.yml"
DOAP="$DEVELOP_WORKTREE/website/resources/plc4x-doap.rdf"
for f in "$PLAYBOOK" "$DOAP"; do
  if [[ ! -f "$f" ]]; then
      echo "❌ '$f' not found, aborting."
      exit 1
  fi
done

CHANGES=()

########################################################################################################################
# 2a. Add the release branch to the Antora playbook
########################################################################################################################

if grep -q "'$RELEASE_BRANCH'" "$PLAYBOOK"; then
    echo "✅ '$RELEASE_BRANCH' is already listed in the Antora playbook."
else
    # Insert directly after the first "start_path", which ends the "develop" source, so the
    # release branches stay listed newest first.
    if ! awk -v branch="$RELEASE_BRANCH" '
        BEGIN { done = 0 }
        { print }
        !done && /^[[:space:]]*start_path: website\/asciidoc[[:space:]]*$/ {
            print "  - url: https://github.com/apache/plc4x.git"
            print "    branches: [\x27" branch "\x27]"
            print "    start_path: website/asciidoc"
            done = 1
        }
    ' "$PLAYBOOK" > "$PLAYBOOK.tmp"; then
        echo "❌ Got non-0 exit code from adding '$RELEASE_BRANCH' to the playbook, aborting."
        exit 1
    fi
    if ! grep -q "'$RELEASE_BRANCH'" "$PLAYBOOK.tmp"; then
        echo "❌ Could not add '$RELEASE_BRANCH' to the Antora playbook, aborting."
        echo "   The 'content.sources' list does not look the way this script expects."
        rm -f "$PLAYBOOK.tmp"
        exit 1
    fi
    mv "$PLAYBOOK.tmp" "$PLAYBOOK"
    CHANGES+=("website/antora-playbook.yml")
    echo "✅ '$RELEASE_BRANCH' added to the Antora playbook."
fi

########################################################################################################################
# 2b. Add the released version to the DOAP file
########################################################################################################################

if grep -q "<revision>$RELEASED_VERSION</revision>" "$DOAP"; then
    echo "✅ $RELEASED_VERSION is already listed in the DOAP file."
else
    RELEASE_DATE=$(date +%Y-%m-%d)
    # The newest release is the one named "Latest", all the others are "Legacy", so the entry
    # that held the title up to now has to hand it over before the new one is inserted at the top.
    if ! sed_in_place "$DOAP" -E "s|<name>Latest ([0-9][^<]*) release</name>|<name>Legacy \1 release</name>|"; then
        echo "❌ Got non-0 exit code from demoting the previous release in the DOAP file, aborting."
        exit 1
    fi
    if ! awk -v version="$RELEASED_VERSION" -v released="$RELEASE_DATE" '
        BEGIN { done = 0 }
        !done && /^[[:space:]]*<release>[[:space:]]*$/ {
            print "        <release>"
            print "            <Version>"
            print "                <name>Latest " version " release</name>"
            print "                <created>" released "</created>"
            print "                <revision>" version "</revision>"
            print "            </Version>"
            print "        </release>"
            done = 1
        }
        { print }
    ' "$DOAP" > "$DOAP.tmp"; then
        echo "❌ Got non-0 exit code from adding $RELEASED_VERSION to the DOAP file, aborting."
        exit 1
    fi
    if ! grep -q "<revision>$RELEASED_VERSION</revision>" "$DOAP.tmp"; then
        echo "❌ Could not add $RELEASED_VERSION to the DOAP file, aborting."
        echo "   The list of <release> entries does not look the way this script expects."
        rm -f "$DOAP.tmp"
        exit 1
    fi
    mv "$DOAP.tmp" "$DOAP"
    CHANGES+=("website/resources/plc4x-doap.rdf")
    echo "✅ $RELEASED_VERSION added to the DOAP file, released on $RELEASE_DATE."
fi

########################################################################################################################
# 2c. Add the release to the download page
########################################################################################################################

# "develop" has its own copy of the download page, published as ".../plc4x/pre-release/...". The
# notes are taken from the RELEASE_NOTES of the release branch, because the ones on "develop" have
# already moved on to the next version.

DOWNLOAD_PAGE="$DEVELOP_WORKTREE/website/asciidoc/modules/users/pages/download.adoc"
DOWNLOAD_PAGE_BEFORE=$(cksum < "$DOWNLOAD_PAGE")
update_download_page "$DOWNLOAD_PAGE" "$RELEASED_VERSION" "$DIRECTORY/RELEASE_NOTES"
if [[ "$(cksum < "$DOWNLOAD_PAGE")" != "$DOWNLOAD_PAGE_BEFORE" ]]; then
    CHANGES+=("website/asciidoc/modules/users/pages/download.adoc")
fi

########################################################################################################################
# 2d. Commit and push whatever changed
########################################################################################################################

if [[ ${#CHANGES[@]} -eq 0 ]]; then
    echo "✅ The website configuration on 'develop' was already up to date."
else
    if ! git -C "$DEVELOP_WORKTREE" add "${CHANGES[@]}"; then
        echo "❌ Got non-0 exit code from adding the changed files, aborting."
        exit 1
    fi
    if ! git -C "$DEVELOP_WORKTREE" commit -m "chore: added $RELEASED_VERSION to the website configuration."; then
        echo "❌ Got non-0 exit code from committing the changed files, aborting."
        exit 1
    fi
    if ! git -C "$DEVELOP_WORKTREE" push origin HEAD:develop; then
        echo "❌ Got non-0 exit code from pushing to 'develop', aborting."
        echo "   Someone probably pushed to 'develop' in the meantime - re-run this script."
        exit 1
    fi
    echo "✅ Pushed to 'develop': ${CHANGES[*]}"
    echo "   Your local 'develop' is now behind by that commit - remember to pull it."
fi

cleanup_develop_worktree
trap - EXIT

# NOTE: Old releases are not removed from the Antora playbook automatically - decide for yourself
#   how many of them the website should keep building.

########################################################################################################################
# 3. Draft the [RESULT] and [ANNOUNCE] emails
########################################################################################################################

# Like the [VOTE] and [DISCUSS] emails drafted by 'release-2-prepare-release.sh', these are
# written to files for you to read and send yourself - nothing is sent from here. They are drafted
# before anything is published, so that a run that stops at the confirmation below still leaves you
# with the [RESULT] email needed to close the vote.


if ! command -v svn > /dev/null; then
    echo "❌ 'svn' not found, but it is needed to publish the release artifacts, aborting."
    exit 1
fi

# Everything below reads SVN to decide what to publish, so make sure we can actually see it. A
# listing that fails because of the network, a proxy or missing credentials looks exactly like
# "nothing is staged", and acting on that would mean publishing the wrong thing.
if ! svn ls "$DIST_DEV/" > /dev/null 2>&1; then
    echo "❌ Could not list $DIST_DEV, aborting."
    echo "   Fix the connection (or your SVN credentials) and run this again."
    exit 1
fi
if ! svn ls "$DIST_RELEASE/" > /dev/null 2>&1; then
    echo "❌ Could not list $DIST_RELEASE, aborting."
    exit 1
fi

# Which release candidate won the vote: the highest numbered one that is still staged. Empty means
# there is genuinely nothing staged any more - the release was published and cleaned up already.
RELEASE_CANDIDATE=$(svn ls "$DIST_DEV/$RELEASED_VERSION/" 2>/dev/null \
    | grep -E '^rc[0-9]+/$' | sed 's|/$||' | sort -V | tail -1)
# Only ever used for the subject of the [RESULT] mail, never to decide what gets published.
RELEASE_CANDIDATE_LABEL=${RELEASE_CANDIDATE:-rc1}

STAGE_DIR="$DIRECTORY/out/stage"
mkdir -p "$STAGE_DIR"

echo
echo "The [RESULT] email needs the vote count."
read -r -p "Number of +1 votes by PMC members: " VOTES_PMC
read -r -p "Number of +1 votes by non-PMC members: " VOTES_NON_PMC
VOTES_PMC=${VOTES_PMC:-0}
VOTES_NON_PMC=${VOTES_NON_PMC:-0}

cat > "$STAGE_DIR/result-email.eml" <<EOF
To: dev@plc4x.apache.org
Subject: [RESULT] [VOTE] Apache PLC4X $RELEASED_VERSION $(echo "$RELEASE_CANDIDATE_LABEL" | tr '[:lower:]' '[:upper:]')
Content-Type: text/plain; charset=UTF-8

So, the vote passes with $VOTES_PMC +1 votes by PMC members and $VOTES_NON_PMC +1 votes by non PMC members.

Thanks to everyone who took the time to validate the release.

I will continue with the release process and will announce the release as soon as the
artifacts have made it to the mirrors.
EOF
echo "✅ Result email generated to $STAGE_DIR/result-email.eml"

cat > "$STAGE_DIR/announce-email.eml" <<EOF
To: announce@apache.org
Cc: dev@plc4x.apache.org
Subject: [ANNOUNCE] Apache PLC4X $RELEASED_VERSION released
Content-Type: text/plain; charset=UTF-8

The Apache PLC4X team is pleased to announce the release of Apache PLC4X $RELEASED_VERSION

PLC4X is a set of libraries for communicating with industrial programmable
logic controllers (PLCs) using a variety of protocols but with a shared API.

The current release contains drivers able to communicate with industrial PLCs using one of the following protocols:

  *   AB-ETH
  *   Beckhoff ADS
  *   CanOpen
  *   EtherNet/IP / EIP
  *   Firmata
  *   KNXNet/IP
  *   Modbus
  *   OPC UA
  *   Siemens S7 (0x32)

Beyond that we also provide integration modules for the following Apache projects and frameworks:

  *   Apache Calcite
  *   Apache Camel
  *   Apache Kafka (Kafka Connect)
  *   Apache NiFi
  *   Logstash

It also provides an \`OPC UA Server\` which can act as a bridge between legacy systems and OPC UA.

Visit the Apache PLC4X website [1] for general information or
the downloads page [2] for release notes and download information.

Regards,
The Apache PLC4X team

[1] https://plc4x.apache.org
[2] https://plc4x.apache.org/plc4x/latest/users/download.html
EOF
echo "✅ Announce email generated to $STAGE_DIR/announce-email.eml"
echo "   Send it from your @apache.org address, and only after the mirrors have had 24 hours."
echo "   Check the list of drivers and integrations in it - it is not generated from anything."

########################################################################################################################
# 4. Publish the release artifacts
########################################################################################################################

# This is the point of no return: moving the release candidate inside the Apache SVN makes it
# available to the mirrors, and releasing the Nexus staging repository pushes the Maven artifacts
# to the Apache release repository, from where they are synced to Maven Central. Neither can be
# undone, which is also why the announcement has to wait 24 hours after this.

########################################################################################################################
# 4a. Work out which Nexus staging repository to release
########################################################################################################################

# The staging repository id was written by 'release-2-prepare-release.sh'. That file lives in
# "out", which does not survive 'release-0-update-generated-code.sh' and is not there at all if
# this script runs on a different machine - so fall back to asking.
DEPLOY_PROPS="$DIRECTORY/out/.local-artifacts-dir/$STAGING_PROFILE_ID.properties"
STAGING_REPO_ID=""
if [[ -f "$DEPLOY_PROPS" ]]; then
    STAGING_REPO_ID=$(grep stagingRepository.id "$DEPLOY_PROPS" | cut -d= -f2 | tr -d ' \r')
fi
if [[ -z "$STAGING_REPO_ID" ]]; then
    echo "Could not determine the Nexus staging repository from '$DEPLOY_PROPS'."
    echo "It is the last part of the URL in [1] of the VOTE email, for example 'orgapacheplc4x-1234'."
    read -r -p "Nexus staging repository id: " STAGING_REPO_ID
fi
if [[ -z "$STAGING_REPO_ID" ]]; then
    echo "❌ No Nexus staging repository id given, aborting."
    exit 1
fi

########################################################################################################################
# 4b. Move the release candidate to the release part of SVN
########################################################################################################################

SVN_ALREADY_PUBLISHED=false
if svn ls "$DIST_RELEASE/$RELEASED_VERSION/" > /dev/null 2>&1; then
    SVN_ALREADY_PUBLISHED=true
    echo "✅ $RELEASED_VERSION is already published under $DIST_RELEASE."
else
    if [[ -z "$RELEASE_CANDIDATE" ]]; then
        echo "❌ Found no release candidate under $DIST_DEV/$RELEASED_VERSION/, aborting."
        echo "   Did 'release-2-prepare-release.sh' stage this release?"
        exit 1
    fi

    echo
    echo "About to publish Apache PLC4X $RELEASED_VERSION. This cannot be undone:"
    echo "  - SVN:   $DIST_DEV/$RELEASED_VERSION/$RELEASE_CANDIDATE"
    echo "        -> $DIST_RELEASE/$RELEASED_VERSION"
    echo "  - Nexus: releasing staging repository '$STAGING_REPO_ID' to the Apache release repo,"
    echo "           from where it is synced to Maven Central."
    read -r -p "Publish? (yes/no) " yn
    case $yn in
        yes ) echo "Publishing.";;
        no ) echo "Nothing published. The documentation changes above are already done.";
            exit 1;;
        * ) echo "invalid response";
            exit 1;;
    esac

    if ! svn move -m "Release Apache PLC4X $RELEASED_VERSION" \
            "$DIST_DEV/$RELEASED_VERSION/$RELEASE_CANDIDATE" \
            "$DIST_RELEASE/$RELEASED_VERSION"; then
        echo "❌ Got non-0 exit code from moving the release candidate in SVN, aborting."
        exit 1
    fi
    echo "✅ $RELEASE_CANDIDATE published as $DIST_RELEASE/$RELEASED_VERSION."
    echo "   The mirrors need up to 24 hours to pick this up - do not announce before that."
fi

########################################################################################################################
# 4c. Release the Nexus staging repository
########################################################################################################################

# Unlike the steps above there is no cheap way to ask Nexus whether this repository has already
# been released, so if SVN says the release is already out, this is most likely a re-run and
# releasing again would just fail. Ask instead of guessing - it is also the way to recover from a
# run where the SVN move succeeded and this step did not.
if [[ "$SVN_ALREADY_PUBLISHED" == "true" ]]; then
    echo
    echo "$RELEASED_VERSION was already published in SVN, so the Nexus staging repository"
    echo "'$STAGING_REPO_ID' has probably been released already too."
    read -r -p "Try to release it anyway? (yes/no) " yn
    if [[ "$yn" != "yes" ]]; then
        echo "✅ Leaving the Nexus staging repository alone."
        STAGING_REPO_ID=""
    fi
fi

if [[ -n "$STAGING_REPO_ID" ]]; then
if ! "$DIRECTORY/mvnw" -f "$DIRECTORY/tools/stage.pom" nexus-staging:rc-release \
        -DstagingRepositoryId="$STAGING_REPO_ID" -DstagingProfileId=$STAGING_PROFILE_ID; then
    echo "❌ Got non-0 exit code from releasing the Nexus staging repository."
    echo "   If it was already released, this is expected and can be ignored - check"
    echo "   $NEXUS_URL under 'Staging Repositories'."
    exit 1
fi
echo "✅ Nexus staging repository '$STAGING_REPO_ID' released."
fi

########################################################################################################################
# 5. Clean up what the release left behind
########################################################################################################################

########################################################################################################################
# 5a. Remove the release candidates that did not make it from SVN
########################################################################################################################

# Only the winning release candidate was moved to the release part of SVN, so every earlier
# attempt is still sitting in the dev part, together with the now pointless version directory.
# None of it is needed once the release is out.

if svn ls "$DIST_RELEASE/$RELEASED_VERSION/" > /dev/null 2>&1 \
        && svn ls "$DIST_DEV/$RELEASED_VERSION/" > /dev/null 2>&1; then
    echo
    echo "These leftovers of $RELEASED_VERSION are still staged under $DIST_DEV/$RELEASED_VERSION:"
    svn ls -R "$DIST_DEV/$RELEASED_VERSION/" | sed 's|^|  |'
    read -r -p "Remove them? (yes/no) " yn
    if [[ "$yn" == "yes" ]]; then
        if ! svn rm -m "Removed the staged release candidates of PLC4X $RELEASED_VERSION" \
                "$DIST_DEV/$RELEASED_VERSION"; then
            echo "❌ Got non-0 exit code from removing the staged release candidates, aborting."
            exit 1
        fi
        echo "✅ Removed $DIST_DEV/$RELEASED_VERSION."
    else
        echo "✅ Leaving the staged release candidates alone."
    fi
else
    echo "✅ Nothing of $RELEASED_VERSION left to clean up in SVN."
fi

########################################################################################################################
# 5b. Remove the previous release from the mirrors
########################################################################################################################

# A lot of mirrors serve our releases, so Apache policy is that only the current one stays in the
# release part of SVN. Older releases remain available from archive.apache.org, and
# downloads.apache.org redirects to it for anything it no longer carries, so the links on the
# download page keep working.
#
# Only directories that look like a version are considered - "KEYS", "build-tools" and
# "plc4x-extras" live in the same place and are none of our business here.

# Only versions OLDER than the one just released. Excluding just "$RELEASED_VERSION" would, for a
# bugfix release off an old branch (0.13.2 while 1.0.0 is current), offer to delete 1.0.0 - the
# release people are actually downloading.
OLD_RELEASES=$(svn ls "$DIST_RELEASE/" 2>/dev/null \
    | grep -E '^[0-9]+\.[0-9]+\.[0-9]+/$' | sed 's|/$||' \
    | while read -r candidate; do
        if [[ "$candidate" != "$RELEASED_VERSION" ]] \
                && [[ "$(printf '%s\n%s\n' "$candidate" "$RELEASED_VERSION" | sort -V | head -1)" == "$candidate" ]]; then
            echo "$candidate"
        fi
      done || true)

if [[ -z "$OLD_RELEASES" ]]; then
    echo "✅ $DIST_RELEASE only carries $RELEASED_VERSION."
else
    echo
    echo "These older releases are still being served by the mirrors:"
    echo "$OLD_RELEASES" | sed 's|^|  |'
    echo "They stay available from https://archive.apache.org/dist/plc4x/ after removal."
    read -r -p "Remove them? (yes/no) " yn
    if [[ "$yn" == "yes" ]]; then
        # shellcheck disable=SC2046
        if ! svn rm -m "Removed the older releases of PLC4X now that $RELEASED_VERSION is out" \
                $(echo "$OLD_RELEASES" | sed "s|^|$DIST_RELEASE/|"); then
            echo "❌ Got non-0 exit code from removing the older releases, aborting."
            exit 1
        fi
        echo "✅ Removed: $(echo "$OLD_RELEASES" | tr '\n' ' ')"
    else
        echo "✅ Leaving the older releases in place."
    fi
fi

########################################################################################################################
# 5c. Drop the staging repositories that were never released
########################################################################################################################

# A vote that had to be repeated leaves a closed staging repository behind for every attempt.
# Nexus does not tell us which of them belong to which attempt in any way we could rely on, so
# rather than guessing this lists what is there and lets you pick.

echo
echo "Staging repositories currently on $NEXUS_URL:"
# No "-q" here: the repository table is printed at INFO level, and quiet mode would leave the
# prompt below asking for ids from a listing that was never shown.
if ! "$DIRECTORY/mvnw" -f "$DIRECTORY/tools/stage.pom" nexus-staging:rc-list; then
    echo "⚠️  Could not list the staging repositories - check them by hand at"
    echo "   $NEXUS_URL under 'Staging Repositories'."
else
    echo
    echo "Enter the ids of the repositories that should be dropped, separated by ','."
    echo "Leave empty to keep all of them. The one released above is already gone."
    read -r -p "Drop: " DROP_REPOS
    DROP_REPOS=${DROP_REPOS// /}
    if [[ -z "$DROP_REPOS" ]]; then
        echo "✅ Keeping all staging repositories."
    else
        read -r -p "Really drop '$DROP_REPOS'? (yes/no) " yn
        if [[ "$yn" != "yes" ]]; then
            echo "✅ Keeping all staging repositories."
        elif ! "$DIRECTORY/mvnw" -f "$DIRECTORY/tools/stage.pom" nexus-staging:rc-drop \
                -DstagingRepositoryId="$DROP_REPOS"; then
            echo "❌ Got non-0 exit code from dropping the staging repositories, aborting."
            exit 1
        else
            echo "✅ Dropped: $DROP_REPOS"
        fi
    fi
fi
