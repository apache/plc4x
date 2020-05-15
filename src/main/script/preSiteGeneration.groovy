/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */

import org.apache.maven.repository.internal.MavenRepositorySystemUtils

print "\nCalculating some additional properties:"
// Get the current year
def currentYear = Calendar.getInstance().get(Calendar.YEAR) as String
print "\nCurrent year:                     " + currentYear
project.properties['current-year'] = currentYear

// Calculate some version related stuff
def currentVersion = project.version as String
def match = (currentVersion =~ /(\d+)\.(\d+)\.(\d+)(-SNAPSHOT)?/)
print "\nCurrent version:                  " + currentVersion
if (match.count >= 1) {
    def majorVersion = match[0][1] as Integer
    def minorVersion = match[0][2] as Integer
    def bugfixVersion = match[0][3] as Integer

    def currentFullVersion = majorVersion + "." + minorVersion + "." + bugfixVersion
    def currentShortVersion = majorVersion + "." + minorVersion
    def currentNextIncrementalVersion = majorVersion + "." + minorVersion + "." + (bugfixVersion + 1)
    def currentNextMinorVersion = majorVersion + "." + (minorVersion + 1) + ".0"

    print "\nCurrent full version:             " + currentFullVersion + " (current-full-version)"
    project.properties['current-full-version'] = currentFullVersion
    print "\nCurrent short version:            " + currentShortVersion + "   (current-short-version)"
    project.properties['current-short-version'] = currentShortVersion
    print "\nCurrent next incremental version: " + currentNextIncrementalVersion + " (current-next-incremental-version)"
    project.properties['current-next-incremental-version'] = currentNextIncrementalVersion
    print "\nCurrent next minor version:       " + currentNextMinorVersion + " (current-next-minor-version)"
    project.properties['current-next-minor-version'] = currentNextMinorVersion
}

// Get the latest released version from our doap file.
def plc4xDoapFile = new File(project.getBasedir(), "src/site/resources-filtered/plc4x-doap.rdf")
if (plc4xDoapFile.exists()) {
    def doapDocument = new XmlSlurper().parse(plc4xDoapFile)
    def lastReleasedVersion = doapDocument.Project.release[0].Version.revision
    print "\nCurrent last released version:    " + lastReleasedVersion + " (current-last-released-version)"
    project.properties['current-last-released-version'] = lastReleasedVersion
}
print "\n"

// Dump all maven properties into a file readable by asciidoc.
print "\nGenerating 'pom.adoc' file in target directory (Use in adoc by including 'include::{pom-adoc} ... however this should not be required)'\n\n"
def propertyFile = new File(project.getBasedir(), "target/pom.adoc")
// Ensure the parent directory is created
propertyFile.getParentFile().mkdirs()
// Make sure the file is deleted so we don't append to an existing file
propertyFile.delete()
// Iterate over all properties and dump them to the file
project.properties.each { entry -> propertyFile << ":$entry.key: $entry.value\n" }

