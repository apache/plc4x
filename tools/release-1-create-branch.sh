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

# 0. Check if there are uncommited changes as these would automatically be committed
if [[ `git status --porcelain` ]]; then
  # Changes
  echo "There are untracked files or changed files, aborting."
  exit 1
fi

# 1. Get and calculate the current veerion
PROJECT_VERSION=$(../mvnw -f ../pom.xml -q -Dexec.executable=echo -Dexec.args='${project.version}' --non-recursive exec:exec)
RELEASE_VERSION=${PROJECT_VERSION%"-SNAPSHOT"}
RELEASE_SHORT_VERSION=${RELEASE_VERSION%".0"}
BRANCH_NAME="rel/$RELEASE_SHORT_VERSION"
IFS='.' read -ra VERSION_SEGMENTS <<< "$RELEASE_VERSION"
NEW_VERSION="${VERSION_SEGMENTS[0]}.$((VERSION_SEGMENTS[1] + 1)).0-SNAPSHOT"
echo "Current Version: '$PROJECT_VERSION'"
echo "Release Version: '$RELEASE_VERSION'"
echo "Release Banch Name: '$BRANCH_NAME'"
echo "New develop Version: '$NEW_VERSION'"

# 2. Do a simple maven branch command with pushChanges=false (inside the Docker container)
docker compose run --rm releaser bash /ws/mvnw -e -P with-c,with-dotnet,with-go,with-python,with-sandbox -Dmaven.repo.local=/ws/out/.repository release:branch -DautoVersionSubmodules=true -DpuchChanges=false -DdevelopmentVersion=$NEW_VERSION -DbranchName=$BRANCH_NAME
if [ $? -ne 0 ]; then
    echo "Got non-0 exit code from docker compose, aborting."
    exit 1
fi

# 3. Push the changes (outside)
git push