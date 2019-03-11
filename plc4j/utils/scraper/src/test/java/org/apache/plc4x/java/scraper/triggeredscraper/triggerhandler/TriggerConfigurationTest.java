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

package org.apache.plc4x.java.scraper.triggeredscraper.triggerhandler;

import org.apache.plc4x.java.scraper.exception.ScraperException;
import org.apache.plc4x.java.scraper.triggeredscraper.TriggeredScrapeJobImpl;
import org.apache.plc4x.test.FastTests;
import org.junit.experimental.categories.Category;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;
import org.mockito.Mockito;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.stream.Stream;

import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.IsNull.notNullValue;
import static org.junit.Assert.assertThat;

/**
 * testing valid and invalid triggerConfigStrings
 */
class TriggerConfigurationTest {
    private static final Logger LOGGER = LoggerFactory.getLogger(TriggerConfigurationTest.class);

    private static Stream<Arguments> validTriggerPattern() {
        return Stream.of(
            Arguments.of("(S7_TRIGGER_VAR,50,(%I0.1:BOOL)==(true))",TriggerConfiguration.TriggerType.S7_TRIGGER_VAR, 50, TriggerConfiguration.Comparators.EQUAL, true),
            Arguments.of("(S7_TRIGGER_VAR,50,(%I0.1:BOOL)!=(0))",TriggerConfiguration.TriggerType.S7_TRIGGER_VAR, 50, TriggerConfiguration.Comparators.UNEQUAL, false),
            Arguments.of("(S7_TRIGGER_VAR,50,(%DB111:DBW10:INT)<=(33))",TriggerConfiguration.TriggerType.S7_TRIGGER_VAR, 50, TriggerConfiguration.Comparators.SMALLER_EQUAL, 33.0),
            Arguments.of("(S7_TRIGGER_VAR,50,(%DB111:DBB10:USINT)>=(33))",TriggerConfiguration.TriggerType.S7_TRIGGER_VAR, 50, TriggerConfiguration.Comparators.GREATER_EQUAL, 33.0),
            Arguments.of("(S7_TRIGGER_VAR,50,(%DB111:DBD10:DINT)<(33))",TriggerConfiguration.TriggerType.S7_TRIGGER_VAR, 50, TriggerConfiguration.Comparators.SMALLER, 33.0),
            Arguments.of("(S7_TRIGGER_VAR,50,(%DB111:DBD10:REAL)>(33.3))",TriggerConfiguration.TriggerType.S7_TRIGGER_VAR, 50, TriggerConfiguration.Comparators.GREATER, 33.3),
            Arguments.of("(S7_TRIGGER_VAR,50,(%DB111:DBD10:REAL)>(33.3))",TriggerConfiguration.TriggerType.S7_TRIGGER_VAR, 50, TriggerConfiguration.Comparators.GREATER, 33.3),
            Arguments.of("(S7_TRIGGER_VAR,50,(%DB111:DBD10:REAL)>(-1))",TriggerConfiguration.TriggerType.S7_TRIGGER_VAR, 50, TriggerConfiguration.Comparators.GREATER, -1.0),
            Arguments.of("(SCHEDULED,1000)",TriggerConfiguration.TriggerType.SCHEDULED, 1000, null, null)
        );
    }

    private static Stream<Arguments> invalidTriggerPattern() {
        return Stream.of(
            Arguments.of("(S7_TRIGGER_VAR,50,(%I0.1:BOOL)(==)(true))"),
            Arguments.of("(SCHEDULED,50,(%I0.1:BOOL)==(true))"),
            Arguments.of("(S7_TRIGGER_VAR,50)"),
            Arguments.of("(S7_TRIGGER_VAR)"),
            Arguments.of("(S7_TRIGGER_VAR,50,(%I0.1:BOOL)==(0.1))"),
            Arguments.of("(S7_TRIGGER_VAR,50,(%DB111:DBW10:BOOL)==(33))"),
            Arguments.of("(S7_TRIGGER_VAR,50,(%DB111:DBX10:BOOL)==(33))"),
            Arguments.of("(S7_TRIGGER_VAR,50,(%DB111:DBX10.1:BOOL)==(33))"),
            Arguments.of("(S7_TRIGGER_VAR,50,(%DB111:DBX10.1:BOOL)<=(true))"),
            Arguments.of("(S7_TRIGGER_VAR,50,(%DB111:DBW10:INT)<=(true))"),
            Arguments.of("(MODBUS_TRIGGER_VAR,50)"),
            Arguments.of("(MODBUS_TRIGGER_VAR,50,(%DB111:DBW10:INT)<=(11))")
        );
    }

    @ParameterizedTest
    @Category(FastTests.class)
    @MethodSource("validTriggerPattern")
    void testValidFieldQueryParsing(String triggerConfig,
                                    TriggerConfiguration.TriggerType triggerType,
                                    long scrapeInterval,
                                    TriggerConfiguration.Comparators comparator,
                                    Object refValue) {
        TriggeredScrapeJobImpl triggeredScrapeJob = Mockito.mock(TriggeredScrapeJobImpl.class);
        TriggerConfiguration triggerConfiguration = null;
        try {
            triggerConfiguration = TriggerConfiguration.createConfiguration(triggerConfig,triggeredScrapeJob);
        } catch (ScraperException e) {
            //should not happen
        }

        assertThat(triggerConfiguration, notNullValue());
        assertThat(triggerConfiguration.getScrapeInterval(), equalTo(scrapeInterval));
        assertThat(triggerConfiguration.getTriggerType(), equalTo(triggerType));
        assertThat(triggerConfiguration.getComparatorType(), equalTo(comparator));
        assertThat(triggerConfiguration.getCompareValue(), equalTo(refValue));
    }


    @ParameterizedTest
    @Category(FastTests.class)
    @MethodSource("invalidTriggerPattern")
    void testInvalidFieldQueryParsing(String triggerConfig) {
        TriggeredScrapeJobImpl triggeredScrapeJob = Mockito.mock(TriggeredScrapeJobImpl.class);
        TriggerConfiguration triggerConfiguration = null;
        try {
            triggerConfiguration = TriggerConfiguration.createConfiguration(triggerConfig,triggeredScrapeJob);
            assertThat(triggerConfiguration,null);
            //NPE should happen when test fails!
        } catch (ScraperException e) {
            LOGGER.info("Exception as expected for positive test result: {}",e.getMessage());
            //should happen
        }


    }

}