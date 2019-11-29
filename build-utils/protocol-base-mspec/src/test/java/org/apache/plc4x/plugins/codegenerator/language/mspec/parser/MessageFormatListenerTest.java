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

import org.antlr.v4.runtime.CommonToken;
import org.apache.plc4x.plugins.codegenerator.language.mspec.MSpecParser;
import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.Map;

import static org.hamcrest.CoreMatchers.not;
import static org.hamcrest.CoreMatchers.nullValue;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.hasSize;

// TODO: implement me
@ExtendWith(MockitoExtension.class)
class MessageFormatListenerTest {

    MessageFormatListener SUT;

    @BeforeEach
    void setUp() {
        SUT = new MessageFormatListener();
    }

    @Test
    void getTypes() {
        Map<String, TypeDefinition> types = SUT.getTypes();
        assertThat(types, nullValue());
    }

    @Test
    void enterFile() {
        assertThat(SUT.getParserContexts(), nullValue());
        assertThat(SUT.getEnumContexts(), nullValue());
        assertThat(SUT.getTypes(), nullValue());
        SUT.enterFile(null);
        assertThat(SUT.getParserContexts(), not(nullValue()));
        assertThat(SUT.getEnumContexts(), not(nullValue()));
        assertThat(SUT.getTypes(), not(nullValue()));
    }

    @Test
    void enterComplexType() {
        assertThat(SUT.getParserContexts(), nullValue());
        assertThat(SUT.getEnumContexts(), nullValue());
        assertThat(SUT.getTypes(), nullValue());
        SUT.enterFile(null);
        MSpecParser.ComplexTypeContext complexTypeContext = new MSpecParser.ComplexTypeContext(null, 0);
        SUT.enterComplexType(complexTypeContext);
        assertThat(SUT.getParserContexts(), hasSize(1));
        assertThat(SUT.getEnumContexts(), hasSize(0));
        assertThat(SUT.getTypes().values(), hasSize(0));
        complexTypeContext.enumValues = new MSpecParser.EnumValueDefinitionContext(null, 0);
        SUT.enterComplexType(complexTypeContext);
        assertThat(SUT.getParserContexts(), hasSize(1));
        assertThat(SUT.getEnumContexts(), hasSize(1));
        assertThat(SUT.getTypes().values(), hasSize(0));
    }

    @Test
    void exitComplexType() {
        MSpecParser.ComplexTypeContext complexTypeContext = new MSpecParser.ComplexTypeContext(null, 0);
        complexTypeContext.name = new MSpecParser.IdExpressionContext(null, 0);
        complexTypeContext.name.id = new CommonToken(0);
        SUT.exitComplexType(complexTypeContext);
    }

    @Test
    void enterArrayField() {
    }

    @Test
    void enterChecksumField() {
    }

    @Test
    void enterConstField() {
    }

    @Test
    void enterDiscriminatorField() {
    }

    @Test
    void enterEnumField() {
    }

    @Test
    void enterImplicitField() {
    }

    @Test
    void enterManualArrayField() {
    }

    @Test
    void enterManualField() {
    }

    @Test
    void enterOptionalField() {
    }

    @Test
    void enterPaddingField() {
    }

    @Test
    void enterReservedField() {
    }

    @Test
    void enterSimpleField() {
    }

    @Test
    void enterTypeSwitchField() {
    }

    @Test
    void enterVirtualField() {
    }

    @Test
    void enterCaseStatement() {
    }

    @Test
    void exitCaseStatement() {
    }

    @Test
    void enterEnumValueDefinition() {
    }

    @Test
    void main() {
    }
}