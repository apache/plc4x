#  Licensed to the Apache Software Foundation (ASF) under one or more
#  contributor license agreements.  See the NOTICE file distributed with
#  this work for additional information regarding copyright ownership.
#  The ASF licenses this file to You under the Apache License, Version 2.0
#  (the "License"); you may not use this file except in compliance with
#  the License.  You may obtain a copy of the License at
#
#      https://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.

##########################################################################################
# Build PLC4X
##########################################################################################

# This is the image we'll use to execute the build (and give it the name 'build').
# (This image is based on Ubuntu)
# Fixed version of this in order to have a fixed JDK version
FROM azul/zulu-openjdk:11 as build

# Install some stuff we need to run the build
RUN apt update -y

# Install general purpose tools
RUN apt install -y make libpcap-dev libc-dev git

# Required for "with-boost" profile
#RUN apt install -y bison flex gcc g++

# Required for "with-cpp" profile
#RUN apt install -y gcc g++

# Required for "with-proxies" and "with-cpp"
#RUN apt install -y clang

# Required for "with-proxies" and "with-cpp"
#RUN apt install -y cmake

# Required for "with-dotnet" profile
RUN apt install -y wget
RUN wget -q https://packages.microsoft.com/config/ubuntu/20.04/packages-microsoft-prod.deb -O packages-microsoft-prod.deb
RUN dpkg -i packages-microsoft-prod.deb
RUN apt install -y software-properties-common
RUN add-apt-repository universe -y
RUN apt install -y apt-transport-https
RUN apt update -y
RUN apt install -y dotnet-sdk-3.1

# Required for "with-go" profile
RUN apt install -y golang

# Required for "with-python" profile
RUN apt install -y python-setuptools python

# Required for running on Windows systems
RUN apt install -y dos2unix

# Copy the project into the docker container
COPY . /ws/

# Change the working directory (where commands are executed) into the new "ws" directory
WORKDIR /ws

# Make the maven wrapper script executalbe (needed when running on Windows)
RUN chmod +x ./mvnw
# Change the line ending to unix-style (needed when running on Windows)
RUN dos2unix ./mvnw
RUN dos2unix .mvn/wrapper/maven-wrapper.properties

# Tell Maven to fetch all needed dependencies first, so they can get cached
# (Tried a patched version of the plugin to allow exclusion of inner artifacts.
# See https://issues.apache.org/jira/browse/MDEP-568 for details)
#RUN ./mvnw -P with-boost,with-c,with-cpp,with-dotnet,with-go,with-logstash,with-python,with-sandbox com.offbytwo.maven.plugins:maven-dependency-plugin:3.1.1.MDEP568:go-offline -DexcludeGroupIds=org.apache.plc4x,org.apache.plc4x.examples,org.apache.plc4x.sandbox
RUN ./mvnw -P with-c,with-dotnet,with-go,with-python,with-sandbox com.offbytwo.maven.plugins:maven-dependency-plugin:3.1.1.MDEP568:go-offline -DexcludeGroupIds=org.apache.plc4x,org.apache.plc4x.examples,org.apache.plc4x.sandbox

# Build everything with all tests
#RUN ./mvnw -P skip-prerequisite-check,with-boost,with-c,with-cpp,with-dotnet,with-go,with-logstash,with-python,with-sandbox install
RUN ./mvnw -P with-c,with-dotnet,with-go,with-python,with-sandbox install

# Get the version of the project and save it in a local file on the container
RUN ./mvnw org.apache.maven.plugins:maven-help-plugin:3.2.0:evaluate -Dexpression=project.version -DforceStdout -q -pl . > project_version

##########################################################################################
# Build a demo container
##########################################################################################

# Move the file to a place we can reference it from without a version
RUN PROJECT_VERSION=`cat project_version`; mv plc4j/examples/hello-integration-iotdb/target/plc4j-hello-integration-iotdb-${PROJECT_VERSION}-uber-jar.jar plc4xdemo.jar

# Build a highly optimized JRE
FROM alpine:3.10 as packager

# Install regular JDK
RUN apk update
RUN apk --no-cache add openjdk11-jdk openjdk11-jmods

# build minimal JRE
ENV JAVA_MINIMAL="/opt/java-minimal"
RUN /usr/lib/jvm/java-11-openjdk/bin/jlink \
    --verbose \
    --add-modules \
        java.base,java.sql,java.naming,java.desktop,java.management,java.security.jgss,java.instrument \
    --compress 2 --strip-debug --no-header-files --no-man-pages \
    --release-info="add:IMPLEMENTOR=radistao:IMPLEMENTOR_VERSION=radistao_JRE" \
    --output "$JAVA_MINIMAL"

# Now create an actual deployment container
FROM alpine:3.10

# Install our optimized JRE
ENV JAVA_HOME=/opt/java-minimal
ENV PATH="$PATH:$JAVA_HOME/bin"
COPY --from=packager "$JAVA_HOME" "$JAVA_HOME"

# Prepare the demo by copying the example artifact from the 'build' container into this new one.
COPY --from=build /ws/plc4xdemo.jar /plc4xdemo.jar

# Let runtime know which ports we will be listening on
EXPOSE 9200 9300

# Allow for extra options to be passed to the jar using PLC4X_OPTIONS env variable
ENV PLC4X_OPTIONS ""

# This will be executed as soon as the container is started.
ENTRYPOINT ["sh", "-c", "[ -f /run/plc4xdemo.env ] && . /run/plc4xdemo.env ; java -jar /plc4xdemo.jar $PLC4X_OPTIONS"]
