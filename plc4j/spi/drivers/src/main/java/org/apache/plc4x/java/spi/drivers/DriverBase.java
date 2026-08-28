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
package org.apache.plc4x.java.spi.drivers;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.api.metadata.Option;
import org.apache.plc4x.java.api.metadata.OptionMetadata;
import org.apache.plc4x.java.api.metadata.PlcDriverMetadata;
import org.apache.plc4x.java.api.types.OptionType;
import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.config.ConfigurationFactory;
import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.Description;
import org.apache.plc4x.java.spi.config.annotations.Required;
import org.apache.plc4x.java.spi.config.annotations.Since;
import org.apache.plc4x.java.spi.config.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.spi.config.annotations.defaults.DoubleDefaultValue;
import org.apache.plc4x.java.spi.config.annotations.defaults.FloatDefaultValue;
import org.apache.plc4x.java.spi.config.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.spi.config.annotations.defaults.LongDefaultValue;
import org.apache.plc4x.java.spi.config.annotations.defaults.ShortDefaultValue;
import org.apache.plc4x.java.spi.config.annotations.defaults.StringDefaultValue;
import org.apache.plc4x.java.spi.drivers.config.ConnectionControlConfiguration;
import org.apache.plc4x.java.spi.drivers.functions.PlcDiscoverer;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcDiscoveryRequest;
import org.apache.plc4x.java.spi.drivers.messages.metadata.DefaultOption;
import org.apache.plc4x.java.spi.drivers.messages.metadata.DefaultOptionMetadata;
import org.apache.plc4x.java.spi.transports.api.DefaultTransportManager;
import org.apache.plc4x.java.spi.transports.api.Transport;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.TransportManager;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogEventType;
import org.apache.plc4x.java.utils.auditlog.api.config.AuditLogConfiguration;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.Field;
import java.lang.reflect.Modifier;
import java.util.ArrayList;
import java.util.Collections;
import java.util.HashSet;
import java.util.List;
import java.util.Optional;
import java.util.Set;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Abstract base class for PLC4X drivers.
 * <p>
 * Handles URI parsing, transport initialization, configuration, and connection factory delegation.
 * Drivers extend this class and implement the abstract methods to provide protocol-specific behavior.
 */
public abstract class DriverBase implements PlcDriver {

    public static final Pattern URI_PATTERN = Pattern.compile(
        "^(?<protocolCode>[a-z0-9\\-]*)(:(?<transportCode>[a-z0-9\\-]*))?://(?<transportConfig>[^?]*)(\\?(?<paramString>.*))?");


    private static final Logger log = LoggerFactory.getLogger(DriverBase.class);

    private final TransportManager transportManager;
    private AuditLog auditLog;

    protected DriverBase() {
        transportManager = new DefaultTransportManager();
    }

    protected abstract Class<? extends Configuration> getConfigurationClass();

    protected Class<? extends TransportConfiguration> getTransportConfigurationClass(Transport<?> transport) {
        return transport.getTransportConfigType();
    }

    public Optional<String> getDefaultTransportCode() {
        return Optional.empty();
    }

    public List<String> getSupportedTransportCodes() {
        return List.of();
    }

    public Set<Integer> defaultPorts(String transportCode) {
        return Collections.emptySet();
    }

    protected abstract ConnectionBase<?> getConnection(Configuration configuration, TransportInstance<?> transportInstance, AuditLog auditLog);

    @Override
    public PlcDiscoveryRequest.Builder discoveryRequestBuilder() {
        if (canDiscover() && (this instanceof PlcDiscoverer plcDiscoverer)) {
            return new DefaultPlcDiscoveryRequest.Builder(plcDiscoverer);
        }
        return PlcDriver.super.discoveryRequestBuilder();
    }

    protected boolean canDiscover() {
        return false;
    }

    protected boolean canBrowse() {
        return false;
    }

    protected boolean canPing() {
        return false;
    }

    protected boolean canRead() {
        return false;
    }

    protected boolean canWrite() {
        return false;
    }

    protected boolean canSubscribe() {
        return false;
    }

    @Override
    public PlcConnection getConnection(String connectionString) throws PlcConnectionException {
        return getConnection(connectionString, null);
    }

    @Override
    public PlcConnection getConnection(String connectionString, PlcAuthentication plcAuthentication) throws PlcConnectionException {
        // Split up the connection string into its individual segments.
        Matcher matcher = URI_PATTERN.matcher(connectionString);
        if (!matcher.matches()) {
            throw new PlcConnectionException(
                "Connection string doesn't match the format '{protocol-code}(:{transport-code})?://{transport-config}(?{parameter-string)?'");
        }
        final String protocolCode = matcher.group("protocolCode");
        String transportCodeMatch = matcher.group("transportCode");
        if (transportCodeMatch == null && getMetadata().getDefaultTransportCode().isEmpty()) {
            throw new PlcConnectionException(
                "This driver has no default transport and no transport code was provided.");
        }
        final String transportCode = (transportCodeMatch != null) ? transportCodeMatch : getMetadata().getDefaultTransportCode().get();
        final String transportConfig = matcher.group("transportConfig");
        final String paramString = matcher.group("paramString");

        // Check if the protocol code matches this driver.
        if (!protocolCode.equals(getProtocolCode())) {
            throw new PlcConnectionException(
                "This driver is not suited to handle this connection string");
        }
        ConfigurationFactory configurationFactory = new ConfigurationFactory();

        // Enforce that the selected transport is one this driver actually supports.
        // Drivers declare their supported transports via getSupportedTransportCodes(); the metadata
        // getter falls back to the single default transport when no explicit list is declared, so a
        // driver that only declares a default still yields a non-empty supported set here. Pairing a
        // driver with a transport it does not support - e.g. a 'tcp' transport with the S7 driver,
        // which speaks COTP - used to be silently accepted and then misbehave; we now fail fast with a
        // clear, actionable message. The 'allow-unsupported-transport' connection option intentionally
        // bypasses ONLY this driver-specific check; it does NOT bypass the 'is the transport registered
        // at all' lookup further below.
        ConnectionControlConfiguration connectionControlConfiguration =
            configurationFactory.createConfiguration(ConnectionControlConfiguration.class, paramString);
        if (!connectionControlConfiguration.isAllowUnsupportedTransport()) {
            List<String> supportedTransportCodes = getMetadata().getSupportedTransportCodes();
            if (!supportedTransportCodes.contains(transportCode)) {
                throw new PlcConnectionException(
                    "Transport '" + transportCode + "' is not supported by driver '" + getProtocolCode()
                        + "'. Supported transports: " + supportedTransportCodes
                        + ". Set 'allow-unsupported-transport=true' in the connection string to use it anyway.");
            }
        }


        // Get the requested transport type.
        Transport<?> transport = transportManager.getTransport(transportCode).orElseThrow(
            () -> new PlcConnectionException("Unsupported transport " + transportCode));

        // Initialize the configuration for the transport.
        Class<? extends TransportConfiguration> transportConfigType = getTransportConfigurationClass(transport);
        TransportConfiguration transportConfiguration = configurationFactory.createPrefixedConfiguration(
            transportConfigType, transportCode, paramString);

        // Create and initialize the audit log.
        AuditLogConfiguration auditLogConfiguration = configurationFactory.createPrefixedConfiguration(
            AuditLogConfiguration.class, "log", paramString);

        // Point out parameters nothing will ever read. The configuration factory ignores what it has
        // no field for, so a misspelt or misplaced parameter looks accepted and silently leaves the
        // default in place - 'remote-slot=2' instead of 'cotp.remote-slot=2' being the classic case.
        warnAboutUnknownParameters(paramString, transportCode, transportConfigType);

        // Logged here rather than on entry: the transport's secrets - a PSK key, a keystore
        // password - are declared on its configuration class, which is only known once the
        // transport has been resolved. Logging earlier would mean redacting without knowing what
        // half the secrets are called.
        log.info("Using connection string: {}", ConnectionStringRedactor.redact(
            connectionString, getConfigurationClass(), transportCode, transportConfigType));
        this.auditLog = AuditLog.builder()
            .withSource(getProtocolCode())
            .withConfiguration(auditLogConfiguration)
            .build();

        // Create an instance for the selected transport
        TransportInstance<?> transportInstance;
        try {
            transportInstance = transport.createTransportInstance(transportConfig, transportConfiguration, auditLog);
        } catch (TransportException e) {
            throw new PlcConnectionException("Unable to create transport instance", e);
        }

        // Create the configuration object for the protocol itself.
        Configuration configuration = configurationFactory.createConfiguration(getConfigurationClass(), paramString);
        if (configuration == null) {
            throw new PlcConnectionException("Unsupported configuration");
        }

        if (auditLog.isEnabled()) {
            auditLog.write(AuditLogEventType.CONFIG, "Starting connection using configuration: " + configuration);
        }

        // Initialize the PlcConnection instance.
        ConnectionBase<?> connection = getConnection(configuration, transportInstance, auditLog);
        connection.setConnectionInfo(getProtocolCode(), getProtocolName(),
            transport.getTransportCode(), transport.getTransportName());
        // Hand the caller-supplied authentication to the connection; previously it was accepted
        // by this method but never propagated, so programmatic credentials were silently ignored.
        connection.setAuthentication(plcAuthentication);
        return connection;
    }

    /**
     * Masks the values of secret-bearing query parameters (passwords, tokens, …) in a connection
     * string so credentials never reach the logs.
     */

    public AuditLog getAuditLog() {
        return auditLog;
    }

    @Override
    public PlcDriverMetadata getMetadata() {
        return new PlcDriverMetadata() {
            @Override
            public Optional<String> getDefaultTransportCode() {
                return DriverBase.this.getDefaultTransportCode();
            }

            @Override
            public List<String> getSupportedTransportCodes() {
                List<String> supportedTransportCodes = DriverBase.this.getSupportedTransportCodes();
                if (supportedTransportCodes.isEmpty() && (getDefaultTransportCode().isPresent())) {
                    return Collections.singletonList(getDefaultTransportCode().get());
                }
                return DriverBase.this.getSupportedTransportCodes();
            }

            @Override
            public Optional<OptionMetadata> getProtocolConfigurationOptionMetadata() {
                Class<? extends Configuration> configurationClass = getConfigurationClass();
                if (configurationClass == null) {
                    return Optional.empty();
                }
                return Optional.of(new DefaultOptionMetadata(extractOptions(configurationClass)));
            }

            @Override
            public Optional<OptionMetadata> getTransportConfigurationOptionMetadata(String transportCode) {
                Transport<?> transport = transportManager.getTransport(transportCode).orElse(null);
                if (transport == null) {
                    return Optional.empty();
                }
                Class<? extends TransportConfiguration> transportConfigurationClass =
                    getTransportConfigurationClass(transport);
                if (transportConfigurationClass == null) {
                    return Optional.empty();
                }
                return Optional.of(new DefaultOptionMetadata(extractOptions(transportConfigurationClass)));
            }

            @Override
            public boolean isDiscoverySupported() {
                return DriverBase.this.canDiscover();
            }
        };
    }

    /**
     * Logs a warning for every parameter in the connection string that none of the configurations
     * involved declares.
     *
     * <p>The logic lives in {@link UnknownParameterReporter}, which is public so the one driver
     * that implements {@code PlcDriver} directly rather than extending this class can report too.
     * The two package-private methods below stay as delegates because this class's tests drive
     * them directly, and keeping them is what shows the extraction changed no behaviour.</p>
     */
    private void warnAboutUnknownParameters(String paramString, String transportCode,
                                            Class<? extends TransportConfiguration> transportConfigType) {
        UnknownParameterReporter.report(
            getProtocolCode(), paramString, transportCode, getConfigurationClass(), transportConfigType);
    }

    static List<String> findUnknownParameters(String paramString, Class<?> driverConfigClass,
                                              Class<?> transportConfigClass, String transportCode) {
        return UnknownParameterReporter.findUnknownParameters(
            paramString, driverConfigClass, transportConfigClass, transportCode);
    }

    static String suggestionFor(String unknown, Set<String> known) {
        return UnknownParameterReporter.suggestionFor(unknown, known);
    }

    /**
     * Walks the full class hierarchy of the given configuration class and returns one
     * {@link Option} for every field annotated with {@link ConfigurationParameter}.
     * Subclass declarations win over superclass declarations for the same parameter key.
     */
    private static List<Option> extractOptions(Class<?> configurationClass) {
        List<Option> options = new ArrayList<>();
        Set<String> seenKeys = new HashSet<>();
        Class<?> current = configurationClass;
        while (current != null && current != Object.class) {
            for (Field field : current.getDeclaredFields()) {
                if (Modifier.isStatic(field.getModifiers())) {
                    continue;
                }
                ConfigurationParameter parameterAnnotation = field.getAnnotation(ConfigurationParameter.class);
                if (parameterAnnotation == null) {
                    continue;
                }
                String key = parameterAnnotation.value();
                if (key == null || key.isEmpty() || !seenKeys.add(key)) {
                    continue;
                }
                options.add(buildOption(field, key));
            }
            current = current.getSuperclass();
        }
        return options;
    }

    private static Option buildOption(Field field, String key) {
        OptionType type = mapOptionType(field.getType());

        Description descriptionAnnotation = field.getAnnotation(Description.class);
        String description = (descriptionAnnotation != null) ? descriptionAnnotation.value() : "";

        boolean required = field.getAnnotation(Required.class) != null;

        Object defaultValue = readDefaultValue(field);

        Since sinceAnnotation = field.getAnnotation(Since.class);
        String since = (sinceAnnotation != null) ? sinceAnnotation.value() : null;

        return new DefaultOption(key, type, description, required, defaultValue, since);
    }

    private static OptionType mapOptionType(Class<?> fieldType) {
        if (fieldType == boolean.class || fieldType == Boolean.class) {
            return OptionType.BOOLEAN;
        }
        if (fieldType == int.class || fieldType == Integer.class
            || fieldType == short.class || fieldType == Short.class) {
            return OptionType.INT;
        }
        if (fieldType == long.class || fieldType == Long.class) {
            return OptionType.LONG;
        }
        if (fieldType == float.class || fieldType == Float.class) {
            return OptionType.FLOAT;
        }
        if (fieldType == double.class || fieldType == Double.class) {
            return OptionType.DOUBLE;
        }
        if (fieldType == String.class || fieldType.isEnum()) {
            return OptionType.STRING;
        }
        if (java.io.File.class.isAssignableFrom(fieldType)) {
            return OptionType.FILE;
        }
        return OptionType.STRUCT;
    }

    private static Object readDefaultValue(Field field) {
        BooleanDefaultValue booleanDefault = field.getAnnotation(BooleanDefaultValue.class);
        if (booleanDefault != null) {
            return booleanDefault.value();
        }
        IntDefaultValue intDefault = field.getAnnotation(IntDefaultValue.class);
        if (intDefault != null) {
            return intDefault.value();
        }
        LongDefaultValue longDefault = field.getAnnotation(LongDefaultValue.class);
        if (longDefault != null) {
            return longDefault.value();
        }
        ShortDefaultValue shortDefault = field.getAnnotation(ShortDefaultValue.class);
        if (shortDefault != null) {
            return shortDefault.value();
        }
        FloatDefaultValue floatDefault = field.getAnnotation(FloatDefaultValue.class);
        if (floatDefault != null) {
            return floatDefault.value();
        }
        DoubleDefaultValue doubleDefault = field.getAnnotation(DoubleDefaultValue.class);
        if (doubleDefault != null) {
            return doubleDefault.value();
        }
        StringDefaultValue stringDefault = field.getAnnotation(StringDefaultValue.class);
        if (stringDefault != null) {
            return stringDefault.value();
        }
        return null;
    }

}
