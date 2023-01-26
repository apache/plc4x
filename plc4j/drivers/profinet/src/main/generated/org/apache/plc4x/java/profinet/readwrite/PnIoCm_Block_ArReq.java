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
package org.apache.plc4x.java.profinet.readwrite;

import static org.apache.plc4x.java.spi.codegen.fields.FieldReaderFactory.*;
import static org.apache.plc4x.java.spi.codegen.fields.FieldWriterFactory.*;
import static org.apache.plc4x.java.spi.codegen.io.DataReaderFactory.*;
import static org.apache.plc4x.java.spi.codegen.io.DataWriterFactory.*;
import static org.apache.plc4x.java.spi.generation.StaticHelper.*;

import java.time.*;
import java.util.*;
import org.apache.plc4x.java.api.exceptions.*;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.spi.codegen.*;
import org.apache.plc4x.java.spi.codegen.fields.*;
import org.apache.plc4x.java.spi.codegen.io.*;
import org.apache.plc4x.java.spi.generation.*;

// Code generated by code-generation. DO NOT EDIT.

public class PnIoCm_Block_ArReq extends PnIoCm_Block implements Message {

  // Accessors for discriminator values.
  public PnIoCm_BlockType getBlockType() {
    return PnIoCm_BlockType.AR_BLOCK_REQ;
  }

  // Properties.
  protected final short blockVersionHigh;
  protected final short blockVersionLow;
  protected final PnIoCm_ArType arType;
  protected final Uuid arUuid;
  protected final int sessionKey;
  protected final MacAddress cmInitiatorMacAddr;
  protected final DceRpc_ObjectUuid cmInitiatorObjectUuid;
  protected final boolean pullModuleAlarmAllowed;
  protected final boolean nonLegacyStartupMode;
  protected final boolean combinedObjectContainerUsed;
  protected final boolean acknowledgeCompanionAr;
  protected final PnIoCm_CompanionArType companionArType;
  protected final boolean deviceAccess;
  protected final boolean cmInitiator;
  protected final boolean supervisorTakeoverAllowed;
  protected final PnIoCm_State state;
  protected final int cmInitiatorActivityTimeoutFactor;
  protected final int cmInitiatorUdpRtPort;
  protected final String cmInitiatorStationName;
  // Reserved Fields
  private Long reservedField0;
  private Byte reservedField1;

  public PnIoCm_Block_ArReq(
      short blockVersionHigh,
      short blockVersionLow,
      PnIoCm_ArType arType,
      Uuid arUuid,
      int sessionKey,
      MacAddress cmInitiatorMacAddr,
      DceRpc_ObjectUuid cmInitiatorObjectUuid,
      boolean pullModuleAlarmAllowed,
      boolean nonLegacyStartupMode,
      boolean combinedObjectContainerUsed,
      boolean acknowledgeCompanionAr,
      PnIoCm_CompanionArType companionArType,
      boolean deviceAccess,
      boolean cmInitiator,
      boolean supervisorTakeoverAllowed,
      PnIoCm_State state,
      int cmInitiatorActivityTimeoutFactor,
      int cmInitiatorUdpRtPort,
      String cmInitiatorStationName) {
    super();
    this.blockVersionHigh = blockVersionHigh;
    this.blockVersionLow = blockVersionLow;
    this.arType = arType;
    this.arUuid = arUuid;
    this.sessionKey = sessionKey;
    this.cmInitiatorMacAddr = cmInitiatorMacAddr;
    this.cmInitiatorObjectUuid = cmInitiatorObjectUuid;
    this.pullModuleAlarmAllowed = pullModuleAlarmAllowed;
    this.nonLegacyStartupMode = nonLegacyStartupMode;
    this.combinedObjectContainerUsed = combinedObjectContainerUsed;
    this.acknowledgeCompanionAr = acknowledgeCompanionAr;
    this.companionArType = companionArType;
    this.deviceAccess = deviceAccess;
    this.cmInitiator = cmInitiator;
    this.supervisorTakeoverAllowed = supervisorTakeoverAllowed;
    this.state = state;
    this.cmInitiatorActivityTimeoutFactor = cmInitiatorActivityTimeoutFactor;
    this.cmInitiatorUdpRtPort = cmInitiatorUdpRtPort;
    this.cmInitiatorStationName = cmInitiatorStationName;
  }

  public short getBlockVersionHigh() {
    return blockVersionHigh;
  }

  public short getBlockVersionLow() {
    return blockVersionLow;
  }

  public PnIoCm_ArType getArType() {
    return arType;
  }

  public Uuid getArUuid() {
    return arUuid;
  }

  public int getSessionKey() {
    return sessionKey;
  }

  public MacAddress getCmInitiatorMacAddr() {
    return cmInitiatorMacAddr;
  }

  public DceRpc_ObjectUuid getCmInitiatorObjectUuid() {
    return cmInitiatorObjectUuid;
  }

  public boolean getPullModuleAlarmAllowed() {
    return pullModuleAlarmAllowed;
  }

  public boolean getNonLegacyStartupMode() {
    return nonLegacyStartupMode;
  }

  public boolean getCombinedObjectContainerUsed() {
    return combinedObjectContainerUsed;
  }

  public boolean getAcknowledgeCompanionAr() {
    return acknowledgeCompanionAr;
  }

  public PnIoCm_CompanionArType getCompanionArType() {
    return companionArType;
  }

  public boolean getDeviceAccess() {
    return deviceAccess;
  }

  public boolean getCmInitiator() {
    return cmInitiator;
  }

  public boolean getSupervisorTakeoverAllowed() {
    return supervisorTakeoverAllowed;
  }

  public PnIoCm_State getState() {
    return state;
  }

  public int getCmInitiatorActivityTimeoutFactor() {
    return cmInitiatorActivityTimeoutFactor;
  }

  public int getCmInitiatorUdpRtPort() {
    return cmInitiatorUdpRtPort;
  }

  public String getCmInitiatorStationName() {
    return cmInitiatorStationName;
  }

  @Override
  protected void serializePnIoCm_BlockChild(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("PnIoCm_Block_ArReq");

    // Implicit Field (blockLength) (Used for parsing, but its value is not stored as it's
    // implicitly given by the objects content)
    int blockLength = (int) ((getLengthInBytes()) - (4));
    writeImplicitField(
        "blockLength",
        blockLength,
        writeUnsignedInt(writeBuffer, 16),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Simple Field (blockVersionHigh)
    writeSimpleField(
        "blockVersionHigh",
        blockVersionHigh,
        writeUnsignedShort(writeBuffer, 8),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Simple Field (blockVersionLow)
    writeSimpleField(
        "blockVersionLow",
        blockVersionLow,
        writeUnsignedShort(writeBuffer, 8),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Simple Field (arType)
    writeSimpleEnumField(
        "arType",
        "PnIoCm_ArType",
        arType,
        new DataWriterEnumDefault<>(
            PnIoCm_ArType::getValue, PnIoCm_ArType::name, writeUnsignedInt(writeBuffer, 16)),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Simple Field (arUuid)
    writeSimpleField(
        "arUuid",
        arUuid,
        new DataWriterComplexDefault<>(writeBuffer),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Simple Field (sessionKey)
    writeSimpleField(
        "sessionKey",
        sessionKey,
        writeUnsignedInt(writeBuffer, 16),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Simple Field (cmInitiatorMacAddr)
    writeSimpleField(
        "cmInitiatorMacAddr",
        cmInitiatorMacAddr,
        new DataWriterComplexDefault<>(writeBuffer),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Simple Field (cmInitiatorObjectUuid)
    writeSimpleField(
        "cmInitiatorObjectUuid",
        cmInitiatorObjectUuid,
        new DataWriterComplexDefault<>(writeBuffer),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Simple Field (pullModuleAlarmAllowed)
    writeSimpleField(
        "pullModuleAlarmAllowed",
        pullModuleAlarmAllowed,
        writeBoolean(writeBuffer),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Simple Field (nonLegacyStartupMode)
    writeSimpleField(
        "nonLegacyStartupMode",
        nonLegacyStartupMode,
        writeBoolean(writeBuffer),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Simple Field (combinedObjectContainerUsed)
    writeSimpleField(
        "combinedObjectContainerUsed",
        combinedObjectContainerUsed,
        writeBoolean(writeBuffer),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Reserved Field (reserved)
    writeReservedField(
        "reserved",
        reservedField0 != null ? reservedField0 : (long) 0x00000,
        writeUnsignedLong(writeBuffer, 17),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Simple Field (acknowledgeCompanionAr)
    writeSimpleField(
        "acknowledgeCompanionAr",
        acknowledgeCompanionAr,
        writeBoolean(writeBuffer),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Simple Field (companionArType)
    writeSimpleEnumField(
        "companionArType",
        "PnIoCm_CompanionArType",
        companionArType,
        new DataWriterEnumDefault<>(
            PnIoCm_CompanionArType::getValue,
            PnIoCm_CompanionArType::name,
            writeUnsignedByte(writeBuffer, 2)),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Simple Field (deviceAccess)
    writeSimpleField(
        "deviceAccess",
        deviceAccess,
        writeBoolean(writeBuffer),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Reserved Field (reserved)
    writeReservedField(
        "reserved",
        reservedField1 != null ? reservedField1 : (byte) 0x0,
        writeUnsignedByte(writeBuffer, 3),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Simple Field (cmInitiator)
    writeSimpleField(
        "cmInitiator",
        cmInitiator,
        writeBoolean(writeBuffer),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Simple Field (supervisorTakeoverAllowed)
    writeSimpleField(
        "supervisorTakeoverAllowed",
        supervisorTakeoverAllowed,
        writeBoolean(writeBuffer),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Simple Field (state)
    writeSimpleEnumField(
        "state",
        "PnIoCm_State",
        state,
        new DataWriterEnumDefault<>(
            PnIoCm_State::getValue, PnIoCm_State::name, writeUnsignedByte(writeBuffer, 3)),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Simple Field (cmInitiatorActivityTimeoutFactor)
    writeSimpleField(
        "cmInitiatorActivityTimeoutFactor",
        cmInitiatorActivityTimeoutFactor,
        writeUnsignedInt(writeBuffer, 16),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Simple Field (cmInitiatorUdpRtPort)
    writeSimpleField(
        "cmInitiatorUdpRtPort",
        cmInitiatorUdpRtPort,
        writeUnsignedInt(writeBuffer, 16),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Implicit Field (stationNameLength) (Used for parsing, but its value is not stored as it's
    // implicitly given by the objects content)
    int stationNameLength = (int) (STR_LEN(getCmInitiatorStationName()));
    writeImplicitField(
        "stationNameLength",
        stationNameLength,
        writeUnsignedInt(writeBuffer, 16),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Simple Field (cmInitiatorStationName)
    writeSimpleField(
        "cmInitiatorStationName",
        cmInitiatorStationName,
        writeString(writeBuffer, (stationNameLength) * (8)),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    writeBuffer.popContext("PnIoCm_Block_ArReq");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    PnIoCm_Block_ArReq _value = this;

    // Implicit Field (blockLength)
    lengthInBits += 16;

    // Simple field (blockVersionHigh)
    lengthInBits += 8;

    // Simple field (blockVersionLow)
    lengthInBits += 8;

    // Simple field (arType)
    lengthInBits += 16;

    // Simple field (arUuid)
    lengthInBits += arUuid.getLengthInBits();

    // Simple field (sessionKey)
    lengthInBits += 16;

    // Simple field (cmInitiatorMacAddr)
    lengthInBits += cmInitiatorMacAddr.getLengthInBits();

    // Simple field (cmInitiatorObjectUuid)
    lengthInBits += cmInitiatorObjectUuid.getLengthInBits();

    // Simple field (pullModuleAlarmAllowed)
    lengthInBits += 1;

    // Simple field (nonLegacyStartupMode)
    lengthInBits += 1;

    // Simple field (combinedObjectContainerUsed)
    lengthInBits += 1;

    // Reserved Field (reserved)
    lengthInBits += 17;

    // Simple field (acknowledgeCompanionAr)
    lengthInBits += 1;

    // Simple field (companionArType)
    lengthInBits += 2;

    // Simple field (deviceAccess)
    lengthInBits += 1;

    // Reserved Field (reserved)
    lengthInBits += 3;

    // Simple field (cmInitiator)
    lengthInBits += 1;

    // Simple field (supervisorTakeoverAllowed)
    lengthInBits += 1;

    // Simple field (state)
    lengthInBits += 3;

    // Simple field (cmInitiatorActivityTimeoutFactor)
    lengthInBits += 16;

    // Simple field (cmInitiatorUdpRtPort)
    lengthInBits += 16;

    // Implicit Field (stationNameLength)
    lengthInBits += 16;

    // Simple field (cmInitiatorStationName)
    lengthInBits += (STR_LEN(getCmInitiatorStationName())) * (8);

    return lengthInBits;
  }

  public static PnIoCm_Block_ArReqBuilder staticParseBuilder(ReadBuffer readBuffer)
      throws ParseException {
    readBuffer.pullContext("PnIoCm_Block_ArReq");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    int blockLength =
        readImplicitField(
            "blockLength",
            readUnsignedInt(readBuffer, 16),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    short blockVersionHigh =
        readSimpleField(
            "blockVersionHigh",
            readUnsignedShort(readBuffer, 8),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    short blockVersionLow =
        readSimpleField(
            "blockVersionLow",
            readUnsignedShort(readBuffer, 8),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    PnIoCm_ArType arType =
        readEnumField(
            "arType",
            "PnIoCm_ArType",
            new DataReaderEnumDefault<>(
                PnIoCm_ArType::enumForValue, readUnsignedInt(readBuffer, 16)),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    Uuid arUuid =
        readSimpleField(
            "arUuid",
            new DataReaderComplexDefault<>(() -> Uuid.staticParse(readBuffer), readBuffer),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    int sessionKey =
        readSimpleField(
            "sessionKey",
            readUnsignedInt(readBuffer, 16),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    MacAddress cmInitiatorMacAddr =
        readSimpleField(
            "cmInitiatorMacAddr",
            new DataReaderComplexDefault<>(() -> MacAddress.staticParse(readBuffer), readBuffer),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    DceRpc_ObjectUuid cmInitiatorObjectUuid =
        readSimpleField(
            "cmInitiatorObjectUuid",
            new DataReaderComplexDefault<>(
                () -> DceRpc_ObjectUuid.staticParse(readBuffer), readBuffer),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    boolean pullModuleAlarmAllowed =
        readSimpleField(
            "pullModuleAlarmAllowed",
            readBoolean(readBuffer),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    boolean nonLegacyStartupMode =
        readSimpleField(
            "nonLegacyStartupMode",
            readBoolean(readBuffer),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    boolean combinedObjectContainerUsed =
        readSimpleField(
            "combinedObjectContainerUsed",
            readBoolean(readBuffer),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    Long reservedField0 =
        readReservedField(
            "reserved",
            readUnsignedLong(readBuffer, 17),
            (long) 0x00000,
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    boolean acknowledgeCompanionAr =
        readSimpleField(
            "acknowledgeCompanionAr",
            readBoolean(readBuffer),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    PnIoCm_CompanionArType companionArType =
        readEnumField(
            "companionArType",
            "PnIoCm_CompanionArType",
            new DataReaderEnumDefault<>(
                PnIoCm_CompanionArType::enumForValue, readUnsignedByte(readBuffer, 2)),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    boolean deviceAccess =
        readSimpleField(
            "deviceAccess",
            readBoolean(readBuffer),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    Byte reservedField1 =
        readReservedField(
            "reserved",
            readUnsignedByte(readBuffer, 3),
            (byte) 0x0,
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    boolean cmInitiator =
        readSimpleField(
            "cmInitiator", readBoolean(readBuffer), WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    boolean supervisorTakeoverAllowed =
        readSimpleField(
            "supervisorTakeoverAllowed",
            readBoolean(readBuffer),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    PnIoCm_State state =
        readEnumField(
            "state",
            "PnIoCm_State",
            new DataReaderEnumDefault<>(
                PnIoCm_State::enumForValue, readUnsignedByte(readBuffer, 3)),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    int cmInitiatorActivityTimeoutFactor =
        readSimpleField(
            "cmInitiatorActivityTimeoutFactor",
            readUnsignedInt(readBuffer, 16),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    int cmInitiatorUdpRtPort =
        readSimpleField(
            "cmInitiatorUdpRtPort",
            readUnsignedInt(readBuffer, 16),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    int stationNameLength =
        readImplicitField(
            "stationNameLength",
            readUnsignedInt(readBuffer, 16),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    String cmInitiatorStationName =
        readSimpleField(
            "cmInitiatorStationName",
            readString(readBuffer, (stationNameLength) * (8)),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    readBuffer.closeContext("PnIoCm_Block_ArReq");
    // Create the instance
    return new PnIoCm_Block_ArReqBuilder(
        blockVersionHigh,
        blockVersionLow,
        arType,
        arUuid,
        sessionKey,
        cmInitiatorMacAddr,
        cmInitiatorObjectUuid,
        pullModuleAlarmAllowed,
        nonLegacyStartupMode,
        combinedObjectContainerUsed,
        acknowledgeCompanionAr,
        companionArType,
        deviceAccess,
        cmInitiator,
        supervisorTakeoverAllowed,
        state,
        cmInitiatorActivityTimeoutFactor,
        cmInitiatorUdpRtPort,
        cmInitiatorStationName,
        reservedField0,
        reservedField1);
  }

  public static class PnIoCm_Block_ArReqBuilder implements PnIoCm_Block.PnIoCm_BlockBuilder {
    private final short blockVersionHigh;
    private final short blockVersionLow;
    private final PnIoCm_ArType arType;
    private final Uuid arUuid;
    private final int sessionKey;
    private final MacAddress cmInitiatorMacAddr;
    private final DceRpc_ObjectUuid cmInitiatorObjectUuid;
    private final boolean pullModuleAlarmAllowed;
    private final boolean nonLegacyStartupMode;
    private final boolean combinedObjectContainerUsed;
    private final boolean acknowledgeCompanionAr;
    private final PnIoCm_CompanionArType companionArType;
    private final boolean deviceAccess;
    private final boolean cmInitiator;
    private final boolean supervisorTakeoverAllowed;
    private final PnIoCm_State state;
    private final int cmInitiatorActivityTimeoutFactor;
    private final int cmInitiatorUdpRtPort;
    private final String cmInitiatorStationName;
    private final Long reservedField0;
    private final Byte reservedField1;

    public PnIoCm_Block_ArReqBuilder(
        short blockVersionHigh,
        short blockVersionLow,
        PnIoCm_ArType arType,
        Uuid arUuid,
        int sessionKey,
        MacAddress cmInitiatorMacAddr,
        DceRpc_ObjectUuid cmInitiatorObjectUuid,
        boolean pullModuleAlarmAllowed,
        boolean nonLegacyStartupMode,
        boolean combinedObjectContainerUsed,
        boolean acknowledgeCompanionAr,
        PnIoCm_CompanionArType companionArType,
        boolean deviceAccess,
        boolean cmInitiator,
        boolean supervisorTakeoverAllowed,
        PnIoCm_State state,
        int cmInitiatorActivityTimeoutFactor,
        int cmInitiatorUdpRtPort,
        String cmInitiatorStationName,
        Long reservedField0,
        Byte reservedField1) {
      this.blockVersionHigh = blockVersionHigh;
      this.blockVersionLow = blockVersionLow;
      this.arType = arType;
      this.arUuid = arUuid;
      this.sessionKey = sessionKey;
      this.cmInitiatorMacAddr = cmInitiatorMacAddr;
      this.cmInitiatorObjectUuid = cmInitiatorObjectUuid;
      this.pullModuleAlarmAllowed = pullModuleAlarmAllowed;
      this.nonLegacyStartupMode = nonLegacyStartupMode;
      this.combinedObjectContainerUsed = combinedObjectContainerUsed;
      this.acknowledgeCompanionAr = acknowledgeCompanionAr;
      this.companionArType = companionArType;
      this.deviceAccess = deviceAccess;
      this.cmInitiator = cmInitiator;
      this.supervisorTakeoverAllowed = supervisorTakeoverAllowed;
      this.state = state;
      this.cmInitiatorActivityTimeoutFactor = cmInitiatorActivityTimeoutFactor;
      this.cmInitiatorUdpRtPort = cmInitiatorUdpRtPort;
      this.cmInitiatorStationName = cmInitiatorStationName;
      this.reservedField0 = reservedField0;
      this.reservedField1 = reservedField1;
    }

    public PnIoCm_Block_ArReq build() {
      PnIoCm_Block_ArReq pnIoCm_Block_ArReq =
          new PnIoCm_Block_ArReq(
              blockVersionHigh,
              blockVersionLow,
              arType,
              arUuid,
              sessionKey,
              cmInitiatorMacAddr,
              cmInitiatorObjectUuid,
              pullModuleAlarmAllowed,
              nonLegacyStartupMode,
              combinedObjectContainerUsed,
              acknowledgeCompanionAr,
              companionArType,
              deviceAccess,
              cmInitiator,
              supervisorTakeoverAllowed,
              state,
              cmInitiatorActivityTimeoutFactor,
              cmInitiatorUdpRtPort,
              cmInitiatorStationName);
      pnIoCm_Block_ArReq.reservedField0 = reservedField0;
      pnIoCm_Block_ArReq.reservedField1 = reservedField1;
      return pnIoCm_Block_ArReq;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof PnIoCm_Block_ArReq)) {
      return false;
    }
    PnIoCm_Block_ArReq that = (PnIoCm_Block_ArReq) o;
    return (getBlockVersionHigh() == that.getBlockVersionHigh())
        && (getBlockVersionLow() == that.getBlockVersionLow())
        && (getArType() == that.getArType())
        && (getArUuid() == that.getArUuid())
        && (getSessionKey() == that.getSessionKey())
        && (getCmInitiatorMacAddr() == that.getCmInitiatorMacAddr())
        && (getCmInitiatorObjectUuid() == that.getCmInitiatorObjectUuid())
        && (getPullModuleAlarmAllowed() == that.getPullModuleAlarmAllowed())
        && (getNonLegacyStartupMode() == that.getNonLegacyStartupMode())
        && (getCombinedObjectContainerUsed() == that.getCombinedObjectContainerUsed())
        && (getAcknowledgeCompanionAr() == that.getAcknowledgeCompanionAr())
        && (getCompanionArType() == that.getCompanionArType())
        && (getDeviceAccess() == that.getDeviceAccess())
        && (getCmInitiator() == that.getCmInitiator())
        && (getSupervisorTakeoverAllowed() == that.getSupervisorTakeoverAllowed())
        && (getState() == that.getState())
        && (getCmInitiatorActivityTimeoutFactor() == that.getCmInitiatorActivityTimeoutFactor())
        && (getCmInitiatorUdpRtPort() == that.getCmInitiatorUdpRtPort())
        && (getCmInitiatorStationName() == that.getCmInitiatorStationName())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(
        super.hashCode(),
        getBlockVersionHigh(),
        getBlockVersionLow(),
        getArType(),
        getArUuid(),
        getSessionKey(),
        getCmInitiatorMacAddr(),
        getCmInitiatorObjectUuid(),
        getPullModuleAlarmAllowed(),
        getNonLegacyStartupMode(),
        getCombinedObjectContainerUsed(),
        getAcknowledgeCompanionAr(),
        getCompanionArType(),
        getDeviceAccess(),
        getCmInitiator(),
        getSupervisorTakeoverAllowed(),
        getState(),
        getCmInitiatorActivityTimeoutFactor(),
        getCmInitiatorUdpRtPort(),
        getCmInitiatorStationName());
  }

  @Override
  public String toString() {
    WriteBufferBoxBased writeBufferBoxBased = new WriteBufferBoxBased(true, true);
    try {
      writeBufferBoxBased.writeSerializable(this);
    } catch (SerializationException e) {
      throw new RuntimeException(e);
    }
    return "\n" + writeBufferBoxBased.getBox().toString() + "\n";
  }
}
