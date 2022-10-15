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
import org.apache.plc4x.plugins.codegenerator.types.fields.ChecksumField;
import org.apache.plc4x.plugins.codegenerator.types.references.SimpleTypeReference;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;

import java.util.*;

public class DefaultChecksumField extends DefaultTypedNamedField implements ChecksumField {

    private final Term checksumExpression;

    public DefaultChecksumField(Map<String, Term> attributes, SimpleTypeReference type, String name, Term checksumExpression) {
        super(attributes, name);
        this.checksumExpression = Objects.requireNonNull(checksumExpression);
        this.type = type;
    }

    public Term getChecksumExpression() {
        return checksumExpression;
    }

    @Override
    public SimpleTypeReference getType() {
        return (SimpleTypeReference) super.getType();
    }

    @Override
    public String toString() {
        return "DefaultChecksumField{" +
            "checksumExpression=" + checksumExpression +
            "} " + super.toString();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        if (!super.equals(o)) return false;
        DefaultChecksumField that = (DefaultChecksumField) o;
        return Objects.equals(checksumExpression, that.checksumExpression);
    }

    @Override
    public int hashCode() {
        return Objects.hash(super.hashCode(), checksumExpression);
    }
}
