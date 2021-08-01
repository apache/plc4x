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
package org.apache.plc4x.language.c;

import freemarker.template.Configuration;
import freemarker.template.Template;
import java.util.Set;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.FreemarkerLanguageOutput;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.FreemarkerLanguageTemplateHelper;
import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;

import java.io.IOException;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.Map;

public class CLanguageOutput extends FreemarkerLanguageOutput {

    @Override
    public String getName() {
        return "C";
    }

    @Override
    public Set<String> supportedOptions() {
        return Collections.emptySet();
    }

    @Override
    public List<String> supportedOutputFlavors() {
        return Collections.singletonList("read-write");
    }

    @Override
    protected List<Template> getSpecTemplates(Configuration freemarkerConfiguration) throws IOException {
        return Collections.emptyList();
    }

    @Override
    protected List<Template> getComplexTypeTemplates(Configuration freemarkerConfiguration) throws IOException {
        return Arrays.asList(
            freemarkerConfiguration.getTemplate("templates/c/pojo-template.h.ftlh"),
            freemarkerConfiguration.getTemplate("templates/c/pojo-template.c.ftlh"));
    }

    @Override
    protected List<Template> getEnumTypeTemplates(Configuration freemarkerConfiguration) throws IOException {
        return Arrays.asList(
            freemarkerConfiguration.getTemplate("templates/c/enum-template.h.ftlh"),
            freemarkerConfiguration.getTemplate("templates/c/enum-template.c.ftlh"));
    }

    @Override
    protected List<Template> getDataIoTemplates(Configuration freemarkerConfiguration) throws IOException {
        return Arrays.asList(
            freemarkerConfiguration.getTemplate("templates/c/data-io-template.h.ftlh"),
            freemarkerConfiguration.getTemplate("templates/c/data-io-template.c.ftlh"));
    }

    @Override
    protected FreemarkerLanguageTemplateHelper getHelper(TypeDefinition thisType, String protocolName, String flavorName, Map<String, TypeDefinition> types,
        Map<String, String> options) {
        return new CLanguageTemplateHelper(thisType, protocolName, flavorName, types);
    }

}
