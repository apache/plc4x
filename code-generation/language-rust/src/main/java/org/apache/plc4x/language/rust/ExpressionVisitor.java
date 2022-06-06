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

package org.apache.plc4x.language.rust;

import org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions.DefaultComplexTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.references.DefaultArrayTypeReference;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.terms.DefaultBinaryTerm;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.terms.DefaultNullLiteral;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.terms.DefaultNumericLiteral;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.terms.DefaultUnaryTerm;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.terms.DefaultVariableLiteral;
import org.apache.plc4x.plugins.codegenerator.types.fields.Field;
import org.apache.plc4x.plugins.codegenerator.types.fields.NamedField;
import org.apache.plc4x.plugins.codegenerator.types.fields.TypedField;
import org.apache.plc4x.plugins.codegenerator.types.references.ComplexTypeReference;
import org.apache.plc4x.plugins.codegenerator.types.references.SimpleTypeReference;
import org.apache.plc4x.plugins.codegenerator.types.references.TypeReference;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;

import java.util.Optional;

public class ExpressionVisitor {

    private final DefaultComplexTypeDefinition type;
    private final RustLanguageTemplateHelper helper;

    public ExpressionVisitor(DefaultComplexTypeDefinition type, RustLanguageTemplateHelper helper) {
        this.type = type;
        this.helper = helper;
    }

    public String[] visitLiteral(DefaultVariableLiteral literal) {
        switch (literal.getName()) {
            case "COUNT":
                return new String[]{this.visit(literal.getArgs().get().get(0))[0] + ".len()"};
            case "ARRAY_SIZE_IN_BYTES":
                String typeName = literal.getArgs().get().get(0).getDiscriminatorName();
                // Get the dimension of the given type
                Optional<Field> any = this.type.getAllFields().stream().filter(field -> field instanceof NamedField).filter(field -> typeName.equals(((NamedField) field).getName())).findAny();
                if (any.isEmpty()) {
                    throw new RuntimeException("Unable to get type of " + typeName);
                }
                TypeReference type = ((TypedField) any.get()).getType();

                if (!(type instanceof DefaultArrayTypeReference)) {
                    throw new RuntimeException("This should not happen!");
                }
                TypeReference innerType = ((DefaultArrayTypeReference) type).getElementTypeReference();
                if (innerType instanceof SimpleTypeReference) {
                    String numBits = Integer.toString(this.helper.getNumBits(((SimpleTypeReference) innerType)));
                    return new String[]{this.visit(literal.getArgs().get().get(0))[0] + ".len() * " + numBits};
                } else if (innerType instanceof ComplexTypeReference) {
                    // sum up all entries
                    return new String[]{"let mut s: u32 = 0;", "for x in &" + this.visit(literal.getArgs().get().get(0))[0] + " {", "\ts += x.get_length_in_bits();",  "};",  "s"};
                } else {
                    throw new RuntimeException("This should not happen!");
                }
            case "lengthInBytes":
                return new String[]{"get_length_in_bytes()"};
            case "lengthInBits":
                return new String[]{"get_length_in_bits()"};
            default:
                if (literal.getChild().isPresent()) {
                    return new String[]{"self." + literal.getName() + "." + this.visitLiteral(((DefaultVariableLiteral) literal.getChild().get()))[0]};
                } else {
                    return new String[]{"self." + literal.getName()};
                }
        }
    }

    public String[] visit(Term exp) {
        if (exp instanceof DefaultUnaryTerm) {
            return this.visit(((DefaultUnaryTerm) exp).getA());
        }
        if (exp instanceof DefaultVariableLiteral) {
            return this.visitLiteral(((DefaultVariableLiteral) exp));
        }
        if (exp instanceof DefaultNumericLiteral) {
            return new String[]{exp.stringRepresentation()};
        }
        if (exp instanceof DefaultBinaryTerm) {
            return new String[]{this.visit(((DefaultBinaryTerm) exp).getA())[0] + " " + ((DefaultBinaryTerm) exp).getOperation() + " " + this.visit(((DefaultBinaryTerm) exp).getB())[0]};
        }
        throw new RuntimeException("Not yet implemented in visit: " + exp);
    }
}
