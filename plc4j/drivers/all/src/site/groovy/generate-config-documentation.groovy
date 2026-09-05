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

import org.apache.plc4x.java.DefaultPlcDriverManager
import org.apache.maven.artifact.Artifact
import org.apache.plc4x.java.api.metadata.Option

/**
 * Configuration descriptions and default values are plain text coming from the
 * "@ConfigurationParameter" annotations - they are not asciidoc. Two characters in them would
 * otherwise be interpreted by asciidoctor and silently corrupt the generated table:
 *   "{" starts an attribute reference, so a description mentioning an address syntax such as
 *       "{IndexGroup}/{IndexOffset}" is reported as a missing attribute and left unresolved,
 *   "|" starts a new table cell, so a description containing one would push the remaining cells
 *       out of the row and asciidoctor would drop them.
 * Escaping both here keeps every description safe, whatever an annotation happens to contain.
 * Backslashes are deliberately left alone: asciidoctor only treats one as an escape when it
 * precedes a construct it would otherwise substitute, and doubling them would corrupt the
 * Windows device paths that appear in some descriptions.
 */
def static escapeAsciidoc(String text) {
    text.replace('{', '\\{').replace('|', '\\|')
}

def static outputOptions(List<Option> options, String prefix, PrintStream printStream) {
    options.each {option->
        def name = prefix?"$prefix.$option.key":option.key
        // Convert java line-breaks into asciidoctor line-breaks.
        def description = escapeAsciidoc(option.description).replaceAll('\n', " +\n")
        option.since.ifPresent {
            description += " +\n*Since: " + option.since.get() + "*"
        }
        def defaultValue = option.defaultValue.map { escapeAsciidoc(it as String) }.orElse(' ')
        printStream.println "|`$name` |$option.type |${defaultValue}|${option.required?'required':''} |$description"
    }
}

// Build a classloader that can access the projects classpath (read from dependencies)
ClassLoader moduleClassloader
try {
    Set<Artifact> artifacts = project.getArtifacts()
    List<URL> classpathElements = new ArrayList<>(artifacts.size() + 1)
    // Add the normal class output (needed for embedded schemas)
    classpathElements.add(new File(project.getBuild().getOutputDirectory()).toURI().toURL())
    // Add all the other artifacts (no matter what scope)
    for (Artifact artifact : artifacts) {
        classpathElements.add(artifact.getFile().toURI().toURL())
    }
    moduleClassloader = new URLClassLoader(classpathElements.toArray(new URL[0]) as URL[], this.class.getClassLoader())
} catch (MalformedURLException e) {
    throw new Exception(
        "Error creating classloader for loading message format schema from module dependencies", e);
}
Thread.currentThread().setContextClassLoader(moduleClassloader)

// Create a driver manager instance, that is using our custom built classloader.
def plcDriverManager = new DefaultPlcDriverManager(moduleClassloader)

// Process all driver information.
for (final def protocolCode in plcDriverManager.getProtocolCodes()) {
    def outputFile = new File(project.getBasedir(), "../../../website/asciidoc/modules/users/partials/" + protocolCode + ".adoc")
    // In order to re-generate this file, make sure it doesn't exist.
    if(outputFile.exists()) {
        outputFile.delete()
    }
    // Make sure all directories exist.
    if(!outputFile.parentFile.exists()) {
        outputFile.parentFile.mkdirs()
    }
    def printStream = new PrintStream(outputFile)

    printStream.println("//\n" +
        "//  Licensed to the Apache Software Foundation (ASF) under one or more\n" +
        "//  contributor license agreements.  See the NOTICE file distributed with\n" +
        "//  this work for additional information regarding copyright ownership.\n" +
        "//  The ASF licenses this file to You under the Apache License, Version 2.0\n" +
        "//  (the \"License\"); you may not use this file except in compliance with\n" +
        "//  the License.  You may obtain a copy of the License at\n" +
        "//\n" +
        "//      https://www.apache.org/licenses/LICENSE-2.0\n" +
        "//\n" +
        "//  Unless required by applicable law or agreed to in writing, software\n" +
        "//  distributed under the License is distributed on an \"AS IS\" BASIS,\n" +
        "//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n" +
        "//  See the License for the specific language governing permissions and\n" +
        "//  limitations under the License.\n" +
        "//\n\n" +
        "// Code generated by code-generation. DO NOT EDIT.\n")

    def driver = plcDriverManager.getDriver(protocolCode)
    printStream.println "[cols=\"2,2a,2a,2a,4a\"]\n" +
        "|===\n" +
        "|Name |Type |Default Value |Required |Description"
    printStream.println "|Name 4+|" + driver.protocolName
    printStream.println "|Code 4+|`" + protocolCode + "`"
    printStream.println "|Maven Dependency 4+|"
    // Find out which jar the driver was loaded from.
    def resource = moduleClassloader.getResource(driver.class.name.replace('.', '/') + ".class")
    def uriString = resource.toExternalForm()
    var moduleName = uriString.substring(uriString.indexOf("plc4j-driver-"), uriString.indexOf("/", uriString.indexOf("plc4j-driver-")))
    printStream.println "\n" +
        "[subs=attributes+]\n" +
        "----\n" +
        "<dependency>\n" +
        "  <groupId>org.apache.plc4x</groupId>\n" +
        "  <artifactId>" + moduleName + "</artifactId>\n" +
        "  <version>{page-component-version}</version>\n" +
        "</dependency>\n" +
        "----"
    if(driver.metadata.defaultTransportCode.isPresent()) {
        printStream.println "|Default Transport 4+|`" + driver.metadata.defaultTransportCode.get() + "`"
    }
    // Filter out the "test" transport - that one is only intended for unit-tests.
    def supportedTransportCodes = driver.metadata.supportedTransportCodes.findAll { it != "test" }
    printStream.println "|Supported Transports 4+|"
    for (final def transportCode in supportedTransportCodes) {
        // TODO: Make it output stuff like the "default port" for UDP and TCP
        printStream.println " - `" + transportCode + "`"
    }
    printStream.println "5+|Config options:"

    // Output the configuration options of the driver itself.
    driver.metadata.protocolConfigurationOptionMetadata.ifPresent { outputOptions(it.options, null, printStream) }

    // Output the configuration options of the transports the driver supports.
    if(!supportedTransportCodes.empty) {
        printStream.println "5+|Transport config options:"
        for (final def transportCode in supportedTransportCodes) {
            printStream.println "5+|\n+++\n" +
                "<h4>$transportCode</h4>\n" +
                "+++"
            driver.metadata.getTransportConfigurationOptionMetadata(transportCode).ifPresent {
                outputOptions(it.options, transportCode, printStream)
            }
        }
    }

    printStream.println "|==="

    printStream.close()
}

// Process all transport information.
//
// The transport pages used to carry hand-maintained tables, and they had drifted far enough to be
// actively misleading: three of them named a transport code that does not exist ("raw" for
// "raw-socket", "pcap" for "pcap-replay", "socketcan" for "can-socketcan"), one named a maven
// artifact that does not exist, one had the wrong transport's name copy-pasted into it, and four
// claimed "Options: none" for transports that have between three and sixteen. They are generated
// from the transports' own metadata now, for the same reason the driver tables are.
//
// The options are read through the same code path the driver tables use, by handing DriverBase a
// transport's configuration class as if it were a protocol configuration. That keeps "read the
// @ConfigurationParameter annotations off a configuration class" a single implementation, so the
// two kinds of table can never disagree about a type, a default or a "Since:" marker.
//
// What is documented here is the transport's *own* configuration class. A driver may extend it -
// the ADS driver adds options to "tcp" - and those additions are driver-specific, so they belong
// in that driver's table, which already lists them under "Transport config options".
def optionsOf(Class<?> configurationClass) {
    def probe = new org.apache.plc4x.java.spi.drivers.DriverBase() {
        String getProtocolCode() { "doc-probe" }
        String getProtocolName() { "Documentation Probe" }
        protected Class getConfigurationClass() { configurationClass }
        protected org.apache.plc4x.java.spi.drivers.ConnectionBase getConnection(
            org.apache.plc4x.java.spi.config.Configuration configuration,
            org.apache.plc4x.java.spi.transports.api.TransportInstance transportInstance,
            org.apache.plc4x.java.utils.auditlog.api.AuditLog auditLog) {
            // Never called: only the metadata of this throwaway driver is ever read.
            null
        }
    }
    probe.metadata.protocolConfigurationOptionMetadata.map { it.options }.orElse([])
}

def transports = ServiceLoader.load(
    Class.forName("org.apache.plc4x.java.spi.transports.api.Transport", true, moduleClassloader),
    moduleClassloader).toList().sort { it.transportCode }

for (final def transport in transports) {
    def transportCode = transport.transportCode
    // The "test" transport is an in-memory stand-in used by PLC4X's own unit tests; the driver
    // tables above filter it out of their supported-transports lists for that reason, and it gets
    // no page here either.
    if (transportCode == "test") {
        continue
    }
    def outputFile = new File(project.getBasedir(), "../../../website/asciidoc/modules/users/partials/transport-" + transportCode + ".adoc")
    if (outputFile.exists()) {
        outputFile.delete()
    }
    if (!outputFile.parentFile.exists()) {
        outputFile.parentFile.mkdirs()
    }
    def printStream = new PrintStream(outputFile)

    printStream.println("//\n" +
        "//  Licensed to the Apache Software Foundation (ASF) under one or more\n" +
        "//  contributor license agreements.  See the NOTICE file distributed with\n" +
        "//  this work for additional information regarding copyright ownership.\n" +
        "//  The ASF licenses this file to You under the Apache License, Version 2.0\n" +
        "//  (the \"License\"); you may not use this file except in compliance with\n" +
        "//  the License.  You may obtain a copy of the License at\n" +
        "//\n" +
        "//      https://www.apache.org/licenses/LICENSE-2.0\n" +
        "//\n" +
        "//  Unless required by applicable law or agreed to in writing, software\n" +
        "//  distributed under the License is distributed on an \"AS IS\" BASIS,\n" +
        "//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n" +
        "//  See the License for the specific language governing permissions and\n" +
        "//  limitations under the License.\n" +
        "//\n\n" +
        "// Code generated by code-generation. DO NOT EDIT.\n")

    printStream.println "[cols=\"2,2a,2a,2a,4a\"]\n" +
        "|===\n" +
        "|Name |Type |Default Value |Required |Description"
    printStream.println "|Name 4+|" + transport.transportName
    printStream.println "|Code 4+|`" + transportCode + "`"
    printStream.println "|Maven Dependency 4+|"
    // Find out which jar the transport was loaded from. Several transport codes can share one
    // artifact ("tls" and "tls-psk" both ship in plc4j-transports-tls), so the artifact name is
    // read off the jar rather than derived from the code.
    def resource = moduleClassloader.getResource(transport.class.name.replace('.', '/') + ".class")
    def uriString = resource.toExternalForm()
    def markerIndex = uriString.indexOf("plc4j-transports-")
    if (markerIndex < 0) {
        throw new IllegalStateException("Unable to tell which artifact transport '" + transportCode +
            "' was loaded from (" + uriString + "). The transport tables are generated from the " +
            "installed jars - run this against a repository build, not a reactor build.")
    }
    def moduleName = uriString.substring(markerIndex, uriString.indexOf("/", markerIndex))
    printStream.println "\n" +
        "[subs=attributes+]\n" +
        "----\n" +
        "<dependency>\n" +
        "  <groupId>org.apache.plc4x</groupId>\n" +
        "  <artifactId>" + moduleName + "</artifactId>\n" +
        "  <version>{page-component-version}</version>\n" +
        "</dependency>\n" +
        "----"

    def configurationClass = transport.transportConfigType
    def options = (configurationClass != null) ? optionsOf(configurationClass) : []
    if (options.empty) {
        printStream.println "|Options 4+|*none*"
    } else {
        printStream.println "5+|Config options:"
        // Transport options are always addressed with the transport code as their prefix, which is
        // what makes "tcp.connect-timeout" and "cotp.tpdu-size" distinguishable in one connection
        // string - so the table shows them the way they have to be typed.
        outputOptions(options, transportCode, printStream)
    }

    printStream.println "|==="

    printStream.close()
}
