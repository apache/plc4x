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

import org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions.DefaultArgument;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions.DefaultComplexTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions.DefaultEnumTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions.DefaultTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.fields.DefaultSimpleField;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.fields.DefaultTypedField;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.fields.DefaultTypedNamedField;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.fields.DefaultVirtualField;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.references.*;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.terms.DefaultBooleanLiteral;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.terms.DefaultVariableLiteral;
import org.apache.plc4x.plugins.codegenerator.protocol.TypeContext;
import org.apache.plc4x.plugins.codegenerator.types.enums.EnumValue;
import org.apache.plc4x.plugins.codegenerator.types.references.SimpleTypeReference;
import org.junit.jupiter.api.Test;

import java.util.Map;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.InstanceOfAssertFactories.type;
import static org.junit.jupiter.api.Assertions.assertThrows;

class MessageFormatParserTest {

    final MessageFormatParser SUT = new MessageFormatParser();

    @Test
    void parseNull() {
        assertThrows(NullPointerException.class, () -> SUT.parse(null));
    }

    @Test
    void parseSomething() {
        TypeContext parse = SUT.parse(getClass().getResourceAsStream("/mspec.example"));
        assertThat(parse)
            .isNotNull();
    }

    @Test
    void parseSomethingElse() {
        TypeContext parse = SUT.parse(getClass().getResourceAsStream("/mspec.example2"));
        assertThat(parse)
            .isNotNull();
    }

    @Test
    void parseNothingElse() {
        TypeContext typeContext = SUT.parse(getClass().getResourceAsStream("/mspec.example3"));
        assertThat(typeContext)
            .extracting(TypeContext::getUnresolvedTypeReferences)
            .extracting(Map::size)
            .isEqualTo(0);

        assertThat(typeContext)
            .extracting(TypeContext::getTypeDefinitions)
            .satisfies(stringTypeDefinitionMap -> {
                assertThat(stringTypeDefinitionMap)
                    .hasEntrySatisfying("A", typeDefinition ->
                        assertThat(typeDefinition)
                            .asInstanceOf(type(DefaultComplexTypeDefinition.class))
                            .satisfies(defaultComplexTypeDefinition -> {
                                assertThat(defaultComplexTypeDefinition)
                                    .extracting(DefaultTypeDefinition::getName)
                                    .isEqualTo("A");
                                assertThat(defaultComplexTypeDefinition)
                                    .extracting(DefaultComplexTypeDefinition::getFields)
                                    .satisfies(fields -> {
                                        assertThat(fields)
                                            .element(0)
                                            .asInstanceOf(type(DefaultSimpleField.class))
                                            .satisfies(defaultSimpleField -> {
                                                assertThat(defaultSimpleField)
                                                    .extracting(DefaultTypedNamedField::getName)
                                                    .isEqualTo("b");
                                                assertThat(defaultSimpleField)
                                                    .extracting(DefaultTypedField::getType)
                                                    .asInstanceOf(type(DefaultComplexTypeReference.class))
                                                    .extracting(DefaultComplexTypeReference::getName)
                                                    .isEqualTo("B");
                                            });
                                    });
                            }))
                    .hasEntrySatisfying("B", typeDefinition ->
                        assertThat(typeDefinition)
                            .asInstanceOf(type(DefaultComplexTypeDefinition.class))
                            .satisfies(defaultComplexTypeDefinition -> {
                                assertThat(defaultComplexTypeDefinition)
                                    .extracting(DefaultTypeDefinition::getName)
                                    .isEqualTo("B");
                                assertThat(defaultComplexTypeDefinition)
                                    .extracting(DefaultComplexTypeDefinition::getFields)
                                    .satisfies(fields -> {
                                        assertThat(fields)
                                            .element(0)
                                            .asInstanceOf(type(DefaultSimpleField.class))
                                            .satisfies(defaultSimpleField -> {
                                                assertThat(defaultSimpleField)
                                                    .extracting(DefaultTypedNamedField::getName)
                                                    .isEqualTo("c");
                                                assertThat(defaultSimpleField)
                                                    .extracting(DefaultTypedField::getType)
                                                    .asInstanceOf(type(DefaultComplexTypeReference.class))
                                                    .extracting(DefaultComplexTypeReference::getName)
                                                    .isEqualTo("C");
                                            });
                                    });
                            })
                    )
                    .hasEntrySatisfying("C", typeDefinition ->
                        assertThat(typeDefinition)
                            .asInstanceOf(type(DefaultComplexTypeDefinition.class))
                            .satisfies(defaultComplexTypeDefinition -> {
                                assertThat(defaultComplexTypeDefinition)
                                    .extracting(DefaultTypeDefinition::getName)
                                    .isEqualTo("C");
                                assertThat(defaultComplexTypeDefinition)
                                    .extracting(DefaultComplexTypeDefinition::getFields)
                                    .satisfies(fields -> {
                                        assertThat(fields)
                                            .element(0)
                                            .asInstanceOf(type(DefaultSimpleField.class))
                                            .satisfies(defaultSimpleField -> {
                                                assertThat(defaultSimpleField)
                                                    .extracting(DefaultTypedNamedField::getName)
                                                    .isEqualTo("onlyOneField");
                                                assertThat(defaultSimpleField)
                                                    .extracting(DefaultTypedField::getType)
                                                    .asInstanceOf(type(DefaultBooleanTypeReference.class))
                                                    .extracting(DefaultBooleanTypeReference::getBaseType)
                                                    .isEqualTo(SimpleTypeReference.SimpleBaseType.BIT);
                                            });
                                        assertThat(fields)
                                            .element(1)
                                            .asInstanceOf(type(DefaultVirtualField.class))
                                            .satisfies(defaultVirtualField -> {
                                                assertThat(defaultVirtualField)
                                                    .extracting(DefaultTypedNamedField::getName)
                                                    .isEqualTo("secondField");
                                                assertThat(defaultVirtualField)
                                                    .extracting(DefaultVirtualField::getValueExpression)
                                                    .asInstanceOf(type(DefaultVariableLiteral.class))
                                                    .satisfies(defaultVariableLiteral -> {
                                                        assertThat(defaultVariableLiteral)
                                                            .extracting(DefaultVariableLiteral::getName)
                                                            .isEqualTo("onlyOneField");
                                                        assertThat(defaultVariableLiteral)
                                                            .extracting(DefaultVariableLiteral::getTypeReference)
                                                            .asInstanceOf(type(DefaultBooleanTypeReference.class))
                                                            .extracting(AbstractSimpleTypeReference::getBaseType)
                                                            .isEqualTo(SimpleTypeReference.SimpleBaseType.BIT);
                                                    });
                                                assertThat(defaultVirtualField)
                                                    .extracting(DefaultTypedField::getType)
                                                    .asInstanceOf(type(DefaultBooleanTypeReference.class))
                                                    .extracting(DefaultBooleanTypeReference::getBaseType)
                                                    .isEqualTo(SimpleTypeReference.SimpleBaseType.BIT);
                                            });
                                        assertThat(fields)
                                            .element(2)
                                            .asInstanceOf(type(DefaultSimpleField.class))
                                            .satisfies(defaultSimpleField -> {
                                                assertThat(defaultSimpleField)
                                                    .extracting(DefaultTypedNamedField::getName)
                                                    .isEqualTo("d");
                                                assertThat(defaultSimpleField)
                                                    .extracting(DefaultTypedField::getType)
                                                    .asInstanceOf(type(DefaultEnumTypeReference.class))
                                                    .extracting(DefaultEnumTypeReference::getName)
                                                    .isEqualTo("D");
                                            });
                                        assertThat(fields)
                                            .element(3)
                                            .asInstanceOf(type(DefaultVirtualField.class))
                                            .satisfies(defaultVirtualField -> {
                                                assertThat(defaultVirtualField)
                                                    .extracting(DefaultTypedNamedField::getName)
                                                    .isEqualTo("thirdField");
                                                assertThat(defaultVirtualField)
                                                    .extracting(DefaultVirtualField::getValueExpression)
                                                    .asInstanceOf(type(DefaultVariableLiteral.class))
                                                    .satisfies(defaultVariableLiteral -> {
                                                        assertThat(defaultVariableLiteral)
                                                            .extracting(DefaultVariableLiteral::getName)
                                                            .isEqualTo("d");
                                                        assertThat(defaultVariableLiteral)
                                                            .extracting(DefaultVariableLiteral::getTypeReference)
                                                            .asInstanceOf(type(DefaultEnumTypeReference.class))
                                                            .extracting(DefaultEnumTypeReference::getName)
                                                            .isEqualTo("D");
                                                    });
                                                assertThat(defaultVirtualField)
                                                    .extracting(DefaultTypedField::getType)
                                                    .asInstanceOf(type(DefaultEnumTypeReference.class))
                                                    .extracting(DefaultEnumTypeReference::getName)
                                                    .isEqualTo("D");
                                            });
                                        assertThat(fields)
                                            .element(4)
                                            .asInstanceOf(type(DefaultSimpleField.class))
                                            .satisfies(defaultSimpleField -> {
                                                assertThat(defaultSimpleField)
                                                    .extracting(DefaultTypedNamedField::getName)
                                                    .isEqualTo("nothingElseMatters");
                                                assertThat(defaultSimpleField)
                                                    .extracting(DefaultTypedField::getType)
                                                    .asInstanceOf(type(DefaultIntegerTypeReference.class))
                                                    .extracting(DefaultIntegerTypeReference::getBaseType, AbstractSimpleTypeReference::getSizeInBits)
                                                    .containsExactly(SimpleTypeReference.SimpleBaseType.UINT, 8);
                                            });
                                        assertThat(fields)
                                            .element(5)
                                            .asInstanceOf(type(DefaultSimpleField.class))
                                            .satisfies(defaultSimpleField -> {
                                                assertThat(defaultSimpleField)
                                                    .extracting(DefaultTypedNamedField::getName)
                                                    .isEqualTo("e");
                                                assertThat(defaultSimpleField)
                                                    .extracting(DefaultTypedField::getType)
                                                    .asInstanceOf(type(DefaultComplexTypeReference.class))
                                                    .satisfies(defaultComplexTypeReference -> {
                                                        assertThat(defaultComplexTypeReference)
                                                            .extracting(DefaultComplexTypeReference::getName)
                                                            .isEqualTo("E");
                                                        assertThat(defaultComplexTypeReference)
                                                            .extracting(DefaultComplexTypeReference::getParams)
                                                            .satisfies(paramsOptional ->
                                                                assertThat(paramsOptional)
                                                                    .hasValueSatisfying(params->
                                                                            assertThat(params)
                                                                                .element(0)
                                                                                .asInstanceOf(type(DefaultBooleanLiteral.class))
                                                                                .extracting(DefaultBooleanLiteral::getValue)
                                                                                .isEqualTo(true)
                                                                        )
                                                            );
                                                    });
                                            });
                                    });
                            })
                    )
                    .hasEntrySatisfying("D", typeDefinition ->
                        assertThat(typeDefinition)
                            .asInstanceOf(type(DefaultEnumTypeDefinition.class))
                            .satisfies(defaultEnumTypeDefinition -> {
                                assertThat(defaultEnumTypeDefinition)
                                    .extracting(DefaultTypeDefinition::getName)
                                    .isEqualTo("D");
                                assertThat(defaultEnumTypeDefinition)
                                    .extracting(DefaultEnumTypeDefinition::getType)
                                    .satisfies(simpleTypeReferenceOptional ->
                                        assertThat(simpleTypeReferenceOptional)
                                            .hasValueSatisfying(simpleTypeReference ->
                                                assertThat(simpleTypeReference)
                                                    .extracting(SimpleTypeReference::getBaseType, SimpleTypeReference::getSizeInBits)
                                                    .containsExactly(SimpleTypeReference.SimpleBaseType.UINT, 32)
                                            )
                                    );
                                assertThat(defaultEnumTypeDefinition)
                                    .extracting(DefaultEnumTypeDefinition::getEnumValues)
                                    .satisfies(enumValues ->
                                        assertThat(enumValues)
                                            .element(0)
                                            .extracting(EnumValue::getName, EnumValue::getValue)
                                            .containsExactly("D", "0x1")
                                    );
                            })
                    )
                    .hasEntrySatisfying("E", typeDefinition ->
                        assertThat(typeDefinition)
                            .asInstanceOf(type(DefaultComplexTypeDefinition.class))
                            .satisfies(defaultComplexTypeDefinition -> {
                                assertThat(defaultComplexTypeDefinition)
                                    .extracting(DefaultComplexTypeDefinition::getName)
                                    .isEqualTo("E");
                                assertThat(defaultComplexTypeDefinition)
                                    .extracting(DefaultComplexTypeDefinition::getFields)
                                    .satisfies(fields ->
                                        assertThat(fields)
                                            .element(0)
                                            .asInstanceOf(type(DefaultSimpleField.class))
                                            .extracting(DefaultTypedNamedField::getName)
                                            .isEqualTo("eField")
                                    );
                                assertThat(defaultComplexTypeDefinition)
                                    .extracting(DefaultComplexTypeDefinition::getParserArguments)
                                    .satisfies(parserArgumentsOptional ->
                                        assertThat(parserArgumentsOptional)
                                            .hasValueSatisfying(arguments -> {
                                                assertThat(arguments)
                                                    .element(0)
                                                    .asInstanceOf(type(DefaultArgument.class))
                                                    .extracting(DefaultArgument::getName)
                                                    .isEqualTo("aBit");
                                                assertThat(arguments)
                                                    .element(0)
                                                    .asInstanceOf(type(DefaultArgument.class))
                                                    .extracting(DefaultArgument::getType)
                                                    .satisfies(typeReference ->
                                                        assertThat(typeReference)
                                                            .asInstanceOf(type(DefaultBooleanTypeReference.class))
                                                            .extracting(AbstractSimpleTypeReference::getBaseType)
                                                            .isEqualTo(SimpleTypeReference.SimpleBaseType.BIT)
                                                    );
                                            })
                                    );
                            })
                    )
                    .hasEntrySatisfying("Root", typeDefinition ->
                        assertThat(typeDefinition)
                            .asInstanceOf(type(DefaultComplexTypeDefinition.class))
                            .satisfies(defaultComplexTypeDefinition -> {
                                assertThat(defaultComplexTypeDefinition)
                                    .extracting(DefaultComplexTypeDefinition::getName)
                                    .isEqualTo("Root");
                                assertThat(defaultComplexTypeDefinition)
                                    .extracting(DefaultComplexTypeDefinition::getFields)
                                    .satisfies(fields -> {
                                        assertThat(fields)
                                            .element(0)
                                            .asInstanceOf(type(DefaultSimpleField.class))
                                            .satisfies(defaultSimpleField -> {
                                                assertThat(defaultSimpleField)
                                                    .extracting(DefaultTypedNamedField::getName)
                                                    .isEqualTo("a");
                                                assertThat(defaultSimpleField)
                                                    .extracting(DefaultTypedField::getType)
                                                    .asInstanceOf(type(DefaultComplexTypeReference.class))
                                                    .extracting(DefaultComplexTypeReference::getName)
                                                    .isEqualTo("A");
                                            });
                                        assertThat(fields)
                                            .element(1)
                                            .asInstanceOf(type(DefaultVirtualField.class))
                                            .satisfies(defaultVirtualField -> {
                                                assertThat(defaultVirtualField)
                                                    .extracting(DefaultTypedNamedField::getName)
                                                    .isEqualTo("doesIt");
                                                assertThat(defaultVirtualField)
                                                    .extracting(DefaultVirtualField::getValueExpression)
                                                    .asInstanceOf(type(DefaultVariableLiteral.class))
                                                    .satisfies(defaultVariableLiteral -> {
                                                        assertThat(defaultVariableLiteral)
                                                            .extracting(DefaultVariableLiteral::getName)
                                                            .isEqualTo("a");
                                                        assertThat(defaultVariableLiteral)
                                                            .extracting(DefaultVariableLiteral::getTypeReference)
                                                            .asInstanceOf(type(DefaultComplexTypeReference.class))
                                                            .extracting(DefaultComplexTypeReference::getName)
                                                            .isEqualTo("A");
                                                        assertThat(defaultVariableLiteral)
                                                            .extracting(DefaultVariableLiteral::getChild)
                                                            .satisfies(childOptional ->
                                                                assertThat(childOptional)
                                                                    .hasValueSatisfying(childVariableLiteral -> {
                                                                        assertThat(childVariableLiteral)
                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                            .extracting(DefaultVariableLiteral::getName)
                                                                            .isEqualTo("b");
                                                                        assertThat(childVariableLiteral)
                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                            .extracting(DefaultVariableLiteral::getTypeReference)
                                                                            .asInstanceOf(type(DefaultComplexTypeReference.class))
                                                                            .extracting(DefaultComplexTypeReference::getName)
                                                                            .isEqualTo("B");
                                                                        assertThat(childVariableLiteral)
                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                            .extracting(DefaultVariableLiteral::getChild)
                                                                            .satisfies(childChildOptional ->
                                                                                assertThat(childChildOptional)
                                                                                    .hasValueSatisfying(childChildVariableLiteral -> {
                                                                                        assertThat(childChildVariableLiteral)
                                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                                            .extracting(DefaultVariableLiteral::getName)
                                                                                            .isEqualTo("c");
                                                                                        assertThat(childChildVariableLiteral)
                                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                                            .extracting(DefaultVariableLiteral::getTypeReference)
                                                                                            .asInstanceOf(type(DefaultComplexTypeReference.class))
                                                                                            .extracting(DefaultComplexTypeReference::getName)
                                                                                            .isEqualTo("C");
                                                                                        assertThat(childChildVariableLiteral)
                                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                                            .extracting(DefaultVariableLiteral::getChild)
                                                                                            .satisfies(childChildChildOptional ->
                                                                                                assertThat(childChildChildOptional)
                                                                                                    .hasValueSatisfying(childChildChildVariableLiteral -> {
                                                                                                        assertThat(childChildChildVariableLiteral)
                                                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                                                            .extracting(DefaultVariableLiteral::getName)
                                                                                                            .isEqualTo("nothingElseMatters");
                                                                                                        assertThat(childChildChildVariableLiteral)
                                                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                                                            .extracting(DefaultVariableLiteral::getTypeReference)
                                                                                                            .asInstanceOf(type(DefaultIntegerTypeReference.class))
                                                                                                            .extracting(DefaultIntegerTypeReference::getBaseType, AbstractSimpleTypeReference::getSizeInBits)
                                                                                                            .containsExactly(SimpleTypeReference.SimpleBaseType.UINT, 8);
                                                                                                    })
                                                                                            );
                                                                                    })
                                                                            );
                                                                    })
                                                            );
                                                    });
                                                assertThat(defaultVirtualField)
                                                    .extracting(DefaultTypedField::getType)
                                                    .asInstanceOf(type(AbstractSimpleTypeReference.class))
                                                    .extracting(AbstractSimpleTypeReference::getBaseType, AbstractSimpleTypeReference::getSizeInBits)
                                                    .containsExactly(SimpleTypeReference.SimpleBaseType.UINT, 8);
                                            });
                                        assertThat(fields)
                                            .element(2)
                                            .asInstanceOf(type(DefaultVirtualField.class))
                                            .satisfies(defaultVirtualField -> {
                                                assertThat(defaultVirtualField)
                                                    .extracting(DefaultTypedNamedField::getName)
                                                    .isEqualTo("thisNow");
                                                assertThat(defaultVirtualField)
                                                    .extracting(DefaultVirtualField::getValueExpression)
                                                    .asInstanceOf(type(DefaultVariableLiteral.class))
                                                    .satisfies(defaultVariableLiteral -> {
                                                        assertThat(defaultVariableLiteral)
                                                            .extracting(DefaultVariableLiteral::getName)
                                                            .isEqualTo("a");
                                                        assertThat(defaultVariableLiteral)
                                                            .extracting(DefaultVariableLiteral::getTypeReference)
                                                            .asInstanceOf(type(DefaultComplexTypeReference.class))
                                                            .extracting(DefaultComplexTypeReference::getName)
                                                            .isEqualTo("A");
                                                        assertThat(defaultVariableLiteral)
                                                            .extracting(DefaultVariableLiteral::getChild)
                                                            .satisfies(childOptional ->
                                                                assertThat(childOptional)
                                                                    .hasValueSatisfying(childVariableLiteral -> {
                                                                        assertThat(childVariableLiteral)
                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                            .extracting(DefaultVariableLiteral::getName)
                                                                            .isEqualTo("b");
                                                                        assertThat(childVariableLiteral)
                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                            .extracting(DefaultVariableLiteral::getTypeReference)
                                                                            .asInstanceOf(type(DefaultComplexTypeReference.class))
                                                                            .extracting(DefaultComplexTypeReference::getName)
                                                                            .isEqualTo("B");
                                                                        assertThat(childVariableLiteral)
                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                            .extracting(DefaultVariableLiteral::getChild)
                                                                            .satisfies(childChildOptional ->
                                                                                assertThat(childChildOptional)
                                                                                    .hasValueSatisfying(childChildVariableLiteral -> {
                                                                                        assertThat(childChildVariableLiteral)
                                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                                            .extracting(DefaultVariableLiteral::getName)
                                                                                            .isEqualTo("c");
                                                                                        assertThat(childChildVariableLiteral)
                                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                                            .extracting(DefaultVariableLiteral::getTypeReference)
                                                                                            .asInstanceOf(type(DefaultComplexTypeReference.class))
                                                                                            .extracting(DefaultComplexTypeReference::getName)
                                                                                            .isEqualTo("C");
                                                                                        assertThat(childChildVariableLiteral)
                                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                                            .extracting(DefaultVariableLiteral::getChild)
                                                                                            .satisfies(childChildChildOptional ->
                                                                                                assertThat(childChildChildOptional)
                                                                                                    .hasValueSatisfying(childChildChildVariableLiteral -> {
                                                                                                        assertThat(childChildChildVariableLiteral)
                                                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                                                            .extracting(DefaultVariableLiteral::getName)
                                                                                                            .isEqualTo("secondField");
                                                                                                        assertThat(childChildChildVariableLiteral)
                                                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                                                            .extracting(DefaultVariableLiteral::getTypeReference)
                                                                                                            .asInstanceOf(type(DefaultBooleanTypeReference.class))
                                                                                                            .extracting(AbstractSimpleTypeReference::getBaseType)
                                                                                                            .isEqualTo(SimpleTypeReference.SimpleBaseType.BIT);
                                                                                                    })
                                                                                            );
                                                                                    })
                                                                            );
                                                                    })
                                                            );
                                                    });
                                                assertThat(defaultVirtualField)
                                                    .extracting(DefaultTypedField::getType)
                                                    .asInstanceOf(type(AbstractSimpleTypeReference.class))
                                                    .extracting(AbstractSimpleTypeReference::getBaseType, AbstractSimpleTypeReference::getSizeInBits)
                                                    .containsExactly(SimpleTypeReference.SimpleBaseType.BIT, 1);
                                            });
                                        assertThat(fields)
                                            .element(3)
                                            .asInstanceOf(type(DefaultVirtualField.class))
                                            .satisfies(defaultVirtualField -> {
                                                assertThat(defaultVirtualField)
                                                    .extracting(DefaultTypedNamedField::getName)
                                                    .isEqualTo("thisNow2");
                                                assertThat(defaultVirtualField)
                                                    .extracting(DefaultVirtualField::getValueExpression)
                                                    .asInstanceOf(type(DefaultVariableLiteral.class))
                                                    .satisfies(defaultVariableLiteral -> {
                                                        assertThat(defaultVariableLiteral)
                                                            .extracting(DefaultVariableLiteral::getName)
                                                            .isEqualTo("a");
                                                        assertThat(defaultVariableLiteral)
                                                            .extracting(DefaultVariableLiteral::getTypeReference)
                                                            .asInstanceOf(type(DefaultComplexTypeReference.class))
                                                            .extracting(DefaultComplexTypeReference::getName)
                                                            .isEqualTo("A");
                                                        assertThat(defaultVariableLiteral)
                                                            .extracting(DefaultVariableLiteral::getChild)
                                                            .satisfies(childOptional ->
                                                                assertThat(childOptional)
                                                                    .hasValueSatisfying(childVariableLiteral -> {
                                                                        assertThat(childVariableLiteral)
                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                            .extracting(DefaultVariableLiteral::getName)
                                                                            .isEqualTo("b");
                                                                        assertThat(childVariableLiteral)
                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                            .extracting(DefaultVariableLiteral::getTypeReference)
                                                                            .asInstanceOf(type(DefaultComplexTypeReference.class))
                                                                            .extracting(DefaultComplexTypeReference::getName)
                                                                            .isEqualTo("B");
                                                                        assertThat(childVariableLiteral)
                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                            .extracting(DefaultVariableLiteral::getChild)
                                                                            .satisfies(childChildOptional ->
                                                                                assertThat(childChildOptional)
                                                                                    .hasValueSatisfying(childChildVariableLiteral -> {
                                                                                        assertThat(childChildVariableLiteral)
                                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                                            .extracting(DefaultVariableLiteral::getName)
                                                                                            .isEqualTo("c");
                                                                                        assertThat(childChildVariableLiteral)
                                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                                            .extracting(DefaultVariableLiteral::getTypeReference)
                                                                                            .asInstanceOf(type(DefaultComplexTypeReference.class))
                                                                                            .extracting(DefaultComplexTypeReference::getName)
                                                                                            .isEqualTo("C");
                                                                                        assertThat(childChildVariableLiteral)
                                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                                            .extracting(DefaultVariableLiteral::getChild)
                                                                                            .satisfies(childChildChildOptional ->
                                                                                                assertThat(childChildChildOptional)
                                                                                                    .hasValueSatisfying(childChildChildVariableLiteral -> {
                                                                                                        assertThat(childChildChildVariableLiteral)
                                                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                                                            .extracting(DefaultVariableLiteral::getName)
                                                                                                            .isEqualTo("thirdField");
                                                                                                        assertThat(childChildChildVariableLiteral)
                                                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                                                            .extracting(DefaultVariableLiteral::getTypeReference)
                                                                                                            .asInstanceOf(type(DefaultEnumTypeReference.class))
                                                                                                            .extracting(DefaultEnumTypeReference::getName)
                                                                                                            .isEqualTo("D");
                                                                                                    })
                                                                                            );
                                                                                    })
                                                                            );
                                                                    })
                                                            );
                                                    });
                                                assertThat(defaultVirtualField)
                                                    .extracting(DefaultTypedField::getType)
                                                    .asInstanceOf(type(AbstractSimpleTypeReference.class))
                                                    .extracting(AbstractSimpleTypeReference::getBaseType, AbstractSimpleTypeReference::getSizeInBits)
                                                    .containsExactly(SimpleTypeReference.SimpleBaseType.BIT, 1);
                                            });
                                        assertThat(fields)
                                            .element(4)
                                            .asInstanceOf(type(DefaultVirtualField.class))
                                            .satisfies(defaultVirtualField -> {
                                                assertThat(defaultVirtualField)
                                                    .extracting(DefaultTypedNamedField::getName)
                                                    .isEqualTo("thisNow3");
                                                assertThat(defaultVirtualField)
                                                    .extracting(DefaultVirtualField::getValueExpression)
                                                    .asInstanceOf(type(DefaultVariableLiteral.class))
                                                    .satisfies(defaultVariableLiteral -> {
                                                        assertThat(defaultVariableLiteral)
                                                            .extracting(DefaultVariableLiteral::getName)
                                                            .isEqualTo("a");
                                                        assertThat(defaultVariableLiteral)
                                                            .extracting(DefaultVariableLiteral::getTypeReference)
                                                            .asInstanceOf(type(DefaultComplexTypeReference.class))
                                                            .extracting(DefaultComplexTypeReference::getName)
                                                            .isEqualTo("A");
                                                        assertThat(defaultVariableLiteral)
                                                            .extracting(DefaultVariableLiteral::getChild)
                                                            .satisfies(childOptional ->
                                                                assertThat(childOptional)
                                                                    .hasValueSatisfying(childVariableLiteral -> {
                                                                        assertThat(childVariableLiteral)
                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                            .extracting(DefaultVariableLiteral::getName)
                                                                            .isEqualTo("b");
                                                                        assertThat(childVariableLiteral)
                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                            .extracting(DefaultVariableLiteral::getTypeReference)
                                                                            .asInstanceOf(type(DefaultComplexTypeReference.class))
                                                                            .extracting(DefaultComplexTypeReference::getName)
                                                                            .isEqualTo("B");
                                                                        assertThat(childVariableLiteral)
                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                            .extracting(DefaultVariableLiteral::getChild)
                                                                            .satisfies(childChildOptional ->
                                                                                assertThat(childChildOptional)
                                                                                    .hasValueSatisfying(childChildVariableLiteral -> {
                                                                                        assertThat(childChildVariableLiteral)
                                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                                            .extracting(DefaultVariableLiteral::getName)
                                                                                            .isEqualTo("c");
                                                                                        assertThat(childChildVariableLiteral)
                                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                                            .extracting(DefaultVariableLiteral::getTypeReference)
                                                                                            .asInstanceOf(type(DefaultComplexTypeReference.class))
                                                                                            .extracting(DefaultComplexTypeReference::getName)
                                                                                            .isEqualTo("C");
                                                                                        assertThat(childChildVariableLiteral)
                                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                                            .extracting(DefaultVariableLiteral::getChild)
                                                                                            .satisfies(childChildChildOptional ->
                                                                                                assertThat(childChildChildOptional)
                                                                                                    .hasValueSatisfying(childChildChildVariableLiteral -> {
                                                                                                        assertThat(childChildChildVariableLiteral)
                                                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                                                            .extracting(DefaultVariableLiteral::getName)
                                                                                                            .isEqualTo("e");
                                                                                                        assertThat(childChildChildVariableLiteral)
                                                                                                            .asInstanceOf(type(DefaultVariableLiteral.class))
                                                                                                            .extracting(DefaultVariableLiteral::getTypeReference)
                                                                                                            .asInstanceOf(type(DefaultComplexTypeReference.class))
                                                                                                            .extracting(DefaultComplexTypeReference::getName)
                                                                                                            .isEqualTo("E");
                                                                                                    })
                                                                                            );
                                                                                    })
                                                                            );
                                                                    })
                                                            );
                                                    });
                                                assertThat(defaultVirtualField)
                                                    .extracting(DefaultTypedField::getType)
                                                    .asInstanceOf(type(AbstractSimpleTypeReference.class))
                                                    .extracting(AbstractSimpleTypeReference::getBaseType, AbstractSimpleTypeReference::getSizeInBits)
                                                    .containsExactly(SimpleTypeReference.SimpleBaseType.BIT, 1);
                                            });
                                    });
                            })
                    );
            });
    }

}