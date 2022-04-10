#  Licensed to the Apache Software Foundation (ASF) under one
#  or more contributor license agreements.  See the NOTICE file
#  distributed with this work for additional information
#  regarding copyright ownership.  The ASF licenses this file
#  to you under the Apache License, Version 2.0 (the
#  "License"); you may not use this file except in compliance
#  with the License.  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing,
#  software distributed under the License is distributed on an
#  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
#  KIND, either express or implied.  See the License for the
#  specific language governing permissions and limitations
#  under the License.
-include .env

#VERSION := $(shell git describe --tags)
#BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")

# csharp related variables.

# Mvn related variables
MVNBIN := ../../mvnw

# Use linker flags to provide version/build settings
#LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# Redirect error output to a file, so we can show it in development mode.
STDERR := /tmp/.$(PROJECTNAME)-stderr.txt

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

## install: Install missing dependencies. Runs `mvn verify` internally.
install: csharp-get

## compile: Compile the binary.
compile: target
	@-touch $(STDERR)
	@-rm $(STDERR)
	@-$(MAKE) -s csharp-compile 2> $(STDERR)
	@cat $(STDERR) | sed -e '1s/.*/\nError:\n/'  | sed 's/make\[.*/ /' | sed "/^/s/^/     /" 1>&2

## test: Runs test. Runs `mvn test` internally.
test: compile
	$(MVNBIN) test

## clean: Clean build files. Runs `mvn clean` internally.
clean:
	@-$(MAKE) csharp-clean

csharp-compile: csharp-get csharp-build

csharp-get:
	@echo "  >  Checking if there is any missing dependencies..."
	$(MVNBIN) verify

csharp-install: csharp-build
	$(MVNBIN) install

csharp-clean:
	@echo "  >  Cleaning build cache"
	$(MVNBIN) clean

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
