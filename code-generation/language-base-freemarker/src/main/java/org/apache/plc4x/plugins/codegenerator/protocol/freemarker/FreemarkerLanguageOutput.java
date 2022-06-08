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
package org.apache.plc4x.plugins.codegenerator.protocol.freemarker;

import freemarker.cache.ClassTemplateLoader;
import freemarker.template.*;
import org.apache.plc4x.plugins.codegenerator.language.LanguageOutput;
import org.apache.plc4x.plugins.codegenerator.types.definitions.EnumTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.definitions.DataIoTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.exceptions.GenerationException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.*;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public abstract class FreemarkerLanguageOutput implements LanguageOutput {

    private static final Logger LOGGER = LoggerFactory.getLogger(FreemarkerLanguageOutput.class);

    @Override
    public void generate(File outputDir, String version, String languageName, String protocolName, String outputFlavor, Map<String, TypeDefinition> types,
                         Map<String, String> options)
        throws GenerationException {

        // Configure the Freemarker template engine
        Configuration freemarkerConfiguration = getFreemarkerConfiguration();

        ClassTemplateLoader classTemplateLoader = new ClassTemplateLoader(FreemarkerLanguageOutput.class, "/");
        freemarkerConfiguration.setTemplateLoader(classTemplateLoader);

        // Initialize all templates
        List<Template> specTemplates;
        List<Template> complexTypesTemplateList;
        List<Template> enumTypesTemplateList;
        List<Template> dataIoTemplateList;
        try {
            specTemplates = getSpecTemplates(freemarkerConfiguration);
            complexTypesTemplateList = getComplexTypeTemplates(freemarkerConfiguration);
            enumTypesTemplateList = getEnumTypeTemplates(freemarkerConfiguration);
            dataIoTemplateList = getDataIoTemplates(freemarkerConfiguration);
        } catch (IOException e) {
            throw new GenerationException("Error getting template", e);
        }

        // Generate output that's global for the entire mspec
        if (!specTemplates.isEmpty()) {
            Map<String, Object> typeContext = new HashMap<>();
            typeContext.put("languageName", languageName);
            typeContext.put("protocolName", protocolName);
            typeContext.put("outputFlavor", outputFlavor);
            typeContext.put("version", version);
            typeContext.put("helper", getHelper(null, protocolName, outputFlavor, types, options));
            typeContext.put("tracer", Tracer.start("global"));
            typeContext.putAll(options);

            for (Template template : specTemplates) {
                try {
                    renderTemplate(outputDir, template, typeContext);
                } catch (IOException | TemplateException e) {
                    throw new GenerationException("Error generating global protocol output.", e);
                }
            }
        }

        // Iterate over the types and have content generated for each one
        for (Map.Entry<String, TypeDefinition> typeEntry : types.entrySet()) {
            // Prepare a new generation context
            Map<String, Object> typeContext = new HashMap<>();
            typeContext.put("languageName", languageName);
            typeContext.put("protocolName", protocolName);
            typeContext.put("outputFlavor", outputFlavor);
            typeContext.put("typeName", typeEntry.getKey());
            typeContext.put("type", typeEntry.getValue());
            typeContext.put("helper", getHelper(typeEntry.getValue(), protocolName, outputFlavor, types, options));
            typeContext.put("tracer", Tracer.start("types"));

            // Depending on the type, get the corresponding list of templates.
            List<Template> templateList;
            if (typeEntry.getValue() instanceof EnumTypeDefinition) {
                templateList = enumTypesTemplateList;
            } else if (typeEntry.getValue() instanceof DataIoTypeDefinition) {
                templateList = dataIoTemplateList;
            } else {
                templateList = complexTypesTemplateList;
            }

            // Generate the output for the given type.
            LOGGER.info("Generating type {}", typeEntry.getKey());
            for (Template template : templateList) {
                LOGGER.debug("Applying template {}", template.getName());
                try {
                    renderTemplate(outputDir, template, typeContext);
                } catch (IOException | TemplateException e) {
                    throw new GenerationException(
                        "Error generating output for type '" + typeEntry.getKey() + "'", e);
                }
            }
        }
    }

    protected void renderTemplate(File outputDir, Template template, Map<String, Object> context)
        throws TemplateException, IOException, GenerationException {
        // Create a variable size output location where the template can generate it's content to
        ByteArrayOutputStream output = new ByteArrayOutputStream();

        // Have Freemarker generate the output
        template.process(context, new OutputStreamWriter(output));

        // Create the means to read in the generated output back in again
        try (BufferedReader input = new BufferedReader(new InputStreamReader(
            new ByteArrayInputStream(output.toByteArray())))) {

            // Extract the output path from the first line of the generated content
            String outputFileName = input.readLine();
            // If there is no outputFileName, this file should be skipped.
            if (outputFileName == null) {
                return;
            }
            File outputFile = new File(outputDir, outputFileName);

            // Create any missing directories
            if (!outputFile.getParentFile().exists() && !outputFile.getParentFile().mkdirs()) {
                throw new GenerationException("Unable to create output directory " + outputFile.getParent());
            }

            // Output the rest to that file
            try (BufferedWriter outputFileWriter = Files.newBufferedWriter(
                outputFile.toPath(), StandardCharsets.UTF_8)) {
                String line;
                while ((line = input.readLine()) != null) {
                    outputFileWriter.write(line);
                    outputFileWriter.newLine();
                }
                outputFileWriter.flush();
            }

            // Apply post-processing to the template
            postProcessTemplateOutput(outputFile);
        }
    }

    protected void postProcessTemplateOutput(File outputFile) {
        // NOOP
    }

    private Configuration getFreemarkerConfiguration() throws GenerationException {
        Configuration configuration = new Configuration(Configuration.VERSION_2_3_28);
        configuration.setDefaultEncoding("UTF-8");
        configuration.setTemplateExceptionHandler(TemplateExceptionHandler.RETHROW_HANDLER);
        configuration.setLogTemplateExceptions(false);
        configuration.setWrapUncheckedExceptions(true);
        try {
            configuration.setDirectoryForTemplateLoading(new File("/"));
        } catch (IOException e) {
            throw new GenerationException("Error setting directory for template loading", e);
        }
        return configuration;
    }

    protected abstract List<Template> getSpecTemplates(Configuration freemarkerConfiguration) throws IOException;

    protected abstract List<Template> getComplexTypeTemplates(Configuration freemarkerConfiguration) throws IOException;

    protected abstract List<Template> getEnumTypeTemplates(Configuration freemarkerConfiguration) throws IOException;

    protected abstract List<Template> getDataIoTemplates(Configuration freemarkerConfiguration) throws IOException;

    protected abstract FreemarkerLanguageTemplateHelper getHelper(TypeDefinition thisType, String protocolName, String flavorName, Map<String, TypeDefinition> types,
                                                                  Map<String, String> options);

}
