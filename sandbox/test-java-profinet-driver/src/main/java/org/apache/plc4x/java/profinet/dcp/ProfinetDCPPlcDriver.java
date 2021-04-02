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
package org.apache.plc4x.java.profinet.dcp;

import java.util.function.Consumer;
import java.util.function.ToIntFunction;
import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.apache.plc4x.java.profinet.dcp.configuration.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.dcp.field.ProfinetFieldHandler;
import org.apache.plc4x.java.profinet.dcp.protocol.ProfinetDCPProtocolLogic;
import org.apache.plc4x.java.profinet.dcp.readwrite.BaseEthernetFrame;
import org.apache.plc4x.java.profinet.dcp.readwrite.EthernetFrame;
import org.apache.plc4x.java.profinet.dcp.readwrite.io.BaseEthernetFrameIO;
import org.apache.plc4x.java.profinet.dcp.readwrite.io.EthernetFrameIO;
import org.apache.plc4x.java.profinet.dcp.readwrite.types.DCPServiceID;
import org.apache.plc4x.java.profinet.dcp.readwrite.types.DCPServiceType;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.spi.values.IEC61131ValueHandler;

/**
 * Proof of concept implementation of Profinet DCP driver.
 */
public class ProfinetDCPPlcDriver extends GeneratedDriverBase<BaseEthernetFrame> {
    static private final Logger LOG = LoggerFactory.getLogger(ProfinetDCPPlcDriver.class);

    @Override
    public String getProtocolCode() {
        return "profinet-dcp";
    }

    @Override
    public String getProtocolName() {
        return "Profinet Discovery and Configuration Protocol";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationType() {
        return ProfinetConfiguration.class;
    }

    @Override
    protected String getDefaultTransport() {
        return "raw";
    }

    @Override
    protected ProfinetFieldHandler getFieldHandler() {
        return new ProfinetFieldHandler();
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return new IEC61131ValueHandler();
    }

    @Override
    protected ProtocolStackConfigurer<BaseEthernetFrame> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(BaseEthernetFrame.class, BaseEthernetFrameIO.class)
            .withProtocol(ProfinetDCPProtocolLogic.class)
            .withPacketSizeEstimator(ProfinetPacketEstimator.class)
            //.withCorruptPacketRemover(CorruptEthernetFrameRemover.class)
            .build();
    }

    public static class ProfinetPacketEstimator implements ToIntFunction<ByteBuf> {
        @Override
        public int applyAsInt(ByteBuf buffer) {
            // When frame uses vlan tag then getShort(12) will return different type, actual ethernet type is shifted by
            // 4 bytes, after type & vlan tag information.
            short protocol = buffer.getShort(12);
            int bufferSize = buffer.readableBytes();
            if (protocol == ProfinetDCPProtocolLogic.VLAN && bufferSize >= 26) {
                return size(4, buffer);
            }

            if (protocol == ProfinetDCPProtocolLogic.PN_DCP && bufferSize >= 24) {
                return size(0, buffer);
            }
            LOG.debug("SIZE NOT CALLED; buffer.getShort(12)= {} and size= {}",protocol,bufferSize);
            return -1;
        }

        private int size(int offset, ByteBuf buffer) {
            int DCPDataLeng =  buffer.getUnsignedShort(24 + offset);
            int val = DCPDataLeng + 26 + offset;
            // profinet frames have fixed overhead of 26 bytes. First 16 is ethernet header, another 10 bytes are
            // needed to reach DCP block length field which gives an base for frame size calculation.
            /*if (buffer.getByte(16 + offset) == DCPServiceID.IDENTIFY.getValue() &&
                buffer.getByte(17 + offset) == DCPServiceType.REQUEST.getValue()) {
                // identify request has fixed size
                return 1024;
            }*/
            // Noticed that all pn_dcp packages (DCPServiceID.IDENTIFY + DCPServiceType.REQUEST) that have size
            // less than 60 will contain 0x00 padding to match min 60 length
            if (val < 60) val=60;
            LOG.info("DCPDataLen= {} and total len(+26)= {}",DCPDataLeng,val);
            return val;
        }
    }

    public static class CorruptEthernetFrameRemover implements Consumer<ByteBuf> {

        @Override
        public void accept(ByteBuf byteBuf) {
            if (byteBuf.getShort(12) != ProfinetDCPProtocolLogic.PN_DCP) {
                // reset buffer
                byteBuf.readBytes(byteBuf.capacity());
            }
        }
    }

}
