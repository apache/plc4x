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


import org.apache.plc4x.plugins.codegenerator.types.definitions.DiscriminatedComplexTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.fields.SwitchField;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;
import org.apache.plc4x.plugins.codegenerator.types.terms.VariableLiteral;

import java.util.*;

public class DefaultSwitchField implements SwitchField {

    private final List<VariableLiteral> variableLiterals;
    private final List<DiscriminatedComplexTypeDefinition> cases;

    public DefaultSwitchField(List<VariableLiteral> variableLiterals) {
        this.variableLiterals = Objects.requireNonNull(variableLiterals);
        this.cases = new LinkedList<>();
    }

    public List<VariableLiteral> getDiscriminatorExpressions() {
        return variableLiterals;
    }

    // TODO: replace with immutable
    public void addCase(DiscriminatedComplexTypeDefinition caseType) {
        cases.add(caseType);
    }

    public List<DiscriminatedComplexTypeDefinition> getCases() {
        return cases;
    }

    @Override
    public Optional<Term> getAttribute(String attributeName) {
        return Optional.empty();
    }

    @Override
    public String toString() {
        return "DefaultSwitchField{" +
            "variableLiterals=" + variableLiterals +
            /*", cases=" + cases +*/
            '}';
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        DefaultSwitchField that = (DefaultSwitchField) o;
        return Objects.equals(variableLiterals, that.variableLiterals) /*&& Objects.equals(cases, that.cases)*/;
    }

    @Override
    public int hashCode() {
        return Objects.hash(variableLiterals/*, cases*/);
    }
}
