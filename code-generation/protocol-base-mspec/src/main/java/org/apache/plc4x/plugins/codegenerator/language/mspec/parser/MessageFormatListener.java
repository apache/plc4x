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

import org.apache.commons.io.IOUtils;
import org.apache.plc4x.plugins.codegenerator.language.mspec.MSpecBaseListener;
import org.apache.plc4x.plugins.codegenerator.language.mspec.MSpecParser;
import org.apache.plc4x.plugins.codegenerator.language.mspec.expression.ExpressionStringParser;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions.*;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.fields.*;
import org.apache.plc4x.plugins.codegenerator.types.definitions.Argument;
import org.apache.plc4x.plugins.codegenerator.types.definitions.DefaultArgument;
import org.apache.plc4x.plugins.codegenerator.types.definitions.DiscriminatedComplexTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.enums.EnumValue;
import org.apache.plc4x.plugins.codegenerator.types.fields.ArrayField;
import org.apache.plc4x.plugins.codegenerator.types.fields.Field;
import org.apache.plc4x.plugins.codegenerator.types.fields.ManualArrayField;
import org.apache.plc4x.plugins.codegenerator.types.fields.SwitchField;
import org.apache.plc4x.plugins.codegenerator.types.references.*;
import org.apache.plc4x.plugins.codegenerator.types.terms.DefaultNumericLiteral;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;

import java.io.InputStream;
import java.nio.charset.Charset;
import java.util.*;
import java.util.stream.Collectors;

public class MessageFormatListener extends MSpecBaseListener {

    private Deque<List<Field>> parserContexts;

    private Deque<List<EnumValue>> enumContexts;

    private Map<String, TypeDefinition> types;

    public Deque<List<Field>> getParserContexts() {
        return parserContexts;
    }

    public Deque<List<EnumValue>> getEnumContexts() {
        return enumContexts;
    }

    public Map<String, TypeDefinition> getTypes() {
        return types;
    }

    @Override
    public void enterFile(MSpecParser.FileContext ctx) {
        parserContexts = new LinkedList<>();
        enumContexts = new LinkedList<>();
        types = new HashMap<>();
    }

    @Override
    public void enterComplexType(MSpecParser.ComplexTypeContext ctx) {
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

        // Handle enum types.
        if (ctx.enumValues != null) {
            TypeReference type = (ctx.type != null) ? getTypeReference(ctx.type) : null;
            List<EnumValue> enumValues = getEnumValues();
            if (type == null) {
                // in case there is no type we default to uint32
                type = new DefaultIntegerTypeReference(SimpleTypeReference.SimpleBaseType.UINT, 32);
            }
            DefaultEnumTypeDefinition enumType = new DefaultEnumTypeDefinition(typeName, type, enumValues,
                parserArguments, null);
            types.put(typeName, enumType);
            enumContexts.pop();
        } else if (ctx.dataIoTypeSwitch != null) {  // Handle data-io types.
            SwitchField switchField = getSwitchField();
            DefaultDataIoTypeDefinition type = new DefaultDataIoTypeDefinition(
                typeName, parserArguments, null, switchField);
            types.put(typeName, type);

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
            DefaultComplexTypeDefinition type = new DefaultComplexTypeDefinition(typeName, parserArguments, null, abstractType, parserContexts.peek());
            types.put(typeName, type);

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
    public void enterAbstractField(MSpecParser.AbstractFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = getIdString(ctx.name);
        MSpecParser.FieldDefinitionContext fieldDefinitionContext = (MSpecParser.FieldDefinitionContext) ctx.parent.parent;
        List<Term> params = getFieldParams(fieldDefinitionContext);
        Field field = new DefaultAbstractField(null, type, name, params);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterArrayField(MSpecParser.ArrayFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = getIdString(ctx.name);
        ArrayField.LoopType loopType = ArrayField.LoopType.valueOf(ctx.loopType.getText().toUpperCase());
        String loopExpressionString = getExprString(ctx.loopExpression);
        Term loopExpression = getExpressionTerm(loopExpressionString);
        MSpecParser.FieldDefinitionContext fieldDefinitionContext = (MSpecParser.FieldDefinitionContext) ctx.parent.parent;
        List<Term> params = getFieldParams(fieldDefinitionContext);
        Field field = new DefaultArrayField(null, type, name, loopType, loopExpression, params);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterChecksumField(MSpecParser.ChecksumFieldContext ctx) {
        SimpleTypeReference type = getSimpleTypeReference(ctx.type);
        String name = getIdString(ctx.name);
        String checksumExpressionString = getExprString(ctx.checksumExpression);
        Term checksumExpression = getExpressionTerm(checksumExpressionString);
        Field field = new DefaultChecksumField(null, type, name, checksumExpression);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterConstField(MSpecParser.ConstFieldContext ctx) {
        TypeReference type = ctx.type.dataType() != null ? getSimpleTypeReference(ctx.type.dataType()) : getTypeReference(ctx.type);
        String name = getIdString(ctx.name);
        String expected = getExprString(ctx.expected);
        Field field = new DefaultConstField(null, type, name, expected);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterDiscriminatorField(MSpecParser.DiscriminatorFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = getIdString(ctx.name);
        Field field = new DefaultDiscriminatorField(null, type, name);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterEnumField(MSpecParser.EnumFieldContext ctx) {
        ComplexTypeReference type = new DefaultComplexTypeReference(ctx.type.complexTypeReference.getText(), null);
        String name = getIdString(ctx.name);
        String fieldName = null;
        if (ctx.fieldName != null) {
            fieldName = getIdString(ctx.fieldName);
        }
        MSpecParser.FieldDefinitionContext fieldDefinitionContext = (MSpecParser.FieldDefinitionContext) ctx.parent.parent;
        List<Term> params = getFieldParams(fieldDefinitionContext);
        Field field = new DefaultEnumField(null, type, name, fieldName, params);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterImplicitField(MSpecParser.ImplicitFieldContext ctx) {
        SimpleTypeReference type = getSimpleTypeReference(ctx.type);
        String name = getIdString(ctx.name);
        String serializeExpressionString = getExprString(ctx.serializeExpression);
        Term serializeExpression = getExpressionTerm(serializeExpressionString);
        Field field = new DefaultImplicitField(null, type, name, serializeExpression);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterAssertField(MSpecParser.AssertFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = getIdString(ctx.name);
        String conditionExpressionString = getExprString(ctx.condition);
        Term conditionExpression = getExpressionTerm(conditionExpressionString);
        MSpecParser.FieldDefinitionContext fieldDefinitionContext = (MSpecParser.FieldDefinitionContext) ctx.parent.parent;
        List<Term> params = getFieldParams(fieldDefinitionContext);
        Field field = new DefaultAssertField(null, type, name, conditionExpression, params);
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
        String loopExpressionString = getExprString(ctx.loopExpression);
        Term loopExpression = getExpressionTerm(loopExpressionString);
        String parseExpressionString = getExprString(ctx.parseExpression);
        Term parseExpression = getExpressionTerm(parseExpressionString);
        String serializeExpressionString = getExprString(ctx.serializeExpression);
        Term serializeExpression = getExpressionTerm(serializeExpressionString);
        String lengthExpressionString = getExprString(ctx.lengthExpression);
        Term lengthExpression = getExpressionTerm(lengthExpressionString);
        MSpecParser.FieldDefinitionContext fieldDefinitionContext = (MSpecParser.FieldDefinitionContext) ctx.parent.parent;
        List<Term> params = getFieldParams(fieldDefinitionContext);
        Field field = new DefaultManualArrayField(null, type, name, loopType, loopExpression,
            parseExpression, serializeExpression, lengthExpression, params);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterManualField(MSpecParser.ManualFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = getIdString(ctx.name);
        String parseExpressionString = getExprString(ctx.parseExpression);
        Term parseExpression = getExpressionTerm(parseExpressionString);
        String serializeExpressionString = getExprString(ctx.serializeExpression);
        Term serializeExpression = getExpressionTerm(serializeExpressionString);
        String lengthExpressionString = getExprString(ctx.lengthExpression);
        Term lengthExpression = getExpressionTerm(lengthExpressionString);
        MSpecParser.FieldDefinitionContext fieldDefinitionContext = (MSpecParser.FieldDefinitionContext) ctx.parent.parent;
        List<Term> params = getFieldParams(fieldDefinitionContext);
        Field field = new DefaultManualField(null, type, name, parseExpression, serializeExpression,
            lengthExpression, params);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterOptionalField(MSpecParser.OptionalFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = getIdString(ctx.name);
        String conditionExpressionString = getExprString(ctx.condition);
        Term conditionExpression = getExpressionTerm(conditionExpressionString);
        MSpecParser.FieldDefinitionContext fieldDefinitionContext = (MSpecParser.FieldDefinitionContext) ctx.parent.parent;
        List<Term> params = getFieldParams(fieldDefinitionContext);
        Field field = new DefaultOptionalField(null, type, name, conditionExpression, params);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterPaddingField(MSpecParser.PaddingFieldContext ctx) {
        SimpleTypeReference type = getSimpleTypeReference(ctx.type);
        String name = getIdString(ctx.name);
        String paddingValueString = getExprString(ctx.paddingValue);
        Term paddingValue = getExpressionTerm(paddingValueString);
        String paddingConditionString = getExprString(ctx.paddingCondition);
        Term paddingCondition = getExpressionTerm(paddingConditionString);
        MSpecParser.FieldDefinitionContext fieldDefinitionContext = (MSpecParser.FieldDefinitionContext) ctx.parent.parent;
        List<Term> params = getFieldParams(fieldDefinitionContext);
        Field field = new DefaultPaddingField(null, type, name, paddingValue, paddingCondition, params);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterReservedField(MSpecParser.ReservedFieldContext ctx) {
        SimpleTypeReference type = getSimpleTypeReference(ctx.type);
        String expected = getExprString(ctx.expected);
        Field field = new DefaultReservedField(null, type, expected);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterSimpleField(MSpecParser.SimpleFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = getIdString(ctx.name);
        MSpecParser.FieldDefinitionContext fieldDefinitionContext = (MSpecParser.FieldDefinitionContext) ctx.parent.parent;
        List<Term> params = getFieldParams(fieldDefinitionContext);
        Field field = new DefaultSimpleField(null, type, name, params);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterTypeSwitchField(MSpecParser.TypeSwitchFieldContext ctx) {
        List<Term> discriminatorExpressions = ctx.discriminators.expression().stream()
            .map(this::getExprString)
            .map(this::getExpressionTerm)
            .collect(Collectors.toList());
        DefaultSwitchField field = new DefaultSwitchField(false, discriminatorExpressions);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterUnknownField(MSpecParser.UnknownFieldContext ctx) {
        SimpleTypeReference type = getSimpleTypeReference(ctx.type);
        Field field = new DefaultUnknownField(null, type);
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterVirtualField(MSpecParser.VirtualFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = getIdString(ctx.name);
        String valueExpressionString = getExprString(ctx.valueExpression);
        Term valueExpression = getExpressionTerm(valueExpressionString);
        Field field = new DefaultVirtualField(null, type, name, valueExpression);
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

        List<String> discriminatorValues;
        if (ctx.discriminatorValues != null) {
            discriminatorValues = ctx.discriminatorValues.expression().stream()
                .map(this::getExprString)
                .collect(Collectors.toList());
        } else {
            discriminatorValues = Collections.emptyList();
        }
        DefaultDiscriminatedComplexTypeDefinition type =
            new DefaultDiscriminatedComplexTypeDefinition(typeName, parserArguments, null,
                discriminatorValues, parserContexts.pop());

        // Add the type to the switch field definition.
        DefaultSwitchField switchField = getSwitchField();
        if (switchField == null) {
            throw new RuntimeException("This shouldn't have happened");
        }
        switchField.addCase(type);

        // Add the type to the type list.
        types.put(typeName, type);
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
                constantName = constantName.substring(1, constantName.length() - 1);
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

    private Term getExpressionTerm(String expressionString) {
        InputStream inputStream = IOUtils.toInputStream(expressionString, Charset.defaultCharset());
        ExpressionStringParser parser = new ExpressionStringParser();
        try {
            return parser.parse(inputStream);
        } catch (Exception e) {
            throw new RuntimeException("Error parsing expression: '" + expressionString + "'", e);
        }
    }

    private TypeReference getTypeReference(MSpecParser.TypeReferenceContext ctx) {
        if (ctx.simpleTypeReference != null) {
            return getSimpleTypeReference(ctx.simpleTypeReference);
        } else {
            return new DefaultComplexTypeReference(ctx.complexTypeReference.getText(), getParams(ctx.params));
        }
    }

    private SimpleTypeReference getSimpleTypeReference(MSpecParser.DataTypeContext ctx) {
        SimpleTypeReference.SimpleBaseType simpleBaseType =
            SimpleTypeReference.SimpleBaseType.valueOf(ctx.base.getText().toUpperCase());
        // String types need an additional "encoding" field and length expression.
        if ((simpleBaseType == SimpleTypeReference.SimpleBaseType.STRING) ||
                (simpleBaseType == SimpleTypeReference.SimpleBaseType.VSTRING)) {
            if (simpleBaseType == SimpleTypeReference.SimpleBaseType.VSTRING) {
                Term lengthExpression = getExpressionTerm(ctx.length.getText().substring(1, ctx.length.getText().length() - 1));
                return new DefaultStringTypeReference(simpleBaseType, lengthExpression, "UTF-8");
            } else {
                int size = Integer.parseInt(ctx.size.getText());
                return new DefaultStringTypeReference(simpleBaseType, new DefaultNumericLiteral(size), "UTF-8");
            }
        }
        // If a size it specified it's a simple integer length based type.
        if (ctx.size != null) {
            int size = Integer.parseInt(ctx.size.getText());
            return new DefaultIntegerTypeReference(simpleBaseType, size);
        }
        // If exponent and mantissa are present, it's a floating point representation.
        else if ((ctx.exponent != null) && (ctx.mantissa != null)) {
            int exponent = Integer.parseInt(ctx.exponent.getText());
            int mantissa = Integer.parseInt(ctx.mantissa.getText());
            return new DefaultFloatTypeReference(simpleBaseType, exponent, mantissa);
        } else if ((simpleBaseType == SimpleTypeReference.SimpleBaseType.TIME) ||
            (simpleBaseType == SimpleTypeReference.SimpleBaseType.DATE) ||
            (simpleBaseType == SimpleTypeReference.SimpleBaseType.DATETIME)) {
            return new DefaultTemporalTypeReference(simpleBaseType);
        } else if (simpleBaseType == SimpleTypeReference.SimpleBaseType.BIT) {
            return new DefaultBooleanTypeReference();
        } else if (simpleBaseType == SimpleTypeReference.SimpleBaseType.BYTE) {
            return new DefaultByteTypeReference();
        }
        // In all other cases (bit) it's just assume it's length it 1.
        else {
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

    private List<Term> getFieldParams(MSpecParser.FieldDefinitionContext parentCtx) {
        return getParams(parentCtx.params);
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
        ExpressionStringParser parser = new ExpressionStringParser();
        try {
            return parser.parse(inputStream);
        } catch (Exception e) {
            throw new RuntimeException("Error parsing expression: '" + expressionString + "'", e);
        }
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

}
