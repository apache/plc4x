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
import org.apache.plc4x.java.profinet.gsdml.ProfinetISO15745Profile;
import org.apache.plc4x.java.profinet.readwrite.DceRpc_ActivityUuid;
import org.apache.plc4x.java.profinet.readwrite.DceRpc_Packet;
import org.apache.plc4x.java.profinet.readwrite.Uuid;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.generation.*;

import java.net.DatagramSocket;
import java.net.SocketException;
import java.util.List;
import java.util.UUID;
import java.util.concurrent.atomic.AtomicInteger;

public class ProfinetDriverContext implements DriverContext {

    public static final int DEFAULT_UDP_PORT = 34964;

    private String deviceType;
    private String deviceName;
    private List<String> roles;
    private int vendorId;
    private int deviceId;

    private String dapId;

    private int localPort;
    private int remotePortImplicitCommunication;

    private ProfinetISO15745Profile deviceProfile;

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
    protected static DceRpc_ActivityUuid generateActivityUuid() {
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

    public String getDapId() {
        return dapId;
    }

    public void setDapId(String dapId) {
        this.dapId = dapId;
    }

    public ProfinetISO15745Profile getDeviceProfile() {
        return deviceProfile;
    }

    public void setDeviceProfile(ProfinetISO15745Profile deviceProfile) {
        this.deviceProfile = deviceProfile;
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

}
