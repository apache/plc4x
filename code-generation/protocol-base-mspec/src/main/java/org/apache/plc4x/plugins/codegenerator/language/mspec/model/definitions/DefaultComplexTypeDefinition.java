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
import org.apache.plc4x.plugins.codegenerator.types.definitions.ComplexTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.fields.*;

import java.util.Collections;
import java.util.LinkedList;
import java.util.List;
import java.util.Objects;
import java.util.stream.Collectors;

public class DefaultComplexTypeDefinition extends DefaultTypeDefinition implements ComplexTypeDefinition {

    private final boolean isAbstract;
    private final List<Field> fields;

    public DefaultComplexTypeDefinition(String name, List<Argument> parserArguments, List<String> tags, boolean isAbstract, List<Field> fields) {
        super(name, parserArguments, tags);
        this.isAbstract = isAbstract;
        this.fields = Objects.requireNonNull(fields);
    }

    public boolean isAbstract() {
        return isAbstract;
    }

    public List<Field> getFields() {
        return fields;
    }

    @Override
    public List<SimpleField> getSimpleFields() {
        return fields.stream()
            .filter(SimpleField.class::isInstance)
            .map(SimpleField.class::cast)
            .collect(Collectors.toList());
    }

    @Override
    public List<ConstField> getConstFields() {
        return fields.stream()
            .filter(ConstField.class::isInstance)
            .map(ConstField.class::cast)
            .collect(Collectors.toList());
    }

    @Override
    public List<AssertField> getAssertFields() {
        return fields.stream()
            .filter(AssertField.class::isInstance)
            .map(AssertField.class::cast)
            .collect(Collectors.toList());
    }

    @Override
    public List<PropertyField> getPropertyFields() {
        return fields.stream()
            .filter(PropertyField.class::isInstance)
            .filter(field -> !(field instanceof ConstField) && !(field instanceof VirtualField))
            .map(PropertyField.class::cast)
            .collect(Collectors.toList());
    }

    @Override
    public List<AbstractField> getAbstractFields() {
        return fields.stream()
            .filter(AbstractField.class::isInstance)
            .map(AbstractField.class::cast)
            .collect(Collectors.toList());
    }

    public List<ImplicitField> getImplicitFields() {
        return fields.stream()
            .filter(ImplicitField.class::isInstance)
            .map(ImplicitField.class::cast)
            .collect(Collectors.toList());
    }

    @Override
    public List<VirtualField> getVirtualFields() {
        return fields.stream()
            .filter(VirtualField.class::isInstance)
            .map(VirtualField.class::cast)
            .collect(Collectors.toList());
    }

    @Override
    public List<PropertyField> getAllPropertyFields() {
        List<PropertyField> fields = new LinkedList<>();
        if (getParentType() != null) {
            fields.addAll(((ComplexTypeDefinition) getParentType()).getAllPropertyFields());
        }
        fields.addAll(getPropertyFields());
        return fields;
    }

    @Override
    public List<PropertyField> getParentPropertyFields() {
        if (getParentType() == null) {
            return Collections.emptyList();
        }
        return ((ComplexTypeDefinition) getParentType()).getAllPropertyFields();
    }

}
