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
package org.apache.plc4x.plugins.codegenerator.language.mspec.model.terms;

import org.apache.plc4x.plugins.codegenerator.types.terms.Term;
import org.apache.plc4x.plugins.codegenerator.types.terms.UnaryTerm;

import java.util.Objects;

public class DefaultUnaryTerm implements UnaryTerm {

    private final Term a;
    private final String operation;

    public DefaultUnaryTerm(Term a, String operation) {
        this.a = Objects.requireNonNull(a);
        this.operation = Objects.requireNonNull(operation);
    }

    public Term getA() {
        return a;
    }

    public String getOperation() {
        return operation;
    }

    @Override
    public String stringRepresentation() {
        return operation + getA().stringRepresentation();
    }

    @Override
    public String toString() {
        return "DefaultUnaryTerm{" +
                "a=" + a +
                ", operation='" + operation + '\'' +
                '}';
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        DefaultUnaryTerm that = (DefaultUnaryTerm) o;
        return a.equals(that.a) && operation.equals(that.operation);
    }

    @Override
    public int hashCode() {
        return Objects.hash(a, operation);
    }
}
