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

# Required for "with-c" profile
RUN apt install -y build-essential

# Required for "with-dotnet" profile
RUN apt install -y wget
RUN wget -q https://packages.microsoft.com/config/ubuntu/20.04/packages-microsoft-prod.deb -O packages-microsoft-prod.deb
RUN dpkg -i packages-microsoft-prod.deb
RUN apt install -y software-properties-common
RUN add-apt-repository universe -y
RUN apt install -y apt-transport-https
RUN apt update -y
RUN apt install -y dotnet-sdk-6.0

# Required for "with-python" profile
RUN apt install -y python3 python3-venv python3-pip
RUN pip3 install wheel

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

# Build everything with all tests
# (Skip signing, as this requires access to the local GPG keys)
RUN ./mvnw -Dskip-pgp-signing=true -P with-c,with-dotnet,with-go,with-python,with-sandbox,enable-all-checks,apache-release install

ENTRYPOINT ["/bin/bash"]
