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

def checkDotnet() {
    print "Detecting Dotnet version:  "
    def output = new StringBuffer()
    def errOutput = new StringBuffer()
    try {
        def proc = "dotnet --version".execute()
        proc.waitForProcessOutput(output, errOutput)
    } catch (IOException e) {
        output = ""
        errOutput.append(e)
    }
    Matcher matcher = extractVersion(output)
    if (matcher.size() > 0) {
        String curVersion = matcher[0][1]
        def result = checkVersionAtLeast(curVersion, "4.5.2")
        if (!result) {
            allConditionsMet = false
        }
    } else {
        println "missing"
        println "--- output of version `dotnet --version` command ---"
        println output
        println errOutput
        println "----------------------------------------------------"
        allConditionsMet = false
    }
}

def checkJavaVersion(String minVersion, String maxVersion) {
    print "Detecting Java version:    "
    String curVersion = System.properties['java.version'].split("-")[0]

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
        }
    }
}

def checkMavenVersion(String minVersion, String maxVersion) {
    print "Detecting Maven version:   "
    String curVersion = project.projectBuilderConfiguration.systemProperties['maven.version']
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
        }
    }
}

def checkGcc() {
    print "Detecting Gcc version:     "
    // TODO: For windows, check that mingw32-make is on the PATH
    def output = new StringBuffer()
    def errOutput = new StringBuffer()
    try {
        def proc = "gcc --version".execute()
        proc.waitForProcessOutput(output, errOutput)
    } catch (IOException ignored) {
        output = ""
        errOutput.append(e)
    }
    Matcher matcher = extractVersion(output)
    if (matcher.size() > 0) {
        String curVersion = matcher[0][1]
        def result = checkVersionAtLeast(curVersion, "1.0.0")
        if (!result) {
            allConditionsMet = false
        }
    } else {
        println "missing"
        println "--- output of version `gcc --version` command ---"
        println output
        println errOutput
        println "-------------------------------------------------"
        allConditionsMet = false
    }
}

def checkGit() {
    print "Detecting Git version:     "
    def output = new StringBuffer()
    def errOutput = new StringBuffer()
    try {
        def proc = "git --version".execute()
        proc.waitForProcessOutput(output, errOutput)
    } catch (IOException ignored) {
        output = ""
        errOutput.append(e)
    }
    Matcher matcher = extractVersion(output)
    if (matcher.size() > 0) {
        String curVersion = matcher[0][1]
        def result = checkVersionAtLeast(curVersion, "1.0.0")
        if (!result) {
            allConditionsMet = false
        }
    } else {
        println "missing"
        println "--- output of version `git --version` command ---"
        println output
        println errOutput
        println "-------------------------------------------------"
        allConditionsMet = false
    }
}

// Remark: We're using venv, which was introduced with python 3.3,
// that's why this is the baseline for python.
def checkPython() {
    String python = project.properties['python.exe.bin']
    println "Using python executable:   " + python.padRight(14) + " OK"
    print "Detecting Python version:  "
    try {
        def process = (python + " --version").execute()
        def stdOut = new StringBuilder()
        def stdErr = new StringBuilder()
        process.waitForProcessOutput(stdOut, stdErr)
        Matcher matcher = extractVersion(stdOut + stdErr)
        if (matcher.size() > 0) {
            String curVersion = matcher[0][1]
            def result = checkVersionAtLeast(curVersion, "3.7.0")
            if (!result) {
                allConditionsMet = false
            }
        } else {
            println "missing (Please install at least version 3.7.0)"
            allConditionsMet = false
        }
    } catch (Exception ignored) {
        println "missing"
        println "--- output of version `${python} --version` command ---"
        println output
        println "----------------------------------------------------"
        allConditionsMet = false
    }
}

// On Ubuntu it seems that venv is generally available, but the 'ensurepip' command fails.
// In this case we need to install the python3-venv package. Unfortunately checking the
// venv is successful in this case, so we need this slightly odd test.
def checkPythonVenv() {
    print "Detecting venv:            "
    try {
        def python = project.properties['python.exe.bin']
        def cmdArray = [python, "-Im", "ensurepip"]
        def process = cmdArray.execute()
        def stdOut = new StringBuilder()
        def stdErr = new StringBuilder()
        process.waitForProcessOutput(stdOut, stdErr)
        if (stdErr.contains("No module named")) {
            println "missing"
            println "--- output of version `python -Im \"ensurepip\"` command ---"
            println output
            println "------------------------------------------------------------"
            allConditionsMet = false
        } else {
            println "               OK"
        }
    } catch (Exception e) {
        println "missing"
        println "--- failed with exception ---"
        println e
        e.printStackTrace()
        println "----------------------------------------------------"
        allConditionsMet = false
    }
}

// When building the StreamPipes modules we need Docker.
// Not only should the docker executable be available, but also should the docker daemon be running.
def checkDocker() {
    print "Detecting Docker version:  "
    def output = new StringBuilder()
    def errOutput = new StringBuilder()
    try {
        def proc = "docker info".execute()
        proc.waitForProcessOutput(output, errOutput)
    } catch (IOException e) {
        output = ""
        errOutput.append(e)
    }
    // Check if Docker is installed at all
    def matcher1 = output =~ /Server:/
    if (matcher1.size() > 0) {
        // If it is check if the daemon is running and if the version is ok
        def matcher2 = output =~ /Server Version: (\d+\.\d+(\.\d+)?).*/
        if (matcher2.size() > 0) {
            String curVersion = matcher2[0][1]
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
        println "--- output of version `docker info` command ---"
        println output
        println errOutput
        println "-----------------------------------------------"
        allConditionsMet = false
    }
    // TODO: Implement the actual check ...
}

def checkLibPcap(String minVersion, String os, String arch) {
    print "Detecting LibPcap version: "
    try {
        // For some reason it doesn't work, if we pass this in from the outside.
        if (os == "mac") {
            // On my Intel Mac I found the libs in: "/usr/local/Cellar/libpcap/1.10.1/lib"
            // On my M1 Mac I found the libs in: "/opt/homebrew/Cellar/libpcap/1.10.1/lib"
            if (new File("/usr/local/Cellar/libpcap/1.10.1/lib").exists()) {
                System.getProperties().setProperty("jna.library.path", "/usr/local/Cellar/libpcap/1.10.1/lib");
            } else if (new File("/opt/homebrew/opt/libpcap/lib").exists()) {
                System.getProperties().setProperty("jna.library.path", "/opt/homebrew/opt/libpcap/lib");
            }
            // java.lang.UnsatisfiedLinkError: Can't load library: /Users/christoferdutz/Library/Caches/JNA/temp/jna877652535357666533.tmp
        }
        // TODO: For some reason this check doesn't work on my M1 mac ... I get unsattisfiedlinkerror from the JNA library.
        if (arch != "aarch64") {
            output = org.pcap4j.core.Pcaps.libVersion()
            String version = output - ~/^libpcap version /
            def result = checkVersionAtLeast(version, minVersion)
            if (!result) {
                //allConditionsMet = false
            }
        } else {
            println "               SKIPPED (on aarch64)"
        }
    } catch (Error e) {
        output = ""
        println "missing"
        println "--- exception ---"
        println e
        e.printStackTrace()
        println "-----------------"
        allConditionsMet = false
    }
}

def checkLibPcapHeaders() {

}

/**
 * Version extraction function/macro. It looks for occurrence of x.y or x.y.z
 * in passed input text (likely output from `program --version` command if found).
 *
 * @param input
 * @return
 */
private static Matcher extractVersion(input) {
    def matcher = input =~ /(\d+\.\d+(\.\d+)?).*/
    matcher
}

/////////////////////////////////////////////////////
// Find out which OS and arch are bring used.
/////////////////////////////////////////////////////
println "Os name:    ${System.getProperty("os.name")}"
println "Os arch:    ${System.getProperty("os.arch")}"
println "Os version: ${System.getProperty("os.version")}"
def osString = project.properties['os.classifier']
def osMatcher = osString =~ /(.*)-(.*)/
if (osMatcher.size() == 0) {
    throw new RuntimeException("Currently unsupported OS. Actual os string: $osString")
}
def os = osMatcher[0][1]
def arch = osMatcher[0][2]
println "Detected OS:   $os"
println "Detected Arch: $arch"

/////////////////////////////////////////////////////
// Find out which profiles are enabled.
/////////////////////////////////////////////////////

def cEnabled = false
def dotnetEnabled = false
def goEnabled = false
// Java is always enabled ...
def javaEnabled = false
def pythonEnabled = false
def apacheReleaseEnabled = false
def activeProfiles = session.request.activeProfiles
for (def activeProfile : activeProfiles) {
    if (activeProfile == "with-c") {
        cEnabled = true
    } else if (activeProfile == "with-dotnet") {
        dotnetEnabled = true
    } else if (activeProfile == "with-go") {
        goEnabled = true
    } else if (activeProfile == "with-java") {
        javaEnabled = true
    } else if (activeProfile == "with-python") {
        pythonEnabled = true
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

checkJavaVersion("11", null)

if (dotnetEnabled) {
    checkDotnet()
}

if (javaEnabled) {
    checkGit()
}

if (cEnabled) {
    checkGcc()
}

if (goEnabled) {
    checkLibPcapHeaders()
}

if (pythonEnabled) {
    checkPython()
    checkPythonVenv()
}

if (cEnabled) {
    // CMake requires at least maven 3.6.0
    checkMavenVersion("3.6.0", null)
}

if (apacheReleaseEnabled) {
    // TODO: Check libpcap is installed
}

if (os == "mac") {
    // The current system version from mac crashes so we assert for a version coming with brew
    checkLibPcap("1.10.1", os, arch)
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

