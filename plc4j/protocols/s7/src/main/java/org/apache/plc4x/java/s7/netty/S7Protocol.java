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
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcProtocolPayloadTooBigException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.isotp.protocol.IsoTPProtocol;
import org.apache.plc4x.java.isotp.protocol.events.IsoTPConnectedEvent;
import org.apache.plc4x.java.isotp.protocol.model.IsoTPMessage;
import org.apache.plc4x.java.isotp.protocol.model.tpdus.DataTpdu;
import org.apache.plc4x.java.s7.netty.events.S7ConnectedEvent;
import org.apache.plc4x.java.s7.netty.model.messages.S7Message;
import org.apache.plc4x.java.s7.netty.model.messages.S7RequestMessage;
import org.apache.plc4x.java.s7.netty.model.messages.S7ResponseMessage;
import org.apache.plc4x.java.s7.netty.model.messages.SetupCommunicationRequestMessage;
import org.apache.plc4x.java.s7.netty.model.params.*;
import org.apache.plc4x.java.s7.netty.model.params.items.S7AnyVarParameterItem;
import org.apache.plc4x.java.s7.netty.model.params.items.VarParameterItem;
import org.apache.plc4x.java.s7.netty.model.payloads.CpuServicesPayload;
import org.apache.plc4x.java.s7.netty.model.payloads.S7Payload;
import org.apache.plc4x.java.s7.netty.model.payloads.VarPayload;
import org.apache.plc4x.java.s7.netty.model.payloads.items.VarPayloadItem;
import org.apache.plc4x.java.s7.netty.model.payloads.ssls.SslDataRecord;
import org.apache.plc4x.java.s7.netty.model.payloads.ssls.SslModuleIdentificationDataRecord;
import org.apache.plc4x.java.s7.netty.model.types.*;
import org.apache.plc4x.java.s7.netty.strategies.S7MessageProcessor;
import org.apache.plc4x.java.s7.netty.util.S7SizeHelper;
import org.apache.plc4x.java.s7.types.S7ControllerType;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.Field;
import java.nio.charset.StandardCharsets;
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
    private S7ControllerType controllerType;

    // For detecting the lower layers.
    private ChannelHandler prevChannelHandler;
    private S7MessageProcessor messageProcessor;

    // For being able to respect the max AMQ restrictions.
    private PendingWriteQueue queue;
    private Map<Short, DataTpdu> sentButUnacknowledgedTpdus;

    public S7Protocol(short requestedMaxAmqCaller, short requestedMaxAmqCallee, short requestedPduSize,
                      S7ControllerType controllerType, S7MessageProcessor messageProcessor) {
        this.maxAmqCaller = requestedMaxAmqCaller;
        this.maxAmqCallee = requestedMaxAmqCallee;
        this.pduSize = requestedPduSize;
        this.controllerType = controllerType;
        this.messageProcessor = messageProcessor;
        sentButUnacknowledgedTpdus = new HashMap<>();
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
    public void write(ChannelHandlerContext ctx, Object msg, ChannelPromise promise) {
        try {
            if(msg instanceof S7Message) {
                S7Message in = (S7Message) msg;

                // Give message processors to process the incoming message.
                Collection<? extends S7Message> messages;
                if ((messageProcessor != null) && (in instanceof S7RequestMessage)) {
                    messages = messageProcessor.processRequest((S7RequestMessage) in, pduSize);
                } else {
                    messages = Collections.singleton(in);
                }

                // Create a promise that has to be called multiple times.
                PromiseCombiner promiseCombiner = new PromiseCombiner();
                for (S7Message message : messages) {
                    ByteBuf buf = Unpooled.buffer();
                    writeS7Message(promise.channel(), promiseCombiner, message, buf);
                }
                promiseCombiner.finish(promise);

                // Start sending the queue content.
                trySendingMessages(ctx);
            }
            // Especially during the phase of connection establishment, we might be sending
            // messages of a lower level protocol, so if it's not S7, we forward it to the next
            // in the pipeline and hope it can handle it. If no layer can handle it Netty will
            // exceptionally complete the future.
            else {
                ctx.write(msg, promise);
            }
        } catch (Exception e) {
            promise.setFailure(e);
        }
    }

    private void writeS7Message(Channel channel, PromiseCombiner promiseCombiner,
                                S7Message message, ByteBuf buf) throws PlcProtocolException {
        encodeHeader(message, buf);
        encodeParameters(message, buf);
        encodePayloads(message, buf);

        // Check if the message doesn't exceed the negotiated maximum size.
        if (buf.writerIndex() > pduSize) {
            throw new PlcProtocolPayloadTooBigException("s7", pduSize, buf.writerIndex(), message);
        } else {
            ChannelPromise subPromise = new DefaultChannelPromise(channel);
            // The tpduRef was 0x01 but had to be changed to 0x00 in order to support Siemens LOGO devices.
            queue.add(new DataTpdu(true, (byte) 0x00, Collections.emptyList(), buf, message), subPromise);
            promiseCombiner.add((Future) subPromise);
            logger.debug("S7 Message with id {} queued", message.getTpduReference());
        }
    }

    private void encodePayloads(S7Message in, ByteBuf buf) throws PlcProtocolException {
        if(in.getPayloads() != null) {
            Iterator<S7Payload> payloadIterator = in.getPayloads().iterator();
            while(payloadIterator.hasNext()) {
                S7Payload payload = payloadIterator.next();
                switch (payload.getType()) {
                    case WRITE_VAR:
                        encodeWriteVarPayload((VarPayload) payload, buf, !payloadIterator.hasNext());
                        break;
                    case CPU_SERVICES:
                        encodeCpuServicesPayload((CpuServicesPayload) payload, buf);
                        break;
                    default:
                        throw new PlcProtocolException("Writing payloads of type " +
                            payload.getType().name() + " not implemented.");
                }
            }
        }
    }

    private void encodeWriteVarPayload(VarPayload varPayload, ByteBuf buf, boolean lastItem) {
        for (VarPayloadItem payloadItem : varPayload.getItems()) {
            buf.writeByte(payloadItem.getReturnCode().getCode());
            buf.writeByte(payloadItem.getDataTransportSize().getCode());
            // TODO: Check if this is correct?!?! Might be problems with sizeInBits = true/false
            buf.writeShort(payloadItem.getData().length);
            buf.writeBytes(payloadItem.getData());
            // if this is not the last item and it's payload is exactly one byte, we need to output a fill-byte.
            if((payloadItem.getData().length == 1) && !lastItem) {
                buf.writeByte(0x00);
            }
        }
    }

    private void encodeCpuServicesPayload(CpuServicesPayload cpuServicesPayload, ByteBuf buf)
            throws PlcProtocolException {
        buf.writeByte(cpuServicesPayload.getReturnCode().getCode());
        // This seems to be constantly set to this.
        buf.writeByte(DataTransportSize.OCTET_STRING.getCode());

        // A request payload is simple.
        if (cpuServicesPayload.getSslDataRecords().isEmpty()) {
            buf.writeShort(4);
            buf.writeShort(cpuServicesPayload.getSslId().getCode());
            buf.writeShort(cpuServicesPayload.getSslIndex());
        }
        // The response payload contains a lot more information.
        else {
            throw new PlcProtocolException("Unexpected SZL Data Records");
            /*short length = 8;
            short sizeOfDataItem = 0;
            for (SslDataRecord sslDataRecord : cpuServicesPayload.getSslDataRecords()) {
                sizeOfDataItem = (short) (sslDataRecord.getLengthInWords() * (short) 2);
                length += sizeOfDataItem;
            }
            buf.writeShort(length);
            buf.writeShort(cpuServicesPayload.getSslId().getCode());
            buf.writeShort(cpuServicesPayload.getSslIndex());
            buf.writeShort(sizeOfDataItem);
            buf.writeShort(cpuServicesPayload.getSslDataRecords().size());
            // Output any sort of ssl list items, if there are any.
            for (SslDataRecord sslDataRecord : cpuServicesPayload.getSslDataRecords()) {
                if(sslDataRecord instanceof SslModuleIdentificationDataRecord) {
                    SslModuleIdentificationDataRecord midr = (SslModuleIdentificationDataRecord) sslDataRecord;
                    buf.writeShort(midr.getIndex());
                    byte[] articleNumberBytes = midr.getArticleNumber().getBytes(StandardCharsets.UTF_8);
                    // An array full of 20 spaces.
                    byte[] data = new byte[]{0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
                        0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20};
                    // Copy max 20 bytes from the article number into the dest array.
                    System.arraycopy(articleNumberBytes, 0, data, 0, 20);
                    buf.writeBytes(data);
                    buf.writeShort(midr.getModuleOrOsVersion());
                    buf.writeShort(midr.getPgDescriptionFileVersion());
                }
            }*/
        }
    }

    private void encodeParameters(S7Message in, ByteBuf buf) throws PlcProtocolException {
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
                case CPU_SERVICES:
                    encodeCpuServicesParameter(buf, (CpuServicesParameter) s7Parameter);
                    break;
                default:
                    throw new PlcProtocolException("Writing parameters of type " +
                        s7Parameter.getType().name() + " not implemented.");
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
        // Not sure why this is implemented, we should never be sending out responses.
        /*if (in instanceof S7ResponseMessage) {
            S7ResponseMessage s7ResponseMessage = (S7ResponseMessage) in;
            buf.writeByte(s7ResponseMessage.getErrorClass());
            buf.writeByte(s7ResponseMessage.getErrorCode());
        }*/
    }

    private void encodeParameterSetupCommunication(ByteBuf buf, SetupCommunicationParameter s7Parameter) {
        // Reserved (is always constant 0x00)
        buf.writeByte((byte) 0x00);
        buf.writeShort(s7Parameter.getMaxAmqCaller());
        buf.writeShort(s7Parameter.getMaxAmqCallee());
        buf.writeShort(s7Parameter.getPduLength());
    }

    private void encodeParameterReadWriteVar(ByteBuf buf, VarParameter s7Parameter) throws PlcProtocolException {
        List<VarParameterItem> items = s7Parameter.getItems();
        // PlcReadRequestItem count (Read one variable at a time)
        buf.writeByte((byte) items.size());
        for (VarParameterItem item : items) {
            VariableAddressingMode addressMode = item.getAddressingMode();
            if (addressMode == VariableAddressingMode.S7ANY) {
                encodeS7AnyParameterItem(buf, (S7AnyVarParameterItem) item);
            } else {
                throw new PlcProtocolException("Writing VarParameterItems with addressing mode " +
                    addressMode.name() + " not implemented");
            }
        }
    }

    private void encodeCpuServicesParameter(ByteBuf buf, CpuServicesParameter parameter) {
        // Output the header for a CPU Services parameter.
        buf.writeByte(0x01);
        buf.writeByte(0x12);
        // Length of the parameter.
        buf.writeByte((parameter instanceof CpuServicesRequestParameter) ? 0x04 : 0x08);
        // Is this a request or a response?
        buf.writeByte((parameter instanceof CpuServicesRequestParameter) ? 0x11 : 0x12);
        // This is a mixture of request/response and function group .
        byte nextByte = (byte) (((parameter instanceof CpuServicesRequestParameter) ?
            (byte) 0x40 : (byte) 0x80) | parameter.getFunctionGroup().getCode());
        buf.writeByte(nextByte);
        buf.writeByte(parameter.getSubFunctionGroup().getCode());
        buf.writeByte(parameter.getSequenceNumber());

        // A response parameter has some more fields.
        // Not sure why this is implemented, we should never be sending out responses.
        /*if(parameter instanceof CpuServicesResponseParameter) {
            CpuServicesResponseParameter responseParameter = (CpuServicesResponseParameter) parameter;
            buf.writeByte(responseParameter.getDataUnitReferenceNumber());
            buf.writeByte(responseParameter.isLastDataUnit() ? 0x00 : 0x01);
            buf.writeShort(responseParameter.getError().getCode());
        }*/
    }

    private void encodeS7AnyParameterItem(ByteBuf buf, S7AnyVarParameterItem s7AnyRequestItem) {
        buf.writeByte(s7AnyRequestItem.getSpecificationType().getCode());
        // Length of this item (excluding spec type and length)
        buf.writeByte((byte) 0x0a);
        buf.writeByte(s7AnyRequestItem.getAddressingMode().getCode());
        buf.writeByte(s7AnyRequestItem.getDataType().getTypeCode());
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
        int i = 0;

        while (i < headerParametersLength) {
            S7Parameter parameter = decodeParameter(userData, isResponse);
            s7Parameters.add(parameter);
            if (parameter instanceof SetupCommunicationParameter) {
                handleSetupCommunications(ctx, (SetupCommunicationParameter) parameter);
            }
            i += S7SizeHelper.getParameterLength(parameter);
        }

        List<S7Payload> s7Payloads = decodePayloads(userData, isResponse, userDataLength, s7Parameters);

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

                    // If this is a USER_DATA packet the probability is high that this is
                    // a response to the identification request, we have to handle that.
                    if(responseMessage.getMessageType() == MessageType.USER_DATA) {
                        for (S7Payload payload : responseMessage.getPayloads()) {
                            if(payload instanceof CpuServicesPayload) {
                                handleIdentifyRemote(ctx, (CpuServicesPayload) payload);
                            }
                        }
                    }
                }

                // Try to send the next message (if there is one).
                trySendingMessages(ctx);
            }

        } else {
            // CpuService responses are encoded as requests.
            for (S7Parameter s7Parameter : s7Parameters) {
                // Only if we have a response parameter, the payload is a response payload.
                if(s7Parameter instanceof CpuServicesResponseParameter) {
                    for (S7Payload s7Payload : s7Payloads) {
                        if(s7Payload instanceof CpuServicesPayload) {
                            CpuServicesPayload cpuServicesPayload = (CpuServicesPayload) s7Payload;
                            handleIdentifyRemote(ctx, cpuServicesPayload);
                        }
                    }
                }
            }
            out.add(new S7RequestMessage(messageType, tpduReference, s7Parameters, s7Payloads, null));
        }
    }

    private void handleSetupCommunications(ChannelHandlerContext ctx, SetupCommunicationParameter setupCommunicationParameter) {
        maxAmqCaller = setupCommunicationParameter.getMaxAmqCaller();
        maxAmqCallee = setupCommunicationParameter.getMaxAmqCallee();
        pduSize = setupCommunicationParameter.getPduLength();

        logger.info("S7Connection established pdu-size {}, max-amq-caller {}, " +
                "max-amq-callee {}", pduSize, maxAmqCaller, maxAmqCallee);

        // Only if the controller type is set to "ANY", then try to identify the PLC type.
        if(controllerType == S7ControllerType.ANY) {
            // Prepare a message to request the remote to identify itself.
            S7RequestMessage identifyRemoteMessage = new S7RequestMessage(MessageType.USER_DATA, (short) 2,
                Collections.singletonList(new CpuServicesRequestParameter(
                    CpuServicesParameterFunctionGroup.CPU_FUNCTIONS,
                    CpuServicesParameterSubFunctionGroup.READ_SSL, (byte) 0)),
                Collections.singletonList(new CpuServicesPayload(DataTransportErrorCode.OK, SslId.MODULE_IDENTIFICATION,
                    (short) 0x0000)), null);
            ctx.channel().writeAndFlush(identifyRemoteMessage);
        }
        // If a concrete type was specified, then we're done here.
        else {
            if(logger.isInfoEnabled()) {
                logger.info(String.format("Successfully connected to S7: %s", controllerType.name()));
                logger.info(String.format("- max amq caller: %s", maxAmqCaller));
                logger.info(String.format("- max amq callee: %s", maxAmqCallee));
                logger.info(String.format("- pdu size: %s", pduSize));
            }

            // Send an event that connection setup is complete.
            ctx.channel().pipeline().fireUserEventTriggered(new S7ConnectedEvent());
        }
    }

    private void handleIdentifyRemote(ChannelHandlerContext ctx, CpuServicesPayload cpuServicesPayload) {
        controllerType = S7ControllerType.ANY;
        for (SslDataRecord sslDataRecord : cpuServicesPayload.getSslDataRecords()) {
            if(sslDataRecord instanceof SslModuleIdentificationDataRecord) {
                SslModuleIdentificationDataRecord sslModuleIdentificationDataRecord =
                    (SslModuleIdentificationDataRecord) sslDataRecord;
                if(sslModuleIdentificationDataRecord.getIndex() == (short) 0x0001) {
                    controllerType = lookupControllerType(sslModuleIdentificationDataRecord.getArticleNumber());
                }
            }
        }
        if(logger.isInfoEnabled()) {
            logger.info(String.format("Successfully connected to S7: %s", controllerType.name()));
            logger.info(String.format("- max amq caller: %s", maxAmqCaller));
            logger.info(String.format("- max amq callee: %s", maxAmqCallee));
            logger.info(String.format("- pdu size: %s", pduSize));
        }

        // Send an event that connection setup is complete.
        ctx.channel().pipeline().fireUserEventTriggered(new S7ConnectedEvent());
    }

    private List<S7Payload> decodePayloads(ByteBuf userData, boolean isResponse, short userDataLength, List<S7Parameter> s7Parameters) {
        List<S7Payload> s7Payloads = new LinkedList<>();
        for (S7Parameter s7Parameter : s7Parameters) {
            if(s7Parameter instanceof VarParameter) {
                VarParameter readWriteVarParameter = (VarParameter) s7Parameter;
                VarPayload varPayload = decodeVarPayload(userData, isResponse, userDataLength, readWriteVarParameter);
                s7Payloads.add(varPayload);
            } else if(s7Parameter instanceof CpuServicesParameter) {
                CpuServicesPayload cpuServicesPayload = decodeCpuServicesPayload(userData);
                s7Payloads.add(cpuServicesPayload);
            }
        }
        return s7Payloads;
    }

    private VarPayload decodeVarPayload(ByteBuf userData, boolean isResponse, short userDataLength,
                                        VarParameter readWriteVarParameter) {
        List<VarPayloadItem> payloadItems = new LinkedList<>();

        // Just keep on reading payloads until the provided length is read.
        int i = 0;
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
                short length = dataTransportSize.isSizeInBits() ?
                    (short) Math.ceil(userData.readShort() / 8.0) : userData.readShort();
                byte[] data = new byte[length];
                userData.readBytes(data);
                // Initialize a rudimentary payload (This is updated in the Plc4XS7Protocol class
                VarPayloadItem payload = new VarPayloadItem(dataTransportErrorCode, dataTransportSize, data);
                payloadItems.add(payload);
                i += S7SizeHelper.getPayloadLength(payload);

                // It seems that odd-byte payloads require a fill byte, but only if it's not the last item.
                if((length % 2== 1) && (userData.readableBytes() > 0)) {
                    userData.readByte();
                    i++;
                }
            }
        }

        return new VarPayload(readWriteVarParameter.getType(), payloadItems);
    }

    private CpuServicesPayload decodeCpuServicesPayload(ByteBuf userData) {
        DataTransportErrorCode returnCode = DataTransportErrorCode.valueOf(userData.readByte());
        DataTransportSize dataTransportSize = DataTransportSize.valueOf(userData.readByte());
        if(dataTransportSize != DataTransportSize.OCTET_STRING) {
            // TODO: Output an error.
        }
        short length = userData.readShort();
        SslId sslId = SslId.valueOf(userData.readShort());
        short sslIndex = userData.readShort();
        // If the length is 4 there is no `partial list length in bytes` and `partial list count` parameters.
        if(length == 4) {
            return new CpuServicesPayload(returnCode, sslId, sslIndex);
        }
        // If the length is not 4, then it has to be at least 8.
        else if(length >= 8) {
            // TODO: We should probably ensure we don't read more than this.
            // Skip the partial list length in words.
            userData.skipBytes(2);
            short partialListCount = userData.readShort();
            List<SslDataRecord> sslDataRecords = new LinkedList<>();
            for(int i = 0; i < partialListCount; i++) {
                short index = userData.readShort();
                byte[] articleNumberBytes = new byte[20];
                userData.readBytes(articleNumberBytes);
                String articleNumber = new String(articleNumberBytes, StandardCharsets.UTF_8).trim();
                short bgType = userData.readShort();
                short moduleOrOsVersion = userData.readShort();
                short pgDescriptionFileVersion = userData.readShort();
                sslDataRecords.add(new SslModuleIdentificationDataRecord(
                    index, articleNumber, bgType, moduleOrOsVersion, pgDescriptionFileVersion));
            }
            return new CpuServicesPayload(returnCode, sslId, sslIndex, sslDataRecords);
        }
        // In all other cases, it's probably an error.
        else {
            // TODO: Output an error.
        }
        return null;
    }

    private S7Parameter decodeParameter(ByteBuf in, boolean isResponse) {
        ParameterType parameterType = ParameterType.valueOf(in.readByte());
        if (parameterType == null) {
            logger.error("Could not find parameter type");
            return null;
        }
        switch (parameterType) {
            case CPU_SERVICES:
                return decodeCpuServicesParameter(in);
            case READ_VAR:
            case WRITE_VAR:
                List<VarParameterItem> varParameterItems;
                byte numItems = in.readByte();
                if (!isResponse) {
                    varParameterItems = decodeReadWriteVarParameter(in, numItems);
                } else {
                    varParameterItems = Collections.singletonList(
                        new S7AnyVarParameterItem(null, null, null, numItems, (short) 0, (short) 0, (byte) 0));
                }
                return new VarParameter(parameterType, varParameterItems);
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

    private CpuServicesParameter decodeCpuServicesParameter(ByteBuf in) {
        if(in.readShort() != 0x0112) {
            if (logger.isErrorEnabled()) {
                logger.error("Expecting 0x0112 for CPU_SERVICES parameter");
            }
            return null;
        }
        byte parameterLength = in.readByte();
        if((parameterLength != 4) && (parameterLength != 8)) {
            if (logger.isErrorEnabled()) {
                logger.error("Parameter length should be 4 or 8, but was {}", parameterLength);
            }
            return null;
        }
        // Skipping this as it sort of contains redundant information.
        in.readByte();
        byte typeAndFunctionGroup = in.readByte();
        // If bit 7 is set, it's a request (if bit 8 is set it's a response).
        boolean requestParameter = (typeAndFunctionGroup & 0x64) != 0;
        // The last 4 bits contain the function group value.
        typeAndFunctionGroup = (byte) (typeAndFunctionGroup & 0xF);
        CpuServicesParameterFunctionGroup functionGroup =
            CpuServicesParameterFunctionGroup.valueOf(typeAndFunctionGroup);
        CpuServicesParameterSubFunctionGroup subFunctionGroup =
            CpuServicesParameterSubFunctionGroup.valueOf(in.readByte());
        byte sequenceNumber = in.readByte();
        if(!requestParameter) {
            return new CpuServicesRequestParameter(functionGroup, subFunctionGroup, sequenceNumber);
        } else {
            byte dataUnitReferenceNumber = in.readByte();
            boolean lastDataUnit = in.readByte() == 0x00;
            ParameterError error = ParameterError.valueOf(in.readShort());
            return new CpuServicesResponseParameter(functionGroup, subFunctionGroup, sequenceNumber,
                dataUnitReferenceNumber, lastDataUnit, error);
        }
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
                TransportSize dataType = TransportSize.valueOf(in.readByte());
                short length = in.readShort();
                short dbNumber = in.readShort();
                byte memoryAreaCode = in.readByte();
                MemoryArea memoryArea = MemoryArea.valueOf(memoryAreaCode);
                if(memoryArea == null) {
                    throw new PlcRuntimeException("Unknown memory area '" + memoryAreaCode + "'");
                }
                short byteAddress = (short) (in.readShort() << 5);
                byte tmp = in.readByte();
                // Only the least 3 bits are the bit address, the
                byte bitAddress = (byte) (tmp & 0x07);
                // Bits 4-8 belong to the byte address
                byteAddress = (short) (byteAddress | (tmp >> 3));
                S7AnyVarParameterItem item = new S7AnyVarParameterItem(
                        specificationType, memoryArea, dataType,
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
                // TODO: Perhaps remove this.
                break;
            } else {
                break;
            }
        }
        ctx.flush();
    }

    private S7ControllerType lookupControllerType(String articleNumber) {
        if(!articleNumber.startsWith("6ES7 ")) {
            return S7ControllerType.ANY;
        }

        String model = articleNumber.substring(articleNumber.indexOf(' ') + 1, articleNumber.indexOf(' ') + 2);
        switch (model) {
            case "2":
                return S7ControllerType.S7_1200;
            case "5":
                return S7ControllerType.S7_1500;
            case "3":
                return S7ControllerType.S7_300;
            case "4":
                return S7ControllerType.S7_400;
            default:
                if(logger.isInfoEnabled()) {
                    logger.info(String.format("Looking up unknown article number %s", articleNumber));
                }
                return S7ControllerType.ANY;
        }
    }

}
