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
import org.apache.plc4x.plugins.codegenerator.types.fields.ReservedField;
import org.apache.plc4x.plugins.codegenerator.types.references.SimpleTypeReference;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;

import java.util.*;

public class DefaultReservedField extends DefaultTypedField implements ReservedField {

    private final Object referenceValue;

    public DefaultReservedField(Map<String, Term> attributes, SimpleTypeReference type, Object referenceValue) {
        super(attributes);
        this.referenceValue = Objects.requireNonNull(referenceValue);
        this.type = type;
    }

    public Object getReferenceValue() {
        return referenceValue;
    }

    @Override
    public String toString() {
        return "DefaultReservedField{" +
            "referenceValue=" + referenceValue +
            "} " + super.toString();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        if (!super.equals(o)) return false;
        DefaultReservedField that = (DefaultReservedField) o;
        return Objects.equals(referenceValue, that.referenceValue);
    }

    @Override
    public int hashCode() {
        return Objects.hash(super.hashCode(), referenceValue);
    }
}
