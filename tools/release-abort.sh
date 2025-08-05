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
PRE_RELEASE_VERSION=0.13.0-SNAPSHOT

# Set the local development version back to the initial one.
docker compose run releaser bash /ws/mvnw -e -P with-c,with-dotnet,with-go,with-java,with-python,update-generated-code -Dmaven.repo.local=/ws/out/.repository versions:set -DnewVersion="$PRE_RELEASE_VERSION"

# Delete left-over files from the last attempt.
find .. -type f -name 'release.properties' -delete
find .. -type f -name 'pom.xml.versionsBackup' -delete
find .. -type f -name 'pom.xml.releaseBackup' -delete

# delete branch locally
echo "git branch -d rel/{version-short}"

# delete branch remotely
echo "git push origin --delete rel/{version-short}"

# delete tag locally
echo "git tag -d v{version}"

# delete tag remotely
echo "git push origin --delete v{version}"
