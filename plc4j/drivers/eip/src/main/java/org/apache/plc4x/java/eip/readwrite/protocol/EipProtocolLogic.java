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
package org.apache.plc4x.java.eip.readwrite.protocol;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.eip.readwrite.*;
import org.apache.plc4x.java.eip.readwrite.configuration.EIPConfiguration;
import org.apache.plc4x.java.eip.readwrite.field.EipField;
import org.apache.plc4x.java.eip.readwrite.types.CIPDataTypeCode;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.messages.*;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.UnsupportedEncodingException;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.time.Duration;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.logging.Level;

public class EipProtocolLogic extends Plc4xProtocolBase<EipPacket>implements HasConfiguration<EIPConfiguration> {

    private static final Logger logger = LoggerFactory.getLogger(EipProtocolLogic.class);
    public static final Duration REQUEST_TIMEOUT = Duration.ofMillis(10000);

    private static final short[] emptySenderContext = new short[] {(short) 0x00 ,(short) 0x00 ,(short) 0x00,
        (short) 0x00,(short) 0x00,(short) 0x00, (short) 0x00,(short) 0x00};
    private short[] senderContext;
    private EIPConfiguration configuration;

    private final AtomicInteger transactionCounterGenerator = new AtomicInteger(10);
    private RequestTransactionManager tm;
    private long sessionHandle;

    @Override
    public void setConfiguration(EIPConfiguration configuration){
        this.configuration = configuration;
        // Set the transaction manager to allow only one message at a time.
        this.tm = new RequestTransactionManager(1);
    }

    @Override
    public void onConnect(ConversationContext<EipPacket> context) {
        /**Send a ENIP Message with Register Session Code '0x0065',
         * empty Session Handle and Sender Context
         * Then we need to accept the response with the same Code
         * and save the assigned Session Handle
         * PS: Check Status for Success : 0x00000000*/
        logger.info("Sending RegisterSession EIP Package");
        EipConnectionRequest connectionRequest =
          new EipConnectionRequest(0L,0L,emptySenderContext,0L);
        context.sendRequest(connectionRequest)
            .expectResponse(EipPacket.class, REQUEST_TIMEOUT).unwrap( p -> (EipPacket) p)
            .check(p -> p instanceof EipConnectionRequest)
            .handle(p -> {
                if(p.getStatus()==0L){
                    sessionHandle = p.getSessionHandle();
                    senderContext= p.getSenderContext();
                    logger.trace("Got assigned with Session {}", sessionHandle);
                    // Send an event that connection setup is complete.
                    context.fireConnected();
                }
                else{
                    logger.warn("Got status code [{}]", p.getStatus());
                }

            });
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        DefaultPlcReadRequest request = (DefaultPlcReadRequest) readRequest;
        List<CipReadRequest> requests = new ArrayList<>(request.getNumberOfFields());
        for(PlcField field : request.getFields()) {
            EipField plcField = (EipField) field;
            String tag = plcField.getTag();
            int elements = 1 ;
            if(plcField.getElementNb()>1){
                elements = plcField.getElementNb();
            }

            //We need the size of the request in words (0x91, tagLength, ... tag + possible pad)
            // Taking half to get word size
            boolean isArray = false;
            String tagIsolated=tag;
            if(tag.contains("[")){
                 isArray = true;
                 tagIsolated = tag.substring(0, tag.indexOf("["));
            }
            int dataLength = (tagIsolated.length() + 2 + (tagIsolated.length() % 2)+(isArray? 2:0));
            byte requestPathSize = (byte) (dataLength/ 2);
            CipReadRequest req = new CipReadRequest(requestPathSize, toAnsi(tag), elements);
            requests.add(req);
        }
        return toPlcReadResponse((InternalPlcReadRequest)readRequest, readInternal(requests));
    }

    private byte[] toAnsi(String tag){
        int arrayIndex = 0;
        boolean isArray = false;
        if(tag.contains("[")){
            isArray = true;
            String index = tag.substring(tag.indexOf("[")+1, tag.indexOf("]"));
            arrayIndex = Integer.parseInt(index);
            tag = tag.substring(0, tag.indexOf("["));
        }

        boolean isPadded = tag.length() % 2 != 0;
        int dataSegLength = 2 + tag.length() + (isPadded==true? 1 : 0) + (isArray==true? 2 : 0);

        ByteBuffer buffer = ByteBuffer.allocate(dataSegLength).order(ByteOrder.LITTLE_ENDIAN);

        buffer.put((byte)0x91);
        buffer.put((byte)tag.length());
        byte[] tagBytes = null;
        try {
            tagBytes = tag.getBytes("US-ASCII");
        } catch (UnsupportedEncodingException ex) {
            java.util.logging.Logger.getLogger(EipProtocolLogic.class.getName()).log(Level.SEVERE, null, ex);
        }

        buffer.put(tagBytes);
        buffer.position(2 +tagBytes.length);


        if(isPadded)
        {
            buffer.put((byte)0x00);
        }

        if(isArray){
            buffer.put((byte)0x28);
            buffer.put((byte)arrayIndex);
        }
        return buffer.array();
    }

    private CompletableFuture<PlcReadResponse> toPlcReadResponse(InternalPlcReadRequest readRequest, CompletableFuture<CipService> response) {
        return response
            .thenApply(p -> {
                try {
                        return ((PlcReadResponse) decodeReadResponse(p, readRequest));

                } catch (PlcProtocolException e) {
                    throw new PlcRuntimeException("Unable to decode Response", e);
                }
            });
    }

    private CompletableFuture<CipService> readInternal(List<CipReadRequest> request) {
        CompletableFuture<CipService> future = new CompletableFuture<>();
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        if(request.size()>1){

            short nb =(short) request.size();
            int[] offsets = new int[nb];
            int offset = 2 + nb*2;
            for(int i = 0; i < nb ; i++){
                offsets[i]=offset;
                offset+=request.get(i).getLengthInBytes();
            }

            CipService[] serviceArr = new CipService[nb];
            for(int i = 0; i < nb ; i++){
                serviceArr[i]=request.get(i);
            }
            Services data = new Services(nb,offsets,serviceArr);
            //Encapsulate the data

            CipRRData pkt = new CipRRData(sessionHandle,0L, emptySenderContext, 0L,
                    new CipExchange(
                        new CipUnconnectedRequest(
                            new MultipleServiceRequest(data),
                            (byte) configuration.getBackplane(),
                            (byte)configuration.getSlot())));


            transaction.submit(() -> context.sendRequest(pkt)
                .expectResponse(EipPacket.class, REQUEST_TIMEOUT)
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .check(p -> p instanceof CipRRData)
                .check(p -> p.getSessionHandle() == sessionHandle)
                //.check(p -> p.getSenderContext() == senderContext)
                .unwrap(p -> (CipRRData) p)
                .unwrap(p -> p.getExchange().getService()).check(p -> p instanceof MultipleServiceResponse)
                .unwrap(p -> (MultipleServiceResponse) p)
                .check(p -> p.getData().getServiceNb() == nb)
                .handle(p -> {
                    future.complete(p);
                    // Finish the request-transaction.
                    transaction.endRequest();
                }));
        }
        else if(request.size()==1) {
            CipExchange exchange = new CipExchange(new CipUnconnectedRequest(request.get(0),(byte)configuration.getBackplane(),(byte)configuration.getSlot()));
            CipRRData pkt = new CipRRData(sessionHandle,0L, emptySenderContext, 0L, exchange);
            transaction.submit(() -> context.sendRequest(pkt)
                .expectResponse(EipPacket.class, REQUEST_TIMEOUT)
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .check(p -> p instanceof CipRRData)
                .check(p -> p.getSessionHandle() == sessionHandle)
                //.check(p -> p.getSenderContext() == senderContext)
                .unwrap(p -> (CipRRData) p)
                .unwrap(p -> p.getExchange().getService()).check(p -> p instanceof CipReadResponse)
                .unwrap(p -> (CipReadResponse) p)
                .handle(p -> {
                    future.complete(p);
                    // Finish the request-transaction.
                    transaction.endRequest();
                }));
        }
        return future;
    }

    private PlcResponse decodeReadResponse(CipService p, InternalPlcReadRequest readRequest) throws PlcProtocolException {
        //TODO Check if this is right
        Map<String, Pair<PlcResponseCode, PlcValue>> values = new HashMap<>();
        // only 1 field
        if(p instanceof CipReadResponse) {
            CipReadResponse resp = (CipReadResponse)p;
            String fieldName = readRequest.getFieldNames().iterator().next();
            EipField field = (EipField) readRequest.getField(fieldName);
            PlcResponseCode code = decodeResponseCode(resp.getStatus());
            PlcValue plcValue = null;
            CIPDataTypeCode type = resp.getDataType();
            ByteBuf data = Unpooled.wrappedBuffer(resp.getData());
            if (code == PlcResponseCode.OK) {
                plcValue = parsePlcValue(field, data, type);
            }
            Pair<PlcResponseCode, PlcValue> result = new ImmutablePair<>(code, plcValue);
            values.put(fieldName, result);
        }
        //Multiple response
        else if(p instanceof MultipleServiceResponse && ((MultipleServiceResponse) p).getStatus()==0){
            MultipleServiceResponse responses = (MultipleServiceResponse)p;
            Services services = responses.getData();
            int nb = services.getServiceNb();
            Iterator<String> it = readRequest.getFieldNames().iterator();
            for(int i = 0; i<nb && it.hasNext() ; i++){
                String fieldName = it.next();
                EipField field = (EipField) readRequest.getField(fieldName);
                PlcValue plcValue = null;
                if(services.getServices()[i] instanceof CipReadResponse){
                    CipReadResponse readResponse = (CipReadResponse)services.getServices()[i];
                    PlcResponseCode code;
                    if (readResponse.getStatus() == 0) {
                        code = PlcResponseCode.OK;
                    } else {
                        code = PlcResponseCode.INTERNAL_ERROR;
                    }
                    CIPDataTypeCode type = readResponse.getDataType();
                    ByteBuf data = Unpooled.wrappedBuffer(readResponse.getData());
                    if (code == PlcResponseCode.OK) {
                        plcValue = parsePlcValue(field, data, type);
                    }
                    Pair<PlcResponseCode, PlcValue> result = new ImmutablePair<>(code, plcValue);
                    values.put(fieldName, result);
                }
            }
        }
        return new DefaultPlcReadResponse(readRequest,values);
    }

    private PlcValue parsePlcValue(EipField field, ByteBuf data, CIPDataTypeCode type) {
        int nb = field.getElementNb();
        if(nb>1){
            int index =0;
            List<PlcValue> list = new ArrayList<>();
            for(int i=0 ; i < nb ; i++){
                switch(type){
                    case DINT:
                    case INT:
                    case SINT:
                        list.add(new PlcInteger(Integer.reverseBytes(data.getInt(index))));
                        index+= type.getSize();
                        break;
                    case REAL:
                        list.add(new PlcDouble(swap(data.getFloat(index))));
                        index+= type.getSize();
                        break;
                    default:
                        return null;
                }
            }
            return new PlcList(list);
        }
        else {
            switch (type) {
                case SINT:
                    return new PlcByte(data.getByte(0));
                case INT:
                    return new PlcInteger(Short.reverseBytes(data.getShort(0)));
                case DINT:
                    return new PlcInteger(Integer.reverseBytes(data.getInt(0)));
                case REAL:
                    return new PlcDouble(swap(data.getFloat(0)));
                default:
                    return null;
            }
        }
    }

    public float swap(float value){
        int bytes = Float.floatToIntBits(value);
        int b1 = (bytes >>  0) & 0xff;
        int b2 = (bytes >>  8) & 0xff;
        int b3 = (bytes >> 16) & 0xff;
        int b4 = (bytes >> 24) & 0xff;
        return Float.intBitsToFloat(b1 << 24 | b2 << 16 | b3 << 8 | b4 << 0);
    }

    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        CompletableFuture<PlcWriteResponse> future = new CompletableFuture<>();
        DefaultPlcWriteRequest request = (DefaultPlcWriteRequest) writeRequest;
        List<CipWriteRequest> items = new ArrayList<>(writeRequest.getNumberOfFields());
        for(String fieldName : request.getFieldNames()){
            final EipField field = (EipField) request.getField(fieldName);
            final PlcValue value = request.getPlcValue(fieldName);
            String tag = field.getTag();
            int elements = 1 ;
            if(field.getElementNb()>1){
                elements = field.getElementNb();
            }

            //We need the size of the request in words (0x91, tagLength, ... tag + possible pad)
            // Taking half to get word size
            boolean isArray = false;
            String tagIsolated=tag;
            if(tag.contains("[")){
                isArray = true;
                tagIsolated = tag.substring(0, tag.indexOf("["));
            }
            int dataLength = (tagIsolated.length() + 2 + (tagIsolated.length() % 2)+(isArray? 2:0));
            byte requestPathSize = (byte) (dataLength/ 2);
            byte[]data = encodeValue(value,field.getType(),(short)elements);
            CipWriteRequest writeReq = new CipWriteRequest(requestPathSize,toAnsi(tag),field.getType(),elements,data);
            items.add(writeReq);
        }

        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        if(items.size()==1){
            tm.startRequest();
            CipRRData rrdata = new CipRRData(sessionHandle,0L,senderContext,0L,
                new CipExchange(items.get(0)));
            transaction.submit(()-> context.sendRequest(rrdata)
            .expectResponse(EipPacket.class, REQUEST_TIMEOUT)
            .onTimeout(future::completeExceptionally)
            .onError((p,e) -> future.completeExceptionally(e))
                .check( p -> p instanceof CipRRData).unwrap(p -> (CipRRData)p)
                .check(p -> p.getSessionHandle() == sessionHandle)
              //.check(p -> p.getSenderContext() == senderContext)
                .check(p -> p.getExchange().getService() instanceof CipWriteResponse)
                .unwrap(p -> (CipWriteResponse)p.getExchange().getService())
                .handle(p ->{
                    future.complete((PlcWriteResponse)decodeWriteResponse(p,((InternalPlcWriteRequest)writeRequest)));
                    transaction.endRequest();
                })
            );
        }
        else {
            tm.startRequest();
            short nb =(short) items.size();
            int[] offsets = new int[nb];
            int offset = 2 + nb*2;
            for(int i = 0; i < nb ; i++){
                offsets[i]=offset;
                offset+=items.get(i).getLengthInBytes();
            }

            CipService[] serviceArr = new CipService[nb];
            for(int i = 0; i < nb ; i++){
                serviceArr[i]=items.get(i);
            }
            Services data = new Services(nb,offsets,serviceArr);
            //Encapsulate the data

            CipRRData pkt = new CipRRData(sessionHandle,0L, emptySenderContext, 0L,
                new CipExchange(
                    new CipUnconnectedRequest(
                        new MultipleServiceRequest(data),
                        (byte) configuration.getBackplane(),
                        (byte)configuration.getSlot())));


            transaction.submit(() -> context.sendRequest(pkt)
                .expectResponse(EipPacket.class, REQUEST_TIMEOUT)
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .check(p -> p instanceof CipRRData)
                .check(p -> p.getSessionHandle() == sessionHandle)
                //.check(p -> p.getSenderContext() == senderContext)
                .unwrap(p -> (CipRRData) p)
                .unwrap(p -> p.getExchange().getService()).check(p -> p instanceof MultipleServiceResponse)
                .unwrap(p -> (MultipleServiceResponse) p)
                .check(p -> p.getData().getServiceNb() == nb)
                .handle(p -> {
                    future.complete((PlcWriteResponse)decodeWriteResponse(p,((InternalPlcWriteRequest)writeRequest)));
                    // Finish the request-transaction.
                    transaction.endRequest();
                }));
        }
        return future;
    }

    private PlcResponse decodeWriteResponse(CipService p, InternalPlcWriteRequest writeRequest) {
        Map<String, PlcResponseCode> responses = new HashMap<>();

        if(p instanceof CipWriteResponse){
            CipWriteResponse resp = (CipWriteResponse)p;
            String fieldName = writeRequest.getFieldNames().iterator().next();
            EipField field =(EipField) writeRequest.getField(fieldName);
            responses.put(fieldName, decodeResponseCode(resp.getStatus()));
            return new DefaultPlcWriteResponse(writeRequest,responses);
        }
        else if(p instanceof MultipleServiceResponse){
            MultipleServiceResponse multResponses = (MultipleServiceResponse)p;
            Services services = multResponses.getData();
            int nb = services.getServiceNb();
            Iterator<String> it = writeRequest.getFieldNames().iterator();
            for(int i = 0; i<nb && it.hasNext() ; i++){
                String fieldName = it.next();
                EipField field = (EipField) writeRequest.getField(fieldName);
                PlcValue plcValue = null;
                if(services.getServices()[i] instanceof CipWriteResponse){
                    CipWriteResponse writeResponse = (CipWriteResponse)services.getServices()[i];
                    PlcResponseCode code = decodeResponseCode(writeResponse.getStatus());
                    responses.put(fieldName, code);
                }
            }
            return new DefaultPlcWriteResponse(writeRequest, responses);
        }
        return null;
    }

    private byte[] encodeValue(PlcValue value, CIPDataTypeCode type, short elements) {
        ByteBuffer buffer = ByteBuffer.allocate(4+type.getSize()).order(ByteOrder.LITTLE_ENDIAN);
        switch(type){
            case SINT:
                buffer.put(value.getByte());
                break;
            case INT:
                buffer.putShort(value.getShort());
                break;
            case DINT:
                buffer.putInt(value.getInteger());
                break;
            case REAL:
                buffer.putDouble(value.getDouble());
                break;
            default:break;
        }
        return buffer.array();

    }

    private PlcResponseCode decodeResponseCode(int status){
        switch (status){
            case 0 : return PlcResponseCode.OK;
            default: return PlcResponseCode.INTERNAL_ERROR;
        }
    }

    @Override
    public void close(ConversationContext<EipPacket> context) {
        /**Send a ENIP Message with Unregister Session Code '0x0066' */
        logger.info("Sending UnregisterSession EIP Pakcet");
        context.sendRequest(new EipDisconnectRequest(sessionHandle, 0L, emptySenderContext, 0L)); //Unregister gets no response
        logger.trace("Unregistred Session {}", sessionHandle);
    }
}
