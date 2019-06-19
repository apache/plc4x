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
package org.apache.plc4x.plugins.codegenerator;

import org.apache.maven.artifact.Artifact;
import org.apache.maven.plugin.AbstractMojo;
import org.apache.maven.plugin.MojoExecutionException;
import org.apache.maven.plugins.annotations.LifecyclePhase;
import org.apache.maven.plugins.annotations.Mojo;
import org.apache.maven.plugins.annotations.Parameter;
import org.apache.maven.plugins.annotations.ResolutionScope;
import org.apache.maven.project.MavenProject;
import org.apache.plc4x.plugins.codegenerator.language.LanguageOutput;
import org.apache.plc4x.plugins.codegenerator.types.definitions.ComplexTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.exceptions.GenerationException;
import org.apache.plc4x.plugins.codegenerator.protocol.Protocol;

import java.io.*;
import java.net.MalformedURLException;
import java.net.URL;
import java.net.URLClassLoader;
import java.util.*;

/**
 * Generate the types, serializer and parser classes based on a DFDL shema.
 */
@Mojo(name = "generate-driver",
    defaultPhase = LifecyclePhase.GENERATE_SOURCES,
    requiresDependencyResolution = ResolutionScope.COMPILE_PLUS_RUNTIME)
public class GenerateMojo extends AbstractMojo {

    @Parameter(defaultValue = "${project}", readonly = true, required = true)
    private MavenProject project;

    /**
     * Where to generate the code.
     */
    @Parameter(defaultValue = "${project.build.directory}/generated-sources/plc4x/", required = true)
    private File outputDir;

    /**
     * The name of the protocol module that will be used to generate the driver.
     */
    @Parameter(required = true)
    private String protocolName;

    /**
     * The name of the language name that will be used to generate the driver.
     */
    @Parameter(required = true)
    private String languageName;

    public void execute()
        throws MojoExecutionException {

        // Make sure the output directory exists.
        if (!outputDir.exists()) {
            if (!outputDir.mkdirs()) {
                throw new MojoExecutionException("Could not generate output directory " + outputDir.getAbsolutePath());
            }
        }

        // Build a classloader that can access the projects classpath (read from dependencies)
        ClassLoader moduleClassloader;
        try {
            Set<Artifact> artifacts = project.getArtifacts();
            List<URL> classpathElements = new ArrayList<>(artifacts.size() + 1);
            // Add the normal class output (needed for embedded schemas)
            classpathElements.add(new File(project.getBuild().getOutputDirectory()).toURI().toURL());
            // Add all the other artifacts (no matter what scope)
            for (Artifact artifact : artifacts) {
                classpathElements.add(artifact.getFile().toURI().toURL());
            }
            moduleClassloader = new URLClassLoader(
                classpathElements.toArray(new URL[0]), GenerateMojo.class.getClassLoader());
        } catch (MalformedURLException e) {
            throw new MojoExecutionException(
                "Error creating classloader for loading message format schema from module dependencies", e);
        }

        // Load the protocol module.
        Protocol protocol = null;
        ServiceLoader<Protocol> protocols = ServiceLoader.load(Protocol.class, moduleClassloader);
        for (Protocol curProtocol : protocols) {
            if(curProtocol.getName().equalsIgnoreCase(protocolName)) {
                protocol = curProtocol;
                break;
            }
        }
        if(protocol == null) {
            throw new MojoExecutionException(
                "Unable to find protocol specification module '" + protocolName + "' on modules classpath");
        }

        // Load the language module.
        LanguageOutput language = null;
        ServiceLoader<LanguageOutput> languages = ServiceLoader.load(LanguageOutput.class, moduleClassloader);
        for (LanguageOutput curLanguage : languages) {
            if(curLanguage.getName().equalsIgnoreCase(languageName)) {
                language = curLanguage;
                break;
            }
        }
        if(language == null) {
            throw new MojoExecutionException(
                "Unable to find language output module '" + languageName + "' on modules classpath");
        }

        try {
            // Parse the type definitions.
            Map<String, ComplexTypeDefinition> types = protocol.getTypeDefinitions();

            // Generate output for the type definitions.
            language.generate(outputDir, "org.apache.plc4x." + languageName.toLowerCase() +
                "." + protocolName.toLowerCase(), types);
        } catch (GenerationException e) {
            getLog().error("Error generating sources", e);
        }

        // Add the generated sources to the project internally.
        project.addCompileSourceRoot(outputDir.getPath());
    }

}
