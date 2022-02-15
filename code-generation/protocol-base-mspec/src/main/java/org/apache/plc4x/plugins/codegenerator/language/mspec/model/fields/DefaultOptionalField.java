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

import org.apache.plc4x.plugins.codegenerator.types.fields.OptionalField;
import org.apache.plc4x.plugins.codegenerator.types.references.TypeReference;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;

import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.Optional;

public class DefaultOptionalField extends DefaultTypedNamedField implements OptionalField {

    private final Term conditionExpression;

    public DefaultOptionalField(Map<String, Term> attributes, TypeReference type, String name, Term conditionExpression) {
        super(attributes, type, name);
        this.conditionExpression = conditionExpression;
    }

    public Optional<Term> getConditionExpression() {
        return Optional.ofNullable(conditionExpression);
    }

    @Override
    public String toString() {
        return "DefaultOptionalField{" +
            "conditionExpression=" + conditionExpression +
            "} " + super.toString();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        if (!super.equals(o)) return false;
        DefaultOptionalField that = (DefaultOptionalField) o;
        return Objects.equals(conditionExpression, that.conditionExpression);
    }

    @Override
    public int hashCode() {
        return Objects.hash(super.hashCode(), conditionExpression);
    }
}
