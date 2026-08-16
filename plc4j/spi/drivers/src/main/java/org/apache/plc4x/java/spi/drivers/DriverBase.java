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
import org.apache.plc4x.java.spi.config.annotations.ComplexConfigurationParameter;
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

import static java.util.stream.Collectors.toList;

/**
 * Abstract base class for PLC4X drivers.
 * <p>
 * Handles URI parsing, transport initialization, configuration, and connection factory delegation.
 * Drivers extend this class and implement the abstract methods to provide protocol-specific behavior.
 */
public abstract class DriverBase implements PlcDriver {

    public static final Pattern URI_PATTERN = Pattern.compile(
        "^(?<protocolCode>[a-z0-9\\-]*)(:(?<transportCode>[a-z0-9\\-]*))?://(?<transportConfig>[^?]*)(\\?(?<paramString>.*))?");

    /**
     * Matches the value of connection-string query parameters that carry secrets so they can be
     * masked before logging. Covers any parameter whose name contains "password", "passwd", "secret"
     * or "token" (optionally with a transport prefix like "tls.").
     */
    private static final Pattern SECRET_PARAM_PATTERN = Pattern.compile(
        "(?i)([?&][^=&]*(?:password|passwd|secret|token)[^=&]*=)[^&]*");

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
        log.info("Using connection string: {}", redactSecrets(connectionString));

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
    static String redactSecrets(String connectionString) {
        if (connectionString == null) {
            return null;
        }
        return SECRET_PARAM_PATTERN.matcher(connectionString).replaceAll("$1***");
    }

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
     * involved declares, listing what was expected.
     * <p>
     * This is deliberately a warning and not an exception: the driver, the transport, the audit log
     * and the connection control options are the configurations this class knows about, but a driver
     * is free to read further prefixed configurations of its own, and those parameters would look
     * unknown from here. Failing the connection over that would break working setups; saying
     * something turns a silent misconfiguration into a visible one.
     */
    private void warnAboutUnknownParameters(String paramString, String transportCode,
                                            Class<? extends TransportConfiguration> transportConfigType) {
        List<String> unknown = findUnknownParameters(
            paramString, getConfigurationClass(), transportConfigType, transportCode);
        if (unknown.isEmpty()) {
            return;
        }
        Set<String> known = knownParameterNames(getConfigurationClass(), transportConfigType, transportCode);
        for (String name : unknown) {
            String suggestion = suggestionFor(name, known);
            if (suggestion == null) {
                log.warn("Connection string parameter '{}' is not known to driver '{}' and is ignored.",
                    name, getProtocolCode());
            } else {
                log.warn("Connection string parameter '{}' is not known to driver '{}' and is ignored - "
                    + "did you mean '{}'?", name, getProtocolCode(), suggestion);
            }
        }
        if (log.isDebugEnabled()) {
            List<String> sorted = new ArrayList<>(known);
            Collections.sort(sorted);
            log.debug("Parameters driver '{}' accepts over transport '{}': {}",
                getProtocolCode(), transportCode, sorted);
        }
    }

    /**
     * The known parameter the given unknown one was most likely meant to be, or {@code null} when
     * nothing is close enough to be worth suggesting.
     * <p>
     * A missing or wrong prefix is the mistake that actually happens - {@code remote-slot} for
     * {@code cotp.remote-slot} - so a name that matches some known parameter's last segment wins
     * outright. Otherwise a short edit distance catches ordinary typos.
     */
    static String suggestionFor(String unknown, Set<String> known) {
        String unknownLeaf = leafOf(unknown);
        List<String> byLeaf = known.stream()
            .filter(candidate -> !candidate.equals(unknown) && leafOf(candidate).equals(unknownLeaf))
            .sorted()
            .collect(toList());
        if (!byLeaf.isEmpty()) {
            return byLeaf.get(0);
        }

        // Allow roughly one edit per four characters, so short names don't match everything.
        int budget = Math.min(3, Math.max(1, unknown.length() / 4));
        String best = null;
        int bestDistance = Integer.MAX_VALUE;
        for (String candidate : known) {
            int distance = editDistance(unknown, candidate);
            if (distance <= budget && ((best == null) || (distance < bestDistance)
                || ((distance == bestDistance) && (candidate.compareTo(best) < 0)))) {
                best = candidate;
                bestDistance = distance;
            }
        }
        return best;
    }

    /** The part after the last dot - the parameter name without any prefix. */
    private static String leafOf(String name) {
        int lastDot = name.lastIndexOf('.');
        return (lastDot < 0) ? name : name.substring(lastDot + 1);
    }

    /**
     * Levenshtein distance. Hand-rolled to keep the SPI free of another dependency; the strings
     * involved are parameter names, so the quadratic cost is irrelevant.
     */
    private static int editDistance(String left, String right) {
        int[] previous = new int[right.length() + 1];
        int[] current = new int[right.length() + 1];
        for (int j = 0; j <= right.length(); j++) {
            previous[j] = j;
        }
        for (int i = 1; i <= left.length(); i++) {
            current[0] = i;
            for (int j = 1; j <= right.length(); j++) {
                int substitution = previous[j - 1] + (left.charAt(i - 1) == right.charAt(j - 1) ? 0 : 1);
                current[j] = Math.min(substitution, Math.min(previous[j] + 1, current[j - 1] + 1));
            }
            int[] swap = previous;
            previous = current;
            current = swap;
        }
        return previous[right.length()];
    }

    /**
     * The parameters in {@code paramString} that none of the configurations involved declares, in the
     * order they were supplied. Package-private so it can be tested without a registered transport.
     */
    static List<String> findUnknownParameters(String paramString, Class<?> driverConfigClass,
                                              Class<?> transportConfigClass, String transportCode) {
        Set<String> supplied = ConfigurationFactory.parameterNames(paramString);
        if (supplied.isEmpty()) {
            return Collections.emptyList();
        }
        Set<String> known = knownParameterNames(driverConfigClass, transportConfigClass, transportCode);
        return supplied.stream().filter(name -> !known.contains(name)).collect(toList());
    }

    private static Set<String> knownParameterNames(Class<?> driverConfigClass, Class<?> transportConfigClass,
                                                   String transportCode) {
        Set<String> known = new HashSet<>();
        collectParameterNames(driverConfigClass, "", known);
        collectParameterNames(transportConfigClass, transportCode + ".", known);
        collectParameterNames(AuditLogConfiguration.class, "log.", known);
        collectParameterNames(ConnectionControlConfiguration.class, "", known);
        return known;
    }

    /**
     * Adds every parameter name the given configuration class declares to {@code names}, prefixed
     * with {@code prefix}. Complex parameters contribute their nested names under their own prefix,
     * so an OPC UA {@code encoding.*} parameter is recognised rather than reported as unknown.
     */
    private static void collectParameterNames(Class<?> configurationClass, String prefix, Set<String> names) {
        if (configurationClass == null) {
            return;
        }
        Class<?> current = configurationClass;
        while (current != null && current != Object.class) {
            for (Field field : current.getDeclaredFields()) {
                if (Modifier.isStatic(field.getModifiers())) {
                    continue;
                }
                ConfigurationParameter parameterAnnotation = field.getAnnotation(ConfigurationParameter.class);
                if (parameterAnnotation != null && !parameterAnnotation.value().isEmpty()) {
                    names.add(prefix + parameterAnnotation.value());
                    continue;
                }
                ComplexConfigurationParameter complexAnnotation =
                    field.getAnnotation(ComplexConfigurationParameter.class);
                if (complexAnnotation != null) {
                    String nestedPrefix = complexAnnotation.prefix().isEmpty()
                        ? prefix : prefix + complexAnnotation.prefix() + ".";
                    collectParameterNames(field.getType(), nestedPrefix, names);
                }
            }
            current = current.getSuperclass();
        }
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
