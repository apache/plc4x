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
package org.apache.plc4x.plugins.codegenerator.language.mspec.parser;

import org.antlr.v4.runtime.RuleContext;
import org.apache.commons.io.IOUtils;
import org.apache.plc4x.plugins.codegenerator.language.mspec.LazyTypeDefinitionConsumer;
import org.apache.plc4x.plugins.codegenerator.language.mspec.MSpecBaseListener;
import org.apache.plc4x.plugins.codegenerator.language.mspec.MSpecParser;
import org.apache.plc4x.plugins.codegenerator.language.mspec.expression.ExpressionStringParser;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions.*;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.fields.*;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.terms.WildcardTerm;
import org.apache.plc4x.plugins.codegenerator.types.definitions.*;
import org.apache.plc4x.plugins.codegenerator.types.enums.EnumValue;
import org.apache.plc4x.plugins.codegenerator.types.fields.ArrayField;
import org.apache.plc4x.plugins.codegenerator.types.fields.Field;
import org.apache.plc4x.plugins.codegenerator.types.fields.ManualArrayField;
import org.apache.plc4x.plugins.codegenerator.types.fields.SwitchField;
import org.apache.plc4x.plugins.codegenerator.types.references.*;
import org.apache.plc4x.plugins.codegenerator.types.terms.Literal;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;
import org.apache.plc4x.plugins.codegenerator.types.terms.VariableLiteral;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.InputStream;
import java.nio.charset.Charset;
import java.util.*;
import java.util.function.Consumer;
import java.util.stream.Collectors;

public class MessageFormatListener extends MSpecBaseListener implements LazyTypeDefinitionConsumer {

    private static final Logger LOGGER = LoggerFactory.getLogger(MessageFormatListener.class);

    private Deque<List<Field>> parserContexts;

    private Deque<List<EnumValue>> enumContexts;

    protected Map<String, TypeDefinition> types;

    protected Map<String, List<Consumer<TypeDefinition>>> typeDefinitionConsumers = new HashMap<>();

    private Stack<Map<String, Term>> batchSetAttributes = new Stack<>();

    public Deque<List<Field>> getParserContexts() {
        return parserContexts;
    }

    public Deque<List<EnumValue>> getEnumContexts() {
        return enumContexts;
    }

    private String currentTypeName;

    @Override
    public void enterFile(MSpecParser.FileContext ctx) {
        parserContexts = new LinkedList<>();
        enumContexts = new LinkedList<>();
        types = new HashMap<>();
    }

    @Override
    public void enterComplexType(MSpecParser.ComplexTypeContext ctx) {
        currentTypeName = getIdString(ctx.name);
        // Set a map of attributes that should be set for all fields.
        Map<String, Term> curBatchSetAttributes = new HashMap<>();
        // Add all attributes defined in the current batchSet field.
        for (MSpecParser.AttributeContext attributeContext : ctx.attributes.attribute()) {
            Term attributeExpression = getExpressionTerm(attributeContext.value);
            curBatchSetAttributes.put(attributeContext.name.getText(), attributeExpression);
        }
        // Make the new Map the top of the stack.
        batchSetAttributes.push(curBatchSetAttributes);

        if (ctx.enumValues != null) {
            List<EnumValue> enumContext = new LinkedList<>();
            enumContexts.push(enumContext);
        } else {
            List<Field> parserContext = new LinkedList<>();
            parserContexts.push(parserContext);
        }
    }

    @Override
    public void exitComplexType(MSpecParser.ComplexTypeContext ctx) {
        String typeName = getIdString(ctx.name);
        List<Argument> parserArguments = null;
        if (ctx.params != null) {
            parserArguments = getParserArguments(ctx.params.argument());
        }

        final Map<String, Term> attributes = batchSetAttributes.peek();
        // Handle enum types.
        if (ctx.enumValues != null) {
            TypeReference type = (ctx.type != null) ? getTypeReference(ctx.type) : null;
            List<EnumValue> enumValues = getEnumValues();
            if (type == null) {
                // in case there is no type we default to uint32
                type = new DefaultIntegerTypeReference(SimpleTypeReference.SimpleBaseType.UINT, 32);
            }
            DefaultEnumTypeDefinition enumType = new DefaultEnumTypeDefinition(typeName, type, attributes, enumValues,
                parserArguments);
            dispatchType(typeName, enumType);
            enumContexts.pop();
        } else if (ctx.dataIoTypeSwitch != null) {  // Handle data-io types.
            SwitchField switchField = getSwitchField();
            DefaultDataIoTypeDefinition type = new DefaultDataIoTypeDefinition(typeName, attributes, parserArguments, switchField);
            dispatchType(typeName, type);

            // Set the parent type for all sub-types.
            if (switchField != null) {
                for (DiscriminatedComplexTypeDefinition subtype : switchField.getCases()) {
                    if (subtype instanceof DefaultDiscriminatedComplexTypeDefinition) {
                        ((DefaultDiscriminatedComplexTypeDefinition) subtype).setParentType(type);
                    }
                }
            }
            parserContexts.pop();
        } else { // Handle all other types.
            // If the type has sub-types it's an abstract type.
            SwitchField switchField = getSwitchField();
            boolean abstractType = switchField != null;
            DefaultComplexTypeDefinition type = new DefaultComplexTypeDefinition(
                typeName, attributes, parserArguments, abstractType, parserContexts.peek());
            dispatchType(typeName, type);

            // Set the parent type for all sub-types.
            if (switchField != null) {
                for (DiscriminatedComplexTypeDefinition subtype : switchField.getCases()) {
                    if (subtype instanceof DefaultDiscriminatedComplexTypeDefinition) {
                        ((DefaultDiscriminatedComplexTypeDefinition) subtype).setParentType(type);
                    }
                }
            }
            parserContexts.pop();
        }
    }

    @Override
    public void enterBatchSetDefinition(MSpecParser.BatchSetDefinitionContext ctx) {
        // Set a map of attributes that should be set for all fields.
        Map<String, Term> curBatchSetAttributes = new HashMap<>();
        // Add all attributes of the lower layers and initialize the new map with it.
        if (!batchSetAttributes.empty()) {
            curBatchSetAttributes.putAll(batchSetAttributes.peek());
        }
        // Add all attributes defined in the current batchSet field.
        for (MSpecParser.AttributeContext attributeContext : ctx.attributes.attribute()) {
            Term attributeExpression = getExpressionTerm(attributeContext.value);
            curBatchSetAttributes.put(attributeContext.name.getText(), attributeExpression);
        }
        // Make the new Map the top of the stack.
        batchSetAttributes.push(curBatchSetAttributes);
    }

    @Override
    public void exitBatchSetDefinition(MSpecParser.BatchSetDefinitionContext ctx) {
        // Clear the map of attributes.
        batchSetAttributes.pop();
    }

    @Override
    public void enterAbstractField(MSpecParser.AbstractFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = getIdString(ctx.name);
        Field field = new DefaultAbstractField(getAttributes(ctx), type, name);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterArrayField(MSpecParser.ArrayFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = getIdString(ctx.name);
        ArrayField.LoopType loopType = ArrayField.LoopType.valueOf(ctx.loopType.getText().toUpperCase());
        Term loopExpression = getExpressionTerm(ctx.loopExpression);
        Field field = new DefaultArrayField(getAttributes(ctx), type, name, loopType, loopExpression);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterChecksumField(MSpecParser.ChecksumFieldContext ctx) {
        SimpleTypeReference type = getSimpleTypeReference(ctx.type);
        String name = getIdString(ctx.name);
        Term checksumExpression = getExpressionTerm(ctx.checksumExpression);
        Field field = new DefaultChecksumField(getAttributes(ctx), type, name, checksumExpression);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterConstField(MSpecParser.ConstFieldContext ctx) {
        TypeReference type = ctx.type.dataType() != null ? getSimpleTypeReference(ctx.type.dataType()) : getTypeReference(ctx.type);
        String name = getIdString(ctx.name);
        Field field = new DefaultConstField(getAttributes(ctx), type, name, getValueLiteral(ctx.expected));
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterDiscriminatorField(MSpecParser.DiscriminatorFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = getIdString(ctx.name);
        Field field = new DefaultDiscriminatorField(getAttributes(ctx), type, name);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterEnumField(MSpecParser.EnumFieldContext ctx) {
        String typeRefName = ctx.type.complexTypeReference.getText();
        DefaultComplexTypeReference type = new DefaultComplexTypeReference(typeRefName, null);
        setOrScheduleTypeDefinitionConsumer(typeRefName, type::setTypeDefinition);
        String name = getIdString(ctx.name);
        String fieldName = null;
        if (ctx.fieldName != null) {
            fieldName = getIdString(ctx.fieldName);
        }
        Field field = new DefaultEnumField(getAttributes(ctx), type, name, fieldName);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterImplicitField(MSpecParser.ImplicitFieldContext ctx) {
        SimpleTypeReference type = getSimpleTypeReference(ctx.type);
        String name = getIdString(ctx.name);
        Term serializeExpression = getExpressionTerm(ctx.serializeExpression);
        Field field = new DefaultImplicitField(getAttributes(ctx), type, name, serializeExpression);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterAssertField(MSpecParser.AssertFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = getIdString(ctx.name);
        Term conditionExpression = getExpressionTerm(ctx.condition);
        Field field = new DefaultAssertField(getAttributes(ctx), type, name, conditionExpression);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterManualArrayField(MSpecParser.ManualArrayFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = getIdString(ctx.name);
        ManualArrayField.LoopType loopType = ManualArrayField.LoopType.valueOf(
            ctx.loopType.getText().toUpperCase());
        Term loopExpression = getExpressionTerm(ctx.loopExpression);
        Term parseExpression = getExpressionTerm(ctx.parseExpression);
        Term serializeExpression = getExpressionTerm(ctx.serializeExpression);
        Term lengthExpression = getExpressionTerm(ctx.lengthExpression);
        Field field = new DefaultManualArrayField(getAttributes(ctx), type, name, loopType, loopExpression,
            parseExpression, serializeExpression, lengthExpression);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterManualField(MSpecParser.ManualFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = getIdString(ctx.name);
        Term parseExpression = getExpressionTerm(ctx.parseExpression);
        Term serializeExpression = getExpressionTerm(ctx.serializeExpression);
        Term lengthExpression = getExpressionTerm(ctx.lengthExpression);
        Field field = new DefaultManualField(getAttributes(ctx), type, name, parseExpression, serializeExpression,
            lengthExpression);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterOptionalField(MSpecParser.OptionalFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = getIdString(ctx.name);
        Term conditionExpression = null;
        if (ctx.condition != null) {
            conditionExpression = getExpressionTerm(ctx.condition);
        }
        Field field = new DefaultOptionalField(getAttributes(ctx), type, name, conditionExpression);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterPeekField(MSpecParser.PeekFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = getIdString(ctx.name);
        Term offsetExpression = null;
        if (ctx.offset != null) {
            offsetExpression = getExpressionTerm(ctx.offset);
        }
        Field field = new DefaultPeekField(getAttributes(ctx), type, name, offsetExpression);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterPaddingField(MSpecParser.PaddingFieldContext ctx) {
        SimpleTypeReference type = getSimpleTypeReference(ctx.type);
        String name = getIdString(ctx.name);
        Term paddingValue = getExpressionTerm(ctx.paddingValue);
        Term paddingCondition = getExpressionTerm(ctx.paddingCondition);
        Field field = new DefaultPaddingField(getAttributes(ctx), type, name, paddingValue, paddingCondition);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterReservedField(MSpecParser.ReservedFieldContext ctx) {
        SimpleTypeReference type = getSimpleTypeReference(ctx.type);
        String expected = getExprString(ctx.expected);
        Field field = new DefaultReservedField(getAttributes(ctx), type, expected);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterSimpleField(MSpecParser.SimpleFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = getIdString(ctx.name);
        Field field = new DefaultSimpleField(getAttributes(ctx), type, name);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterTypeSwitchField(MSpecParser.TypeSwitchFieldContext ctx) {
        List<VariableLiteral> variableLiterals = ctx.discriminators.variableLiteral().stream()
            .map(this::getVariableLiteral)
            .collect(Collectors.toList());
        DefaultSwitchField field = new DefaultSwitchField(variableLiterals);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterUnknownField(MSpecParser.UnknownFieldContext ctx) {
        SimpleTypeReference type = getSimpleTypeReference(ctx.type);
        Field field = new DefaultUnknownField(getAttributes(ctx), type);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterVirtualField(MSpecParser.VirtualFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = getIdString(ctx.name);
        Term valueExpression = getExpressionTerm(ctx.valueExpression);
        Field field = new DefaultVirtualField(getAttributes(ctx), type, name, valueExpression);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterValidationField(MSpecParser.ValidationFieldContext ctx) {
        Term validationExpression = getExpressionTerm(ctx.validationExpression);
        Field field = new DefaultValidationField(validationExpression, ctx.description.getText());
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterCaseStatement(MSpecParser.CaseStatementContext ctx) {
        List<Field> parserContext = new LinkedList<>();
        parserContexts.push(parserContext);
    }

    @Override
    public void exitCaseStatement(MSpecParser.CaseStatementContext ctx) {
        String typeName = ctx.name.getText();

        final Map<String, Term> attributes = batchSetAttributes.peek();

        List<Argument> parserArguments = new LinkedList<>();
        // For DataIO types, add all the arguments from the parent type.
        if (!(ctx.parent.parent.parent.parent instanceof MSpecParser.ComplexTypeContext)
            && ((MSpecParser.ComplexTypeContext) ctx.parent.parent.parent).params != null) {
            parserArguments.addAll(getParserArguments(
                ((MSpecParser.ComplexTypeContext) ctx.parent.parent.parent).params.argument()));
        }
        // Add all eventually existing local arguments.
        if (ctx.argumentList() != null) {
            parserArguments.addAll(getParserArguments(ctx.argumentList().argument()));
        }

        List<Term> discriminatorValues;
        if (ctx.discriminatorValues != null) {
            discriminatorValues = ctx.discriminatorValues.expression().stream()
                .map(this::getExpressionTerm)
                .collect(Collectors.toList());
        } else {
            discriminatorValues = Collections.emptyList();
        }
        DefaultDiscriminatedComplexTypeDefinition type =
            new DefaultDiscriminatedComplexTypeDefinition(typeName, attributes, parserArguments,
                discriminatorValues, parserContexts.pop());
        dispatchType(typeName, type);

        // Add the type to the switch field definition.
        DefaultSwitchField switchField = getSwitchField();
        if (switchField == null) {
            throw new RuntimeException("This shouldn't have happened");
        }
        switchField.addCase(type);
    }

    @Override
    public void enterEnumValueDefinition(MSpecParser.EnumValueDefinitionContext ctx) {
        String value = (ctx.valueExpression != null) ? unquoteString(ctx.valueExpression.getText()) : null;
        String name = ctx.name.getText();
        Map<String, String> constants = null;
        if (ctx.constantValueExpressions != null) {
            MSpecParser.ComplexTypeContext parentCtx = (MSpecParser.ComplexTypeContext) ctx.parent;
            int numConstantValues = parentCtx.params.argument().size();
            int numExpressionValues = ctx.constantValueExpressions.expression().size();
            // This only works if we provide exactly the same number of expressions as we defined constants
            if (numConstantValues != numExpressionValues) {
                throw new RuntimeException("Number of constant value expressions doesn't match the number of " +
                    "defined constants. Expecting " + numConstantValues + " but got " + numExpressionValues);
            }

            // Build a map of the constant expressions (With the constant name as key)
            constants = new HashMap<>();
            for (int i = 0; i < numConstantValues; i++) {
                MSpecParser.ArgumentContext argumentContext = parentCtx.params.argument(i);
                String constantName = argumentContext.name.getText();
                MSpecParser.ExpressionContext expression = ctx.constantValueExpressions.expression(i);
                String constant = unquoteString(expression.getText());
                // String expressions are double escaped
                if (constant != null && constant.startsWith("\"")) {
                    constant = unquoteString(constant);
                }
                constants.put(constantName, constant);
            }
        }
        List<EnumValue> enumValues = Objects.requireNonNull(this.enumContexts.peek());
        if (value == null) {
            // If no values are specified we count
            String counted = "0";
            if (enumValues.size() > 0) {
                String previousValue = enumValues.get(enumValues.size() - 1).getValue();
                int parsedPreviousValue = Integer.parseInt(previousValue);
                counted = "" + (parsedPreviousValue + 1);
            }
            value = counted;
        }
        final DefaultEnumValue enumValue = new DefaultEnumValue(value, name, constants);
        enumValues.add(enumValue);
    }

    private Term getExpressionTerm(MSpecParser.ExpressionContext expressionContext) {
        if (expressionContext.ASTERISK() != null) {
            return WildcardTerm.INSTANCE;
        }
        String expressionString = getExprString(expressionContext);
        Objects.requireNonNull(expressionString, "Expression string should not be null");
        InputStream inputStream = IOUtils.toInputStream(expressionString, Charset.defaultCharset());

        Objects.requireNonNull(currentTypeName, "expression term can only occur within a type");
        ExpressionStringParser parser = new ExpressionStringParser(this, currentTypeName);
        try {
            return parser.parse(inputStream);
        } catch (Exception e) {
            throw new RuntimeException(String.format("Error parsing expression: '%s' at line %d column %d",
                expressionString, expressionContext.start.getLine(), expressionContext.start.getStartIndex()), e);
        }
    }

    private VariableLiteral getVariableLiteral(MSpecParser.VariableLiteralContext variableLiteralContext) {
        // TODO: make nullsafe
        final String variableLiteral = variableLiteralContext.getText();
        InputStream inputStream = IOUtils.toInputStream(variableLiteral, Charset.defaultCharset());
        ExpressionStringParser parser = new ExpressionStringParser(this, currentTypeName);
        try {
            // As this come from a VariableLiteralContext we know that it is a VariableLiteral
            return (VariableLiteral) parser.parse(inputStream);
        } catch (Exception e) {
            throw new RuntimeException(String.format("Error parsing variable literal: '%s' at line %d column %d",
                variableLiteral, variableLiteralContext.start.getLine(), variableLiteralContext.start.getStartIndex()), e);
        }
    }

    private Literal getValueLiteral(MSpecParser.ValueLiteralContext valueLiteralContext) {
        // TODO: make nullsafe
        final String valueLiteralContextText = valueLiteralContext.getText();
        InputStream inputStream = IOUtils.toInputStream(valueLiteralContextText, Charset.defaultCharset());
        ExpressionStringParser parser = new ExpressionStringParser(this, currentTypeName);
        try {
            // As this come from a ValueLiteralContext we know that it is a Literal
            return (Literal) parser.parse(inputStream);
        } catch (Exception e) {
            throw new RuntimeException(String.format("Error parsing variable literal: '%s' at line %d column %d",
                valueLiteralContextText, valueLiteralContext.start.getLine(), valueLiteralContext.start.getStartIndex()), e);
        }
    }

    private TypeReference getTypeReference(MSpecParser.TypeReferenceContext ctx) {
        if (ctx.simpleTypeReference != null) {
            return getSimpleTypeReference(ctx.simpleTypeReference);
        } else {
            String typeRefName = ctx.complexTypeReference.getText();
            DefaultComplexTypeReference type = new DefaultComplexTypeReference(typeRefName, getParams(ctx.params));
            setOrScheduleTypeDefinitionConsumer(typeRefName, type::setTypeDefinition);
            return type;
        }
    }

    private SimpleTypeReference getSimpleTypeReference(MSpecParser.DataTypeContext ctx) {
        SimpleTypeReference.SimpleBaseType simpleBaseType =
            SimpleTypeReference.SimpleBaseType.valueOf(ctx.base.getText().toUpperCase());
        // String types need an additional length expression.
        if (simpleBaseType == SimpleTypeReference.SimpleBaseType.VSTRING) {
            if (ctx.length != null) {
                Term lengthExpression = getExpressionTerm(ctx.length);
                return new DefaultVstringTypeReference(simpleBaseType, lengthExpression);
            } else {
                return new DefaultVstringTypeReference(simpleBaseType, null);
            }
        }
        switch (simpleBaseType) {
            case INT:
            case UINT:
                int integerSize = Integer.parseInt(ctx.size.getText());
                return new DefaultIntegerTypeReference(simpleBaseType, integerSize);
            case FLOAT:
            case UFLOAT:
                int floatSize = Integer.parseInt(ctx.size.getText());
                return new DefaultFloatTypeReference(simpleBaseType, floatSize);
            case TIME:
            case DATE:
            case DATETIME:
                return new DefaultTemporalTypeReference(simpleBaseType);
            case BIT:
                return new DefaultBooleanTypeReference();
            case BYTE:
                return new DefaultByteTypeReference();
            case STRING:
                int stringSize = Integer.parseInt(ctx.size.getText());
                return new DefaultStringTypeReference(simpleBaseType, stringSize);
            default:
                return new DefaultIntegerTypeReference(simpleBaseType, 1);
        }
    }

    private DefaultSwitchField getSwitchField() {
        for (Field field : Objects.requireNonNull(parserContexts.peek())) {
            if (field instanceof DefaultSwitchField) {
                return (DefaultSwitchField) field;
            }
        }
        return null;
    }

    private List<EnumValue> getEnumValues() {
        return Objects.requireNonNull(enumContexts.peek());
    }

    private List<Argument> getParserArguments(List<MSpecParser.ArgumentContext> params) {
        return params.stream()
            .map(argumentContext -> new DefaultArgument(getTypeReference(argumentContext.type), getIdString(argumentContext.name)))
            .collect(Collectors.toList());
    }

    private List<Term> getParams(MSpecParser.MultipleExpressionsContext params) {
        if (params == null) {
            return null;
        }
        return params.expression().stream()
            .map(this::getExprString)
            .map(this::parseExpression)
            .collect(Collectors.toList());
    }

    private Term parseExpression(String expressionString) {
        InputStream inputStream = IOUtils.toInputStream(expressionString, Charset.defaultCharset());
        ExpressionStringParser parser = new ExpressionStringParser(this, currentTypeName);
        try {
            Term term = parser.parse(inputStream);
            return term;
        } catch (Exception e) {
            throw new RuntimeException("Error parsing expression: '" + expressionString + "'", e);
        }
    }

    private Map<String, Term> getAttributes(RuleContext ctx) {
        Map<String, Term> attributes = new HashMap<>();
        // Add any attributes from the current batchSet context.
        if (!batchSetAttributes.empty()) {
            attributes.putAll(batchSetAttributes.peek());
        }
        // Add any attributes set on the current field itself.
        if (ctx.parent.parent instanceof MSpecParser.FieldDefinitionContext) {
            MSpecParser.FieldDefinitionContext fieldDefinitionContext = (MSpecParser.FieldDefinitionContext) ctx.parent.parent;
            for (MSpecParser.AttributeContext attributeContext : fieldDefinitionContext.attributes.attribute()) {
                Term attributeExpression = getExpressionTerm(attributeContext.value);
                attributes.put(attributeContext.name.getText(), attributeExpression);
            }
        }
        return attributes;
    }

    private String unquoteString(String quotedString) {
        if (quotedString != null && quotedString.length() >= 2) {
            return quotedString.substring(1, quotedString.length() - 1);
        }
        return quotedString;
    }

    private String getIdString(MSpecParser.IdExpressionContext ctx) {
        if (ctx.id != null) {
            return ctx.id.getText();
        }
        return null;
    }

    private String getExprString(MSpecParser.ExpressionContext ctx) {
        if (ctx.expr != null) {
            return ctx.expr.getText();
        }
        return null;
    }

    public void dispatchType(String typeName, TypeDefinition type) {
        LOGGER.debug("dispatching {}:{}", typeName, type);
        List<Consumer<TypeDefinition>> waitingConsumers = typeDefinitionConsumers.getOrDefault(typeName, new LinkedList<>());
        LOGGER.debug("{} waiting for {}", waitingConsumers.size(), typeName);
        Iterator<Consumer<TypeDefinition>> consumerIterator = waitingConsumers.iterator();
        while (consumerIterator.hasNext()) {
            Consumer<TypeDefinition> setter = consumerIterator.next();
            LOGGER.debug("setting {} for {}", typeName, setter);
            setter.accept(type);
            consumerIterator.remove();
        }
        typeDefinitionConsumers.remove(typeName);
        types.put(typeName, type);
    }

    @Override
    public void setOrScheduleTypeDefinitionConsumer(String typeRefName, Consumer<TypeDefinition> setTypeDefinition) {
        LOGGER.debug("set or schedule {}", typeRefName);
        TypeDefinition typeDefinition = types.get(typeRefName);
        if (typeDefinition != null) {
            LOGGER.debug("{} present so setting for {}", typeRefName, setTypeDefinition);
            setTypeDefinition.accept(typeDefinition);
        } else {
            // put up order
            if (LOGGER.isDebugEnabled()) {
                LOGGER.debug("{} already waiting for {}", typeDefinitionConsumers.getOrDefault(typeRefName, new LinkedList<>()).size(), typeRefName);
            }
            typeDefinitionConsumers.putIfAbsent(typeRefName, new LinkedList<>());
            typeDefinitionConsumers.get(typeRefName).add(setTypeDefinition);
        }
    }

}
