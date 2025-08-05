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

package org.apache.plc4x.java.spi.configuration;

import org.apache.plc4x.java.spi.configuration.config.*;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

public class ConfigurationFactoryTest {

    @Test
    public void simpleConfigurationWithEmptyParamStringTest() {
        TestConfigurationSimple configuration = new ConfigurationFactory().createConfiguration(
            TestConfigurationSimple.class,
            "protocolCode", "transportCode", "transportConfig", "");
        assertFalse(configuration.isBooleanField());
        assertEquals(configuration.getIntegerField(), 0);
        assertEquals(configuration.getLongField(), 0L);
        assertEquals(configuration.getFloatField(), 0.0f);
        assertEquals(configuration.getDoubleField(), 0.0d);
        assertNull(configuration.getStringField());
    }

    @Test
    public void simpleConfigurationTest() {
        TestConfigurationSimple configuration = new ConfigurationFactory().createConfiguration(
            TestConfigurationSimple.class,
            "protocolCode", "transportCode", "transportConfig",
            "booleanField=true&integerField=42&longField=232323232323232323&floatField=3.1415927&doubleField=2.718281828459045&stringField=Hurz");
        assertTrue(configuration.isBooleanField());
        assertEquals(configuration.getIntegerField(), 42);
        assertEquals(configuration.getLongField(), 232323232323232323L);
        assertEquals(configuration.getFloatField(), 3.1415927f);
        assertEquals(configuration.getDoubleField(), 2.718281828459045d);
        assertEquals(configuration.getStringField(), "Hurz");
    }

    @Test
    public void simpleConfigurationWithDefaultsTest() {
        TestConfigurationDefaults configuration = new ConfigurationFactory().createConfiguration(
            TestConfigurationDefaults.class,
            "protocolCode", "transportCode", "transportConfig",
            "booleanField=true&integerField=42&longField=232323232323232323&floatField=3.1415927&doubleField=2.718281828459045&stringField=Hurz");
        assertTrue(configuration.isBooleanField());
        assertEquals(configuration.getIntegerField(), 42);
        assertEquals(configuration.getLongField(), 232323232323232323L);
        assertEquals(configuration.getFloatField(), 3.1415927f);
        assertEquals(configuration.getDoubleField(), 2.718281828459045d);
        assertEquals(configuration.getStringField(), "Hurz");
    }

    @Test
    public void requiredConfigurationTest() {
        // Test with missing of required fields.
        try {
            new ConfigurationFactory().createConfiguration(
                TestConfigurationRequired.class,
                "protocolCode", "transportCode", "transportConfig",
                "");
            fail("Parsing should have failed.");
        } catch (IllegalArgumentException e) {
            assertEquals(e.getMessage(), "Missing required fields: [floatField, longField, doubleField, booleanField, integerField, stringField]");
        }

        // Test with all fields provided.
        TestConfigurationRequired configuration = new ConfigurationFactory().createConfiguration(
            TestConfigurationRequired.class,
            "protocolCode", "transportCode", "transportConfig",
            "booleanField=true&integerField=42&longField=232323232323232323&floatField=3.1415927&doubleField=2.718281828459045&stringField=Hurz");
        assertTrue(configuration.isBooleanField());
        assertEquals(configuration.getIntegerField(), 42);
        assertEquals(configuration.getLongField(), 232323232323232323L);
        assertEquals(configuration.getFloatField(), 3.1415927f);
        assertEquals(configuration.getDoubleField(), 2.718281828459045d);
        assertEquals(configuration.getStringField(), "Hurz");
    }

    @Test
    public void parameterConverterTest() {
        TestConfigurationParameterConverter configuration = new ConfigurationFactory().createConfiguration(
            TestConfigurationParameterConverter.class,
            "protocolCode", "transportCode", "transportConfig",
            "parameterConverterTypeField=lalala");
        assertNotNull(configuration.getParameterConverterTypeField());
        assertEquals(configuration.getParameterConverterTypeField().getConvertedType(), "lalala");
    }

    @Test
    public void complexParameterTest() {
        TestConfigurationComplex configuration = new ConfigurationFactory().createConfiguration(
            TestConfigurationComplex.class,
            "protocolCode", "transportCode", "transportConfig",
            "simple.booleanField=true&simple.integerField=42&simple.longField=232323232323232323&simple.floatField=3.1415927&simple.doubleField=2.718281828459045&simple.stringField=Hurz");
        assertNotNull(configuration.getComplexSimple());
        TestConfigurationSimple complexSimple = configuration.getComplexSimple();
        assertTrue(complexSimple.isBooleanField());
        assertEquals(complexSimple.getIntegerField(), 42);
        assertEquals(complexSimple.getLongField(), 232323232323232323L);
        assertEquals(complexSimple.getFloatField(), 3.1415927f);
        assertEquals(complexSimple.getDoubleField(), 2.718281828459045d);
        assertEquals(complexSimple.getStringField(), "Hurz");
    }

}
