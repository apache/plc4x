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
package org.apache.plc4x.java.knxnetip;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.ValueSource;

import java.util.regex.Pattern;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * A browse query names group addresses with {@code *} standing in for a level, so it is a pattern
 * in its own small language rather than a regular expression. Anything the caller sends that is not
 * in that language has to be turned away, not handed to the regex engine.
 */
class BrowseQueryGrammarTest {

    @ParameterizedTest
    @ValueSource(strings = {"*", "**", "*/*/*", "1/*", "1/2/*", "1/2/3", "*/2/3", "1/*/3", "31/7/255"})
    void theDocumentedQueriesAreAccepted(String query) {
        assertTrue(KnxNetIpConnection.isSupportedBrowseQuery(query),
            "'" + query + "' is in the documented grammar and must be accepted");
    }

    @Test
    void anAbsentQueryMeansEverything() {
        assertTrue(KnxNetIpConnection.isSupportedBrowseQuery(null));
        assertTrue(KnxNetIpConnection.isSupportedBrowseQuery(""));
        assertNull(KnxNetIpConnection.compileQueryPattern("*"));
        assertNull(KnxNetIpConnection.compileQueryPattern("**"));
        assertNull(KnxNetIpConnection.compileQueryPattern("*/*/*"));
    }

    @ParameterizedTest
    @ValueSource(strings = {"(", ")", "1/2/[", "*)", "a/b/c", "1/2/3/4", "1//3", "-1/2/3", "1/2/3 "})
    void aQueryOutsideTheGrammarIsTurnedAway(String query) {
        assertFalse(KnxNetIpConnection.isSupportedBrowseQuery(query),
            "'" + query + "' is not in the documented grammar and must be turned away");
    }

    /**
     * The query used to be pasted into a regex, so a caller could widen their own query past the
     * addresses they asked for. The alternation is what does it - and note it needs no {@code *},
     * because {@code *} was being substituted for {@code \d+} and that defused it by accident.
     */
    @Test
    void aQueryCannotWidenItselfPastTheAddressesItNames() {
        for (String query : new String[]{"1/2/3|.+", "1/2/3|[0-9/]+", "1/2/3|\\d/\\d/\\d", "1/2/3|.*"}) {
            assertFalse(KnxNetIpConnection.isSupportedBrowseQuery(query),
                "'" + query + "' must be turned away rather than matching addresses it does not name");
        }
    }

    /**
     * Whatever the caller sends, working out what it matches must not raise anything the browse
     * path was not written to expect.
     */
    @ParameterizedTest
    @ValueSource(strings = {"(", "1/2/[", "1/2/3|.+", "\\", "[a-", "1/2/3", "*"})
    void decidingOnAQueryNeverThrows(String query) {
        assertDoesNotThrow(() -> {
            if (KnxNetIpConnection.isSupportedBrowseQuery(query)) {
                Pattern pattern = KnxNetIpConnection.compileQueryPattern(query);
                if (pattern != null) {
                    pattern.matcher("1/2/3").matches();
                }
            }
        });
    }

    @Test
    void anAcceptedQueryStillMatchesWhatItUsedTo() {
        assertTrue(KnxNetIpConnection.compileQueryPattern("1/*").matcher("1/2/3").matches());
        assertFalse(KnxNetIpConnection.compileQueryPattern("1/*").matcher("2/2/3").matches());
        assertTrue(KnxNetIpConnection.compileQueryPattern("1/2/*").matcher("1/2/9").matches());
        assertFalse(KnxNetIpConnection.compileQueryPattern("1/2/*").matcher("1/3/9").matches());
        assertTrue(KnxNetIpConnection.compileQueryPattern("1/2/3").matcher("1/2/3").matches());
        assertFalse(KnxNetIpConnection.compileQueryPattern("1/2/3").matcher("1/2/4").matches());
        assertTrue(KnxNetIpConnection.compileQueryPattern("*/2/3").matcher("9/2/3").matches());
        assertFalse(KnxNetIpConnection.compileQueryPattern("*/2/3").matcher("9/2/4").matches());
    }
}
