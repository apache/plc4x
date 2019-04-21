#!/usr/bin/env bash

# Clean the directory
rm -f lib/*.jar

# Builds the interop server module with maven
mvn clean package -f ../plc4j/utils/interop/pom.xml

# Moves the server from the target folder here
cp ../plc4j/utils/interop/target/apache-plc4x-incubating-*-SNAPSHOT-jar-with-dependencies.jar lib/interop-server.jar