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

# 0. Check if the release properties file exists.
#if [[ `git status --porcelain` ]]; then
#  # Changes
#  echo "There are untracked files or changed files, aborting."
#  exit 1
#fi

# 1. Do a simple release-perfoem command skip sigining of artifacts and deploy to local directory (inside the Docker container)
docker compose run --rm releaser bash /ws/mvnw -e -Dmaven.repo.local=/ws/out/.repository -DaltDeploymentRepository=snapshot-repo::default::file:/ws/out/.local-artifacts-dir release:perform

# 2. Sign the artifacts
# TODO: it seems the gpg sign plugin only signs one artifact ... gotta find out how to sign every jar in out/.local-artifacts-dir

# 3. Deploy the artifacts to Nexus
# TODO: Use the same technique we use on Jenkins to deploy everything in a local repo

# 4. Prepare a directory for the release candidate
# TODO: Implement ...

# 5. Upload the release candidate artifacts to SVN
# TODO: Implement ...
