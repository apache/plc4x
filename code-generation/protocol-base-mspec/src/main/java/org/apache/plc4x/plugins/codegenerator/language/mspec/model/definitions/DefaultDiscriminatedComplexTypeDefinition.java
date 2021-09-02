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
package org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions;

import org.apache.plc4x.plugins.codegenerator.types.definitions.Argument;
import org.apache.plc4x.plugins.codegenerator.types.definitions.DiscriminatedComplexTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.fields.DiscriminatorField;
import org.apache.plc4x.plugins.codegenerator.types.fields.Field;

import java.util.List;

public class DefaultDiscriminatedComplexTypeDefinition extends DefaultComplexTypeDefinition implements DiscriminatedComplexTypeDefinition {

    private final String[] discriminatorValues;

    public DefaultDiscriminatedComplexTypeDefinition(String name, Argument[] parserArguments, String[] tags, String[] discriminatorValues, List<Field> fields) {
        super(name, parserArguments, tags, false, fields);
        this.discriminatorValues = discriminatorValues;
    }

    public DiscriminatorField getDiscriminatorField() {
        // For a discriminated type, the discriminator is always defined in the parent type,
        // which is always a DefaultComplexTypeDefinition instance.
        return ((DefaultComplexTypeDefinition) getParentType()).getFields().stream().filter(
            field -> field instanceof DiscriminatorField).map(
            field -> (DiscriminatorField) field).findFirst().orElse(null);
    }

    public String[] getDiscriminatorValues() {
        return discriminatorValues;
    }

}
