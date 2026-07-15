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

package org.apache.plc4x.codegeneration.language.java;

import org.apache.plc4x.codegeneration.language.java.generators.ComplexTypeGenerator;
import org.apache.plc4x.codegeneration.language.java.generators.ConstantsTypeGenerator;
import org.apache.plc4x.codegeneration.language.java.generators.DataIoGenerator;
import org.apache.plc4x.codegeneration.language.java.generators.EnumGenerator;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.plugins.codegenerator.language.LanguageOutput;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions.DefaultComplexTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions.DefaultDataIoTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions.DefaultDiscriminatedComplexTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions.DefaultEnumTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.terms.DefaultBooleanLiteral;
import org.apache.plc4x.plugins.codegenerator.types.definitions.ConstantsTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.exceptions.GenerationException;
import org.apache.plc4x.plugins.codegenerator.types.terms.BooleanLiteral;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.util.*;
import java.util.stream.Collectors;

public class JavaLanguageOutput implements LanguageOutput {

    private static final Logger LOGGER = LoggerFactory.getLogger(JavaLanguageOutput.class);

    private final Formatter formatter = new Formatter();

    @Override
    public String getName() {
        return "Java";
    }

    @Override
    public Set<String> supportedOptions() {
        return new HashSet<>(Arrays.asList(
            // Overrides the package name generated.
            "package",
            // Generates additional properties to save values of reserved fields for the case that the value differs from the expected value.
            "generate-properties-for-reserved-fields",
            // Map containing the type-names for external types.
            "externalTypes"));
    }

    @Override
    @SuppressWarnings("unchecked")
    public void generate(File outputRootDir, String version, String languageName, String protocolName, String outputFlavor, Map<String, TypeDefinition> types, Map<String, Object> options) throws GenerationException {
        Map<String, String> externalTypes = (Map<String, String>) options.get("externalTypes");

        String targetPackage = (String) options.getOrDefault("package", "org.apache.plc4x." + languageName + "." + protocolName.replace("-", "") + "." + outputFlavor.replace("-", ""));
        File outputDirectory = new File(outputRootDir, targetPackage.replace('.', '/'));
        if (!outputDirectory.exists() && !outputDirectory.mkdirs()) {
            throw new GenerationException("Unable to create output directory " + outputDirectory);
        }

        // Output a general plan of what we are planning on generating.
        outputGenerationPlan(types);

        try {
            ArrayList<ConstantsTypeDefinition> constantTypes = types.values().stream().filter(typeDefinition -> typeDefinition instanceof ConstantsTypeDefinition).map(typeDefinition -> (ConstantsTypeDefinition) typeDefinition).collect(Collectors.toCollection(ArrayList::new));
            ConstantsTypeGenerator constantsTypeGenerator = new ConstantsTypeGenerator(targetPackage, outputDirectory, outputRootDir, types, externalTypes, options);
            for (ConstantsTypeDefinition constantType : constantTypes) {
                Term external = constantType.getAttribute("external").orElse(new DefaultBooleanLiteral(false));
                // Only output a file if it's not an "external" type.
                if (external.isLiteral() &&
                    ((external.asLiteral().orElseThrow() instanceof BooleanLiteral) &&
                        !((BooleanLiteral) external.asLiteral().orElseThrow()).getValue())) {
                    constantsTypeGenerator.generate(constantType);
                }
            }

            ArrayList<DefaultComplexTypeDefinition> rootComplexTypes = types.values().stream().filter(typeDefinition -> typeDefinition instanceof DefaultComplexTypeDefinition).map(typeDefinition -> (DefaultComplexTypeDefinition) typeDefinition).filter(defaultComplexTypeDefinition -> !((defaultComplexTypeDefinition instanceof DefaultDiscriminatedComplexTypeDefinition) || (defaultComplexTypeDefinition instanceof DefaultDataIoTypeDefinition))).collect(Collectors.toCollection(ArrayList::new));
            ComplexTypeGenerator complexTypeGenerator = new ComplexTypeGenerator(targetPackage, outputDirectory, outputRootDir, types, externalTypes, options);
            for (DefaultComplexTypeDefinition rootComplexType : rootComplexTypes) {
                Term external = rootComplexType.getAttribute("external").orElse(new DefaultBooleanLiteral(false));
                // Only output a file if it's not an "external" type.
                if (external.isLiteral() &&
                    ((external.asLiteral().orElseThrow() instanceof BooleanLiteral) &&
                        !((BooleanLiteral) external.asLiteral().orElseThrow()).getValue())) {
                    complexTypeGenerator.generate(rootComplexType);

                    // Generate the children of this type (if he has any):
                    if (rootComplexType.isAbstract()) {
                        generateChildren(rootComplexType, types, complexTypeGenerator);
                    }
                }
            }

            ArrayList<DefaultEnumTypeDefinition> enumTypes = types.values().stream().filter(typeDefinition -> typeDefinition instanceof DefaultEnumTypeDefinition).map(typeDefinition -> (DefaultEnumTypeDefinition) typeDefinition).collect(Collectors.toCollection(ArrayList::new));
            EnumGenerator enumGenerator = new EnumGenerator(targetPackage, outputDirectory, outputRootDir, types, externalTypes, options);
            for (DefaultEnumTypeDefinition enumType : enumTypes) {
                Term external = enumType.getAttribute("external").orElse(new DefaultBooleanLiteral(false));
                // Only output a file if it's not an "external" type.
                if (external.isLiteral() &&
                    ((external.asLiteral().orElseThrow() instanceof BooleanLiteral) &&
                        !((BooleanLiteral) external.asLiteral().orElseThrow()).getValue())) {
                    enumGenerator.generate(enumType);
                }
            }

            ArrayList<DefaultDataIoTypeDefinition> dataIoTypes = types.values().stream().filter(typeDefinition -> typeDefinition instanceof DefaultDataIoTypeDefinition).map(typeDefinition -> (DefaultDataIoTypeDefinition) typeDefinition).collect(Collectors.toCollection(ArrayList::new));
            DataIoGenerator dataIoGenerator = new DataIoGenerator(targetPackage, outputDirectory, outputRootDir, types, externalTypes, options);
            for (DefaultDataIoTypeDefinition ioType : dataIoTypes) {
                Term external = ioType.getAttribute("external").orElse(new DefaultBooleanLiteral(false));
                // Only output a file if it's not an "external" type.
                if (external.isLiteral() &&
                    ((external.asLiteral().orElseThrow() instanceof BooleanLiteral) &&
                        !((BooleanLiteral) external.asLiteral().orElseThrow()).getValue())) {
                    dataIoGenerator.generate(ioType);
                }
            }
        } catch (BufferException e) {
            throw new GenerationException("Error while generating code", e);
        }
    }

    @Override
    public List<String> supportedOutputFlavors() {
        return Arrays.asList("read-write", "read-only", "passive");
    }

    protected void outputGenerationPlan(Map<String, TypeDefinition> types) {
        LOGGER.info("");
        LOGGER.info("Complex Types:");
        ArrayList<DefaultComplexTypeDefinition> rootComplexTypes = types.values().stream().filter(typeDefinition -> typeDefinition instanceof DefaultComplexTypeDefinition).map(typeDefinition -> (DefaultComplexTypeDefinition) typeDefinition).filter(defaultComplexTypeDefinition -> !((defaultComplexTypeDefinition instanceof DefaultDiscriminatedComplexTypeDefinition) || (defaultComplexTypeDefinition instanceof DefaultDataIoTypeDefinition))).collect(Collectors.toCollection(ArrayList::new));
        rootComplexTypes.forEach(defaultComplexTypeDefinition -> {
            LOGGER.info(" - {}", defaultComplexTypeDefinition.getName());
            // List the children of this type:
            if (defaultComplexTypeDefinition.isAbstract()) {
                listChildren(defaultComplexTypeDefinition, types, "   - ");
            }
        });
        LOGGER.info("");

        LOGGER.info("Enum Types:");
        ArrayList<DefaultEnumTypeDefinition> enumTypes = types.values().stream().filter(typeDefinition -> typeDefinition instanceof DefaultEnumTypeDefinition).map(typeDefinition -> (DefaultEnumTypeDefinition) typeDefinition).collect(Collectors.toCollection(ArrayList::new));
        enumTypes.forEach(defaultEnumTypeDefinition -> {
            LOGGER.info(" - {}", defaultEnumTypeDefinition.getName());
        });
        LOGGER.info("");

        LOGGER.info("DataIo Types:");
        ArrayList<DefaultDataIoTypeDefinition> dataIoTypes = types.values().stream().filter(typeDefinition -> typeDefinition instanceof DefaultDataIoTypeDefinition).map(typeDefinition -> (DefaultDataIoTypeDefinition) typeDefinition).collect(Collectors.toCollection(ArrayList::new));
        dataIoTypes.forEach(defaultEnumTypeDefinition -> {
            LOGGER.info(" - {}", defaultEnumTypeDefinition.getName());
        });
    }

    protected void listChildren(DefaultComplexTypeDefinition parent, Map<String, TypeDefinition> types, String indent) {
        ArrayList<DefaultComplexTypeDefinition> children = types.values().stream().filter(typeDefinition -> typeDefinition instanceof DefaultComplexTypeDefinition).map(typeDefinition -> (DefaultComplexTypeDefinition) typeDefinition).filter(defaultComplexTypeDefinition -> defaultComplexTypeDefinition.getParentType().isPresent() && defaultComplexTypeDefinition.getParentType().get() == parent).collect(Collectors.toCollection(ArrayList::new));
        children.forEach(child -> {
            LOGGER.info("{}{}", indent, child.getName());
            if (child.isAbstract()) {
                listChildren(child, types, "  " + indent);
            }
        });
    }

    protected void generateChildren(DefaultComplexTypeDefinition parent, Map<String, TypeDefinition> types, ComplexTypeGenerator complexTypeGenerator) throws BufferException {
        ArrayList<DefaultComplexTypeDefinition> children = types.values().stream().filter(typeDefinition -> typeDefinition instanceof DefaultComplexTypeDefinition).map(typeDefinition -> (DefaultComplexTypeDefinition) typeDefinition).filter(defaultComplexTypeDefinition -> defaultComplexTypeDefinition.getParentType().isPresent() && defaultComplexTypeDefinition.getParentType().get() == parent).collect(Collectors.toCollection(ArrayList::new));
        for (DefaultComplexTypeDefinition child : children) {
            complexTypeGenerator.generate(child);
            if (child.isAbstract()) {
                generateChildren(child, types, complexTypeGenerator);
            }
        }
    }

}
