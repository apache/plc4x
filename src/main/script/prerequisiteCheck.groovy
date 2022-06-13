import java.util.regex.Matcher

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

allConditionsMet = true

baseDirectory = project.model.pomFile.parent

/*
 Checks if a given version number is at least as high as a given reference version.
*/

def checkVersionAtLeast(String current, String minimum) {
    def currentSegments = current.tokenize('.')
    def minimumSegments = minimum.tokenize('.')
    def numSegments = Math.min(currentSegments.size(), minimumSegments.size())
    for (int i = 0; i < numSegments; ++i) {
        def currentSegment = currentSegments[i].toInteger()
        def minimumSegment = minimumSegments[i].toInteger()
        if (currentSegment < minimumSegment) {
            println current.padRight(14) + " FAILED (required min " + minimum + " but got " + current + ")"
            return false
        } else if (currentSegment > minimumSegment) {
            println current.padRight(14) + " OK"
            return true
        }
    }
    def curNotShorter = currentSegments.size() >= minimumSegments.size()
    if (curNotShorter) {
        println current.padRight(14) + " OK"
    } else {
        println current.padRight(14) + " (required min " + minimum + " but got " + current + ")"
    }
    curNotShorter
}

def checkVersionAtMost(String current, String maximum) {
    def currentSegments = current.tokenize('.')
    def maximumSegments = maximum.tokenize('.')
    def numSegments = Math.min(currentSegments.size(), maximumSegments.size())
    for (int i = 0; i < numSegments; ++i) {
        def currentSegment = currentSegments[i].toInteger()
        def maximumSegment = maximumSegments[i].toInteger()
        if (currentSegment > maximumSegment) {
            println current.padRight(14) + " FAILED (required max " + maximum + " but got " + current + ")"
            return false
        } else if (currentSegment < maximumSegment) {
            println current.padRight(14) + " OK"
            return true
        }
    }
    def curNotShorter = currentSegments.size() >= maximumSegments.size()
    if (curNotShorter) {
        println current.padRight(14) + " OK"
    } else {
        println current.padRight(14) + " (required max " + maximum + " but got " + current + ")"
    }
    curNotShorter
}

def checkBison() {
    print "Detecting Bison version:   "
    def output
    try {
        output = "bison --version".execute().text
    } catch (IOException e) {
        output = ""
    }
    Matcher matcher = extractVersion(output)
    if (matcher.size() > 0) {
        def curVersion = matcher[0][1]
        def result = checkVersionAtLeast(curVersion, "2.4.0")
        if (!result) {
            allConditionsMet = false
        }

        // TODO: Ensure the path of the `bison` binary doesn't contain any spaces.
    } else {
        println "missing"
        allConditionsMet = false
    }
}

def checkDotnet() {
    print "Detecting Dotnet version:  "
    def output
    try {
        output = "dotnet --version".execute().text
    } catch (IOException e) {
        output = ""
    }
    Matcher matcher = extractVersion(output)
    if (matcher.size() > 0) {
        def curVersion = matcher[0][1]
        def result = checkVersionAtLeast(curVersion, "4.5.2")
        if (!result) {
            allConditionsMet = false
        }
    } else {
        println "missing"
        allConditionsMet = false
    }
}

def checkGo() {
    print "Detecting Go version:      "
    def output
    try {
        output = "go version".execute().text
    } catch (IOException e) {
        output = ""
    }
    Matcher matcher = extractVersion(output)
    if (matcher.size() > 0) {
        def curVersion = matcher[0][1]
        def result = checkVersionAtLeast(curVersion, "1.0.0")
        if (!result) {
            allConditionsMet = false
        }
    } else {
        println "missing"
        allConditionsMet = false
    }
}


def checkJavaVersion(String minVersion, String maxVersion) {
    print "Detecting Java version:    "
    def curVersion = System.properties['java.version']
    def result
    if (minVersion != null) {
        result = checkVersionAtLeast(curVersion, minVersion)
        if (!result) {
            allConditionsMet = false
            return
        }
    }
    if (maxVersion != null) {
        result = checkVersionAtMost(curVersion, maxVersion)
        if (!result) {
            allConditionsMet = false
            return
        }
    }
}

def checkMavenVersion(String minVersion, String maxVersion) {
    print "Detecting Maven version:   "
    def curVersion = project.projectBuilderConfiguration.systemProperties['maven.version']
    def result
    if (minVersion != null) {
        result = checkVersionAtLeast(curVersion, minVersion)
        if (!result) {
            allConditionsMet = false
            return
        }
    }
    if (maxVersion != null) {
        result = checkVersionAtMost(curVersion, maxVersion)
        if (!result) {
            allConditionsMet = false
            return
        }
    }
}

def checkFlex() {
    print "Detecting Flex version:    "
    def output
    try {
        output = "flex --version".execute().text
    } catch (IOException e) {
        output = ""
    }
    Matcher matcher = extractVersion(output)
    if (matcher.size() > 0) {
        def curVersion = matcher[0][1]
        def result = checkVersionAtLeast(curVersion, "2.0.0")
        if (!result) {
            allConditionsMet = false
        }
    } else {
        println "missing"
        allConditionsMet = false
    }
}

def checkGcc() {
    print "Detecting Gcc version:     "
    // TODO: For windows, check that mingw32-make is on the PATH
    def output
    try {
        output = "gcc --version".execute().text
    } catch (IOException e) {
        output = ""
    }
    Matcher matcher = extractVersion(output)
    if (matcher.size() > 0) {
        def curVersion = matcher[0][1]
        def result = checkVersionAtLeast(curVersion, "1.0.0")
        if (!result) {
            allConditionsMet = false
        }
    } else {
        println "missing"
        allConditionsMet = false
    }
}

def checkGit() {
    print "Detecting Git version:     "
    def output
    try {
        output = "git --version".execute().text
    } catch (IOException e) {
        output = ""
    }
    Matcher matcher = extractVersion(output)
    if (matcher.size() > 0) {
        def curVersion = matcher[0][1]
        def result = checkVersionAtLeast(curVersion, "1.0.0")
        if (!result) {
            allConditionsMet = false
        }
    } else {
        println "missing"
        allConditionsMet = false
    }
}

def checkGpp() {
    print "Detecting G++ version:     "
    def output
    try {
        output = "g++ --version".execute().text
    } catch (IOException e) {
        output = ""
    }
    Matcher matcher = extractVersion(output)
    if (matcher.size() > 0) {
        def curVersion = matcher[0][1]
        def result = checkVersionAtLeast(curVersion, "1.0.0")
        if (!result) {
            allConditionsMet = false
        }
    } else {
        println "missing"
        allConditionsMet = false
    }
}

def checkClang() {
    print "Detecting clang version:   "
    def output
    try {
        output = "clang --version".execute().text
    } catch (IOException e) {
        output = ""
    }
    Matcher matcher = extractVersion(output)
    if (matcher.size() > 0) {
        def curVersion = matcher[0][1]
        def result = checkVersionAtLeast(curVersion, "1.0.0")
        if (!result) {
            allConditionsMet = false
        }
    } else {
        println "missing"
        allConditionsMet = false
    }
}

def checkCmake() {
    print "Detecting cmake version:   "
    def output
    try {
        output = "cmake --version".execute().text
    } catch (IOException e) {
        output = ""
    }
    Matcher matcher = extractVersion(output)
    if (matcher.size() > 0) {
        def curVersion = matcher[0][1]
        def result = checkVersionAtLeast(curVersion, "3.0.0")
        if (!result) {
            allConditionsMet = false
        }
    } else {
        println "missing"
        allConditionsMet = false
    }
}

def checkPython() {
    print "Detecting Python version:  "
    try {
        def process = ("python --version").execute()
        def stdOut = new StringBuilder()
        def stdErr = new StringBuilder()
        process.consumeProcessOutput(stdOut, stdErr)
        process.waitForOrKill(500)
        Matcher matcher = extractVersion(stdOut + stdErr)
        if (matcher.size() > 0) {
            def curVersion = matcher[0][1]
            def result = checkVersionAtLeast(curVersion, "3.6.0")
            if (!result) {
                allConditionsMet = false
            }
        } else {
            println "missing (Please install at least version 3.6.0)"
            allConditionsMet = false
        }
    } catch (Exception e) {
        println "missing"
        allConditionsMet = false
    }
}

def checkSetupTools() {
    print "Detecting setuptools:      "
    try {
        def cmdArray = ["python", "-c", "import setuptools"]
        def process = cmdArray.execute()
        def stdOut = new StringBuilder()
        def stdErr = new StringBuilder()
        process.consumeProcessOutput(stdOut, stdErr)
        process.waitForOrKill(500)
        if(stdErr.contains("No module named setuptools")) {
            println "missing"
            allConditionsMet = false
        } else {
            println "               OK"
        }
    } catch (Exception e) {
        println "missing"
        allConditionsMet = false
    }
}

/*
 * This check does an extremely simple check, if the boost library exists in the maven local repo.
 * We're not checking if it could be resolved.
 */

def checkBoost() {
    print "Detecting Boost library:   "
    def localRepoBaseDir = session.getLocalRepository().getBasedir()
    def expectedFile = new File(localRepoBaseDir, "org/apache/plc4x/plc4x-tools-boost/" + project.version +
        "/plc4x-tools-boost-" + project.version + "-lib-" + project.properties["os.classifier"] + ".zip")
    if (!expectedFile.exists()) {
        println "              missing"
        println ""
        println "Missing the Boost library. This has to be built by activating the Maven profile 'with-boost'. This only has to be built once."
        println ""
        allConditionsMet = false
    } else {
        println "              OK"
    }
}

def checkOpenSSL() {
    print "Detecting OpenSSL version: "
    def output
    try {
        output = "openssl version".execute().text
    } catch (IOException e) {
        output = ""
    }
    Matcher matcher = extractVersion(output)
    if (matcher.size() > 0) {
        def curVersion = matcher[0][1]
        def result = checkVersionAtLeast(curVersion, "1.0.0")
        if (!result) {
            allConditionsMet = false
        }
    } else {
        println "missing"
        allConditionsMet = false
    }
}

// When building the StreamPipes modules we need Docker.
// Not only should the docker executable be available, but also should the docker daemon be running.
def checkDocker() {
    print "Detecting Docker version:  "
    def output
    try {
        output = "docker info".execute().text
    } catch (IOException e) {
        output = ""
    }
    // Check if Docker is installed at all
    def matcher1 = output =~ /Server:/
    if (matcher1.size() > 0) {
        // If it is check if the daemon is running and if the version is ok
        def matcher2 = output =~ /Server Version: (\d+\.\d+(\.\d+)?).*/
        if (matcher2.size() > 0) {
            def curVersion = matcher2[0][1]
            def result = checkVersionAtLeast(curVersion, "1.0.0")
            if (!result) {
                allConditionsMet = false
            }
        } else {
            println "Docker daemon probably not running"
            allConditionsMet = false
        }
    } else {
        println "missing"
        allConditionsMet = false
    }
    // TODO: Implement the actual check ...
}

def checkLibPcap(String minVersion) {
    print "Detecting LibPcap version: "
    try {
        output = org.pcap4j.core.Pcaps.libVersion()
        String version = output - ~/^libpcap version /
        def result =  checkVersionAtLeast(version, minVersion)
        if (!result) {
            // TODO: only on mac we need the minimum version so we need to refine this
            // allConditionsMet = false
            println "This will probably a problem on mac"
        }
    } catch (Error e) {
        output = ""
        println "missing"
        allConditionsMet = false
    }
}

/**
 * Version extraction function/macro. It looks for occurrence of x.y or x.y.z
 * in passed input text (likely output from `program --version` command if found).
 *
 * @param input
 * @return
 */
private Matcher extractVersion(input) {
    def matcher = input =~ /(\d+\.\d+(\.\d+)?).*/
    matcher
}

/////////////////////////////////////////////////////
// Find out which OS and arch are bring used.
/////////////////////////////////////////////////////

def osString = project.properties['os.classifier']
def osMatcher = osString =~ /(.*)-(.*)/
if (osMatcher.size() == 0) {
    throw new RuntimeException("Currently unsupported OS")
}
def os = osMatcher[0][1]
def arch = osMatcher[0][2]
println "Detected OS:   " + os
println "Detected Arch: " + arch

/////////////////////////////////////////////////////
// Find out which profiles are enabled.
/////////////////////////////////////////////////////

def boostEnabled = false
def cEnabled = false
def cppEnabled = false
def dockerEnabled = false
def dotnetEnabled = false
def goEnabled = false
// Java is always enabled ...
def javaEnabled = true
def pythonEnabled = false
def sandboxEnabled = false
def apacheReleaseEnabled = false
def activeProfiles = session.request.activeProfiles
for (def activeProfile : activeProfiles) {
    if (activeProfile == "with-boost") {
        boostEnabled = true
    } else if (activeProfile == "with-c") {
        cEnabled = true
    } else if (activeProfile == "with-cpp") {
        cppEnabled = true
    } else if (activeProfile == "with-docker") {
        dockerEnabled = true
    } else if (activeProfile == "with-dotnet") {
        dotnetEnabled = true
    } else if (activeProfile == "with-go") {
        goEnabled = true
    } else if (activeProfile == "with-python") {
        pythonEnabled = true
    } else if (activeProfile == "with-sandbox") {
        sandboxEnabled = true
    } else if (activeProfile == "apache-release") {
        apacheReleaseEnabled = true
    }
}
println ""

// - Windows:
//     - Check the length of the path of the base dir as we're having issues with the length of paths being too long.
if (os == "windows") {
    File pomFile = project.model.pomFile
    if (pomFile.absolutePath.length() > 100) {
        println "On Windows we encounter problems with maximum path lengths. " +
            "Please move the project to a place it has a shorter base path " +
            "and run the build again."
        allConditionsMet = false;
    }
}

/////////////////////////////////////////////////////
// Do the actual checks depending on the enabled
// profiles.
/////////////////////////////////////////////////////

// Codegen requires at least java 9
checkJavaVersion("9", null)

if (dotnetEnabled) {
    checkDotnet()
}

if (goEnabled) {
    checkGo()
}

if (cppEnabled) {
    checkClang()
    // The cmake-maven-plugin requires at least java 11
    checkJavaVersion("11", null)
    checkGcc()
}

if (javaEnabled) {
    checkGit()
}

if (cEnabled) {
    // The cmake-maven-plugin requires at least java 11
    checkJavaVersion("11", null)
    checkGcc()
}

if (cppEnabled) {
    checkGpp()
}

if (pythonEnabled) {
    checkPython()
    checkSetupTools()
}

// Boost needs the visual-studio `cl` compiler to compile the boostrap.
if (boostEnabled && (os == "windows")) {
    // TODO: checkVisualStudio()
}

// We only need this check, if boost is not enabled but we're enabling cpp.
if (!boostEnabled && cppEnabled) {
    checkBoost()
}

if (sandboxEnabled && dockerEnabled) {
    checkDocker()
}

if (cppEnabled || cEnabled) {
    // CMake requires at least maven 3.6.0
    checkMavenVersion("3.6.0", null)
}

if (apacheReleaseEnabled) {
    // TODO: Check libpcap is installed
}

if (cppEnabled && (os == "windows")) {
    print "Unfortunately currently we don't support building the 'with-cpp' profile on windows. This will definitely change in the future."
    allConditionsMet = false
}

if (os == "mac") {
    // The current system version from mac crashes so we assert for a version coming with brew
    checkLibPcap("1.10.1")
}

if (!allConditionsMet) {
    throw new RuntimeException("Not all conditions met, see log for details.")
}
println ""
println "All known conditions met successfully."
println ""

// Things we could possibly check:
// - DNS Providers that return a default ip on unknown host-names
// - Availability and version of LibPCAP/NPCAP

