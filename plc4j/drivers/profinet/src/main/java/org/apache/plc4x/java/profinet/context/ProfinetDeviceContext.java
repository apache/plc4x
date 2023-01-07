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
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.profinet.config.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.device.ProfinetCallable;
import org.apache.plc4x.java.profinet.device.ProfinetChannel;
import org.apache.plc4x.java.profinet.device.ProfinetSubscriptionHandle;
import org.apache.plc4x.java.profinet.gsdml.*;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.generation.*;

import java.net.DatagramSocket;
import java.net.InetAddress;
import java.util.*;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class ProfinetDeviceContext implements DriverContext, HasConfiguration<ProfinetConfiguration> {

    public static final int DEFAULT_UDP_PORT = 34964;
    public static final int DEFAULT_ARGS_MAXIMUM = 16696;
    public static final int DEFAULT_MAX_ARRAY_COUNT = 16696;
    public static final int DEFAULT_ACTIVITY_TIMEOUT = 600;
    public static final int UDP_RT_PORT = 0x8892;
    public static final short BLOCK_VERSION_HIGH = 1;
    public static final short BLOCK_VERSION_LOW = 0;
    public static final MacAddress DEFAULT_EMPTY_MAC_ADDRESS;
    public static final Pattern RANGE_PATTERN = Pattern.compile("(?<from>\\d+)\\.\\.(?<to>\\d+)");

    static {
        try {
            DEFAULT_EMPTY_MAC_ADDRESS = new MacAddress(Hex.decodeHex("000000000000"));
        } catch (DecoderException e) {
            throw new RuntimeException(e);
        }
    }
    public static final Uuid ARUUID;
    static {
        try {
            ARUUID = new Uuid(Hex.decodeHex(UUID.randomUUID().toString().replace("-", "")));
        } catch (DecoderException e) {
            throw new RuntimeException(e);
        }
    }
    public static final int DEFAULT_SEND_UDP_PORT = 50000;
    public static final String DEFAULT_PLC4X_STATION_NAME = "plc4x";
    public static final int DEFAULT_IO_DATA_SIZE = 40;

    private DceRpc_ActivityUuid dceRpc_activityUuid;
    private MacAddress localMacAddress;
    private DceRpc_ActivityUuid uuid;
    private ProfinetConfiguration configuration;
    private InetAddress localIpAddress;
    private DatagramSocket socket;
    private ProfinetChannel channel;
    private MacAddress macAddress;
    private ConversationContext<Ethernet_Frame> context;
    private ProfinetDeviceState state = ProfinetDeviceState.IDLE;
    private boolean lldpReceived = false;
    private boolean dcpReceived = false;
    private String ipAddress;
    private String portId;
    PnIoCm_Block_IoCrReq inputReq = null;
    PnIoCm_Block_IoCrReq outputReq = null;
    private String[] subModules;
    private AtomicInteger sessionKeyGenerator = new AtomicInteger(1);
    private AtomicInteger identificationGenerator = new AtomicInteger(1);
    List<PnIoCm_IoDataObject> inputIoDataApiBlocks = new ArrayList<>();
    List<PnIoCm_IoCs> inputIoCsApiBlocks = new ArrayList<>();
    List<PnIoCm_IoDataObject> outputIoDataApiBlocks = new ArrayList<>();
    List<PnIoCm_IoCs> outputIoCsApiBlocks = new ArrayList<>();
    List<PnIoCm_Submodule> expectedSubModuleApiBlocks = new ArrayList<>();
    List<PnIoCm_Block_ExpectedSubmoduleReq> expectedSubmoduleReq = new ArrayList<>();
    private String deviceTypeName;
    private String vendorId;
    private String deviceId;
    private String deviceName;
    private ProfinetISO15745Profile gsdFile;
    private boolean startupMode = false;
    private int frameId = 0xBBF0;
    private Map<Long, ProfinetCallable<DceRpc_Packet>> queue = new HashMap<>();
    private int sessionKey;
    private int sourcePort = DEFAULT_SEND_UDP_PORT;
    private int destinationPort = DEFAULT_UDP_PORT;
    private ProfinetSubscriptionHandle subscriptionHandle;

    private String deviceAccess;
    private ProfinetDeviceAccessPointItem deviceAccessItem;

    public ProfinetDeviceContext() {
        // Generate a new Activity Id, which will be used throughout the connection.
        this.uuid = generateActivityUuid();
    }

    protected static DceRpc_ActivityUuid generateActivityUuid() {
        UUID number = UUID.randomUUID();
        try {
            WriteBufferByteBased wb = new WriteBufferByteBased(128);
            wb.writeLong(64, number.getMostSignificantBits());
            wb.writeLong(64, number.getLeastSignificantBits());

            ReadBuffer rb = new ReadBufferByteBased(wb.getData());
            return new DceRpc_ActivityUuid(rb.readLong(32), rb.readInt(16), rb.readInt(16), rb.readByteArray(8));
        } catch (SerializationException | ParseException e) {
            // Ignore ... this should actually never happen.
        }
        return null;
    }

    public int getAndIncrementIdentification() {
        int id = getIdentificationGenerator().getAndIncrement();
        if (id >= 65535) {
            getIdentificationGenerator().set(1);
        }
        return id;
    }

    public int getIncrementAndGetFrameId() {
        frameId += 1;
        return frameId;
    }

    @Override
    public void setConfiguration(ProfinetConfiguration configuration) {
        this.configuration = configuration;
    }

    public DceRpc_ActivityUuid getDceRpc_activityUuid() {
        return dceRpc_activityUuid;
    }

    public void setDceRpc_activityUuid(DceRpc_ActivityUuid dceRpc_activityUuid) {
        this.dceRpc_activityUuid = dceRpc_activityUuid;
    }

    public MacAddress getLocalMacAddress() {
        return localMacAddress;
    }

    public void setLocalMacAddress(MacAddress localMacAddress) {
        this.localMacAddress = localMacAddress;
    }

    public DceRpc_ActivityUuid getUuid() {
        return uuid;
    }

    public void setUuid(DceRpc_ActivityUuid uuid) {
        this.uuid = uuid;
    }

    public ProfinetConfiguration getConfiguration() {
        return configuration;
    }

    public InetAddress getLocalIpAddress() {
        return localIpAddress;
    }

    public void setLocalIpAddress(InetAddress localIpAddress) {
        this.localIpAddress = localIpAddress;
    }

    public DatagramSocket getSocket() {
        return socket;
    }

    public void setSocket(DatagramSocket socket) {
        this.socket = socket;
    }

    public ProfinetChannel getChannel() {
        return channel;
    }

    public void setChannel(ProfinetChannel channel) {
        this.channel = channel;
    }

    public MacAddress getMacAddress() {
        return macAddress;
    }

    public void setMacAddress(MacAddress macAddress) {
        this.macAddress = macAddress;
    }

    public ConversationContext<Ethernet_Frame> getContext() {
        return context;
    }

    public void setContext(ConversationContext<Ethernet_Frame> context) {
        this.context = context;
    }

    public ProfinetDeviceState getState() {
        return state;
    }

    public void setState(ProfinetDeviceState state) {
        this.state = state;
    }

    public boolean isLldpReceived() {
        return lldpReceived;
    }

    public void setLldpReceived(boolean lldpReceived) {
        this.lldpReceived = lldpReceived;
    }

    public boolean isDcpReceived() {
        return dcpReceived;
    }

    public void setDcpReceived(boolean dcpReceived) {
        this.dcpReceived = dcpReceived;
    }

    public String getIpAddress() {
        return ipAddress;
    }

    public void setIpAddress(String ipAddress) {
        this.ipAddress = ipAddress;
    }

    public String getPortId() {
        return portId;
    }

    public void setPortId(String portId) {
        this.portId = portId;
    }

    public PnIoCm_Block_IoCrReq getInputReq() {
        return inputReq;
    }

    public void setInputReq(PnIoCm_Block_IoCrReq inputReq) {
        this.inputReq = inputReq;
    }

    public PnIoCm_Block_IoCrReq getOutputReq() {
        return outputReq;
    }

    public void setOutputReq(PnIoCm_Block_IoCrReq outputReq) {
        this.outputReq = outputReq;
    }

    public String[] getSubModules() {
        return subModules;
    }

    public void setSubModules(String subModules) {
        String[] splitModules = subModules.split("[, ]");
        ArrayList<String> arrayList = new ArrayList<>();
        for (String s : splitModules) {
            if (s.length() > 0) {
                arrayList.add(s);
            }
        }
        this.subModules = new String[arrayList.size()];
        arrayList.toArray(this.subModules);
    }

    public AtomicInteger getSessionKeyGenerator() {
        return sessionKeyGenerator;
    }

    public void setSessionKeyGenerator(AtomicInteger sessionKeyGenerator) {
        this.sessionKeyGenerator = sessionKeyGenerator;
    }

    public AtomicInteger getIdentificationGenerator() {
        return identificationGenerator;
    }

    public void setIdentificationGenerator(AtomicInteger identificationGenerator) {
        this.identificationGenerator = identificationGenerator;
    }

    public List<PnIoCm_IoDataObject> getInputIoDataApiBlocks() {
        return inputIoDataApiBlocks;
    }

    public void setInputIoDataApiBlocks(List<PnIoCm_IoDataObject> inputIoDataApiBlocks) {
        this.inputIoDataApiBlocks = inputIoDataApiBlocks;
    }

    public List<PnIoCm_IoCs> getInputIoCsApiBlocks() {
        return inputIoCsApiBlocks;
    }

    public void setInputIoCsApiBlocks(List<PnIoCm_IoCs> inputIoCsApiBlocks) {
        this.inputIoCsApiBlocks = inputIoCsApiBlocks;
    }

    public List<PnIoCm_IoDataObject> getOutputIoDataApiBlocks() {
        return outputIoDataApiBlocks;
    }

    public void setOutputIoDataApiBlocks(List<PnIoCm_IoDataObject> outputIoDataApiBlocks) {
        this.outputIoDataApiBlocks = outputIoDataApiBlocks;
    }

    public List<PnIoCm_IoCs> getOutputIoCsApiBlocks() {
        return outputIoCsApiBlocks;
    }

    public void setOutputIoCsApiBlocks(List<PnIoCm_IoCs> outputIoCsApiBlocks) {
        this.outputIoCsApiBlocks = outputIoCsApiBlocks;
    }

    public List<PnIoCm_Submodule> getExpectedSubModuleApiBlocks() {
        return expectedSubModuleApiBlocks;
    }

    public void setExpectedSubModuleApiBlocks(List<PnIoCm_Submodule> expectedSubModuleApiBlocks) {
        this.expectedSubModuleApiBlocks = expectedSubModuleApiBlocks;
    }

    public List<PnIoCm_Block_ExpectedSubmoduleReq> getExpectedSubmoduleReq() {
        return expectedSubmoduleReq;
    }

    public void setExpectedSubmoduleReq(List<PnIoCm_Block_ExpectedSubmoduleReq> expectedSubmoduleReq) {
        this.expectedSubmoduleReq = expectedSubmoduleReq;
    }

    public String getDeviceTypeName() {
        return deviceTypeName;
    }

    public void setDeviceTypeName(String deviceTypeName) {
        this.deviceTypeName = deviceTypeName;
    }

    public String getVendorId() {
        return vendorId;
    }

    public void setVendorId(String vendorId) {
        this.vendorId = vendorId;
    }

    public String getDeviceId() {
        return deviceId;
    }

    public void setDeviceId(String deviceId) {
        this.deviceId = deviceId;
    }

    public String getDeviceName() {
        return deviceName;
    }

    public void setDeviceName(String deviceName) {
        this.deviceName = deviceName;
    }

    public ProfinetISO15745Profile getGsdFile() {
        return gsdFile;
    }

    public void setGsdFile(ProfinetISO15745Profile gsdFile) throws PlcConnectionException {
        this.gsdFile = gsdFile;
        extractGSDFileInfo(this.gsdFile);
    }

    private void extractGSDFileInfo(ProfinetISO15745Profile gsdFile) throws PlcConnectionException {
        ProfinetDeviceAccessPointItem foundDeviceAccessItem = null;
        for (ProfinetDeviceAccessPointItem deviceAccessItem : gsdFile.getProfileBody().getApplicationProcess().getDeviceAccessPointList()) {
            if (deviceAccess.equals(deviceAccessItem.getId())) {
                foundDeviceAccessItem = deviceAccessItem;
            }
        }
        if (foundDeviceAccessItem == null) {
            throw new PlcConnectionException("Unable to find Device Access Item - " + this.deviceAccess);
        }

        this.deviceAccessItem = foundDeviceAccessItem;

        List<ProfinetModuleItemRef> usableSubModules = this.deviceAccessItem.getUseableModules();
        for (String subModule : this.subModules) {
            boolean found = false;
            for (ProfinetModuleItemRef useableModule : usableSubModules) {
                if (useableModule.getModuleItemTarget().equals(subModule)) {
                    found = true;
                    break;
                }
            }
            if (!found) {
                throw new PlcConnectionException("Sub Module not Found in allowed Modules");
            }
        }

        List<ProfinetInterfaceSubmoduleItem> interfaceSubModules = foundDeviceAccessItem.getSystemDefinedSubmoduleList().getInterfaceSubmodules();
        if (interfaceSubModules != null && interfaceSubModules.size() > 0) {
            if (interfaceSubModules.get(0).getApplicationRelations().getStartupMode() != null && interfaceSubModules.get(0).getApplicationRelations().getStartupMode().toLowerCase().contains("advanced")) {
                this.startupMode = true;
                this.frameId = 0x8001;
            }
        }
    }

    public boolean isStartupMode() {
        return startupMode;
    }

    public void setStartupMode(boolean startupMode) {
        this.startupMode = startupMode;
    }

    public int getFrameId() {
        return frameId;
    }

    public void setFrameId(int FrameId) {
        this.frameId = frameId;
    }

    public Map<Long, ProfinetCallable<DceRpc_Packet>> getQueue() {
        return queue;
    }

    public void setQueue(Map<Long, ProfinetCallable<DceRpc_Packet>> queue) {
        this.queue = queue;
    }

    public int getSessionKey() {
        return sessionKey;
    }

    public void setSessionKey(int sessionKey) {
        this.sessionKey = sessionKey;
    }

    public int getSourcePort() {
        return sourcePort;
    }

    public void setSourcePort(int port) {
        this.sourcePort = port;
    }

    public int getDestinationPort() {
        return destinationPort;
    }

    public void setDestinationPort(int port) {
        this.destinationPort = port;
    }

    public ProfinetSubscriptionHandle getSubscriptionHandle() {
        return subscriptionHandle;
    }

    public void setSubscriptionHandle(ProfinetSubscriptionHandle subscriptionHandle) {
        this.subscriptionHandle = subscriptionHandle;
    }

    public String getDeviceAccess() {
        return deviceAccess;
    }

    public void setDeviceAccess(String deviceAccess) {
        this.deviceAccess = deviceAccess;
    }

    public void populateNode() throws PlcException {
        extractGSDFileInfo(this.gsdFile);
        Matcher matcher = RANGE_PATTERN.matcher(deviceAccessItem.getPhysicalSlots());
        if (!matcher.matches()) {
            throw new PlcConnectionException("Physical Slots Range is not in the correct format " + deviceAccessItem.getPhysicalSlots());
        }

        int inputIoDataOffsetCount = 0;
        int outputIoCsOffsetCount = 0;

        for (ProfinetVirtualSubmoduleItem virtualItem : this.deviceAccessItem.getVirtualSubmoduleList()) {
            Integer identNumber = Integer.decode(virtualItem.getSubmoduleIdentNumber());
            inputIoDataApiBlocks.add(new PnIoCm_IoDataObject(
                0,
                identNumber,
                inputIoDataOffsetCount));
            outputIoCsApiBlocks.add(new PnIoCm_IoCs(
                0,
                identNumber,
                outputIoCsOffsetCount));
            expectedSubModuleApiBlocks.add(new PnIoCm_Submodule_NoInputNoOutputData(
                identNumber,
                identNumber,
                false,
                false,
                false,
                false));
            inputIoDataOffsetCount += 1;
            outputIoCsOffsetCount += 1;
        }

        for (ProfinetInterfaceSubmoduleItem interfaceItem : this.deviceAccessItem.getSystemDefinedSubmoduleList().getInterfaceSubmodules()) {
            Integer identNumber = Integer.decode(interfaceItem.getSubmoduleIdentNumber());
            inputIoDataApiBlocks.add(new PnIoCm_IoDataObject(
                0,
                identNumber,
                inputIoDataOffsetCount));
            outputIoCsApiBlocks.add(new PnIoCm_IoCs(
                0,
                identNumber,
                outputIoCsOffsetCount));
            expectedSubModuleApiBlocks.add(new PnIoCm_Submodule_NoInputNoOutputData(
                identNumber,
                identNumber,
                false,
                false,
                false,
                false));
            inputIoDataOffsetCount += 1;
            outputIoCsOffsetCount += 1;
        }

        for (
            ProfinetPortSubmoduleItem portItem : this.deviceAccessItem.getSystemDefinedSubmoduleList().getPortSubmodules()) {
            Integer identNumber = Integer.decode(portItem.getSubmoduleIdentNumber());
            inputIoDataApiBlocks.add(new PnIoCm_IoDataObject(
                0,
                identNumber,
                inputIoDataOffsetCount));
            outputIoCsApiBlocks.add(new PnIoCm_IoCs(
                0,
                identNumber,
                outputIoCsOffsetCount));
            expectedSubModuleApiBlocks.add(new PnIoCm_Submodule_NoInputNoOutputData(
                identNumber,
                identNumber,
                false,
                false,
                false,
                false));
            inputIoDataOffsetCount += 1;
            outputIoCsOffsetCount += 1;
        }
        expectedSubmoduleReq.add(
            new PnIoCm_Block_ExpectedSubmoduleReq((short) 1, (short) 0,
                Collections.singletonList(
                    new PnIoCm_ExpectedSubmoduleBlockReqApi(0,
                        0x00000001,
                        0x00000000,
                        expectedSubModuleApiBlocks
                    )
                )
            )
        );

        int slot = 1;
        for (String submodule : subModules) {
            ProfinetModuleItem foundModule = null;
            for (ProfinetModuleItem module : gsdFile.getProfileBody().getApplicationProcess().getModuleList()) {
                if (module.getId().equals(submodule)) {
                    foundModule = module;
                    break;
                }
            }
            if (foundModule == null) {
                throw new PlcException("Unable to find module id in configured devices");
            }

            Integer identNumber = Integer.decode(foundModule.getModuleIdentNumber());
            if (foundModule.getInputDataLength() != 0) {
                inputIoDataApiBlocks.add(new PnIoCm_IoDataObject(
                    slot,
                    0x01,
                    inputIoDataOffsetCount));
                inputIoDataOffsetCount += 1 + foundModule.getInputDataLength();
            }
            if (foundModule.getInputDataLength() != 0) {
                outputIoCsApiBlocks.add(new PnIoCm_IoCs(
                    slot,
                    0x01,
                    outputIoCsOffsetCount));
                outputIoCsOffsetCount += 1;
            }
            slot += 1;
        }
        slot = 1;
        for (String submodule : this.subModules) {
            ProfinetModuleItem foundModule = null;
            for (ProfinetModuleItem module : gsdFile.getProfileBody().getApplicationProcess().getModuleList()) {
                if (module.getId().equals(submodule)) {
                    foundModule = module;
                    break;
                }
            }

            Integer identNumber = Integer.decode(foundModule.getModuleIdentNumber());
            if (foundModule.getOutputDataLength() != 0) {
                inputIoCsApiBlocks.add(new PnIoCm_IoCs(
                    slot,
                    0x01,
                    inputIoDataOffsetCount));
                inputIoDataOffsetCount += foundModule.getOutputDataLength();
            }

            if (foundModule.getOutputDataLength() != 0) {
                outputIoDataApiBlocks.add(new PnIoCm_IoDataObject(
                    slot,
                    0x01,
                    outputIoCsOffsetCount));
                outputIoCsOffsetCount += 1 + foundModule.getOutputDataLength();
            }

            if (foundModule.getInputDataLength() != 0 && foundModule.getOutputDataLength() != 0) {
                expectedSubmoduleReq.add(new PnIoCm_Block_ExpectedSubmoduleReq((short) 1, (short) 0,
                    Collections.singletonList(
                        new PnIoCm_ExpectedSubmoduleBlockReqApi(slot,
                            identNumber,
                            0x00000000,
                            Collections.singletonList(new PnIoCm_Submodule_InputAndOutputData(
                                0x01,
                                Long.decode(foundModule.getVirtualSubmoduleList().get(0).getSubmoduleIdentNumber()),
                                false,
                                false,
                                false,
                                false,
                                foundModule.getInputDataLength(),
                                (short) 0x01,
                                (short) 0x01,
                                foundModule.getOutputDataLength(),
                                (short) 0x01,
                                (short) 0x01
                            ))
                        )
                    )));
            } else if (foundModule.getInputDataLength() != 0) {
                expectedSubmoduleReq.add(new PnIoCm_Block_ExpectedSubmoduleReq((short) 1, (short) 0,
                    Collections.singletonList(
                        new PnIoCm_ExpectedSubmoduleBlockReqApi(slot,
                            identNumber,
                            0x00000000,
                            Collections.singletonList(new PnIoCm_Submodule_InputData(
                                0x01,
                                Long.decode(foundModule.getVirtualSubmoduleList().get(0).getSubmoduleIdentNumber()),
                                false,
                                false,
                                false,
                                false,
                                foundModule.getInputDataLength(),
                                (short) 0x01,
                                (short) 0x01))
                        )
                    )));
            } else if (foundModule.getOutputDataLength() != 0) {
                expectedSubmoduleReq.add(new PnIoCm_Block_ExpectedSubmoduleReq((short) 1, (short) 0,
                    Collections.singletonList(
                        new PnIoCm_ExpectedSubmoduleBlockReqApi(slot,
                            identNumber,
                            0x00000000,
                            Collections.singletonList(new PnIoCm_Submodule_OutputData(
                                0x01,
                                Long.decode(foundModule.getVirtualSubmoduleList().get(0).getSubmoduleIdentNumber()),
                                false,
                                false,
                                false,
                                false,
                                foundModule.getOutputDataLength(),
                                (short) 0x01,
                                (short) 0x01))
                        )
                    )));
            }
            slot += 1;
        }
    }
}
