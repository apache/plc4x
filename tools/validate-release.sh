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
# 1. Check if the release properties file exists.
########################################################################################################################

########################################################################################################################
# 2. Do a simple release-perform command skip signing of artifacts and deploy to local directory
#    (inside the Docker container)
########################################################################################################################

echo "Validate Release:"
docker compose build
if ! docker compose run releaser bash /ws/mvnw  -e -P with-java -Dmaven.repo.local=/ws/out/.repository -Dreference.repo=https://repository.apache.org/content/repositories/staging/ -Dbuildinfo.reproducible verify artifact:compare;
then
    echo "❌ Got non-0 exit code from docker compose, aborting."
    exit 1
fi
