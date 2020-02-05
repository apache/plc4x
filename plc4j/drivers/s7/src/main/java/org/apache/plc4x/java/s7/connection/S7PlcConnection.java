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
package org.apache.plc4x.java.s7.connection;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.*;
import java.net.InetAddress;
import java.time.Instant;
import java.util.Collection;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ArrayBlockingQueue;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.function.Consumer;
import org.apache.commons.configuration2.Configuration;
import org.apache.commons.configuration2.SystemConfiguration;
import org.apache.commons.lang3.StringUtils;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.connection.ChannelFactory;
import org.apache.plc4x.java.base.connection.NettyPlcConnection;
import org.apache.plc4x.java.base.events.ConnectEvent;
import org.apache.plc4x.java.base.events.ConnectedEvent;
import org.apache.plc4x.java.base.messages.*;
import org.apache.plc4x.java.isoontcp.protocol.IsoOnTcpProtocol;
import org.apache.plc4x.java.isotp.protocol.IsoTPProtocol;
import org.apache.plc4x.java.isotp.protocol.model.tpdus.DisconnectRequestTpdu;
import org.apache.plc4x.java.isotp.protocol.model.types.DeviceGroup;
import org.apache.plc4x.java.isotp.protocol.model.types.DisconnectReason;
import org.apache.plc4x.java.isotp.protocol.model.types.TpduSize;
import org.apache.plc4x.java.s7.model.S7Field;
import org.apache.plc4x.java.s7.netty.Plc4XS7Protocol;
import org.apache.plc4x.java.s7.netty.S7Protocol;
import org.apache.plc4x.java.s7.netty.model.messages.S7PushMessage;
import org.apache.plc4x.java.s7.netty.model.params.CpuDiagnosticPushParameter;
import org.apache.plc4x.java.s7.netty.model.params.CpuServicesPushParameter;
import org.apache.plc4x.java.s7.netty.model.payloads.AlarmMessagePayload;
import org.apache.plc4x.java.s7.netty.model.payloads.CpuCyclicServicesResponsePayload;
import org.apache.plc4x.java.s7.netty.model.payloads.CpuDiagnosticMessagePayload;
import org.apache.plc4x.java.s7.netty.model.payloads.items.AssociatedValueItem;
import org.apache.plc4x.java.s7.netty.model.types.CpuServicesParameterSubFunctionGroup;
import org.apache.plc4x.java.s7.netty.model.types.MemoryArea;
import org.apache.plc4x.java.s7.netty.model.types.SubscribedEventType;
import org.apache.plc4x.java.s7.netty.strategies.DefaultS7MessageProcessor;
import org.apache.plc4x.java.s7.netty.util.S7PlcFieldEventHandler;
import org.apache.plc4x.java.s7.netty.util.S7PlcFieldHandler;
import org.apache.plc4x.java.s7.protocol.S7CyclicServicesSubscriptionHandle;
import org.apache.plc4x.java.s7.protocol.S7DiagnosticSubscriptionHandle;
import org.apache.plc4x.java.s7.protocol.event.S7AlarmEvent;
import org.apache.plc4x.java.s7.protocol.event.S7Event;
import org.apache.plc4x.java.s7.protocol.event.S7ModeEvent;
import org.apache.plc4x.java.s7.protocol.event.S7SysEvent;
import org.apache.plc4x.java.s7.types.S7ControllerType;
import org.apache.plc4x.java.s7.utils.S7TsapIdEncoder;
import org.apache.plc4x.java.tcp.connection.TcpSocketChannelFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Class implementing the Connection handling for Siemens S7.
 * The adressing of Values in S7 works as follows:
 * <p>
 * For adressing values from Datablocks the following syntax is used:
 * <pre>
 *     DATA_BLOCKS/{blockNumer}/{byteOffset}
 * </pre>
 * <p>
 * For adressing data from other memory segments like I/O, Markers, ...
 * <pre>
 *     {memory area}/{byte offset}
 *     or
 *     {memory area}/{byte offset}/{bit offset}
 * </pre>
 * where the {bit-offset} is optional.
 * All Available Memory Areas for this mode are defined in the {@link MemoryArea} enum.
 */
public class S7PlcConnection extends NettyPlcConnection implements PlcReader, PlcWriter, PlcSubscriber {

    private static final int ISO_ON_TCP_PORT = 102;

    // Fetch values from configuration
    private static final Configuration CONF = new SystemConfiguration();
    private static final long CLOSE_DEVICE_TIMEOUT_MS = CONF.getLong("plc4x.s7connection.close.device,timeout", 1_000);

    private static final Logger logger = LoggerFactory.getLogger(S7PlcConnection.class);

    private final int rack;
    private final int slot;

    private final short paramPduSize;
    private final short paramMaxAmqCaller;
    private final short paramMaxAmqCallee;
    private final S7ControllerType paramControllerType;

    private BlockingQueue<S7PushMessage> alarmsqueue;
    
    Map<Consumer, Collection<PlcSubscriptionHandle>> cyclicServicesSubscriptions = new HashMap();
    
    Map<SubscribedEventType, Map<Short, PlcSubscriptionHandle>> pushEventHandles = new HashMap();    
    Map<Short, PlcSubscriptionHandle> cyclicServicesHandles = new HashMap();
    
    private EventLoop alarmsloopthread;

    public S7PlcConnection(InetAddress address, int rack, int slot, String params) {
        this(new TcpSocketChannelFactory(address, ISO_ON_TCP_PORT), rack, slot, params);

        logger.info("Setting up S7 Connection with: host-name {}, rack {}, slot {}, pdu-size {}, max-amq-caller {}, " +
                "max-amq-callee {}", address.getHostAddress(), rack, slot,
            paramPduSize, paramMaxAmqCaller, paramMaxAmqCallee);
    }

    public S7PlcConnection(ChannelFactory channelFactory, int rack, int slot, String params) {
        super(channelFactory, true);

        this.rack = rack;
        this.slot = slot;

        short curParamPduSize = 1024;
        short curParamMaxAmqCaller = 8;
        short curParamMaxAmqCallee = 8;
        S7ControllerType curParamControllerType = S7ControllerType.ANY;

        if (!StringUtils.isEmpty(params)) {
            for (String param : params.split("&")) {
                String[] paramElements = param.split("=");
                String paramName = paramElements[0];
                if (paramElements.length == 2) {
                    String paramValue = paramElements[1];
                    switch (paramName) {
                        case "pdu-size":
                            curParamPduSize = Short.parseShort(paramValue);
                            break;
                        case "max-amq-caller":
                            curParamMaxAmqCaller = Short.parseShort(paramValue);
                            break;
                        case "max-amq-callee":
                            curParamMaxAmqCallee = Short.parseShort(paramValue);
                            break;
                        case "controller-type":
                            curParamControllerType = S7ControllerType.valueOf(paramValue);
                            break;
                        default:
                            logger.debug("Unknown parameter {} with value {}", paramName, paramValue);
                    }
                } else {
                    logger.debug("Unknown no-value parameter {}", paramName);
                }
            }
        }

        // It seems that the LOGO devices are a little picky about the pdu-size.
        // Instead of handling this out, they just hang up without any error message.
        // So in case of a LOGO controller, set this to a known working value.
        if(curParamControllerType == S7ControllerType.LOGO && curParamPduSize == 1024) {
            curParamPduSize = 480;
        }

        // IsoTP uses pre defined sizes. Find the smallest box,
        // that would be able to contain the requested pdu size.
        this.paramPduSize = curParamPduSize;
        this.paramMaxAmqCaller = curParamMaxAmqCaller;
        this.paramMaxAmqCallee = curParamMaxAmqCallee;
        this.paramControllerType = curParamControllerType;
                    
        /*
         * Take into account that the size of this buffer depends on the final device.
         * S7-300 goes from 20 to 300 and for S7-400 it goes from 300 to 10000. 
         * Depending on the configuration of the alarm system, a large number of 
         * them should be expected when starting the connection. 
         * (Examples of this are PCS7 and Braumat).
         * Alarm filtering, ack, etc. must be performed by the client application.
        */
        this.alarmsqueue = new ArrayBlockingQueue<>(1024);
        
        Map<Short, PlcSubscriptionHandle> usrEventHandle = new HashMap();
        Map<Short, PlcSubscriptionHandle> sysEventHandle = new HashMap();
        Map<Short, PlcSubscriptionHandle> modeEventHandle = new HashMap();
        Map<Short, PlcSubscriptionHandle> almEventHandle = new HashMap();
        
        pushEventHandles.put(SubscribedEventType.USR, usrEventHandle);
        pushEventHandles.put(SubscribedEventType.SYS, sysEventHandle);
        pushEventHandles.put(SubscribedEventType.MODE, modeEventHandle);
        pushEventHandles.put(SubscribedEventType.ALM_S, almEventHandle);
        pushEventHandles.put(SubscribedEventType.ALM_8, almEventHandle);
        
        alarmsloopthread = new EventLoop(channel,
                                this.alarmsqueue,
                                cyclicServicesSubscriptions,
                                pushEventHandles,
                                cyclicServicesHandles);

    }

    @Override
    public boolean canRead() {
        return true;
    }

    @Override
    public boolean canWrite() {
        return true;
    }
    
    @Override
    protected ChannelHandler getChannelHandler(CompletableFuture<Void> sessionSetupCompleteFuture) {
        //short calledTsapId = S7TsapIdEncoder.encodeS7TsapId(DeviceGroup.OS, 0, 0);
        //short callingTsapId = S7TsapIdEncoder.encodeS7TsapId(DeviceGroup.PG_OR_PC, rack, slot);
        short calledTsapId = S7TsapIdEncoder.encodeS7TsapId(DeviceGroup.PG_OR_PC, rack, slot);
        short callingTsapId = S7TsapIdEncoder.encodeS7TsapId(DeviceGroup.OS, 0, 0);
        
        return new ChannelInitializer() {
            @Override
            protected void initChannel(Channel channel) {
                // Build the protocol stack for communicating with the s7 protocol.
                ChannelPipeline pipeline = channel.pipeline();
                pipeline.addLast(new ChannelInboundHandlerAdapter() {
                    @Override
                    public void userEventTriggered(ChannelHandlerContext ctx, Object evt) throws Exception {
                        if (evt instanceof ConnectedEvent) {
                            sessionSetupCompleteFuture.complete(null);
                        } else {
                            super.userEventTriggered(ctx, evt);
                        }
                    }
                });
                pipeline.addLast(new IsoOnTcpProtocol());
                pipeline.addLast(new IsoTPProtocol(callingTsapId, calledTsapId, TpduSize.valueForGivenSize(paramPduSize)));
                pipeline.addLast(new S7Protocol(paramMaxAmqCaller, paramMaxAmqCallee, paramPduSize, paramControllerType,
                    new DefaultS7MessageProcessor()));
                pipeline.addLast(new Plc4XS7Protocol(alarmsqueue));
            }
        };
    }

    @Override
    protected void sendChannelCreatedEvent() {
        // Send an event to the pipeline telling the Protocol filters what's going on.
        channel.pipeline().fireUserEventTriggered(new ConnectEvent());
    }

    @Override
    public PlcField prepareField(String fieldQuery) throws PlcInvalidFieldException {
        return S7Field.of(fieldQuery);
    }

    public int getRack() {
        return rack;
    }

    public int getSlot() {
        return slot;
    }

    public short getParamPduSize() {
        return paramPduSize;
    }

    public int getParamMaxAmqCaller() {
        return paramMaxAmqCaller;
    }

    public int getParamMaxAmqCallee() {
        return paramMaxAmqCallee;
    }

    public S7ControllerType getParamControllerType() {
        return paramControllerType;
    }

    @Override
    public void connect() throws PlcConnectionException {
        super.connect(); 
        alarmsloopthread.start();
    }

    
    @Override
    public void close() throws PlcConnectionException {
        alarmsloopthread.cancel();
        if ((channel != null) && channel.isOpen()) {
            // Send the PLC a message that the connection is being closed.
            DisconnectRequestTpdu disconnectRequest = new DisconnectRequestTpdu(
                (short) 0x0000, (short) 0x000F, DisconnectReason.NORMAL, Collections.emptyList(),
                Unpooled.EMPTY_BUFFER);

            // In case of an ISO TP Class 0 connection, the remote is usually expected to actively
            // close the connection. So we add a listener waiting for this to happen.
            CompletableFuture<Void> disconnectFuture = new CompletableFuture<>();
            channel.closeFuture().addListener(
                (ChannelFutureListener) future -> disconnectFuture.complete(null));

            // Send the disconnect request.
            channel.writeAndFlush(disconnectRequest);
            // Wait for the configured time for the remote to close the session.
            try {
                disconnectFuture.get(CLOSE_DEVICE_TIMEOUT_MS, TimeUnit.MILLISECONDS);
            }
            // If the remote didn't close the connection within the given time-frame, we have to take
            // care of closing the connection.
            catch (TimeoutException e) {
                logger.debug("Remote didn't close connection within the configured timeout of {} ms, shutting down actively.", CLOSE_DEVICE_TIMEOUT_MS, e);
                channel.close();
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            } catch (ExecutionException e) {
                throw new PlcConnectionException(e);
            }

            // Do some additional cleanup operations ...
            // In normal operation, the channels event loop has a parent, however when running with
            // the embedded channel for unit tests, parent is null.
            if (channel.eventLoop().parent() != null) {
                channel.eventLoop().parent().shutdownGracefully();
            }
        }
        super.close();
    }

    @Override
    public PlcReadRequest.Builder readRequestBuilder() {
        return new DefaultPlcReadRequest.Builder(this, new S7PlcFieldHandler());
    }

    @Override
    public PlcWriteRequest.Builder writeRequestBuilder() {
        return new DefaultPlcWriteRequest.Builder(this, new S7PlcFieldHandler());
    }

    @Override
    public PlcSubscriptionRequest.Builder subscriptionRequestBuilder() {
        return new DefaultPlcSubscriptionRequest.Builder(this, new S7PlcFieldEventHandler());
    }    
    
    @Override
    public PlcUnsubscriptionRequest.Builder unsubscriptionRequestBuilder() {
        return new DefaultPlcUnsubscriptionRequest.Builder(this);
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        InternalPlcReadRequest internalReadRequest = checkInternal(readRequest, InternalPlcReadRequest.class);
        CompletableFuture<InternalPlcReadResponse> future = new CompletableFuture<>();
        PlcRequestContainer<InternalPlcReadRequest, InternalPlcReadResponse> container =
            new PlcRequestContainer<>(internalReadRequest, future);
        channel.writeAndFlush(container).addListener(f -> {
            if (!f.isSuccess()) {
                future.completeExceptionally(f.cause());
            }
        });
        return future
            .thenApply(PlcReadResponse.class::cast);
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        InternalPlcWriteRequest internalWriteRequest = checkInternal(writeRequest, InternalPlcWriteRequest.class);
        CompletableFuture<InternalPlcWriteResponse> future = new CompletableFuture<>();
        PlcRequestContainer<InternalPlcWriteRequest, InternalPlcWriteResponse> container =
            new PlcRequestContainer<>(internalWriteRequest, future);
        channel.writeAndFlush(container).addListener(f -> {
            if (!f.isSuccess()) {
                future.completeExceptionally(f.cause());
            }
        });
        return future
            .thenApply(PlcWriteResponse.class::cast);
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        
        InternalPlcSubscriptionRequest internalSubsRequest = checkInternal(subscriptionRequest, InternalPlcSubscriptionRequest.class);
        
        CompletableFuture<InternalPlcSubscriptionResponse> future = new CompletableFuture<>();
        
        PlcRequestContainer<InternalPlcSubscriptionRequest, InternalPlcSubscriptionResponse> container =
            new PlcRequestContainer<>(internalSubsRequest, future);
        
        channel.writeAndFlush(container).addListener(f -> {
            if (!f.isSuccess()) {
                future.completeExceptionally(f.cause());
            }
        });       
        
        return future.thenApply(PlcSubscriptionResponse.class::cast);       
    }

    @Override
    public CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        
        InternalPlcUnsubscriptionRequest internalUnsubsRequest = checkInternal(unsubscriptionRequest, InternalPlcUnsubscriptionRequest.class);
        
        CompletableFuture<InternalPlcUnsubscriptionResponse> future = new CompletableFuture<>();
        
        PlcRequestContainer<InternalPlcUnsubscriptionRequest, InternalPlcUnsubscriptionResponse> container =
            new PlcRequestContainer<>(internalUnsubsRequest, future);
        
        channel.writeAndFlush(container).addListener(f -> {
            if (!f.isSuccess()) {
                future.completeExceptionally(f.cause());
            }
        });      
        
        return future.thenApply(PlcUnsubscriptionResponse.class::cast);        
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles) {
        //Add any handler to 
        handles.forEach(handle ->{
            if (handle instanceof S7CyclicServicesSubscriptionHandle) {
                S7CyclicServicesSubscriptionHandle s7handle = (S7CyclicServicesSubscriptionHandle) handle;
                if (!cyclicServicesHandles.containsKey(s7handle.getJobId())){
                    cyclicServicesHandles.put((short) s7handle.getJobId(), s7handle);
                }
            }
            
            if (handle instanceof S7DiagnosticSubscriptionHandle) {
                S7DiagnosticSubscriptionHandle s7handle = (S7DiagnosticSubscriptionHandle) handle; 
                s7handle.getConsumers().add(consumer);
                s7handle.getSubscribedevents().forEach((event)->{
                    pushEventHandles.get(event).put(s7handle.getJobId(), s7handle);
                });
                
            }
        });
        
        if (!cyclicServicesSubscriptions.containsKey(consumer)){
            cyclicServicesSubscriptions.put(consumer, handles);
        }
        
        return null;
    }

    @Override
    public void unregister(PlcConsumerRegistration registration) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    private class EventLoop extends Thread {
        private volatile boolean cancelled;
        private final Channel channel;
        private boolean alarmquery;
        private int delay;
        private final Map<Consumer, Collection<PlcSubscriptionHandle>> cyclicServicesSubscriptions; 
        private final Map<SubscribedEventType, Map<Short, PlcSubscriptionHandle>> pushEventHandles;        
        private final Map<Short, PlcSubscriptionHandle> cyclicServicesHandles;
        private final BlockingQueue<S7PushMessage> alarmsqueue;
        
        EventLoop(Channel channel, 
                BlockingQueue<S7PushMessage> alarmsqueue, 
                Map<Consumer, Collection<PlcSubscriptionHandle>> cyclicServicesSubscriptions,
                Map<SubscribedEventType, Map<Short, PlcSubscriptionHandle>> pushEventHandles,
                Map<Short, PlcSubscriptionHandle> cyclicServicesHandles) {
            this.channel = channel;
            this.alarmsqueue = alarmsqueue;
            this.cyclicServicesSubscriptions = cyclicServicesSubscriptions;
            this.pushEventHandles = pushEventHandles;
            this.cyclicServicesHandles = cyclicServicesHandles;
            this.alarmquery = true;
            this.delay = 1000;
        }
        
        @Override
        public void run() {
            while (!cancelled) { 
                try {
                    S7PushMessage msg = alarmsqueue.poll(delay, TimeUnit.MILLISECONDS);
                    
                    if (msg != null){
                        if (msg instanceof AlarmMessagePayload){
                            AlarmMessagePayload themsg = (AlarmMessagePayload) msg;
                            logger.debug("AlarmMessagePayload: " + themsg);
                            dispathAlmEvents(pushEventHandles.get(SubscribedEventType.ALM_8), themsg); 
                            
                        } else if (msg instanceof CpuDiagnosticPushParameter) {                            
                            CpuDiagnosticPushParameter themsg = (CpuDiagnosticPushParameter) msg;
                            logger.debug("CpuDiagnosticPushParameter: " + themsg);  
                            dispathModeEvents(pushEventHandles.get(SubscribedEventType.MODE), themsg);
                            
                        } else if (msg instanceof CpuDiagnosticMessagePayload) {
                            CpuDiagnosticMessagePayload themsg = (CpuDiagnosticMessagePayload) msg;
                            logger.debug("CpuDiagnosticMessagePayload: " + themsg);
                            int EventID = Short.toUnsignedInt(themsg.getMsg().getEventID());
                            if ((EventID >= 0x0A000) & (EventID <= 0x0BFFF)) {
                                dispathSysEvents(pushEventHandles.get(SubscribedEventType.USR), themsg);  
                            } else {
                                dispathSysEvents(pushEventHandles.get(SubscribedEventType.SYS), themsg);                                
                            }
                            
                        } else if (msg instanceof CpuServicesPushParameter) {
                            CpuServicesPushParameter themsg = (CpuServicesPushParameter) msg;
                            logger.info("CpuServicesPushParameter: " + themsg);
                            
                        } else if (msg instanceof CpuCyclicServicesResponsePayload) {
                            CpuCyclicServicesResponsePayload themsg = (CpuCyclicServicesResponsePayload) msg;
                            logger.debug("CpuCyclicServicesResponsePayload: " + themsg + " JobId:" + themsg.getJobId());
                            
                            S7CyclicServicesSubscriptionHandle handle = 
                                    (S7CyclicServicesSubscriptionHandle) cyclicServicesHandles.get(themsg.getJobId());
                            if (handle != null) {
                                UpdateCyclicServicesData(handle, themsg);                            
                            };
                            
                        } else {
                            logger.info("Object type: " + msg.getClass());
                        }                       
                    } else {
                        if (alarmquery){
                            //TODO Send alarm query to plc
                        } else {
                            
                        }
                    }
                                                            
                } catch (InterruptedException ex) {
                    logger.info(ex.getLocalizedMessage());
                }
            }
            logger.info("Closing the alarm loop.");
        }
        
        private void dispathAlmEvents(Map<Short, PlcSubscriptionHandle> handles, AlarmMessagePayload payload)
        {          
            List<S7Event> alarmsEvent = null;
            Object object = payload.getMsgtype();
            if (object != null) {
                if (object instanceof CpuServicesParameterSubFunctionGroup) {
                    CpuServicesParameterSubFunctionGroup subFunction = (CpuServicesParameterSubFunctionGroup) object;
                    alarmsEvent = S7AlarmEvent.getAlarmsEvents(subFunction, payload);
                }
            }     
            
            for (PlcSubscriptionHandle handle:handles.values()){ 
                S7DiagnosticSubscriptionHandle s7handle = (S7DiagnosticSubscriptionHandle) handle;
                for (Consumer consumer:s7handle.getConsumers()){ 
                    if (alarmsEvent != null) {
                        for (S7Event event:alarmsEvent) {
                            consumer.accept(event);
                        }
                    }
                };
            };              
        }         
        
        private void dispathSysEvents(Map<Short, PlcSubscriptionHandle> handles, CpuDiagnosticMessagePayload payload)
        {
            handles.forEach((index,handle)->{
                S7DiagnosticSubscriptionHandle thehandle = (S7DiagnosticSubscriptionHandle) handle;
                thehandle.getConsumers().forEach((consumer)->{
                    S7SysEvent event = new S7SysEvent(payload.getMsg());
                    consumer.accept(event);
                });
            });              
        }        
        
        private void dispathModeEvents(Map<Short, PlcSubscriptionHandle> handles, CpuDiagnosticPushParameter parameter)
        {
            handles.forEach((index,handle)->{
                S7DiagnosticSubscriptionHandle thehandle = (S7DiagnosticSubscriptionHandle) handle;
                thehandle.getConsumers().forEach((consumer)->{
                    S7ModeEvent event = new S7ModeEvent(Instant.now(), parameter);
                    consumer.accept(event);
                });
            });              
        }              
        
        private void UpdateCyclicServicesData(S7CyclicServicesSubscriptionHandle handle,CpuCyclicServicesResponsePayload themsg){
            
            Map<String, AssociatedValueItem> items = handle.getItems();
            Collection<AssociatedValueItem> values = items.values();
            List<AssociatedValueItem> newvalues = themsg.getItems();
            int i = 0;
            synchronized(values) {
                for (AssociatedValueItem value:values){
                    AssociatedValueItem newvalue = newvalues.get(i);
                    value.getData().setBytes(0, newvalue.getData());
                    i++;                
                }
            }
            
            cyclicServicesSubscriptions.forEach((consumer,handles)->{
                if (handles.contains(handle)){
                    handles.forEach((exehandle)->{
                        Map<String, Pair<PlcResponseCode, ByteBuf>> fields = new HashMap<>();
                        S7CyclicServicesSubscriptionHandle exe2handle = (S7CyclicServicesSubscriptionHandle) exehandle;
                        Map<String, AssociatedValueItem> values2consumer = exe2handle.getItems();
                        
                        //PlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(Instant.now(), fields);
                        
                        values2consumer.forEach((index, itemvalue)->{
                            if (consumer != null) {
                                consumer.accept(null);
                            };
                            //logger.info("Procesando valores : " + index + "\r\n" + ByteBufUtil.prettyHexDump(itemvalue.getData()));
                        });
                        //Pair<PlcResponseCode, ByteBuf> newPair = new ImmutablePair<>(PlcResponseCode, stringItem);
                        //exehandle
                    });
                }
            });
            
        }              
        
        public void cancel() {
            cancelled = true;  
        }

        public boolean isCancelled() {
            return cancelled;
        }
        
    }    
    
}
