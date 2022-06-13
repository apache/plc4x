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
package org.apache.plc4x.java.opm;

import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

public class OpmUtilsTest implements WithAssertions {

    @Test
    public void expression_matches() {
        assertTrue(OpmUtils.pattern.matcher("${Hallo}").matches());
        assertTrue(OpmUtils.pattern.matcher("${Hallo:Hallo}").matches());
        // ...
        assertTrue(OpmUtils.pattern.matcher("${Ha{}llo}").matches());
    }

    @Test
    public void getAlias_matches() {
        String alias = OpmUtils.getAlias("${hallo}");

        assertEquals("hallo", alias);
    }

    @Test
    public void isAlias_bothCases() {
        // True
        assertTrue(OpmUtils.isAlias("${hallo}"));
        assertTrue(OpmUtils.isAlias("${hal{}lo}"));
        assertTrue(OpmUtils.isAlias("${hallo:hallo}"));
        // False
        assertFalse(OpmUtils.isAlias("hallo"));
        assertFalse(OpmUtils.isAlias("${hallo"));
        assertFalse(OpmUtils.isAlias("${ha}llo"));
    }

    @Test
    public void isValidExpression_startingDollar_false() {
        assertFalse(OpmUtils.isValidExpression("${hallo"));
    }

    @Test
    public void getAlias_illegalString_throws() {
        assertThatThrownBy(() -> OpmUtils.getAlias("hallo"))
            .isInstanceOf(IllegalArgumentException.class);
    }
}