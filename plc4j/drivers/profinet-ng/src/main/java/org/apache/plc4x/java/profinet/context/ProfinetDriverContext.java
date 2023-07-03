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

package org.apache.plc4x.java.profinet.context;

import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.profinet.gsdml.ProfinetDeviceAccessPointItem;
import org.apache.plc4x.java.profinet.gsdml.ProfinetModuleItem;
import org.apache.plc4x.java.profinet.gsdml.ProfinetVirtualSubmoduleItem;
import org.apache.plc4x.java.profinet.readwrite.DceRpc_ActivityUuid;
import org.apache.plc4x.java.profinet.readwrite.DceRpc_ObjectUuid;
import org.apache.plc4x.java.profinet.readwrite.MacAddress;
import org.apache.plc4x.java.profinet.readwrite.Uuid;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.generation.*;

import java.util.List;
import java.util.Map;
import java.util.UUID;
import java.util.concurrent.atomic.AtomicInteger;

public class ProfinetDriverContext implements DriverContext {

    public static final int DEFAULT_UDP_PORT = 34964;
    public static final int DEFAULT_ARGS_MAXIMUM = 16696;
    public static final int DEFAULT_MAX_ARRAY_COUNT = 16696;
    public static final int DEFAULT_ACTIVITY_TIMEOUT = 600;
    public static final int UDP_RT_PORT = 0x8892;
    public static final short BLOCK_VERSION_HIGH = 1;
    public static final short BLOCK_VERSION_LOW = 0;
    public static final MacAddress DEFAULT_EMPTY_MAC_ADDRESS;
    static {
        try {
            DEFAULT_EMPTY_MAC_ADDRESS = new MacAddress(Hex.decodeHex("000000000000"));
        } catch (DecoderException e) {
            throw new RuntimeException(e);
        }
    }

    public static final int DEFAULT_IO_DATA_SIZE = 40;

    private String deviceType;
    private String deviceName;
    private List<String> roles;
    private int vendorId;
    private int deviceId;

    private boolean nonLegacyStartupMode;
    private int sessionKey;
    private int sendClockFactor;
    private int reductionRatio;
    private int watchdogFactor;
    private int dataHoldFactor;

    private String dapId;
    private ProfinetDeviceAccessPointItem dap;
    private Map<Integer, ProfinetModuleItem> moduleIndex;
    private Map<Integer, Map<Integer, ProfinetVirtualSubmoduleItem>> submoduleIndex;

    private int localPort;
    private int remotePortImplicitCommunication;

    // PN-CM Related:
    private final DceRpc_ActivityUuid activityUuid;
    private final Uuid applicationRelationUuid;
    private final AtomicInteger identificationGenerator;

    public ProfinetDriverContext() {
        this.activityUuid = generateActivityUuid();
        this.applicationRelationUuid = generateApplicationRelationUuid();
        this.identificationGenerator = new AtomicInteger(1);
        this.localPort = DEFAULT_UDP_PORT;
        this.remotePortImplicitCommunication = DEFAULT_UDP_PORT;
    }

    /**
     * Generates a new UUID for this connection.
     * @return UUID
     */
    public static DceRpc_ActivityUuid generateActivityUuid() {
        UUID number = UUID.randomUUID();
        try {
            WriteBufferByteBased wb = new WriteBufferByteBased(128);
            wb.writeLong(64, number.getMostSignificantBits());
            wb.writeLong(64, number.getLeastSignificantBits());

            ReadBuffer rb = new ReadBufferByteBased(wb.getBytes());
            return new DceRpc_ActivityUuid(rb.readLong(32), rb.readInt(16), rb.readInt(16), rb.readByteArray(8));
        } catch (SerializationException | ParseException e) {
            // Ignore ... this should actually never happen.
        }
        return null;
    }

    protected static Uuid generateApplicationRelationUuid() {
        try {
            return new Uuid(Hex.decodeHex(UUID.randomUUID().toString().replace("-", "")));
        } catch (DecoderException e) {
            throw new RuntimeException(e);
        }
    }

    public String getDeviceType() {
        return deviceType;
    }

    public void setDeviceType(String deviceType) {
        this.deviceType = deviceType;
    }

    public String getDeviceName() {
        return deviceName;
    }

    public void setDeviceName(String deviceName) {
        this.deviceName = deviceName;
    }

    public List<String> getRoles() {
        return roles;
    }

    public void setRoles(List<String> roles) {
        this.roles = roles;
    }

    public int getVendorId() {
        return vendorId;
    }

    public void setVendorId(int vendorId) {
        this.vendorId = vendorId;
    }

    public int getDeviceId() {
        return deviceId;
    }

    public void setDeviceId(int deviceId) {
        this.deviceId = deviceId;
    }

    public boolean isNonLegacyStartupMode() {
        return nonLegacyStartupMode;
    }

    // TODO: Setup the nonLegacyStartupMode variable.
    public void setNonLegacyStartupMode(boolean nonLegacyStartupMode) {
        this.nonLegacyStartupMode = nonLegacyStartupMode;
    }

    public int getSessionKey() {
        return sessionKey;
    }

    // TODO: Setup the sessionKey variable.
    public void setSessionKey(int sessionKey) {
        this.sessionKey = sessionKey;
    }

    public int getSendClockFactor() {
        return sendClockFactor;
    }

    public void setSendClockFactor(int sendClockFactor) {
        this.sendClockFactor = sendClockFactor;
    }

    public int getReductionRatio() {
        return reductionRatio;
    }

    public void setReductionRatio(int reductionRatio) {
        this.reductionRatio = reductionRatio;
    }

    public int getWatchdogFactor() {
        return watchdogFactor;
    }

    public void setWatchdogFactor(int watchdogFactor) {
        this.watchdogFactor = watchdogFactor;
    }

    public int getDataHoldFactor() {
        return dataHoldFactor;
    }

    public void setDataHoldFactor(int dataHoldFactor) {
        this.dataHoldFactor = dataHoldFactor;
    }

    public String getDapId() {
        return dapId;
    }

    public void setDapId(String dapId) {
        this.dapId = dapId;
    }

    public ProfinetDeviceAccessPointItem getDap() {
        return dap;
    }

    public void setDap(ProfinetDeviceAccessPointItem dap) {
        this.dap = dap;
    }

    public Map<Integer, ProfinetModuleItem> getModuleIndex() {
        return moduleIndex;
    }

    public void setModuleIndex(Map<Integer, ProfinetModuleItem> moduleIndex) {
        this.moduleIndex = moduleIndex;
    }

    public Map<Integer, Map<Integer, ProfinetVirtualSubmoduleItem>> getSubmoduleIndex() {
        return submoduleIndex;
    }

    public void setSubmoduleIndex(Map<Integer, Map<Integer, ProfinetVirtualSubmoduleItem>> submoduleIndex) {
        this.submoduleIndex = submoduleIndex;
    }

    public DceRpc_ActivityUuid getActivityUuid() {
        return activityUuid;
    }

    public Uuid getApplicationRelationUuid() {
        return applicationRelationUuid;
    }

    public int getLocalPort() {
        return localPort;
    }

    public int getRemotePortImplicitCommunication() {
        return remotePortImplicitCommunication;
    }

    public int getAndIncrementIdentification() {
        int id = identificationGenerator.getAndIncrement();
        if (id == 0xFFFF) {
            identificationGenerator.set(1);
        }
        return id;
    }

    public DceRpc_ObjectUuid getCmInitiatorObjectUuid() {
        return new DceRpc_ObjectUuid((byte) 0x00, 0x0001,
            Integer.decode("0x" + getDeviceId()),
            Integer.decode("0x" + getVendorId()));
    }

    public Uuid generateUuid() {
        try {
            return new Uuid(Hex.decodeHex(UUID.randomUUID().toString().replace("-", "")));
        } catch (DecoderException e) {
            throw new RuntimeException(e);
        }
    }

}
