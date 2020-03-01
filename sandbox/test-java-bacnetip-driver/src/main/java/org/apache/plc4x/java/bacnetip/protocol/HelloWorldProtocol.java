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
package org.apache.plc4x.java.bacnetip.protocol;

import io.netty.channel.ChannelHandlerContext;
import org.apache.plc4x.java.bacnetip.readwrite.*;
import org.apache.plc4x.java.base.PlcMessageToMessageCodec;
import org.apache.plc4x.java.base.messages.PlcRequestContainer;

import java.util.List;

public class HelloWorldProtocol extends PlcMessageToMessageCodec<BVLC, PlcRequestContainer> {

    private int packetCount = 0;
    private long startTime = -1;

    @Override
    protected void encode(ChannelHandlerContext channelHandlerContext, PlcRequestContainer plcRequestContainer, List<Object> list) throws Exception {
        System.out.println(plcRequestContainer);
    }

    @Override
    protected void decode(ChannelHandlerContext channelHandlerContext, BVLC packet, List<Object> list) throws Exception {
        if(startTime == -1) {
            startTime = System.currentTimeMillis();
        }
        packetCount++;
        if(packetCount % 10000 == 0) {
            long curTime = System.currentTimeMillis();
            System.out.println("Read " + packetCount + " packets in " + (curTime - startTime) + "ms");
        }
/*        if(packet instanceof BVLCOriginalUnicastNPDU) {
            final NPDU npdu = ((BVLCOriginalUnicastNPDU) packet).getNpdu();
            final APDU apdu = npdu.getApdu();
            if(apdu instanceof APDUSimpleAck) {
                APDUSimpleAck ack = (APDUSimpleAck) apdu;
                System.out.println("Simple-ACK(" + ack.getOriginalInvokeId() + ")");
            } else if(apdu instanceof APDUConfirmedRequest) {
                APDUConfirmedRequest request = (APDUConfirmedRequest) apdu;
                final BACnetConfirmedServiceRequest serviceRequest = request.getServiceRequest();
                if(serviceRequest instanceof BACnetConfirmedServiceRequestConfirmedCOVNotification) {
                    BACnetConfirmedServiceRequestConfirmedCOVNotification covNotification = (BACnetConfirmedServiceRequestConfirmedCOVNotification) serviceRequest;
                    final BACnetTagWithContent[] notifications = covNotification.getNotifications();
                    System.out.println("Simple-ACK(" + request.getInvokeId() + "): Confirmed COV Notification [" + notifications.length + "]");
                } else {
                    System.out.println("Simple-ACK(" + request.getInvokeId() + "): Other");
                }
            } else {
                System.out.println("Other");
            }
        } else {
            System.out.println("Other");
        }*/
    }

}
