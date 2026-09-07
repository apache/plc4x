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

import org.apache.plc4x.java.spi.config.ConfigurationFactory;
import org.apache.plc4x.java.spi.config.annotations.ComplexConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.drivers.config.ConnectionControlConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.config.AuditLogConfiguration;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.Field;
import java.lang.reflect.Modifier;
import java.util.ArrayList;
import java.util.Collections;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

import static java.util.stream.Collectors.toList;

/**
 * Reports connection-string parameters that no configuration involved declares.
 *
 * <p>This lived inside {@link DriverBase} and was private, so the one driver that implements
 * {@code PlcDriver} directly - CtrlX - got no reporting at all. The logic is unchanged; it moved
 * here so both entry points can use it, and {@code DriverBase} delegates.</p>
 *
 * <p>One connection string feeds several consumers: the driver configuration, the transport
 * configuration under its code, the audit log under {@code log.} and the connection-control
 * options. A parameter belonging to any one of them looks unknown to the others, so the check is
 * only correct once every consumer has been accounted for - which is why it takes the
 * configuration classes rather than reading them from somewhere.</p>
 */
public final class UnknownParameterReporter {

    private static final Logger log = LoggerFactory.getLogger(UnknownParameterReporter.class);

    private UnknownParameterReporter() {
        // Utility class.
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
    public static void report(String protocolCode, String paramString, String transportCode,
                              Class<?> driverConfigType, Class<?> transportConfigType) {
        List<String> unknown = findUnknownParameters(
            paramString, driverConfigType, transportConfigType, transportCode);
        if (unknown.isEmpty()) {
            return;
        }
        Set<String> known = knownParameterNames(driverConfigType, transportConfigType, transportCode);
        for (String name : unknown) {
            String suggestion = suggestionFor(name, known);
            if (suggestion == null) {
                log.warn("Connection string parameter '{}' is not known to driver '{}' and is ignored.",
                    name, protocolCode);
            } else {
                log.warn("Connection string parameter '{}' is not known to driver '{}' and is ignored - "
                    + "did you mean '{}'?", name, protocolCode, suggestion);
            }
        }
        if (log.isDebugEnabled()) {
            List<String> sorted = new ArrayList<>(known);
            Collections.sort(sorted);
            log.debug("Parameters driver '{}' accepts over transport '{}': {}",
                protocolCode, transportCode, sorted);
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
    public static String suggestionFor(String unknown, Set<String> known) {
        String unknownLeaf = leafOf(unknown);
        List<String> byLeaf = known.stream()
            .filter(candidate -> !candidate.equals(unknown) && leafOf(candidate).equals(unknownLeaf))
            .sorted()
            .toList();
        if (!byLeaf.isEmpty()) {
            return byLeaf.getFirst();
        }

        // Allow roughly one edit per four characters, so short names don't match everything.
        int budget = Math.clamp(unknown.length() / 4, 1, 3);
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
    public static List<String> findUnknownParameters(String paramString, Class<?> driverConfigClass,
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
}
