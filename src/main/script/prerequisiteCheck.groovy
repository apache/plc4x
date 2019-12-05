import java.util.regex.Matcher

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
        if(currentSegment < minimumSegment) {
            println current.padRight(14) + "FAILED (required " + minimum + ")"
            return false
        } else if(currentSegment > minimumSegment) {
            println current.padRight(14) + "OK"
            return true
        }
    }
    def curNotShorter = currentSegments.size() >= minimumSegments.size()
    if(curNotShorter) {
        println current.padRight(14) + " OK"
    } else {
        println current.padRight(14) + " (required " + minimum + ")"
    }
    curNotShorter
}

def checkBison() {
    print "Detecting Bison version:   "
    def output = "bison --version".execute().text
    Matcher matcher = extractVersion(output)
    if(matcher.size() > 0) {
        def curVersion = matcher[0][1]
        def result = checkVersionAtLeast(curVersion, "2.4.0")
        if(!result) {
            allConditionsMet = false
        }
    } else {
        println "missing"
        allConditionsMet = false
    }
}

def checkDotnet() {
    print "Detecting Dotnet version:  "
    def output = "dotnet --version".execute().text
    Matcher matcher = extractVersion(output)
    if(matcher.size() > 0) {
        def curVersion = matcher[0][1]
        def result = checkVersionAtLeast(curVersion, "2.0.0")
        if(!result) {
            allConditionsMet = false
        }
    } else {
        println "missing"
        allConditionsMet = false
    }
}

def checkFlex() {
    print "Detecting Flex version:    "
    def output = "flex --version".execute().text
    Matcher matcher = extractVersion(output)
    if(matcher.size() > 0) {
        def curVersion = matcher[0][1]
        def result = checkVersionAtLeast(curVersion, "2.0.0")
        if(!result) {
            allConditionsMet = false
        }
    } else {
        println "missing"
        allConditionsMet = false
    }
}

def checkGcc() {
    print "Detecting Gcc version:     "
    def output = "gcc --version".execute().text
    Matcher matcher = extractVersion(output)
    if(matcher.size() > 0) {
        def curVersion = matcher[0][1]
        def result = checkVersionAtLeast(curVersion, "1.0.0")
        if(!result) {
            allConditionsMet = false
        }
    } else {
        println "missing"
        allConditionsMet = false
    }
}

def checkGit() {
    print "Detecting Git version:     "
    def output = "git --version".execute().text
    Matcher matcher = extractVersion(output)
    if(matcher.size() > 0) {
        def curVersion = matcher[0][1]
        def result = checkVersionAtLeast(curVersion, "1.0.0")
        if(!result) {
            allConditionsMet = false
        }
    } else {
        println "missing"
        allConditionsMet = false
    }
}

def checkGpp() {
    print "Detecting G++ version:     "
    def output = "g++ --version".execute().text
    Matcher matcher = extractVersion(output)
    if(matcher.size() > 0) {
        def curVersion = matcher[0][1]
        def result = checkVersionAtLeast(curVersion, "1.0.0")
        if(!result) {
            allConditionsMet = false
        }
    } else {
        println "missing"
        allConditionsMet = false
    }
}

def checkClang() {
    print "Detecting clang version:   "
    def output = "clang --version".execute().text
    Matcher matcher = extractVersion(output)
    if(matcher.size() > 0) {
        def curVersion = matcher[0][1]
        def result = checkVersionAtLeast(curVersion, "1.0.0")
        if(!result) {
            allConditionsMet = false
        }
    } else {
        println "missing"
        allConditionsMet = false
    }
}

def checkPython() {
    print "Detecting Python version:  "
    def process = ("python --version").execute()
    def stdOut = new StringBuilder()
    def stdErr = new StringBuilder()
    process.consumeProcessOutput(stdOut, stdErr)
    process.waitForOrKill(500)
    Matcher matcher = extractVersion(stdErr)
    if(matcher.size() > 0) {
        def curVersion = matcher[0][1]
        def result = checkVersionAtLeast(curVersion, "2.7.0")
        if(!result) {
            allConditionsMet = false
        }
    } else {
        println "missing"
        // For debugging regular build failures on our build vm
        println "StdOut: " + stdOut
        println "StrErr: " + stdErr
        println "matcher size: " + matcher.size()
        for(int i = 0; i < matcher.size(); i++) {
            println "matcher[" + i + "]=" + matcher[i]
        }
        // Example for a failed python detection:
        //
        //Detecting Python version: missing
        //StdOut:
        //StrErr: Python 2.7.12
        // Example of a successful detection
        //StrErr:
        //2.7.15        OK
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
    if(!expectedFile.exists()) {
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
    def output = "openssl version".execute().text
    Matcher matcher = extractVersion(output)
    if(matcher.size() > 0) {
        def curVersion = matcher[0][1]
        def result = checkVersionAtLeast(curVersion, "1.0.0")
        if(!result) {
            allConditionsMet = false
        }
    } else {
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
if(osMatcher.size() == 0) {
    throw new RuntimeException("Currently unsupported OS")
}
def os = osMatcher[0][1]
def arch = osMatcher[0][2]
println "Detected OS:   " + os
println "Detected Arch: " + arch

/////////////////////////////////////////////////////
// Find out which profiles are enabled.
/////////////////////////////////////////////////////

println "Enabled profiles:"
def boostEnabled = false
def cppEnabled = false
def dotnetEnabled = false
def javaEnabled = true
def pythonEnabled = false
def proxiesEnabled = false
def sandboxEnabled = false
def activeProfiles = session.request.activeProfiles
for (def activeProfile : activeProfiles) {
    if(activeProfile == "with-boost") {
        boostEnabled = true
        println "boost"
    } else if(activeProfile == "with-cpp") {
        cppEnabled = true
        println "cpp"
    } else if(activeProfile == "with-dotnet") {
        dotnetEnabled = true
        println "dotnet"
    } else if(activeProfile == "with-python") {
        pythonEnabled = true
        println "python"
    } else if(activeProfile == "with-proxies") {
        proxiesEnabled = true
        println "proxies"
    } else if(activeProfile == "with-sandbox") {
        sandboxEnabled = true
        println "sandbox"
    }
}
println ""

// - Windows:
//     - Check the length of the path of the base dir as we're having issues with the length of paths being too long.
if(os == "win") {
    File pomFile = project.model.pomFile
    if(pomFile.absolutePath.length() > 100) {
        println "On Windows we encounter problems with maximum path lengths. " +
            "Please move the project to a place it has a shorter base path " +
            "and run the build again."
        allConditionsMet = false;
    }
}

if(pythonEnabled && !proxiesEnabled) {
    println "Currently the build of the python modules require the `with-proxies` profile to be enabled tpo."
    allConditionsMet = false;
}

/////////////////////////////////////////////////////
// Do the actual checks depending on the enabled
// profiles.
/////////////////////////////////////////////////////

if(proxiesEnabled) {
    checkBison()
}

if(dotnetEnabled) {
    checkDotnet()
}

if(proxiesEnabled) {
    checkFlex()
    checkOpenSSL()
}

if(proxiesEnabled || cppEnabled) {
    checkClang()
    checkGcc()
}

if(javaEnabled) {
    checkGit()
}

if(proxiesEnabled || cppEnabled) {
    checkGpp()
}

if(pythonEnabled) {
    checkPython()
}

// We only need this check, if boost is not enabled but we're enabling cpp.
if(!boostEnabled && cppEnabled) {
    checkBoost()
}

if(!allConditionsMet) {
    throw new RuntimeException("Not all conditions met, see log for details.")
}
println ""
println "All known conditions met successfully."
println ""

// Things we could possibly check:
// - DNS Providers that return a default ip on unknown host-names
// - Availability and version of LibPCAP/WinPCAP

