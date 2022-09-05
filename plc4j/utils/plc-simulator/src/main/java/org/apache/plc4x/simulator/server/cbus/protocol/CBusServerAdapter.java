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
package org.apache.plc4x.simulator.server.cbus.protocol;

import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInboundHandlerAdapter;
import org.apache.plc4x.java.cbus.readwrite.*;
import org.apache.plc4x.simulator.model.Context;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Arrays;
import java.util.Collections;
import java.util.LinkedList;
import java.util.List;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

public class CBusServerAdapter extends ChannelInboundHandlerAdapter {

    private static final List<Byte> AVAILABLE_UNITS = Arrays.asList((byte) 3, (byte) 23, (byte) 48);

    private static final Logger LOGGER = LoggerFactory.getLogger(CBusServerAdapter.class);

    private Context context;

    private static final RequestContext requestContext = new RequestContext(false);

    private static boolean connect;
    private static boolean smart;
    private static boolean idmon;
    private static boolean exstat;
    private static boolean monitor;
    private static boolean monall;
    private static boolean pun;
    private static boolean pcn;
    private static boolean srchk;

    private static byte monitorApplicationAddress1;

    private static byte monitorApplicationAddress2;
    private static CBusOptions cBusOptions;

    private final Lock outputLock = new ReentrantLock();

    private ScheduledFuture<?> salMonitorFuture;

    private ScheduledFuture<?> mmiMonitorFuture;

    public CBusServerAdapter(Context context) {
        LOGGER.info("Creating adapter with context {}", context);
        this.context = context;
        cBusOptions = new CBusOptions(connect, smart, idmon, exstat, monitor, monall, pun, pcn, srchk);
    }

    private static void buildCBusOptions() {
        LOGGER.info("Updating options {}", cBusOptions);
        cBusOptions = new CBusOptions(connect, smart, idmon, exstat, monitor, monall, pun, pcn, srchk);
        LOGGER.info("Updated options {}", cBusOptions);
    }

    @Override
    public void channelInactive(ChannelHandlerContext ctx) throws Exception {
        if (salMonitorFuture != null)
            salMonitorFuture.cancel(false);
        if (mmiMonitorFuture != null)
            mmiMonitorFuture.cancel(false);
        super.channelInactive(ctx);
    }

    @Override
    public void channelRead(ChannelHandlerContext ctx, Object msg) throws Exception {
        if (!(msg instanceof CBusMessage)) {
            return;
        }
        try {
            outputLock.lock();
            syncChannelRead(ctx, msg);
        } finally {
            ctx.flush();
            outputLock.unlock();
        }
    }

    private void syncChannelRead(ChannelHandlerContext ctx, Object msg) throws Exception {
        LOGGER.debug("Working with cBusOptions\n{}", cBusOptions);
        // Serial is slow
        TimeUnit.MILLISECONDS.sleep(100);
        if (!smart && !connect) {
            // In this mode every message will be echoed
            LOGGER.info("Sending echo");
            ctx.write(msg);
        }
        CBusMessage packet = (CBusMessage) msg;
        if (packet instanceof CBusMessageToClient) {
            LOGGER.info("Message to client not supported\n{}", packet);
            return;
        }
        CBusMessageToServer cBusMessageToServer = (CBusMessageToServer) packet;
        Request request = cBusMessageToServer.getRequest();
        if (request instanceof RequestEmpty || request instanceof RequestNull) {
            LOGGER.debug("Ignoring\n{}", request);
        } else if (request instanceof RequestDirectCommandAccess) {
            handleDirectCommandAccess(ctx, (RequestDirectCommandAccess) request);
        } else if (request instanceof RequestCommand) {
            handleRequestCommand(ctx, (RequestCommand) request);
        } else if (request instanceof RequestObsolete) {
            RequestObsolete requestObsolete = (RequestObsolete) request;
            LOGGER.info("Handling RequestObsolete\n{}", requestObsolete);
            // TODO: handle this
        } else if (request instanceof RequestReset) {
            handleReset((RequestReset) request);
        } else if (request instanceof RequestSmartConnectShortcut) {
            handleSmartConnect((RequestSmartConnectShortcut) request);
        }
    }

    private void handleDirectCommandAccess(ChannelHandlerContext ctx, RequestDirectCommandAccess requestDirectCommandAccess) {
        CALData calData = requestDirectCommandAccess.getCalData();
        LOGGER.info("Handling RequestDirectCommandAccess\n{}", requestDirectCommandAccess);

        // TODO: handle other cal data type
        if (calData instanceof CALDataWrite) {
            CALDataWrite calDataWrite = (CALDataWrite) calData;
            Runnable acknowledger = () -> {
                CALDataAcknowledge calDataAcknowledge = new CALDataAcknowledge(CALCommandTypeContainer.CALCommandAcknowledge, null, calDataWrite.getParamNo(), (short) 0x0, requestContext);
                CALReplyShort calReply = new CALReplyShort((byte) 0x0, calDataAcknowledge, cBusOptions, requestContext);
                EncodedReplyCALReply encodedReply = new EncodedReplyCALReply((byte) 0x0, calReply, cBusOptions, requestContext);
                ReplyEncodedReply replyEncodedReply = new ReplyEncodedReply((byte) 0x0, encodedReply, null, cBusOptions, requestContext);
                ReplyOrConfirmation replyOrConfirmation = new ReplyOrConfirmationReply((byte) 0x0, replyEncodedReply, new ResponseTermination(), cBusOptions, requestContext);
                if (requestDirectCommandAccess.getAlpha() != null) {
                    replyOrConfirmation = new ReplyOrConfirmationConfirmation((byte) 0x0, new Confirmation(requestDirectCommandAccess.getAlpha(), null, ConfirmationType.CONFIRMATION_SUCCESSFUL), replyOrConfirmation, cBusOptions, requestContext);
                }
                CBusMessageToClient cBusMessageToClient = new CBusMessageToClient(replyOrConfirmation, requestContext, cBusOptions);
                LOGGER.info("Sending ack\n{}", cBusMessageToClient);
                ctx.writeAndFlush(cBusMessageToClient);
            };
            switch (calDataWrite.getParamNo().getParameterType()) {
                case APPLICATION_ADDRESS_1:
                    ApplicationAddress1 applicationAddress1 = ((ParameterValueApplicationAddress1) calDataWrite.getParameterValue()).getValue();
                    monitorApplicationAddress1 = applicationAddress1.getAddress();
                    acknowledger.run();
                    return;
                case APPLICATION_ADDRESS_2:
                    ApplicationAddress2 applicationAddress2 = ((ParameterValueApplicationAddress2) calDataWrite.getParameterValue()).getValue();
                    monitorApplicationAddress2 = applicationAddress2.getAddress();
                    acknowledger.run();
                    return;
                case INTERFACE_OPTIONS_1:
                    InterfaceOptions1 interfaceOptions1 = ((ParameterValueInterfaceOptions1) calDataWrite.getParameterValue()).getValue();
                    idmon = interfaceOptions1.getIdmon();
                    monitor = interfaceOptions1.getMonitor();
                    if (monitor) startMMIMonitor(ctx);
                    else stopMMIMonitor();
                    smart = interfaceOptions1.getSmart();
                    srchk = interfaceOptions1.getSrchk();
                    // TODO: add support for xonxoff
                    // xonxoff = interfaceOptions1.getXonXoff();
                    connect = interfaceOptions1.getConnect();
                    if (connect) startSALMonitor(ctx);
                    else stopSALMonitor();
                    buildCBusOptions();
                    acknowledger.run();
                    return;
                case INTERFACE_OPTIONS_2:
                    InterfaceOptions2 interfaceOptions2 = ((ParameterValueInterfaceOptions2) calDataWrite.getParameterValue()).getValue();
                    // TODO: add support for burden
                    // burden =  interfaceOptions2.getBurden();
                    // TODO: add support for clockgen
                    // clockgen = interfaceOptions2.getClockGen();
                    buildCBusOptions();
                    acknowledger.run();
                    return;
                case INTERFACE_OPTIONS_3:
                    InterfaceOptions3 interfaceOptions3Value = ((ParameterValueInterfaceOptions3) calDataWrite.getParameterValue()).getValue();
                    exstat = interfaceOptions3Value.getExstat();
                    pun = interfaceOptions3Value.getPun();
                    // TODO: add support for localsal
                    // localsal = interfaceOptions3Value.getLocalSal();
                    pcn = interfaceOptions3Value.getPcn();
                    buildCBusOptions();
                    acknowledger.run();
                    return;
                case BAUD_RATE_SELECTOR:
                    BaudRateSelector baudRateSelector = ((ParameterValueBaudRateSelector) calDataWrite.getParameterValue()).getValue();
                    // TODO: add support for baudrate
                    // baudrate = baudRateSelector.getValue();
                    buildCBusOptions();
                    acknowledger.run();
                    return;
                case INTERFACE_OPTIONS_1_POWER_UP_SETTINGS:
                    InterfaceOptions1 interfaceOptions1PowerUpSettings = ((ParameterValueInterfaceOptions1PowerUpSettings) calDataWrite.getParameterValue()).getValue().getInterfaceOptions1();
                    idmon = interfaceOptions1PowerUpSettings.getIdmon();
                    monitor = interfaceOptions1PowerUpSettings.getMonitor();
                    if (monitor) startMMIMonitor(ctx);
                    else stopMMIMonitor();
                    smart = interfaceOptions1PowerUpSettings.getSmart();
                    srchk = interfaceOptions1PowerUpSettings.getSrchk();
                    // TODO: add support for xonxoff
                    // xonxoff = interfaceOptions1PowerUpSettings.getXonXoff();
                    connect = interfaceOptions1PowerUpSettings.getConnect();
                    if (connect) startSALMonitor(ctx);
                    else stopSALMonitor();
                    buildCBusOptions();
                    acknowledger.run();
                    return;
                case CUSTOM_MANUFACTURER:
                    // TODO: handle other param typed
                    acknowledger.run();
                    return;
                case SERIAL_NUMBER:
                    // TODO: handle other param typed
                    acknowledger.run();
                    return;
                case CUSTOM_TYPE:
                    // TODO: handle other param typed
                    acknowledger.run();
                    return;
                default:
                    throw new IllegalStateException("Unmapped type");
            }
        } else if (calData instanceof CALDataIdentify) {
            handleCalDataIdentify(ctx, (CALDataIdentify) calData, requestDirectCommandAccess.getAlpha());
        }
    }

    private void handleRequestCommand(ChannelHandlerContext ctx, RequestCommand requestCommand) {
        LOGGER.info("Handling RequestCommand\n{}", requestCommand);
        CBusCommand cbusCommand = requestCommand.getCbusCommand();
        LOGGER.info("Handling CBusCommand\n{}", cbusCommand);
        if (cbusCommand instanceof CBusCommandPointToPoint) {
            CBusCommandPointToPoint cBusCommandPointToPoint = (CBusCommandPointToPoint) cbusCommand;
            CBusPointToPointCommand command = cBusCommandPointToPoint.getCommand();
            UnitAddress unitAddress = null;
            if (command instanceof CBusPointToPointCommandIndirect) {
                CBusPointToPointCommandIndirect cBusPointToPointCommandIndirect = (CBusPointToPointCommandIndirect) command;
                // TODO: handle bridgeAddress
                // TODO: handle networkRoute
                unitAddress = cBusPointToPointCommandIndirect.getUnitAddress();
            }
            if (command instanceof CBusPointToPointCommandDirect) {
                CBusPointToPointCommandDirect cBusPointToPointCommandDirect = (CBusPointToPointCommandDirect) command;
                unitAddress = cBusPointToPointCommandDirect.getUnitAddress();
            }
            if (unitAddress == null) {
                throw new IllegalStateException("Unit address should be set at this point");
            }
            boolean knownUnit = AVAILABLE_UNITS.contains(unitAddress.getAddress());
            if (!knownUnit) {
                LOGGER.warn("{} not a known unit", unitAddress);
                ReplyOrConfirmation replyOrConfirmation = new ServerErrorReply((byte) 0x0, cBusOptions, requestContext);
                CBusMessageToClient cBusMessageToClient = new CBusMessageToClient(replyOrConfirmation, requestContext, cBusOptions);
                ctx.writeAndFlush(cBusMessageToClient);
                return;
            }
            CALData calData = command.getCalData();
            // TODO: handle other Datatypes
            if (calData instanceof CALDataIdentify) {
                handleCalDataIdentify(ctx, (CALDataIdentify) calData, requestCommand.getAlpha());
            }
            return;
        } else if (cbusCommand instanceof CBusCommandPointToMultiPoint) {
            CBusCommandPointToMultiPoint cBusCommandPointToMultiPoint = (CBusCommandPointToMultiPoint) cbusCommand;
            CBusPointToMultiPointCommand command = cBusCommandPointToMultiPoint.getCommand();
            if (command instanceof CBusPointToMultiPointCommandStatus) {
                CBusPointToMultiPointCommandStatus cBusPointToMultiPointCommandStatus = (CBusPointToMultiPointCommandStatus) command;
                StatusRequest statusRequest = cBusPointToMultiPointCommandStatus.getStatusRequest();
                if (statusRequest instanceof StatusRequestBinaryState) {
                    StatusRequestBinaryState statusRequestBinaryState = (StatusRequestBinaryState) statusRequest;
                    LOGGER.info("Handling StatusRequestBinaryState\n{}", statusRequestBinaryState);
                    if (statusRequestBinaryState.getApplication() == ApplicationIdContainer.NETWORK_CONTROL) {
                        LOGGER.info("Handling installation MMI Request");
                        sendInstallationMMIResponse(ctx, requestCommand, statusRequestBinaryState.getApplication());
                        return;
                    }
                    CALReply calReply;
                    if (exstat) {
                        // TODO: map actuall values from simulator
                        List<StatusByte> statusBytes = new LinkedList<>();
                        for (int i = 0; i < 22; i++) {
                            statusBytes.add(new StatusByte(GAVState.ON, GAVState.ERROR, GAVState.OFF, GAVState.DOES_NOT_EXIST));
                        }
                        CALData calData = new CALDataStatusExtended(CALCommandTypeContainer.CALCommandStatusExtended_25Bytes, null, StatusCoding.BINARY_BY_ELSEWHERE, ApplicationIdContainer.LIGHTING_38, (byte) 0x00, statusBytes, null, requestContext);
                        calReply = new CALReplyLong((byte) 0x86, calData, 0x00, new UnitAddress((byte) 0x04), null, new SerialInterfaceAddress((byte) 0x02), (byte) 0x00, null, cBusOptions, requestContext);
                    } else {
                        List<StatusByte> statusBytes = new LinkedList<>();
                        // TODO: map actuall values from simulator
                        for (int i = 0; i < 23; i++) {
                            statusBytes.add(new StatusByte(GAVState.ON, GAVState.ERROR, GAVState.OFF, GAVState.DOES_NOT_EXIST));
                        }
                        CALData calData = new CALDataStatus(CALCommandTypeContainer.CALCommandStatus_25Bytes, null, ApplicationIdContainer.LIGHTING_38, (byte) 0x00, statusBytes, requestContext);
                        calReply = new CALReplyShort((byte) 0x0, calData, cBusOptions, requestContext);
                    }
                    EncodedReply encodedReply = new EncodedReplyCALReply((byte) 0x0, calReply, cBusOptions, requestContext);
                    ReplyEncodedReply replyEncodedReply = new ReplyEncodedReply((byte) 0xC0, encodedReply, null, cBusOptions, requestContext);
                    ReplyOrConfirmation replyOrConfirmation = new ReplyOrConfirmationReply((byte) 0xFF, replyEncodedReply, new ResponseTermination(), cBusOptions, requestContext);
                    Alpha alpha = requestCommand.getAlpha();
                    if (alpha != null) {
                        Confirmation confirmation = new Confirmation(alpha, null, ConfirmationType.CONFIRMATION_SUCCESSFUL);
                        replyOrConfirmation = new ReplyOrConfirmationConfirmation(alpha.getCharacter(), confirmation, replyOrConfirmation, cBusOptions, requestContext);
                    }
                    CBusMessage response = new CBusMessageToClient(replyOrConfirmation, requestContext, cBusOptions);
                    LOGGER.info("Send binary status response\n{}", response);
                    ctx.writeAndFlush(response);
                    return;
                }
                if (statusRequest instanceof StatusRequestBinaryStateDeprecated) {
                    StatusRequestBinaryStateDeprecated statusRequestBinaryStateDeprecated = (StatusRequestBinaryStateDeprecated) statusRequest;
                    LOGGER.info("Handling StatusRequestBinaryStateDeprecated\n{}", statusRequestBinaryStateDeprecated);
                    if (statusRequestBinaryStateDeprecated.getApplication() == ApplicationIdContainer.NETWORK_CONTROL) {
                        LOGGER.info("Handling installation MMI Request");
                        sendInstallationMMIResponse(ctx, requestCommand, statusRequestBinaryStateDeprecated.getApplication());
                        return;
                    }
                    CALReply calReply;
                    if (exstat) {
                        // TODO: map actuall values from simulator
                        List<StatusByte> statusBytes = new LinkedList<>();
                        for (int i = 0; i < 22; i++) {
                            statusBytes.add(new StatusByte(GAVState.ON, GAVState.ERROR, GAVState.OFF, GAVState.DOES_NOT_EXIST));
                        }
                        CALData calData = new CALDataStatusExtended(CALCommandTypeContainer.CALCommandStatusExtended_25Bytes, null, StatusCoding.BINARY_BY_ELSEWHERE, ApplicationIdContainer.LIGHTING_38, (byte) 0x00, statusBytes, null, requestContext);
                        calReply = new CALReplyLong((byte) 0x86, calData, 0x00, new UnitAddress((byte) 0x04), null, new SerialInterfaceAddress((byte) 0x02), (byte) 0x00, null, cBusOptions, requestContext);
                    } else {
                        List<StatusByte> statusBytes = new LinkedList<>();
                        // TODO: map actuall values from simulator
                        for (int i = 0; i < 23; i++) {
                            statusBytes.add(new StatusByte(GAVState.ON, GAVState.ERROR, GAVState.OFF, GAVState.DOES_NOT_EXIST));
                        }
                        CALData calData = new CALDataStatus(CALCommandTypeContainer.CALCommandStatus_25Bytes, null, ApplicationIdContainer.LIGHTING_38, (byte) 0x00, statusBytes, requestContext);
                        calReply = new CALReplyShort((byte) 0x0, calData, cBusOptions, requestContext);
                    }
                    EncodedReply encodedReply = new EncodedReplyCALReply((byte) 0x0, calReply, cBusOptions, requestContext);
                    ReplyEncodedReply replyEncodedReply = new ReplyEncodedReply((byte) 0xC0, encodedReply, null, cBusOptions, requestContext);
                    ReplyOrConfirmation replyOrConfirmation = new ReplyOrConfirmationReply((byte) 0xFF, replyEncodedReply, new ResponseTermination(), cBusOptions, requestContext);
                    Alpha alpha = requestCommand.getAlpha();
                    if (alpha != null) {
                        Confirmation confirmation = new Confirmation(alpha, null, ConfirmationType.CONFIRMATION_SUCCESSFUL);
                        replyOrConfirmation = new ReplyOrConfirmationConfirmation(alpha.getCharacter(), confirmation, replyOrConfirmation, cBusOptions, requestContext);
                    }
                    CBusMessage response = new CBusMessageToClient(replyOrConfirmation, requestContext, cBusOptions);
                    LOGGER.info("Send binary status response\n{}", response);
                    ctx.writeAndFlush(response);
                    return;
                }
                if (statusRequest instanceof StatusRequestLevel) {
                    StatusRequestLevel statusRequestLevel = (StatusRequestLevel) statusRequest;
                    StatusCoding coding = StatusCoding.LEVEL_BY_THIS_SERIAL_INTERFACE;
                    // TODO: map actuall values from simulator
                    byte blockStart = statusRequestLevel.getStartingGroupAddressLabel();
                    List<LevelInformation> levelInformations = Collections.singletonList(new LevelInformationNormal(0x5555, LevelInformationNibblePair.Value_F, LevelInformationNibblePair.Value_F));
                    CALData calData = new CALDataStatusExtended(CALCommandTypeContainer.CALCommandReply_4Bytes, null, coding, statusRequestLevel.getApplication(), blockStart, null, levelInformations, requestContext);
                    CALReply calReply = new CALReplyLong((byte) 0x0, calData, (byte) 0x0, new UnitAddress((byte) 0x04), null, new SerialInterfaceAddress((byte) 0x02), (byte) 0x0, null, cBusOptions, requestContext);
                    EncodedReply encodedReply = new EncodedReplyCALReply((byte) 0x0, calReply, cBusOptions, requestContext);
                    ReplyEncodedReply replyEncodedReply = new ReplyEncodedReply((byte) 0xC0, encodedReply, null, cBusOptions, requestContext);
                    ReplyOrConfirmation replyOrConfirmation = new ReplyOrConfirmationReply((byte) 0xFF, replyEncodedReply, new ResponseTermination(), cBusOptions, requestContext);
                    Alpha alpha = requestCommand.getAlpha();
                    if (alpha != null) {
                        Confirmation confirmation = new Confirmation(alpha, null, ConfirmationType.CONFIRMATION_SUCCESSFUL);
                        replyOrConfirmation = new ReplyOrConfirmationConfirmation(alpha.getCharacter(), confirmation, replyOrConfirmation, cBusOptions, requestContext);
                    }
                    CBusMessage response = new CBusMessageToClient(replyOrConfirmation, requestContext, cBusOptions);
                    LOGGER.info("Send level status response\n{}", response);
                    ctx.writeAndFlush(response);
                    return;
                }
                // TODO: handle this
                return;
            }
            if (command instanceof CBusPointToMultiPointCommandNormal) {
                CBusPointToMultiPointCommandNormal cBusPointToMultiPointCommandNormal = (CBusPointToMultiPointCommandNormal) command;
                LOGGER.info("Handling CBusPointToMultiPointCommandNormal\n{}", cBusPointToMultiPointCommandNormal);
                // TODO: handle this
                return;
            }
            // TODO: handle this
            return;
        } else if (cbusCommand instanceof CBusCommandPointToPointToMultiPoint) {
            CBusCommandPointToPointToMultiPoint cBusCommandPointToPointToMultiPoint = (CBusCommandPointToPointToMultiPoint) cbusCommand;
            LOGGER.info("Handling CBusCommandPointToPointToMultiPoint\n{}", cBusCommandPointToPointToMultiPoint);
            // TODO: handle this
            return;
        } else if (cbusCommand instanceof CBusCommandDeviceManagement) {
            CBusCommandDeviceManagement cBusCommandDeviceManagement = (CBusCommandDeviceManagement) cbusCommand;
            LOGGER.info("Handling CBusCommandDeviceManagement\n{}", cBusCommandDeviceManagement);
            // TODO: handle this
            return;
        }

        Alpha alpha = requestCommand.getAlpha();
        if (alpha != null) {
            Confirmation confirmation = new Confirmation(alpha, null, ConfirmationType.NOT_TRANSMITTED_CORRUPTION);
            ReplyOrConfirmationConfirmation replyOrConfirmationConfirmation = new ReplyOrConfirmationConfirmation(alpha.getCharacter(), confirmation, null, cBusOptions, requestContext);
            CBusMessage response = new CBusMessageToClient(replyOrConfirmationConfirmation, requestContext, cBusOptions);
            LOGGER.info("Send response\n{}", response);
            ctx.writeAndFlush(response);
        }
    }

    private static void sendInstallationMMIResponse(ChannelHandlerContext ctx, RequestCommand requestCommand, ApplicationIdContainer application) {
        LOGGER.info("Send installation mmis");
        {
            byte blockStart = 0x0;
            List<StatusByte> unitStatusBytes = new LinkedList<>();
            for (int i = blockStart; i <= 88 - 4; i = i + 4) {
                LOGGER.debug("Handling units 0-88 {},{},{},{}", i, (i + 1), (i + 2), (i + 3));
                unitStatusBytes.add(
                    new StatusByte(
                        AVAILABLE_UNITS.contains((byte) (i + 3)) ? GAVState.ON : GAVState.DOES_NOT_EXIST,
                        AVAILABLE_UNITS.contains((byte) (i + 2)) ? GAVState.ON : GAVState.DOES_NOT_EXIST,
                        AVAILABLE_UNITS.contains((byte) (i + 1)) ? GAVState.ON : GAVState.DOES_NOT_EXIST,
                        AVAILABLE_UNITS.contains((byte) (i + 0)) ? GAVState.ON : GAVState.DOES_NOT_EXIST
                    )
                );
            }
            LOGGER.debug("Produced {}, status bytes which equates to {} status", unitStatusBytes.size(), unitStatusBytes.size() * 4);
            CALData calData = new CALDataStatusExtended(CALCommandTypeContainer.CALCommandStatusExtended_25Bytes, null, StatusCoding.BINARY_BY_THIS_SERIAL_INTERFACE, application, blockStart, unitStatusBytes, null, requestContext);
            CALReply calReply = new CALReplyShort((byte) 0x0, calData, cBusOptions, requestContext);
            EncodedReply encodedReply = new EncodedReplyCALReply((byte) 0x0, calReply, cBusOptions, requestContext);
            ReplyEncodedReply replyEncodedReply = new ReplyEncodedReply((byte) 0xC0, encodedReply, null, cBusOptions, requestContext);
            ReplyOrConfirmation replyOrConfirmation = new ReplyOrConfirmationReply((byte) 0xFF, replyEncodedReply, new ResponseTermination(), cBusOptions, requestContext);
            Alpha alpha = requestCommand.getAlpha();
            if (alpha != null) {
                Confirmation confirmation = new Confirmation(alpha, null, ConfirmationType.CONFIRMATION_SUCCESSFUL);
                replyOrConfirmation = new ReplyOrConfirmationConfirmation(alpha.getCharacter(), confirmation, replyOrConfirmation, cBusOptions, requestContext);
            }
            CBusMessage response = new CBusMessageToClient(replyOrConfirmation, requestContext, cBusOptions);
            LOGGER.info("Sending first part {}", response);
            ctx.writeAndFlush(response);
        }
        {
            byte blockStart = 88;
            List<StatusByte> unitStatusBytes = new LinkedList<>();
            for (int i = 88; i <= 88 + 88 - 4; i = i + 4) {
                LOGGER.debug("Handling units 88-176 {},{},{},{}", i, (i + 1), (i + 2), (i + 3));
                unitStatusBytes.add(
                    new StatusByte(
                        AVAILABLE_UNITS.contains((byte) (i + 3)) ? GAVState.ON : GAVState.DOES_NOT_EXIST,
                        AVAILABLE_UNITS.contains((byte) (i + 2)) ? GAVState.ON : GAVState.DOES_NOT_EXIST,
                        AVAILABLE_UNITS.contains((byte) (i + 1)) ? GAVState.ON : GAVState.DOES_NOT_EXIST,
                        AVAILABLE_UNITS.contains((byte) (i + 0)) ? GAVState.ON : GAVState.DOES_NOT_EXIST
                    )
                );
            }
            LOGGER.debug("Produced {}, status bytes which equates to {} status", unitStatusBytes.size(), unitStatusBytes.size() * 4);
            CALData calData = new CALDataStatusExtended(CALCommandTypeContainer.CALCommandStatusExtended_25Bytes, null, StatusCoding.BINARY_BY_THIS_SERIAL_INTERFACE, application, blockStart, unitStatusBytes, null, requestContext);
            CALReply calReply = new CALReplyShort((byte) 0x0, calData, cBusOptions, requestContext);
            EncodedReply encodedReply = new EncodedReplyCALReply((byte) 0x0, calReply, cBusOptions, requestContext);
            ReplyEncodedReply replyEncodedReply = new ReplyEncodedReply((byte) 0xC0, encodedReply, null, cBusOptions, requestContext);
            ReplyOrConfirmation replyOrConfirmation = new ReplyOrConfirmationReply((byte) 0xFF, replyEncodedReply, new ResponseTermination(), cBusOptions, requestContext);
            CBusMessage response = new CBusMessageToClient(replyOrConfirmation, requestContext, cBusOptions);
            LOGGER.info("Sending second part {}", response);
            ctx.writeAndFlush(response);
        }
        {
            byte blockStart = (byte) 176;
            List<StatusByte> unitStatusBytes = new LinkedList<>();
            for (int i = 176; i <= 176 + 80 - 4; i = i + 4) {
                LOGGER.debug("Handling units 176-256 {},{},{},{}", i, (i + 1), (i + 2), (i + 3));
                unitStatusBytes.add(
                    new StatusByte(
                        AVAILABLE_UNITS.contains((byte) (i + 3)) ? GAVState.ON : GAVState.DOES_NOT_EXIST,
                        AVAILABLE_UNITS.contains((byte) (i + 2)) ? GAVState.ON : GAVState.DOES_NOT_EXIST,
                        AVAILABLE_UNITS.contains((byte) (i + 1)) ? GAVState.ON : GAVState.DOES_NOT_EXIST,
                        AVAILABLE_UNITS.contains((byte) (i + 0)) ? GAVState.ON : GAVState.DOES_NOT_EXIST
                    )
                );
            }
            LOGGER.debug("Produced {}, status bytes which equates to {} status", unitStatusBytes.size(), unitStatusBytes.size() * 4);
            CALData calData = new CALDataStatusExtended(CALCommandTypeContainer.CALCommandStatusExtended_23Bytes, null, StatusCoding.BINARY_BY_THIS_SERIAL_INTERFACE, application, blockStart, unitStatusBytes, null, requestContext);
            CALReply calReply = new CALReplyShort((byte) 0x0, calData, cBusOptions, requestContext);
            EncodedReply encodedReply = new EncodedReplyCALReply((byte) 0x0, calReply, cBusOptions, requestContext);
            ReplyEncodedReply replyEncodedReply = new ReplyEncodedReply((byte) 0xC0, encodedReply, null, cBusOptions, requestContext);
            ReplyOrConfirmation replyOrConfirmation = new ReplyOrConfirmationReply((byte) 0xFF, replyEncodedReply, new ResponseTermination(), cBusOptions, requestContext);
            CBusMessage response = new CBusMessageToClient(replyOrConfirmation, requestContext, cBusOptions);
            LOGGER.info("Sending third part {}", response);
            ctx.writeAndFlush(response);
        }
    }

    private void handleCalDataIdentify(ChannelHandlerContext ctx, CALDataIdentify calDataIdentify, Alpha alpha) {
        short numBytes;
        IdentifyReplyCommand identifyReplyCommand;
        switch (calDataIdentify.getAttribute()) {
            case Manufacturer:
                numBytes = 0x08;
                identifyReplyCommand = new IdentifyReplyCommandManufacturer("Apache  ", numBytes);
                break;
            case Type:
                numBytes = 0x08;
                identifyReplyCommand = new IdentifyReplyCommandType("plc4x-si", numBytes);
                break;
            case FirmwareVersion:
                numBytes = 0x08;
                identifyReplyCommand = new IdentifyReplyCommandFirmwareVersion("  0.09  ", numBytes);
                break;
            case Summary:
                numBytes = 0x09;
                identifyReplyCommand = new IdentifyReplyCommandSummary("NOIDEA", (byte) 0xAF, "0900", numBytes);
                break;
            case ExtendedDiagnosticSummary:
                numBytes = 0x0C;
                identifyReplyCommand = new IdentifyReplyCommandExtendedDiagnosticSummary(ApplicationIdContainer.FREE_USAGE_01, ApplicationIdContainer.FREE_USAGE_0F, (byte) 0x0, 0x0, 4711L, (byte) 0x13, false, false, false, true, false, false, false, false, false, false, false, false, false, numBytes);
                break;
            case NetworkTerminalLevels:
                numBytes = 0x0C;
                identifyReplyCommand = new IdentifyReplyCommandNetworkTerminalLevels(new byte[]{0x13, 0x13, 0x13, 0x13, 0x13, 0x13, 0x13, 0x13, 0x13, 0x13, 0x13, 0x13}, numBytes);
                break;
            case TerminalLevel:
                numBytes = 0x0C;
                identifyReplyCommand = new IdentifyReplyCommandTerminalLevels(new byte[]{0x13, 0x13, 0x13, 0x13, 0x13, 0x13, 0x13, 0x13, 0x13, 0x13, 0x13, 0x13}, numBytes);
                break;
            case NetworkVoltage:
                numBytes = 0x05;
                identifyReplyCommand = new IdentifyReplyCommandNetworkVoltage("48", "7", numBytes);
                break;
            case GAVValuesCurrent:
                numBytes = 0x10;
                identifyReplyCommand = new IdentifyReplyCommandGAVValuesCurrent(new byte[]{
                    0x13, 0x13, 0x13, 0x13, 0x13, 0x13, 0x13, 0x13,
                    0x13, 0x13, 0x13, 0x13, 0x13, 0x13, 0x13, 0x13,
                }, numBytes);
                break;
            case GAVValuesStored:
                numBytes = 0x10;
                identifyReplyCommand = new IdentifyReplyCommandGAVValuesStored(new byte[]{
                    0x13, 0x13, 0x13, 0x13, 0x13, 0x13, 0x13, 0x13,
                    0x13, 0x13, 0x13, 0x13, 0x13, 0x13, 0x13, 0x13,
                }, numBytes);
                break;
            case GAVPhysicalAddresses:
                numBytes = 0x10;
                identifyReplyCommand = new IdentifyReplyCommandGAVPhysicalAddresses(new byte[]{
                    0x13, 0x13, 0x13, 0x13, 0x13, 0x13, 0x13, 0x13,
                    0x13, 0x13, 0x13, 0x13, 0x13, 0x13, 0x13, 0x13,
                }, numBytes);
                break;
            case LogicalAssignment:
                numBytes = 0x0E;
                identifyReplyCommand = new IdentifyReplyCommandLogicalAssignment(Collections.singletonList(new LogicAssignment(false, true, true, true, true, true)), numBytes);
                break;
            case Delays:
                numBytes = 0x0F;
                identifyReplyCommand = new IdentifyReplyCommandDelays(new byte[]{0x3}, (byte) 0x13, numBytes);
                break;
            case MinimumLevels:
                numBytes = 0x0E;
                identifyReplyCommand = new IdentifyReplyCommandMinimumLevels(new byte[]{0x3}, numBytes);
                break;
            case MaximumLevels:
                numBytes = 0x0F;
                identifyReplyCommand = new IdentifyReplyCommandMaximumLevels(new byte[]{0xF}, numBytes);
                break;
            case CurrentSenseLevels:
                numBytes = 0x10;
                identifyReplyCommand = new IdentifyReplyCommandCurrentSenseLevels(new byte[]{0xF}, numBytes);
                break;
            case OutputUnitSummary:
                numBytes = 0x12;
                identifyReplyCommand = new IdentifyReplyCommandOutputUnitSummary(new IdentifyReplyCommandUnitSummary(false, false, false, false, false, false, false, false), (byte) 0x4, (byte) 0x4, (short) 45, numBytes);
                break;
            case DSIStatus:
                numBytes = 0x12;
                identifyReplyCommand = new IdentifyReplyCommandDSIStatus(ChannelStatus.OK, ChannelStatus.OK, ChannelStatus.OK, ChannelStatus.OK, ChannelStatus.OK, ChannelStatus.OK, ChannelStatus.OK, ChannelStatus.OK, UnitStatus.OK, (byte) 0x34, numBytes);
                break;
            default:
                throw new IllegalStateException("unmapped type " + calDataIdentify.getAttribute());
        }

        CALData calData = new CALDataIdentifyReply(getReplyCommandType(numBytes + 1), null, calDataIdentify.getAttribute(), identifyReplyCommand, requestContext);
        CALReply calReply;
        if (exstat) {
            calReply = new CALReplyLong((byte) 0x0, calData, (byte) 0x0, new UnitAddress((byte) 0x0), null, new SerialInterfaceAddress((byte) 0x02), (byte) 0x0, null, cBusOptions, requestContext);
        } else {
            calReply = new CALReplyShort((byte) 0x0, calData, cBusOptions, requestContext);
        }
        EncodedReply encodedReply = new EncodedReplyCALReply((byte) 0x0, calReply, cBusOptions, requestContext);
        ReplyEncodedReply replyEncodedReply = new ReplyEncodedReply((byte) 0xC0, encodedReply, null, cBusOptions, requestContext);
        ReplyOrConfirmation replyOrConfirmation = new ReplyOrConfirmationReply((byte) 0xFF, replyEncodedReply, new ResponseTermination(), cBusOptions, requestContext);
        if (alpha != null) {
            Confirmation confirmation = new Confirmation(alpha, null, ConfirmationType.CONFIRMATION_SUCCESSFUL);
            replyOrConfirmation = new ReplyOrConfirmationConfirmation(alpha.getCharacter(), confirmation, replyOrConfirmation, cBusOptions, requestContext);
        }
        CBusMessage response = new CBusMessageToClient(replyOrConfirmation, requestContext, cBusOptions);
        LOGGER.info("Send identify response\n{}", response);
        ctx.writeAndFlush(response);
    }

    private static void handleSmartConnect(RequestSmartConnectShortcut requestSmartConnectShortcut) {
        LOGGER.info("Handling RequestSmartConnectShortcut\n{}", requestSmartConnectShortcut);
        smart = true;
        connect = true;
        buildCBusOptions();
    }

    private void handleReset(RequestReset requestReset) {
        LOGGER.info("Handling RequestReset\n{}", requestReset);
        connect = false;
        smart = false;
        idmon = false;
        exstat = false;
        monitor = false;
        monall = false;
        pun = false;
        pcn = false;
        srchk = false;
        stopSALMonitor();
        stopMMIMonitor();
    }

    private void startSALMonitor(ChannelHandlerContext ctx) {
        if (salMonitorFuture != null) {
            LOGGER.debug("SAL Monitor already running");
            return;
        }
        LOGGER.info("Starting monitor");
        salMonitorFuture = ctx.executor().scheduleAtFixedRate(() -> {
            if (monitorApplicationAddress1 != 0x38 && monitorApplicationAddress2 != 0x38 && monitorApplicationAddress1 != (byte) 0xFF && monitorApplicationAddress2 != (byte) 0xFF) {
                LOGGER.debug("Filtered because monitor application address 1 {} monitor application address 1 {}", monitorApplicationAddress1, monitorApplicationAddress2);
                return;
            }
            try {
                outputLock.lock();
                MonitoredSAL monitoredSAL;
                if (cBusOptions.getExstat()) {
                    LightingData lightingData;
                    double random = Math.random();
                    if (random < 0.25) {
                        lightingData = new LightingDataOn(LightingCommandTypeContainer.LightingCommandOn, (byte) 0xAF);
                    } else if (random > 0.25 && random < 0.5) {
                        lightingData = new LightingDataOff(LightingCommandTypeContainer.LightingCommandOff, (byte) 0xAF);
                    } else if (random > 0.5 && random < 0.75) {
                        lightingData = new LightingDataRampToLevel(LightingCommandTypeContainer.LightingCommandRampToLevel_20Second, (byte) 0xAF, (byte) 0xE0);
                    } else {
                        lightingData = new LightingDataTerminateRamp(LightingCommandTypeContainer.LightingCommandTerminateRamp, (byte) 0xAF);
                    }
                    SALData salData = new SALDataLighting(null, lightingData);
                    monitoredSAL = new MonitoredSALLongFormSmartMode((byte) 0x05, (byte) 0x00, new UnitAddress((byte) 0x0), null, ApplicationIdContainer.LIGHTING_38, (byte) 0x00, null, salData, cBusOptions);
                } else {
                    LightingData lightingData;
                    double random = Math.random();
                    if (random < 0.25) {
                        lightingData = new LightingDataOn(LightingCommandTypeContainer.LightingCommandOn, (byte) 0xAF);
                    } else if (random > 0.25 && random < 0.5) {
                        lightingData = new LightingDataOff(LightingCommandTypeContainer.LightingCommandOff, (byte) 0xAF);
                    } else if (random > 0.5 && random < 0.75) {
                        lightingData = new LightingDataRampToLevel(LightingCommandTypeContainer.LightingCommandRampToLevel_20Second, (byte) 0xAF, (byte) 0xE0);
                    } else {
                        lightingData = new LightingDataTerminateRamp(LightingCommandTypeContainer.LightingCommandTerminateRamp, (byte) 0xAF);
                    }
                    SALData salData = new SALDataLighting(null, lightingData);
                    monitoredSAL = new MonitoredSALShortFormBasicMode((byte) 0x0, (byte) 0x0, (short) 0x0, (short) 0x0, (byte) 0x0, ApplicationIdContainer.LIGHTING_38, salData, cBusOptions);
                }
                EncodedReply encodedReply = new MonitoredSALReply((byte) 0x0, monitoredSAL, cBusOptions, requestContext);
                Reply reply = new ReplyEncodedReply((byte) 0x0, encodedReply, null, cBusOptions, requestContext);
                ReplyOrConfirmation replyOrConfirmation = new ReplyOrConfirmationReply((byte) 0x00, reply, new ResponseTermination(), cBusOptions, requestContext);
                CBusMessage message = new CBusMessageToClient(replyOrConfirmation, requestContext, cBusOptions);
                LOGGER.info("[SAL Monitor] Sending out\n{}", message);
                ctx.writeAndFlush(message);
            } finally {
                outputLock.unlock();
            }
        }, 5, 5, TimeUnit.SECONDS);
    }

    private void stopSALMonitor() {
        if (salMonitorFuture == null) {
            return;
        }
        LOGGER.info("Stopping SAL monitor");
        salMonitorFuture.cancel(false);
        salMonitorFuture = null;
    }

    private void startMMIMonitor(ChannelHandlerContext ctx) {
        if (mmiMonitorFuture != null) {
            LOGGER.debug("MMI Monitor already running");
            return;
        }
        LOGGER.info("Starting MMI monitor");
        mmiMonitorFuture = ctx.executor().scheduleAtFixedRate(() -> {
            // TODO: for whatever reason those are not send with a crc
            CBusOptions cBusOptions = new CBusOptions(connect, smart, idmon, exstat, monitor, monall, pun, pcn, false);
            try {
                outputLock.lock();
                CALReply calReply;
                if (exstat) {
                    // TODO: map actuall values from simulator
                    List<StatusByte> statusBytes = new LinkedList<>();
                    for (int i = 0; i < 22; i++) {
                        statusBytes.add(new StatusByte(GAVState.ON, GAVState.ERROR, GAVState.OFF, GAVState.DOES_NOT_EXIST));
                    }
                    CALData calData = new CALDataStatusExtended(CALCommandTypeContainer.CALCommandStatusExtended_25Bytes, null, StatusCoding.BINARY_BY_ELSEWHERE, ApplicationIdContainer.LIGHTING_38, (byte) 0x00, statusBytes, null, requestContext);
                    calReply = new CALReplyLong((byte) 0x86, calData, 0x00, new UnitAddress((byte) 0x04), null, new SerialInterfaceAddress((byte) 0x02), (byte) 0x00, null, cBusOptions, requestContext);
                } else {
                    List<StatusByte> statusBytes = new LinkedList<>();
                    // TODO: map actuall values from simulator
                    for (int i = 0; i < 23; i++) {
                        statusBytes.add(new StatusByte(GAVState.ON, GAVState.ERROR, GAVState.OFF, GAVState.DOES_NOT_EXIST));
                    }
                    CALData calData = new CALDataStatus(CALCommandTypeContainer.CALCommandStatus_25Bytes, null, ApplicationIdContainer.LIGHTING_38, (byte) 0x00, statusBytes, requestContext);
                    calReply = new CALReplyShort((byte) 0x0, calData, cBusOptions, requestContext);
                }
                EncodedReply encodedReply = new EncodedReplyCALReply((byte) 0x0, calReply, cBusOptions, requestContext);
                Reply reply = new ReplyEncodedReply((byte) 0x0, encodedReply, null, cBusOptions, requestContext);
                ReplyOrConfirmation replyOrConfirmation = new ReplyOrConfirmationReply((byte) 0x00, reply, new ResponseTermination(), cBusOptions, requestContext);
                CBusMessage message = new CBusMessageToClient(replyOrConfirmation, requestContext, cBusOptions);
                LOGGER.info("[MMI Monitor] Sending out\n{}", message);
                ctx.writeAndFlush(message);
            } finally {
                outputLock.unlock();
            }
        }, 5, 5, TimeUnit.SECONDS);
    }

    private void stopMMIMonitor() {
        if (mmiMonitorFuture == null) {
            return;
        }
        LOGGER.info("Stopping monitor");
        mmiMonitorFuture.cancel(false);
        mmiMonitorFuture = null;
    }

    private CALCommandTypeContainer getReplyCommandType(int numBytes) {
        for (CALCommandTypeContainer value : CALCommandTypeContainer.values()) {
            if (value.getCommandType() == CALCommandType.REPLY && value.getNumBytes() == numBytes) {
                return value;
            }
        }
        throw new IllegalArgumentException("No reply type for " + numBytes);
    }

}
