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

package org.apache.plc4x.plugins.codegenerator.parser;

import org.apache.plc4x.codegenerator.parser.imaginary.ImaginaryBaseListener;
import org.apache.plc4x.codegenerator.parser.imaginary.ImaginaryParser;
import org.apache.plc4x.language.fields.ArrayField;
import org.apache.plc4x.language.fields.Field;
import org.apache.plc4x.language.fields.SwitchField;
import org.apache.plc4x.language.references.SimpleTypeReference;
import org.apache.plc4x.language.references.TypeReference;
import org.apache.plc4x.language.definitions.ComplexTypeDefinition;
import org.apache.plc4x.language.definitions.DiscriminatedComplexTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.model.definitions.DefaultComplexTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.model.definitions.DefaultDiscriminatedComplexTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.model.references.DefaultComplexTypeReference;
import org.apache.plc4x.plugins.codegenerator.model.references.DefaultSimpleTypeReference;
import org.apache.plc4x.plugins.codegenerator.model.fields.*;

import java.util.*;

public class MessageFormatListener extends ImaginaryBaseListener {

    private Stack<List<Field>> parserContexts;
    private Map<String, ComplexTypeDefinition> complexTypes;

    public Map<String, ComplexTypeDefinition> getComplexTypes() {
        return complexTypes;
    }

    @Override
    public void enterFile(ImaginaryParser.FileContext ctx) {
        parserContexts = new Stack<>();
        complexTypes = new HashMap<>();
    }

    @Override
    public void enterComplexType(ImaginaryParser.ComplexTypeContext ctx) {
        List<Field> parserContext = new LinkedList<>();
        parserContexts.push(parserContext);
    }

    @Override
    public void exitComplexType(ImaginaryParser.ComplexTypeContext ctx) {
        String typeName = ctx.name.id.getText();

        // If the type has sub-types it's an abstract type.
        SwitchField switchField = getSwitchField();
        boolean abstractType = switchField != null;
        DefaultComplexTypeDefinition type = new DefaultComplexTypeDefinition(typeName, abstractType, parserContexts.peek());
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
    public void enterArrayField(ImaginaryParser.ArrayFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = ctx.name.id.getText();
        ArrayField.LengthType lengthType;
        if(ctx.lengthType.K_COUNT() != null) {
            lengthType = ArrayField.LengthType.COUNT;
        } else {
            lengthType = ArrayField.LengthType.LENGTH;
        }
        String lengthExpression = ctx.lengthExpression.getText();
        Field field = new DefaultArrayField(type, name, lengthType, lengthExpression);
        parserContexts.peek().add(field);
    }

    @Override
    public void enterConstField(ImaginaryParser.ConstFieldContext ctx) {
        SimpleTypeReference type = getSimpleTypeReference(ctx.type);
        String name = ctx.name.id.getText();
        String expected = ctx.expected.expr.getText();
        Field field = new DefaultConstField(type, name, expected);
        parserContexts.peek().add(field);
    }

    @Override
    public void enterContextField(ImaginaryParser.ContextFieldContext ctx) {
        // This field type just adds something to the context (No parsing involved)
        // TODO: Implement
    }

    @Override
    public void enterDiscriminatorField(ImaginaryParser.DiscriminatorFieldContext ctx) {
        SimpleTypeReference type = getSimpleTypeReference(ctx.type);
        String name = ctx.name.id.getText();
        Field field = new DefaultDiscriminatorField(type, name);
        parserContexts.peek().add(field);
    }

    @Override
    public void enterEmbeddedField(ImaginaryParser.EmbeddedFieldContext ctx) {
        super.enterEmbeddedField(ctx);
        // TODO: Implement
    }

    @Override
    public void enterSimpleField(ImaginaryParser.SimpleFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = ctx.name.id.getText();
        Field field = new DefaultSimpleField(type, name);
        parserContexts.peek().add(field);
    }

    @Override
    public void enterImplicitField(ImaginaryParser.ImplicitFieldContext ctx) {
        SimpleTypeReference type = getSimpleTypeReference(ctx.type);
        String name = ctx.name.id.getText();
        String serializationExpression = ctx.serializationExpression.getText();
        Field field = new DefaultImplicitField(type, name, serializationExpression);
        parserContexts.peek().add(field);
    }

    @Override
    public void enterOptionalField(ImaginaryParser.OptionalFieldContext ctx) {
        TypeReference type = getTypeReference(ctx.type);
        String name = ctx.name.id.getText();
        String conditionExpression = ctx.condition.expr.getText();
        Field field = new DefaultOptionalField(type, name, conditionExpression);
        parserContexts.peek().add(field);
    }

    @Override
    public void enterReservedField(ImaginaryParser.ReservedFieldContext ctx) {
        SimpleTypeReference type = getSimpleTypeReference(ctx.type);
        String expected = ctx.expected.expr.getText();
        Field field = new DefaultReservedField(type, expected);
        parserContexts.peek().add(field);
    }

    @Override
    public void enterTypeSwitchField(ImaginaryParser.TypeSwitchFieldContext ctx) {
        int numDiscriminators = ctx.discriminators.expression().size();
        String[] discriminatorNames = new String[numDiscriminators];
        for(int i = 0; i < numDiscriminators; i++) {
            discriminatorNames[i] = ctx.discriminators.expression().get(i).expr.getText();
        }
        DefaultSwitchField switchField = new DefaultSwitchField(discriminatorNames);
        parserContexts.peek().add(switchField);
    }

    @Override
    public void enterCaseStatement(ImaginaryParser.CaseStatementContext ctx) {
        List<Field> parserContext = new LinkedList<>();
        parserContexts.push(parserContext);
    }

    @Override
    public void exitCaseStatement(ImaginaryParser.CaseStatementContext ctx) {
        String typeName = ctx.name.getText();
        List<ImaginaryParser.ExpressionContext> expressions = ctx.discriminatorValues.expression();
        String[] discriminatorValues = new String[expressions.size()];
        for (int i = 0; i < expressions.size(); i++) {
            discriminatorValues[i] = expressions.get(i).expr.getText();
        }
        DefaultDiscriminatedComplexTypeDefinition type =
            new DefaultDiscriminatedComplexTypeDefinition(typeName, discriminatorValues, parserContexts.pop());

        // Add the type to the switch field definition.
        DefaultSwitchField switchField = getSwitchField();
        if(switchField == null) {
            throw new RuntimeException("This shouldn't have happened");
        }
        switchField.addCase(type);

        // Add the type to the type list.
        complexTypes.put(typeName, type);
    }

    private TypeReference getTypeReference(ImaginaryParser.TypeReferenceContext ctx) {
        if(ctx.simpleTypeReference != null) {
            SimpleTypeReference.SimpleBaseType simpleBaseType = SimpleTypeReference.SimpleBaseType.valueOf(
                ctx.simpleTypeReference.base.getText().toUpperCase());
            if(ctx.simpleTypeReference.size != null) {
                int size = Integer.valueOf(ctx.simpleTypeReference.size.getText());
                return new DefaultSimpleTypeReference(simpleBaseType, size);
            } else {
                return new DefaultSimpleTypeReference(simpleBaseType, 1);
            }
        } else {
            return new DefaultComplexTypeReference(ctx.complexTypeReference.getText());
        }
    }

    private SimpleTypeReference getSimpleTypeReference(ImaginaryParser.DataTypeContext ctx) {
        SimpleTypeReference.SimpleBaseType simpleBaseType =
            SimpleTypeReference.SimpleBaseType.valueOf(ctx.base.getText().toUpperCase());
        if(ctx.size != null) {
            int size = Integer.valueOf(ctx.size.getText());
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
}
