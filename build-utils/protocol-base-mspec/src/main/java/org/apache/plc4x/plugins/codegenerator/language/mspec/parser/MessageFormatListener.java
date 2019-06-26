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
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.fields.*;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.references.DefaultComplexTypeReference;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.references.DefaultSimpleTypeReference;
import org.apache.plc4x.plugins.codegenerator.types.definitions.Argument;
import org.apache.plc4x.plugins.codegenerator.types.definitions.ComplexTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.definitions.DiscriminatedComplexTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.fields.ArrayField;
import org.apache.plc4x.plugins.codegenerator.types.fields.Field;
import org.apache.plc4x.plugins.codegenerator.types.fields.SwitchField;
import org.apache.plc4x.plugins.codegenerator.types.references.SimpleTypeReference;
import org.apache.plc4x.plugins.codegenerator.types.references.TypeReference;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;

import java.io.InputStream;
import java.util.*;

public class MessageFormatListener extends MSpecBaseListener {

    private Stack<List<Field>> parserContexts;
    private Map<String, ComplexTypeDefinition> complexTypes;

    public Map<String, ComplexTypeDefinition> getComplexTypes() {
        return complexTypes;
    }

    @Override
    public void enterFile(MSpecParser.FileContext ctx) {
        parserContexts = new Stack<>();
        complexTypes = new HashMap<>();
    }

    @Override
    public void enterComplexType(MSpecParser.ComplexTypeContext ctx) {
        List<Field> parserContext = new LinkedList<>();
        parserContexts.push(parserContext);
    }

    @Override
    public void exitComplexType(MSpecParser.ComplexTypeContext ctx) {
        String typeName = ctx.name.id.getText();
        Argument[] parserArguments = null;
        if(ctx.params != null) {
            parserArguments = getParserArguments(ctx.params.argument());
        }

        // If the type has sub-types it's an abstract type.
        SwitchField switchField = getSwitchField();
        boolean abstractType = switchField != null;
        DefaultComplexTypeDefinition type = new DefaultComplexTypeDefinition(typeName, parserArguments,
            abstractType, parserContexts.peek());
        complexTypes.put(typeName, type);

        // Set the parent type for all sub-types.
        if(switchField != null) {
            for (DiscriminatedComplexTypeDefinition subtype : switchField.getCases()) {
                if(subtype instanceof DefaultDiscriminatedComplexTypeDefinition) {
                    ((DefaultDiscriminatedComplexTypeDefinition) subtype).setParentType(type);
                }
            }
        }

        parserContexts.pop();
    }

    @Override
    public void enterArrayField(MSpecParser.ArrayFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = ctx.name.id.getText();
        ArrayField.LengthType lengthType;
        if(ctx.lengthType.K_COUNT() != null) {
            lengthType = ArrayField.LengthType.COUNT;
        } else {
            lengthType = ArrayField.LengthType.LENGTH;
        }
        String lengthExpressionString = ctx.lengthExpression.expr.getText();
        InputStream inputStream = IOUtils.toInputStream(lengthExpressionString);
        ExpressionStringParser parser = new ExpressionStringParser();
        Term lengthExpression =  parser.parse(inputStream);
        String[] params = getFieldParams((MSpecParser.FieldDefinitionContext) ctx.parent.parent);
        Field field = new DefaultArrayField(type, name, lengthType, lengthExpression, params);
        parserContexts.peek().add(field);
    }

    @Override
    public void enterConstField(MSpecParser.ConstFieldContext ctx) {
        SimpleTypeReference type = getSimpleTypeReference(ctx.type);
        String name = ctx.name.id.getText();
        String expected = ctx.expected.expr.getText();
        Field field = new DefaultConstField(type, name, expected);
        parserContexts.peek().add(field);
    }

    @Override
    public void enterDiscriminatorField(MSpecParser.DiscriminatorFieldContext ctx) {
        SimpleTypeReference type = getSimpleTypeReference(ctx.type);
        String name = ctx.name.id.getText();
        Field field = new DefaultDiscriminatorField(type, name);
        parserContexts.peek().add(field);
    }

    @Override
    public void enterSimpleField(MSpecParser.SimpleFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = ctx.name.id.getText();
        String[] params = getFieldParams((MSpecParser.FieldDefinitionContext) ctx.parent.parent);
        Field field = new DefaultSimpleField(type, name, params);
        parserContexts.peek().add(field);
    }

    @Override
    public void enterImplicitField(MSpecParser.ImplicitFieldContext ctx) {
        SimpleTypeReference type = getSimpleTypeReference(ctx.type);
        String name = ctx.name.id.getText();
        String serializationExpressionString = ctx.serializationExpression.expr.getText();
        InputStream inputStream = IOUtils.toInputStream(serializationExpressionString);
        ExpressionStringParser parser = new ExpressionStringParser();
        Term serializationExpression =  parser.parse(inputStream);
        Field field = new DefaultImplicitField(type, name, serializationExpression);
        parserContexts.peek().add(field);
    }

    @Override
    public void enterOptionalField(MSpecParser.OptionalFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = ctx.name.id.getText();
        String conditionExpressionString = ctx.condition.expr.getText();
        InputStream inputStream = IOUtils.toInputStream(conditionExpressionString);
        ExpressionStringParser parser = new ExpressionStringParser();
        Term conditionExpression =  parser.parse(inputStream);
        String[] params = getFieldParams((MSpecParser.FieldDefinitionContext) ctx.parent.parent);
        Field field = new DefaultOptionalField(type, name, conditionExpression, params);
        parserContexts.peek().add(field);
    }

    @Override
    public void enterReservedField(MSpecParser.ReservedFieldContext ctx) {
        SimpleTypeReference type = getSimpleTypeReference(ctx.type);
        String expected = ctx.expected.expr.getText();
        Field field = new DefaultReservedField(type, expected);
        parserContexts.peek().add(field);
    }

    @Override
    public void enterTypeSwitchField(MSpecParser.TypeSwitchFieldContext ctx) {
        int numDiscriminators = ctx.discriminators.expression().size();
        String[] discriminatorNames = new String[numDiscriminators];
        for(int i = 0; i < numDiscriminators; i++) {
            discriminatorNames[i] = ctx.discriminators.expression().get(i).expr.getText();
        }
        DefaultSwitchField switchField = new DefaultSwitchField(discriminatorNames);
        parserContexts.peek().add(switchField);
    }

    @Override
    public void enterCaseStatement(MSpecParser.CaseStatementContext ctx) {
        List<Field> parserContext = new LinkedList<>();
        parserContexts.push(parserContext);
    }

    @Override
    public void exitCaseStatement(MSpecParser.CaseStatementContext ctx) {
        String typeName = ctx.name.getText();
        Argument[] parserArguments = null;
        if(((MSpecParser.ComplexTypeContext) ctx.parent.parent.parent.parent).params != null) {
            parserArguments = getParserArguments(
                ((MSpecParser.ComplexTypeContext) ctx.parent.parent.parent.parent).params.argument());
        }

        List<MSpecParser.ExpressionContext> expressions = ctx.discriminatorValues.expression();
        String[] discriminatorValues = new String[expressions.size()];
        for (int i = 0; i < expressions.size(); i++) {
            discriminatorValues[i] = expressions.get(i).expr.getText();
        }
        DefaultDiscriminatedComplexTypeDefinition type =
            new DefaultDiscriminatedComplexTypeDefinition(typeName, parserArguments,
                discriminatorValues, parserContexts.pop());

        // Add the type to the switch field definition.
        DefaultSwitchField switchField = getSwitchField();
        if(switchField == null) {
            throw new RuntimeException("This shouldn't have happened");
        }
        switchField.addCase(type);

        // Add the type to the type list.
        complexTypes.put(typeName, type);
    }

    /*@Override
    public void exitExpression(MSpecParser.ExpressionContext ctx) {
        MSpecParser.InnerExpressionContext expressionContext = ctx.innerExpression();
        super.exitExpression(ctx);
    }*/

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
        for (Field field : parserContexts.peek()) {
            if(field instanceof DefaultSwitchField) {
                return (DefaultSwitchField) field;
            }
        }
        return null;
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

    private String[] getFieldParams(MSpecParser.FieldDefinitionContext parentCtx) {
        String[] params = null;
        if(parentCtx.params != null) {
            params = new String[parentCtx.params.expression().size()];
            for(int i = 0; i < parentCtx.params.expression().size(); i++) {
                params[i] = parentCtx.params.expression().get(i).expr.getText();
            }
        }
        return params;
    }

    private Term parseExpression(String expressionString) {
        return null;
    }

}
