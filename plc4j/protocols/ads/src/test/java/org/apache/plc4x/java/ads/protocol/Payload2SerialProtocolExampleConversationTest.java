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
package org.apache.plc4x.java.ads.protocol;

import io.netty.buffer.Unpooled;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelFutureListener;
import io.netty.channel.ChannelHandlerContext;
import io.netty.util.concurrent.EventExecutor;
import io.netty.util.concurrent.ScheduledFuture;
import org.apache.commons.lang3.ArrayUtils;
import org.apache.commons.lang3.reflect.FieldUtils;
import org.apache.plc4x.java.ads.api.commands.AdsReadRequest;
import org.apache.plc4x.java.ads.api.commands.AdsReadResponse;
import org.apache.plc4x.java.ads.api.commands.types.*;
import org.apache.plc4x.java.ads.api.generic.AmsPacket;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;
import org.apache.plc4x.java.ads.api.serial.AmsSerialAcknowledgeFrame;
import org.apache.plc4x.java.ads.api.serial.AmsSerialFrame;
import org.apache.plc4x.java.ads.api.serial.types.FragmentNumber;
import org.apache.plc4x.java.ads.api.serial.types.UserData;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.hamcrest.core.IsEqual;
import org.junit.Before;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.ErrorCollector;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.concurrent.Callable;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;
import java.util.concurrent.atomic.AtomicInteger;

import static org.mockito.Mockito.*;

/**
 * @see <a href="https://infosys.beckhoff.com/content/1033/tcadsamsserialspec/html/tcamssericalspec,0xsample.htm?id=60692407917020132">example</a>
 */
public class Payload2SerialProtocolExampleConversationTest {

    @Rule
    public ErrorCollector errorCollector = new ErrorCollector();

    private Payload2SerialProtocol SUT;

    private ChannelHandlerContext channelHandlerContextMock;

    @Before
    public void setUp() {
        ExecutorService executorService = Executors.newFixedThreadPool(10);
        channelHandlerContextMock = mock(ChannelHandlerContext.class, RETURNS_DEEP_STUBS);
        when(channelHandlerContextMock.toString()).thenReturn("ChannelHandlerContextMock");
        when(channelHandlerContextMock.executor()).then(_ign -> {
            EventExecutor eventExecutor = mock(EventExecutor.class);
            when(eventExecutor.schedule(any(Callable.class), anyLong(), any()))
                .then(invocation -> {
                    Future<Object> submit = executorService.submit((Callable<Object>) invocation.getArgument(0));
                    ScheduledFuture scheduledFuture = mock(ScheduledFuture.class);
                    when(scheduledFuture.cancel(anyBoolean()))
                        .then(invocation1 -> submit.cancel(invocation1.getArgument(0)));
                    return scheduledFuture;
                });
            return eventExecutor;
        });
        ChannelFuture channelFuture = mock(ChannelFuture.class);
        when(channelFuture.addListener(any())).then(invocation -> {
            ChannelFutureListener channelFutureListener = invocation.getArgument(0);
            ChannelFuture mock = mock(ChannelFuture.class);
            when(mock.isSuccess()).thenReturn(true);
            channelFutureListener.operationComplete(mock);
            return mock(ChannelFuture.class);
        });
        when(channelHandlerContextMock.writeAndFlush(any())).thenReturn(channelFuture);
        SUT = new Payload2SerialProtocol();
    }

    @Test(expected = PlcProtocolException.class)
    public void testWrongCrc() throws Exception {
        int[] exampleRequestInt = {
            /*Magic Cookie    */    0x01, 0xA5,
            /*Sender          */    0x00,
            /*Empfaenger      */    0x00,
            /*Fragmentnummer  */    0x06,
            /*Datenlaenge     */    0x2C,
            /*NetID Empfaenger*/    0xC0, 0xA8, 0x64, 0xAE, 0x01, 0x01,
            /*Port Nummer     */    0x21, 0x03,
            /*NetID Sender    */    0xC0, 0xA8, 0x64, 0x9C, 0x01, 0x01,
            /*Portnummer      */    0x01, 0x80,
            /*Kommando lesen  */    0x02, 0x00,
            /*Status          */    0x04, 0x00,
            /*Anzahl Datenbyte*/    0x0C, 0x00, 0x00, 0x00,
            /*Fehlercode      */    0x00, 0x00, 0x00, 0x00,
            /*InvokeID        */    0x07, 0x00, 0x00, 0x00,
            /*Index Gruppe    */    0x05, 0xF0, 0x00, 0x00,
            /*Index Offset    */    0x04, 0x00, 0x00, 0x9D,
            /*Anzahl Byte     */    0x02, 0x00, 0x00, 0x00,
            // This Checksum is flipped to provoke exception
            /*Checksumme      */    0x28, 0x79,
        };
        byte[] exampleRequest = ArrayUtils.toPrimitive(Arrays
            .stream(exampleRequestInt)
            .mapToObj(value -> (byte) value)
            .toArray(Byte[]::new));
        SUT.decode(channelHandlerContextMock, Unpooled.wrappedBuffer(exampleRequest), new ArrayList<>());
    }

    @Test
    public void exampleConversation() throws Exception {
        FieldUtils.writeDeclaredField(SUT, "fragmentCounter", new AtomicInteger(6), true);
        // 1. Terminal --> PLC : Request of 2 bytre data
        int[] exampleRequestInt = {
            /*Magic Cookie    */    0x01, 0xA5,
            /*Sender          */    0x00,
            /*Empfaenger      */    0x00,
            /*Fragmentnummer  */    0x06,
            /*Datenlaenge     */    0x2C,
            /*NetID Empfaenger*/    0xC0, 0xA8, 0x64, 0xAE, 0x01, 0x01,
            /*Port Nummer     */    0x21, 0x03,
            /*NetID Sender    */    0xC0, 0xA8, 0x64, 0x9C, 0x01, 0x01,
            /*Portnummer      */    0x01, 0x80,
            /*Kommando lesen  */    0x02, 0x00,
            /*Status          */    0x04, 0x00,
            /*Anzahl Datenbyte*/    0x0C, 0x00, 0x00, 0x00,
            /*Fehlercode      */    0x00, 0x00, 0x00, 0x00,
            /*InvokeID        */    0x07, 0x00, 0x00, 0x00,
            /*Index Gruppe    */    0x05, 0xF0, 0x00, 0x00,
            /*Index Offset    */    0x04, 0x00, 0x00, 0x9D,
            /*Anzahl Byte     */    0x02, 0x00, 0x00, 0x00,
            /*Checksumme      */    0x82, 0x97,
        };
        byte[] exampleRequest = ArrayUtils.toPrimitive(Arrays
            .stream(exampleRequestInt)
            .mapToObj(value -> (byte) value)
            .toArray(Byte[]::new));
        AmsPacket amsPacket = AdsReadRequest.of(
            AmsNetId.of(Arrays.copyOfRange(exampleRequest, 6, 12)),
            AmsPort.of(Arrays.copyOfRange(exampleRequest, 12, 14)),
            AmsNetId.of(Arrays.copyOfRange(exampleRequest, 14, 20)),
            AmsPort.of(Arrays.copyOfRange(exampleRequest, 20, 22)),
            Invoke.of(0x07),
            IndexGroup.of(0xF0_05),
            IndexOffset.of(0x9D_00_00_04L),
            Length.of(0x2)
        );
        AmsSerialFrame amsSerialFrame = AmsSerialFrame.of(FragmentNumber.of((byte) 0x06), UserData.of(amsPacket.getBytes()));
        errorCollector.checkThat("example request not same", amsSerialFrame.getBytes(), IsEqual.equalTo(exampleRequest));
        SUT.encode(channelHandlerContextMock, amsPacket.getByteBuf(), new ArrayList<>());

        // PLC --> Terminal : Acknowledge:
        int[] exampleResponsAckInt = {
            /*Magic Cookie   */      0x01, 0x5A,
            /*Sender         */      0x00,
            /*Empfaenger     */      0x00,
            /*Fragmentnummer */      0x06,
            /*Datenlaenge    */      0x00,
            /*Checksumme     */      0x67, 0x5A
        };
        byte[] exampleAckResponse = ArrayUtils.toPrimitive(Arrays
            .stream(exampleResponsAckInt)
            .mapToObj(value -> (byte) value)
            .toArray(Byte[]::new));

        AmsSerialAcknowledgeFrame amsSerialAcknowledgeFrame = AmsSerialAcknowledgeFrame.of(amsSerialFrame.getTransmitterAddress(), amsSerialFrame.getReceiverAddress(), amsSerialFrame.getFragmentNumber());
        errorCollector.checkThat("ack response not same", amsSerialAcknowledgeFrame.getBytes(), IsEqual.equalTo(exampleRequest));
        SUT.decode(channelHandlerContextMock, Unpooled.wrappedBuffer(exampleAckResponse), new ArrayList<>());
        SUT.decode(channelHandlerContextMock, amsSerialAcknowledgeFrame.getByteBuf(), new ArrayList<>());

        // PLC sends data:
        int[] exampleResponseInt = {
            /*Magic Cookie     */    0x01, 0xA5,
            /*Sender           */    0x00,
            /*Empfaenger       */    0x00,
            /*Fragmentnummer   */    0xEC,
            /*Anzahl Daten     */    0x2A,
            /*NetID Empfaenger */    0xC0, 0xA8, 0x64, 0x9C, 0x01, 0x01,
            /*Portnummer       */    0x01, 0x80,
            /*NetID Sender     */    0xC0, 0xA8, 0x64, 0xAE, 0x01, 0x01,
            /*Portnummer       */    0x21, 0x03,
            /*Response Lesen   */    0x02, 0x00,
            /*Status           */    0x05, 0x00,
            /*Anzahl Daten     */    0x0A, 0x00, 0x00, 0x00,
            /*Fehlercode       */    0x00, 0x00, 0x00, 0x00,
            /*InvokeID         */    0x07, 0x00, 0x00, 0x00,
            /*Ergebnis         */    0x00, 0x00, 0x00, 0x00,
            /*Anzahl Daten     */    0x02, 0x00, 0x00, 0x00,
            /*Daten            */    0xAF, 0x27,
            /*Checksumme       */    0x04, 0xA9,
        };
        byte[] exampleResponse = ArrayUtils.toPrimitive(Arrays
            .stream(exampleResponseInt)
            .mapToObj(value -> (byte) value)
            .toArray(Byte[]::new));
        AmsPacket amsResponsePacket = AdsReadResponse.of(
            AmsNetId.of(Arrays.copyOfRange(exampleResponse, 6, 12)),
            AmsPort.of(Arrays.copyOfRange(exampleResponse, 12, 14)),
            AmsNetId.of(Arrays.copyOfRange(exampleResponse, 14, 20)),
            AmsPort.of(Arrays.copyOfRange(exampleResponse, 20, 22)),
            Invoke.of(0x07),
            Result.of(0x00),
            Data.of((byte) 0xAF, (byte) 0x27)
        );
        AmsSerialFrame amsResponseSerialFrame = AmsSerialFrame.of(FragmentNumber.of((byte) 0xEC), UserData.of(amsResponsePacket.getBytes()));
        errorCollector.checkThat("read response not same", amsResponseSerialFrame.getBytes(), IsEqual.equalTo(exampleRequest));
        SUT.decode(channelHandlerContextMock, Unpooled.wrappedBuffer(exampleResponse), new ArrayList<>());
        SUT.decode(channelHandlerContextMock, amsResponseSerialFrame.getByteBuf(), new ArrayList<>());
    }
}
