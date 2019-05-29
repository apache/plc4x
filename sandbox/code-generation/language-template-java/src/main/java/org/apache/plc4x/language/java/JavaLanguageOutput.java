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

package org.apache.plc4x.language.java;

import freemarker.cache.ClassTemplateLoader;
import freemarker.core.ParseException;
import freemarker.template.*;
import org.apache.plc4x.language.LanguageOutput;
import org.apache.plc4x.language.exceptions.GenerationException;
import org.apache.plc4x.language.definitions.ComplexTypeDefinition;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.*;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class JavaLanguageOutput implements LanguageOutput {

    private static final Logger LOGGER = LoggerFactory.getLogger(JavaLanguageOutput.class);

    @Override
    public String getName() {
        return "Java";
    }

    @Override
    public void generate(File outputDir, String packageName, Map<String, ComplexTypeDefinition> types)
        throws GenerationException {

        try {
            // Configure the Freemarker template engine
            Configuration freemarkerConfiguration = getFreemarkerConfiguration();

            // Initialize all templates
            List<Template> templateList = getTemplates(freemarkerConfiguration);

            // Iterate over the types and have content generated for each one
            for (Map.Entry<String, ComplexTypeDefinition> typeEntry : types.entrySet()) {
                // Prepare a new generation context
                Map<String, Object> typeContext = new HashMap<>();
                typeContext.put("packageName", packageName);
                typeContext.put("typeName", typeEntry.getKey());
                typeContext.put("type", typeEntry.getValue());
                typeContext.put("helper", getHelper());

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
                            throw new GenerationException(
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
                LOGGER.info("Generating type " + typeEntry.getKey());
            }
        } catch (TemplateNotFoundException | TemplateException | MalformedTemplateNameException | ParseException e) {
            throw new GenerationException("Error resolving template", e);
        } catch (IOException e) {
            throw new GenerationException("Error generating sources", e);
        }
    }

    private Configuration getFreemarkerConfiguration() throws IOException {
        Configuration configuration = new Configuration(Configuration.VERSION_2_3_28);
        configuration.setDefaultEncoding("UTF-8");
        configuration.setTemplateExceptionHandler(TemplateExceptionHandler.RETHROW_HANDLER);
        configuration.setLogTemplateExceptions(false);
        configuration.setWrapUncheckedExceptions(true);
        configuration.setDirectoryForTemplateLoading(new File("/"));
        return configuration;
    }

    private List<Template> getTemplates(Configuration freemarkerConfiguration) throws IOException {
        ClassTemplateLoader classTemplateLoader = new ClassTemplateLoader(JavaLanguageOutput.class, "/");
        freemarkerConfiguration.setTemplateLoader(classTemplateLoader);
        return Collections.singletonList(freemarkerConfiguration.getTemplate("templates/java/pojo-template.ftlh"));
    }

    private JavaLanguageTemplateHelper getHelper() {
        return new JavaLanguageTemplateHelper();
    }

}
