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


import org.apache.maven.execution.MavenSession
import org.apache.maven.plugin.MojoExecution
import org.apache.maven.plugin.descriptor.PluginDescriptor
import org.apache.maven.project.MavenProject
import org.apache.maven.settings.Settings
import org.codehaus.gmaven.plugin.FailClosureTarget
import org.codehaus.gmaven.plugin.util.ContainerHelper
import org.slf4j.Logger

import java.nio.file.Files
import java.nio.file.StandardCopyOption
import java.util.zip.GZIPInputStream

def project = project as MavenProject
def basedir = basedir as File
def properties = properties as Properties
def ant = ant as AntBuilder
def fail = fail as FailClosureTarget
def log = log as Logger

def container = container as ContainerHelper
def plugin = plugin as PluginDescriptor
def pluginContext = pluginContext as Map<String, String>
def mojo = mojo as MojoExecution
def session = session as MavenSession
def settings = settings as Settings

////////
//

def download(String url, String localFileName) {
    // Make sure the cache directory exists in the used maven local repo
    def localRepoBaseDir = session.localRepository.basedir
    def cacheDir = new File(localRepoBaseDir, ".cache/opcua-cache")
    if (!cacheDir.exists()) {
        cacheDir.mkdirs()
    }

    def file = new File(cacheDir, localFileName)
    def update = true
    if (file.exists()) {
        // If the last update was less than 24h before, don't update it again.
        if (file.lastModified() > (new Date().getTime() - 86400000)) {
            update = false
        }
    }
    // TODO: temp
    update = true

    if (update) {
        try {
            URLConnection connection = new URL(url).openConnection();
            InputStream input = connection.getInputStream();
            if ("gzip".equals(connection.contentEncoding)) {
                log.info("handling sudden gzip")
                input = new GZIPInputStream(input);
            }
            Files.copy(input, file.toPath(), StandardCopyOption.REPLACE_EXISTING)
            log.info("Successfully updated {}", localFileName)
        } catch (Exception e) {
            log.info("Got an error updating {}. Intentionally not failing the build as we might just be offline", localFileName, e)
        }
    } else {
        log.info("Skipped updating {} as it's fresh enough", localFileName)
    }

    def targetDir = new File(project.getBasedir(), "target/downloads")
    if (!targetDir.exists()) {
        targetDir.mkdirs()
    }
    def targetFile = new File(targetDir, localFileName)
    Files.copy(file.toPath(), targetFile.toPath(), StandardCopyOption.REPLACE_EXISTING)
}

download("https://github.com/OPCFoundation/UA-Nodeset/raw/latest/Schema/Opc.Ua.Types.bsd", "Opc.Ua.Types.bsd")
download("https://github.com/OPCFoundation/UA-Nodeset/raw/latest/Schema/StatusCode.csv", "StatusCode.csv")
download("https://github.com/OPCFoundation/UA-Nodeset/raw/latest/Schema/Opc.Ua.NodeSet2.Services.xml", "Opc.Ua.NodeSet2.Services.xml")
download("https://github.com/OPCFoundation/UA-Nodeset/raw/latest/Schema/NodeIds.csv", "Opc.Ua.NodeIds.Services.csv")
download("https://github.com/OPCFoundation/UA-Nodeset/raw/latest/Schema/AttributeIds.csv", "AttributeIds.csv")


def servicesFile = new File(project.getBasedir(), "target/downloads/Opc.Ua.NodeIds.Services.csv")
def servicesFileTmp = new File(project.getBasedir(), "target/downloads/Opc.Ua.NodeIds.Services.csv.tmp")

def pattern = ~"^([A-Z][a-z]*)(.*),Variable\$"

servicesFileTmp.withWriter { writer ->
    servicesFile.withReader { reader ->
        while ((line = reader.readLine()) != null) {
            String line = line
            if (line.endsWith(",Variable")) {
                line = line.replaceFirst(pattern, "\$1\$2,Variable\$1")
            }
            writer.write("$line\n")
        }
    }
}

log.info("overwriting services file {}", servicesFile)
Files.copy(servicesFileTmp.toPath(), servicesFile.toPath(), StandardCopyOption.REPLACE_EXISTING)

servicesFileTmp.delete()
