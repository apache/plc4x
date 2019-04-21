#!/usr/bin/env bash

# Builds the interop server module with maven
mvn clean package -f ../plc4j/utils/interop/pom.xml

# Moves the server from the target folder here
cp ../plc4j/utils/interop/target/apache-plc4x-incubating-* lib/