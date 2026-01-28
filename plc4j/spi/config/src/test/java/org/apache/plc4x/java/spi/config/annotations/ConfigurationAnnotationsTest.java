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

package org.apache.plc4x.java.spi.config.annotations;

import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.config.ConfigurationParameterConverter;
import org.apache.plc4x.java.spi.config.annotations.defaults.*;
import org.junit.jupiter.api.Test;

import java.lang.reflect.Field;

import static org.junit.jupiter.api.Assertions.*;

class ConfigurationAnnotationsTest {

    @Test
    void testConfigurationParameterAnnotation() throws NoSuchFieldException {
        Field hostField = TestConfig.class.getDeclaredField("host");
        Field customNameField = TestConfig.class.getDeclaredField("customNameField");

        assertTrue(hostField.isAnnotationPresent(ConfigurationParameter.class));
        assertTrue(customNameField.isAnnotationPresent(ConfigurationParameter.class));

        ConfigurationParameter hostAnnotation = hostField.getAnnotation(ConfigurationParameter.class);
        ConfigurationParameter customAnnotation = customNameField.getAnnotation(ConfigurationParameter.class);

        assertEquals("", hostAnnotation.value()); // default value
        assertEquals("custom-name", customAnnotation.value());
    }

    @Test
    void testRequiredAnnotation() throws NoSuchFieldException {
        Field hostField = TestConfig.class.getDeclaredField("host");
        Field optionalField = TestConfig.class.getDeclaredField("optionalField");

        assertTrue(hostField.isAnnotationPresent(Required.class));
        assertFalse(optionalField.isAnnotationPresent(Required.class));
    }

    @Test
    void testDescriptionAnnotation() throws NoSuchFieldException {
        Field descriptionField = TestConfig.class.getDeclaredField("descriptionField");

        assertTrue(descriptionField.isAnnotationPresent(Description.class));

        Description description = descriptionField.getAnnotation(Description.class);
        assertEquals("This is a test description", description.value());
    }

    @Test
    void testSinceAnnotation() throws NoSuchFieldException {
        Field sinceField = TestConfig.class.getDeclaredField("sinceField");

        assertTrue(sinceField.isAnnotationPresent(Since.class));

        Since since = sinceField.getAnnotation(Since.class);
        assertEquals("1.2.0", since.value());
    }

    @Test
    void testParameterConverterAnnotation() throws NoSuchFieldException {
        Field converterField = TestConfig.class.getDeclaredField("converterField");

        assertTrue(converterField.isAnnotationPresent(ParameterConverter.class));

        ParameterConverter converter = converterField.getAnnotation(ParameterConverter.class);
        assertEquals(TestConverter.class, converter.value());
    }

    @Test
    void testComplexConfigurationParameterAnnotation() throws NoSuchFieldException {
        Field complexField = TestConfig.class.getDeclaredField("complexField");

        assertTrue(complexField.isAnnotationPresent(ComplexConfigurationParameter.class));

        ComplexConfigurationParameter complex = complexField.getAnnotation(ComplexConfigurationParameter.class);
        assertEquals("ssl", complex.prefix());
        assertEquals(0, complex.defaultOverrides().length);
        assertEquals(0, complex.requiredOverrides().length);
    }

    @Test
    void testDefaultValueAnnotations() throws NoSuchFieldException {
        Field stringField = TestConfig.class.getDeclaredField("stringDefaultField");
        Field shortField = TestConfig.class.getDeclaredField("shortDefaultField");
        Field intField = TestConfig.class.getDeclaredField("intDefaultField");
        Field longField = TestConfig.class.getDeclaredField("longDefaultField");
        Field booleanField = TestConfig.class.getDeclaredField("booleanDefaultField");
        Field floatField = TestConfig.class.getDeclaredField("floatDefaultField");
        Field doubleField = TestConfig.class.getDeclaredField("doubleDefaultField");

        // Test StringDefaultValue
        assertTrue(stringField.isAnnotationPresent(StringDefaultValue.class));
        StringDefaultValue stringDefault = stringField.getAnnotation(StringDefaultValue.class);
        assertEquals("default-string", stringDefault.value());

        // Test ShortDefaultValue
        assertTrue(shortField.isAnnotationPresent(ShortDefaultValue.class));
        ShortDefaultValue shortDefault = shortField.getAnnotation(ShortDefaultValue.class);
        assertEquals(23, shortDefault.value());

        // Test IntDefaultValue
        assertTrue(intField.isAnnotationPresent(IntDefaultValue.class));
        IntDefaultValue intDefault = intField.getAnnotation(IntDefaultValue.class);
        assertEquals(42, intDefault.value());

        // Test LongDefaultValue
        assertTrue(longField.isAnnotationPresent(LongDefaultValue.class));
        LongDefaultValue longDefault = longField.getAnnotation(LongDefaultValue.class);
        assertEquals(123456789L, longDefault.value());

        // Test BooleanDefaultValue
        assertTrue(booleanField.isAnnotationPresent(BooleanDefaultValue.class));
        BooleanDefaultValue booleanDefault = booleanField.getAnnotation(BooleanDefaultValue.class);
        assertEquals(true, booleanDefault.value());

        // Test FloatDefaultValue
        assertTrue(floatField.isAnnotationPresent(FloatDefaultValue.class));
        FloatDefaultValue floatDefault = floatField.getAnnotation(FloatDefaultValue.class);
        assertEquals(3.14f, floatDefault.value(), 0.001f);

        // Test DoubleDefaultValue
        assertTrue(doubleField.isAnnotationPresent(DoubleDefaultValue.class));
        DoubleDefaultValue doubleDefault = doubleField.getAnnotation(DoubleDefaultValue.class);
        assertEquals(2.718281828, doubleDefault.value(), 0.000001);
    }

    @Test
    void testComplexConfigurationParameterOverrides() throws NoSuchFieldException {
        Field complexWithOverridesField = TestConfig.class.getDeclaredField("complexWithOverridesField");

        assertTrue(complexWithOverridesField.isAnnotationPresent(ComplexConfigurationParameter.class));

        ComplexConfigurationParameter complex = complexWithOverridesField.getAnnotation(ComplexConfigurationParameter.class);
        assertEquals("tls", complex.prefix());

        // Test default overrides
        assertEquals(1, complex.defaultOverrides().length);
        ComplexConfigurationParameterDefaultOverride defaultOverride = complex.defaultOverrides()[0];
        assertEquals("enabled", defaultOverride.name());
        assertEquals("true", defaultOverride.value());

        // Test required overrides
        assertEquals(1, complex.requiredOverrides().length);
        ComplexConfigurationParameterRequiredOverride requiredOverride = complex.requiredOverrides()[0];
        assertEquals("keystore", requiredOverride.name());
    }

    @Test
    void testAnnotationTargetAndRetention() throws NoSuchFieldException {
        // Test that annotations have correct target and retention
        ConfigurationParameter configParam = TestConfig.class.getDeclaredField("host").getAnnotation(ConfigurationParameter.class);
        assertNotNull(configParam);

        // Test annotation inheritance
        Field[] fields = TestConfig.class.getDeclaredFields();
        long annotatedFields = java.util.Arrays.stream(fields)
            .filter(field -> field.isAnnotationPresent(ConfigurationParameter.class) ||
                            field.isAnnotationPresent(ComplexConfigurationParameter.class))
            .count();

        assertTrue(annotatedFields > 0, "Should have annotated fields");
    }

    // Test configuration class with all annotation types
    public static class TestConfig implements Configuration {
        @ConfigurationParameter
        @Required
        public String host;

        @ConfigurationParameter("custom-name")
        public String customNameField;

        @ConfigurationParameter
        public String optionalField;

        @ConfigurationParameter
        @Description("This is a test description")
        public String descriptionField;

        @ConfigurationParameter
        @Since("1.2.0")
        public String sinceField;

        @ConfigurationParameter
        @ParameterConverter(TestConverter.class)
        public TestValue converterField;

        @ComplexConfigurationParameter(prefix = "ssl", defaultOverrides = {}, requiredOverrides = {})
        public SubConfig complexField;

        @ComplexConfigurationParameter(
            prefix = "tls",
            defaultOverrides = @ComplexConfigurationParameterDefaultOverride(name = "enabled", value = "true"),
            requiredOverrides = @ComplexConfigurationParameterRequiredOverride(name = "keystore")
        )
        public SubConfig complexWithOverridesField;

        @ConfigurationParameter
        @StringDefaultValue("default-string")
        public String stringDefaultField;

        @ConfigurationParameter
        @ShortDefaultValue(23)
        public int shortDefaultField;

        @ConfigurationParameter
        @IntDefaultValue(42)
        public int intDefaultField;

        @ConfigurationParameter
        @LongDefaultValue(123456789L)
        public long longDefaultField;

        @ConfigurationParameter
        @BooleanDefaultValue(true)
        public boolean booleanDefaultField;

        @ConfigurationParameter
        @FloatDefaultValue(3.14f)
        public float floatDefaultField;

        @ConfigurationParameter
        @DoubleDefaultValue(2.718281828)
        public double doubleDefaultField;
    }

    public static class SubConfig implements Configuration {
        @ConfigurationParameter
        public String keystore;

        @ConfigurationParameter
        @BooleanDefaultValue(false)
        public boolean enabled;
    }

    public static class TestValue {
        public final String value;
        public TestValue(String value) { this.value = value; }
    }

    public static class TestConverter implements ConfigurationParameterConverter<TestValue> {
        @Override
        public Class<TestValue> getType() { return TestValue.class; }

        @Override
        public TestValue convert(String value) { return new TestValue(value); }
    }
}