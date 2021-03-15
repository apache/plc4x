/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */
package org.apache.plc4x.java.profinet.dcp.protocol;

import java.time.Duration;
import java.util.concurrent.atomic.AtomicInteger;
import org.apache.plc4x.java.profinet.dcp.configuration.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.dcp.readwrite.AllSelector;
import org.apache.plc4x.java.profinet.dcp.readwrite.DCPBlock;
import org.apache.plc4x.java.profinet.dcp.readwrite.DcpIdentRequestPDU;
import org.apache.plc4x.java.profinet.dcp.readwrite.DcpIdentResponsePDU;
import org.apache.plc4x.java.profinet.dcp.readwrite.DeviceProperties;
import org.apache.plc4x.java.profinet.dcp.readwrite.EthernetFrame;
import org.apache.plc4x.java.profinet.dcp.readwrite.IP;
import org.apache.plc4x.java.profinet.dcp.readwrite.IPv4Address;
import org.apache.plc4x.java.profinet.dcp.readwrite.MacAddress;
import org.apache.plc4x.java.profinet.dcp.readwrite.ProfinetFrame;
import org.apache.plc4x.java.profinet.dcp.readwrite.types.DCPBlockOption;
import org.apache.plc4x.java.profinet.dcp.readwrite.types.DCPServiceID;
import org.apache.plc4x.java.profinet.dcp.readwrite.types.DCPServiceType;
import org.apache.plc4x.java.profinet.dcp.readwrite.types.FrameType;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Driver logic for handling Profinet-DCP packets.
 */
public class ProfinetDCPProtocolLogic extends Plc4xProtocolBase<EthernetFrame> implements
    HasConfiguration<ProfinetConfiguration> {

    public static MacAddress PROFINET_BROADCAST = createAddress(0x01, 0x0E, 0xCF, 0x00, 0x00, 0x00);
    public static int PN_DCP = 0x8892;

    public static final Duration REQUEST_TIMEOUT = Duration.ofMillis(10000);
    private final Logger logger = LoggerFactory.getLogger(ProfinetDCPProtocolLogic.class);

    private AtomicInteger invokeId = new AtomicInteger(0);
    private ProfinetConfiguration configuration;

    @Override
    public void setConfiguration(ProfinetConfiguration configuration) {
        this.configuration = configuration;
    }

    @Override
    public void onConnect(ConversationContext<EthernetFrame> context) {
        DCPServiceID serviceId = DCPServiceID.IDENTIFY;
        DCPServiceType serviceType = DCPServiceType.REQUEST;
        long xid = invokeId.incrementAndGet();
        int responseDelay = 1000;
        DCPBlock[] blocks = new DCPBlock[] { new AllSelector(DCPBlockOption.ALL_SELECTOR)};
        int dcpDataLength = blocks[0].getLengthInBytes();

        DcpIdentRequestPDU requestPDU = new DcpIdentRequestPDU(serviceId, serviceType, xid, responseDelay, dcpDataLength, blocks);
        EthernetFrame ethernetFrame = new EthernetFrame(PROFINET_BROADCAST, configuration.getSender(), PN_DCP,
            new ProfinetFrame(FrameType.IDENTIFY_MULTICAST_REQUEST, requestPDU));

        // this is broadcast thus reply might come from multiple participants
        context.sendToWire(ethernetFrame);
        context.fireConnected();
    }

    @Override
    public void onDisconnect(ConversationContext<EthernetFrame> context) {
        context.fireDisconnected();
    }

    @Override
    protected void decode(ConversationContext<EthernetFrame> context, EthernetFrame msg) throws Exception {
        if (msg.getEthernetType() != PN_DCP) {
            logger.trace("Discarding unwanted frame type {}", msg.getEthernetType());
        }

        ProfinetFrame profinetFrame = msg.getPayload();
        if (profinetFrame.getFrameType() == FrameType.IDENTIFY_RESPONSE) {
            logger.info("Ident response from Profinet device:");
            if (profinetFrame.getFrame() instanceof DcpIdentResponsePDU) {
                DcpIdentResponsePDU response = (DcpIdentResponsePDU) profinetFrame.getFrame();
                DCPBlock[] blocks = response.getBlocks();
                for (int index = 0; index < blocks.length; index++) {
                    DCPBlock block = blocks[index];
                    if (block instanceof IP) {
                        IP ip = (IP) block;
                        logger.info("Device IP: {}, mask: {}, gateway: {}", addressString(ip.getIpAddress()), addressString(ip.getSubnetMask()), addressString(ip.getStandardGateway()));
                    }
                    if (block instanceof DeviceProperties) {
                        DeviceProperties properties = (DeviceProperties) block;
                        logger.info("Device option: {}, value: {}", properties.getSubOption().name(), properties.getProperties());
                    }
                }
            } else {
                logger.error("Unexpected ident response {}", profinetFrame.getFrame().getClass());
            }
        }
    }

    private String addressString(IPv4Address address) {
        return address.getOctet1() + "." + address.getOctet2() + "." + address.getOctet3() + "." + address.getOctet4();
    }

    @Override
    public void close(ConversationContext<EthernetFrame> context) {

    }

    public final static MacAddress createAddress(int octet1, int octet2, int octet3, int octet4, int octet5, int octet6) {
        return new MacAddress((short) octet1, (short) octet2, (short) octet3, (short) octet4, (short) octet5, (short) octet6);
    }

}
