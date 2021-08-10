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

import freemarker.template.*;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.FreemarkerLanguageOutput;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.FreemarkerLanguageTemplateHelper;
import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;

import java.io.*;
import java.util.*;
import org.apache.plc4x.plugins.codegenerator.types.exceptions.GenerationException;

public class JavaLanguageOutput extends FreemarkerLanguageOutput {

    @Override
    public String getName() {
        return "Java";
    }

    @Override
    public List<String> supportedOutputFlavors() {
        return Arrays.asList("read-write", "read-only", "passive");
    }

    @Override
    protected List<Template> getSpecTemplates(Configuration freemarkerConfiguration) throws IOException {
        return Collections.singletonList(
            freemarkerConfiguration.getTemplate("templates/java/enum-package-info-template.java.ftlh"));
    }

    @Override
    protected List<Template> getComplexTypeTemplates(Configuration freemarkerConfiguration) throws IOException {
        return Arrays.asList(
            freemarkerConfiguration.getTemplate("templates/java/pojo-template.java.ftlh"),
            freemarkerConfiguration.getTemplate("templates/java/io-template.java.ftlh"));
    }

    @Override
    protected List<Template> getEnumTypeTemplates(Configuration freemarkerConfiguration) throws IOException {
        return Collections.singletonList(
            freemarkerConfiguration.getTemplate("templates/java/enum-template.java.ftlh"));
    }

    @Override
    protected List<Template> getDataIoTemplates(Configuration freemarkerConfiguration) throws IOException {
        return Collections.singletonList(
            freemarkerConfiguration.getTemplate("templates/java/data-io-template.java.ftlh"));
    }

    @Override
    protected FreemarkerLanguageTemplateHelper getHelper(TypeDefinition thisType, String protocolName, String flavorName, Map<String, TypeDefinition> types) {
        return new JavaLanguageTemplateHelper(thisType, protocolName, flavorName, types);
    }


}
