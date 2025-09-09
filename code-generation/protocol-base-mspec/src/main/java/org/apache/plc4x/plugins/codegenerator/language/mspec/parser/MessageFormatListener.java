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
package org.apache.plc4x.plugins.codegenerator.language.mspec.parser;

import org.antlr.v4.runtime.CommonTokenStream;
import org.antlr.v4.runtime.ParserRuleContext;
import org.antlr.v4.runtime.RuleContext;
import org.antlr.v4.runtime.Token;
import org.apache.commons.io.IOUtils;
import org.apache.commons.lang3.StringUtils;
import org.apache.plc4x.plugins.codegenerator.language.mspec.LazyTypeDefinitionConsumer;
import org.apache.plc4x.plugins.codegenerator.language.mspec.MSpecBaseListener;
import org.apache.plc4x.plugins.codegenerator.language.mspec.MSpecLexer;
import org.apache.plc4x.plugins.codegenerator.language.mspec.MSpecParser;
import org.apache.plc4x.plugins.codegenerator.language.mspec.expression.ExpressionStringParser;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions.*;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.fields.*;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.references.*;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.terms.DefaultVariableLiteral;
import org.apache.plc4x.plugins.codegenerator.protocol.TypeContext;
import org.apache.plc4x.plugins.codegenerator.types.definitions.Argument;
import org.apache.plc4x.plugins.codegenerator.types.definitions.DiscriminatedComplexTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.enums.EnumValue;
import org.apache.plc4x.plugins.codegenerator.types.fields.ArrayField;
import org.apache.plc4x.plugins.codegenerator.types.fields.Field;
import org.apache.plc4x.plugins.codegenerator.types.fields.ManualArrayField;
import org.apache.plc4x.plugins.codegenerator.types.fields.SwitchField;
import org.apache.plc4x.plugins.codegenerator.types.references.SimpleTypeReference;
import org.apache.plc4x.plugins.codegenerator.types.references.TypeReference;
import org.apache.plc4x.plugins.codegenerator.types.terms.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.nio.charset.Charset;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.CompletionStage;
import java.util.function.Consumer;
import java.util.stream.Collectors;

public class MessageFormatListener extends MSpecBaseListener implements LazyTypeDefinitionConsumer {

    private static final Logger LOGGER = LoggerFactory.getLogger(MessageFormatListener.class);
    protected final Map<String, TypeDefinition> types;
    protected final Map<String, List<Consumer<TypeDefinition>>> typeDefinitionConsumers;
    private final CommonTokenStream tokenStream;
    private final Stack<Map<String, Term>> batchSetAttributes = new Stack<>();
    private final Stack<String> currentTypeName = new Stack<>();
    // For storing comments for complex types between enter/exit
    private final Map<MSpecParser.ComplexTypeContext, String> complexTypeComments = new HashMap<>();
    private Deque<List<Field>> parserContexts;
    private Deque<List<EnumValue>> enumContexts;

    public MessageFormatListener(CommonTokenStream tokenStream) {
        this.tokenStream = tokenStream;
        types = new HashMap<>();
        typeDefinitionConsumers = new HashMap<>();
    }

    public MessageFormatListener(CommonTokenStream tokenStream, TypeContext exitingTypeContext) {
        this.tokenStream = tokenStream;
        types = new HashMap<>(exitingTypeContext.getTypeDefinitions());
        typeDefinitionConsumers = new HashMap<>(exitingTypeContext.getUnresolvedTypeReferences());
    }

    public Deque<List<Field>> getParserContexts() {
        return parserContexts;
    }

    public Deque<List<EnumValue>> getEnumContexts() {
        return enumContexts;
    }

    @Override
    public void enterFile(MSpecParser.FileContext ctx) {
        parserContexts = new LinkedList<>();
        enumContexts = new LinkedList<>();
    }

    // Helper: Attach comment if directly above or on the same line
    private String consumePendingComment(ParserRuleContext ctx) {
        Token startToken = ctx.getStart();
        while (startToken.getType() != MSpecLexer.LBRACKET) { // Bit hacky workaround to assume start is always the bracket but should work
            ctx = ctx.getParent();          // Switch to parent context...
            startToken = ctx.getStart();    // ...and use this as start
        }

        List<Token> frontComments = new LinkedList<>();
        int tokenIndex = startToken.getTokenIndex();
        List<Token> hiddenToLeft = tokenStream.getHiddenTokensToLeft(tokenIndex);
        if (hiddenToLeft != null) {
            hiddenToLeft.reversed().stream()
                .takeWhile(token -> frontComments.isEmpty() || (token.getType() != MSpecLexer.EmptyLine))
                .filter(token -> token.getType() == MSpecLexer.LINE_COMMENT || token.getType() == MSpecLexer.BLOCK_COMMENT)
                .forEach(frontComments::add);
        }
        Token eolComment = null;
        List<Token> hiddenToRight = tokenStream.getHiddenTokensToRight(tokenIndex);
        if (hiddenToRight != null) {
            eolComment = hiddenToRight.stream()
                .filter(token -> token.getType() == MSpecLexer.LINE_COMMENT || token.getType() == MSpecLexer.BLOCK_COMMENT)
                .findFirst()
                .orElse(null);
        }

        StringBuilder comment = new StringBuilder();

        var defLine = startToken.getLine();
        if (!frontComments.isEmpty() && frontComments.getFirst().getLine() == defLine - 1) {
            frontComments.reversed().stream().map(Token::getText).forEach(tokenText -> comment.append(tokenText).append('\n'));
        }

        if (eolComment != null && eolComment.getLine() == defLine) {
            comment.append(eolComment.getText());
        }

        if (comment.isEmpty()) {
            return null;
        }
        String commentString = comment.toString();
        if (StringUtils.isBlank(commentString)) {
            return null;
        }
        commentString = StringUtils.stripEnd(commentString, "\n");
        return commentString;
    }

    @Override
    public void enterContantsDefinition(MSpecParser.ContantsDefinitionContext ctx) {
        currentTypeName.push("Constants");
        // Set a map of attributes that should be set for all fields.
        Map<String, Term> curBatchSetAttributes = new HashMap<>();
        // Make the new Map the top of the stack.
        batchSetAttributes.push(curBatchSetAttributes);

        List<Field> parserContext = new LinkedList<>();
        parserContexts.push(parserContext);
        super.enterContantsDefinition(ctx);
    }

    @Override
    public void exitContantsDefinition(MSpecParser.ContantsDefinitionContext ctx) {
        String typeName = "Constants";

        // If the type has subtypes, it's an abstract type.
        var fields = parserContexts.pop();
        var type = new DefaultConstantsTypeDefinition(typeName, fields, consumePendingComment(ctx));
        // Link the fields and the complex types.
        if (fields != null) {
            fields.forEach(field -> ((DefaultField) field).setOwner(type));
        }
        dispatchType(typeName, type);

        currentTypeName.pop();
    }

    @Override
    public void enterGlobalsDefinition(MSpecParser.GlobalsDefinitionContext ctx) {
        super.enterGlobalsDefinition(ctx);
    }

    @Override
    public void exitGlobalsDefinition(MSpecParser.GlobalsDefinitionContext ctx) {
        super.exitGlobalsDefinition(ctx);
    }

    @Override
    public void enterContextDefintion(MSpecParser.ContextDefintionContext ctx) {
        super.enterContextDefintion(ctx);
    }

    @Override
    public void exitContextDefintion(MSpecParser.ContextDefintionContext ctx) {
        super.exitContextDefintion(ctx);
    }

    @Override
    public void enterComplexType(MSpecParser.ComplexTypeContext ctx) {
        currentTypeName.push(getIdString(ctx.name));
        // Set a map of attributes that should be set for all fields.
        Map<String, Term> curBatchSetAttributes = new HashMap<>();
        // Add all attributes defined in the current batchSet field.
        for (var attributeContext : ctx.attributes.attribute()) {
            var attributeExpression = getExpressionTerm(attributeContext.value);
            curBatchSetAttributes.put(attributeContext.name.getText(), attributeExpression);
        }
        // Make the new Map the top of the stack.
        batchSetAttributes.push(curBatchSetAttributes);

        // Attach comment if present (store for use in exit)
        var comment = consumePendingComment(ctx);
        if (comment != null) {
            complexTypeComments.put(ctx, comment);
        }

        if ("enum".equals(ctx.getChild(0).getText())) {
            var enumContext = new LinkedList<EnumValue>();
            enumContexts.push(enumContext);
        } else {
            var parserContext = new LinkedList<Field>();
            parserContexts.push(parserContext);
        }
    }

    @Override
    public void exitComplexType(MSpecParser.ComplexTypeContext ctx) {
        var typeName = getIdString(ctx.name);
        var parserArguments = ctx.params != null ? getParserArguments(ctx.params.argument()) : null;

        var attributes = batchSetAttributes.peek();
        // Handle enum types.
        if ("enum".equals(ctx.getChild(0).getText())) {
            var type = ctx.type != null ? getSimpleTypeReference(ctx.type) : null;
            var enumValues = getEnumValues();
            if (type == null) {
                // in case there is no type we default to uint32
                type = new DefaultIntegerTypeReference(SimpleTypeReference.SimpleBaseType.UINT, 32);
            }
            var enumType = new DefaultEnumTypeDefinition(typeName, type, attributes, enumValues, parserArguments, consumePendingComment(ctx));
            dispatchType(typeName, enumType);
            enumContexts.pop();
        } else if (ctx.dataIoTypeSwitch != null) { // Handle data-io types.
            var switchField = getSwitchField();
            var type = new DefaultDataIoTypeDefinition(typeName, attributes, parserArguments, switchField, consumePendingComment(ctx));
            dispatchType(typeName, type);

            // Set the parent type for all sub-types.
            if (switchField != null) {
                for (DiscriminatedComplexTypeDefinition subtype : switchField.getCases()) {
                    if (subtype instanceof DefaultDiscriminatedComplexTypeDefinition ddctd) {
                        LOGGER.debug("Setting parent {} for {}", type, ddctd);
                        ddctd.setParentType(type);
                    }
                }
            }
            parserContexts.pop();
        } else { // Handle all other types.
            // If the type has sub-types it's an abstract type.
            var switchField = getSwitchField();
            var abstractType = switchField != null;
            var fields = parserContexts.pop();
            var type = new DefaultComplexTypeDefinition(typeName, attributes, parserArguments, abstractType, fields, complexTypeComments.remove(ctx));
            // Link the fields and the complex types.
            if (fields != null) {
                fields.forEach(field -> ((DefaultField) field).setOwner(type));
            }
            dispatchType(typeName, type);

            // Set the parent type for all sub-types.
            setParentRelationship(type);
        }
        currentTypeName.pop();
    }

    protected void setParentRelationship(DefaultComplexTypeDefinition type) {
        Optional<SwitchField> switchField = type.getSwitchField();
        if (switchField.isPresent()) {
            for (DiscriminatedComplexTypeDefinition subtype : switchField.get().getCases()) {
                if (subtype instanceof DefaultDiscriminatedComplexTypeDefinition ddct) {
                    LOGGER.debug("Setting parent {} for {}", type, ddct);
                    ddct.setParentType(type);
                    setParentRelationship(ddct);
                }
            }
        }
    }

    @Override
    public void enterBatchSetDefinition(MSpecParser.BatchSetDefinitionContext ctx) {
        // Set a map of attributes that should be set for all fields.
        var curBatchSetAttributes = new HashMap<String, Term>();
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
        var name = getIdString(ctx.name);
        var field = new DefaultAbstractField(getAttributes(ctx), name, consumePendingComment(ctx));
        getTypeReference(ctx.type).whenComplete((typeReference, throwable) -> {
            if (throwable != null) {
                // TODO: proper error collection in type context error bucket
                LOGGER.debug("Error setting type for {}", field, throwable);
                return;
            }
            field.setType(typeReference);
        });
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterArrayField(MSpecParser.ArrayFieldContext ctx) {
        var name = getIdString(ctx.name);
        var loopType = ArrayField.LoopType.valueOf(ctx.loopType.getText().toUpperCase());
        var loopExpression = getExpressionTerm(ctx.loopExpression);
        var field = new DefaultArrayField(getAttributes(ctx), name, loopType, loopExpression, consumePendingComment(ctx));
        getTypeReference(ctx.type).whenComplete((typeReference, throwable) -> {
            if (throwable != null) {
                // TODO: proper error collection in type context error bucket
                LOGGER.debug("Error setting type for {}", field, throwable);
                return;
            }
            field.setType(new DefaultArrayTypeReference(typeReference));
        });
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterAssertField(MSpecParser.AssertFieldContext ctx) {
        var name = getIdString(ctx.name);
        var conditionExpression = getExpressionTerm(ctx.condition);
        var field = new DefaultAssertField(getAttributes(ctx), name, conditionExpression, consumePendingComment(ctx));
        getTypeReference(ctx.type).whenComplete((typeReference, throwable) -> {
            if (throwable != null) {
                // TODO: proper error collection in type context error bucket
                LOGGER.debug("Error setting type for {}", field, throwable);
                return;
            }
            field.setType(typeReference);
        });
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterChecksumField(MSpecParser.ChecksumFieldContext ctx) {
        var type = getSimpleTypeReference(ctx.type);
        var name = getIdString(ctx.name);
        var checksumExpression = getExpressionTerm(ctx.checksumExpression);
        var field = new DefaultChecksumField(getAttributes(ctx), type, name, checksumExpression, consumePendingComment(ctx));
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterConstField(MSpecParser.ConstFieldContext ctx) {
        var name = getIdString(ctx.name);
        var field = new DefaultConstField(getAttributes(ctx), name, getValueLiteral(ctx.expected), consumePendingComment(ctx));
        if (ctx.type.dataType() != null) {
            field.setType(getSimpleTypeReference(ctx.type.dataType()));
        } else {
            getTypeReference(ctx.type).whenComplete((typeReference, throwable) -> {
                if (throwable != null) {
                    // TODO: proper error collection in type context error bucket
                    LOGGER.debug("Error setting type for {}", field, throwable);
                    return;
                }
                field.setType(typeReference);
            });
        }
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterDiscriminatorField(MSpecParser.DiscriminatorFieldContext ctx) {
        var name = getIdString(ctx.name);
        var field = new DefaultDiscriminatorField(getAttributes(ctx), name, consumePendingComment(ctx));
        getTypeReference(ctx.type).whenComplete((typeReference, throwable) -> {
            if (throwable != null) {
                // TODO: proper error collection in type context error bucket
                LOGGER.debug("Error setting type for {}", field, throwable);
                return;
            }
            field.setType(typeReference);
        });
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterEnumField(MSpecParser.EnumFieldContext ctx) {
        var typeRefName = ctx.type.complexTypeReference.getText();
        var type = new DefaultEnumTypeReference(typeRefName, null);
        setOrScheduleTypeDefinitionConsumer(typeRefName, type::setTypeDefinition);
        var name = getIdString(ctx.name);
        String fieldName = null;
        if (ctx.fieldName != null) {
            fieldName = getIdString(ctx.fieldName);
        }
        var field = new DefaultEnumField(getAttributes(ctx), type, name, fieldName, consumePendingComment(ctx));
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterImplicitField(MSpecParser.ImplicitFieldContext ctx) {
        var type = getSimpleTypeReference(ctx.type);
        var name = getIdString(ctx.name);
        var serializeExpression = getExpressionTerm(ctx.serializeExpression);
        var field = new DefaultImplicitField(getAttributes(ctx), type, name, serializeExpression, consumePendingComment(ctx));
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterManualArrayField(MSpecParser.ManualArrayFieldContext ctx) {
        String name = getIdString(ctx.name);
        ManualArrayField.LoopType loopType = ManualArrayField.LoopType.valueOf(
            ctx.loopType.getText().toUpperCase());
        Term loopExpression = getExpressionTerm(ctx.loopExpression);
        Term parseExpression = getExpressionTerm(ctx.parseExpression);
        Term serializeExpression = getExpressionTerm(ctx.serializeExpression);
        Term lengthExpression = getExpressionTerm(ctx.lengthExpression);
        DefaultManualArrayField field = new DefaultManualArrayField(getAttributes(ctx), name, loopType, loopExpression,
            parseExpression, serializeExpression, lengthExpression, consumePendingComment(ctx));
        getTypeReference(ctx.type).whenComplete((typeReference, throwable) -> {
            if (throwable != null) {
                // TODO: proper error collection in type context error bucket
                LOGGER.debug("Error setting type for {}", field, throwable);
                return;
            }
            field.setType(new DefaultArrayTypeReference(typeReference));
        });
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterManualField(MSpecParser.ManualFieldContext ctx) {
        String name = getIdString(ctx.name);
        Term parseExpression = getExpressionTerm(ctx.parseExpression);
        Term serializeExpression = getExpressionTerm(ctx.serializeExpression);
        Term lengthExpression = getExpressionTerm(ctx.lengthExpression);
        DefaultManualField field = new DefaultManualField(getAttributes(ctx), name, parseExpression,
            serializeExpression,
            lengthExpression, consumePendingComment(ctx));
        getTypeReference(ctx.type).whenComplete((typeReference, throwable) -> {
            if (throwable != null) {
                // TODO: proper error collection in type context error bucket
                LOGGER.debug("Error setting type for {}", field, throwable);
                return;
            }
            field.setType(typeReference);
        });
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterOptionalField(MSpecParser.OptionalFieldContext ctx) {
        String name = getIdString(ctx.name);
        Term conditionExpression = null;
        if (ctx.condition != null) {
            conditionExpression = getExpressionTerm(ctx.condition);
        }
        var field = new DefaultOptionalField(getAttributes(ctx), name, conditionExpression, consumePendingComment(ctx));
        getTypeReference(ctx.type).whenComplete((typeReference, throwable) -> {
            if (throwable != null) {
                // TODO: proper error collection in type context error bucket
                LOGGER.debug("Error setting type for {}", field, throwable);
                return;
            }
            field.setType(typeReference);
        });
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterPaddingField(MSpecParser.PaddingFieldContext ctx) {
        var type = getSimpleTypeReference(ctx.type);
        var name = getIdString(ctx.name);
        var paddingValue = getExpressionTerm(ctx.paddingValue);
        var timesPadding = getExpressionTerm(ctx.timesPadding);
        var field = new DefaultPaddingField(getAttributes(ctx), type, name, paddingValue, timesPadding, consumePendingComment(ctx));
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterPeekField(MSpecParser.PeekFieldContext ctx) {
        var name = getIdString(ctx.name);
        Term offsetExpression = null;
        if (ctx.offset != null) {
            offsetExpression = getExpressionTerm(ctx.offset);
        }
        var field = new DefaultPeekField(getAttributes(ctx), name, offsetExpression, consumePendingComment(ctx));
        getTypeReference(ctx.type).whenComplete((typeReference, throwable) -> {
            if (throwable != null) {
                // TODO: proper error collection in type context error bucket
                LOGGER.debug("Error setting type for {}", field, throwable);
                return;
            }
            field.setType(typeReference);
        });
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterReservedField(MSpecParser.ReservedFieldContext ctx) {
        var type = getSimpleTypeReference(ctx.type);
        var expected = getExprString(ctx.expected);
        var field = new DefaultReservedField(getAttributes(ctx), type, expected, consumePendingComment(ctx));
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterSimpleField(MSpecParser.SimpleFieldContext ctx) {
        var name = getIdString(ctx.name);
        var field = new DefaultSimpleField(getAttributes(ctx), name, consumePendingComment(ctx));
        getTypeReference(ctx.type).whenComplete((typeReference, throwable) -> {
            if (throwable != null) {
                // TODO: proper error collection in type context error bucket
                LOGGER.debug("Error setting type for {}", field, throwable);
                return;
            }
            field.setType(typeReference);
        });
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterStateField(MSpecParser.StateFieldContext ctx) {
        var name = getIdString(ctx.name);
        // Get the type information from the parents arguments.
        var parent = ctx.getParent().getParent().getParent();
        if (!(parent instanceof MSpecParser.ComplexTypeContext)
            && !(parent instanceof MSpecParser.CaseStatementContext)) {
            throw new RuntimeException("state fields must only be defined in complex types");
        }
        Optional<MSpecParser.ArgumentContext> argumentContext;
        if (parent instanceof MSpecParser.ComplexTypeContext complexTypeContext) {
            argumentContext = complexTypeContext.params.argument().stream()
                .filter(argContext -> argContext.name.getText().equalsIgnoreCase(name)).findFirst();
        } else {
            RuleContext curContext = parent;
            while ((curContext.parent != null) && !(curContext.parent instanceof MSpecParser.ComplexTypeContext)) {
                curContext = curContext.parent;
            }
            if (curContext.parent == null) {
                throw new RuntimeException(
                    "state fields must refer to arguments by using the same name. Parent context is null.");
            }
            var complexTypeContext = (MSpecParser.ComplexTypeContext) curContext.parent;
            argumentContext = complexTypeContext.params.argument().stream()
                .filter(argContext -> argContext.name.getText().equalsIgnoreCase(name)).findFirst();
        }
        var type = argumentContext.orElseThrow(
            () -> new RuntimeException("state fields must refer to arguments by using the same name.")).type;
        // The variable term is always just a direct reference to the parser argument.
        var valueExpression = new DefaultVariableLiteral(name, null, null, null);
        var field = new DefaultStateField(getAttributes(ctx), name, valueExpression, consumePendingComment(ctx));
        getTypeReference(type).whenComplete((typeReference, throwable) -> {
            if (throwable != null) {
                // TODO: proper error collection in type context error bucket
                LOGGER.debug("Error setting type for {}", field, throwable);
                return;
            }
            field.setType(typeReference);
        });
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterTypeSwitchField(MSpecParser.TypeSwitchFieldContext ctx) {
        var variableLiterals = ctx.discriminators.variableLiteral().stream()
            .map(this::getVariableLiteral)
            .collect(Collectors.toList());
        var field = new DefaultSwitchField(getAttributes(ctx), variableLiterals, consumePendingComment(ctx));
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterUnknownField(MSpecParser.UnknownFieldContext ctx) {
        var type = getSimpleTypeReference(ctx.type);
        var field = new DefaultUnknownField(getAttributes(ctx), type, consumePendingComment(ctx));
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterValidationField(MSpecParser.ValidationFieldContext ctx) {
        var validationExpression = getExpressionTerm(ctx.validationExpression);
        var shouldFail = true;
        if (ctx.shouldFail != null) {
            shouldFail = "true".equalsIgnoreCase(ctx.shouldFail.getText());
        }
        String description = null;
        if (ctx.description != null) {
            description = ctx.description.getText();
        }
        var field = new DefaultValidationField(getAttributes(ctx), validationExpression, description, shouldFail, consumePendingComment(ctx));
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterVirtualField(MSpecParser.VirtualFieldContext ctx) {
        var name = getIdString(ctx.name);
        var valueExpression = getExpressionTerm(ctx.valueExpression);
        var field = new DefaultVirtualField(getAttributes(ctx), name, valueExpression, consumePendingComment(ctx));
        getTypeReference(ctx.type).whenComplete((typeReference, throwable) -> {
            if (throwable != null) {
                // TODO: proper error collection in type context error bucket
                LOGGER.debug("Error setting type for {}", field, throwable);
                return;
            }
            field.setType(typeReference);
        });
        if (parserContexts.peek() != null) {
            parserContexts.peek().add(field);
        }
    }

    @Override
    public void enterCaseStatement(MSpecParser.CaseStatementContext ctx) {
        List<Field> parserContext = new LinkedList<>();

        // Calculate the name of the current type
        var namePrefix = "";
        if (ctx.nameWildcard != null) {
            namePrefix = getCurrentTypeName();
        }
        var typeName = namePrefix + ctx.name.getText();

        currentTypeName.push(typeName);
        // For DataIo we don't generate types.
        if (ctx.parent.parent instanceof MSpecParser.DataIoDefinitionContext) {
            currentTypeName.pop();
        }

        parserContexts.push(parserContext);
    }

    @Override
    public void exitCaseStatement(MSpecParser.CaseStatementContext ctx) {
        var typeName = currentTypeName.pop();
        // For DataIo we don't generate types.
        if (ctx.parent.parent instanceof MSpecParser.DataIoDefinitionContext) {
            currentTypeName.push(typeName);
            typeName = ctx.name.getText();
        }

        var abstractType = getSwitchField() != null;

        var attributes = batchSetAttributes.peek();

        var parserArguments = new LinkedList<Argument>();
        // For DataIO types, add all the arguments from the parent type.
        /*
         * if (!(ctx.parent.parent.parent.parent instanceof
         * MSpecParser.ComplexTypeContext)
         * && ((MSpecParser.ComplexTypeContext) ctx.parent.parent.parent).params !=
         * null) {
         * parserArguments.addAll(getParserArguments(
         * ((MSpecParser.ComplexTypeContext)
         * ctx.parent.parent.parent).params.argument()));
         * }
         */
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
        var fields = parserContexts.pop();
        var type = new DefaultDiscriminatedComplexTypeDefinition(typeName, attributes, parserArguments, abstractType, fields, discriminatorValues, consumePendingComment(ctx));
        // Link the fields and the complex types.
        if (fields != null) {
            fields.forEach(field -> ((DefaultField) field).setOwner(type));
        }

        // For DataIO we don't need to generate the sub-types as these will be
        // PlcValues.
        if (!(ctx.parent.parent instanceof MSpecParser.DataIoDefinitionContext)) {
            dispatchType(typeName, type);
        }

        // Add the type to the switch field definition.
        var switchField = getSwitchField();
        if (switchField == null) {
            throw new RuntimeException("This shouldn't have happened");
        }
        switchField.addCase(type);
    }

    @Override
    public void enterEnumValueDefinition(MSpecParser.EnumValueDefinitionContext ctx) {
        var value = (ctx.valueExpression != null) ? unquoteString(ctx.valueExpression.getText()) : null;
        var name = ctx.name.getText();
        Map<String, String> constants = null;
        if (ctx.constantValueExpressions != null) {
            var parentCtx = (MSpecParser.ComplexTypeContext) ctx.parent;
            var numConstantValues = parentCtx.params.argument().size();
            var numExpressionValues = ctx.constantValueExpressions.expression().size();
            // This only works if we provide exactly the same number of expressions as we
            // defined constants
            if (numConstantValues != numExpressionValues) {
                throw new RuntimeException("Number of constant value expressions doesn't match the number of " +
                    "defined constants. Expecting " + numConstantValues + " but got " + numExpressionValues);
            }

            // Build a map of the constant expressions (With the constant name as key)
            constants = new HashMap<>();
            for (int i = 0; i < numConstantValues; i++) {
                var argumentContext = parentCtx.params.argument(i);
                var constantName = argumentContext.name.getText();
                var expression = ctx.constantValueExpressions.expression(i);
                var constant = unquoteString(expression.getText());
                // String expressions are double escaped
                if (constant != null && constant.startsWith("\"")) {
                    constant = unquoteString(constant);
                }
                constants.put(constantName, constant);
            }
        }
        var enumValues = Objects.requireNonNull(this.enumContexts.peek());
        if (value == null) {
            // If no values are specified we count
            var counted = "0";
            if (!enumValues.isEmpty()) {
                var previousValue = enumValues.getLast().getValue();
                var parsedPreviousValue = Integer.parseInt(previousValue);
                counted = "" + (parsedPreviousValue + 1);
            }
            value = counted;
        }
        var enumValue = new DefaultEnumValue(value, name, constants, consumePendingComment(ctx));
        enumValues.add(enumValue);
    }

    private Term getExpressionTerm(MSpecParser.ExpressionContext expressionContext) {
        if (expressionContext.ASTERISK() != null) {
            return WildcardTerm.INSTANCE;
        }
        var expressionString = getExprString(expressionContext);
        Objects.requireNonNull(expressionString, "Expression string should not be null");
        var inputStream = IOUtils.toInputStream(expressionString, Charset.defaultCharset());

        Objects.requireNonNull(getCurrentTypeName(), "expression term can only occur within a type");
        var parser = new ExpressionStringParser(this, getCurrentTypeName());
        try {
            return parser.parse(inputStream);
        } catch (Exception e) {
            throw new RuntimeException(String.format("Error parsing expression: '%s' at line %d column %d",
                expressionString, expressionContext.start.getLine(), expressionContext.start.getStartIndex()), e);
        }
    }

    private VariableLiteral getVariableLiteral(MSpecParser.VariableLiteralContext variableLiteralContext) {
        // TODO: make nullsafe
        var variableLiteral = variableLiteralContext.getText();
        var inputStream = IOUtils.toInputStream(variableLiteral, Charset.defaultCharset());
        var parser = new ExpressionStringParser(this, getCurrentTypeName());
        try {
            // As this come from a VariableLiteralContext we know that it is a
            // VariableLiteral
            return (VariableLiteral) parser.parse(inputStream);
        } catch (Exception e) {
            throw new RuntimeException(String.format("Error parsing variable literal: '%s' at line %d column %d",
                variableLiteral, variableLiteralContext.start.getLine(),
                variableLiteralContext.start.getStartIndex()), e);
        }
    }

    private Literal getValueLiteral(MSpecParser.ValueLiteralContext valueLiteralContext) {
        // TODO: make nullsafe
        var valueLiteralContextText = valueLiteralContext.getText();
        var inputStream = IOUtils.toInputStream(valueLiteralContextText, Charset.defaultCharset());
        var parser = new ExpressionStringParser(this, getCurrentTypeName());
        try {
            // As this come from a ValueLiteralContext we know that it is a Literal
            return (Literal) parser.parse(inputStream);
        } catch (Exception e) {
            throw new RuntimeException(String.format("Error parsing variable literal: '%s' at line %d column %d",
                valueLiteralContextText, valueLiteralContext.start.getLine(),
                valueLiteralContext.start.getStartIndex()), e);
        }
    }

    private CompletionStage<TypeReference> getTypeReference(MSpecParser.TypeReferenceContext ctx) {
        if (ctx.simpleTypeReference != null) {
            return CompletableFuture.completedFuture(getSimpleTypeReference(ctx.simpleTypeReference));
        } else {
            CompletableFuture<TypeReference> typeReferenceCompletableFuture = new CompletableFuture<>();
            var typeRefName = ctx.complexTypeReference.getText();
            setOrScheduleTypeDefinitionConsumer(typeRefName, typeDefinition -> {
                if (typeDefinition.isDataIoTypeDefinition()) {
                    var value = new DefaultDataIoTypeReference(typeRefName,
                        getParams(ctx.params));
                    value.setTypeDefinition(typeDefinition);
                    typeReferenceCompletableFuture.complete(value);
                } else if (typeDefinition.isComplexTypeDefinition()) {
                    var value = new DefaultComplexTypeReference(typeRefName,
                        getParams(ctx.params));
                    value.setTypeDefinition(typeDefinition);
                    typeReferenceCompletableFuture.complete(value);
                } else if (typeDefinition.isEnumTypeDefinition()) {
                    var value = new DefaultEnumTypeReference(typeRefName, getParams(ctx.params));
                    value.setTypeDefinition(typeDefinition);
                    typeReferenceCompletableFuture.complete(value);
                } else {
                    throw new RuntimeException("Support for " + typeDefinition.getClass() + " not implemented yet");
                }
            });
            return typeReferenceCompletableFuture;
        }
    }

    private SimpleTypeReference getSimpleTypeReference(MSpecParser.DataTypeContext ctx) {
        var simpleBaseType = SimpleTypeReference.SimpleBaseType
            .valueOf(ctx.base.getText().toUpperCase());
        // String types need an additional length expression.
        if (simpleBaseType == SimpleTypeReference.SimpleBaseType.VSTRING) {
            if (ctx.length != null) {
                var lengthExpression = getExpressionTerm(ctx.length);
                return new DefaultVstringTypeReference(simpleBaseType, lengthExpression);
            } else {
                return new DefaultVstringTypeReference(simpleBaseType, null);
            }
        }
        switch (simpleBaseType) {
            case INT:
            case UINT:
                var integerSize = Integer.parseInt(ctx.size.getText());
                return new DefaultIntegerTypeReference(simpleBaseType, integerSize);
            case VINT: {
                var attributes = getAttributes(ctx.parent.parent);
                SimpleTypeReference propertyType;
                var propertySizeInBits = 32;
                if (attributes.containsKey("propertySizeInBits")) {
                    var propertySizeInBitsTerm = attributes.get("propertySizeInBits");
                    if (!(propertySizeInBitsTerm instanceof NumericLiteral propertySizeInBitsLiteral)) {
                        throw new RuntimeException(
                            "'propertySizeInBits' attribute is required to be a numeric literal");
                    }
                    propertySizeInBits = propertySizeInBitsLiteral.getNumber().intValue();
                }
                propertyType = new DefaultIntegerTypeReference(SimpleTypeReference.SimpleBaseType.INT,
                    propertySizeInBits);
                return new DefaultVintegerTypeReference(simpleBaseType, propertyType);
            }
            case VUINT: {
                final Map<String, Term> attributes = getAttributes(ctx.parent.parent);
                SimpleTypeReference propertyType;
                int propertySizeInBits = 32;
                if (attributes.containsKey("propertySizeInBits")) {
                    final Term propertySizeInBitsTerm = attributes.get("propertySizeInBits");
                    if (!(propertySizeInBitsTerm instanceof NumericLiteral propertySizeInBitsLiteral)) {
                        throw new RuntimeException(
                            "'propertySizeInBits' attribute is required to be a numeric literal");
                    }
                    propertySizeInBits = propertySizeInBitsLiteral.getNumber().intValue();
                }
                propertyType = new DefaultIntegerTypeReference(SimpleTypeReference.SimpleBaseType.UINT,
                    propertySizeInBits);
                return new DefaultVintegerTypeReference(simpleBaseType, propertyType);
            }
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
        for (var field : Objects.requireNonNull(parserContexts.peek())) {
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
            .map(argumentContext -> {
                DefaultArgument argument = new DefaultArgument(getIdString(argumentContext.name));
                getTypeReference(argumentContext.type).whenComplete((typeReference, throwable) -> {
                    if (throwable != null) {
                        // TODO: proper error collection in type context error bucket
                        LOGGER.debug("Error setting type for {}", argument, throwable);
                        return;
                    }
                    argument.setType(typeReference);
                });
                return argument;
            })
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
        var inputStream = IOUtils.toInputStream(expressionString, Charset.defaultCharset());
        var parser = new ExpressionStringParser(this, getCurrentTypeName());
        try {
            return parser.parse(inputStream);
        } catch (Exception e) {
            throw new RuntimeException("Error parsing expression: '" + expressionString + "'", e);
        }
    }

    private String getCurrentTypeName() {
        if (currentTypeName.isEmpty()) {
            return null;
        }
        return currentTypeName.peek();
    }

    private Map<String, Term> getAttributes(RuleContext ctx) {
        Map<String, Term> attributes = new HashMap<>();
        // Add any attributes from the current batchSet context.
        if (!batchSetAttributes.empty()) {
            attributes.putAll(batchSetAttributes.peek());
        }
        // Add any attributes set on the current field itself.
        if (ctx.parent.parent instanceof MSpecParser.FieldDefinitionContext fieldDefinitionContext) {
            for (MSpecParser.AttributeContext attributeContext : fieldDefinitionContext.attributes.attribute()) {
                Term attributeExpression = getExpressionTerm(attributeContext.value);
                attributes.put(attributeContext.name.getText(), attributeExpression);
            }
        }
        return attributes;
    }

    private String unquoteString(String quotedString) {
        if (quotedString == null || quotedString.length() < 2) {
            return quotedString;
        }
        return quotedString.substring(1, quotedString.length() - 1);
    }

    private String getIdString(MSpecParser.IdExpressionContext ctx) {
        if (ctx.id == null) {
            return null;
        }
        return ctx.id.getText();
    }

    private String getExprString(MSpecParser.ExpressionContext ctx) {
        if (ctx.expr == null) {
            return null;
        }
        return ctx.expr.getText();
    }

    public void dispatchType(String typeName, TypeDefinition type) {
        LOGGER.debug("dispatching {}:{}", typeName, type);

        if (types.containsKey(typeName)) {
            LOGGER.warn("{} being overridden", typeName);
            // TODO: we need to implement replace logic... means we need to replace all old
            // references with the new one in that case otherwise we just get an exception
        }

        types.put(typeName, type);

        while (!typeDefinitionConsumers.getOrDefault(typeName, new LinkedList<>()).isEmpty()) {

            consumerDispatchType(typeName, type);
        }

        typeDefinitionConsumers.remove(typeName);
    }

    private void consumerDispatchType(String typeName, TypeDefinition type) {
        var waitingConsumers = typeDefinitionConsumers.getOrDefault(typeName, new LinkedList<>());
        LOGGER.debug("{} waiting for {}", waitingConsumers.size(), typeName);

        var consumerIterator = waitingConsumers.iterator();
        var removedItems = new ArrayList<Consumer<TypeDefinition>>();

        while (consumerIterator.hasNext()) {
            Consumer<TypeDefinition> setter = consumerIterator.next();
            LOGGER.debug("setting {} for {}", typeName, setter);
            removedItems.add(setter);
        }

        waitingConsumers.removeAll(removedItems);

        for (Consumer<TypeDefinition> setter : removedItems) {
            setter.accept(type);
        }
    }

    @Override
    public void setOrScheduleTypeDefinitionConsumer(String typeRefName, Consumer<TypeDefinition> setTypeDefinition) {
        LOGGER.debug("set or schedule {}", typeRefName);

        var typeDefinition = types.get(typeRefName);
        if (typeDefinition != null) {
            LOGGER.debug("{} present so setting for {}", typeRefName, setTypeDefinition);
            setTypeDefinition.accept(typeDefinition);
        } else {
            // put up order
            if (LOGGER.isDebugEnabled()) {
                LOGGER.debug("{} already waiting for {}",
                    typeDefinitionConsumers.getOrDefault(typeRefName, new LinkedList<>()).size(), typeRefName);
            }
            typeDefinitionConsumers.putIfAbsent(typeRefName, new LinkedList<>());
            typeDefinitionConsumers.get(typeRefName).add(setTypeDefinition);
        }
    }

}
