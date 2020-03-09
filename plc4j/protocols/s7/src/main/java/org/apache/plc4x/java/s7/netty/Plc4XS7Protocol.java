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
package org.apache.plc4x.java.s7.netty;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.PooledByteBufAllocator;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelFutureListener;
import io.netty.channel.ChannelHandlerContext;
import io.netty.util.AttributeKey;
import java.io.IOException;
import java.lang.reflect.Array;
import java.math.BigInteger;
import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.*;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.Function;
import java.util.logging.Level;
import java.util.stream.Collectors;
import java.util.stream.IntStream;
import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.exceptions.*;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcRequest;
import org.apache.plc4x.java.api.messages.PlcResponse;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.PlcMessageToMessageCodec;
import org.apache.plc4x.java.base.events.ConnectedEvent;
import org.apache.plc4x.java.base.messages.*;
import org.apache.plc4x.java.base.messages.items.*;
import org.apache.plc4x.java.base.model.InternalPlcSubscriptionHandle;
import org.apache.plc4x.java.base.model.SubscriptionPlcField;
import org.apache.plc4x.java.s7.model.S7Field;
import org.apache.plc4x.java.s7.model.S7SslField;
import org.apache.plc4x.java.s7.model.S7SubscriptionField;
import org.apache.plc4x.java.s7.netty.events.S7ConnectedEvent;
import org.apache.plc4x.java.s7.netty.model.messages.S7Message;
import org.apache.plc4x.java.s7.netty.model.messages.S7PushMessage;
import org.apache.plc4x.java.s7.netty.model.messages.S7RequestMessage;
import org.apache.plc4x.java.s7.netty.model.messages.S7ResponseMessage;
import org.apache.plc4x.java.s7.netty.model.params.CpuCyclicServicesRequestParameter;
import org.apache.plc4x.java.s7.netty.model.params.CpuServicesParameter;
import org.apache.plc4x.java.s7.netty.model.params.CpuServicesRequestParameter;
import org.apache.plc4x.java.s7.netty.model.params.CpuServicesResponseParameter;
import org.apache.plc4x.java.s7.netty.model.params.S7Parameter;
import org.apache.plc4x.java.s7.netty.model.params.VarParameter;
import org.apache.plc4x.java.s7.netty.model.params.items.S7AnyVarParameterItem;
import org.apache.plc4x.java.s7.netty.model.params.items.VarParameterItem;
import org.apache.plc4x.java.s7.netty.model.payloads.AlarmMessagePayload;
import org.apache.plc4x.java.s7.netty.model.payloads.CpuCyclicServicesRequestPayload;
import org.apache.plc4x.java.s7.netty.model.payloads.CpuCyclicServicesResponsePayload;
import org.apache.plc4x.java.s7.netty.model.payloads.CpuMessageSubscriptionServicePayload;
import org.apache.plc4x.java.s7.netty.model.payloads.CpuServicesPayload;
import org.apache.plc4x.java.s7.netty.model.payloads.S7Payload;
import org.apache.plc4x.java.s7.netty.model.payloads.VarPayload;
import org.apache.plc4x.java.s7.netty.model.payloads.items.AlarmMessageItem;
import org.apache.plc4x.java.s7.netty.model.payloads.items.AssociatedValueItem;
import org.apache.plc4x.java.s7.netty.model.payloads.items.MessageObjectItem;
import org.apache.plc4x.java.s7.netty.model.payloads.items.S7AnyVarPayloadItem;
import org.apache.plc4x.java.s7.netty.model.payloads.items.VarPayloadItem;
import org.apache.plc4x.java.s7.netty.model.types.*;
import org.apache.plc4x.java.s7.protocol.S7ByteReadResponse;
import org.apache.plc4x.java.s7.protocol.S7CyclicServicesSubscriptionHandle;
import org.apache.plc4x.java.s7.protocol.S7DiagnosticSubscriptionHandle;
import org.apache.plc4x.java.s7.utils.S7Helper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * This layer transforms between {@link PlcRequestContainer}s {@link S7Message}s.
 * And stores all "in-flight" requests in an internal structure ({@link Plc4XS7Protocol#requests}).
 * <p>
 * While sending a request, a {@link S7RequestMessage} is generated and send downstream (to the {@link S7Protocol}.
 * <p>
 * When a {@link S7ResponseMessage} is received it takes the existing request container from its Map and finishes
 * the {@link PlcRequestContainer}s future with the {@link PlcResponse}.
 */
public class Plc4XS7Protocol extends PlcMessageToMessageCodec<S7Message, PlcRequestContainer> {
    private static final Logger logger = LoggerFactory.getLogger( Plc4XS7Protocol.class );

    public static final AtomicInteger tpduGenerator = new AtomicInteger(10);
    private final AttributeKey<AtomicInteger> tpduGeneratorKey = AttributeKey.valueOf("tpduGeneratorKey");

    private Map<Short, PlcRequestContainer> requests;
    
    /* 
     * Each layer of Netty must do its work, it sends the Requests and must 
     * verify the answers.
     * Special case of the array bits, I need to know if additional messages 
     * were created.
     * Key:Tpdu, Left:Item row, Rigth: Num of Items created-1
    */
    private final  Map<Short, Map<Short, Integer>> requestBitArray;
    
    private final BlockingQueue<S7PushMessage> alarmsqueue;

    public Plc4XS7Protocol() {
        this.requests = new HashMap<>();
        this.requestBitArray = new HashMap<>();      
        this.alarmsqueue = null;
    }
    /*
    * @param s7Type
    * 
    */
    public Plc4XS7Protocol(BlockingQueue<S7PushMessage> alarmsqueue) {
        this.requests = new HashMap<>();
        this.requestBitArray = new HashMap<>();
        //We need check the device here
        this.alarmsqueue = alarmsqueue;
    }
        
    /**
     * If this protocol layer catches an {@link S7ConnectedEvent} from the protocol layer beneath,
     * the connection establishment is finished.
     *
     * @param ctx the current protocol layers context
     * @param evt the event
     * @throws Exception throws an exception if something goes wrong internally
     */
    @Override
    public void userEventTriggered(ChannelHandlerContext ctx, Object evt) throws Exception {
        ctx.channel().attr(tpduGeneratorKey).set(tpduGenerator);
        if (evt instanceof S7ConnectedEvent) {
            ctx.channel().pipeline().fireUserEventTriggered(new ConnectedEvent());
        } else {
            super.userEventTriggered(ctx, evt);
        }
    }

    /**
     * When receiving an error inside the pipeline, we have to find out which {@link PlcRequestContainer}
     * correlates needs to be notified about the problem. If a container is found, we can relay the
     * exception to that by calling completeExceptionally and passing in the exception.
     *
     * @param ctx   the current protocol layers context
     * @param cause the exception that was caught
     * @throws Exception throws an exception if something goes wrong internally
     */
    @Override
    public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) throws Exception {
        if (cause instanceof PlcProtocolPayloadTooBigException) {
            PlcProtocolPayloadTooBigException pptbe = (PlcProtocolPayloadTooBigException) cause;
            if (pptbe.getPayload() instanceof S7RequestMessage) {
                S7RequestMessage request = (S7RequestMessage) pptbe.getPayload();
                if (request.getParent() instanceof PlcRequestContainer) {
                    PlcRequestContainer requestContainer = (PlcRequestContainer) request.getParent();

                    // Remove the current request from the unconfirmed requests list.
                    requests.remove(request.getTpduReference());

                    requestContainer.getResponseFuture().completeExceptionally(cause);
                }
            }
        } else if ((cause instanceof IOException) && (cause.getMessage().contains("Connection reset by peer") ||
            cause.getMessage().contains("Operation timed out"))) {
            String reason = cause.getMessage().contains("Connection reset by peer") ?
                "Connection terminated unexpectedly" : "Remote host not responding";
            if (!requests.isEmpty()) {
                // If the connection is hung up, all still pending requests can be closed.
                for (PlcRequestContainer requestContainer : requests.values()) {
                    requestContainer.getResponseFuture().completeExceptionally(new PlcIoException(reason));
                }
                // Clear the list
                requests.clear();
            }
        } else {
            super.exceptionCaught(ctx, cause);
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Encoding
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected void encode(ChannelHandlerContext ctx, PlcRequestContainer msg, List<Object> out) throws PlcException {
        PlcRequest request = msg.getRequest();

        if (request instanceof PlcReadRequest) {
            encodeReadRequest(msg, out);
        } else if (request instanceof PlcWriteRequest) {
            encodeWriteRequest(msg, out);
        } else if (request instanceof PlcSubscriptionRequest) {
            encodeSubcriptionRequest(msg, out);
        } else if (request instanceof PlcUnsubscriptionRequest) {
            encodeUnSubcriptionRequest(msg, out);
        }
    }

    //TODO: Todo el proceso de suscripcion a eventos, suscripcion de variable y reconocimientos de alarmas
    // se ejecuta aqui. Se debe separar en funciones bien definidas.
    // Subscripcion a Eventos
    // Reconocimiento de alarmas
    // Consulta de alarmas.
    // Subscripcion a eventos ciclicos
    // 
    private void encodeSubcriptionRequest(PlcRequestContainer msg, List<Object> out) throws PlcException {
        //Chequear el campo prueba que ya incluyte la información del tipo
                
        PlcSubscriptionRequest subsRequest = (PlcSubscriptionRequest)  msg.getRequest();        

        if ( subsRequest.getFields().get(0) instanceof S7SubscriptionField){  

            S7SubscriptionField event = (S7SubscriptionField) subsRequest.getFields().get(0);

            switch(event.getFieldtype()){
                case EVENT_SUBSCRIPTION:;
                    encodeEventSubcriptionRequest(msg, out);
                    break;
                case EVENT_UNSUBSCRIPTION:;
                    encodeEventUnSubcriptionRequest(msg, out);
                    break;
                case ALARM_ACK:;
                    encodeAlarmAckRequest(msg, out);
                    break;
                case ALARM_QUERY:;
                    encodeAlarmQueryRequest(msg, out);
                    break;
                case CYCLIC_SUBSCRIPTION:;
                    encodeCycledSubscriptionRequest(msg, out);
                    break;
                case CYCLIC_UNSUBSCRIPTION:;
                    encodeCycledUnSubscriptionRequest(msg, out);
                    break;                     
                default:;                     
             };

        };            
    }
    
    //TODO: Suppport for ALARM_8
    private void encodeEventSubcriptionRequest(PlcRequestContainer msg, List<Object> out) throws PlcException {
        byte subsevent = 0;
        List<S7Parameter> parameterItems = new LinkedList<>();
        List<S7Payload> payloadItems = new LinkedList<>();
        AlarmType alarmType = AlarmType.ALARM_S_INITIATE;
         
        PlcSubscriptionRequest subsRequest = (PlcSubscriptionRequest)  msg.getRequest();
        
        for (String fieldName : subsRequest.getFieldNames()) {                        
            if ( subsRequest.getField(fieldName) instanceof S7SubscriptionField){
                 S7SubscriptionField event = (S7SubscriptionField) subsRequest.getField(fieldName);
                 subsevent = (byte) (subsevent | event.getEventtype().getCode());
                 if (event.getEventtype() == SubscribedEventType.ALM_8) {
                    alarmType = AlarmType.ALARM_INITIATE;
                 }
            }
        }
        
        CpuServicesParameter cpuservice = new CpuServicesRequestParameter(CpuServicesParameterFunctionGroup.CPU_FUNCTIONS,
            CpuServicesParameterSubFunctionGroup.MESSAGE_SERVICE, (byte) 0x00);
        
        parameterItems.add(cpuservice);
        
       
        S7Payload Data = new CpuMessageSubscriptionServicePayload(DataTransportErrorCode.OK,
                                DataTransportSize.OCTET_STRING,
                                subsevent,
                                new String("HmiRtm  "),
                                alarmType);
  
        
        payloadItems.add(Data);
      

        
        S7RequestMessage s7ReadRequest = new S7RequestMessage(MessageType.USER_DATA,
            (short) tpduGenerator.getAndIncrement(), 
            parameterItems,
            payloadItems, 
            msg);

        requests.put(s7ReadRequest.getTpduReference(), msg);

        out.add(s7ReadRequest);         
    }    
    
 
    //TODO: Check for ALARM_D and ALARM_8
    private void encodeAlarmAckRequest(PlcRequestContainer msg, List<Object> out) throws PlcException {
        byte subsevent = 0;
        List<S7Parameter> parameterItems = new LinkedList<>();
        List<S7Payload> payloadItems = new LinkedList<>();
         
        PlcSubscriptionRequest subsRequest = (PlcSubscriptionRequest)  msg.getRequest();
        S7SubscriptionField event = (S7SubscriptionField) subsRequest.getFields().get(0);

        
        CpuServicesParameter cpuservice = new CpuServicesRequestParameter(CpuServicesParameterFunctionGroup.CPU_FUNCTIONS,
            CpuServicesParameterSubFunctionGroup.ALARM_ACK, (byte) 0x00);

        parameterItems.add(cpuservice);
        
        
        ArrayList<Object> items = new ArrayList<>();
        for(Integer eventid:event.getAckalarms()){
            items.add(new MessageObjectItem((byte) 0x12,
                            (byte) 0x08,
                            VariableAddressingMode.ALARM_ACK, 
                            (byte) 0x00,
                            eventid,
                            (byte) 0x01,
                            (byte) 0x01)); //ONLY for Alarm_SQ
        }

        AlarmMessageItem messageitem = new  AlarmMessageItem((byte) 0x09,
                                                (byte) items.size(),
                                                items)     ;
        
        S7Payload Data = new AlarmMessagePayload(DataTransportErrorCode.OK,
                                DataTransportSize.OCTET_STRING,
                                (2+items.size()*10),
                                messageitem);
        payloadItems.add(Data);

        
        S7RequestMessage s7ReadRequest = new S7RequestMessage(MessageType.USER_DATA,
            (short) tpduGenerator.getAndIncrement(), 
            parameterItems,
            payloadItems, 
            msg);

        requests.put(s7ReadRequest.getTpduReference(), msg);

        out.add(s7ReadRequest);         
    }   
    
    //TODO: Check for ALARM_D and ALARM_8
    private void encodeAlarmQueryRequest(PlcRequestContainer msg, List<Object> out) throws PlcException {

        byte subsevent = 0;
        List<S7Parameter> parameterItems = new LinkedList<>();
        List<S7Payload> payloadItems = new LinkedList<>();
         
        PlcSubscriptionRequest subsRequest = (PlcSubscriptionRequest)  msg.getRequest();
        S7SubscriptionField event = (S7SubscriptionField) subsRequest.getFields().get(0);
        
        CpuServicesParameter cpuservice = new CpuServicesRequestParameter(CpuServicesParameterFunctionGroup.CPU_FUNCTIONS,
            CpuServicesParameterSubFunctionGroup.ALARM_QUERY, (byte) 0x00);

        parameterItems.add(cpuservice);
        
        ArrayList<Object> items = new ArrayList<>();
        
        //TODO: Multiple AlarmQueryType from the same request.
        items.add(new MessageObjectItem((byte) 0x12,
                        (byte) 0x08,
                        VariableAddressingMode.ALARM_QUERYREQ, 
                        QueryType.BYALARMTYPE,
                        event.getAlarmquerytype())); //ONLY for Alarm_S


        AlarmMessageItem messageitem = new  AlarmMessageItem((byte) 0x00,
                                                (byte) items.size(),
                                                items)     ;
        
        S7Payload Data = new AlarmMessagePayload(DataTransportErrorCode.OK,
                                DataTransportSize.OCTET_STRING,
                                (2+items.size()*10),
                                messageitem);
        payloadItems.add(Data);


        S7RequestMessage s7ReadRequest = new S7RequestMessage(MessageType.USER_DATA,
            (short) tpduGenerator.getAndIncrement(), 
            parameterItems,
            payloadItems, 
            msg);

        requests.put(s7ReadRequest.getTpduReference(), msg);

        out.add(s7ReadRequest);         
    }    
    
    private void encodeCycledSubscriptionRequest(PlcRequestContainer msg, List<Object> out) throws PlcException {
        List<S7Parameter> parameterItems = new LinkedList<>();
        List<S7Payload> payloads = new LinkedList<>();        
        DefaultPlcSubscriptionRequest prueba = (DefaultPlcSubscriptionRequest)  msg.getRequest();
        List<SubscriptionPlcField> fields = prueba.getSubscriptionFields();
        Optional<Duration> optionalduration = prueba.getSubscriptionFields().get(0).getDuration();
        Duration duration  = optionalduration.get();
        
        List<S7AnyVarPayloadItem> payloadItems = new LinkedList<>();
        
        fields.forEach((action) ->{
            PlcField field = action.getPlcField();
            if (!(field instanceof S7SubscriptionField)) {
                logger.info("The field should have been of type S7SubscriptionField");
                return;
            }  
            S7SubscriptionField subsS7Field = (S7SubscriptionField) field; 
            S7Field s7Field = subsS7Field.getS7field();
            S7AnyVarPayloadItem varPayloadItem = new S7AnyVarPayloadItem(
                SpecificationType.VARIABLE_SPECIFICATION, s7Field.getMemoryArea(),
                s7Field.getDataType(),
                s7Field.getNumElements(), s7Field.getBlockNumber(), s7Field.getByteOffset(), (byte) s7Field.getBitOffset());
            payloadItems.add(varPayloadItem);            
        });
        
        //TODO: Play with Time factor. Actually only 1..9
        CpuCyclicServiceTimeBaseType timeBase = null;
        byte timefactor = 0x00;
        if (duration.getSeconds() == 0) {
            timeBase = CpuCyclicServiceTimeBaseType.TB_100_MS;
            timefactor = (byte) (duration.toMillis() / 100);
        } else if (duration.getSeconds() <= 9) {
            timeBase = CpuCyclicServiceTimeBaseType.TB_1_SEC;
            timefactor = (byte) duration.getSeconds();            
        } else if (duration.getSeconds() <= 90) {
            timeBase = CpuCyclicServiceTimeBaseType.TB_10_SEC;
            timefactor = (byte) (duration.getSeconds() / 10);  
        }
        
        CpuCyclicServicesRequestPayload payload = new CpuCyclicServicesRequestPayload(DataTransportErrorCode.OK,
                                                        DataTransportSize.OCTET_STRING,
                                                        (4+payloadItems.size()*12),
                                                        payloadItems.size(),
                                                        timeBase,
                                                        timefactor,
                                                        payloadItems);
                                                        
        payloads.add(payload);
        
        CpuCyclicServicesRequestParameter parameter = new CpuCyclicServicesRequestParameter(CpuUserDataMethodType.REQUEST,
                                                            CpuServicesParameterFunctionGroup.CYCLIC_SERVICES,
                                                            CpuCyclicServicesParameterSubFunctionGroupType.CYCLIC_TRANSFER,
                                                            (byte) 0x00);
        parameterItems.add(parameter);
        

        
        // Assemble the request.
        S7RequestMessage s7ReadRequest = new S7RequestMessage(MessageType.USER_DATA,
            (short) tpduGenerator.getAndIncrement(), parameterItems,
            payloads, msg);

        requests.put(s7ReadRequest.getTpduReference(), msg);

        out.add(s7ReadRequest);        
    }
    
    private void encodeUnSubcriptionRequest(PlcRequestContainer msg, List<Object> out) throws PlcException {
        //Chequear el campo prueba que ya incluyte la información del tipo
        InternalPlcUnsubscriptionRequest subsRequest = (InternalPlcUnsubscriptionRequest)  msg.getRequest();  
        Collection<? extends InternalPlcSubscriptionHandle> handles = subsRequest.getInternalPlcSubscriptionHandles();
        handles.forEach((handle)->{
            if (handle instanceof S7CyclicServicesSubscriptionHandle) {
                try {
                    encodeCycledUnSubscriptionRequest(msg,out);
                } catch (PlcException ex) {
                    logger.info(ex.toString());
                }
            }
        });

    } 
    
    private void encodeEventUnSubcriptionRequest(PlcRequestContainer msg, List<Object> out) throws PlcException {
        byte subsevent = 0;
        List<S7Parameter> parameterItems = new LinkedList<>();
        List<S7Payload> payloadItems = new LinkedList<>();
        AlarmType alarmType = AlarmType.ALARM_S_INITIATE;
         
        PlcSubscriptionRequest subsRequest = (PlcSubscriptionRequest)  msg.getRequest();
        
        for (String fieldName : subsRequest.getFieldNames()) {                        
            if ( subsRequest.getField(fieldName) instanceof S7SubscriptionField){
                 S7SubscriptionField event = (S7SubscriptionField) subsRequest.getField(fieldName);
                 subsevent = (byte) (subsevent | event.getEventtype().getCode());
                 if (event.getEventtype() == SubscribedEventType.ALM_8) {
                    alarmType = AlarmType.ALARM_INITIATE;
                 }                 
            }
        }
        
        CpuServicesParameter cpuservice = new CpuServicesRequestParameter(CpuServicesParameterFunctionGroup.CPU_FUNCTIONS,
            CpuServicesParameterSubFunctionGroup.MESSAGE_SERVICE, (byte) 0x00);
        
        parameterItems.add(cpuservice);
        
       
        S7Payload Data = new CpuMessageSubscriptionServicePayload(DataTransportErrorCode.OK,
                                DataTransportSize.OCTET_STRING,
                                subsevent,
                                new String("HmiRtm  "),
                                alarmType);
        
        payloadItems.add(Data);

        
        S7RequestMessage s7ReadRequest = new S7RequestMessage(MessageType.USER_DATA,
            (short) tpduGenerator.getAndIncrement(), 
            parameterItems,
            payloadItems, 
            msg);

        requests.put(s7ReadRequest.getTpduReference(), msg);

        out.add(s7ReadRequest);         
    }  
       
    private void encodeCycledUnSubscriptionRequest(PlcRequestContainer msg, List<Object> out) throws PlcException {
 
        List<S7Parameter> parameterItems = new LinkedList<>();
        List<S7Payload> payloads = new LinkedList<>();
        byte jobId = 0x00;
        
        DefaultPlcUnsubscriptionRequest prueba = (DefaultPlcUnsubscriptionRequest)  msg.getRequest();
   
        Collection<? extends InternalPlcSubscriptionHandle> handles = prueba.getInternalPlcSubscriptionHandles();
 
        for(InternalPlcSubscriptionHandle ihandle:handles){
            S7CyclicServicesSubscriptionHandle s7handle = (S7CyclicServicesSubscriptionHandle) ihandle;
            jobId = s7handle.getJobId();
            break;
        } 
 

        CpuCyclicServicesRequestPayload payload = new CpuCyclicServicesRequestPayload(DataTransportErrorCode.OK,
                                                        DataTransportSize.OCTET_STRING,
                                                        (byte) 0x02,
                                                        (byte) 0x05,
                                                        jobId);
                                                        
        payloads.add(payload);

        CpuCyclicServicesRequestParameter parameter = new CpuCyclicServicesRequestParameter(CpuUserDataMethodType.REQUEST,
                                                            CpuServicesParameterFunctionGroup.CYCLIC_SERVICES,
                                                            CpuCyclicServicesParameterSubFunctionGroupType.CYCLIC_UNSUBSCRIBE,
                                                            (byte) 0x00);
        parameterItems.add(parameter);
        

        
        // Assemble the request.
        S7RequestMessage s7ReadRequest = new S7RequestMessage(MessageType.USER_DATA,
            (short) tpduGenerator.getAndIncrement(), parameterItems,
            payloads, msg);

        requests.put(s7ReadRequest.getTpduReference(), msg);

        out.add(s7ReadRequest);
       
    }
        
    private void encodeReadRequest(PlcRequestContainer msg, List<Object> out) throws PlcException {

        List<VarParameterItem> parameterItems = new LinkedList<>();

        PlcReadRequest readRequest = (PlcReadRequest) msg.getRequest();
        
        Short itemRow = 0;
        Short tpdu = (short) tpduGenerator.getAndIncrement();
        
        for (String fieldName : readRequest.getFieldNames()) {
            PlcField field = readRequest.getField(fieldName);
            
            if ((field instanceof S7SslField)) {
                S7SslField szlfield = (S7SslField) field;
                encodeSslReadRequest(msg,szlfield,out);
                return;
            } else {
                if (!(field instanceof S7Field)) {
                    throw new PlcProtocolException("The field should have been of type S7Field");
                }

                S7Field s7Field = (S7Field) field;
                if(!(readRequest instanceof DefaultPlcReadRequest)) {
                    throw new PlcException("The writeRequest should have been of type DefaultPlcWriteRequest");
                }
                
                if (s7Field.getDataType() == TransportSize.BOOL){
                    encodeWriteRequestBitFieldAndArray( tpdu,
                            itemRow,
                            null,
                            s7Field, 
                            parameterItems, 
                            Collections.EMPTY_LIST);                                         
                } else {
                
                    VarParameterItem varParameterItem = new S7AnyVarParameterItem(
                        SpecificationType.VARIABLE_SPECIFICATION, s7Field.getMemoryArea(),
                        s7Field.getDataType(),
                        s7Field.getNumElements(), s7Field.getBlockNumber(), s7Field.getByteOffset(), (byte) s7Field.getBitOffset());
                    parameterItems.add(varParameterItem);
                }
                

            }
            itemRow++;
        }
        
        VarParameter readVarParameter = new VarParameter(ParameterType.READ_VAR, parameterItems);

        // Assemble the request.
        S7RequestMessage s7ReadRequest = new S7RequestMessage(MessageType.JOB,
            tpdu, Collections.singletonList(readVarParameter),
            Collections.emptyList(), msg);

        requests.put(s7ReadRequest.getTpduReference(), msg);

        out.add(s7ReadRequest);
    }

    private void encodeSslReadRequest(PlcRequestContainer msg, S7SslField sslfield, List<Object> out) throws PlcException {
        
        S7RequestMessage szlRequestMessage = new S7RequestMessage(MessageType.USER_DATA,
                (short) tpduGenerator.getAndIncrement(),
            Collections.singletonList(new CpuServicesRequestParameter(
                CpuServicesParameterFunctionGroup.CPU_FUNCTIONS,
                CpuServicesParameterSubFunctionGroup.READ_SSL, (byte) 0)),
            Collections.singletonList(new CpuServicesPayload(DataTransportErrorCode.OK, SslId.valueOf((short) (sslfield.getSslId() & 0x01FF)),
                (short) sslfield.getIndex())), null); 
        
        requests.put(szlRequestMessage.getTpduReference(), msg);

        out.add(szlRequestMessage);        
    }    
    
    // TODO: 
    private void encodeWriteRequest(PlcRequestContainer msg, List<Object> out) throws PlcException {
        List<VarParameterItem> parameterItems = new LinkedList<>();
        List<VarPayloadItem> payloadItems = new LinkedList<>();

        PlcWriteRequest writeRequest = (PlcWriteRequest) msg.getRequest();
        
        Short itemRow = 0;
        Short tpdu = (short) tpduGenerator.getAndIncrement();
        for (String fieldName : writeRequest.getFieldNames()) {
            PlcField field = writeRequest.getField(fieldName);
            if (!(field instanceof S7Field)) {
                throw new PlcException("The field should have been of type S7Field");
            }
            S7Field s7Field = (S7Field) field;
            if(!(writeRequest instanceof DefaultPlcWriteRequest)) {
                throw new PlcException("The writeRequest should have been of type DefaultPlcWriteRequest");
            }
            BaseDefaultFieldItem fieldItem = ((DefaultPlcWriteRequest) writeRequest).getFieldItem(fieldName);

            // The number of elements provided in the request must match the number defined in the field, or
            // bad things are going to happen.
            // An exception is STRINGS, as they are implemented as byte arrays
            // TODO: Calculate the number of correct elements for a "STRING", 
            //       since it will be written to the PLC as a string of "CHAR" elements.
            if (s7Field.getDataType() != TransportSize.STRING &&
                writeRequest.getNumberOfValues(fieldName) != s7Field.getNumElements()) {
                throw new PlcException("The number of values provided doesn't match the number specified by the field.");
            }
            
            //TODO: Bit Arrays: One item request for every bit
            
            if (s7Field.getDataType() == TransportSize.BOOL) {
                
                encodeWriteRequestBitFieldAndArray( tpdu,
                        itemRow,
                        fieldItem,
                        s7Field, 
                        parameterItems, 
                        payloadItems);            
                
            } else {

                VarParameterItem varParameterItem = new S7AnyVarParameterItem(
                    SpecificationType.VARIABLE_SPECIFICATION, s7Field.getMemoryArea(),
                    s7Field.getDataType(),
                    (s7Field.getDataType() == TransportSize.STRING)?fieldItem.getString(0).length()+2:s7Field.getNumElements(), 
                    s7Field.getBlockNumber(), 
                    s7Field.getByteOffset(),
                    (byte) s7Field.getBitOffset());

                parameterItems.add(varParameterItem);

                DataTransportSize dataTransportSize = s7Field.getDataType().getDataTransportSize();

                // TODO: Checkout if the payload items are sort of a flatMap of all request items.
                byte[] byteData = null;
                switch(s7Field.getDataType()) {
                    // -----------------------------------------
                    // Bit
                    // -----------------------------------------
                    case BOOL:
                        //byteData = encodeWriteRequestBitField(fieldItem);
                        break;
                    // -----------------------------------------
                    // Signed integer values
                    // -----------------------------------------
                    case BYTE:
                    case SINT:
                    case CHAR:  // 1 byte
                        byteData = encodeWriteRequestByteField(fieldItem, true);
                        break;                
                    case WORD:
                    case INT:
                    case WCHAR:  // 2 byte (16 bit)
                        byteData = encodeWriteRequestShortField(fieldItem, true);
                        break;
                    case DWORD:
                    case DINT:  // 4 byte (32 bit)
                        byteData = encodeWriteRequestIntegerField(fieldItem, true);
                        break;
                    case LWORD:
                    case LINT:  // 8 byte (64 bit)
                        byteData = encodeWriteRequestLongField(fieldItem, true);
                        break;
                    // -----------------------------------------
                    // Unsigned integer values
                    // -----------------------------------------
                    // 8 bit:
                    case USINT:
                        byteData = encodeWriteRequestByteField(fieldItem, false);
                        break;
                    // 16 bit:
                    case UINT:
                        byteData = encodeWriteRequestShortField(fieldItem, false);
                        break;
                    // 32 bit:
                    case UDINT:
                        byteData = encodeWriteRequestIntegerField(fieldItem, false);
                        break;
                    // 64 bit:
                    case ULINT:
                        byteData = encodeWriteRequestLongField(fieldItem, false);
                        break;
                    // -----------------------------------------
                    // Floating point values
                    // -----------------------------------------
                    case REAL:
                        byteData = encodeWriteRequestFloatField(fieldItem);
                        break;
                    case LREAL:
                        byteData = encodeWriteRequestDoubleField(fieldItem);
                        break;
                    // -----------------------------------------
                    // Characters & Strings
                    // -----------------------------------------
                    case STRING:
                        byteData = encodeWriteRequestStringField(fieldItem, false);
                        break;
                    case WSTRING:
                        byteData = encodeWriteRequestStringField(fieldItem, true);
                        break;
                    // -----------------------------------------
                    // Date & Time
                    // -----------------------------------------
                    case S5TIME:
                        byteData = encodeWriteRequestS5TimeField(fieldItem, false);
                        break;
                    case TIME:  // 4 byte (32 bit)
                        byteData = encodeWriteRequestTimeField(fieldItem, false);
                        break;
                    case TOD:
                    case TIME_OF_DAY:
                        byteData = encodeWriteRequestTodField(fieldItem, true);
                        break;
                    case DATE:  // 2 byte (16 bit)
                        byteData = encodeWriteRequestDateField(fieldItem, true);
                        break; 
                    case DT:                        
                    case DATE_AND_TIME:  // 8 byte (64 bit)
                        byteData = encodeWriteRequestDateTimeField(fieldItem, true);
                        break;                    
                    default:
                        throw new PlcProtocolException("Unsupported type.. " + s7Field.getDataType());
                }

                VarPayloadItem varPayloadItem = new VarPayloadItem(
                DataTransportErrorCode.RESERVED, dataTransportSize, byteData);

                payloadItems.add(varPayloadItem);
            }
            itemRow++;
        };
        
        VarParameter writeVarParameter = new VarParameter(ParameterType.WRITE_VAR, parameterItems);
        VarPayload writeVarPayload = new VarPayload(ParameterType.WRITE_VAR, payloadItems);

        // Assemble the request.
        S7RequestMessage s7WriteRequest = new S7RequestMessage(MessageType.JOB,
            tpdu, Collections.singletonList(writeVarParameter),
            Collections.singletonList(writeVarPayload), msg);

        requests.put(s7WriteRequest.getTpduReference(), msg);

        out.add(s7WriteRequest);
    }
    
    
    void encodeWriteRequestBitFieldAndArray(Short tpdu,
            Short itemRow,
            BaseDefaultFieldItem fieldItem,
            S7Field s7Field, 
            List<VarParameterItem> parameterItems, 
            List<VarPayloadItem> payloadItems){
        
        int curParameterByteOffset = s7Field.getByteOffset();
        byte curParameterBitOffset = (byte) s7Field.getBitOffset();
        byte curPayloadBitOffset = 0; 
        
        //Only BITS arrays
        if (s7Field.getNumElements() > 1) {
            requestBitArray.computeIfAbsent(tpdu, 
                    k->new HashMap()).put(itemRow,s7Field.getNumElements()-1);
        }

        for (int i = 0; i < s7Field.getNumElements(); i++) {

            VarParameterItem varParameterItem = new S7AnyVarParameterItem(
                SpecificationType.VARIABLE_SPECIFICATION, s7Field.getMemoryArea(),
                s7Field.getDataType(),
                (short) 1, 
                s7Field.getBlockNumber(), 
                curParameterByteOffset,
                curParameterBitOffset);

            // Create a one-byte byte array and set it to 0x01, if the corresponding bit
            // was 1 else set it to 0x00.
            byte[] elementData = new byte[1];
                                   
            parameterItems.add(varParameterItem);
            
            if (!(payloadItems == Collections.EMPTY_LIST)){
                elementData[0] = (byte) (fieldItem.getBoolean(i)?0x01:0x00);                  
                VarPayloadItem itemPayload = new VarPayloadItem(
                    DataTransportErrorCode.RESERVED,
                    s7Field.getDataType().getDataTransportSize(),
                        elementData);            
                payloadItems.add(itemPayload);
            }
            
            // Calculate the new memory and bit offset.
            curParameterBitOffset++;
            if ((i > 0) && ((curParameterBitOffset % 8) == 0)) {
                curParameterByteOffset++;
                curParameterBitOffset = 0;
            }
            curPayloadBitOffset++;
            if ((i > 0) && ((curPayloadBitOffset % 8) == 0)) {
                curPayloadBitOffset = 0;
            }                                    
        }            
    }
    

    byte[] encodeWriteRequestBitField(BaseDefaultFieldItem fieldItem) {
        int numBytes = fieldItem.getNumberOfValues() >> 3 / 8;
        byte[] byteData = new byte[numBytes];
        BitSet bitSet = new BitSet();
        for (int i = 0; i < fieldItem.getNumberOfValues(); i++) {
            bitSet.set(i, fieldItem.getBoolean(i));
        }
        byte[] src = bitSet.toByteArray();
        System.arraycopy(src, 0, byteData, 0, Math.min(src.length, numBytes));
        return byteData;
    }

    byte[] encodeWriteRequestByteField(BaseDefaultFieldItem fieldItem, boolean signed) {
        int numBytes = fieldItem.getNumberOfValues();
        ByteBuffer buffer = ByteBuffer.allocate(numBytes);
        for (int i = 0; i < fieldItem.getNumberOfValues(); i++) {
            if(signed) {
                buffer.put(fieldItem.getByte(i));
            } else {
                buffer.put((byte) (short) fieldItem.getShort(i));
            }
        }
        return buffer.array();
    }

    byte[] encodeWriteRequestShortField(BaseDefaultFieldItem fieldItem, boolean signed) {
        int numBytes = fieldItem.getNumberOfValues() * 2;
        ByteBuffer buffer = ByteBuffer.allocate(numBytes);
        for (int i = 0; i < fieldItem.getNumberOfValues(); i++) {
            if(signed) {
                buffer.putShort(fieldItem.getShort(i));
            } else {
                buffer.putShort((short) (int) fieldItem.getInteger(i));
            }
        }
        return buffer.array();
    }

    byte[] encodeWriteRequestIntegerField(BaseDefaultFieldItem fieldItem, boolean signed) {
        int numBytes = fieldItem.getNumberOfValues() * 4;
        ByteBuffer buffer = ByteBuffer.allocate(numBytes);
        for (int i = 0; i < fieldItem.getNumberOfValues(); i++) {
            if(signed) {
                buffer.putInt(fieldItem.getInteger(i));
            } else {
                buffer.putInt((int) (long) fieldItem.getLong(i));
            }
        }
        return buffer.array();
    }

    byte[] encodeWriteRequestLongField(BaseDefaultFieldItem fieldItem, boolean signed) {
        int numBytes = fieldItem.getNumberOfValues() * 8;
        ByteBuffer buffer = ByteBuffer.allocate(numBytes);
        for (int i = 0; i < fieldItem.getNumberOfValues(); i++) {
            if(signed) {
                buffer.putLong(fieldItem.getLong(i));
            } else {
                // TODO: Implement this ...
            }
        }
        return buffer.array();
    }

    byte[] encodeWriteRequestFloatField(BaseDefaultFieldItem fieldItem) {
        int numBytes = fieldItem.getNumberOfValues() * 4;
        ByteBuffer buffer = ByteBuffer.allocate(numBytes);
        for (int i = 0; i < fieldItem.getNumberOfValues(); i++) {
            buffer.putFloat(fieldItem.getFloat(i));
        }
        return buffer.array();
    }

    byte[] encodeWriteRequestDoubleField(BaseDefaultFieldItem fieldItem) {
        int numBytes = fieldItem.getNumberOfValues() * 8;
        ByteBuffer buffer = ByteBuffer.allocate(numBytes);
        for (int i = 0; i < fieldItem.getNumberOfValues(); i++) {
            buffer.putDouble(fieldItem.getDouble(i));
        }
        return buffer.array();
    }
    
    byte[] encodeWriteRequestStringField(BaseDefaultFieldItem fieldItem, boolean isUtf16) {
        // TODO: Implement this ...
        String str = fieldItem.getString(0);
        ByteBuffer buffer = ByteBuffer.allocate(str.length()+2);
        buffer.put((byte)254);
        buffer.put((byte)str.length());
        buffer.put(str.getBytes());
        return buffer.array();
    }
    
    byte[] encodeWriteRequestS5TimeField(BaseDefaultFieldItem fieldItem, boolean signed) {
        int numBytes = fieldItem.getNumberOfValues() * 2;
        ByteBuffer buffer = ByteBuffer.allocate(numBytes);
        for (int i = 0; i < fieldItem.getNumberOfValues(); i++) {
            buffer.putShort(S7Helper.DurationToS5Time(fieldItem.getDuration(i)));
        }
        return buffer.array();
    }    
    
    byte[] encodeWriteRequestTimeField(BaseDefaultFieldItem fieldItem, boolean signed) {
        int numBytes = fieldItem.getNumberOfValues() * 4;
        ByteBuffer buffer = ByteBuffer.allocate(numBytes);
        for (int i = 0; i < fieldItem.getNumberOfValues(); i++) {
            buffer.putInt(S7Helper.DurationToS7Time(fieldItem.getDuration(i)));
        }
        return buffer.array();
    } 
    
    byte[] encodeWriteRequestTodField(BaseDefaultFieldItem fieldItem, boolean signed) {
        int numBytes = fieldItem.getNumberOfValues() * 4;
        ByteBuffer buffer = ByteBuffer.allocate(numBytes);
        for (int i = 0; i < fieldItem.getNumberOfValues(); i++) {
            buffer.putInt(S7Helper.LocalTimeToS7Tod(fieldItem.getTime(i)));
        }
        return buffer.array();
    }
    
     byte[] encodeWriteRequestDateField(BaseDefaultFieldItem fieldItem, boolean signed) {
        int numBytes = fieldItem.getNumberOfValues() * 2;
        ByteBuffer buffer = ByteBuffer.allocate(numBytes);
        for (int i = 0; i < fieldItem.getNumberOfValues(); i++) {
            buffer.putShort(S7Helper.LocalDateToS7Date(fieldItem.getDate(i)));
        }
        return buffer.array();
    }  
     
     byte[] encodeWriteRequestDateTimeField(BaseDefaultFieldItem fieldItem, boolean signed) {
        int numBytes = fieldItem.getNumberOfValues() * 8;
        ByteBuffer buffer = ByteBuffer.allocate(numBytes);
        for (int i = 0; i < fieldItem.getNumberOfValues(); i++) {
            buffer.put(S7Helper.LocalDateTimeToS7DateTime(fieldItem.getDateTime(i)));
        }
        return buffer.array();
    }      
    
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Decoding
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @SuppressWarnings("unchecked")
    @Override
    protected void decode(ChannelHandlerContext ctx, S7Message msg, List<Object> out) throws PlcException {
        // We're currently just expecting responses.

        if (!(msg instanceof S7ResponseMessage)) {
            logger.info("Aborta decode: " + msg.getMessageType());
            return;
        }
        
        S7ResponseMessage responseMessage = (S7ResponseMessage) msg;
        short tpduReference = responseMessage.getTpduReference();

        if (requests.containsKey(tpduReference)) {
            // As every response has a matching request, get this request based on the tpdu.
            PlcRequestContainer requestContainer = requests.remove(tpduReference);
            PlcRequest request = requestContainer.getRequest();

            PlcResponse response = null;
            if (request instanceof PlcReadRequest) {
                response = decodeReadResponse(responseMessage, requestContainer);
                //If responseMessage is SSL type and response is null
                //At moments Alarms and SSL request must be fragmented
                if (response == null){ //Fragmented message? Query again!
                    if (responseMessage.getPayloads().get(0) instanceof CpuServicesPayload){
                        requests.put(tpduReference, requestContainer);
                        S7RequestMessage sslRequestMessage = new S7RequestMessage(MessageType.USER_DATA,
                                msg.getTpduReference(),
                            Collections.singletonList(new  CpuServicesResponseParameter(
                                CpuServicesParameterFunctionGroup.CPU_FUNCTIONS,
                                CpuServicesParameterSubFunctionGroup.READ_SSL, 
                                ((CpuServicesParameter) msg.getParameters().get(0)).getSequenceNumber())),
                            Collections.singletonList(new CpuServicesPayload(DataTransportErrorCode.NOT_FOUND, 
                                    SslId.NULL,
                                (short) 0x0000)), null); 
                        try {
                            ChannelFuture future = ctx.writeAndFlush(sslRequestMessage);
                            future.addListener(new ChannelFutureListener() {
                                public void operationComplete(ChannelFuture future) {
                                    logger.info("Request fragment READ_SSL done: " + future.isSuccess());
                                }
                            });                        
                        } catch (Exception ex) {
                            java.util.logging.Logger.getLogger(Plc4XS7Protocol.class.getName()).log(Level.SEVERE, null, ex);
                        }
                    }
                }
                
            } else if (request instanceof PlcWriteRequest) {
                response = decodeWriteResponse(responseMessage, requestContainer);
            } else if (request instanceof PlcSubscriptionRequest) {
                response = decodeSubscriptionResponse(responseMessage, requestContainer);
                
                //At moments Alarms and SSL request must be fragmented
                if (response == null){ //Fragmented message? Query again!
                    requests.put(tpduReference, requestContainer);
                    List<S7Parameter> parameterItems = new LinkedList<>();
                    List<S7Payload> payloadItems = new LinkedList<>();
                    
                    AlarmMessagePayload payload = new AlarmMessagePayload(DataTransportErrorCode.NOT_FOUND,
                                                                            DataTransportSize.NULL,
                                                                            0x0000,
                                                                            null);
                    payloadItems.add(payload);

                    //TODO: Check multiple parameters
                    CpuServicesParameter cpuservice = new CpuServicesResponseParameter(CpuServicesParameterFunctionGroup.CPU_FUNCTIONS,
                    CpuServicesParameterSubFunctionGroup.ALARM_QUERY, ((CpuServicesParameter) msg.getParameters().get(0)).getSequenceNumber());

                    parameterItems.add(cpuservice); 

                    S7RequestMessage s7ReadRequest = new S7RequestMessage(MessageType.USER_DATA,
                        msg.getTpduReference(), 
                        parameterItems,
                        payloadItems, 
                        null);  

                    try {
                        ChannelFuture future = ctx.writeAndFlush(s7ReadRequest);
                        future.addListener(new ChannelFutureListener() {
                            public void operationComplete(ChannelFuture future) {
                                logger.debug("Request fragment ALARM_QUERY done: " + future.isSuccess());
                            }
                        });                        
                        //out.add(s7ReadRequest);
                        //ctx.pipeline().
                    } catch (Exception ex) {
                        java.util.logging.Logger.getLogger(Plc4XS7Protocol.class.getName()).log(Level.SEVERE, null, ex);
                    }
                    
                }
            } else if (request instanceof PlcUnsubscriptionRequest) {
                response = decodeUnSubscriptionResponse(responseMessage, requestContainer);
            } else {
                logger.info("There is the client's request, but it is not a valid response....");
            }

            // Confirm the response being handled.
            if (response != null) {
                requestContainer.getResponseFuture().complete(response);
            } else {
                logger.info("The message could not be processed...");
            }
        } else {
            //PUSH Message
            List<S7Parameter> parameters = msg.getParameters();
            for (S7Parameter parameter:parameters){
                if (parameter instanceof S7PushMessage){
                    if (!alarmsqueue.offer((S7PushMessage) parameter)){
                        logger.info("Event queue buffer is full.");
                    };                    
                }    
            }
            
            List<S7Payload> payloads = msg.getPayloads();  
            //TODO: Use alarmsqueue.addAll() method using streams
            for (S7Payload payload:payloads){
                if (payload instanceof S7PushMessage){
                    if (!alarmsqueue.offer((S7PushMessage) payload)){
                        logger.info("Event queue buffer is full.");
                    };
                }
            }            
        }
    }
    
    // TODO: Check if it is a CPU service response and proceed as applicable.
    @SuppressWarnings("unchecked")
    private PlcResponse decodeReadResponse(S7ResponseMessage responseMessage, PlcRequestContainer requestContainer) throws PlcProtocolException {
        InternalPlcReadRequest plcReadRequest = (InternalPlcReadRequest) requestContainer.getRequest();
        Map<Short, Integer> splitItems  = null;
                
        if (responseMessage.getPayloads().get(0) instanceof CpuServicesPayload){
            return decodeSslReadResponse(responseMessage, requestContainer);
        }
        
        VarPayload payload = responseMessage.getPayload(VarPayload.class)
            .orElseThrow(() -> new PlcProtocolException("No VarPayload supplied"));

        // TODO: Checkout if the payload items are sort of a flatMap of all request items.

        // If the numbers of items don't match, we're in big trouble as the only
        // way to know how to interpret the responses is by aligning them with the
        // items from the request as this information is not returned by the PLC.

        short tpdu = responseMessage.getTpduReference();
        
        int[] fieldOffset= new int[1];
        fieldOffset[0] = 0;
        if (requestBitArray.containsKey(tpdu)){
            splitItems = requestBitArray.remove(tpdu);
            splitItems.forEach((key,length) -> {
                fieldOffset[0] += length;
            });
        }
        
        if ((plcReadRequest.getNumberOfFields()+fieldOffset[0]) != payload.getItems().size()) {
            logger.info("Excepcion: " + plcReadRequest.getFieldNames().toString());
            throw new PlcProtocolException(
                "The number of requested items doesn't match the number of returned items Request: " + 
                        (plcReadRequest.getNumberOfFields()+fieldOffset[0]) +
                        " Payload: " + 
                        payload.getItems().size());
        }


        
        Map<String, Pair<PlcResponseCode, BaseDefaultFieldItem>> values = new HashMap<>();
        List<VarPayloadItem> payloadItems = payload.getItems();
        short index = 0;
        short itemRow = 0;
        for (String fieldName : plcReadRequest.getFieldNames()) {
            S7Field field = (S7Field) plcReadRequest.getField(fieldName);
            VarPayloadItem payloadItem = payloadItems.get(itemRow);

            PlcResponseCode responseCode = decodeResponseCode(payloadItem.getReturnCode());
            BaseDefaultFieldItem fieldItem = null;

            //ByteBuf data = Unpooled.copiedBuffer(payloadItem.getData());
            ByteBuf data = PooledByteBufAllocator.DEFAULT.buffer();
            if (payloadItem.getData() != null){
                data.writeBytes(payloadItem.getData());
            }
            
            if (responseCode == PlcResponseCode.OK) {
                try {
                    switch (field.getDataType()) {
                        // -----------------------------------------
                        // Bit
                        // 1 responseMessage
                        // -----------------------------------------
                        case BOOL:
                            if (splitItems != null){
                                if (splitItems.containsKey(index)){
                                    short endRow = (short)(itemRow + splitItems.get(index));
                                    for(short i=itemRow; i<=endRow; i++){
                                        payloadItem = payloadItems.get(i);
                                        data.writeBytes(payloadItem.getData());
                                    }
                                    itemRow = endRow;
                                }
                            }
                            fieldItem = decodeReadResponseBitField(field, data);                                                            
                            break;
                        // -----------------------------------------
                        // Bit-strings
                        // -----------------------------------------
                        case BYTE:  // 1 byte
                            fieldItem = decodeReadResponseSignedByteField(field, data);
                            break;
                        case WORD:  // 2 byte (16 bit)
                            fieldItem = decodeReadResponseSignedShortField(field, data);
                            break;
                        case DWORD:  // 4 byte (32 bit)
                            fieldItem = decodeReadResponseSignedIntegerField(field, data);
                            break;
                        case LWORD:  // 8 byte (64 bit)
                            fieldItem = decodeReadResponseSignedLongField(field, data);
                            break;
                        // -----------------------------------------
                        // Integers
                        // -----------------------------------------
                        // 8 bit:
                        case SINT:
                            fieldItem = decodeReadResponseSignedByteField(field, data);
                            break;
                        case USINT:
                            fieldItem = decodeReadResponseUnsignedByteField(field, data);
                            break;
                        // 16 bit:
                        case INT:
                            fieldItem = decodeReadResponseSignedShortField(field, data);
                            break;
                        case UINT:
                            fieldItem = decodeReadResponseUnsignedShortField(field, data);
                            break;
                        case S5TIME:
                            fieldItem = decodeReadResponseS5Time(field, data);
                            break;
                        // 32 bit:
                        case TIME:
                            fieldItem = decodeReadResponseTime(field, data);
                            break;
                        case DINT:
                            fieldItem = decodeReadResponseSignedIntegerField(field, data);
                            break;
                        case UDINT:
                            fieldItem = decodeReadResponseUnsignedIntegerField(field, data);
                            break;
                        // 64 bit:
                        case LINT:
                            fieldItem = decodeReadResponseSignedLongField(field, data);
                            break;
                        case ULINT:
                            fieldItem = decodeReadResponseUnsignedLongField(field, data);
                            break;
                        // -----------------------------------------
                        // Floating point values
                        // -----------------------------------------
                        case REAL:
                            fieldItem = decodeReadResponseFloatField(field, data);
                            break;
                        case LREAL:
                            fieldItem = decodeReadResponseDoubleField(field, data);
                            break;
                        // -----------------------------------------
                        // Characters & Strings
                        // -----------------------------------------
                        case CHAR: // 1 byte (8 bit)
                            fieldItem = decodeReadResponseFixedLengthStringField(1, false, data);
                            break;
                        case WCHAR: // 2 byte
                            fieldItem = decodeReadResponseFixedLengthStringField(1, true, data);
                            break;
                        case STRING:
                            fieldItem = decodeReadResponseVarLengthStringField(false, data);
                            break;
                        case WSTRING:
                            fieldItem = decodeReadResponseVarLengthStringField(true, data);
                            break;
                        // -----------------------------------------
                        // TIA Date-Formats
                        // -----------------------------------------
                        case DT:
                        case DATE_AND_TIME:
                            fieldItem = decodeReadResponseDateAndTime(field, data);
                            break;
                        case TOD:
                        case TIME_OF_DAY:
                            fieldItem = decodeReadResponseTimeOfDay(field, data);                            
                            break;
                        case DATE:
                            fieldItem = decodeReadResponseDate(field, data);
                            break;
                        default:
                            throw new PlcProtocolException("Unsupported type " + field.getDataType());
                    }
                }
                catch (Plc4XNettyException e){
                    logger.warn("Problem during casting of field {}: Exception: {}; FieldInformation: {}",fieldName,e.getMessage(),field);
                }
                catch (Exception e){
                    logger.warn("Some other error occurred casting field {}, FieldInformation: {}",fieldName, field,e);
                }
            }
            Pair<PlcResponseCode, BaseDefaultFieldItem> result = new ImmutablePair<>(responseCode, fieldItem);
            values.put(fieldName, result);
            index++;
            itemRow++;
        }

        return new DefaultPlcReadResponse(plcReadRequest, values);
    }
    
    private PlcResponse decodeSslReadResponse(S7ResponseMessage responseMessage, PlcRequestContainer requestContainer) throws PlcProtocolException {
        InternalPlcReadRequest plcReadRequest = (InternalPlcReadRequest) requestContainer.getRequest();
        
        Map<String, Pair<PlcResponseCode, ByteBuf>> values = new HashMap<>();                
        List<S7Payload> payloads = responseMessage.getPayloads();
        
        if (payloads.get(0) instanceof CpuServicesPayload){
            CpuServicesPayload payload = (CpuServicesPayload) payloads.get(0); 
            if (payload.getSslPayload() == null) {
                logger.debug("decodeSslReadResponse: Check for fragment payload");
                return null;
            }
        }
        
        payloads.forEach((payload)->{
            CpuServicesPayload s7payload = (CpuServicesPayload) payload;});    
             
        int index = 0;
        for (String requestName : plcReadRequest.getFieldNames()) {

            CpuServicesPayload payload = (CpuServicesPayload) payloads.get(index);  

            Pair<PlcResponseCode, ByteBuf> result = new ImmutablePair<>( 
                    decodeResponseCode(payload.getReturnCode()), 
                    payload.getSslPayload());
            values.put(requestName, result);

            index++;
        }
     
        return new S7ByteReadResponse(plcReadRequest, values);
    }    

    BaseDefaultFieldItem decodeReadResponseBitField(S7Field field, ByteBuf data) {
        Boolean[] booleans = readAllValues(Boolean.class, field, i -> data.readByte() != 0x00);
        return new DefaultBooleanFieldItem(booleans);
    }

    BaseDefaultFieldItem decodeReadResponseByteBitStringField(S7Field field, ByteBuf data) {
        byte[] bytes = new byte[field.getNumElements()];
        data.readBytes(bytes);
        return decodeBitStringField(bytes);
    }

    BaseDefaultFieldItem decodeReadResponseShortBitStringField(S7Field field, ByteBuf data) {
        byte[] bytes = new byte[field.getNumElements() * 2];
        data.readBytes(bytes);
        return decodeBitStringField(bytes);
    }

    BaseDefaultFieldItem decodeReadResponseIntegerBitStringField(S7Field field, ByteBuf data) {
        byte[] bytes = new byte[field.getNumElements() * 4];
        data.readBytes(bytes);
        return decodeBitStringField(bytes);
    }

    BaseDefaultFieldItem decodeReadResponseLongBitStringField(S7Field field, ByteBuf data) {
        byte[] bytes = new byte[field.getNumElements() * 8];
        data.readBytes(bytes);
        return decodeBitStringField(bytes);
    }

    BaseDefaultFieldItem decodeBitStringField(byte[] bytes) {
        BitSet bitSet = BitSet.valueOf(bytes);
        Boolean[] booleanValues = new Boolean[8 * bytes.length];
        int k = 0;
        for(int i = bytes.length - 1; i >= 0; i--) {
            for(int j = 0; j < 8; j++) {
                booleanValues[k++] = bitSet.get(8 * i + j);
            }
        }
        return new DefaultBooleanFieldItem(booleanValues);
    }

    BaseDefaultFieldItem decodeReadResponseSignedByteField(S7Field field, ByteBuf data) {
        Byte[] bytes = readAllValues(Byte.class, field, i -> data.readByte());
        return new DefaultByteFieldItem(bytes);
    }

    BaseDefaultFieldItem decodeReadResponseUnsignedByteField(S7Field field, ByteBuf data) {
        Short[] shorts = readAllValues(Short.class, field, i -> data.readUnsignedByte());
        return new DefaultShortFieldItem(shorts);
    }

    BaseDefaultFieldItem decodeReadResponseSignedShortField(S7Field field, ByteBuf data) {
        Short[] shorts = readAllValues(Short.class, field, i -> data.readShort());
        return new DefaultShortFieldItem(shorts);
    }

    BaseDefaultFieldItem decodeReadResponseUnsignedShortField(S7Field field, ByteBuf data) {
        Integer[] ints = readAllValues(Integer.class, field, i -> data.readUnsignedShort());
        return new DefaultIntegerFieldItem(ints);
    }

    BaseDefaultFieldItem decodeReadResponseSignedIntegerField(S7Field field, ByteBuf data) {
        Integer[] ints = readAllValues(Integer.class, field, i -> data.readInt());
        return new DefaultIntegerFieldItem(ints);
    }

    BaseDefaultFieldItem decodeReadResponseUnsignedIntegerField(S7Field field, ByteBuf data) {
        Long[] longs = readAllValues(Long.class, field, i -> data.readUnsignedInt());
        return new DefaultLongFieldItem(longs);
    }

    BaseDefaultFieldItem decodeReadResponseSignedLongField(S7Field field, ByteBuf data) {
        Long[] longs = readAllValues(Long.class, field, i -> data.readLong());
        return new DefaultLongFieldItem(longs);
    }

    BaseDefaultFieldItem decodeReadResponseUnsignedLongField(S7Field field, ByteBuf data) {
        BigInteger[] bigIntegers = readAllValues(BigInteger.class, field, i -> readUnsigned64BitInteger(data));
        return new DefaultBigIntegerFieldItem(bigIntegers);
    }

    BaseDefaultFieldItem decodeReadResponseFloatField(S7Field field, ByteBuf data) {
        Float[] floats = readAllValues(Float.class, field, i -> data.readFloat());
        return new DefaultFloatFieldItem(floats);
    }

    BaseDefaultFieldItem decodeReadResponseDoubleField(S7Field field, ByteBuf data) {
        Double[] doubles = readAllValues(Double.class, field, i -> data.readDouble());
        return new DefaultDoubleFieldItem(doubles);
    }

    BaseDefaultFieldItem decodeReadResponseFixedLengthStringField(int numChars, boolean isUtf16, ByteBuf data) {
        int numBytes = isUtf16 ? numChars * 2 : numChars;
        String stringValue = data.readCharSequence(numBytes, StandardCharsets.UTF_8).toString();
        return new DefaultStringFieldItem(stringValue);
    }

    BaseDefaultFieldItem decodeReadResponseVarLengthStringField(boolean isUtf16, ByteBuf data) {
        // Max length ... ignored.
        data.skipBytes(1);

        //reading out byte and transforming that to an unsigned byte within an integer, otherwise longer strings are failing
        byte currentLengthByte = data.readByte();
        int currentLength = currentLengthByte & 0xFF;
        return decodeReadResponseFixedLengthStringField(currentLength, isUtf16, data);
    }

    BaseDefaultFieldItem decodeReadResponseS5Time(S7Field field,ByteBuf data) {
        Duration[] durations = readAllValues(Duration.class,field, i -> readS5TimeToDuration(data));
        return new DefaultDurationFieldItem(durations);
    }
    
    BaseDefaultFieldItem decodeReadResponseTime(S7Field field,ByteBuf data) {
        Duration[] durations = readAllValues(Duration.class,field, i -> readS7TimeToDuration(data));
        return new DefaultDurationFieldItem(durations);
    }
    
    BaseDefaultFieldItem decodeReadResponseTimeOfDay(S7Field field,ByteBuf data) {
        LocalTime[] localTimes = readAllValues(LocalTime.class,field, i -> readTimeOfDay(data));
        return new DefaultLocalTimeFieldItem(localTimes);
    }

    BaseDefaultFieldItem decodeReadResponseDate(S7Field field,ByteBuf data) {
        LocalDate[] localDates = readAllValues(LocalDate.class,field, i -> readDate(data));
        return new DefaultLocalDateFieldItem(localDates);
    }
    
    BaseDefaultFieldItem decodeReadResponseDateAndTime(S7Field field,ByteBuf data) {
        LocalDateTime[] localDateTimes = readAllValues(LocalDateTime.class,field, i -> readDateAndTime(data));
        return new DefaultLocalDateTimeFieldItem(localDateTimes);
    }    

    // Returns a 32 bit unsigned value : from 0 to 4294967295 (2^32-1)
    public static int getUDIntAt(byte[] buffer, int pos) {
        int result;
        result = buffer[pos] & 0x0FF;
        result <<= 8;
        result += buffer[pos + 1] & 0x0FF;
        result <<= 8;
        result += buffer[pos + 2] & 0x0FF;
        result <<= 8;
        result += buffer[pos + 3] & 0x0FF;
        return result;
    }

    private static <T> T[] readAllValues(Class<T> clazz, S7Field field, Function<Integer, T> extract) {
        try {
            return IntStream.rangeClosed(1, field.getNumElements())
                .mapToObj(extract::apply)
                .collect(Collectors.toList())
                .toArray((T[])Array.newInstance(clazz, 0));
        } catch (IndexOutOfBoundsException e) {
            throw new PlcRuntimeException("To few bytes in the buffer to read requested type", e);
        }
    }

    @SuppressWarnings("unchecked")
    private PlcResponse decodeWriteResponse(S7ResponseMessage responseMessage, PlcRequestContainer requestContainer) throws PlcProtocolException {
        InternalPlcWriteRequest plcWriteRequest = (InternalPlcWriteRequest) requestContainer.getRequest();
        Map<Short, Integer> splitItems  = null;
        
        
        //TODO: if the PLC Send a error, It send only the Header with PDUR and error code.
        //      in this condition we send a error to the client, no a exception.
        VarPayload payload = responseMessage.getPayload(VarPayload.class)
            .orElseThrow(() -> new PlcProtocolException("No VarPayload supplied"));

        // TODO: Checkout if the payload items are sort of a flatMap of all request items.
        // TODO: Check for array response, no payload items presents, Why?

        // If the numbers of items don't match, we're in big trouble as the only
        // way to know how to interpret the responses is by aligning them with the
        // items from the request as this information is not returned by the PLC.
        // If we write BitsArray -> We need to add the number of items created 
        // for individual requests.
        int[] numberOfFields = new int[1];
        numberOfFields[0] = 0;
        short tpdu = responseMessage.getTpduReference();
        if (requestBitArray.containsKey(tpdu)) {
            splitItems = requestBitArray.remove(tpdu);
            splitItems.forEach((key,length)->{                
                numberOfFields[0] += length;
            }); 
            numberOfFields[0] += plcWriteRequest.getNumberOfFields();
        } else {
            numberOfFields[0] = plcWriteRequest.getNumberOfFields();
        }

        if (numberOfFields[0] != payload.getItems().size()) {
            throw new PlcProtocolException(
                "The number of requested items doesn't match the number of returned items Request: " +
                       numberOfFields[0] +
                        " " +
                       payload.getItems().size());
        }
        
        Map<String, PlcResponseCode> values = new HashMap<>();
        List<VarPayloadItem> payloadItems = payload.getItems();

        short sFieldCount = 0;
        short sPayloadCount = 0;
        
        for (String fieldName : plcWriteRequest.getFieldNames()) {
            
            VarPayloadItem payloadItem = payloadItems.get(sPayloadCount);
            PlcResponseCode responseCode = decodeResponseCode(payloadItem.getReturnCode());            
            
            if (splitItems != null) {     
             
                if (splitItems.containsKey(sFieldCount)) {
  
                    int numberOfItems = splitItems.get(sFieldCount);                                        
                    responseCode = decodeResponseCode(payloadItem.getReturnCode());
                    for (short i=sPayloadCount; i<=(short)(sPayloadCount+numberOfItems); i++){
                        if (payloadItem.getReturnCode() != DataTransportErrorCode.OK) {
                            responseCode = decodeResponseCode(payloadItem.getReturnCode());
                        };                        
                        payloadItem = payloadItems.get(i);
                    }
                    sPayloadCount += (short) numberOfItems+1; 
                } else {
                    responseCode = decodeResponseCode(payloadItem.getReturnCode());
                    sPayloadCount++;                   
                }
                
            } else {
                sPayloadCount++; 
            }

            values.put(fieldName, responseCode);

            sFieldCount++; 
        }
        

        return new DefaultPlcWriteResponse(plcWriteRequest, values);
    }


    private PlcResponse decodeSubscriptionResponse(S7ResponseMessage responseMessage, PlcRequestContainer requestContainer) throws PlcProtocolException {

        InternalPlcSubscriptionRequest subsRequest = (InternalPlcSubscriptionRequest) requestContainer.getRequest(); 
        //TODO: Try multiple ALARM_ACK in the same request and multiple EVENT_ID.
        S7SubscriptionField event = (S7SubscriptionField) subsRequest.getFields().get(0);
        
        //TODO: Multiple parameters. Try with S7400. At this point if we have embedded more parameter will be lost.
        CpuServicesResponseParameter parameter = (CpuServicesResponseParameter) responseMessage.getParameters().get(0);

        List<S7Payload> payloads = responseMessage.getPayloads();
        
        Map<String, Pair<PlcResponseCode, PlcSubscriptionHandle>> values = new HashMap<>();

        //May be only one iteration?:
        for (S7Payload payload:payloads){

            if (payload instanceof AlarmMessagePayload) {
                AlarmMessagePayload payloadItem = (AlarmMessagePayload) payload;   
                switch(parameter.getSubFunctionGroup()){
                    case MESSAGE_SERVICE:{
                        List<SubscribedEventType> items = new ArrayList();
                        PlcResponseCode responseCode = decodeResponseCode(((AlarmMessagePayload) payload).getReturnCode());

                        if (responseCode == PlcResponseCode.OK) {
                            subsRequest.getFields().forEach((field)->{
                                S7SubscriptionField s7field = (S7SubscriptionField) field;
                                items.add(s7field.getEventtype());
                            });
                        };    
                        
                        S7DiagnosticSubscriptionHandle handler = new S7DiagnosticSubscriptionHandle(
                                                                        "MS",
                                                                        responseMessage.getTpduReference(),
                                                                        responseCode,
                                                                        items);
                            
                        LinkedHashSet<String> fieldnames = subsRequest.getFieldNames();  
                        
                        fieldnames.forEach((fieldname)->{
                            values.put(fieldname, new ImmutablePair(responseCode, handler)); 
                        });
                            
                        }    
                    break;
                    case ALARM_ACK:{
                        if (payloadItem.getLength() != 0) {
                            List<Object> items = payloadItem.getMsg().getMsgItems();                            
                            ArrayList<Integer> eventids = event.getAckalarms();
                            int index = 0;
                            for (Object item:items){
                                // A ack response contains only the return code for every item.
                                // Is direct response for the ALARM_ACK, don't pass to
                                // push message handler.
                                //subsBuilder2.addEventField("MyAck", "ACK:16#60000001,16#60000002,16#60000003");
                                values.put("16#"+Integer.toHexString(eventids.get(index)), new ImmutablePair(decodeResponseCode((DataTransportErrorCode) item),null));                    
                                index++;
                            }    
                        } else {
                            values.put("ERROR", new ImmutablePair(PlcResponseCode.NOT_FOUND,null));                                                
                        }
                    }
                    break;
                    case ALARM_QUERY:{
                        if (payloadItem.getMsg() == null) {
                            return null;   
                        } 

                        if (payloadItem instanceof S7PushMessage){
                            if (payloadItem.getMsg().getMsgItems().size() > 0){
                                if (!alarmsqueue.offer((S7PushMessage) payload)){
                                    logger.info("decodeSubscriptionResponse: Alarm queue buffer is full.");
                                };
                            };
                        } else {
                            logger.debug("Check for ALARM_QUERY not attended: " + payloadItem);   
                        }                    
                    }
                    break;
                    default:;
                } 
            } else if (payload instanceof CpuCyclicServicesResponsePayload) {
                              
                    CpuCyclicServicesResponsePayload  cyclicPayloadItem = (CpuCyclicServicesResponsePayload) payload;
                    //
                    //
                    if (cyclicPayloadItem.getItems().size() != 0) {
                    
                        List<AssociatedValueItem> valueitems = cyclicPayloadItem.getItems();
                        Map<String,AssociatedValueItem> items = new LinkedHashMap();
                        int i=0;
                        if (parameter.getError().getCode() == 0x0000){                    
                            for (String fieldname:subsRequest.getFieldNames()){
                                items.put(fieldname, valueitems.get(i));
                                i++;                        
                            };
                        };

                        S7CyclicServicesSubscriptionHandle handler = new S7CyclicServicesSubscriptionHandle("UNO", 
                                                                            parameter.getSequenceNumber(),
                                                                            parameter.getError().getCode(),
                                                                            items);
                        // Return quality code for alarm 
                        for (String fieldname:subsRequest.getFieldNames()){
                            values.put(fieldname, new ImmutablePair(decodeResponseCode(items.get(fieldname).getReturnCode()), handler));                    
                        }

                    
                    } else {
                        logger.info("Must return values...?");
                    } 
                    
                } else {
                    logger.info("Payload not defined...");
            } 
        }
        return new DefaultPlcSubscriptionResponse(subsRequest, values);
    }
    
    
    private PlcResponse decodeUnSubscriptionResponse(S7ResponseMessage responseMessage, PlcRequestContainer requestContainer) throws PlcProtocolException {
        DefaultPlcUnsubscriptionRequest unSubsRequest = (DefaultPlcUnsubscriptionRequest) requestContainer.getRequest();
        //TODO: Multiple parameters. Try with S7400. At this point if we have embedded more parameter will be lost.
        CpuServicesResponseParameter parameter = (CpuServicesResponseParameter) responseMessage.getParameters().get(0);
        InternalPlcUnsubscriptionRequest request = (InternalPlcUnsubscriptionRequest) requestContainer.getRequest();

        return new DefaultPlcUnsubscriptionResponse(request);
    }    
    
    private PlcResponseCode decodeResponseCode(DataTransportErrorCode dataTransportErrorCode) {
        if (dataTransportErrorCode == null) {
            return PlcResponseCode.INTERNAL_ERROR;
        }
        switch (dataTransportErrorCode) {
            case OK:
                return PlcResponseCode.OK;
            case NOT_FOUND:
                return PlcResponseCode.NOT_FOUND;
            case INVALID_ADDRESS:
                return PlcResponseCode.INVALID_ADDRESS;
            case DATA_TYPE_NOT_SUPPORTED:
                return PlcResponseCode.INVALID_DATATYPE;
            case ACCESS_DENIED:
                return PlcResponseCode.ACCESS_DENIED;
            default:
                return PlcResponseCode.INTERNAL_ERROR;
        }
    }

    private static BigInteger readUnsignedLong(ByteBuf data) {
        // as there is no unsigned long primitive, we have to switch to
        // BigDecimal and manually convert the bytes to a BigDecimal.
        // In order to be unsigned 4 bytes, we create an array of 5 bytes
        // where the 5th byte is set to 0. The most significant bit being
        // 0 we are guaranteed to interpret the input a positive value.
        byte[] bytes = new byte[5];
        // Set the first byte to 0
        bytes[0] = 0;
        // Read the next 4 bytes into the rest.
        data.readBytes(bytes, 1, 4);
        return new BigInteger(bytes);
    }

    private static BigInteger readSigned64BitInteger(ByteBuf data) {
        byte[] bytes = new byte[8];
        data.readBytes(bytes, 0, 8);
        return new BigInteger(bytes);
    }

    private static BigInteger readUnsigned64BitInteger(ByteBuf data) {
        byte[] bytes = new byte[9];
        // Set the first byte to 0
        bytes[0] = 0;
        // Read the next 8 bytes into the rest.
        data.readBytes(bytes, 1, 8);
        return new BigInteger(bytes);
    }

    //TODO: millisoconds
    LocalDateTime readDateAndTime(ByteBuf data) {
        return S7Helper.S7DateTimeToLocalDateTime(data);
    }

    Duration readS5TimeToDuration(ByteBuf data){
        return S7Helper.S5TimeToDuration(data.readShort());        
    };
    
    Duration readS7TimeToDuration(ByteBuf data){
        return S7Helper.S7TimeToDuration(data.readInt());        
    };
    
    LocalTime readTimeOfDay(ByteBuf data) {
        return S7Helper.S7TodToLocalTime(data.readInt());
    }

    LocalDate readDate(ByteBuf data) {
        return S7Helper.S7DateToLocalDate(data.readShort());
    }

    /**
     * converts incoming byte to an integer regarding used BCD format
     * @param incomingByte
     * @return converted BCD number
     */
    private static int convertByteToBcd(byte incomingByte) {
        int dec = (incomingByte >> 4) * 10;
        return dec + (incomingByte & 0x0f);
    }
    


}
