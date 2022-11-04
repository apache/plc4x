/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.simulator.server.bacnet.protocol;

import io.netty.channel.ChannelFutureListener;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInboundHandlerAdapter;
import org.apache.plc4x.java.bacnetip.readwrite.Error;
import org.apache.plc4x.java.bacnetip.readwrite.*;
import org.apache.plc4x.java.bacnetip.readwrite.utils.StaticHelper;
import org.apache.plc4x.simulator.model.Context;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class BacnetServerAdapter extends ChannelInboundHandlerAdapter {

    private static final Logger LOGGER = LoggerFactory.getLogger(BacnetServerAdapter.class);

    private Context context;

    // TODO: make configurable
    private static int DEVICE_INSTANCE = 4711;

    // TODO: make configurable
    private static int DEVICE_ID = 815;

    public BacnetServerAdapter(Context context) {
        LOGGER.info("Creating adapter with context {}", context);
    }


    @Override
    public void channelInactive(ChannelHandlerContext ctx) throws Exception {
        super.channelInactive(ctx);
    }

    @Override
    public void channelRead(ChannelHandlerContext ctx, Object msg) throws Exception {
        System.out.println("Got request");
        System.out.println(msg);
        if (!(msg instanceof BVLC)) {
            return;
        }
        BVLC bvlc = (BVLC) msg;
        if (!(bvlc instanceof BVLCOriginalUnicastNPDU)) {
            // TODO: write useful error
            ctx.writeAndFlush(new BVLCOriginalUnicastNPDU(
                new NPDU(
                    (short) 1,
                    new NPDUControl(true, false, false, false, NPDUNetworkPriority.NORMAL_MESSAGE),
                    0,
                    (short) 0,
                    null,
                    0,
                    (short) 0,
                    null,
                    (short) 0,
                    null,
                    new APDUError(
                        (short) 0,
                        BACnetConfirmedServiceChoice.READ_PROPERTY,
                        new BACnetErrorGeneral(new Error(
                            new ErrorClassTagged(new BACnetTagHeader((byte) 0, TagClass.APPLICATION_TAGS, (byte) 1, (short) 0, (short) 0, 0, 0L), ErrorClass.COMMUNICATION, 0, (short) 0, TagClass.APPLICATION_TAGS),
                            new ErrorCodeTagged(new BACnetTagHeader((byte) 0, TagClass.APPLICATION_TAGS, (byte) 1, (short) 0, (short) 0, 0, 0L), ErrorCode.VENDOR_PROPRIETARY_VALUE, 0, (short) 0, TagClass.APPLICATION_TAGS)
                        )),
                        0
                    ),
                    0
                ),
                0
            )).addListener((ChannelFutureListener) f -> {
                if (!f.isSuccess()) {
                    f.cause().printStackTrace();
                }
            });
            return;
        }
        BVLCOriginalUnicastNPDU bvlcOriginalUnicastNPDU = (BVLCOriginalUnicastNPDU) bvlc;
        // TODO: get messageTypeField
        APDU apdu = bvlcOriginalUnicastNPDU.getNpdu().getApdu();
        if (apdu instanceof APDUUnconfirmedRequest) {
            APDUUnconfirmedRequest apduUnconfirmedRequest = (APDUUnconfirmedRequest) apdu;
            BACnetUnconfirmedServiceRequest serviceRequest = apduUnconfirmedRequest.getServiceRequest();
            if (serviceRequest instanceof BACnetUnconfirmedServiceRequestWhoIs) {
                BACnetUnconfirmedServiceRequestWhoIs baCnetUnconfirmedServiceRequestWhoIs = (BACnetUnconfirmedServiceRequestWhoIs) serviceRequest;
                if (baCnetUnconfirmedServiceRequestWhoIs.getDeviceInstanceRangeLowLimit() != null) {
                    if (DEVICE_INSTANCE < baCnetUnconfirmedServiceRequestWhoIs.getDeviceInstanceRangeLowLimit().getActualValue().longValue()) {
                        // Ignoring because we out if limit
                        return;
                    }
                }
                if (baCnetUnconfirmedServiceRequestWhoIs.getDeviceInstanceRangeHighLimit() != null) {
                    if (DEVICE_INSTANCE > baCnetUnconfirmedServiceRequestWhoIs.getDeviceInstanceRangeHighLimit().getActualValue().longValue()) {
                        // Ignoring because we out if limit
                        return;
                    }
                }
                BVLCOriginalUnicastNPDU response = new BVLCOriginalUnicastNPDU(
                    new NPDU(
                        (short) 1,
                        new NPDUControl(false, false, false, false, NPDUNetworkPriority.NORMAL_MESSAGE),
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        null,
                        new APDUUnconfirmedRequest(
                            new BACnetUnconfirmedServiceRequestIAm(
                                StaticHelper.createBACnetApplicationTagObjectIdentifier(BACnetObjectType.DEVICE.getValue(), DEVICE_INSTANCE),
                                StaticHelper.createBACnetApplicationTagUnsignedInteger(1024),
                                StaticHelper.creatBACnetSegmentationTagged(BACnetSegmentation.NO_SEGMENTATION),
                                StaticHelper.createBACnetVendorIdApplicationTagged(BACnetVendorId.MAPPED.getVendorId()),
                                0
                            ),
                            0
                        ),
                        0
                    ),
                    0
                );
                System.out.println("Writing response");
                System.out.println(response);
                ctx.writeAndFlush(response).addListener((ChannelFutureListener) f -> {
                    if (!f.isSuccess()) {
                        f.cause().printStackTrace();
                    }
                });
            } else {
                throw new Exception(apdu.getClass() + " not set supported");
            }
        } else if (apdu instanceof APDUConfirmedRequest) {
            APDUConfirmedRequest apduConfirmedRequest = (APDUConfirmedRequest) apdu;
            // TODO: just handle read for now
            BVLCOriginalUnicastNPDU response = new BVLCOriginalUnicastNPDU(
                new NPDU(
                    (short) 1,
                    new NPDUControl(false, false, false, false, NPDUNetworkPriority.NORMAL_MESSAGE),
                    null,
                    null,
                    null,
                    null,
                    null,
                    null,
                    null,
                    null,
                    new APDUComplexAck(
                        false,
                        false,
                        apduConfirmedRequest.getInvokeId(),
                        null,
                        null,
                        new BACnetServiceAckReadProperty(
                            StaticHelper.createBACnetContextTagObjectIdentifier((byte) 0, 2, 1L),
                            StaticHelper.createBACnetPropertyIdentifierTagged((byte) 1, 85),
                            null,
                            new BACnetConstructedDataAnalogValuePresentValue(
                                StaticHelper.createBACnetOpeningTag((short) 3),
                                StaticHelper.createBACnetTagHeaderBalanced(true, (short) 3, 3L),
                                StaticHelper.createBACnetClosingTag((short) 3),
                                StaticHelper.createBACnetApplicationTagReal(101L),
                                null,
                                null
                            ),
                            0L
                        ),
                        null,
                        null,
                        0
                    ),
                    0
                ),
                0
            );
            System.out.println("Writing response");
            System.out.println(response);
            ctx.writeAndFlush(response).addListener((ChannelFutureListener) f -> {
                if (!f.isSuccess()) {
                    f.cause().printStackTrace();
                }
            });
        }
    }

}
