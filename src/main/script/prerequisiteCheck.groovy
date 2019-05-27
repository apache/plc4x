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
    print "Detecting Bison version:  "
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
    print "Detecting Dotnet version: "
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
    print "Detecting Flex version:   "
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
    print "Detecting Gcc version:    "
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
    print "Detecting Git version:    "
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
    print "Detecting G++ version:    "
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

def checkPython() {
    print "Detecting Python version: "
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
        allConditionsMet = false
    }
}

/**
 * Version extraction function/macro. It looks for occurance of x.y or x.y.z
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
def cppEnabled = false
def dotnetEnabled = false
def javaEnabled = false
def pythonEnabled = false
def proxiesEnabled = false
def sandboxEnabled = false
def activeProfiles = session.request.activeProfiles
for (def activeProfile : activeProfiles) {
    if(activeProfile == "with-cpp") {
        cppEnabled = true
        println "cpp"
    } else if(activeProfile == "with-dotnet") {
        dotnetEnabled = true
        println "dotnet"
    } else if(activeProfile == "with-java") {
        javaEnabled = true
        println "java"
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
}

if(proxiesEnabled || cppEnabled) {
    checkGcc()
}

if(javaEnabled) {
    checkGit()
}

if(proxiesEnabled || cppEnabled) {
    checkGpp()
}

// TODO: Doesn't work yet
if(pythonEnabled) {
    checkPython()
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

