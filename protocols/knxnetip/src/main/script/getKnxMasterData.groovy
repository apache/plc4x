/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
import java.nio.file.Files
import java.nio.file.StandardCopyOption

// Make sure the cache directory exists in the used maven local repo
def localRepoBaseDir = session.getLocalRepository().getBasedir()
def cacheDir = new File(localRepoBaseDir, ".cache/knx-masterdata")
if (!cacheDir.exists()) {
    cacheDir.mkdirs()
}

// Check if a previous version exists and check if we need to re-download
// If the file is less than 24h old, we won't re-download in order to avoid
// being banned on the KNX server.
def knxMasterDataFile = new File(cacheDir, "knx-mater-data.xml")
def update = true
if (knxMasterDataFile.exists()) {
    // If the last update was less than 24h before, don't update it again.
    if (knxMasterDataFile.lastModified() > (new Date().getTime() - 86400000)) {
        update = false
    }
}

// If we need to update the master-data
if (update) {
    try {
        InputStream inputStream = new URL("https://update.knx.org/data/XML/project-20/knx_master.xml").openStream()
        Files.copy(inputStream, knxMasterDataFile.toPath(), StandardCopyOption.REPLACE_EXISTING)
        println "Successfully updated knx-master-data.xml"
    } catch(Exception e) {
        println "Got an error updating knx-master-data.xml. Intentionally not failing the build as we might just be offline: " + e.getMessage()
    }
} else {
    println "Skipped updating knx-master-data.xml as it's fresh enough"
}

// Copy the knx-master-data to the current target directory
def targetDir = new File(project.getBasedir(),"target/downloads")
if (!targetDir.exists()) {
    targetDir.mkdirs()
}
def targetFile = new File(targetDir, "knx-master-data.xml")
Files.copy(knxMasterDataFile.toPath(), targetFile.toPath(), StandardCopyOption.REPLACE_EXISTING)
