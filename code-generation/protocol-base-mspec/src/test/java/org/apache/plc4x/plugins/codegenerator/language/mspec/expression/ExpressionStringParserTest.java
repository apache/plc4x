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
package org.apache.plc4x.plugins.codegenerator.language.mspec.expression;

import org.apache.commons.io.IOUtils;
import org.apache.plc4x.plugins.codegenerator.types.terms.*;
import org.junit.jupiter.api.Test;

import java.nio.charset.Charset;
import java.util.List;
import java.util.Optional;
import java.util.function.Consumer;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.*;

class ExpressionStringParserTest {

    ExpressionStringParser SUT = new ExpressionStringParser((typeRefName, setTypeDefinition) -> {}, "hurz");

    @Test
    void parseSimple() {
        Term term = SUT.parse(IOUtils.toInputStream("payload.lengthInBytes+4", Charset.defaultCharset()));
        assertThat(term, not(nullValue()));
        assertThat(term, instanceOf(BinaryTerm.class));
        BinaryTerm binaryTerm = (BinaryTerm) term;
        assertVariableLiteral(
            binaryTerm.getA(),
            "payload",
            null,
            lengthInBytesVariableLiteral -> assertVariableLiteral(
                lengthInBytesVariableLiteral,
                "lengthInBytes"
            )
        );
        assertNumericLiteral(binaryTerm.getB(), 4L);
        assertThat(binaryTerm.getOperation(), is("+"));
    }


    @Test
    void parseDoubleCast() {
        Term term = SUT.parse(IOUtils.toInputStream("CAST(CAST(parameter,S7ParameterUserData).items(hurz)[0],S7ParameterUserDataItemCPUFunctions).cpuFunctionType", Charset.defaultCharset()));
        assertVariableLiteral(
            term,
            "CAST",
            outerCast -> {
                assertThat(outerCast, hasSize(2));
                assertVariableLiteral(
                    outerCast.get(0),
                    "CAST",
                    innerCast -> {
                        assertThat(innerCast, hasSize(2));
                        assertVariableLiteral(
                            innerCast.get(0),
                            "parameter"
                        );
                        assertVariableLiteral(
                            innerCast.get(1),
                            "S7ParameterUserData"
                        );
                    },
                    items -> {
                        assertVariableLiteral(
                            items,
                            "items",
                            hurzes -> {
                                assertThat(hurzes, hasSize(1));
                                assertVariableLiteral(
                                    hurzes.get(0),
                                    "hurz"
                                );
                            },
                            null,
                            0
                        );
                    });
                assertVariableLiteral(
                    outerCast.get(1),
                    "S7ParameterUserDataItemCPUFunctions"
                );
            },
            variableLiteral -> {
                assertVariableLiteral(
                    variableLiteral,
                    "cpuFunctionType"
                );
            }
        );
    }

    @Test
    void parseCast() {
        Term term = SUT.parse(IOUtils.toInputStream("CAST(parameter,S7ParameterUserData).items(hurz)[0]", Charset.defaultCharset()));
        assertVariableLiteral(
            term,
            "CAST",
            terms -> {
                assertThat(terms, hasSize(2));
                assertVariableLiteral(
                    terms.get(0),
                    "parameter"
                );
                assertVariableLiteral(
                    terms.get(1),
                    "S7ParameterUserData"
                );
            },
            variableLiteral -> {
                assertVariableLiteral(
                    variableLiteral,
                    "items",
                    terms -> {
                        assertThat(terms, hasSize(1));
                        assertVariableLiteral(
                            terms.get(0),
                            "hurz"
                        );
                    },
                    null,
                    0
                );
            }
        );
    }

    void assertNumericLiteral(Term term, Number number) {
        assertThat(term, not(nullValue()));
        assertThat(term, instanceOf(NumericLiteral.class));
        NumericLiteral numericLiteral = (NumericLiteral) term;
        assertThat(numericLiteral.getNumber(), is(number));
    }

    void assertVariableLiteral(Term term, String name) {
        assertVariableLiteral(term, name, null, null, null);
    }

    void assertVariableLiteral(Term term, String name, Consumer<List<Term>> argsAsserter, Consumer<VariableLiteral> childAsserter) {
        assertVariableLiteral(term, name, argsAsserter, childAsserter, null);
    }

    void assertVariableLiteral(Term term, String name, Consumer<List<Term>> argsAsserter, Consumer<VariableLiteral> childAsserter, Integer index) {
        assertThat(term, not(nullValue()));
        assertThat(term, instanceOf(VariableLiteral.class));
        VariableLiteral variableLiteral = (VariableLiteral) term;
        assertThat(variableLiteral.getName(), is(name));
        assertThat(variableLiteral.getIndex().orElse(null), is(index));
        if (argsAsserter != null) {
            argsAsserter.accept(variableLiteral.getArgs().orElseThrow(RuntimeException::new));
        } else {
            assertThat(variableLiteral.getArgs().orElse(null), nullValue());
        }
        if (childAsserter != null) {
            childAsserter.accept(variableLiteral.getChild().orElseThrow(RuntimeException::new));
        } else {
            assertThat(variableLiteral.getChild().orElse(null), nullValue());
        }
    }


}