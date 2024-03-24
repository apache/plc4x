#!/bin/bash
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

# TODO: Reproducible builds
# - Creating an RC
#   - We should run the following builds on the normal developer environment:
#     - build with all "with-*" profiles and "update-generated-code" profile enabled
#     - release:branch
#     - release:prepare
#   - The release build (release:perform) should then be run inside the docker-compose build in order to build the jars in the controlled environment
#   - Run the signing on all locally staged artifacts
#   - Transfer all locally staged artifacts to nexus
#
# - Validating an RC
#   - Download the entire content of the staging repo locally
#   - Programmatically check the Hashes and the signatures with a script
#   - Run a local build in docker-compose but with the "clean validate {compare-artifact}" config



# Run a standard build
function build() {
  echo "Building ..."
  mvn -U -P apache-release,with-c,with-dotnet,with-go,with-java,with-python -DaltDeploymentRepository=snapshot-repo::default::file:./local-snapshots-dir clean deploy
  echo "Done"
}

# Just remove all the metadata information and all the hashes and signatures
function clean() {
  echo "Cleaning ..."
  cd local-snapshots-dir || exit
  find ./ -name "*.asc" -exec rm -rf {} \;
  find ./ -name "*.md5" -exec rm -rf {} \;
  find ./ -name "*.sha1" -exec rm -rf {} \;
  find ./ -name "maven-metadata.xml" -exec rm -rf {} \;
  find ./ -name "maven-metadata.xml.md5" -exec rm -rf {} \;
  find ./ -name "maven-metadata.xml.sha1" -exec rm -rf {} \;
  cd ..
  echo "Done"
  pwd
}

# Rename all snapshot files to not contain the timestamps
function renameArtifacts() {
  echo "Renaming ..."
  cd local-snapshots-dir || exit
  # Remove the timestamp from the file-names
  find . -type f | rename 's/-\d{8}\.\d{6}-\d{1,2}//'
  cd ..
  echo "Done"
  pwd
}

# Package the remaining files into one tgz archive
function packageDirectory() {
  echo "Packaging ..."
  tar -cvf reproducible-build-candidate.tgz local-snapshots-dir/
  echo "Done"
  pwd
}

# Remove any pre-existing directory
rm -r local-snapshots-dir

build
clean
renameArtifacts
packageDirectory
