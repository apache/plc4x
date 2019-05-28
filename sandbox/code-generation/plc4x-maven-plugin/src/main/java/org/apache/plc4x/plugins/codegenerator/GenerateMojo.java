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

import freemarker.cache.ClassTemplateLoader;
import freemarker.core.ParseException;
import freemarker.template.*;
import org.apache.maven.artifact.Artifact;
import org.apache.maven.plugin.AbstractMojo;
import org.apache.maven.plugin.MojoExecutionException;
import org.apache.maven.plugins.annotations.LifecyclePhase;
import org.apache.maven.plugins.annotations.Mojo;
import org.apache.maven.plugins.annotations.Parameter;
import org.apache.maven.plugins.annotations.ResolutionScope;
import org.apache.maven.project.MavenProject;
import org.apache.plc4x.language.LanguageTemplate;
import org.apache.plc4x.plugins.codegenerator.model.types.ComplexType;
import org.apache.plc4x.plugins.codegenerator.parser.MessageFormatParser;
import org.apache.plc4x.protocol.Protocol;

import java.io.*;
import java.net.MalformedURLException;
import java.net.URL;
import java.net.URLClassLoader;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
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
     * The name of the language template that will be used to generate the driver.
     */
    @Parameter(required = true)
    private String languageTemplateName;

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

        // Load the language template.
        LanguageTemplate languageTemplate = null;
        ServiceLoader<LanguageTemplate> languageTemplates = ServiceLoader.load(LanguageTemplate.class, moduleClassloader);
        for (LanguageTemplate curLanguageTemplate : languageTemplates) {
            if(curLanguageTemplate.getName().equalsIgnoreCase(languageTemplateName)) {
                languageTemplate = curLanguageTemplate;
                break;
            }
        }
        if(languageTemplate == null) {
            throw new MojoExecutionException(
                "Unable to find language template '" + languageTemplateName + "' on modules classpath");
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
                "Unable to find protocol module '" + protocolName + "' on modules classpath");
        }

        // Try loading the file directly (Without classloader)
        InputStream schemaInputStream = protocol.getMessageFormatSchema();
        if(schemaInputStream == null) {
            throw new MojoExecutionException("Error loading message-format schema for protocol '" + protocolName + "'");
        }

        Map<String, ComplexType> types = new MessageFormatParser().parse(schemaInputStream);
        try {
            // Configure the Freemarker template engine
            Configuration freemarkerConfiguration = getFreemarkerConfiguration();
            freemarkerConfiguration.setClassLoaderForTemplateLoading(moduleClassloader, "/");

            // Initialize all templates
            List<Template> templateList = languageTemplate.getTemplates(freemarkerConfiguration);

            // Iterate over the types and have content generated for each one
            for (Map.Entry<String, ComplexType> typeEntry : types.entrySet()) {
                // Prepare a new generation context
                Map<String, Object> typeContext = new HashMap<>();
                typeContext.put("packageName", "org.apache.plc4x." +
                    languageTemplateName.toLowerCase() + "." + protocolName.toLowerCase());
                typeContext.put("typeName", typeEntry.getKey());
                typeContext.put("type", typeEntry.getValue());
                typeContext.put("helper", languageTemplate.getHelper());

                for(Template template : templateList) {
                    // Create a variable size output location where the template can generate it's content to
                    ByteArrayOutputStream output = new ByteArrayOutputStream();

                    // Have Freemarker generate the output
                    template.process(typeContext, new OutputStreamWriter(output));

                    // Create the means to read in the generated output back in again
                    BufferedReader input = new BufferedReader(
                        new InputStreamReader(new ByteArrayInputStream(output.toByteArray())));

                    // Extract the output path from the first line of the generated content
                    String outputFileName = input.readLine();
                    File outputFile = new File(outputDir, outputFileName);

                    // Create any missing directories
                    if(!outputFile.getParentFile().exists()) {
                        if(!outputFile.getParentFile().mkdirs()) {
                            throw new MojoExecutionException(
                                "Unable to create output directory " + outputFile.getParent());
                        }
                    }

                    // Output the rest to that file
                    BufferedWriter outputFileWriter = Files.newBufferedWriter(
                        outputFile.toPath(), StandardCharsets.UTF_8);
                    String line;
                    while ((line = input.readLine()) != null) {
                        outputFileWriter.write(line);
                        outputFileWriter.newLine();
                    }
                    outputFileWriter.flush();
                }
                getLog().info("Generating type " + typeEntry.getKey());
            }
        } catch (TemplateNotFoundException | TemplateException | MalformedTemplateNameException | ParseException e) {
            throw new MojoExecutionException("Error resolving template", e);
        } catch (IOException e) {
            throw new MojoExecutionException("Error generating sources", e);
        }
    }

    private Configuration getFreemarkerConfiguration() throws MojoExecutionException {
        Configuration configuration = new Configuration(Configuration.VERSION_2_3_28);
        configuration.setDefaultEncoding("UTF-8");
        configuration.setTemplateExceptionHandler(TemplateExceptionHandler.RETHROW_HANDLER);
        configuration.setLogTemplateExceptions(false);
        configuration.setWrapUncheckedExceptions(true);
        try {
            configuration.setDirectoryForTemplateLoading(new File("/"));
        } catch (IOException e) {
            throw new MojoExecutionException("Error initializing freemarker configuration");
        }
        return configuration;
    }

}
