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
package org.apache.plc4x.plugins.codegenerator.language.mspec.model.terms;

import org.apache.plc4x.plugins.codegenerator.types.terms.BinaryTerm;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;

import java.util.Objects;

public class DefaultBinaryTerm implements BinaryTerm {

    private final Term a;
    private final Term b;
    private final String operation;

    public DefaultBinaryTerm(Term a, Term b, String operation) {
        this.a = Objects.requireNonNull(a);
        this.b = Objects.requireNonNull(b);
        this.operation = Objects.requireNonNull(operation);
    }

    public Term getA() {
        return a;
    }

    public Term getB() {
        return b;
    }

    public String getOperation() {
        return operation;
    }

    @Override
    public String stringRepresentation() {
        return a + operation + b;
    }

    @Override
    public String toString() {
        return "DefaultBinaryTerm{" +
                "a=" + a +
                ", b=" + b +
                ", operation='" + operation + '\'' +
                '}';
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        DefaultBinaryTerm that = (DefaultBinaryTerm) o;
        return a.equals(that.a) && b.equals(that.b) && operation.equals(that.operation);
    }

    @Override
    public int hashCode() {
        return Objects.hash(a, b, operation);
    }
}
