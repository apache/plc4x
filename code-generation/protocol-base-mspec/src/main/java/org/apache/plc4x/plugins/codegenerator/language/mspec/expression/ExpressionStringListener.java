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
package org.apache.plc4x.plugins.codegenerator.language.mspec.expression;

import org.apache.plc4x.plugins.codegenerator.language.mspec.LazyTypeDefinitionConsumer;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions.DefaultArgument;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.fields.DefaultTypedField;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.terms.*;
import org.apache.plc4x.plugins.codegenerator.types.definitions.ComplexTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.references.TypeReference;
import org.apache.plc4x.plugins.codegenerator.types.terms.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.stream.Collectors;

public class ExpressionStringListener extends ExpressionBaseListener {

    private static final Logger LOGGER = LoggerFactory.getLogger(ExpressionStringListener.class);

    private final LazyTypeDefinitionConsumer lazyTypeDefinitionConsumer;

    private final String rootTypeName;

    private Stack<List<Term>> parserContexts;

    private Stack<CompletableFuture<TypeReference>> futureStack;

    private Term root;

    public ExpressionStringListener(LazyTypeDefinitionConsumer lazyTypeDefinitionConsumer, String rootTypeName) {
        this.lazyTypeDefinitionConsumer = lazyTypeDefinitionConsumer;
        this.rootTypeName = rootTypeName;
    }

    public Term getRoot() {
        return root;
    }

    @Override
    public void enterExpressionString(ExpressionParser.ExpressionStringContext ctx) {
        parserContexts = new Stack<>();
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitExpressionString(ExpressionParser.ExpressionStringContext ctx) {
        List<Term> roots = parserContexts.pop();
        if (roots.isEmpty()) {
            throw new RuntimeException("Empty Expression not supported.");
        }
        if (roots.size() != 1) {
            throw new RuntimeException("Expression can only contain one root term.");
        }
        root = roots.get(0);
    }

    /////////////////////////////////////////////////////////////////////////////////////////
    // Literals
    /////////////////////////////////////////////////////////////////////////////////////////

    @Override
    public void exitNullExpression(ExpressionParser.NullExpressionContext ctx) {
        parserContexts.peek().add(new DefaultNullLiteral());
    }

    @Override
    public void exitBoolExpression(ExpressionParser.BoolExpressionContext ctx) {
        parserContexts.peek().add(new DefaultBooleanLiteral(Boolean.parseBoolean(ctx.getText())));
    }

    @Override
    public void exitNumberExpression(ExpressionParser.NumberExpressionContext ctx) {
        String strValue = ctx.Number().getText();
        if (strValue.contains(".")) {
            parserContexts.peek().add(new DefaultNumericLiteral(Double.valueOf(strValue)));
        } else {
            parserContexts.peek().add(new DefaultNumericLiteral(Long.valueOf(strValue)));
        }
    }

    @Override
    public void exitHexExpression(ExpressionParser.HexExpressionContext ctx) {
        String hexValue = ctx.HexExpression().getText();
        parserContexts.peek().add(new DefaultHexadecimalLiteral(hexValue));
    }

    @Override
    public void exitStringExpression(ExpressionParser.StringExpressionContext ctx) {
        parserContexts.peek().add(new DefaultStringLiteral(ctx.getText().substring(1, ctx.getText().length() - 1)));
    }

    @Override
    public void enterIdentifierSegment(ExpressionParser.IdentifierSegmentContext ctx) {
        String propertyName = ctx.name.getText();

        CompletableFuture<TypeReference> typeReferenceFuture = new CompletableFuture<>();
        if (futureStack == null) {
            schedulePropertyResolution(propertyName, typeReferenceFuture, rootTypeName);
            futureStack = new Stack<>();
        } else {
            futureStack.peek().whenComplete((typeReference, throwable) -> {
                if (throwable != null) {
                    LOGGER.debug("Error processing variables", throwable);
                    return;
                }
                String typeName = typeReference.asNonSimpleTypeReference().orElseThrow().getName();
                schedulePropertyResolution(propertyName, typeReferenceFuture, typeName);
            });
        }
        futureStack.push(typeReferenceFuture);
        parserContexts.push(new LinkedList<>());
    }

    private void schedulePropertyResolution(String propertyName, CompletableFuture<TypeReference> typeReferenceFuture, String typeName) {
        lazyTypeDefinitionConsumer.setOrScheduleTypeDefinitionConsumer(typeName, (TypeDefinition typeDefinition) -> {
            if (!typeDefinition.isComplexTypeDefinition()) {
                typeReferenceFuture.completeExceptionally(new RuntimeException("is not a complex type"));
                return;
            }
            final ComplexTypeDefinition complexTypeDefinition = typeDefinition
                .asComplexTypeDefinition()
                .orElseThrow();
            // Check for property fields context
            Optional<DefaultTypedField> propertyFieldByName = complexTypeDefinition
                .getPropertyFieldByName(propertyName)
                .map(DefaultTypedField.class::cast);
            // Check for virtual fields context
            if (propertyFieldByName.isEmpty()) {
                // TODO: do we need all virtual fields (from parent too)
                propertyFieldByName = complexTypeDefinition.getVirtualFields().stream()
                    .filter(virtualField -> propertyName.equals(virtualField.getName()))
                    .map(DefaultTypedField.class::cast)
                    .findAny();
            }
            // Check for arguments context
            if (propertyFieldByName.isEmpty() && complexTypeDefinition.getAllParserArguments().isPresent()) {
                Optional<DefaultArgument> defaultArgument = complexTypeDefinition.getAllParserArguments().orElseThrow().stream()
                    .filter(argument -> propertyName.equals(argument.getName()))
                    .map(DefaultArgument.class::cast)
                    .findAny();
                if (defaultArgument.isPresent()) {
                    defaultArgument.get().getTypeReferenceCompletionStage().whenComplete((typeReference, throwable) -> {
                        if (throwable != null) {
                            typeReferenceFuture.completeExceptionally(throwable);
                        } else {
                            typeReferenceFuture.complete(typeReference);
                        }
                    });
                    return;
                }
            }
            if (propertyFieldByName.isEmpty()) {
                typeReferenceFuture.completeExceptionally(new RuntimeException("Field with name " + propertyName + " not found on " + typeName));
                return;
            }
            DefaultTypedField propertyField = propertyFieldByName.orElseThrow();
            propertyField.getTypeReferenceCompletionStage().whenComplete((propertyTypeReference, throwable) -> {
                if (throwable != null) {
                    typeReferenceFuture.completeExceptionally(throwable);
                } else {
                    typeReferenceFuture.complete(propertyTypeReference);
                }
            });
        });
    }

    @Override
    public void exitIdentifierSegment(ExpressionParser.IdentifierSegmentContext ctx) {
        List<Term> args = parserContexts.pop();
        ArgsContext argsContext = null;
        IndexContext indexContext = null;
        RestContext restContext = null;
        for (Term arg : args) {
            if (arg instanceof ArgsContext) {
                argsContext = (ArgsContext) arg;
            } else if (arg instanceof IndexContext) {
                indexContext = (IndexContext) arg;
            } else if (arg instanceof RestContext) {
                restContext = (RestContext) arg;
            }
        }

        String name = ctx.name.getText();
        // TODO: Based on the current context type-definition, get the type of the property with name ctx.name.getText()

        int index = VariableLiteral.NO_INDEX;
        if (indexContext != null) {
            // TODO: Add a check, that the field providing the property is an "array" or "manualArray" field.
            index = indexContext.getFirst().getNumber().intValue();
        }
        VariableLiteral rest = null;
        if (restContext != null) {
            // TODO: Add a check, that the field providing the property references a complex type (or uses one of the built-ins)
            rest = restContext.getFirst();
        }

        final DefaultVariableLiteral variableLiteral = new DefaultVariableLiteral(name, argsContext, index, rest);
        futureStack.pop().whenComplete((typeReference, throwable) -> {
            if (throwable != null) {
                LOGGER.error("Error setting type", throwable);
                return;
            }
            if (typeReference.isSimpleTypeReference()) {
                variableLiteral.setTypeReference(typeReference);
            } else {
                variableLiteral.setTypeReference(typeReference);
                String typeName = typeReference.asNonSimpleTypeReference().orElseThrow().getName();
                lazyTypeDefinitionConsumer.setOrScheduleTypeDefinitionConsumer(typeName, variableLiteral::setTypeDefinition);
            }
        });
        if (futureStack.empty()) {
            futureStack = null;
        }

        parserContexts.peek().add(variableLiteral);
    }

    @Override
    public void enterIdentifierSegmentArguments(ExpressionParser.IdentifierSegmentArgumentsContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitIdentifierSegmentArguments(ExpressionParser.IdentifierSegmentArgumentsContext ctx) {
        List<Term> args = parserContexts.pop();
        parserContexts.peek().add(new ArgsContext(args));
    }

    @Override
    public void enterIdentifierSegmentIndexes(ExpressionParser.IdentifierSegmentIndexesContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitIdentifierSegmentIndexes(ExpressionParser.IdentifierSegmentIndexesContext ctx) {
        List<Term> args = parserContexts.pop();
        List<NumericLiteral> numericLiterals = args.stream().map(NumericLiteral.class::cast).collect(Collectors.toList());
        parserContexts.peek().add(new IndexContext(numericLiterals));
    }

    @Override
    public void enterIdentifierSegmentRest(ExpressionParser.IdentifierSegmentRestContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitIdentifierSegmentRest(ExpressionParser.IdentifierSegmentRestContext ctx) {
        List<Term> args = parserContexts.pop();
        List<VariableLiteral> variableLiterals = args.stream().map(VariableLiteral.class::cast).collect(Collectors.toList());
        parserContexts.peek().add(new RestContext(variableLiterals));
    }

    /////////////////////////////////////////////////////////////////////////////////////////
    // Unary Terms
    /////////////////////////////////////////////////////////////////////////////////////////

    @Override
    public void enterNotExpression(ExpressionParser.NotExpressionContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitNotExpression(ExpressionParser.NotExpressionContext ctx) {
        UnaryTerm ut = getUnaryTerm("!", parserContexts.pop());
        parserContexts.peek().add(ut);
    }

    @Override
    public void enterUnaryMinusExpression(ExpressionParser.UnaryMinusExpressionContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitUnaryMinusExpression(ExpressionParser.UnaryMinusExpressionContext ctx) {
        UnaryTerm ut = getUnaryTerm("-", parserContexts.pop());
        parserContexts.peek().add(ut);
    }

    @Override
    public void enterExpressionExpression(ExpressionParser.ExpressionExpressionContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitExpressionExpression(ExpressionParser.ExpressionExpressionContext ctx) {
        UnaryTerm ut = getUnaryTerm("()", parserContexts.pop());
        parserContexts.peek().add(ut);
    }

    /////////////////////////////////////////////////////////////////////////////////////////
    // Binary Terms
    /////////////////////////////////////////////////////////////////////////////////////////

    @Override
    public void enterOrExpression(ExpressionParser.OrExpressionContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitOrExpression(ExpressionParser.OrExpressionContext ctx) {
        BinaryTerm bt = getBinaryTerm("||", parserContexts.pop());
        parserContexts.peek().add(bt);
    }

    @Override
    public void enterPowerExpression(ExpressionParser.PowerExpressionContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitPowerExpression(ExpressionParser.PowerExpressionContext ctx) {
        BinaryTerm bt = getBinaryTerm("^", parserContexts.pop());
        parserContexts.peek().add(bt);
    }

    @Override
    public void enterEqExpression(ExpressionParser.EqExpressionContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitEqExpression(ExpressionParser.EqExpressionContext ctx) {
        BinaryTerm bt = getBinaryTerm(ctx.op.getText(), parserContexts.pop());
        parserContexts.peek().add(bt);
    }

    @Override
    public void enterAndExpression(ExpressionParser.AndExpressionContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitAndExpression(ExpressionParser.AndExpressionContext ctx) {
        BinaryTerm bt = getBinaryTerm("&&", parserContexts.pop());
        parserContexts.peek().add(bt);
    }

    @Override
    public void enterAddExpression(ExpressionParser.AddExpressionContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitAddExpression(ExpressionParser.AddExpressionContext ctx) {
        BinaryTerm bt = getBinaryTerm(ctx.op.getText(), parserContexts.pop());
        parserContexts.peek().add(bt);
    }

    @Override
    public void enterBitShiftExpression(ExpressionParser.BitShiftExpressionContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitBitShiftExpression(ExpressionParser.BitShiftExpressionContext ctx) {
        BinaryTerm bt = getBinaryTerm(ctx.op.getText(), parserContexts.pop());
        parserContexts.peek().add(bt);
    }

    @Override
    public void enterBitAndExpression(ExpressionParser.BitAndExpressionContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitBitAndExpression(ExpressionParser.BitAndExpressionContext ctx) {
        BinaryTerm bt = getBinaryTerm("&", parserContexts.pop());
        parserContexts.peek().add(bt);
    }

    @Override
    public void enterBitOrExpression(ExpressionParser.BitOrExpressionContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitBitOrExpression(ExpressionParser.BitOrExpressionContext ctx) {
        BinaryTerm bt = getBinaryTerm("|", parserContexts.pop());
        parserContexts.peek().add(bt);
    }

    @Override
    public void enterCompExpression(ExpressionParser.CompExpressionContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitCompExpression(ExpressionParser.CompExpressionContext ctx) {
        BinaryTerm bt = getBinaryTerm(ctx.op.getText(), parserContexts.pop());
        parserContexts.peek().add(bt);
    }

    @Override
    public void enterMultExpression(ExpressionParser.MultExpressionContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitMultExpression(ExpressionParser.MultExpressionContext ctx) {
        BinaryTerm bt = getBinaryTerm(ctx.op.getText(), parserContexts.pop());
        parserContexts.peek().add(bt);
    }

    /////////////////////////////////////////////////////////////////////////////////////////
    // Ternary Terms
    /////////////////////////////////////////////////////////////////////////////////////////

    @Override
    public void enterIfExpression(ExpressionParser.IfExpressionContext ctx) {
        parserContexts.push(new LinkedList<>());
    }

    @Override
    public void exitIfExpression(ExpressionParser.IfExpressionContext ctx) {
        TernaryTerm tt = getTernaryTerm("if", parserContexts.pop());
        parserContexts.peek().add(tt);
    }

    /////////////////////////////////////////////////////////////////////////////////////////
    // Helpers
    /////////////////////////////////////////////////////////////////////////////////////////

    private UnaryTerm getUnaryTerm(String op, List<Term> terms) {
        if (terms.size() != 1) {
            throw new RuntimeException(op + " should be a unary operation");
        }
        Term a = terms.get(0);
        return new DefaultUnaryTerm(a, op);
    }

    private BinaryTerm getBinaryTerm(String op, List<Term> terms) {
        if (terms.size() != 2) {
            throw new RuntimeException(op + " should be a binary operation");
        }
        Term a = terms.get(0);
        Term b = terms.get(1);
        return new DefaultBinaryTerm(a, b, op);
    }

    private TernaryTerm getTernaryTerm(String op, List<Term> terms) {
        if (terms.size() != 3) {
            throw new RuntimeException(op + " should be a ternary operation");
        }
        Term a = terms.get(0);
        Term b = terms.get(1);
        Term c = terms.get(2);
        return new DefaultTernaryTerm(a, b, c, op);
    }

    static class ArgsContext extends LinkedList<Term> implements Term {
        ArgsContext(Collection<Term> c) {
            super(c);
        }

        @Override
        public boolean contains(String str) {
            return false;
        }

        @Override
        public String stringRepresentation() {
            return "";
        }
    }

    static class IndexContext extends LinkedList<NumericLiteral> implements Term {
        IndexContext(Collection<NumericLiteral> c) {
            super(c);
        }

        @Override
        public boolean contains(String str) {
            return false;
        }

        @Override
        public String stringRepresentation() {
            return "";
        }
    }

    static class RestContext extends LinkedList<VariableLiteral> implements Term {
        RestContext(Collection<VariableLiteral> c) {
            super(c);
        }

        @Override
        public boolean contains(String str) {
            return false;
        }

        @Override
        public String stringRepresentation() {
            return "";
        }
    }

}
