/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */

package org.apache.plc4x.plugins.codegenerator.language.mspec.parser;

import org.apache.commons.io.IOUtils;
import org.apache.plc4x.plugins.codegenerator.language.mspec.MSpecBaseListener;
import org.apache.plc4x.plugins.codegenerator.language.mspec.MSpecParser;
import org.apache.plc4x.plugins.codegenerator.language.mspec.expression.ExpressionStringParser;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions.DefaultComplexTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions.DefaultDiscriminatedComplexTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions.DefaultEnumTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions.DefaultEnumValue;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.fields.*;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.references.DefaultComplexTypeReference;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.references.DefaultSimpleTypeReference;
import org.apache.plc4x.plugins.codegenerator.types.definitions.Argument;
import org.apache.plc4x.plugins.codegenerator.types.definitions.DiscriminatedComplexTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.enums.EnumValue;
import org.apache.plc4x.plugins.codegenerator.types.fields.ArrayField;
import org.apache.plc4x.plugins.codegenerator.types.fields.Field;
import org.apache.plc4x.plugins.codegenerator.types.fields.ManualArrayField;
import org.apache.plc4x.plugins.codegenerator.types.fields.SwitchField;
import org.apache.plc4x.plugins.codegenerator.types.references.ComplexTypeReference;
import org.apache.plc4x.plugins.codegenerator.types.references.SimpleTypeReference;
import org.apache.plc4x.plugins.codegenerator.types.references.TypeReference;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;

import java.io.InputStream;
import java.nio.charset.Charset;
import java.util.*;

public class MessageFormatListener extends MSpecBaseListener {

    private Deque<List<Field>> parserContexts;

    private Deque<List<EnumValue>> enumContexts;

    private Map<String, TypeDefinition> types;

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
        if(ctx.enumValues != null) {
            List<EnumValue> enumContext = new LinkedList<>();
            enumContexts.push(enumContext);
        } else {
            List<Field> parserContext = new LinkedList<>();
            parserContexts.push(parserContext);
        }
    }

    @Override
    public void exitComplexType(MSpecParser.ComplexTypeContext ctx) {
        String typeName = ctx.name.id.getText();
        Argument[] parserArguments = null;
        if(ctx.params != null) {
            parserArguments = getParserArguments(ctx.params.argument());
        }

        if(ctx.enumValues != null) {
            TypeReference type = getTypeReference(ctx.type);
            EnumValue[] enumValues = getEnumValues();
            DefaultEnumTypeDefinition enumType = new DefaultEnumTypeDefinition(typeName, type, enumValues,
                parserArguments, null);
            types.put(typeName, enumType);
            enumContexts.pop();
        } else {
            // If the type has sub-types it's an abstract type.
            SwitchField switchField = getSwitchField();
            boolean abstractType = switchField != null;
            DefaultComplexTypeDefinition type = new DefaultComplexTypeDefinition(typeName, parserArguments, null,
                abstractType, parserContexts.peek());
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
    public void enterArrayField(MSpecParser.ArrayFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = ctx.name.id.getText();
        ArrayField.LoopType loopType = ArrayField.LoopType.valueOf(ctx.loopType.getText().toUpperCase());
        String loopExpressionString = ctx.loopExpression.expr.getText();
        Term loopExpression =  getExpressionTerm(loopExpressionString);
        Term[] params = getFieldParams((MSpecParser.FieldDefinitionContext) ctx.parent.parent);
        Field field = new DefaultArrayField(null, type, name, loopType, loopExpression, params);
        if(parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterChecksumField(MSpecParser.ChecksumFieldContext ctx) {
        SimpleTypeReference type = getSimpleTypeReference(ctx.type);
        String name = ctx.name.id.getText();
        String checksumExpressionString = ctx.checksumExpression.expr.getText();
        Term checksumExpression = getExpressionTerm(checksumExpressionString);
        Field field = new DefaultChecksumField(null, type, name, checksumExpression);
        if(parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterConstField(MSpecParser.ConstFieldContext ctx) {
        SimpleTypeReference type = getSimpleTypeReference(ctx.type);
        String name = ctx.name.id.getText();
        String expected = ctx.expected.expr.getText();
        Field field = new DefaultConstField(null, type, name, expected);
        if(parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterDiscriminatorField(MSpecParser.DiscriminatorFieldContext ctx) {
        SimpleTypeReference type = getSimpleTypeReference(ctx.type);
        String name = ctx.name.id.getText();
        Field field = new DefaultDiscriminatorField(null, type, name);
        if(parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterEnumField(MSpecParser.EnumFieldContext ctx) {
        ComplexTypeReference type = new DefaultComplexTypeReference(ctx.type.complexTypeReference.getText());
        String name = ctx.name.id.getText();
        Term[] params = getFieldParams((MSpecParser.FieldDefinitionContext) ctx.parent.parent);
        Field field = new DefaultEnumField(null, type, name, params);
        if(parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterImplicitField(MSpecParser.ImplicitFieldContext ctx) {
        SimpleTypeReference type = getSimpleTypeReference(ctx.type);
        String name = ctx.name.id.getText();
        String serializationExpressionString = ctx.serializationExpression.expr.getText();
        Term serializationExpression = getExpressionTerm(serializationExpressionString);
        Field field = new DefaultImplicitField(null, type, name, serializationExpression);
        if(parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterManualArrayField(MSpecParser.ManualArrayFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = ctx.name.id.getText();
        ManualArrayField.LoopType loopType = ManualArrayField.LoopType.valueOf(
            ctx.loopType.getText().toUpperCase());
        String loopExpressionString = ctx.loopExpression.expr.getText();
        Term loopExpression =  getExpressionTerm(loopExpressionString);
        String serializationExpressionString = ctx.serializationExpression.expr.getText();
        Term serializationExpression = getExpressionTerm(serializationExpressionString);
        String deserializationExpressionString = ctx.deserializationExpression.expr.getText();
        Term deserializationExpression = getExpressionTerm(deserializationExpressionString);
        String lengthExpressionString = ctx.lengthExpression.expr.getText();
        Term lengthExpression =  getExpressionTerm(lengthExpressionString);
        Term[] params = getFieldParams((MSpecParser.FieldDefinitionContext) ctx.parent.parent);
        Field field = new DefaultManualArrayField(null, type, name, loopType, loopExpression, serializationExpression,
            deserializationExpression, lengthExpression, params);
        if(parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterManualField(MSpecParser.ManualFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = ctx.name.id.getText();
        String serializationExpressionString = ctx.serializationExpression.expr.getText();
        Term serializationExpression = getExpressionTerm(serializationExpressionString);
        String deserializationExpressionString = ctx.deserializationExpression.expr.getText();
        Term deserializationExpression = getExpressionTerm(deserializationExpressionString);
        String lengthExpressionString = ctx.lengthExpression.expr.getText();
        Term lengthExpression =  getExpressionTerm(lengthExpressionString);
        Term[] params = getFieldParams((MSpecParser.FieldDefinitionContext) ctx.parent.parent);
        Field field = new DefaultManualField(null, type, name, serializationExpression, deserializationExpression,
            lengthExpression, params);
        if(parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterOptionalField(MSpecParser.OptionalFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = ctx.name.id.getText();
        String conditionExpressionString = ctx.condition.expr.getText();
        Term conditionExpression = getExpressionTerm(conditionExpressionString);
        Term[] params = getFieldParams((MSpecParser.FieldDefinitionContext) ctx.parent.parent);
        Field field = new DefaultOptionalField(null, type, name, conditionExpression, params);
        if(parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterPaddingField(MSpecParser.PaddingFieldContext ctx) {
        SimpleTypeReference type = getSimpleTypeReference(ctx.type);
        String name = ctx.name.id.getText();
        String paddingValueString = ctx.paddingValue.expr.getText();
        Term paddingValue = getExpressionTerm(paddingValueString);
        String paddingConditionString = ctx.paddingCondition.expr.getText();
        Term paddingCondition = getExpressionTerm(paddingConditionString);
        Term[] params = getFieldParams((MSpecParser.FieldDefinitionContext) ctx.parent.parent);
        Field field = new DefaultPaddingField(null, type, name, paddingValue, paddingCondition, params);
        if(parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterReservedField(MSpecParser.ReservedFieldContext ctx) {
        SimpleTypeReference type = getSimpleTypeReference(ctx.type);
        String expected = ctx.expected.expr.getText();
        Field field = new DefaultReservedField(null, type, expected);
        if(parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterSimpleField(MSpecParser.SimpleFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = ctx.name.id.getText();
        Term[] params = getFieldParams((MSpecParser.FieldDefinitionContext) ctx.parent.parent);
        Field field = new DefaultSimpleField(null, type, name, params);
        if(parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterTypeSwitchField(MSpecParser.TypeSwitchFieldContext ctx) {
        int numDiscriminators = ctx.discriminators.expression().size();
        String[] discriminatorNames = new String[numDiscriminators];
        for(int i = 0; i < numDiscriminators; i++) {
            discriminatorNames[i] = ctx.discriminators.expression().get(i).expr.getText();
        }
        DefaultSwitchField field = new DefaultSwitchField(discriminatorNames);
        if(parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterVirtualField(MSpecParser.VirtualFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = ctx.name.id.getText();
        String valueExpressionString = ctx.valueExpression.expr.getText();
        Term valueExpression = getExpressionTerm(valueExpressionString);
        Field field = new DefaultVirtualField(null, type, name, valueExpression);
        if(parserContexts.peek() != null) {
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
        // Add all the arguments from the parent type.
        if(((MSpecParser.ComplexTypeContext) ctx.parent.parent.parent.parent).params != null) {
            parserArguments.addAll(Arrays.asList(getParserArguments(
                ((MSpecParser.ComplexTypeContext) ctx.parent.parent.parent.parent).params.argument())));
        }
        // Add all eventually existing local arguments.
        if(ctx.argumentList() != null) {
            parserArguments.addAll(Arrays.asList(getParserArguments(ctx.argumentList().argument())));
        }

        String[] discriminatorValues;
        if(ctx.discriminatorValues != null) {
            List<MSpecParser.ExpressionContext> expressions = ctx.discriminatorValues.expression();
            discriminatorValues = new String[expressions.size()];
            for (int i = 0; i < expressions.size(); i++) {
                discriminatorValues[i] = expressions.get(i).expr.getText();
            }
        } else {
            discriminatorValues = new String[0];
        }
        DefaultDiscriminatedComplexTypeDefinition type =
            new DefaultDiscriminatedComplexTypeDefinition(typeName, parserArguments.toArray(new Argument[0]), null,
                discriminatorValues, parserContexts.pop());

        // Add the type to the switch field definition.
        DefaultSwitchField switchField = getSwitchField();
        if(switchField == null) {
            throw new RuntimeException("This shouldn't have happened");
        }
        switchField.addCase(type);

        // Add the type to the type list.
        types.put(typeName, type);
    }

    @Override
    public void enterEnumValueDefinition(MSpecParser.EnumValueDefinitionContext ctx) {
        String value = unquoteString(ctx.valueExpression.getText());
        String name = ctx.name.getText();
        Map<String, String> constants = null;
        if(ctx.constantValueExpressions != null) {
            MSpecParser.ComplexTypeContext parentCtx = (MSpecParser.ComplexTypeContext) ctx.parent;
            int numConstantValues = parentCtx.params.argument().size();
            int numExpressionValues = ctx.constantValueExpressions.expression().size();
            // This only works if we provide exactly the same number of expressions as we defined constants
            if(numConstantValues != numExpressionValues) {
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
                constants.put(constantName, constant);
            }
        }
        final DefaultEnumValue enumValue = new DefaultEnumValue(value, name, constants);
        this.enumContexts.peek().add(enumValue);
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
        if(ctx.simpleTypeReference != null) {
            SimpleTypeReference.SimpleBaseType simpleBaseType = SimpleTypeReference.SimpleBaseType.valueOf(
                ctx.simpleTypeReference.base.getText().toUpperCase());
            if(ctx.simpleTypeReference.size != null) {
                int size = Integer.parseInt(ctx.simpleTypeReference.size.getText());
                return new DefaultSimpleTypeReference(simpleBaseType, size);
            } else {
                return new DefaultSimpleTypeReference(simpleBaseType, 1);
            }
        } else {
            return new DefaultComplexTypeReference(ctx.complexTypeReference.getText());
        }
    }

    private SimpleTypeReference getSimpleTypeReference(MSpecParser.DataTypeContext ctx) {
        SimpleTypeReference.SimpleBaseType simpleBaseType =
            SimpleTypeReference.SimpleBaseType.valueOf(ctx.base.getText().toUpperCase());
        if(ctx.size != null) {
            int size = Integer.parseInt(ctx.size.getText());
            return new DefaultSimpleTypeReference(simpleBaseType, size);
        } else {
            return new DefaultSimpleTypeReference(simpleBaseType, 1);
        }
    }

    private DefaultSwitchField getSwitchField() {
        for (Field field : Objects.requireNonNull(parserContexts.peek())) {
            if(field instanceof DefaultSwitchField) {
                return (DefaultSwitchField) field;
            }
        }
        return null;
    }

    private EnumValue[] getEnumValues() {
        return Objects.requireNonNull(enumContexts.peek()).toArray(new EnumValue[0]);
    }

    private Argument[] getParserArguments(List<MSpecParser.ArgumentContext> params) {
        Argument[] parserArguments = new Argument[params.size()];
        for (int i = 0; i < params.size(); i++) {
            TypeReference type = getTypeReference(params.get(i).type);
            String name = params.get(i).name.id.getText();
            parserArguments[i] = new Argument(type, name);
        }
        return parserArguments;
    }

    private Term[] getFieldParams(MSpecParser.FieldDefinitionContext parentCtx) {
        Term[] params = null;
        if(parentCtx.params != null) {
            params = new Term[parentCtx.params.expression().size()];
            for(int i = 0; i < parentCtx.params.expression().size(); i++) {
                params[i] = parseExpression(parentCtx.params.expression().get(i).expr.getText());
            }
        }
        return params;
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
        if(quotedString != null && quotedString.length() > 2) {
            return quotedString.substring(1, quotedString.length() - 1);
        }
        return quotedString;
    }

    public static void main(String[] args) {
        //Term term = new ExpressionStringParser().parse(IOUtils.toInputStream("CAST(CAST(parameter,S7ParameterUserData).items(hurz)[0],S7ParameterUserDataItemCPUFunctions).cpuFunctionType", Charset.defaultCharset()));
        //Term term = new ExpressionStringParser().parse(IOUtils.toInputStream("CAST(parameter,S7ParameterUserData).items(hurz)[0]", Charset.defaultCharset()));
        Term term = new ExpressionStringParser().parse(IOUtils.toInputStream("payload.lengthInBytes+4", Charset.defaultCharset()));

        System.out.println(term);
    }

}
