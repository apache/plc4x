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
# directory it is started from.
DIRECTORY="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# Values shared with the other release scripts (Nexus staging profile, dist.apache.org URLs).
if [[ ! -f "$DIRECTORY/tools/release-common.sh" ]]; then
    echo "❌ '$DIRECTORY/tools/release-common.sh' not found, aborting."
    exit 1
fi
# shellcheck source=release-common.sh
source "$DIRECTORY/tools/release-common.sh"


########################################################################################################################
# 0. Check Docker Memory Availability
########################################################################################################################

# Minimum required memory in bytes (12 GB)
REQUIRED_MEM=$((12 * 1024 * 1024 * 1024))

# Extract total memory from `docker system info`
TOTAL_MEM=$(docker system info --format '{{.MemTotal}}')

# Check if TOTAL_MEM was retrieved successfully
if [[ -z "$TOTAL_MEM" || "$TOTAL_MEM" -eq 0 ]]; then
    echo "❌ Unable to determine total Docker memory. Is Docker running?"
    exit 1
fi

# Compare and exit if not enough memory
if (( TOTAL_MEM < REQUIRED_MEM )); then
    echo "❌ Docker runtime has insufficient memory: $(awk "BEGIN {printf \"%.2f\", $TOTAL_MEM/1024/1024/1024}") GB"
    echo "   At least 12 GB is required. Aborting."
    exit 1
fi

########################################################################################################################
# 1. Check that this is actually a release, and not a development checkout
########################################################################################################################

# "artifact:compare" compares what is built here against what is staged in Nexus. That only means
# anything if this really is the release: run against a "develop" checkout it would compare
# SNAPSHOT artifacts against a release repository and report differences that say nothing about
# the release candidate.

# Maven 4 prefixes even quiet output with "[INFO] [stdout] ", so take the last token of the
# last line rather than the whole output.
PROJECT_VERSION=$("$DIRECTORY/mvnw" -f "$DIRECTORY/pom.xml" -q -Dexec.executable=echo -Dexec.args="\${project.version}" --non-recursive exec:exec | tail -n 1 | awk '{print $NF}')
if [[ -z "$PROJECT_VERSION" ]]; then
    echo "❌ Could not determine the project version, aborting."
    exit 1
fi
if [[ "$PROJECT_VERSION" =~ -SNAPSHOT$ ]]; then
    echo "❌ This is a SNAPSHOT checkout ($PROJECT_VERSION), aborting."
    echo "   Unpack the staged apache-plc4x-<version>-source-release.zip and run this in there,"
    echo "   or check out the release tag."
    exit 1
fi
echo "✅ Validating Apache PLC4X $PROJECT_VERSION"

########################################################################################################################
# 2. Do a simple release-perform command skip signing of artifacts and deploy to local directory
#    (inside the Docker container)
########################################################################################################################

echo "Validate Release:"
if ! docker compose -f "$DIRECTORY/tools/docker-compose.yaml" build; then
    echo "❌ Got non-0 exit code from building the release docker container, aborting."
    exit 1
fi

# Only the Java artifacts are compared: the C, .Net and Python ones are either platform specific or
# not published to Maven at all, so there is nothing in the staging repository to compare them to.
if ! docker compose -f "$DIRECTORY/tools/docker-compose.yaml" run releaser \
        bash /ws/mvnw -e -P with-java -Dmaven.repo.local=/ws/out/.repository \
        -Dreference.repo="$NEXUS_URL/content/repositories/staging/" \
        -Dbuildinfo.reproducible verify artifact:compare; then
    echo "❌ Got non-0 exit code from docker compose, aborting."
    exit 1
fi
echo "✅ The build of $PROJECT_VERSION matches the staged artifacts."
