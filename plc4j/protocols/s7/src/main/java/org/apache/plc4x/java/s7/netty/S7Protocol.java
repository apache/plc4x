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
import io.netty.buffer.ByteBufUtil;
import io.netty.buffer.Unpooled;
import io.netty.channel.*;
import io.netty.handler.codec.MessageToMessageDecoder;
import io.netty.util.concurrent.Future;
import io.netty.util.concurrent.PromiseCombiner;
import org.apache.commons.lang3.reflect.FieldUtils;
import org.apache.plc4x.java.api.exceptions.PlcProtocolPayloadTooBigException;
import org.apache.plc4x.java.isotp.netty.IsoTPProtocol;
import org.apache.plc4x.java.isotp.netty.events.IsoTPConnectedEvent;
import org.apache.plc4x.java.isotp.netty.model.IsoTPMessage;
import org.apache.plc4x.java.isotp.netty.model.tpdus.DataTpdu;
import org.apache.plc4x.java.s7.netty.events.S7ConnectedEvent;
import org.apache.plc4x.java.s7.netty.model.messages.*;
import org.apache.plc4x.java.s7.netty.model.params.VarParameter;
import org.apache.plc4x.java.s7.netty.model.params.S7Parameter;
import org.apache.plc4x.java.s7.netty.model.params.SetupCommunicationParameter;
import org.apache.plc4x.java.s7.netty.model.params.items.VarParameterItem;
import org.apache.plc4x.java.s7.netty.model.params.items.S7AnyVarParameterItem;
import org.apache.plc4x.java.s7.netty.model.payloads.S7Payload;
import org.apache.plc4x.java.s7.netty.model.payloads.VarPayload;
import org.apache.plc4x.java.s7.netty.model.payloads.items.VarPayloadItem;
import org.apache.plc4x.java.s7.netty.model.types.*;
import org.apache.plc4x.java.s7.netty.strategies.DefaultS7MessageProcessor;
import org.apache.plc4x.java.s7.netty.strategies.S7MessageProcessor;
import org.apache.plc4x.java.s7.netty.util.S7SizeHelper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.Field;
import java.util.*;

/**
 * Communication Layer between the Application level ({@link Plc4XS7Protocol}) and the lower level (tcp) that sends and receives {@link S7Message}s.
 * This layer also handles the control over the "wire", i.e., the queues of incoming and outgoing messages.
 * Furthermore, here {@link S7Message}s are marshalled and unmarshalled to {@link ByteBuf}s to be send over wire.
 *
 * Before messages are send to the wire an optional {@link S7MessageProcessor} can be applied.
 *
 * @see S7MessageProcessor
 */
public class S7Protocol extends ChannelDuplexHandler {

    private static final byte S7_PROTOCOL_MAGIC_NUMBER = 0x32;

    private static final Logger logger = LoggerFactory.getLogger(S7Protocol.class);

    private final MessageToMessageDecoder<Object> decoder = new MessageToMessageDecoder<Object>() {

        @Override
        public boolean acceptInboundMessage(Object msg) {
            return msg instanceof IsoTPMessage;
        }

        @Override
        @SuppressWarnings("unchecked")
        protected void decode(ChannelHandlerContext ctx, Object msg, List<Object> out) {
            S7Protocol.this.decode(ctx, (IsoTPMessage) msg, out);
        }
    };

    private short maxAmqCaller;
    private short maxAmqCallee;
    private short pduSize;

    // For detecting the lower layers.
    private ChannelHandler prevChannelHandler;
    private S7MessageProcessor messageProcessor;

    // For being able to respect the max AMQ restrictions.
    private PendingWriteQueue queue;
    private Map<Short, DataTpdu> sentButUnacknowledgedTpdus;

    public S7Protocol(short requestedMaxAmqCaller, short requestedMaxAmqCallee, short requestedPduSize) {
        this.maxAmqCaller = requestedMaxAmqCaller;
        this.maxAmqCallee = requestedMaxAmqCallee;
        this.pduSize = requestedPduSize;
        sentButUnacknowledgedTpdus = new HashMap<>();
        messageProcessor = new DefaultS7MessageProcessor();
    }

    @Override
    public void channelRegistered(ChannelHandlerContext ctx) {
        this.queue = new PendingWriteQueue(ctx);
        try {
            Field prevField = FieldUtils.getField(ctx.getClass(), "prev", true);
            if(prevField != null) {
                ChannelHandlerContext prevContext = (ChannelHandlerContext) prevField.get(ctx);
                prevChannelHandler = prevContext.handler();
            }
        } catch(Exception e) {
            logger.error("Error accessing field 'prev'", e);
        }
    }

    @Override
    public void channelUnregistered(ChannelHandlerContext ctx) throws Exception {
        this.queue.removeAndWriteAll();
        super.channelUnregistered(ctx);
    }

    @Override
    public void channelInactive(ChannelHandlerContext ctx) throws Exception {
        // Send everything so we get a proper failure for those pending writes
        this.queue.removeAndWriteAll();
        super.channelInactive(ctx);
    }

    /**
     * If the S7 protocol layer is used over Iso TP, then after receiving a {@link IsoTPConnectedEvent} the
     * corresponding S7 setup communication message has to be sent in order to negotiate the S7 protocol layer.
     *
     * @param ctx the current protocol layers context
     * @param evt the event
     * @throws Exception throws an exception if something goes wrong internally
     */
    @Override
    public void userEventTriggered(ChannelHandlerContext ctx, Object evt) throws Exception {
        // If we are using S7 inside of IsoTP, then we need to intercept IsoTPs connected events.
        if ((prevChannelHandler instanceof IsoTPProtocol) && (evt instanceof IsoTPConnectedEvent)) {
            // Setup Communication
            SetupCommunicationRequestMessage setupCommunicationRequest =
                new SetupCommunicationRequestMessage((short) 0, maxAmqCaller, maxAmqCallee, pduSize, null);

            ctx.channel().writeAndFlush(setupCommunicationRequest);
        }

        else {
            super.userEventTriggered(ctx, evt);
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Encoding
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    public void write(ChannelHandlerContext ctx, Object msg, ChannelPromise promise) throws Exception {
        if(msg instanceof S7Message) {
            S7Message in = (S7Message) msg;

            // Give message processors to process the incoming message.
            Collection<? extends S7Message> messages;
            if((messageProcessor != null) && (in instanceof S7RequestMessage)) {
                try {
                    messages = messageProcessor.processRequest((S7RequestMessage) in, pduSize);
                } catch(Exception e) {
                    logger.error("Error processing message", e);
                    ctx.fireExceptionCaught(e);
                    return;
                }
            } else {
                messages = Collections.singleton(in);
            }

            // Create a promise that has to be called multiple times.
            PromiseCombiner promiseCombiner = new PromiseCombiner();
            if(messages != null) {
                for (S7Message message : messages) {
                    ByteBuf buf = Unpooled.buffer();

                    encodeHeader(message, buf);
                    encodeParameters(message, buf);
                    encodePayloads(message, buf);

                    // Check if the message doesn't exceed the negotiated maximum size.
                    if (buf.writerIndex() > pduSize) {
                        ctx.fireExceptionCaught(new PlcProtocolPayloadTooBigException("s7", pduSize, buf.writerIndex(), message));
                    } else {
                        ChannelPromise subPromise = new DefaultChannelPromise(promise.channel());
                        queue.add(new DataTpdu(true, (byte) 0x01, Collections.emptyList(), buf, message), subPromise);
                        promiseCombiner.add((Future) subPromise);
                        logger.debug("S7 Message with id {} queued", message.getTpduReference());
                    }
                }
            }
            promiseCombiner.finish(promise);

            // Start sending the queue content.
            trySendingMessages(ctx);
        } else {
            super.write(ctx, msg, promise);
        }
    }

    private void encodePayloads(S7Message in, ByteBuf buf) {
        for (S7Payload payload : in.getPayloads()) {
            ParameterType parameterType = payload.getType();

            // When sending requests currently only write var has payloads.
            if (parameterType == ParameterType.WRITE_VAR) {
                VarPayload varPayload = (VarPayload) payload;
                for (VarPayloadItem payloadItem : varPayload.getItems()) {
                    buf.writeByte(payloadItem.getReturnCode().getCode());
                    buf.writeByte(payloadItem.getDataTransportSize().getCode());
                    // TODO: Check if this is correct?!?! Might be problems with sizeInBits = true/false
                    buf.writeShort(payloadItem.getData().length);
                    buf.writeBytes(payloadItem.getData());
                    // TODO: It looks as if BIT type reads require a 0x00 fill byte at the end ...
                }
            }
        }
    }

    private void encodeParameters(S7Message in, ByteBuf buf) {
        for (S7Parameter s7Parameter : in.getParameters()) {
            buf.writeByte(s7Parameter.getType().getCode());
            switch (s7Parameter.getType()) {
                case READ_VAR:
                case WRITE_VAR:
                    encodeParameterReadWriteVar(buf, (VarParameter) s7Parameter);
                    break;
                case SETUP_COMMUNICATION:
                    encodeParameterSetupCommunication(buf, (SetupCommunicationParameter) s7Parameter);
                    break;
                default:
                    logger.error("writing this parameter type not implemented");
            }
        }
    }

    private void encodeHeader(S7Message in, ByteBuf buf) {
        buf.writeByte(S7_PROTOCOL_MAGIC_NUMBER);
        buf.writeByte(in.getMessageType().getCode());
        // Reserved (is always constant 0x0000)
        buf.writeShort((short) 0x0000);
        // PDU Reference (Request Id, generated by the initiating node)
        buf.writeShort(in.getTpduReference());
        // S7 message parameters length
        buf.writeShort(S7SizeHelper.getParametersLength(in.getParameters()));
        // Data field length
        buf.writeShort(S7SizeHelper.getPayloadsLength(in.getPayloads()));
        if (in instanceof S7ResponseMessage) {
            S7ResponseMessage s7ResponseMessage = (S7ResponseMessage) in;
            buf.writeByte(s7ResponseMessage.getErrorClass());
            buf.writeByte(s7ResponseMessage.getErrorCode());
        }
    }

    private void encodeParameterSetupCommunication(ByteBuf buf, SetupCommunicationParameter s7Parameter) {
        // Reserved (is always constant 0x00)
        buf.writeByte((byte) 0x00);
        buf.writeShort(s7Parameter.getMaxAmqCaller());
        buf.writeShort(s7Parameter.getMaxAmqCallee());
        buf.writeShort(s7Parameter.getPduLength());
    }

    private void encodeParameterReadWriteVar(ByteBuf buf, VarParameter s7Parameter) {
        List<VarParameterItem> items = s7Parameter.getItems();
        // ReadRequestItem count (Read one variable at a time)
        buf.writeByte((byte) items.size());
        for (VarParameterItem item : items) {
            VariableAddressingMode addressMode = item.getAddressingMode();
            if (addressMode == VariableAddressingMode.S7ANY) {
                encodeS7AnyParameterItem(buf, (S7AnyVarParameterItem) item);
            } else {
                logger.error("writing this item type not implemented");
            }
        }
    }

    private void encodeS7AnyParameterItem(ByteBuf buf, S7AnyVarParameterItem s7AnyRequestItem) {
        buf.writeByte(s7AnyRequestItem.getSpecificationType().getCode());
        // Length of this item (excluding spec type and length)
        buf.writeByte((byte) 0x0a);
        buf.writeByte(s7AnyRequestItem.getAddressingMode().getCode());
        buf.writeByte(s7AnyRequestItem.getTransportSize().getCode());
        buf.writeShort(s7AnyRequestItem.getNumElements());
        buf.writeShort(s7AnyRequestItem.getDataBlockNumber());
        buf.writeByte(s7AnyRequestItem.getMemoryArea().getCode());
        // A S7 address is 3 bytes long. Unfortunately the byte-offset is NOT located in
        // byte 1 and byte 2 and the bit offset in byte 3. Siemens used the last 3 bits of
        // byte 3 for the bit-offset and the remaining 5 bits of byte 3 to contain the lowest
        // 5 bits of the byte-offset. The highest 5 bits of byte 1 are probably left unused
        // for future extensions.
        buf.writeShort((short) (s7AnyRequestItem.getByteOffset() >> 5));
        buf.writeByte((byte) ((
                (s7AnyRequestItem.getByteOffset() & 0x1F) << 3)
                | (s7AnyRequestItem.getBitOffset() & 0x07)));
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Decoding
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    public void channelRead(ChannelHandlerContext ctx, Object msg) throws Exception {
        decoder.channelRead(ctx, msg);
        super.channelRead(ctx, msg);
    }

    protected void decode(ChannelHandlerContext ctx, IsoTPMessage in, List<Object> out) {
        if (logger.isTraceEnabled()) {
            logger.trace("Got Data: {}", ByteBufUtil.hexDump(in.getUserData()));
        }
        ByteBuf userData = in.getUserData();
        if (userData.readableBytes() == 0) {
            return;
        }

        if (userData.readByte() != S7_PROTOCOL_MAGIC_NUMBER) {
            logger.warn("Expecting S7 protocol magic number.");
            if (logger.isDebugEnabled()) {
                logger.debug("Got Data: {}", ByteBufUtil.hexDump(userData));
            }
            return;
        }

        MessageType messageType = MessageType.valueOf(userData.readByte());
        boolean isResponse = messageType == MessageType.ACK_DATA;
        userData.readShort();  // Reserved (is always constant 0x0000)
        short tpduReference = userData.readShort();
        short headerParametersLength = userData.readShort();
        short userDataLength = userData.readShort();
        byte errorClass = 0;
        byte errorCode = 0;
        if (isResponse) {
            errorClass = userData.readByte();
            errorCode = userData.readByte();
        }

        List<S7Parameter> s7Parameters = new LinkedList<>();
        VarParameter readWriteVarParameter = null;
        int i = 0;

        while (i < headerParametersLength) {
            S7Parameter parameter = decodeParameter(userData, isResponse, headerParametersLength - i);
            s7Parameters.add(parameter);
            if (parameter instanceof SetupCommunicationParameter) {
                handleSetupCommunications(ctx, (SetupCommunicationParameter) parameter);
            }
            if (parameter instanceof VarParameter) {
                readWriteVarParameter = (VarParameter) parameter;
            }
            i += S7SizeHelper.getParameterLength(parameter);
        }

        List<S7Payload> s7Payloads = decodePayloads(userData, isResponse, userDataLength, readWriteVarParameter);

        logger.debug("S7 Message with id {} received", tpduReference);

        if (isResponse) {
            S7ResponseMessage responseMessage = new S7ResponseMessage(
                messageType, tpduReference, s7Parameters, s7Payloads, errorClass, errorCode);

            // Remove the current response from the list of unconfirmed messages.
            DataTpdu requestTpdu = sentButUnacknowledgedTpdus.remove(tpduReference);

            // Get the corresponding request message.
            S7RequestMessage requestMessage = (requestTpdu != null) ? (S7RequestMessage) requestTpdu.getParent() : null;

            if(requestMessage != null) {
                // Set this individual request to "acknowledged".
                requestMessage.setAcknowledged(true);

                // Give the request and response to a message processor to process the incoming message.
                if(messageProcessor != null) {
                    try {
                        responseMessage = messageProcessor.processResponse(requestMessage, responseMessage);
                    } catch(Exception e) {
                        logger.error("Error processing message", e);
                        ctx.fireExceptionCaught(e);
                        return;
                    }
                }

                if(responseMessage != null) {
                    out.add(responseMessage);
                }

                // Eventually send the next message (if there is one).
                trySendingMessages(ctx);
            }
        } else {
            // TODO: Find out if there is any situation in which a request is sent from the PLC
            out.add(new S7RequestMessage(messageType, tpduReference, s7Parameters, s7Payloads, null));
        }
    }

    private void handleSetupCommunications(ChannelHandlerContext ctx, SetupCommunicationParameter setupCommunicationParameter) {
        maxAmqCaller = setupCommunicationParameter.getMaxAmqCaller();
        maxAmqCallee = setupCommunicationParameter.getMaxAmqCallee();
        pduSize = setupCommunicationParameter.getPduLength();

        logger.info("S7Connection established pdu-size {}, max-amq-caller {}, " +
                "max-amq-callee {}", pduSize, maxAmqCaller, maxAmqCallee);

        // Send an event that setup is complete.
        ctx.channel().pipeline().fireUserEventTriggered(new S7ConnectedEvent());
    }

    private List<S7Payload> decodePayloads(ByteBuf userData, boolean isResponse, short userDataLength, VarParameter readWriteVarParameter) {
        int i = 0;
        List<S7Payload> s7Payloads = new LinkedList<>();
        if (readWriteVarParameter != null) {
            List<VarPayloadItem> payloadItems = new LinkedList<>();

            // Just keep on reading payloads until the provided length is read.
            while (i < userDataLength) {
                DataTransportErrorCode dataTransportErrorCode = DataTransportErrorCode.valueOf(userData.readByte());
                // This is a response to a WRITE_VAR request (It only contains the return code for every sent item.
                if ((readWriteVarParameter.getType() == ParameterType.WRITE_VAR) && isResponse) {
                    // Initialize a rudimentary payload (This is updated in the Plc4XS7Protocol class
                    VarPayloadItem payload = new VarPayloadItem(dataTransportErrorCode, null, null);
                    payloadItems.add(payload);
                    i += 1;
                }
                // This is a response to a READ_VAR request.
                else if ((readWriteVarParameter.getType() == ParameterType.READ_VAR) && isResponse) {
                    DataTransportSize dataTransportSize = DataTransportSize.valueOf(userData.readByte());
                    short length = (dataTransportSize.isSizeInBits()) ?
                        (short) Math.ceil(userData.readShort() / 8.0) : userData.readShort();
                    byte[] data = new byte[length];
                    userData.readBytes(data);
                    // Initialize a rudimentary payload (This is updated in the Plc4XS7Protocol class
                    VarPayloadItem payload = new VarPayloadItem(dataTransportErrorCode, dataTransportSize, data);
                    payloadItems.add(payload);
                    i += S7SizeHelper.getPayloadLength(payload);
                }
            }

            VarPayload varPayload = new VarPayload(readWriteVarParameter.getType(), payloadItems);
            s7Payloads.add(varPayload);
        }
        return s7Payloads;
    }

    private S7Parameter decodeParameter(ByteBuf in, boolean isResponse, int restLength) {
        ParameterType parameterType = ParameterType.valueOf(in.readByte());
        if (parameterType == null) {
            logger.error("Could not find parameter type");
            return null;
        }
        switch (parameterType) {
            case CPU_SERVICES:
                // Just read in the rest of the header as content of this parameter.
                // Will have to do a lot more investigation on how this parameter is
                // constructed.
                byte[] cpuServices = new byte[restLength - 1];
                in.readBytes(cpuServices);
                return null;
            case READ_VAR:
            case WRITE_VAR:
                List<VarParameterItem> varParamameter;
                byte numItems = in.readByte();
                if (!isResponse) {
                    varParamameter = decodeReadWriteVarParameter(in, numItems);
                } else {
                    varParamameter = Collections.emptyList();
                }
                return new VarParameter(parameterType, varParamameter);
            case SETUP_COMMUNICATION:
                // Reserved (is always constant 0x00)
                in.readByte();
                short callingMaxAmq = in.readShort();
                short calledMaxAmq = in.readShort();
                short pduLength = in.readShort();
                return new SetupCommunicationParameter(callingMaxAmq, calledMaxAmq, pduLength);
            default:
                if (logger.isErrorEnabled()) {
                    logger.error("Unimplemented parameter type: {}", parameterType.name());
                }
        }
        return null;
    }

    private List<VarParameterItem> decodeReadWriteVarParameter(ByteBuf in, byte numItems) {
        List<VarParameterItem> items = new LinkedList<>();
        for (int i = 0; i < numItems; i++) {
            SpecificationType specificationType = SpecificationType.valueOf(in.readByte());
            // Length of the rest of this item.
            byte itemLength = in.readByte();
            if (itemLength != 0x0a) {
                logger.warn("Expecting a length of 10 here.");
                return items;
            }
            VariableAddressingMode variableAddressingMode = VariableAddressingMode.valueOf(in.readByte());
            if (variableAddressingMode == VariableAddressingMode.S7ANY) {
                TransportSize transportSize = TransportSize.valueOf(in.readByte());
                short length = in.readShort();
                short dbNumber = in.readShort();
                MemoryArea memoryArea = MemoryArea.valueOf(in.readByte());
                short byteAddress = (short) (in.readShort() << 5);
                byte tmp = in.readByte();
                // Only the least 3 bits are the bit address, the
                byte bitAddress = (byte) (tmp & 0x07);
                // Bits 4-8 belong to the byte address
                byteAddress = (short) (byteAddress | (tmp >> 3));
                S7AnyVarParameterItem item = new S7AnyVarParameterItem(
                        specificationType, memoryArea, transportSize,
                        length, dbNumber, byteAddress, bitAddress);
                items.add(item);
            } else {
                logger.error("Error parsing item type");
                return items;
            }
        }

        return items;
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Helpers
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////


    @Override
    public void flush(ChannelHandlerContext ctx) throws Exception {
        super.flush(ctx);
    }

    private synchronized void trySendingMessages(ChannelHandlerContext ctx) {
        while(sentButUnacknowledgedTpdus.size() < maxAmqCaller) {
            // Get the TPDU that is up next in the queue.
            DataTpdu curTpdu = (DataTpdu) queue.current();

            if (curTpdu != null) {
                // Send the TPDU.
                try {
                    ChannelFuture channelFuture = queue.removeAndWrite();
                    ctx.flush();
                    if (channelFuture == null) {
                        break;
                    }
                } catch (Exception e) {
                    logger.error("Error sending more queues messages", e);
                    ctx.fireExceptionCaught(e);
                }

                if(curTpdu.getParent() != null) {
                    // Add it to the list of sentButUnacknowledgedTpdus.
                    // (It seems that the S7 drops the value of the COTP reference id, so we have to use the S7 one)
                    S7RequestMessage s7RequestMessage = (S7RequestMessage) curTpdu.getParent();
                    sentButUnacknowledgedTpdus.put(s7RequestMessage.getTpduReference(), curTpdu);

                    logger.debug("S7 Message with id {} sent", s7RequestMessage.getTpduReference());
                }
                // TODO: Eventually remove this.
                break;
            } else {
                break;
            }
        }
        ctx.flush();
    }

}
