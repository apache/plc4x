/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import org.apache.maven.repository.internal.MavenRepositorySystemUtils

print "\nCalculating some additional properties:"
// Get the current year
def currentYear = Calendar.getInstance().get(Calendar.YEAR) as String
print "\nCurrent year:                     " + currentYear
project.properties['current-year'] = currentYear

// Calculate some version related stuff for the main project
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

// Get the version of the code-generation
codeGenerationVersion = project.properties['plc4x-code-generation.version']
def codeGenMatch = (codeGenerationVersion =~ /(\d+)\.(\d+)\.(\d+)(-SNAPSHOT)?/)
if (codeGenMatch.count >= 1) {
    def codeGenMajorVersion = codeGenMatch[0][1] as Integer
    def codeGenMinorVersion = codeGenMatch[0][2] as Integer
    def codeGenBugfixVersion = codeGenMatch[0][3] as Integer
    // If this is not a snapshot version, increment the minor version
    if(codeGenMatch[0][4] == null) {
        codeGenMinorVersion++
    }

    def codeGenerationReleaseFullVersion = codeGenMajorVersion + "." + codeGenMinorVersion + "." + codeGenBugfixVersion
    def codeGenerationReleaseShortVersion = codeGenMajorVersion + "." + codeGenMinorVersion
    def codeGenerationNextDevelopmentVersion = codeGenMajorVersion + "." + (codeGenMinorVersion + 1) + ".0"
    def codeGenerationBugfixShortVersion = codeGenMajorVersion + "." + codeGenMinorVersion + "." + (codeGenBugfixVersion + 1)

    print "\nNext code-generation full version:             " + codeGenerationReleaseFullVersion + " (code-generation-full-version)"
    project.properties['code-generation-full-version'] = codeGenerationReleaseFullVersion
    print "\nNext code-generation short version:             " + codeGenerationReleaseShortVersion + " (code-generation-short-version)"
    project.properties['code-generation-short-version'] = codeGenerationReleaseShortVersion
    print "\nNext code-generation development version:             " + codeGenerationNextDevelopmentVersion + " (code-generation-development-version)"
    project.properties['code-generation-development-version'] = codeGenerationNextDevelopmentVersion
    print "\nNext code-generation bugfix version:             " + codeGenerationBugfixShortVersion + " (code-generation-bugfix-version)"
    project.properties['code-generation-bugfix-version'] = codeGenerationBugfixShortVersion
}

// Get the latest released version from our doap file.
def plc4xDoapFile = new File(project.getBasedir(), "website/resources/plc4x-doap.rdf")
if (plc4xDoapFile.exists()) {
    def doapDocument = new XmlSlurper().parse(plc4xDoapFile)
    def lastReleasedVersion = doapDocument.Project.release[0].Version.revision
    print "\nCurrent last released version:    " + lastReleasedVersion + " (current-last-released-version)"
    project.properties['current-last-released-version'] = lastReleasedVersion
}
print "\n"

// Dump all maven properties into a file readable by asciidoc.
print "\nGenerating 'pom.adoc' file in target directory (Use in adoc by including 'include::{pom-adoc} ... however this should not be required)'\n\n"
def propertyFile = new File(project.getBasedir(), "website/target/keys/pom.adoc")
// Ensure the parent directory is created
propertyFile.getParentFile().mkdirs()
// Make sure the file is deleted so we don't append to an existing file
propertyFile.delete()
// Iterate over all properties and dump them to the file
project.properties.each { entry -> propertyFile << ":$entry.key: $entry.value\n" }

