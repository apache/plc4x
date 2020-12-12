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
package org.apache.plc4x.java.s7.netty.strategies;

import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.MessageToMessageCodec;
import io.netty.util.AttributeKey;
import java.util.Collections;
import java.util.Iterator;
import java.util.LinkedHashMap;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;
import java.util.Map.Entry;
import java.util.Optional;
import java.util.Set;
import java.util.concurrent.atomic.AtomicInteger;
import org.apache.plc4x.java.base.messages.PlcProtocolMessage;
import org.apache.plc4x.java.s7.netty.model.messages.S7RequestMessage;
import org.apache.plc4x.java.s7.netty.model.messages.S7ResponseMessage;
import org.apache.plc4x.java.s7.netty.model.params.S7Parameter;
import org.apache.plc4x.java.s7.netty.model.params.VarParameter;
import org.apache.plc4x.java.s7.netty.model.params.items.S7AnyVarParameterItem;
import org.apache.plc4x.java.s7.netty.model.params.items.VarParameterItem;
import org.apache.plc4x.java.s7.netty.model.payloads.S7Payload;
import org.apache.plc4x.java.s7.netty.model.payloads.VarPayload;
import org.apache.plc4x.java.s7.netty.model.payloads.items.VarPayloadItem;
import org.apache.plc4x.java.s7.netty.model.types.DataTransportErrorCode;
import org.apache.plc4x.java.s7.netty.model.types.MessageType;
import org.apache.plc4x.java.s7.netty.model.types.ParameterType;
import org.apache.plc4x.java.s7.netty.model.types.TransportSize;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 *
 * @author cgarcia
 */
public class LongS7MessageProcessorCodec extends
             MessageToMessageCodec<S7ResponseMessage, S7RequestMessage>{
    
    private static final Logger LOGGER = LoggerFactory.getLogger(LongS7MessageProcessorCodec.class );
    private final AttributeKey<Short> pduSizeKey = AttributeKey.valueOf("pduSizeKey");
    private final AttributeKey<AtomicInteger> tpduGeneratorKey = AttributeKey.valueOf("tpduGeneratorKey");
    
    private final Map<Short, S7CompositeRequestMessage> compRequestMessages;
    
    private Short pduSize = 0;
    private AtomicInteger tpduGenerator;

    public LongS7MessageProcessorCodec() {
        this.compRequestMessages = new LinkedHashMap<>(10);
    }
             
    @Override
    protected void encode(ChannelHandlerContext ctx, S7RequestMessage request, List<Object> out) throws Exception {
        pduSize = ctx.channel().attr(pduSizeKey).get();
        tpduGenerator = ctx.channel().attr(tpduGeneratorKey).get();
        
        Optional<VarParameter> varParameterOptional = request.getParameter(VarParameter.class);
        if (varParameterOptional.isPresent()) {
            VarParameter varParameter = varParameterOptional.get();
            /* Requests whose size is larger than the psuSize should be 
            *processed or if it is smaller than this but contains a number 
            *of elements greater than 18 */
            if ((pduSize < 1024) && (tpduGenerator != null) && (
                    (request.getRequestSize() > pduSize) ||
                    ((request.getRequestSize() < pduSize) && (varParameter.getItems().size() > 18)))){
                if (varParameter.getType() == ParameterType.READ_VAR){
                    LinkedHashMap<Short, S7RequestMessage> splitRequest = SplitLongReadRequestMessageGroup(request);
                    out.addAll(splitRequest.values());
                    return;                    
                } else if(varParameter.getType() == ParameterType.WRITE_VAR){
                    LinkedHashMap<Short, S7RequestMessage> splitRequest = SplitLongWriteRequestMessageGroup(request);
                    out.addAll(splitRequest.values());
                    return;
                }
            }       
        }
        out.add(request);
    }

    @Override
    protected void decode(ChannelHandlerContext ctx, S7ResponseMessage response, List<Object> out) throws Exception {
        S7CompositeRequestMessage compRequestMessage = compRequestMessages.get(response.getTpduReference());
        if (compRequestMessage != null) {
            compRequestMessage.addResponseMessage(response);
            if (compRequestMessage.isCompleted()){
                S7ResponseMessage res = AssemblyResponseMessage(compRequestMessage, response);
                out.add(res);
                return;
            }
        }
        out.add(response);
    }
    
    /************************************************************************
    * SPLIT SECTION
    ************************************************************************/
    
    /*
    * This method split original S7RequestMessage using 80% rule.
    * Another problem when tpduGenerator returns to zero (short type).
    * TODO: split for individual item when item.size > pduSize
    */
    private LinkedHashMap<Short, S7RequestMessage> SplitLongWriteRequestMessageGroup(S7RequestMessage request){
        
        boolean firtsRow = true;
        S7CompositeRequestMessage compRequestMessage = new S7CompositeRequestMessage(request);

        compRequestMessages.put(request.getTpduReference(), compRequestMessage);

        List<S7RequestMessage> splitRequests = new LinkedList<>();

        Optional<VarParameter> varParameterOptional = request.getParameter(VarParameter.class);        
        VarParameter parameters =  varParameterOptional.get();
                
        Optional<VarPayload> varPayloadOptional = request.getPayload(VarPayload.class); 
        VarPayload payloads = varPayloadOptional.get();
               
        short itemSize = 0;
        short actualSize = 0;
        short maxItemSize = (short) Math.round(pduSize * 0.8);
        
        S7Parameter parameter = null;
        S7Payload payload = null;
       
        List<VarParameterItem> parameterItems = (LinkedList)parameters.getItems();
        List<VarPayloadItem> payloadItems  = (LinkedList)payloads.getItems();
        List<VarParameterItem> splitParameters = new LinkedList<>();
        List<VarPayloadItem> splitPayloads = new LinkedList<>();
                       
        S7AnyVarParameterItem s7AnyVarParameterItem = null;
        
        //No se estan generando en orden
        Iterator<VarParameterItem> itParameter = parameterItems.iterator();
        Iterator<VarPayloadItem> itPayload = payloadItems.iterator();
        int i = 0;
        while (itParameter.hasNext()) {
                                    
            VarParameterItem varParameterItem = itParameter.next();
            VarPayloadItem varPayloadItem = itPayload.next();            
                        
            if (varParameterItem instanceof S7AnyVarParameterItem) {
                s7AnyVarParameterItem = (S7AnyVarParameterItem) varParameterItem;
            }

            //TODO: for odd size count the fill 0x00 byte
            itemSize = (short)(28 + varPayloadItem.getData().length); //Header+Parameter+data
            actualSize = (short)(28 + 14*(splitPayloads.size()));//Header+Parameter
            for(VarPayloadItem item:splitPayloads){
                actualSize = (short) (actualSize +item.getData().length); //+=data.length
            };
            
            LOGGER.debug("SplitLongWriteRequestMessageGroup itemSize: " + itemSize + " actualSize: " + actualSize);
            
            //TODO: Item size is > pduSize the Split is in the user layer. 
            if (itemSize > pduSize){
                //TODO: Evaluate that the fractional message should be blocked 
                //      by the next layer. The response message indicating poor 
                //      quality must also be generated at this point.
                if (!splitParameters.isEmpty()){
                    PrepareMessagesToNextLayer(request.getMessageType(),
                    parameters.getType(),
                    splitParameters,
                    splitPayloads,
                    compRequestMessage); 
                    
                    splitParameters = new LinkedList<>();
                    splitPayloads = new LinkedList<>();                    
                }
                
                short tpdu = PrepareMessagesToNextLayer(request.getMessageType(),
                                parameters.getType(),
                                Collections.singletonList(varParameterItem),
                                Collections.singletonList(varPayloadItem),
                                compRequestMessage);                 
                
                VarParameter writeVarParameter = new VarParameter(parameters.getType(), 
                        Collections.singletonList(varParameterItem));
                VarPayload writeVarPayload = new VarPayload(parameters.getType(), 
                        Collections.singletonList(varPayloadItem));
                //This message generates an "error" from the PLC, which is discarded.
                S7RequestMessage splitRequest = new S7RequestMessage(request.getMessageType(),
                    tpdu, Collections.singletonList(writeVarParameter),
                    Collections.EMPTY_LIST, compRequestMessage);
                                              
                VarParameter ackVarParameter = new VarParameter(parameters.getType(), 
                    Collections.singletonList(
                        new S7AnyVarParameterItem(null, null, null, 0x01, (short) 0, (short) 0, (byte) 0)));
                VarPayload ackVarPayload = new VarPayload(parameters.getType(), 
                        Collections.singletonList( new VarPayloadItem(DataTransportErrorCode.ACCESS_DENIED, null, null)));                
                
                S7ResponseMessage responseMessage = new S7ResponseMessage(
                    MessageType.ACK_DATA, 
                    splitRequest.getTpduReference(), 
                    Collections.singletonList(ackVarParameter), 
                    Collections.singletonList(ackVarPayload), 
                    (byte) 0x81, //Application relationship 
                    (byte) 0x04); //Service no implemented or frame error
               /* 
                compRequestMessage.getResponseMessages().put(
                        responseMessage.getTpduReference(), 
                        responseMessage);
                */
               
                compRequestMessage.getRequestMessages().put(splitRequest.getTpduReference(),splitRequest);  
                compRequestMessages.put(splitRequest.getTpduReference(), compRequestMessage);
                //A fault response is generated, with poor quality, allowing 
                //to follow the writing process for the other Items.
                
                //throw new Plc4XNettyException("Write item request size bigger than " + pduSize+
                //            " bytes. Split the Item in smaller units."); 
                
            } else if ((itemSize > maxItemSize) && ((itemSize <  pduSize))){
                //TODO: Add a S7RequestMessage to the list                               
                PrepareMessagesToNextLayer(request.getMessageType(),
                        parameters.getType(),
                        splitParameters,
                        splitPayloads,
                        compRequestMessage); 
                
            } else if (!firtsRow){
               
                if ((((itemSize + actualSize) >= maxItemSize) && 
                     ((itemSize + actualSize) <= pduSize)) || 
                        (!itParameter.hasNext())){  
                    splitParameters.add(varParameterItem);
                    splitPayloads.add(varPayloadItem);  
                    
                    PrepareMessagesToNextLayer(request.getMessageType(),
                            parameters.getType(),
                            splitParameters,
                            splitPayloads,
                            compRequestMessage); 
                    
                    splitParameters = new LinkedList<>();
                    splitPayloads = new LinkedList<>();
                } else if (((itemSize + actualSize) >= pduSize)){  
                    PrepareMessagesToNextLayer(request.getMessageType(),
                            parameters.getType(),
                            splitParameters,
                            splitPayloads,
                            compRequestMessage); 
                    splitParameters = new LinkedList<>();
                    splitPayloads = new LinkedList<>(); 
                    splitParameters.add(varParameterItem);
                    splitPayloads.add(varPayloadItem);                      
                }
                else {
                    splitParameters.add(varParameterItem);
                    splitPayloads.add(varPayloadItem);                    
                }                
            } else {
                splitParameters.add(varParameterItem);
                splitPayloads.add(varPayloadItem);
            }
            firtsRow = false;
        }
                        
        if (!splitParameters.isEmpty()){
            PrepareMessagesToNextLayer(request.getMessageType(),
                    parameters.getType(),
                    splitParameters,
                    splitPayloads,
                    compRequestMessage);        
        }
   
        return compRequestMessage.getRequestMessages();
    };

    private LinkedHashMap<Short, S7RequestMessage> SplitLongReadRequestMessageGroup(S7RequestMessage request){

        S7CompositeRequestMessage compRequestMessage = new S7CompositeRequestMessage(request);      
        compRequestMessages.put(request.getTpduReference(), compRequestMessage);
        
        Optional<VarParameter> varParameterOptional = request.getParameter(VarParameter.class);        
        VarParameter parameters =  varParameterOptional.get();                
        
        boolean firtsRow = true;
        short itemSize = 0;
        short actualSize = 0;
        short maxItemSize = (short) Math.round(pduSize * 0.8);        
        
        
        List<VarParameterItem> parameterItems = parameters.getItems();

        List<VarParameterItem> splitParameters = new LinkedList<>();
        List<VarPayloadItem> splitPayloads = new LinkedList<>();

        for (int i = 0; i < parameterItems.size(); i++) {
            
            VarParameterItem varParameterItem = parameterItems.get(i);
            
            itemSize = (short)(2 + getEstimatedResponseReadVarPayloadItemSize(varParameterItem));
            actualSize = 0;
            for(VarParameterItem item:splitParameters){
                actualSize = (short) (actualSize + getEstimatedResponseReadVarPayloadItemSize(item)); //+=data.length
            };     
            
            //TODO: Item size is > pduSize the Split is in the user layer. 
            if (itemSize > pduSize){
                //TODO: Evaluate that the fractional message should be blocked 
                //      by the next layer. The response message indicating poor 
                //      quality must also be generated at this point.
                if (!splitParameters.isEmpty()){
                    PrepareMessagesToNextLayer(request.getMessageType(),
                    parameters.getType(),
                    splitParameters,
                    Collections.EMPTY_LIST,
                    compRequestMessage); 
                    
                    splitParameters = new LinkedList<>();
                    splitPayloads = new LinkedList<>();                    
                }                
                
                short tpdu = PrepareMessagesToNextLayer(request.getMessageType(),
                                parameters.getType(),
                                Collections.singletonList(varParameterItem),
                                Collections.EMPTY_LIST,
                                compRequestMessage);                 
                
                //throw new Plc4XNettyException("Read item request size bigger than " + pduSize+
                //            " bytes. Split the Item in smaller units."); 
                
            } else if ((itemSize > maxItemSize) && ((itemSize <  pduSize))){
                
                if (!splitParameters.isEmpty()){
                    PrepareMessagesToNextLayer(request.getMessageType(),
                    parameters.getType(),
                    splitParameters,
                    Collections.EMPTY_LIST,
                    compRequestMessage); 
                    
                    splitParameters = new LinkedList<>();
                    splitPayloads = new LinkedList<>();                    
                }                                     
                
                VarParameter writeVarParameter = new VarParameter(parameters.getType(), 
                        Collections.singletonList(varParameterItem ));
                
                S7RequestMessage splitRequest = new S7RequestMessage(request.getMessageType(),
                    (short) tpduGenerator.getAndIncrement(), Collections.singletonList(writeVarParameter),
                    Collections.EMPTY_LIST, compRequestMessage);
                
                compRequestMessage.getRequestMessages().put(splitRequest.getTpduReference(),splitRequest);  
                compRequestMessages.put(splitRequest.getTpduReference(), compRequestMessage); 
                itemSize = 0;
                
            }  else if (!firtsRow){
                              
                if ((((itemSize + actualSize) >= maxItemSize) && 
                     ((itemSize + actualSize) <= pduSize)) ||
                       (splitParameters.size() > 16) ){  
                      
                    splitParameters.add(varParameterItem);
                    
                    VarParameter writeVarParameter = new VarParameter(parameters.getType(), splitParameters);

                    S7RequestMessage splitRequest = new S7RequestMessage(request.getMessageType(),
                        (short) tpduGenerator.getAndIncrement(), Collections.singletonList(writeVarParameter),
                        Collections.EMPTY_LIST, compRequestMessage);

                    compRequestMessage.getRequestMessages().put(splitRequest.getTpduReference(),splitRequest); 
                    compRequestMessages.put(splitRequest.getTpduReference(), compRequestMessage);
                    
                    splitParameters = new LinkedList<>();                                   
                    splitPayloads = new LinkedList<>();
                    itemSize = 0;
                    
                } else if ((((itemSize + actualSize) >= maxItemSize) && 
                     ((itemSize + actualSize) > pduSize))){ 
                    
                    if (!splitParameters.isEmpty()){
                        PrepareMessagesToNextLayer(request.getMessageType(),
                        parameters.getType(),
                        splitParameters,
                        Collections.EMPTY_LIST,
                        compRequestMessage); 

                        splitParameters = new LinkedList<>();
                        splitPayloads = new LinkedList<>();                    
                    }                                  
                    
                    VarParameter writeVarParameter = new VarParameter(parameters.getType(), 
                            Collections.singletonList(varParameterItem ));

                    S7RequestMessage splitRequest = new S7RequestMessage(request.getMessageType(),
                        (short) tpduGenerator.getAndIncrement(), Collections.singletonList(writeVarParameter),
                        Collections.EMPTY_LIST, compRequestMessage);

                    compRequestMessage.getRequestMessages().put(splitRequest.getTpduReference(),splitRequest);  
                    compRequestMessages.put(splitRequest.getTpduReference(), compRequestMessage); 
                    itemSize = 0;
                    
                }  else {
                    splitParameters.add(varParameterItem);                   
                }                 
                
            } else {
                splitParameters.add(varParameterItem);
            }
            
            if((i == (parameterItems.size() - 1)) && (splitParameters.size() !=0)){
                VarParameter writeVarParameter = new VarParameter(parameters.getType(), splitParameters);

                S7RequestMessage splitRequest = new S7RequestMessage(request.getMessageType(),
                    (short) tpduGenerator.getAndIncrement(), Collections.singletonList(writeVarParameter),
                    Collections.EMPTY_LIST, compRequestMessage);

                compRequestMessage.getRequestMessages().put(splitRequest.getTpduReference(),splitRequest); 
                compRequestMessages.put(splitRequest.getTpduReference(), compRequestMessage);

                splitParameters = new LinkedList<>();                                   
                splitPayloads = new LinkedList<>();
                itemSize = 0;                   
            }
            
            firtsRow = false;          
     
        }        
                
        return compRequestMessage.getRequestMessages();       
    }    
    
    
    /************************************************************************
    * ASSEMBLY RESPONSE MESSAGE
    ************************************************************************/
    
    /*
    * Try to assemble a response message from the response fragments.
    * Since the PLC can handle at least two requests and depending on 
    * the load of its CPUs, the response order is not granted.
    * TODO: If the connection is lost, the queue is invalid and must be deleted.
    * TODO: The PLC reporta error, only the header is send. 
    *       Example Error Class 0x81, Error code 0x04
    */
    private S7ResponseMessage AssemblyResponseMessage(S7CompositeRequestMessage compRequestMessage, 
            S7ResponseMessage response){
        
        List<S7Parameter> splitParameters = new LinkedList<>();
        List<S7Payload> splitPayloads = new LinkedList<>();
        List<VarPayloadItem> payloadItems = new LinkedList<>();
        
        ParameterType parameterType = ParameterType.WRITE_VAR;
        int numParameters = 0x000;
        short tpduReference = 0x0000;
                
        LinkedHashMap<Short, S7RequestMessage> requestMessages =
                compRequestMessage.getRequestMessages();
        
        LinkedHashMap<Short, S7ResponseMessage> responseMessages =
                compRequestMessage.getResponseMessages();
        
        tpduReference =((S7RequestMessage) compRequestMessage.getParent()).getTpduReference();
        //Always the firts element is the response the pdu number    
        Iterator<Entry<Short, S7RequestMessage>> iterRequest = requestMessages.entrySet().iterator();
        
        while(iterRequest.hasNext()){

            Entry<Short, S7RequestMessage> itemRequest = iterRequest.next();
            S7ResponseMessage msgResponse = (S7ResponseMessage) responseMessages.get(itemRequest.getKey());
                        
            //If no present we have a error code in the Header
            Optional<VarParameter> varParameterOptional = msgResponse.getParameter(VarParameter.class); 
            if (varParameterOptional.isPresent()){
                VarParameter parameters =  varParameterOptional.get();
                parameterType = parameters.getType();
                //Como extraer los parametros del mensaje de respuesta            
                List<VarParameterItem> varParameterItems = parameters.getItems();

                S7AnyVarParameterItem s7AnyVarParameterItem = (S7AnyVarParameterItem) varParameterItems.get(0);

                numParameters += s7AnyVarParameterItem.getNumElements();

                List<S7Payload> payloads = msgResponse.getPayloads();
                
                splitPayloads.addAll(payloads);                
            } else {
                //Error from PLC
                LOGGER.info("Error for tpud " + msgResponse.getTpduReference() +
                            " Error class: " + msgResponse.getErrorClass() +
                            " Error code" + msgResponse.getErrorCode());
                S7RequestMessage msgRequest = (S7RequestMessage) requestMessages.get(msgResponse.getTpduReference());

                VarParameter requestParameters = (VarParameter) msgRequest.getParameters().get(0);
                parameterType = requestParameters.getType();
                List<VarParameterItem> varParameterItems = requestParameters.getItems();
                S7AnyVarParameterItem s7AnyVarParameterItem = (S7AnyVarParameterItem) varParameterItems.get(0);  
                numParameters += s7AnyVarParameterItem.getNumElements();  
                VarPayload ackVarPayload = new VarPayload(requestParameters.getType(), 
                    Collections.singletonList( new VarPayloadItem(DataTransportErrorCode.ACCESS_DENIED, null, null)));                     
                splitPayloads.addAll(Collections.singletonList(ackVarPayload)); 
                                
            }

        }

        //Clean request messages
        Set<Short> setKeys = requestMessages.keySet();
        setKeys.forEach(key->{compRequestMessages.remove(key);});
        compRequestMessages.remove(tpduReference);
                
        List<VarParameterItem> varParameterItems;
        varParameterItems = Collections.singletonList(
                        new S7AnyVarParameterItem(null, null, null, numParameters , (short) 0, (short) 0, (byte) 0));

        S7Parameter s7Parameter = new VarParameter(parameterType, varParameterItems);
        
        splitParameters.add(s7Parameter); 
        
        for(S7Payload s7Payload:splitPayloads){
            VarPayload varPayload = (VarPayload) s7Payload;
            payloadItems.addAll(varPayload.getItems());            
        }
                       
        S7ResponseMessage res = new S7ResponseMessage(response.getMessageType(),
                tpduReference,
                splitParameters,
                Collections.singletonList(new VarPayload(parameterType, payloadItems)), 
                (byte)0x00, (byte)0x00);

        return res;
    }
    
    /************************************************************************
    * UTILITYS
    ************************************************************************/    
    

    static class S7CompositeRequestMessage implements PlcProtocolMessage {

        private S7RequestMessage originalRequest;
        private LinkedHashMap<Short, S7RequestMessage> requestMessages;
        private LinkedHashMap<Short, S7ResponseMessage> responseMessages;

        S7CompositeRequestMessage(S7RequestMessage originalRequest) {
            this.originalRequest = originalRequest;
            this.requestMessages = new LinkedHashMap<>();
            this.responseMessages = new LinkedHashMap<>();
        }

        @Override
        public PlcProtocolMessage getParent() {
            return originalRequest;
        }

        /**
         * A {@link S7CompositeRequestMessage} is only acknowledged, if all children are acknowledged.
         *
         * @return true if all children are acknowledged.
         */
        private boolean isAcknowledged() {
            /*
            for (S7RequestMessage requestMessage : requestMessages) {
                if(!requestMessage.isAcknowledged()) {
                    return false;
                }
            }*/
            return true;
        }

        void addRequestMessage(S7RequestMessage requestMessage) {
            requestMessages.put(requestMessage.getTpduReference(),requestMessage);
        }

        public LinkedHashMap getRequestMessages() {
            return requestMessages;
        }

        private void addResponseMessage(S7ResponseMessage responseMessage) {
            responseMessages.putIfAbsent(responseMessage.getTpduReference(),responseMessage);
        }

        public LinkedHashMap getResponseMessages() {
            return  responseMessages;
        }
        
        public boolean isCompleted(){            
            return ((requestMessages.size() - responseMessages.size()) == 0) ;
        }
    }
    
    //TODO: Use BaseType for length calculation. (Ok)
    //TODO: For S7 STRING, maximum length is 254 -> One STRING request -> S7RequestMessage. HMI used CHAR array's 
    private static short getEstimatedResponseReadVarPayloadItemSize(VarParameterItem varParameterItem) {
        // A var payload item always has a minimum size of 4 bytes (return code, transport size, size (two bytes))
        short length = 4;
        short base =0;
        S7AnyVarParameterItem s7AnyVarParameterItem = (S7AnyVarParameterItem) varParameterItem;
        
        base = (s7AnyVarParameterItem.getDataType().getBaseType() == null)?
                (short) s7AnyVarParameterItem.getDataType().getSizeInBytes():
                (short)(s7AnyVarParameterItem.getDataType().getSizeInBytes() *
                s7AnyVarParameterItem.getDataType().getBaseType().getSizeInBytes());
        
        length +=
            s7AnyVarParameterItem.getNumElements() * base;
        // It seems that bit payloads need a additional separating 0x00 byte.
        if(s7AnyVarParameterItem.getDataType() == TransportSize.BOOL) {
            length += 1;
        }
        return length;
    }


    private short PrepareMessagesToNextLayer(MessageType msgType,
            ParameterType parameter,
            List<VarParameterItem> splitParameters,
            List<VarPayloadItem> splitPayloads,
            S7CompositeRequestMessage compMsg){
            
            VarParameter writeVarParameter = new VarParameter(parameter, splitParameters);
            VarPayload writeVarPayload = new VarPayload(parameter, splitPayloads); 
            
            short tpdu = (short) tpduGenerator.getAndIncrement();                        
            
            LOGGER.debug("SendMessagesToNextLayer tpdu: " + tpdu + 
                    " Parameters:" + splitParameters.size() +
                    " Payloads: " + splitPayloads.size());
            
            S7RequestMessage splitRequest = new S7RequestMessage(msgType,
                tpdu,
                Collections.singletonList(writeVarParameter),
                (splitPayloads == Collections.EMPTY_LIST)?
                        Collections.EMPTY_LIST:
                        Collections.singletonList(writeVarPayload),
                compMsg);            

            compMsg.getRequestMessages().put(splitRequest.getTpduReference(),splitRequest);  
            compRequestMessages.put(splitRequest.getTpduReference(), compMsg);
            
            return tpdu;
    }
    
    
    
}
