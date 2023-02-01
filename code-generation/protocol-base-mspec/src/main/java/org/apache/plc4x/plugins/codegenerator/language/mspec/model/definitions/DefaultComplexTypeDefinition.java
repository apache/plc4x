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
package org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions;

import org.apache.plc4x.plugins.codegenerator.types.definitions.Argument;
import org.apache.plc4x.plugins.codegenerator.types.definitions.ComplexTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.fields.*;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;

import java.util.*;
import java.util.stream.Collectors;

public class DefaultComplexTypeDefinition extends DefaultTypeDefinition implements ComplexTypeDefinition {

    private final boolean isAbstract;
    private final List<Field> fields;
    protected ComplexTypeDefinition parentType;

    public DefaultComplexTypeDefinition(String name, Map<String, Term> attributes, List<Argument> parserArguments, boolean isAbstract, List<Field> fields) {
        super(name, attributes, parserArguments);
        this.isAbstract = isAbstract;
        this.fields = Objects.requireNonNull(fields);
    }

    public Optional<ComplexTypeDefinition> getParentType() {
        return Optional.ofNullable(parentType);
    }

    public void setParentType(ComplexTypeDefinition parentType) {
        this.parentType = parentType;
    }

    public Optional<List<Argument>> getAllParserArguments() {
        List<Argument> allArguments = new ArrayList<>();
        allArguments.addAll(getParserArguments().orElse(Collections.emptyList()));
        if(getParentType().isPresent()) {
            ComplexTypeDefinition parent = getParentType().get();
            allArguments.addAll(parent.getAllParserArguments().orElse(Collections.emptyList()));
        }
        return Optional.of(allArguments);
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
    public List<Field> getAllFields() {
        List<Field> fields = new LinkedList<>();
        getParentType()
            .map(ComplexTypeDefinition::getAllFields)
            .map(fields::addAll);
        fields.addAll(getFields());
        return fields;
    }

    @Override
    public List<PropertyField> getAllPropertyFields() {
        List<PropertyField> fields = new LinkedList<>();
        getParentType()
            .map(ComplexTypeDefinition::getAllPropertyFields)
            .map(fields::addAll);
        fields.addAll(getPropertyFields());
        return fields;
    }

    @Override
    public List<VirtualField> getAllVirtualFields() {
        List<VirtualField> fields = new LinkedList<>();
        getParentType()
            .map(ComplexTypeDefinition::getAllVirtualFields)
            .map(fields::addAll);
        fields.addAll(getVirtualFields());
        return fields;
    }

    @Override
    public List<PropertyField> getParentPropertyFields() {
        return getParentType().map(ComplexTypeDefinition::getAllPropertyFields).orElse(Collections.emptyList());
    }

    @Override
    public String toString() {
        return "DefaultComplexTypeDefinition{" +
            "isAbstract=" + isAbstract +
            ", fields=" + fields +
            ", parentType=" + (parentType != null ? parentType.getName() : null) +
            "} " + super.toString();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        if (!super.equals(o)) return false;
        DefaultComplexTypeDefinition that = (DefaultComplexTypeDefinition) o;
        return isAbstract == that.isAbstract && Objects.equals(fields, that.fields);
    }

    @Override
    public int hashCode() {
        return Objects.hash(super.hashCode(), isAbstract, fields);
    }
}
