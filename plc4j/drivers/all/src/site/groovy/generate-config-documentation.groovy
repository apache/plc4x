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
import org.apache.plc4x.java.spi.configuration.annotations.ComplexConfigurationParameter
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter
import org.apache.plc4x.java.spi.configuration.annotations.Description
import org.apache.plc4x.java.spi.configuration.annotations.defaults.BooleanDefaultValue
import org.apache.plc4x.java.spi.configuration.annotations.defaults.DoubleDefaultValue
import org.apache.plc4x.java.spi.configuration.annotations.defaults.FloatDefaultValue
import org.apache.plc4x.java.spi.configuration.annotations.defaults.IntDefaultValue
import org.apache.plc4x.java.spi.configuration.annotations.defaults.LongDefaultValue
import org.apache.plc4x.java.spi.configuration.annotations.defaults.StringDefaultValue
import org.apache.maven.artifact.Artifact

import java.lang.reflect.Array
import java.lang.reflect.Field
import java.lang.reflect.Modifier

def static getAllFields(Class<?> type) {
    def fields = new ArrayList();
    fields.addAll(type.getDeclaredFields())
    if(type.superclass) {
        fields.addAll(getAllFields(type.superclass))
    }
    return fields
}

def static outputFields(List<Field> fields, String prefix, PrintStream printStream) {
    for (final def field in fields) {
        // Skip constants.
        if(Modifier.isFinal(field.getModifiers())) {
            continue
        }

        var name = ((prefix) ? prefix + "." : "") + field.name
        var configurationParameterAnnotation = field.annotations.find( annotation -> annotation.annotationType().name.endsWith("ConfigurationParameter") )
        if(!configurationParameterAnnotation) {
            continue;
        }
        if(configurationParameterAnnotation instanceof ComplexConfigurationParameter) {
            def parameterPrefix = ((ComplexConfigurationParameter) configurationParameterAnnotation).prefix()
            def parameterType = field.type
            def parameterFields = getAllFields(parameterType)
            outputFields(parameterFields, ((prefix) ? prefix + "." : "") + parameterPrefix, printStream)
            return
        } else {
            def parameterName = ((ConfigurationParameter) configurationParameterAnnotation).value().toString()
            if (parameterName && parameterName.length() > 0) {
                name = ((prefix) ? prefix + "." : "") + ((ConfigurationParameter) configurationParameterAnnotation).value().toString()
            }
        }
        var type = field.type.name
        if (type == "java.lang.String") {
            type = "string"
        }
        var defaultValueAnnotation = field.annotations.find { annotation -> annotation.annotationType().name.endsWith("DefaultValue") }
        var defaultValue = ""
        if (defaultValueAnnotation) {
            switch (defaultValueAnnotation.annotationType().name) {
                case "org.apache.plc4x.java.spi.configuration.annotations.defaults.BooleanDefaultValue":
                    defaultValue = ((BooleanDefaultValue) defaultValueAnnotation).value().booleanValue()
                    break
                case "org.apache.plc4x.java.spi.configuration.annotations.defaults.DoubleDefaultValue":
                    defaultValue = ((DoubleDefaultValue) defaultValueAnnotation).value().doubleValue()
                    break
                case "org.apache.plc4x.java.spi.configuration.annotations.defaults.FloatDefaultValue":
                    defaultValue = ((FloatDefaultValue) defaultValueAnnotation).value().floatValue()
                    break
                case "org.apache.plc4x.java.spi.configuration.annotations.defaults.IntDefaultValue":
                    defaultValue = ((IntDefaultValue) defaultValueAnnotation).value().intValue()
                    break
                case "org.apache.plc4x.java.spi.configuration.annotations.defaults.LongDefaultValue":
                    defaultValue = ((LongDefaultValue) defaultValueAnnotation).value().longValue()
                    break
                case "org.apache.plc4x.java.spi.configuration.annotations.defaults.StringDefaultValue":
                    defaultValue = ((StringDefaultValue) defaultValueAnnotation).value().toString()
                    break
            }
        }
        var requiredAnnotation = field.annotations.find { annotation -> annotation.annotationType().name == "org.apache.plc4x.java.spi.configuration.annotations.Required" }
        var required = (requiredAnnotation) ? "required" : ""
        var descriptionAnnotation = field.annotations.find { annotation -> annotation.annotationType().name == "org.apache.plc4x.java.spi.configuration.annotations.Description" }
        var description = ""
        if (descriptionAnnotation != null) {
            description = ((Description) descriptionAnnotation).value()
        }
        printStream.println "|`" + name + "` |" + type + " |" + defaultValue + " |" + required + " |" + description
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
    moduleClassloader = new URLClassLoader(
        classpathElements.toArray(new URL[0]), this.class.getClassLoader())
} catch (MalformedURLException e) {
    throw new Exception(
        "Error creating classloader for loading message format schema from module dependencies", e);
}

// Create a driver manager instance, that is using our custom built classloader.
def plcDriverManager = new DefaultPlcDriverManager(moduleClassloader)

// Process all driver information.
for (final def protocolCode in plcDriverManager.listProtocolCodes()) {
    def outputFile = new File(project.getBasedir(), "src/site/generated/" + protocolCode + ".adoc")
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
        "----\n" +
        "<dependency>\n" +
        "  <groupId>org.apache.plc4x</groupId>\n" +
        "  <artifactId>" + moduleName + "</artifactId>\n" +
        "  <version>{current-last-released-version}</version>\n" +
        "</dependency>\n" +
        "----"
    if(driver.defaultTransportCode.isPresent()) {
        printStream.println "|Default Transport 4+|`" + driver.defaultTransportCode.get() + "`"
    }
    printStream.println "|Supported Transports 4+|"
    for (final def transportCode in driver.supportedTransportCodes) {
        // TODO: Make it output stuff like the "default port" for UDP and TCP
        printStream.println " - `" + transportCode + "`"
    }
    printStream.println "5+|Config options:"
    def configurationType = driver.getConfigurationType()

    // Output the configuration options of the driver itself.
    def fields = getAllFields(configurationType)
    outputFields(fields, null, printStream)

    // Output the configuration options of the transports the driver supports.
    if(!driver.supportedTransportCodes.empty) {
        printStream.println "5+|Transport config options:"
        for (final def transportCode in driver.supportedTransportCodes) {
            printStream.println "5+| - `" + transportCode + "`"
            def transportConfigurationType = driver.getTransportConfigurationType(transportCode)
            if (transportConfigurationType.present) {
                def transportConfigurationFields = getAllFields(transportConfigurationType.get())
                outputFields(transportConfigurationFields, transportCode, printStream)
            }
        }
    }

    printStream.println "|==="

    printStream.close()
}