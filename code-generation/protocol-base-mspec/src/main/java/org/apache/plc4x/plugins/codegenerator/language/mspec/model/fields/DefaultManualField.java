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

import org.apache.plc4x.plugins.codegenerator.types.fields.ManualField;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;

import java.util.Map;
import java.util.Objects;

public class DefaultManualField extends DefaultTypedNamedField implements ManualField {

    private final Term parseExpression;
    private final Term serializeExpression;
    private final Term lengthExpression;

    public DefaultManualField(Map<String, Term> attributes, String name, Term parseExpression, Term serializeExpression, Term lengthExpression) {
        super(attributes, name);
        this.parseExpression = Objects.requireNonNull(parseExpression);
        this.serializeExpression = Objects.requireNonNull(serializeExpression);
        this.lengthExpression = Objects.requireNonNull(lengthExpression);
    }

    public Term getParseExpression() {
        return parseExpression;
    }

    public Term getSerializeExpression() {
        return serializeExpression;
    }

    public Term getLengthExpression() {
        return lengthExpression;
    }

    @Override
    public String toString() {
        return "DefaultManualField{" +
            "parseExpression=" + parseExpression +
            ", serializeExpression=" + serializeExpression +
            ", lengthExpression=" + lengthExpression +
            "} " + super.toString();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        if (!super.equals(o)) return false;
        DefaultManualField that = (DefaultManualField) o;
        return Objects.equals(parseExpression, that.parseExpression) && Objects.equals(serializeExpression, that.serializeExpression) && Objects.equals(lengthExpression, that.lengthExpression);
    }

    @Override
    public int hashCode() {
        return Objects.hash(super.hashCode(), parseExpression, serializeExpression, lengthExpression);
    }
}
