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
package org.apache.plc4x.java.spi.connection;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.spi.configuration.PlcConfiguration;
import org.apache.plc4x.java.spi.configuration.PlcConnectionConfiguration;
import org.apache.plc4x.java.spi.configuration.PlcTransportConfiguration;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.metadata.Option;
import org.apache.plc4x.java.api.metadata.OptionMetadata;
import org.apache.plc4x.java.api.types.OptionType;
import org.apache.plc4x.java.api.metadata.PlcDriverMetadata;
import org.apache.plc4x.java.spi.configuration.ConfigurationFactory;
import org.apache.plc4x.java.spi.configuration.annotations.*;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.*;
import org.apache.plc4x.java.spi.generation.Message;
import org.apache.plc4x.java.spi.metadata.DefaultOption;
import org.apache.plc4x.java.spi.metadata.DefaultOptionMetadata;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.apache.plc4x.java.spi.transport.Transport;
import org.apache.plc4x.java.spi.values.DefaultPlcValueHandler;
import org.apache.plc4x.java.spi.values.PlcValueHandler;

import java.lang.reflect.Field;
import java.util.*;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import java.util.stream.Collectors;

import static org.apache.plc4x.java.spi.configuration.ConfigurationFactory.configure;

public abstract class GeneratedDriverBase<BASE_PACKET extends Message> implements PlcDriver {

    public static final String PROPERTY_PLC4X_FORCE_FIRE_DISCOVER_EVENT = "PLC4X_FORCE_FIRE_DISCOVER_EVENT";
    public static final String PROPERTY_PLC4X_FORCE_AWAIT_SETUP_COMPLETE = "PLC4X_FORCE_AWAIT_SETUP_COMPLETE";
    public static final String PROPERTY_PLC4X_FORCE_AWAIT_DISCONNECT_COMPLETE = "PLC4X_FORCE_AWAIT_DISCONNECT_COMPLETE";
    public static final String PROPERTY_PLC4X_FORCE_AWAIT_DISCOVER_COMPLETE = "PLC4X_FORCE_AWAIT_DISCOVER_COMPLETE";

    public static final Pattern URI_PATTERN = Pattern.compile(
        "^(?<protocolCode>[a-z0-9\\-]*)(:(?<transportCode>[a-z0-9]*))?://(?<transportConfig>[^?]*)(\\?(?<paramString>.*))?");

    /**
     * Configuration class for configuration introspection
     *
     * @return the configuration class
     */
    protected abstract Class<? extends PlcConnectionConfiguration> getConfigurationClass();

    protected Optional<Class<? extends PlcTransportConfiguration>> getTransportConfigurationClass(String transportCode) {
        return Optional.empty();
    }

    protected Optional<String> getDefaultTransportCode() {
        return Optional.empty();
    }

    protected List<String> getSupportedTransportCodes() {
        return Collections.emptyList();
    }

    @Override
    public PlcDriverMetadata getMetadata() {
        return new PlcDriverMetadata() {
            @Override
            public Optional<String> getDefaultTransportCode() {
                return GeneratedDriverBase.this.getDefaultTransportCode();
            }

            @Override
            public List<String> getSupportedTransportCodes() {
                List<String> supportedTransportCodes = GeneratedDriverBase.this.getSupportedTransportCodes();
                if (supportedTransportCodes.isEmpty() && (getDefaultTransportCode().isPresent())) {
                    return Collections.singletonList(getDefaultTransportCode().get());
                }
                return GeneratedDriverBase.this.getSupportedTransportCodes();
            }

            @Override
            public Optional<OptionMetadata> getProtocolConfigurationOptionMetadata() {
                var clazz = getConfigurationClass();
                if (clazz == null) {
                    return Optional.empty();
                }
                var options = getOptions(clazz);
                return Optional.of(new DefaultOptionMetadata(options));
            }

            @Override
            public Optional<OptionMetadata> getTransportConfigurationOptionMetadata(String transportCode) {
                var clazzOption = resolveTransportConfigurationClass(transportCode);
                if (clazzOption.isEmpty()) {
                    return Optional.empty();
                }
                var clazz = clazzOption.get();
                var options = getOptions(clazz);
                return Optional.of(new DefaultOptionMetadata(options));
            }

            private List<Option> getOptions(Class<? extends PlcConfiguration> clazz) {
                return getAllFields(clazz).stream()
                    .map(this::optionsForField)
                    .flatMap(Collection::stream)
                    .filter(Objects::nonNull)
                    .map(Option.class::cast)
                    .collect(Collectors.toList());
            }

            private List<DefaultOption> optionsForField(Field field) {
                // check if this is a complex configuration parameter and bail early
                var complexConfigurationParameterAnnotation = field.getAnnotation(ComplexConfigurationParameter.class);
                if (complexConfigurationParameterAnnotation != null) {
                    var prefix = complexConfigurationParameterAnnotation.prefix();
                    if (PlcConfiguration.class.isAssignableFrom(field.getType())) {
                        return getOptions((Class<? extends PlcConfiguration>) field.getType())
                            .stream()
                            .map(option -> new DefaultOption(
                                prefix + "." + option.getKey(),
                                option.getType(),
                                option.getDescription(),
                                option.isRequired(),
                                option.getDefaultValue().orElse(null),
                                option.getSince().orElse(null)
                            ))
                            .collect(Collectors.toList());
                    }
                }
                String key = null;
                var configurationParameterAnnotation = field.getAnnotation(ConfigurationParameter.class);
                if (configurationParameterAnnotation != null) {
                    key = configurationParameterAnnotation.value();
                    if (key.isEmpty()) {
                        key = field.getName();
                    }
                }
                if (key == null) {
                    return Collections.emptyList();
                }
                String description = "";
                var descriptionAnnotation = field.getAnnotation(Description.class);
                if (descriptionAnnotation != null) {
                    description = descriptionAnnotation.value();
                }
                boolean required = false;
                var requiredAnnotation = field.getAnnotation(Required.class);
                if (requiredAnnotation != null) {
                    required = true;
                }
                OptionType type;
                switch (field.getType().getSimpleName()) {
                    case "boolean":
                    case "Boolean":
                        type = OptionType.BOOLEAN;
                        break;
                    case "float":
                    case "Float":
                        type = OptionType.FLOAT;
                        break;
                    case "double":
                    case "Double":
                        type = OptionType.DOUBLE;
                        break;
                    case "int":
                    case "Integer":
                        type = OptionType.INT;
                        break;
                    case "long":
                    case "Long":
                        type = OptionType.LONG;
                        break;
                    case "String":
                        type = OptionType.STRING;
                        break;
                    default:
                        // If there's a property-converter, use "STRING" as type.
                        var parameterConverterAnnotation = field.getAnnotation(ParameterConverter.class);
                        if (parameterConverterAnnotation != null) {
                            type = OptionType.STRING;
                        } else {
                            type = OptionType.STRUCT;
                        }
                        break;
                }
                Object defaultValue = null;
                var booleanDefaultValueAnnotation = field.getAnnotation(BooleanDefaultValue.class);
                if (booleanDefaultValueAnnotation != null) {
                    defaultValue = booleanDefaultValueAnnotation.value();
                }
                var doubleDefaultValueAnnotation = field.getAnnotation(DoubleDefaultValue.class);
                if (doubleDefaultValueAnnotation != null) {
                    defaultValue = doubleDefaultValueAnnotation.value();
                }
                var floatDefaultValueAnnotation = field.getAnnotation(FloatDefaultValue.class);
                if (floatDefaultValueAnnotation != null) {
                    defaultValue = floatDefaultValueAnnotation.value();
                }
                var intDefaultValueAnnotation = field.getAnnotation(IntDefaultValue.class);
                if (intDefaultValueAnnotation != null) {
                    defaultValue = intDefaultValueAnnotation.value();
                }
                var longDefaultValueAnnotation = field.getAnnotation(LongDefaultValue.class);
                if (longDefaultValueAnnotation != null) {
                    defaultValue = longDefaultValueAnnotation.value();
                }
                var stringDefaultValueAnnotation = field.getAnnotation(StringDefaultValue.class);
                if (stringDefaultValueAnnotation != null) {
                    type = OptionType.STRING;
                    defaultValue = stringDefaultValueAnnotation.value();
                }
                String since = null;
                var sinceAnnotation = field.getAnnotation(Since.class);
                if(sinceAnnotation != null) {
                    since = sinceAnnotation.value();
                }
                return Collections.singletonList(new DefaultOption(key, type, description, required, defaultValue, since));
            }

            @Override
            public boolean isDiscoverySupported() {
                return canDiscover();
            }

        };
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

    protected boolean canBrowse() {
        return false;
    }

    protected boolean canDiscover() {
        return false;
    }

    protected boolean fireDiscoverEvent() {
        return false;
    }

    protected boolean awaitSetupComplete() {
        return true;
    }

    protected boolean awaitDisconnectComplete() {
        return false;
    }

    protected boolean awaitDiscoverComplete() {
        return false;
    }

    protected BaseOptimizer getOptimizer() {
        return null;
    }

    protected PlcValueHandler getValueHandler() {
        return new DefaultPlcValueHandler();
    }

    protected abstract ProtocolStackConfigurer<BASE_PACKET> getStackConfigurer();

    protected ProtocolStackConfigurer<BASE_PACKET> getStackConfigurer(Transport transport) {
        return getStackConfigurer();
    }

    protected void initializePipeline(ChannelFactory channelFactory) {
        // Override in derived drivers.
    }

    @Override
    public PlcConnection getConnection(String connectionString) throws PlcConnectionException {
        return getConnection(connectionString, null);
    }

    @Override
    public PlcConnection getConnection(String connectionString, PlcAuthentication authentication) throws PlcConnectionException {
        ConfigurationFactory configurationFactory = new ConfigurationFactory();
        // Split up the connection string into its individual segments.
        Matcher matcher = URI_PATTERN.matcher(connectionString);
        if (!matcher.matches()) {
            throw new PlcConnectionException(
                "Connection string doesn't match the format '{protocol-code}:({transport-code})?//{transport-config}(?{parameter-string)?'");
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
            // Actually this shouldn't happen as the DriverManager should have not used this driver in the first place.
            throw new PlcConnectionException(
                "This driver is not suited to handle this connection string");
        }

        // Create the configuration object.
        PlcConnectionConfiguration configuration = configurationFactory
            .createConfiguration(getConfigurationClass(), protocolCode, transportCode, transportConfig, paramString);
        if (configuration == null) {
            throw new PlcConnectionException("Unsupported configuration");
        }

        // Try to find a suitable transport-type for creating the communication channel.
        Transport transport = null;
        ServiceLoader<Transport> transportLoader = ServiceLoader.load(
            Transport.class, Thread.currentThread().getContextClassLoader());
        for (Transport curTransport : transportLoader) {
            if (curTransport.getTransportCode().equals(transportCode)) {
                transport = curTransport;
                break;
            }
        }
        if (transport == null) {
            throw new PlcConnectionException("Unsupported transport " + transportCode);
        }

        // Find out the type of the transport configuration.
        Class<? extends PlcTransportConfiguration> transportConfigurationType = transport.getTransportConfigType();
        if (getTransportConfigurationClass(transportCode).isPresent()) {
            transportConfigurationType = getTransportConfigurationClass(transportCode).get();
        }
        // Use the transport configuration type to actually configure the transport instance.
        PlcTransportConfiguration plcTransportConfiguration = configurationFactory
            .createTransportConfiguration(transportConfigurationType,
                protocolCode, transportCode, transportConfig, paramString);
        configure(plcTransportConfiguration, transport);

        // Create an instance of the communication channel which the driver should use.
        ChannelFactory channelFactory = transport.createChannelFactory(transportConfig);
        if (channelFactory == null) {
            throw new PlcConnectionException("Unable to get channel factory from url " + transportConfig);
        }
        configure(configuration, channelFactory);

        // Give drivers the option to customize the channel.
        initializePipeline(channelFactory);

        // Make the "fire discover event" overridable via system property.
        boolean fireDiscoverEvent = fireDiscoverEvent();
        if (System.getProperty(PROPERTY_PLC4X_FORCE_FIRE_DISCOVER_EVENT) != null) {
            fireDiscoverEvent = Boolean.parseBoolean(System.getProperty(PROPERTY_PLC4X_FORCE_FIRE_DISCOVER_EVENT));
        }

        // Make the "await setup complete" overridable via system property.
        boolean awaitSetupComplete = awaitSetupComplete();
        if (System.getProperty(PROPERTY_PLC4X_FORCE_AWAIT_SETUP_COMPLETE) != null) {
            awaitSetupComplete = Boolean.parseBoolean(System.getProperty(PROPERTY_PLC4X_FORCE_AWAIT_SETUP_COMPLETE));
        }

        // Make the "await disconnect complete" overridable via system property.
        boolean awaitDisconnectComplete = awaitDisconnectComplete();
        if (System.getProperty(PROPERTY_PLC4X_FORCE_AWAIT_DISCONNECT_COMPLETE) != null) {
            awaitDisconnectComplete = Boolean.parseBoolean(System.getProperty(PROPERTY_PLC4X_FORCE_AWAIT_DISCONNECT_COMPLETE));
        }

        // Make the "await disconnect complete" overridable via system property.
        boolean awaitDiscoverComplete = awaitDiscoverComplete();
        if (System.getProperty(PROPERTY_PLC4X_FORCE_AWAIT_DISCOVER_COMPLETE) != null) {
            awaitDiscoverComplete = Boolean.parseBoolean(System.getProperty(PROPERTY_PLC4X_FORCE_AWAIT_DISCOVER_COMPLETE));
        }

        return new DefaultNettyPlcConnection(
            canPing(), canRead(), canWrite(), canSubscribe(), canBrowse(),
            getValueHandler(),
            configuration,
            channelFactory,
            fireDiscoverEvent,
            awaitSetupComplete,
            awaitDisconnectComplete,
            awaitDiscoverComplete,
            getStackConfigurer(transport),
            getOptimizer(),
            authentication);
    }

    protected Optional<Class<? extends PlcTransportConfiguration>> resolveTransportConfigurationClass(String transportCode) {
        if (getTransportConfigurationClass(transportCode).isPresent()) {
            return Optional.of(getTransportConfigurationClass(transportCode).get());
        }

        // Try to find a suitable transport-type for creating the communication channel.
        Transport transport = null;
        ServiceLoader<Transport> transportLoader = ServiceLoader.load(
            Transport.class, Thread.currentThread().getContextClassLoader());
        for (Transport curTransport : transportLoader) {
            if (curTransport.getTransportCode().equals(transportCode)) {
                transport = curTransport;
                break;
            }
        }
        if (transport == null) {
            return Optional.empty();
        }

        // Find out the type of the transport configuration.
        Class<? extends PlcTransportConfiguration> transportConfigurationType = transport.getTransportConfigType();
        return Optional.of(transportConfigurationType);
    }

    protected List<Field> getAllFields(Class<?> type) {
        List<Field> fields = new ArrayList<>(Arrays.asList(type.getDeclaredFields()));
        if (type.getSuperclass() != null) {
            fields.addAll(getAllFields(type.getSuperclass()));
        }
        return fields;
    }

}
