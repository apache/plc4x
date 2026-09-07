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

package org.apache.plc4x.java.spi.config;

import org.apache.plc4x.java.spi.config.annotations.ComplexConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.ParameterConverter;
import org.apache.plc4x.java.spi.config.annotations.Required;
import org.apache.plc4x.java.spi.config.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.spi.config.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.spi.config.annotations.defaults.StringDefaultValue;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.io.File;

import static org.junit.jupiter.api.Assertions.*;

class ConfigurationFactoryTest {

    private ConfigurationFactory factory;

    @BeforeEach
    void setUp() {
        factory = new ConfigurationFactory();
    }

    @Test
    void testCreateConfiguration_withBasicParameters() {
        String paramString = "host=192.168.1.100&port=502&timeout=5000";

        TestConfiguration config = factory.createConfiguration(TestConfiguration.class, paramString);

        assertEquals("192.168.1.100", config.host);
        assertEquals(502, config.port);
        assertEquals(5000, config.timeout);
        assertTrue(config.enabled); // default value
    }

    @Test
    void testCreateConfiguration_withDefaultValues() {
        String paramString = "host=192.168.1.100";

        TestConfiguration config = factory.createConfiguration(TestConfiguration.class, paramString);

        assertEquals("192.168.1.100", config.host);
        assertEquals(502, config.port); // default value
        assertEquals(1000, config.timeout); // default value
        assertTrue(config.enabled); // default value
        assertEquals("default", config.protocol); // default value
    }

    @Test
    void testCreateConfiguration_withMissingRequiredFields() {
        String paramString = "port=502";

        IllegalArgumentException exception = assertThrows(IllegalArgumentException.class,
            () -> factory.createConfiguration(TestConfiguration.class, paramString));

        assertTrue(exception.getMessage().contains("Missing required fields"));
        assertTrue(exception.getMessage().contains("host"));
    }

    @Test
    void testCreateConfiguration_withCustomParameterNames() {
        String paramString = "custom-host=test.com&port=8080";

        TestConfigurationCustomNames config = factory.createConfiguration(TestConfigurationCustomNames.class, paramString);

        assertEquals("test.com", config.hostname);
        assertEquals(8080, config.port);
    }

    @Test
    void testCreateConfiguration_withEnumParameters() {
        String paramString = "mode=CLIENT";

        TestConfigurationWithEnum config = factory.createConfiguration(TestConfigurationWithEnum.class, paramString);

        assertEquals(TestMode.CLIENT, config.mode);
        assertEquals(TestMode.SERVER, config.defaultMode); // default enum value
    }

    @Test
    void testCreateConfiguration_withFileParameter() {
        String paramString = "configFile=/path/to/config.xml";

        TestConfigurationWithFile config = factory.createConfiguration(TestConfigurationWithFile.class, paramString);

        assertEquals(new File("/path/to/config.xml"), config.configFile);
    }

    @Test
    void testCreateConfiguration_withCustomConverter() {
        String paramString = "customValue=test_value";

        TestConfigurationWithConverter config = factory.createConfiguration(TestConfigurationWithConverter.class, paramString);

        assertEquals("CONVERTED_test_value", config.customValue.value);
    }

    @Test
    void testCreateConfiguration_withComplexParameter() {
        String paramString = "ssl.keystore=/path/to/keystore&ssl.password=secret&ssl.enabled=true";

        TestConfigurationWithComplex config = factory.createConfiguration(TestConfigurationWithComplex.class, paramString);

        assertNotNull(config.sslConfig);
        assertEquals("/path/to/keystore", config.sslConfig.keystore);
        assertEquals("secret", config.sslConfig.password);
        assertTrue(config.sslConfig.enabled);
    }

    @Test
    void testCreatePrefixedConfiguration() {
        String paramString = "prefix.host=test.com&prefix.port=8080&other.param=value";

        TestConfiguration config = factory.createPrefixedConfiguration(TestConfiguration.class, "prefix", paramString);

        assertEquals("test.com", config.host);
        assertEquals(8080, config.port);
    }

    @Test
    void testCreateConfiguration_withEmptyParameterString() {
        // This test expects a required field exception since host is required
        IllegalArgumentException exception = assertThrows(IllegalArgumentException.class,
            () -> factory.createConfiguration(TestConfiguration.class, ""));

        assertTrue(exception.getMessage().contains("Missing required fields"));
        assertTrue(exception.getMessage().contains("host"));
    }

    @Test
    void testCreateConfiguration_withUrlEncodedParameters() {
        String paramString = "host=test%2Ecom&description=Hello%20World";

        TestConfigurationWithDescription config = factory.createConfiguration(TestConfigurationWithDescription.class, paramString);

        assertEquals("test.com", config.host);
        assertEquals("Hello World", config.description);
    }

    @Test
    void testCreateConfiguration_withBooleanParameters() {
        String paramString = "enabled=false&debug=true";

        TestConfigurationWithBooleans config = factory.createConfiguration(TestConfigurationWithBooleans.class, paramString);

        assertFalse(config.enabled);
        assertTrue(config.debug);
        assertTrue(config.defaultBool); // default value
    }

    @Test
    void testCreateConfiguration_withNumericTypes() {
        String paramString = "byteVal=127&shortVal=32767&longVal=9223372036854775807&floatVal=3.14&doubleVal=2.718281828";

        TestConfigurationWithNumericTypes config = factory.createConfiguration(TestConfigurationWithNumericTypes.class, paramString);

        assertEquals((byte) 127, config.byteVal);
        assertEquals((short) 32767, config.shortVal);
        assertEquals(9223372036854775807L, config.longVal);
        assertEquals(3.14f, config.floatVal, 0.001f);
        assertEquals(2.718281828, config.doubleVal, 0.000001);
    }

    @Test
    void testGetConfigurationName_withDefaultFieldName() throws NoSuchFieldException {
        java.lang.reflect.Field field = TestConfiguration.class.getDeclaredField("timeout");

        String name = ConfigurationFactory.getConfigurationName(field);

        assertEquals("timeout", name);
    }

    @Test
    void testGetConfigurationName_withCustomName() throws NoSuchFieldException {
        java.lang.reflect.Field field = TestConfigurationCustomNames.class.getDeclaredField("hostname");

        String name = ConfigurationFactory.getConfigurationName(field);

        assertEquals("custom-host", name);
    }

    @Test
    void testConfigureObject_withHasConfiguration() {
        TestConfiguration config = new TestConfiguration();
        config.host = "test.com";

        TestConfigurableObject obj = new TestConfigurableObject();
        TestConfigurableObject result = ConfigurationFactory.configure(config, obj);

        assertSame(obj, result);
        assertEquals(config, obj.getConfiguration());
    }

    // Test Configuration Classes

    public static class TestConfiguration implements Configuration {
        @ConfigurationParameter
        @Required
        public String host;

        @ConfigurationParameter
        @IntDefaultValue(502)
        public int port;

        @ConfigurationParameter
        @IntDefaultValue(1000)
        public int timeout;

        @ConfigurationParameter
        @BooleanDefaultValue(true)
        public boolean enabled;

        @ConfigurationParameter
        @StringDefaultValue("default")
        public String protocol;
    }

    public static class TestConfigurationCustomNames implements Configuration {
        @ConfigurationParameter("custom-host")
        @Required
        public String hostname;

        @ConfigurationParameter
        public int port;
    }

    public static class TestConfigurationWithEnum implements Configuration {
        @ConfigurationParameter
        public TestMode mode;

        @ConfigurationParameter
        @StringDefaultValue("SERVER")
        public TestMode defaultMode;
    }

    public static class TestConfigurationWithFile implements Configuration {
        @ConfigurationParameter
        public File configFile;
    }

    public static class TestConfigurationWithConverter implements Configuration {
        @ConfigurationParameter
        @ParameterConverter(TestConverter.class)
        public TestValue customValue;
    }

    public static class TestConfigurationWithComplex implements Configuration {
        @ComplexConfigurationParameter(prefix = "ssl", defaultOverrides = {}, requiredOverrides = {})
        public SslConfiguration sslConfig;
    }

    public static class TestConfigurationWithDescription implements Configuration {
        @ConfigurationParameter
        @Required
        public String host;

        @ConfigurationParameter
        public String description;
    }

    public static class TestConfigurationWithBooleans implements Configuration {
        @ConfigurationParameter
        public boolean enabled;

        @ConfigurationParameter
        public boolean debug;

        @ConfigurationParameter
        @BooleanDefaultValue(true)
        public boolean defaultBool;
    }

    public static class TestConfigurationWithNumericTypes implements Configuration {
        @ConfigurationParameter
        public byte byteVal;

        @ConfigurationParameter
        public short shortVal;

        @ConfigurationParameter
        public long longVal;

        @ConfigurationParameter
        public float floatVal;

        @ConfigurationParameter
        public double doubleVal;
    }

    public static class SslConfiguration implements Configuration {
        @ConfigurationParameter
        public String keystore;

        @ConfigurationParameter
        public String password;

        @ConfigurationParameter
        @BooleanDefaultValue(false)
        public boolean enabled;
    }

    public enum TestMode {
        CLIENT, SERVER
    }

    public static class TestValue {
        public final String value;

        public TestValue(String value) {
            this.value = value;
        }
    }

    public static class TestConverter implements ConfigurationParameterConverter<TestValue> {
        @Override
        public Class<TestValue> getType() {
            return TestValue.class;
        }

        @Override
        public TestValue convert(String value) {
            return new TestValue("CONVERTED_" + value);
        }
    }

    public static class TestConfigurableObject implements HasConfiguration<TestConfiguration> {
        private TestConfiguration configuration;

        public TestConfiguration getConfiguration() {
            return configuration;
        }

        @Override
        public void setConfiguration(TestConfiguration configuration) {
            this.configuration = configuration;
        }
    }
}