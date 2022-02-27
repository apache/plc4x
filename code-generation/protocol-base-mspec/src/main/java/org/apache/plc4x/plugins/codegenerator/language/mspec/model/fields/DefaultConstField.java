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
package org.apache.plc4x.plugins.codegenerator.language.mspec.model.fields;

import org.apache.plc4x.plugins.codegenerator.types.fields.ConstField;
import org.apache.plc4x.plugins.codegenerator.types.terms.Literal;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;

import java.util.*;

public class DefaultConstField extends DefaultTypedNamedField implements ConstField {

    private final Literal referenceValue;

    public DefaultConstField(Map<String, Term> attributes, String name, Literal referenceValue) {
        super(attributes, name);
        this.referenceValue = Objects.requireNonNull(referenceValue);
    }

    public Literal getReferenceValue() {
        return referenceValue;
    }

    @Override
    public String toString() {
        return "DefaultConstField{" +
            "referenceValue=" + referenceValue +
            "} " + super.toString();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        if (!super.equals(o)) return false;
        DefaultConstField that = (DefaultConstField) o;
        return Objects.equals(referenceValue, that.referenceValue);
    }

    @Override
    public int hashCode() {
        return Objects.hash(super.hashCode(), referenceValue);
    }
}
