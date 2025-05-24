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

DIRECTORY=..

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
# 1. Check if there are uncommitted changes as these would automatically be committed (local)
########################################################################################################################

if [[ $(git status --porcelain) ]]; then
  # Changes
  echo "❌ There are untracked files or changed files, aborting."
  exit 1
fi

########################################################################################################################
# 2. Delete the pre-exising "out" directory that contains the maven local repo and deployments (local)
########################################################################################################################

echo "Deleting the maven local repo and previous deployments"
rm -r $DIRECTORY/out

########################################################################################################################
# 3. Delete all generated sources (local)
########################################################################################################################

echo "Deleting generated-sources:"
find "$DIRECTORY" -path "*/src/main/generated" -print0 | while IFS= read -r -d '' f; do
    echo " - Deleting: $f"
    rm -r "$f"
done
# Delete the PLC4C code (local)
echo " - Deleting:  $DIRECTORY/plc4c/generated-sources"
rm -r "$DIRECTORY/plc4c/generated-sources"
# Delete the PLC4Go code (local)
echo " - Deleting:  generated files in $DIRECTORY/plc4c/generated-sources"
find "$DIRECTORY/plc4go/protocols" -mindepth 2 -type f ! \( -name 'StaticHelper.go' -o -name 'StaticHelper_test.go' \) -exec rm -v {} \;
# Delete the PLC4Net code (local)
echo " - Deleting:  generated files in $DIRECTORY/plc4net/drivers"
for dir in "$DIRECTORY/plc4net/drivers"/*; do
    # Delete generated classes
    if [[ -d "$dir" && ! "$(basename "$dir")" =~ -test$ ]]; then
        SRC_DIR="$dir/src"
        if [[ -d "$SRC_DIR" ]]; then
            echo "🧹 Deleting files in: $SRC_DIR"
            find "$SRC_DIR" -type f -exec rm -v {} \;
        fi
    else
        SRC_DIR="$dir/resources"
        if [[ -d "$SRC_DIR" ]]; then
            echo "🧹 Deleting files in: $SRC_DIR"
            find "$SRC_DIR" -type f -exec rm -v {} \;
        fi
    fi
done
# Delete the PLC4Py code (local)
echo " - Deleting:  generated files in $DIRECTORY/plc4py/plc4py/protocols"
find "$DIRECTORY/plc4py/plc4py/protocols" -mindepth 2 -type f ! \( -name '__init__.py' -o -name 'StaticHelper.py' \) -exec rm -v {} \;

########################################################################################################################
# 4. Make sure the NOTICE file has the current year in the second line
########################################################################################################################

NOTICE_FILE="../NOTICE"
CURRENT_YEAR=$(date +%Y)
EXPECTED="Copyright 2017-${CURRENT_YEAR} The Apache Software Foundation"

# Extract the second line
SECOND_LINE=$(sed -n '2p' "$NOTICE_FILE")

if [[ "$SECOND_LINE" != "$EXPECTED" ]]; then
    echo "✏️  Updating $NOTICE_FILE"

    # Replace line 2 with the expected text
    awk -v expected="$EXPECTED" 'NR==2 {$0=expected} {print}' "$NOTICE_FILE" > "$NOTICE_FILE.tmp" &&
    mv "$NOTICE_FILE.tmp" "$NOTICE_FILE"
else
    echo "✅ $NOTICE_FILE is already up to date."
fi

########################################################################################################################
# 5 Run the maven build for all modules with "update-generated-code" enabled (Docker container)
########################################################################################################################

if ! docker compose build; then
    echo "❌ Got non-0 exit code from building the release docker container, aborting."
    exit 1
fi

if ! docker compose run releaser bash /ws/mvnw -e -P with-c,with-dotnet,with-go,with-java,with-python,enable-all-checks,update-generated-code -Dmaven.repo.local=/ws/out/.repository clean package -DskipTests; then
    echo "❌ Got non-0 exit code from running the code-generation inside docker, aborting."
    exit 1
fi

########################################################################################################################
# 6. Make sure the generated driver documentation is up-to-date.
########################################################################################################################

if ! docker compose run releaser bash /ws/mvnw -e -P with-java -Dmaven.repo.local=/ws/out/.repository clean site -pl :plc4j-driver-all; then
    echo "❌ Got non-0 exit code from running the site code-generation inside docker, aborting."
    exit 1
fi

########################################################################################################################
# 7. Commit and push any changed files
########################################################################################################################

if [[ $(git status --porcelain) ]]; then
  echo "Committing changes."
  git add --all
  git commit -m "chore: updated generated code"
  git push
else
  echo "No changes."
fi

echo "✅ Pre-release updates complete. Please continue with 'release-1-create-branch.sh' next."