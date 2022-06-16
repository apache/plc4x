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
import org.apache.plc4x.plugins.codegenerator.types.fields.Field;
import org.apache.plc4x.plugins.codegenerator.types.fields.ValidationField;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;

import java.util.Map;
import java.util.Objects;
import java.util.Optional;

public class DefaultValidationField extends DefaultField implements ValidationField, Field {

    private final Term validationExpression;
    private final String description;

    private final boolean shouldFail;

    public DefaultValidationField(Map<String, Term> attributes, Term validationExpression, String description, boolean shouldFail) {
        super(attributes);
        this.validationExpression = Objects.requireNonNull(validationExpression);
        this.description = description;
        this.shouldFail = shouldFail;
    }

    @Override
    public Term getValidationExpression() {
        return validationExpression;
    }

    @Override
    public Optional<String> getDescription() {
        return Optional.ofNullable(description);
    }

    @Override
    public boolean shouldFail() {
        return shouldFail;
    }

    // TODO: dummy implementation as this is not really a field
    @Override
    public String getTypeName() {
        return "validation";
    }

    // TODO: dummy implementation as this is not really a field
    @Override
    public Optional<Term> getAttribute(String s) {
        return Optional.empty();
    }

    @Override
    public String toString() {
        return "DefaultValidationField{" +
            "validationExpression=" + validationExpression +
            ", description='" + description + '\'' +
            '}';
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        DefaultValidationField that = (DefaultValidationField) o;
        return Objects.equals(validationExpression, that.validationExpression) && Objects.equals(description, that.description);
    }

    @Override
    public int hashCode() {
        return Objects.hash(validationExpression, description);
    }
}
