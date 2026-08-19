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

import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.config.annotations.ComplexConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.junit.jupiter.api.Test;

import java.util.List;
import java.util.Set;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * The configuration factory ignores any parameter it has no field for, so a misspelt or misplaced
 * parameter used to look accepted while silently leaving the default in place - writing
 * {@code remote-slot=2} instead of {@code cotp.remote-slot=2} being the case that keeps costing
 * people an afternoon (GH-2247). These parameters are now reported.
 */
class DriverBaseUnknownParameterTest {

    private static List<String> unknownIn(String paramString) {
        return DriverBase.findUnknownParameters(
            paramString, DriverConfig.class, TransportConfig.class, "cotp");
    }

    @Test
    void acceptsParametersTheDriverDeclares() {
        assertEquals(List.of(), unknownIn("controller-type=S7_300&read-timeout=5000"));
    }

    @Test
    void acceptsTransportParametersUnderTheirPrefix() {
        assertEquals(List.of(), unknownIn("cotp.remote-slot=2&cotp.local-tsap=512"));
    }

    /**
     * The case from the issue: a transport parameter written without its prefix reaches nothing.
     */
    @Test
    void reportsATransportParameterWithoutItsPrefix() {
        assertEquals(List.of("remote-slot"), unknownIn("controller-type=S7_200&remote-slot=2"));
    }

    @Test
    void reportsAMisspeltParameter() {
        assertEquals(List.of("controler-type"), unknownIn("controler-type=S7_300"));
    }

    @Test
    void reportsSeveralInTheOrderTheyWereSupplied() {
        assertEquals(List.of("transport-protocol", "remote-rack"),
            unknownIn("transport-protocol=pniop&remote-rack=0&controller-type=S7_200"));
    }

    /**
     * Parameters inherited from a superclass count - transport configurations are usually a
     * protocol-specific subclass of a shared base.
     */
    @Test
    void acceptsInheritedTransportParameters() {
        assertEquals(List.of(), unknownIn("cotp.remote-tsap=512"));
    }

    /**
     * A complex parameter contributes its nested names under its own prefix. Warning about these
     * would be a false positive, which is the failure mode that would make the whole check a
     * nuisance - OPC UA has exactly this shape with its {@code encoding.*} options.
     */
    @Test
    void acceptsNestedComplexParameters() {
        assertEquals(List.of(), unknownIn("encoding.charset=utf-8"));
    }

    /**
     * The options the SPI itself consumes are not the driver's, but they are still legitimate.
     */
    @Test
    void acceptsSpiOwnedParameters() {
        assertEquals(List.of(), unknownIn("allow-unsupported-transport=true"));
    }

    @Test
    void saysNothingAboutAnEmptyParameterString() {
        assertEquals(List.of(), unknownIn(null));
        assertEquals(List.of(), unknownIn(""));
    }

    /**
     * A driver that declares no configuration at all must not turn every parameter into a warning
     * beyond the ones that are genuinely unknown.
     */
    @Test
    void handlesADriverWithoutAConfigurationClass() {
        List<String> unknown = DriverBase.findUnknownParameters(
            "allow-unsupported-transport=true&nonsense=1", null, TransportConfig.class, "cotp");
        assertEquals(List.of("nonsense"), unknown);
    }

    @Test
    void keepsTheValueOutOfTheReportedName() {
        assertTrue(unknownIn("nonsense=secret-value").contains("nonsense"),
            "the name is reported, not the value");
        assertEquals(List.of("nonsense"), unknownIn("nonsense=secret-value"));
    }

    private static final Set<String> KNOWN = Set.of(
        "controller-type", "read-timeout", "allow-unsupported-transport",
        "cotp.remote-slot", "cotp.remote-rack", "cotp.local-tsap", "log.audit-log-file");

    /**
     * The mistake that actually happens: the right name, without its transport prefix.
     */
    @Test
    void suggestsThePrefixedNameForAnUnprefixedOne() {
        assertEquals("cotp.remote-slot", DriverBase.suggestionFor("remote-slot", KNOWN));
    }

    @Test
    void suggestsTheRightNameForATypo() {
        assertEquals("controller-type", DriverBase.suggestionFor("controler-type", KNOWN));
    }

    /**
     * Nothing is close to this, and inventing a suggestion would be worse than none.
     */
    @Test
    void suggestsNothingWhenNothingIsClose() {
        assertNull(DriverBase.suggestionFor("transport-protocol", KNOWN));
    }

    /**
     * A short name must not match half the parameter list just because everything is within a
     * couple of edits of it.
     */
    @Test
    void doesNotStretchForShortNames() {
        assertNull(DriverBase.suggestionFor("x", KNOWN));
    }

    @Test
    void suggestsDeterministicallyWhenSeveralMatchEquallyWell() {
        Set<String> ambiguous = Set.of("cotp.remote-slot", "tcp.remote-slot");
        assertEquals("cotp.remote-slot", DriverBase.suggestionFor("remote-slot", ambiguous));
        assertEquals("cotp.remote-slot", DriverBase.suggestionFor("remote-slot", ambiguous));
    }

    public static class EncodingConfig implements Configuration {
        @ConfigurationParameter("charset")
        protected String charset = "utf-8";
    }

    public static class DriverConfig implements Configuration {
        @ConfigurationParameter("controller-type")
        protected String controllerType = "ANY";

        @ConfigurationParameter("read-timeout")
        protected int readTimeout = 10000;

        @ComplexConfigurationParameter(prefix = "encoding", defaultOverrides = {}, requiredOverrides = {})
        protected EncodingConfig encoding;
    }

    public static class BaseTransportConfig implements TransportConfiguration {
        @ConfigurationParameter("remote-tsap")
        protected int remoteTsap;
    }

    public static class TransportConfig extends BaseTransportConfig {
        @ConfigurationParameter("remote-slot")
        protected int remoteSlot;

        @ConfigurationParameter("local-tsap")
        protected int localTsap;
    }
}
