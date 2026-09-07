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
