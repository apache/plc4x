#  Licensed to the Apache Software Foundation (ASF) under one or more
#  contributor license agreements.  See the NOTICE file distributed with
#  this work for additional information regarding copyright ownership.
#  The ASF licenses this file to You under the Apache License, Version 2.0
#  (the "License"); you may not use this file except in compliance with
#  the License.  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
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
FROM azul/zulu-openjdk:latest as build

# Install some stuff we need to run the build
RUN apt update -y

# Install general purpose tools
RUN apt install -y make libpcap-dev libc-dev

# Requied for "with-boost" profile
RUN apt install -y bison flex gcc g++

# Required for "with-cpp" profile
RUN apt install -y gcc g++

# Required for "with-dotnet" profile
RUN apt install -y wget
RUN wget -q https://packages.microsoft.com/config/ubuntu/18.04/packages-microsoft-prod.deb -O packages-microsoft-prod.deb
RUN dpkg -i packages-microsoft-prod.deb
RUN apt install -y software-properties-common
RUN add-apt-repository universe -y
RUN apt install -y apt-transport-https
RUN apt update -y
RUN apt install -y dotnet-sdk-2.2

# Required for "with-java" profile
RUN apt install -y git

# Required for "with-proxies"
RUN apt install -y bison flex gcc g++

# Required for "with-python" profile
RUN apt install -y python-setuptools python

# Copy the project into the docker container
COPY . /ws/

# Change the working directory (where commands are executed) into the new "ws" directory
WORKDIR /ws

# Tell Maven to fetch all needed dependencies first, so they can get cached
# (Tried a patched version of the plugin to allow exclusion of inner artifacts.
# See https://issues.apache.org/jira/browse/MDEP-568 for details)
RUN ./mvnw -P with-java,with-cpp,with-boost,with-dotnet,with-python,with-proxies,with-sandbox com.offbytwo.maven.plugins:maven-dependency-plugin:3.1.1.MDEP568:go-offline -DexcludeGroupIds=org.apache.plc4x,org.apache.plc4x.examples,org.apache.plc4x.sandbox
# Build everything with all tests
RUN ./mvnw -P with-java,with-cpp,with-boost,with-dotnet,with-python,with-proxies,with-sandbox install

#RUN ./mvnw -P with-java -DskipTests install

# Get the version of the project and save it in a local file on the container
RUN ./mvnw org.apache.maven.plugins:maven-help-plugin:3.2.0:evaluate -Dexpression=project.version -DforceStdout -q -pl . > project_version

# Move the file to a place we can reference it from without a version
RUN PROJECT_VERSION=`cat project_version`; mv plc4j/examples/hello-storage-elasticsearch/target/plc4j-hello-storage-elasticsearch-$PROJECT_VERSION-uber-jar.jar plc4xdemo.jar

##########################################################################################
# Build a demo container
##########################################################################################

# Now create an actuall deployment container
FROM azul/zulu-openjdk:latest

# Prepare the demo by copying the example artifact from the 'build' container into this new one.
COPY --from=build /ws/plc4xdemo.jar /plc4xdemo.jar

# This will be executed as soon as the container is started.
ENTRYPOINT ["java", "-jar", "/plc4xdemo.jar"]