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
import io.vavr.CheckedRunnable;
import org.apache.plc4x.java.cbus.readwrite.*;
import org.apache.plc4x.simulator.model.Context;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

public class CBusServerAdapter extends ChannelInboundHandlerAdapter {

    private static final Logger LOGGER = LoggerFactory.getLogger(CBusServerAdapter.class);

    private Context context;

    private static final RequestContext requestContext = new RequestContext(false, false, false);

    private static boolean connect;
    private static boolean smart;
    private static boolean idmon;
    private static boolean exstat;
    private static boolean monitor;
    private static boolean monall;
    private static boolean pun;
    private static boolean pcn;
    private static boolean srchk;
    private static CBusOptions cBusOptions;

    private final Lock writeLock = new ReentrantLock();

    private ScheduledFuture<?> sf;

    public CBusServerAdapter(Context context) {
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
        if (sf != null)
            sf.cancel(false);
        super.channelInactive(ctx);
    }

    @Override
    public void channelRead(ChannelHandlerContext ctx, Object msg) throws Exception {
        if (!(msg instanceof CBusMessage)) {
            return;
        }
        if (!smart && !connect) {
            // In this mode every message will be echoed
            LOGGER.info("Sending echo");
            ctx.writeAndFlush(msg);
        }
        try {
            writeLock.lock();
            CBusMessage packet = (CBusMessage) msg;
            if (packet instanceof CBusMessageToClient) {
                LOGGER.info("Message to client not supported\n{}", packet);
                return;
            }
            CBusMessageToServer cBusMessageToServer = (CBusMessageToServer) packet;
            Request request = cBusMessageToServer.getRequest();
            if (request instanceof RequestEmpty || request instanceof RequestNull) {
                LOGGER.debug("Ignoring\n{}", request);
                return;
            }
            if (request instanceof RequestDirectCommandAccess) {
                RequestDirectCommandAccess requestDirectCommandAccess = (RequestDirectCommandAccess) request;
                CALData calData = requestDirectCommandAccess.getCalData();
                LOGGER.info("Handling RequestDirectCommandAccess\n{}\n{}", requestDirectCommandAccess, calData);

                // TODO: handle other cal data type
                if (calData instanceof CALDataWrite) {
                    CALDataWrite calDataWrite = (CALDataWrite) calData;
                    Runnable acknowledger = () -> {
                        CALDataAcknowledge calDataAcknowledge = new CALDataAcknowledge(CALCommandTypeContainer.CALCommandAcknowledge, null, calDataWrite.getParamNo(), (short) 0x0, requestContext);
                        CALReplyShort calReply = new CALReplyShort((byte) 0x0, calDataAcknowledge, cBusOptions, requestContext);
                        EncodedReplyCALReply encodedReply = new EncodedReplyCALReply((byte) 0x0, calReply, cBusOptions, requestContext);
                        ReplyEncodedReply replyEncodedReply = new ReplyEncodedReply((byte) 0x0, encodedReply, null, cBusOptions, requestContext);
                        ReplyOrConfirmationReply replyOrConfirmationReply = new ReplyOrConfirmationReply((byte) 0x0, replyEncodedReply, new ResponseTermination(), cBusOptions, requestContext);
                        CBusMessageToClient cBusMessageToClient = new CBusMessageToClient(replyOrConfirmationReply, requestContext, cBusOptions);
                        LOGGER.info("Sending ack\n{}\n{}", cBusMessageToClient, encodedReply);
                        ctx.writeAndFlush(cBusMessageToClient);
                    };
                    switch (calDataWrite.getParamNo().getParameterType()) {
                        case APPLICATION_ADDRESS_1:
                            // TODO: check settings for subscription etc.
                            acknowledger.run();
                            return;
                        case APPLICATION_ADDRESS_2:
                            // TODO: check settings for subscription etc.
                            acknowledger.run();
                            return;
                        case INTERFACE_OPTIONS_1:
                            InterfaceOptions1 interfaceOptions1 = ((ParameterValueInterfaceOptions1) calDataWrite.getParameterValue()).getValue();
                            idmon = interfaceOptions1.getIdmon();
                            monitor = interfaceOptions1.getMonitor();
                            if (monitor) startMonitor(ctx);
                            else stopMonitor();
                            smart = interfaceOptions1.getSmart();
                            srchk = interfaceOptions1.getSrchk();
                            // TODO: add support for xonxoff
                            // xonxoff = interfaceOptions1.getXonXoff();
                            connect = interfaceOptions1.getConnect();
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
                            if (monitor) startMonitor(ctx);
                            else stopMonitor();
                            smart = interfaceOptions1PowerUpSettings.getSmart();
                            srchk = interfaceOptions1PowerUpSettings.getSrchk();
                            // TODO: add support for xonxoff
                            // xonxoff = interfaceOptions1PowerUpSettings.getXonXoff();
                            connect = interfaceOptions1PowerUpSettings.getConnect();
                            buildCBusOptions();
                            acknowledger.run();
                            return;
                        case CUSTOM_MANUFACTURER:
                            // TODO: handle other parm typed
                            acknowledger.run();
                            return;
                        case SERIAL_NUMBER:
                            // TODO: handle other parm typed
                            acknowledger.run();
                            return;
                        case CUSTOM_TYPE:
                            // TODO: handle other parm typed
                            acknowledger.run();
                            return;
                        default:
                            throw new IllegalStateException("Unmapped type");
                    }
                }
                return;
            }
            if (request instanceof RequestCommand) {
                RequestCommand requestCommand = (RequestCommand) request;
                LOGGER.info("Handling RequestCommand\n{}", requestCommand);
                CBusCommand cbusCommand = requestCommand.getCbusCommand();
                LOGGER.info("Handling CBusCommand\n{}", cbusCommand);
                if (cbusCommand instanceof CBusCommandPointToPoint) {
                    CBusCommandPointToPoint cBusCommandPointToPoint = (CBusCommandPointToPoint) cbusCommand;
                    LOGGER.info("Handling CBusCommandPointToPoint\n{}", cBusCommandPointToPoint);
                    // TODO: handle this
                    return;
                }
                if (cbusCommand instanceof CBusCommandPointToMultiPoint) {
                    CBusCommandPointToMultiPoint cBusCommandPointToMultiPoint = (CBusCommandPointToMultiPoint) cbusCommand;
                    CBusPointToMultiPointCommand command = cBusCommandPointToMultiPoint.getCommand();
                    if (command instanceof CBusPointToMultiPointCommandStatus) {
                        CBusPointToMultiPointCommandStatus cBusPointToMultiPointCommandStatus = (CBusPointToMultiPointCommandStatus) command;
                        StatusRequest statusRequest = cBusPointToMultiPointCommandStatus.getStatusRequest();
                        if (statusRequest instanceof StatusRequestBinaryState) {
                            StatusRequestBinaryState statusRequestBinaryState = (StatusRequestBinaryState) statusRequest;
                            StatusHeader statusHeader = new StatusHeader((short) (2 + 1)); // 2 we have always + 1 as we got one status byte
                            // TODO: map actuall values from simulator
                            byte blockStart = 0x0;
                            List<StatusByte> statusBytes = List.of(new StatusByte(GAVState.ON, GAVState.ERROR, GAVState.OFF, GAVState.DOES_NOT_EXIST));
                            // TODO: this might be extended or standard depeding on exstat
                            StandardFormatStatusReply standardFormatStatusReply = new StandardFormatStatusReply(statusHeader, statusRequestBinaryState.getApplication(), blockStart, statusBytes);
                            EncodedReply encodedReply = new EncodedReplyStandardFormatStatusReply((byte) 0xC0, standardFormatStatusReply, cBusOptions, requestContext);
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
                            // TODO: handle this
                            return;
                        }
                        if (statusRequest instanceof StatusRequestLevel) {
                            StatusRequestLevel statusRequestLevel = (StatusRequestLevel) statusRequest;
                            ExtendedStatusHeader statusHeader = new ExtendedStatusHeader((short) (3 + 2)); // 3 we have always (coding is extra opposed to the standard) + 2 as we got one level information
                            StatusCoding coding = StatusCoding.LEVEL_BY_THIS_SERIAL_INTERFACE;
                            // TODO: map actuall values from simulator
                            byte blockStart = statusRequestLevel.getStartingGroupAddressLabel();
                            List<LevelInformation> levelInformations = List.of(new LevelInformationNormal(0x5555, LevelInformationNibblePair.Value_F, LevelInformationNibblePair.Value_F));
                            ExtendedFormatStatusReply extendedFormatStatusReply = new ExtendedFormatStatusReply(statusHeader, coding, statusRequestLevel.getApplication(), blockStart, null, levelInformations);
                            EncodedReply encodedReply = new EncodedReplyExtendedFormatStatusReply((byte) 0xC0, extendedFormatStatusReply, cBusOptions, requestContext);
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
                        return;
                    }
                    // TODO: handle this
                    return;
                }
                if (cbusCommand instanceof CBusCommandPointToPointToMultiPoint) {
                    CBusCommandPointToPointToMultiPoint cBusCommandPointToPointToMultiPoint = (CBusCommandPointToPointToMultiPoint) cbusCommand;
                    LOGGER.info("Handling CBusCommandPointToPointToMultiPoint\n{}", cBusCommandPointToPointToMultiPoint);
                    // TODO: handle this
                    return;
                }
                if (cbusCommand instanceof CBusCommandDeviceManagement) {
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
                return;
            }
            if (request instanceof RequestObsolete) {
                RequestObsolete requestObsolete = (RequestObsolete) request;
                LOGGER.info("Handling RequestObsolete\n{}", requestObsolete);
                // TODO: handle this
                return;
            }
            if (request instanceof RequestReset) {
                RequestReset requestReset = (RequestReset) request;
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
                return;
            }
            if (request instanceof RequestSmartConnectShortcut) {
                RequestSmartConnectShortcut requestSmartConnectShortcut = (RequestSmartConnectShortcut) request;
                LOGGER.info("Handling RequestSmartConnectShortcut\n{}", requestSmartConnectShortcut);
                // TODO: handle this
                return;
            }
        } finally {
            writeLock.unlock();
        }
    }

    private void startMonitor(ChannelHandlerContext ctx) {
        if (sf != null) {
            LOGGER.debug("Monitor already running");
            return;
        }
        LOGGER.info("Starting monitor");
        sf = ctx.executor().scheduleAtFixedRate(() -> {
            try {
                writeLock.lock();
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
                LOGGER.info("[Monitor] Sending out\n{}\n{}", message, encodedReply);
                ctx.writeAndFlush(message);
            } finally {
                writeLock.unlock();
            }
        }, 5, 5, TimeUnit.SECONDS);
    }

    private void stopMonitor() {
        LOGGER.info("Stopping monitor");
        sf.cancel(false);
        sf = null;
    }

}
