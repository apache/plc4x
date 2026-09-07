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
package org.apache.plc4x.java.s7.tag;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class S7AlarmTagTest {

    @Test
    void matchesAlmAddressIgnoringCase() {
        assertTrue(S7AlarmTag.matches("ALM"));
        assertTrue(S7AlarmTag.matches("alm"));
        assertFalse(S7AlarmTag.matches("ALMS"));
        assertFalse(S7AlarmTag.matches(""));
    }

    @Test
    void tagHandlerRoutesAlmToAlarmTag() {
        S7PlcTagHandler handler = new S7PlcTagHandler();
        assertInstanceOf(S7AlarmTag.class, handler.parseTag("ALM"));
    }

    @Test
    void ofRejectsNonMatchingAddress() {
        assertThrows(IllegalArgumentException.class, () -> S7AlarmTag.of("ALMS"));
    }

    @Test
    void matchesQuerySyntax() {
        assertTrue(S7AlarmTag.matches("QUERY:ALARM_S"));
        assertTrue(S7AlarmTag.matches("query:alarm_s"));    // case-insensitive
        assertTrue(S7AlarmTag.matches("QUERY:ALARM_8"));
        assertFalse(S7AlarmTag.matches("QUERY:ALARM_X"));
    }

    @Test
    void queryAlarmS_carriesQueryTypeAndKind() {
        S7AlarmTag tag = S7AlarmTag.of("QUERY:ALARM_S");
        assertEquals(S7AlarmTag.Kind.QUERY, tag.getKind());
        assertEquals(org.apache.plc4x.java.s7.readwrite.QueryType.ALARM_S, tag.getQueryType());
        assertEquals("QUERY:ALARM_S", tag.getAddressString());
    }

    @Test
    void queryAlarm8_carriesQueryTypeAndKind() {
        S7AlarmTag tag = S7AlarmTag.of("QUERY:ALARM_8");
        assertEquals(S7AlarmTag.Kind.QUERY, tag.getKind());
        assertEquals(org.apache.plc4x.java.s7.readwrite.QueryType.ALARM_8, tag.getQueryType());
    }

    @Test
    void pushTag_hasNullQueryType() {
        S7AlarmTag tag = S7AlarmTag.of("ALM");
        assertEquals(S7AlarmTag.Kind.PUSH, tag.getKind());
        assertNull(tag.getQueryType());
    }
}
