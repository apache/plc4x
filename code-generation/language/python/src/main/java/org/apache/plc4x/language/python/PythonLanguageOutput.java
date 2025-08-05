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
package org.apache.plc4x.language.python;

import freemarker.template.Configuration;
import freemarker.template.Template;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.FreemarkerLanguageOutput;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.FreemarkerLanguageTemplateHelper;
import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;

import java.io.IOException;
import java.util.*;

public class PythonLanguageOutput extends FreemarkerLanguageOutput {

    /**
     * Returns the name of the language, which is used for outputting the generated code.
     *
     * @return The name of the language.
     */
    @Override
    public String getName() {
        return "Python";
    }

    /**
     * Returns a set of options that are supported by this language output flavor.
     * Options are key-value pairs that can be passed to the code generator to customize the generated code.
     *
     * @return A set of supported options.
     */
    @Override
    public Set<String> supportedOptions() {
        return Collections.singleton("package");
    }


    /**
     * Returns a list of output flavors that are supported by this language output flavor.
     * An output flavor is a customization option that can be passed to the code generator to
     * produce code for a specific use case.
     *
     * @return A list of supported output flavors.
     */
    @Override
    public List<String> supportedOutputFlavors() {
        return Arrays.asList("read-write", "read-only", "passive");
    }

    /**
     * Returns a list of templates that are used to generate a specification of the generated code.
     * The specification is usually a markdown file that contains information about the structure of the
     * generated code, such as method signatures, comments, and so on.
     *
     * @return A list of templates that are used to generate the specification.
     */
    @Override
    protected List<Template> getSpecTemplates(Configuration freemarkerConfiguration) {
        /**
         * Python does not need a spec as it generates the code directly for the user.
         * Therefore, this method returns an empty list.
         */
        return Collections.emptyList();
    }

    /**
     * Returns a list of templates that are used to generate code for a complex type.
     * A complex type is a type that is composed of multiple fields.
     *
     * @param freemarkerConfiguration The configuration of the FreeMarker engine.
     * @return A list of templates that are used to generate code for a complex type.
     * @throws IOException Thrown if an error occurs while reading a template file.
     */
    @Override
    protected List<Template> getComplexTypeTemplates(Configuration freemarkerConfiguration) throws IOException {
        /**
         * The template used to generate code for a complex type is currently a single file.
         * The template is stored at 'templates/python/complex-type-template.python.ftlh'.
         *
         * The template is used to generate code for a complex type by passing the complex type
         * object as a variable named 'type' to the template.
         *
         * The template is also passed information about the fields of the complex type,
         * such as their names, types and so on.
         */
        return List.of(freemarkerConfiguration.getTemplate("templates/python/complex-type-template.python.ftlh"));
    }

    /**
     * Returns a list of templates that are used to generate code for an enumeration type.
     * An enumeration type is a type that is composed of a set of constants.
     *
     * @param freemarkerConfiguration The configuration of the FreeMarker engine.
     * @return A list of templates that are used to generate code for an enumeration type.
     * @throws IOException Thrown if an error occurs while reading a template file.
     */
    @Override
    protected List<Template> getEnumTypeTemplates(Configuration freemarkerConfiguration) throws IOException {
        /**
         * The template used to generate code for an enumeration type is currently a single file.
         * The template is stored at 'templates/python/enum-template.python.ftlh'.
         *
         * The template is used to generate code for an enumeration type by passing the enumeration type
         * object as a variable named 'type' to the template.
         *
         * The template is also passed information about the constants of the enumeration type,
         * such as their names and so on.
         */
        return Collections.singletonList(
            freemarkerConfiguration.getTemplate("templates/python/enum-template.python.ftlh"));
    }

    /**
     * Returns a list of templates that are used to generate code for a data IO class.
     *
     * A data IO class is a class that is responsible for reading and/or writing
     * a specific type of data from and/or to a PLC (Programmable Logic Controller).
     *
     * The templates are used to generate code for a data IO class by passing the
     * type of data that is being read and/or written as a variable named 'type'
     * to the templates.
     *
     * The templates are also passed information about the fields of the type
     * that is being read and/or written, such as their names, types and so on.
     *
     * @param freemarkerConfiguration The configuration of the FreeMarker engine.
     * @return A list of templates that are used to generate code for a data IO class.
     * @throws IOException Thrown if an error occurs while reading a template file.
     */
    @Override
    protected List<Template> getDataIoTemplates(Configuration freemarkerConfiguration) throws IOException {
        /**
         * The template used to generate code for a data IO class is currently a single file.
         * The template is stored at 'templates/python/data-io-template.python.ftlh'.
         *
         * The template is used to generate code for a data IO class by passing the
         * type of data that is being read and/or written as a variable named 'type' to the
         * template.
         *
         * The template is also passed information about the fields of the type
         * that is being read and/or written, such as their names, types and so on.
         */
        return Collections.singletonList(
            freemarkerConfiguration.getTemplate("templates/python/data-io-template.python.ftlh"));
    }

    /**
     * Returns an instance of a helper class that is used by the code generation
     * templates to generate code.
     *
     * The helper class is responsible for generating code that is dependent on the
     * language and the type of code being generated (e.g. data IO, enum, etc.).
     *
     * The helper class is also responsible for generating code that is dependent on
     * the PLC protocol that is being generated code for.
     *
     * @param thisType     The type that is being generated code for.
     * @param protocolName The name of the PLC protocol that the code is being generated
     *                     for.
     * @param flavorName  The name of the code generation flavor that the code is being
     *                     generated for.
     * @param types        A map of type definitions.
     * @param options      A map of options that are passed to the code generator.
     * @return An instance of a helper class.
     */
    @Override
    protected FreemarkerLanguageTemplateHelper getHelper(TypeDefinition thisType, String protocolName, String flavorName, Map<String, TypeDefinition> types,
                                                         Map<String, Object> options) {
        return new PythonLanguageTemplateHelper(thisType, protocolName, flavorName, types);
    }

}
