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

var cmakeRoot = new File(project.properties["cmake.root"])
if(!cmakeRoot.exists()) {
    exit(0)
}

var cmakeExecutable = new File(cmakeRoot, "cmake")
if(!cmakeRoot.exists()) {
    exit(0)
}

// Rename the original 'cmake' file to 'cmakeOrig'
var cmakeRenamedExecutable = new File(cmakeRoot, "cmakeOrig")
cmakeExecutable.renameTo(cmakeRenamedExecutable)

// Create a little bash-script that calls the build-wrapper and then calls the original cmake file
cmakeExecutable.write("#!/bin/bash\n")
cmakeExecutable.append("echo \"Arguments: \$@\"\n")
cmakeExecutable.append("if [[ \"\$1\" == \"--build\" ]]; then\n")
cmakeExecutable.append("  " + cmakeRoot.absolutePath + "/../../build-wrapper-linux-x86/build-wrapper-linux-x86-64 --out-dir " + project.properties["sonar.cfamily.build-wrapper-output"] + " " + cmakeRoot.absolutePath + "/cmakeOrig \"\$@\"\n")
cmakeExecutable.append("else\n")
cmakeExecutable.append("  " + cmakeRoot.absolutePath + "/cmakeOrig \"\$@\"\n")
cmakeExecutable.append("fi\n")
// Make the script executable
("chmod +x " + cmakeExecutable.absolutePath).execute()
