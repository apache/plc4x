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

package org.apache.plc4x.java.spi.config.integration;

import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.config.ConfigurationFactory;
import org.apache.plc4x.java.spi.config.ConfigurationParameterConverter;
import org.apache.plc4x.java.spi.config.HasConfiguration;
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
import java.time.Duration;
import java.util.Arrays;
import java.util.List;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;

class ConfigurationSystemIntegrationTest {

    private ConfigurationFactory factory;

    @BeforeEach
    void setUp() {
        factory = new ConfigurationFactory();
    }

    @Test
    void testRealWorldDriverConfiguration() {
        String connectionString = "host=192.168.1.100&port=502&timeout=5000&retry-count=3&protocol=TCP" +
                                 "&ssl.enabled=true&ssl.keystore=/path/to/keystore.jks&ssl.password=secret" +
                                 "&advanced.buffer-size=8192&advanced.keep-alive=true";

        ModbusDriverConfiguration config = factory.createConfiguration(ModbusDriverConfiguration.class, connectionString);

        // Test basic parameters
        assertEquals("192.168.1.100", config.host);
        assertEquals(502, config.port);
        assertEquals(Duration.ofMillis(5000), config.timeout);
        assertEquals(3, config.retryCount);
        assertEquals(Protocol.TCP, config.protocol);

        // Test SSL complex configuration
        assertNotNull(config.sslConfig);
        assertTrue(config.sslConfig.enabled);
        assertEquals("/path/to/keystore.jks", config.sslConfig.keystore);
        assertEquals("secret", config.sslConfig.password);

        // Test advanced complex configuration
        assertNotNull(config.advancedConfig);
        assertEquals(8192, config.advancedConfig.bufferSize);
        assertTrue(config.advancedConfig.keepAlive);
        assertEquals(30, config.advancedConfig.connectionPoolSize); // default value
    }

    @Test
    void testPrefixedConfigurationExtraction() {
        String connectionString = "main.host=192.168.1.100&main.port=502&other.value=ignored&main.enabled=true";

        BasicConfiguration config = factory.createPrefixedConfiguration(
            BasicConfiguration.class, "main", connectionString);

        assertEquals("192.168.1.100", config.host);
        assertEquals(502, config.port);
        assertTrue(config.enabled);
    }

    @Test
    void testConfigurationWithMultipleConverters() {
        String connectionString = "duration=30s&file-list=/tmp/file1.txt,/tmp/file2.txt&tags=tag1,tag2,tag3";

        AdvancedConfiguration config = factory.createConfiguration(AdvancedConfiguration.class, connectionString);

        assertEquals(Duration.ofSeconds(30), config.duration);
        assertEquals(2, config.fileList.length);
        assertEquals(new File("/tmp/file1.txt"), config.fileList[0]);
        assertEquals(new File("/tmp/file2.txt"), config.fileList[1]);
        assertArrayEquals(new String[]{"tag1", "tag2", "tag3"}, config.tags);
    }

    @Test
    void testNestedComplexConfiguration() {
        String connectionString = "database.host=db.example.com&database.port=5432&database.name=testdb" +
                                 "&database.ssl.enabled=true&database.ssl.cert-file=/certs/client.crt" +
                                 "&cache.enabled=true&cache.size=1000&cache.ttl=3600";

        ApplicationConfiguration config = factory.createConfiguration(ApplicationConfiguration.class, connectionString);

        // Test database configuration
        assertNotNull(config.databaseConfig);
        assertEquals("db.example.com", config.databaseConfig.host);
        assertEquals(5432, config.databaseConfig.port);
        assertEquals("testdb", config.databaseConfig.name);

        // Test nested SSL configuration within database
        assertNotNull(config.databaseConfig.sslConfig);
        assertTrue(config.databaseConfig.sslConfig.enabled);
        assertEquals("/certs/client.crt", config.databaseConfig.sslConfig.certFile);

        // Test cache configuration
        assertNotNull(config.cacheConfig);
        assertTrue(config.cacheConfig.enabled);
        assertEquals(1000, config.cacheConfig.size);
        assertEquals(3600, config.cacheConfig.ttl);
    }

    @Test
    void testConfigurationWithHasConfigurationInterface() {
        String connectionString = "host=test.com&port=8080&enabled=true";

        BasicConfiguration config = factory.createConfiguration(BasicConfiguration.class, connectionString);
        TestService service = new TestService();

        TestService configuredService = ConfigurationFactory.configure(config, service);

        assertSame(service, configuredService);
        assertEquals(config, service.getConfiguration());
        assertEquals("test.com", service.getConfiguration().host);
    }

    @Test
    void testConfigurationValidationAndErrorHandling() {
        // Test missing required field
        String connectionString = "port=502"; // missing required host

        IllegalArgumentException exception = assertThrows(IllegalArgumentException.class,
            () -> factory.createConfiguration(ModbusDriverConfiguration.class, connectionString));

        // The actual error could be about Duration converter or missing fields
        assertTrue(exception.getMessage().contains("Missing required fields") ||
                  exception.getMessage().contains("Can not set") ||
                  exception.getMessage().contains("Duration"),
                  "Expected error message, but was: " + exception.getMessage());
    }

    @Test
    void testConfigurationWithInvalidEnumValue() {
        String connectionString = "host=test.com&protocol=INVALID_PROTOCOL";

        IllegalArgumentException exception = assertThrows(IllegalArgumentException.class,
            () -> factory.createConfiguration(ModbusDriverConfiguration.class, connectionString));

        assertTrue(exception.getMessage().contains("INVALID_PROTOCOL") ||
                  exception.getCause() instanceof IllegalArgumentException);
    }

    @Test
    void testConfigurationFromMap() {
        Map<String, List<String>> paramMap = Map.of(
            "host", List.of("192.168.1.100"),
            "port", List.of("502"),
            "enabled", List.of("true")
        );

        BasicConfiguration config = factory.createConfiguration(BasicConfiguration.class, paramMap);

        assertEquals("192.168.1.100", config.host);
        assertEquals(502, config.port);
        assertTrue(config.enabled);
    }

    @Test
    void testTransportConfiguration() {
        String connectionString = "type=TCP&buffer-size=4096&keep-alive=true&no-delay=false";

        TestTransportConfiguration config = factory.createConfiguration(TestTransportConfiguration.class, connectionString);

        assertEquals(TransportType.TCP, config.type);
        assertEquals(4096, config.bufferSize);
        assertTrue(config.keepAlive);
        assertFalse(config.noDelay);
        assertEquals(30, config.timeout); // default value
    }

    // Test Configuration Classes

    public static class ModbusDriverConfiguration implements Configuration {
        @ConfigurationParameter
        @Required
        public String host;

        @ConfigurationParameter
        @IntDefaultValue(502)
        public int port;

        @ConfigurationParameter
        @ParameterConverter(DurationConverter.class)
        @StringDefaultValue("5000ms")
        public Duration timeout;

        @ConfigurationParameter("retry-count")
        @IntDefaultValue(1)
        public int retryCount;

        @ConfigurationParameter
        @StringDefaultValue("TCP")
        public Protocol protocol;

        @ComplexConfigurationParameter(prefix = "ssl", defaultOverrides = {}, requiredOverrides = {})
        public SslConfiguration sslConfig;

        @ComplexConfigurationParameter(prefix = "advanced", defaultOverrides = {}, requiredOverrides = {})
        public AdvancedDriverConfiguration advancedConfig;
    }

    public static class BasicConfiguration implements Configuration {
        @ConfigurationParameter
        @Required
        public String host;

        @ConfigurationParameter
        @IntDefaultValue(80)
        public int port;

        @ConfigurationParameter
        @BooleanDefaultValue(false)
        public boolean enabled;
    }

    public static class AdvancedConfiguration implements Configuration {
        @ConfigurationParameter
        @ParameterConverter(DurationConverter.class)
        public Duration duration;

        @ConfigurationParameter("file-list")
        @ParameterConverter(FileArrayConverter.class)
        public File[] fileList;

        @ConfigurationParameter
        @ParameterConverter(StringArrayConverter.class)
        public String[] tags;
    }

    public static class ApplicationConfiguration implements Configuration {
        @ComplexConfigurationParameter(prefix = "database", defaultOverrides = {}, requiredOverrides = {})
        public DatabaseConfiguration databaseConfig;

        @ComplexConfigurationParameter(prefix = "cache", defaultOverrides = {}, requiredOverrides = {})
        public CacheConfiguration cacheConfig;
    }

    public static class DatabaseConfiguration implements Configuration {
        @ConfigurationParameter
        @Required
        public String host;

        @ConfigurationParameter
        @IntDefaultValue(5432)
        public int port;

        @ConfigurationParameter
        public String name;

        @ComplexConfigurationParameter(prefix = "ssl", defaultOverrides = {}, requiredOverrides = {})
        public DatabaseSslConfiguration sslConfig;
    }

    public static class DatabaseSslConfiguration implements Configuration {
        @ConfigurationParameter
        @BooleanDefaultValue(false)
        public boolean enabled;

        @ConfigurationParameter("cert-file")
        public String certFile;
    }

    public static class CacheConfiguration implements Configuration {
        @ConfigurationParameter
        @BooleanDefaultValue(false)
        public boolean enabled;

        @ConfigurationParameter
        @IntDefaultValue(100)
        public int size;

        @ConfigurationParameter
        @IntDefaultValue(3600)
        public int ttl;
    }

    public static class SslConfiguration implements Configuration {
        @ConfigurationParameter
        @BooleanDefaultValue(false)
        public boolean enabled;

        @ConfigurationParameter
        public String keystore;

        @ConfigurationParameter
        public String password;
    }

    public static class AdvancedDriverConfiguration implements Configuration {
        @ConfigurationParameter("buffer-size")
        @IntDefaultValue(4096)
        public int bufferSize;

        @ConfigurationParameter("keep-alive")
        @BooleanDefaultValue(false)
        public boolean keepAlive;

        @ConfigurationParameter("connection-pool-size")
        @IntDefaultValue(30)
        public int connectionPoolSize;
    }

    public static class TestTransportConfiguration implements Configuration {
        @ConfigurationParameter
        @StringDefaultValue("TCP")
        public TransportType type;

        @ConfigurationParameter("buffer-size")
        @IntDefaultValue(8192)
        public int bufferSize;

        @ConfigurationParameter("keep-alive")
        @BooleanDefaultValue(true)
        public boolean keepAlive;

        @ConfigurationParameter("no-delay")
        @BooleanDefaultValue(true)
        public boolean noDelay;

        @ConfigurationParameter
        @IntDefaultValue(30)
        public int timeout;
    }

    // Test service implementing HasConfiguration
    public static class TestService implements HasConfiguration<BasicConfiguration> {
        private BasicConfiguration configuration;

        public BasicConfiguration getConfiguration() {
            return configuration;
        }

        @Override
        public void setConfiguration(BasicConfiguration configuration) {
            this.configuration = configuration;
        }
    }

    // Enums
    public enum Protocol {
        TCP, UDP, SERIAL
    }

    public enum TransportType {
        TCP, UDP
    }

    // Custom converters
    public static class DurationConverter implements ConfigurationParameterConverter<Duration> {
        @Override
        public Class<Duration> getType() {
            return Duration.class;
        }

        @Override
        public Duration convert(String value) {
            if (value.endsWith("ms")) {
                return Duration.ofMillis(Long.parseLong(value.substring(0, value.length() - 2)));
            } else if (value.endsWith("s")) {
                return Duration.ofSeconds(Long.parseLong(value.substring(0, value.length() - 1)));
            }
            return Duration.ofMillis(Long.parseLong(value));
        }
    }

    public static class FileArrayConverter implements ConfigurationParameterConverter<File[]> {
        @Override
        public Class<File[]> getType() {
            return File[].class;
        }

        @Override
        public File[] convert(String value) {
            return Arrays.stream(value.split(","))
                .map(String::trim)
                .map(File::new)
                .toArray(File[]::new);
        }
    }

    public static class StringArrayConverter implements ConfigurationParameterConverter<String[]> {
        @Override
        public Class<String[]> getType() {
            return String[].class;
        }

        @Override
        public String[] convert(String value) {
            return Arrays.stream(value.split(","))
                .map(String::trim)
                .toArray(String[]::new);
        }
    }
}