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
import org.apache.plc4x.plugins.codegenerator.model.ComplexType;
import org.apache.plc4x.plugins.codegenerator.model.DiscriminatedComplexType;
import org.apache.plc4x.plugins.codegenerator.model.SimpleType;
import org.apache.plc4x.plugins.codegenerator.model.fields.*;

import java.util.*;

public class MessageFormatListener extends ImaginaryBaseListener {

    private Stack<List<Field>> parserContexts;
    private Map<String, ComplexType> complexTypes;

    public Map<String, ComplexType> getComplexTypes() {
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
        ComplexType type = new ComplexType(typeName, abstractType, parserContexts.peek());
        complexTypes.put(typeName, type);

        // Set the parent type for all sub-types.
        if(switchField != null) {
            for (DiscriminatedComplexType subtype : switchField.getCases()) {
                subtype.setParentType(type);
            }
        }

        parserContexts.pop();
    }

    @Override
    public void enterArrayField(ImaginaryParser.ArrayFieldContext ctx) {
        String typeName = ctx.type.getText();
        String name = ctx.name.IDENTIFIER().getText();
        ArrayField.LengthType lengthType;
        if(ctx.lengthType.K_COUNT() != null) {
            lengthType = ArrayField.LengthType.COUNT;
        } else {
            lengthType = ArrayField.LengthType.LENGTH;
        }
        String lengthExpression = ctx.lengthExpression.getText();
        Field field = new ArrayField(typeName, name, lengthType, lengthExpression);
        parserContexts.peek().add(field);
    }

    @Override
    public void enterConstField(ImaginaryParser.ConstFieldContext ctx) {
        SimpleType type = new SimpleType(ctx.type.getText());
        // TODO use a visitor or something to get this without TICKs
        String expected = ctx.expected.expr.getText();
        Field field = new ConstField(type, expected);
        parserContexts.peek().add(field);
    }

    @Override
    public void enterContextField(ImaginaryParser.ContextFieldContext ctx) {
        // This field type just adds something to the context (No parsing involved)
        // TODO: Implement
    }

    @Override
    public void enterDiscriminatorField(ImaginaryParser.DiscriminatorFieldContext ctx) {
        SimpleType type = new SimpleType(ctx.type.getText());
        String name = ctx.name.getText();
        Field field = new DiscriminatorField(type, name);
        parserContexts.peek().add(field);
    }

    @Override
    public void enterEmbeddedField(ImaginaryParser.EmbeddedFieldContext ctx) {
        super.enterEmbeddedField(ctx);
        // TODO: Implement
    }

    @Override
    public void enterSimpleField(ImaginaryParser.SimpleFieldContext ctx) {
        SimpleType type = new SimpleType(ctx.type.getText());
        String name = ctx.name.IDENTIFIER().getText();
        Field field = new SimpleField(type, name);
        parserContexts.peek().add(field);
    }

    @Override
    public void enterImplicitField(ImaginaryParser.ImplicitFieldContext ctx) {
        SimpleType type = new SimpleType(ctx.type.getText());
        String serializationExpression = ctx.serializationExpression.innerExpression().getText();
        Field field = new ImplicitField(type, serializationExpression);
        parserContexts.peek().add(field);
    }

    @Override
    public void enterOptionalField(ImaginaryParser.OptionalFieldContext ctx) {
        SimpleType type = new SimpleType(ctx.type.getText());
        String name = ctx.name.IDENTIFIER().getText();
        String conditionExpression = ctx.condition.expr.getText();
        Field field = new OptionalField(type, name, conditionExpression);
        parserContexts.peek().add(field);
    }

    @Override
    public void enterReservedField(ImaginaryParser.ReservedFieldContext ctx) {
        SimpleType type = new SimpleType(ctx.type.getText());
        String expected = ctx.expected.getText();
        Field field = new ReservedField(type, expected);
        parserContexts.peek().add(field);
    }

    @Override
    public void enterTypeSwitchField(ImaginaryParser.TypeSwitchFieldContext ctx) {
        int numDescriminators = ctx.discriminators.expression().size();
        String[] discriminators = new String[numDescriminators];
        for(int i = 0; i < numDescriminators; i++) {
            discriminators[i] = ctx.discriminators.expression().get(i).getText();
        }
        SwitchField switchField = new SwitchField(discriminators);
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
        int numDiscriminatorValues = (ctx.discriminatorValues.children != null) ?
            ctx.discriminatorValues.children.size() : 0;
        String[] discriminatorValues = new String[numDiscriminatorValues];
        for (int i = 0; i < numDiscriminatorValues; i++) {
            discriminatorValues[i] = ctx.discriminatorValues.children.get(i).getText();
        }
        DiscriminatedComplexType type =
            new DiscriminatedComplexType(typeName, discriminatorValues, parserContexts.pop());

        // Add the type to the switch field definition.
        SwitchField switchField = getSwitchField();
        if(switchField == null) {
            throw new RuntimeException("This shouldn't have happened");
        }
        switchField.addCase(type);

        // Add the type to the type list.
        complexTypes.put(typeName, type);
    }


    private SwitchField getSwitchField() {
        for (Field field : parserContexts.peek()) {
            if(field instanceof SwitchField) {
                return (SwitchField) field;
            }
        }
        return null;
    }
}
