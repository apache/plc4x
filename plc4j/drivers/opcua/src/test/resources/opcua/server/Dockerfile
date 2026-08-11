#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

# OPC UA test server used by OpcuaPlcDriverTest.
#
# Eclipse Milo is deliberately NOT a Maven dependency of this module: it would drag its
# whole transitive tree (netty, jackson, guava, bouncycastle, ...) onto the test classpath,
# where every one of them shows up as a dependency-update / vulnerability alert for a jar
# that never leaves this container. Instead the server-examples uber jar - which already
# bundles the complete Milo stack plus its dependencies - is fetched here and the handful
# of PLC4X test-server classes under src/ are compiled against it inside the image.
#
# Consequence: the sources in src/ are not compiled by Maven, so the IDE will not flag
# breakage in them. Bumping MILO_VERSION is a manual step; a compile error surfaces as a
# failed image build when the integration test runs.

ARG MILO_VERSION=0.6.16

FROM eclipse-temurin:21-jdk AS build
ARG MILO_VERSION
ADD "https://repo1.maven.org/maven2/org/eclipse/milo/server-examples/${MILO_VERSION}/server-examples-${MILO_VERSION}.jar" /opt/milo/server-examples.jar
COPY src /src
RUN javac -nowarn -cp /opt/milo/server-examples.jar -d /opt/milo/classes $(find /src -name '*.java')

FROM eclipse-temurin:21-jre
COPY --from=build /opt/milo /opt/milo

# The server binds this port and also advertises its endpoints on it, so the test container
# maps it 1:1 onto the host (see MiloTestContainer).
EXPOSE 12686

CMD ["java", "-cp", "/opt/milo/server-examples.jar:/opt/milo/classes", "org.eclipse.milo.examples.server.TestMiloServer"]
