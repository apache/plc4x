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

import freemarker.core.ParseException;
import freemarker.template.*;
import org.apache.maven.artifact.DependencyResolutionRequiredException;
import org.apache.maven.plugin.AbstractMojo;
import org.apache.maven.plugin.MojoExecutionException;
import org.apache.maven.plugins.annotations.LifecyclePhase;
import org.apache.maven.plugins.annotations.Mojo;
import org.apache.maven.plugins.annotations.Parameter;
import org.apache.maven.plugins.annotations.ResolutionScope;
import org.apache.maven.project.MavenProject;
import org.dom4j.*;
import org.dom4j.io.SAXReader;

import java.io.*;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.util.*;

/**
 * Generate the types, serializer and parser classes based on a DFDL shema.
 */
@Mojo(name = "generate-driver",
    defaultPhase = LifecyclePhase.GENERATE_SOURCES,
    requiresDependencyResolution = ResolutionScope.RUNTIME)
public class GenerateMojo extends AbstractMojo {

    private static final Namespace xsNamespace = new Namespace("xs", "http://www.w3.org/2001/XMLSchema");
    private static final QName complexType = new QName("complexType", xsNamespace);


    @Parameter(defaultValue = "${project}", readonly = true, required = true)
    private MavenProject project;

    /**
     * Where to generate the code.
     */
    @Parameter(defaultValue = "${project.build.directory}/generated-sources/plc4x/", required = true)
    private File outputDir;

    /**
     * The path to the DFDL schema file in the modules classpath (Ideally as runtime dependency)
     */
    @Parameter(required = true)
    private String dfdlSchema;

    @Parameter(required = true)
    private List<String> templates;

    @Parameter(required = true)
    private String packageName;

    public void execute()
        throws MojoExecutionException {

        // Make sure the output directory exists.
        if (!outputDir.exists()) {
            if (!outputDir.mkdirs()) {
                throw new MojoExecutionException("Could not generate output directory " + outputDir.getAbsolutePath());
            }
        }

        try {
            // Configure the Freemarker template engine
            Configuration freemarkerConfiguration = getFreemarkerConfiguration();

            // Initialize all templates
            List<Template> templateList = new ArrayList<>(this.templates.size());
            for(String templateName : templates) {
                templateList.add(freemarkerConfiguration.getTemplate(templateName));
            }

            // Try to find a schema file in one of the following locations:
            // - The absolute path the dfdlSchema property references
            // - Any directories in the classpath (usually target/classes)
            // - Any jar in the classpath
            File schemaFile = new File(dfdlSchema);
            if(!schemaFile.exists()) {
                for (String runtimeClasspathElement : project.getRuntimeClasspathElements()) {
                    File classPathElement = new File(runtimeClasspathElement);
                    if (classPathElement.exists()) {
                        // This is usually the "target/classes" directory.
                        if (classPathElement.isDirectory()) {
                            schemaFile = new File(classPathElement, dfdlSchema);
                            if (schemaFile.exists()) {
                                break;
                            }
                        }
                        // These are usually jar files.
                        else {
                            // TODO: Implement ...
                        }
                    }
                }
            }

            if(!schemaFile.exists()) {
                throw new MojoExecutionException("Unable to find 'dfdl-schema' at" + dfdlSchema);
            }

            // Load the DFDL schema file
            Document dfdlSchema = parseDFDLSchema(schemaFile);

            // Get the list of main types in the schema file (complexType and root level)
            Iterator<Element> types = getMainTypes(dfdlSchema);

            // Iterate over the types and have content generated for each one
            while (types.hasNext()){
                Element typeElement = types.next();
                String typeName = typeElement.attributeValue("name");

                // Prepare a new generation context
                Map<String, Object> typeContext = new HashMap<>();
                typeContext.put("typeName", typeName);
                typeContext.put("typeElement", typeElement);
                typeContext.put("packageName", packageName);

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
                }
                getLog().info("Generating type " + typeName);
            }
        } catch (DependencyResolutionRequiredException e) {
            throw new MojoExecutionException("Error resolving dependencies", e);
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

    private Document parseDFDLSchema(File schemaFile) throws MojoExecutionException {
        try {
            SAXReader reader = new SAXReader();
            return reader.read(schemaFile);
        } catch (DocumentException e) {
            throw new MojoExecutionException("Unable to parse DFDL schema at " + schemaFile.getAbsolutePath(), e);
        }
    }

    private Iterator<Element> getMainTypes(Document dfdlSchema) {
        Element rootElement = dfdlSchema.getRootElement();
        return rootElement.elementIterator(complexType);
    }

}
