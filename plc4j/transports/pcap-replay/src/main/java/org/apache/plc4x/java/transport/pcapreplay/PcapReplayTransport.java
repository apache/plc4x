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
package org.apache.plc4x.java.transport.pcapreplay;

import org.apache.plc4x.java.spi.transports.api.Transport;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.pcapreplay.config.PcapReplayTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.pcap4j.core.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class PcapReplayTransport implements Transport<PcapReplayTransportConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(PcapReplayTransport.class);

    @Override
    public String getTransportCode() {
        return "pcap-replay";
    }

    @Override
    public String getTransportName() {
        return "PCAP Replay";
    }

    @Override
    public Class<PcapReplayTransportConfiguration> getTransportConfigType() {
        return PcapReplayTransportConfiguration.class;
    }

    @Override
    public TransportInstance<PcapReplayTransportConfiguration> createTransportInstance(
        String transportUrl, TransportConfiguration configuration, AuditLog auditLog) throws TransportException {
        if (!(configuration instanceof PcapReplayTransportConfiguration pcapReplayTransportConfiguration)) {
            throw new IllegalArgumentException(String.format("Expected configuration of type %s but got %s",
                PcapReplayTransportConfiguration.class.getSimpleName(), configuration.getClass().getSimpleName()));
        }

        // The address segment of the connection string names the capture -
        // "pcap-replay:///captures/line-3.pcapng" - the way every other transport is addressed by
        // it. The "pcap-file" parameter stays as the alternative for callers that assemble a
        // connection string out of options alone; naming the file in both places is contradictory,
        // and the address segment wins.
        if (transportUrl != null && !transportUrl.trim().isEmpty()) {
            pcapReplayTransportConfiguration.pcapFile = transportUrl.trim();
        }
        // Checked here rather than with @Required on the configuration field: that check runs while
        // the configuration is being built, before the address segment has been seen, and would
        // reject the very form this method accepts.
        if (pcapReplayTransportConfiguration.pcapFile == null
                || pcapReplayTransportConfiguration.pcapFile.trim().isEmpty()) {
            throw new TransportException("No PCAP file given. Name it in the connection string "
                + "('pcap-replay:///captures/line-3.pcapng') or with the 'pcap-replay.pcap-file' option.");
        }

        LOGGER.debug("Creating PCAP replay transport for file: {}", pcapReplayTransportConfiguration.pcapFile);
        return new PcapReplayTransportInstance(pcapReplayTransportConfiguration, auditLog);
    }

}
