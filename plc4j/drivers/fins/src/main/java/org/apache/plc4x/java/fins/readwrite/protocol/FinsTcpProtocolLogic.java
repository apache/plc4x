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
package org.apache.plc4x.java.fins.readwrite.protocol;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.Channel;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelFutureListener;
import io.netty.util.AttributeKey;
import org.apache.commons.lang3.concurrent.BasicThreadFactory;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcResponse;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.fins.readwrite.FinsTcpMessage;
import org.apache.plc4x.java.fins.readwrite.context.FinsTcpDriverContext;
import org.apache.plc4x.java.fins.readwrite.tag.FinsTcpPlcTagHandler;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.connection.PlcTagHandler;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.generation.Message;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.messages.utils.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.messages.utils.PlcResponseItem;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;
import org.apache.plc4x.java.spi.values.PlcDATE_AND_LTIME;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcNull;
import org.apache.plc4x.java.spi.values.PlcSINT;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.time.LocalDateTime;
import java.util.HashMap;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * The S7 Protocol states that there can not be more then {min(maxAmqCaller, maxAmqCallee} "ongoing" requests.
 * So we need to limit those.
 * Thus, each request goes to a Work Queue and this Queue ensures, that only 3 are open at the same time.
 */
public class FinsTcpProtocolLogic extends Plc4xProtocolBase<FinsTcpMessage> {

    private static final Logger logger = LoggerFactory.getLogger(FinsTcpProtocolLogic.class);

    public static final Duration REQUEST_TIMEOUT = Duration.ofMillis(10000);
    private final AtomicInteger tpduGenerator = new AtomicInteger(10);

    /*
     * Task group for managing connection redundancy.
     */
    private final ExecutorService clientExecutorService = Executors.newFixedThreadPool(4, new BasicThreadFactory.Builder()
        .namingPattern("plc4x-app-thread-%d")
        .daemon(true)
        .priority(Thread.MAX_PRIORITY)
        .build());



    /*
     * This array stores the cyclic subscription requests between the driver
     * and the PLC. The purpose is to document the tags associated with the
     * request. Each subscription uses a 'JobID' that is managed by the PLC and
     * obtained from the response to the request. In the following,
     * the values sent PUSH from the PLC to the driver refer to this JobID.
     */
    private final Map<Short, PlcSubscriptionRequest> cycRequests = new HashMap<>();
    
    /*
    * This data structure stores the last value associated with a cyclic 
    * subscription request. In each event received, the values of the internal 
    * PlcValue are compared and if any of them are different, the new value is 
    * transferred to the event stack and the value is updated in this HashMap.
    */
    private FinsTcpDriverContext finsTcpDriverContext;
    private RequestTransactionManager tm;

    @Override
    public void setDriverContext(DriverContext driverContext) {
        super.setDriverContext(driverContext);
        this.finsTcpDriverContext = (FinsTcpDriverContext) driverContext;

        // Initialize Transaction Manager.
        // Until the number of concurrent requests is successfully negotiated we set it to a
        // maximum of only one request being able to be sent at a time. During the login process
        // No concurrent requests can be sent anyway. It will be updated when receiving the
        // S7ParameterSetupCommunication response.
        this.tm = new RequestTransactionManager(1);
    }

    @Override
    public PlcTagHandler getTagHandler() {
        return new FinsTcpPlcTagHandler();
    }

    @Override
    public void close(ConversationContext<FinsTcpMessage> context) {
        // TODO: Find out how to close this prior to Java 19
        //clientExecutorService.close();
        tm.shutdown();
    }

    @Override
    public void onConnect(ConversationContext<FinsTcpMessage> context) {

        //Set feature for all handlers in the pipeline from
        //the driver configuration.
        setChannelFeatures();

        String order = "46494E530000000C0000000000000000000000";

        String hexString = Integer.toHexString(10);
        String handshakeMsg =  order + String.format("%2s", hexString).replace(' ', '0');
        logger.info("握手的消息:{}", handshakeMsg);
        ByteBuf byteBuf = Unpooled.wrappedBuffer(hexToByteArray(handshakeMsg));

        Channel channel = context.getChannel();

        if(channel.isOpen() && channel.isActive() && channel.isWritable())
        {
            //必须有这句话
            byteBuf.retain();
            AttributeKey<Boolean> IS_CONNECTED = AttributeKey.valueOf("IS_CONNECTED");
            ChannelFuture future = channel.writeAndFlush(byteBuf);
            future.addListener((ChannelFutureListener) future1 -> {
                if (future1.isSuccess()) {
                    System.out.println("Message sent successfully");
                    conversationContext.getChannel().attr(IS_CONNECTED).set(true);
                } else {
                    System.err.println("Message send failed");
                    conversationContext.getChannel().attr(IS_CONNECTED).set(false);
                    future1.cause().printStackTrace();
                }
            });


            context.fireConnected();
        }

    }


    /*
     * It performs the sequential and safe shutdown of the driver.
     * Completion of pending requests, executors and associated tasks.
     */
    @Override
    public void onDisconnect(ConversationContext<FinsTcpMessage> context) {
        // 1. Here we shut down the local task executor.
        clientExecutorService.shutdown();
        // 2. Performs the shutdown of the transaction executor.
        tm.shutdown();
        // 3. Finish the execution of the tasks for the handling of Events.
        //eventLogic.stop();
        // 4. Executes the closing of the main channel.
        context.getChannel().close();
        // 5. Here is the stop of any task or state machine that is added.
    }


    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        // If we're not connected, just abort with an error.
        if (!isConnected()) {
            CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();
            future.completeExceptionally(new PlcRuntimeException("Disconnected"));
            return future;
        }

        DefaultPlcReadRequest request = (DefaultPlcReadRequest) readRequest;
        CompletableFuture<FinsTcpMessage> responseFuture =null;



        // Just send a single response and chain it as Response
        return toPlcReadResponse(readRequest, responseFuture);
    }


    /**
     * Maps the S7ReadResponse of a PlcReadRequest to a PlcReadResponse
     */
    private CompletableFuture<PlcReadResponse> toPlcReadResponse(PlcReadRequest readRequest, CompletableFuture<FinsTcpMessage> responseFuture) {
        CompletableFuture<PlcReadResponse> clientFuture = new CompletableFuture<>();
        //Pointers
        FinsTcpMessage[] responseMessage = new FinsTcpMessage[1];
        PlcReadRequest[] plcReadRequest = new PlcReadRequest[1];

        responseFuture.whenComplete((finsTcpMessage, throwable) -> {
            if (throwable != null) {
                clientFuture.completeExceptionally(new PlcProtocolException("Error reading", throwable));
            } else {
                try {
                    responseMessage[0] = finsTcpMessage;
                    plcReadRequest[0] = readRequest;
                    clientExecutorService.submit(() -> {
                        try {
                            PlcReadResponse response = (PlcReadResponse) decodeReadResponse(responseMessage[0], plcReadRequest[0]);
                            clientFuture.complete(response);
                        } catch (Exception ex){

                        }
                    });

                } catch (Exception ex) {
                    logger.info(ex.toString());
                }
            }
        });

        return clientFuture;
    }

    private PlcResponse decodeReadResponse(FinsTcpMessage responseMessage, PlcReadRequest plcReadRequest) throws PlcProtocolException {
        Map<String, PlcResponseItem<PlcValue>> values = new HashMap<>();
        short errorClass;
        short errorCode;

        return new DefaultPlcReadResponse(plcReadRequest, values);
    }


    private void setChannelFeatures() {
        conversationContext.getChannel().attr(AttributeKey.valueOf("READ_TIME_OUT")).set(finsTcpDriverContext.getReadTimeout());
        conversationContext.getChannel().attr(AttributeKey.valueOf("IS_PIN_ACTIVE")).set(finsTcpDriverContext.getPing());
        conversationContext.getChannel().attr(AttributeKey.valueOf("PING_TIME")).set(finsTcpDriverContext.getPingTime());
        conversationContext.getChannel().attr(AttributeKey.valueOf("RETRY_TIME")).set(finsTcpDriverContext.getRetryTime());
    }

    private boolean isConnected() {
        AttributeKey<Boolean> IS_CONNECTED = AttributeKey.valueOf("IS_CONNECTED");
        return conversationContext.getChannel().attr(IS_CONNECTED).get();
        //return true;
    }


    /**
     * 16进制字符串转化为byte数组
     *
     * @param inHex 要转16进制字节流数组的16进制字符
     * @return 16进制字节流
     */
    private byte[] hexToByteArray(String inHex) {
        int hexlen = inHex.length();
        byte[] result;
        if (hexlen % 2 == 1) {
            // 奇数
            hexlen++;
            result = new byte[(hexlen / 2)];
            inHex = "0" + inHex;
        } else {
            // 偶数
            result = new byte[(hexlen / 2)];
        }
        int j = 0;
        for (int i = 0; i < hexlen; i += 2) {
            result[j] = hexToByte(inHex.substring(i, i + 2));
            j++;
        }
        return result;
    }

    private byte hexToByte(String inHex) {
        return (byte) Integer.parseInt(inHex, 16);
    }

}



