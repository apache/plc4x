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
package org.apache.plc4x.plugins.codegenerator.language.mspec.model.fields;

import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.fields.ArrayField;
import org.apache.plc4x.plugins.codegenerator.types.references.ArrayTypeReference;
import org.apache.plc4x.plugins.codegenerator.types.references.TypeReference;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;

import java.util.Map;
import java.util.Objects;

public class DefaultArrayField extends DefaultTypedNamedField implements ArrayField {

    private final LoopType loopType;
    private final Term loopExpression;

    public DefaultArrayField(Map<String, Term> attributes, String name, LoopType loopType, Term loopExpression) {
        super(attributes, name);
        this.loopType = Objects.requireNonNull(loopType);
        this.loopExpression = Objects.requireNonNull(loopExpression);
    }

    public LoopType getLoopType() {
        return loopType;
    }

    public Term getLoopExpression() {
        return loopExpression;
    }

    @Override
    public void setType(TypeReference typeReference) {
        if(!(typeReference instanceof ArrayTypeReference)) {
            throw new IllegalArgumentException("Array fields can only have ArrayTypeReferences");
        }
        super.setType(typeReference);
    }

    @Override
    public ArrayTypeReference getType() {
        return (ArrayTypeReference) super.getType();
    }

    @Override
    public String toString() {
        return "DefaultArrayField{" +
            "loopType=" + loopType +
            ", loopExpression=" + loopExpression +
            "} " + super.toString();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        if (!super.equals(o)) return false;
        DefaultArrayField that = (DefaultArrayField) o;
        return loopType == that.loopType && Objects.equals(loopExpression, that.loopExpression);
    }

    @Override
    public int hashCode() {
        return Objects.hash(super.hashCode(), loopType, loopExpression);
    }
}
