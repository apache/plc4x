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

package org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions;

import org.apache.plc4x.plugins.codegenerator.types.definitions.Argument;
import org.apache.plc4x.plugins.codegenerator.types.definitions.ComplexTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.fields.*;
import org.apache.plc4x.plugins.codegenerator.types.references.IntegerTypeReference;

import java.util.ArrayList;
import java.util.Collections;
import java.util.LinkedList;
import java.util.List;
import java.util.Optional;
import java.util.stream.Collectors;

public class DefaultComplexTypeDefinition extends DefaultTypeDefinition implements ComplexTypeDefinition {

    private final boolean isAbstract;
    private final List<Field> fields;

    public DefaultComplexTypeDefinition(String name, Argument[] parserArguments, String[] tags, boolean littleEndian, boolean isAbstract, List<Field> fields) {
        super(name, parserArguments, tags, littleEndian);
        this.isAbstract = isAbstract;
        this.fields = fields;
    }

    public boolean isAbstract() {
        return isAbstract;
    }

    public List<Field> getFields() {
        return fields;
    }

    @Override
    public List<SimpleField> getSimpleFields() {
        return fields.stream().filter(field -> field instanceof SimpleField).map(
            field -> (SimpleField) field).collect(Collectors.toList());
    }

    @Override
    public List<ConstField> getConstFields() {
        return fields.stream().filter(field -> field instanceof ConstField).map(
            field -> (ConstField) field).collect(Collectors.toList());
    }

    @Override
    public List<PropertyField> getPropertyFields() {
        return fields.stream().filter(field -> ((field instanceof PropertyField) && !(field instanceof ConstField) && !(field instanceof VirtualField))).map(field -> (PropertyField) field)
            .collect(Collectors.toList());
    }

    @Override
    public List<VirtualField> getVirtualFields() {
        return fields.stream().filter(field -> (field instanceof VirtualField)).map(field -> (VirtualField) field)
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
        if (getParentType() != null) {
            return ((ComplexTypeDefinition) getParentType()).getAllPropertyFields();
        }
        return Collections.emptyList();
    }

    public List<Field> getFieldsInOrder() {
        if (!isLittleEndian()) {
            return getFields();
        }

        List<Field> ordered = new ArrayList<>();
        List<Field> group = new ArrayList<>();
        int groupSize = 0;
        int length = this.fields.size();
        for (int index = length; index > 0; index--) {
            Field field = fields.get(index - 1);

            if (field instanceof TypedField) {
                TypedField typed = (TypedField) field;
                if (typed.getType() instanceof IntegerTypeReference) {
                    IntegerTypeReference size = (IntegerTypeReference) typed.getType();
                    group.add(field);
                    groupSize += size.getSizeInBits();

                    // we can have odd sizes then its best to
                    if ((groupSize % 8) == 0) {
                        Collections.reverse(group);
                        ordered.addAll(group);
                        group.clear();
                        groupSize = 0;
                    }
                } else {
                    ordered.add(field);
                }
            } else {
                ordered.add(field);
            }
        }

        return ordered;
    }
}
