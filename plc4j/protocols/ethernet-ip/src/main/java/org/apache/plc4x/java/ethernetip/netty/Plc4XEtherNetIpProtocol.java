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
package org.apache.plc4x.java.ethernetip.netty;

import com.digitalpetri.enip.EnipPacket;
import com.digitalpetri.enip.EnipStatus;
import com.digitalpetri.enip.cip.epath.EPath;
import com.digitalpetri.enip.cip.epath.LogicalSegment;
import com.digitalpetri.enip.cip.services.GetAttributeSingleService;
import com.digitalpetri.enip.cip.structs.MessageRouterResponse;
import com.digitalpetri.enip.commands.*;
import com.digitalpetri.enip.cpf.*;
import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.MessageToMessageCodec;
import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcRequest;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.events.ConnectEvent;
import org.apache.plc4x.java.base.events.ConnectedEvent;
import org.apache.plc4x.java.base.messages.*;
import org.apache.plc4x.java.base.messages.items.BaseDefaultFieldItem;
import org.apache.plc4x.java.base.messages.items.DefaultLongFieldItem;
import org.apache.plc4x.java.ethernetip.model.EtherNetIpField;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.nio.charset.Charset;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicLong;

public class Plc4XEtherNetIpProtocol extends MessageToMessageCodec<EnipPacket, PlcRequestContainer<InternalPlcRequest, InternalPlcResponse>> {

    private static final Logger LOGGER = LoggerFactory.getLogger(Plc4XEtherNetIpProtocol.class);

    private static final int SERVICE_COMMUNICATIONS_TYPE_CODE = 0x0100;

    private long sessionHandle = 0;
    private static final AtomicLong messageId = new AtomicLong();

    // General information about the remote communication endpoint.
    private CipIdentityItem identityItem;
    // Flag to signal, if the remote communication endpoint supports encapsulation of CIP data.
    private boolean supportsCipEncapsulation = false;
    // Flag to indicate, if implicit IO (subscription) is generally supported by the remote communication endpoint.
    // This is handled via separate UDP socket, which would have to be established in parallel.
    private boolean supportsClass0Or1UdpConnections = false;
    // Map of non-cip interfaces, that might be used for specialized IO in future versions.
    private Map<String, Integer> nonCipInterfaces = null;
    // In CIP we are doing explicit connected messaging, this requires every used address to be registered at the
    // remote server and to use that Addresses connectionId for accessing data. We are saving the references to
    // these here.
    // REMARK: Perhaps we should add a timeout to these so we unregister them after not being used
    // for quire some time. Hereby freeing resources on both client and server.
    private Map<PlcField, Long> fieldConnectionMap = new ConcurrentHashMap<>();

    private final Map<Long, PlcRequestContainer<InternalPlcRequest, InternalPlcResponse>> requestsMap = new ConcurrentHashMap<>();

    /**
     * If the IsoTP protocol is used on top of the ISO on TCP protocol, then as soon as the pipeline receives the
     * request to connect, an IsoTP connection request TPDU must be sent in order to initialize the connection.
     *
     * @param ctx the current protocol layers context
     * @param evt the event
     * @throws Exception throws an exception if something goes wrong internally
     */
    @Override
    public void userEventTriggered(ChannelHandlerContext ctx, Object evt) throws Exception {
        // If the connection has just been established, start setting up the connection
        // by sending a connection request to the plc.
        if (evt instanceof ConnectEvent) {
            LOGGER.debug("EtherNet/IP Protocol Sending Connection Request");

            EnipPacket packet = new EnipPacket(CommandCode.RegisterSession, 0, EnipStatus.EIP_SUCCESS,
                messageId.getAndIncrement(), new RegisterSession());

            ctx.channel().writeAndFlush(packet);
        } else {
            super.userEventTriggered(ctx, evt);
        }
    }

    @Override
    public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) throws Exception {
        LOGGER.trace("(-->ERR): {}", ctx, cause);
        super.exceptionCaught(ctx, cause);
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Encoding
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected void encode(ChannelHandlerContext ctx, PlcRequestContainer<InternalPlcRequest, InternalPlcResponse> msg, List<Object> out) {
        LOGGER.trace("(<--OUT): {}, {}, {}", ctx, msg, out);
        // Reset transactionId on overflow
        messageId.compareAndSet(Short.MAX_VALUE + 1L, 0);
        PlcRequest request = msg.getRequest();
        if (request instanceof PlcReadRequest) {
            encodeReadRequest(msg, out);
        } else if (request instanceof PlcWriteRequest) {
            encodeWriteRequest(msg, out);
        } /*else if(request instanceof PlcSubscriptionRequest) {
            encodeSubscriptionRequest(msg, out);
        } else if(request instanceof PlcUnsubscriptionRequest) {
            TODO: Implement this and refactor PlcUnsubscriptionRequest first ...
        }*/
    }

    private void encodeWriteRequest(PlcRequestContainer<InternalPlcRequest, InternalPlcResponse> msg, List<Object> out) {
        if (!supportsCipEncapsulation) {
            LOGGER.warn("CIP Encapsulation not supported by remote, payload encapsulation must be handled by target and originator");
        }

        /*PlcWriteRequest request = (PlcWriteRequest) msg.getRequest();

        // Create a ForwardOpen CIP request

        // Create EIP UnconnectedDataItemRequest
        UnconnectedDataItemRequest dataItem = new UnconnectedDataItemRequest(dataEncoder);
        CpfPacket packet = new CpfPacket(new NullAddressItem(), dataItem);

        // Send that via EIP SendRRData packet
        CompletableFuture<T> future = new CompletableFuture<>();

        sendRRData(new SendRRData(packet)).whenComplete((command, ex) -> {
            if (command != null) {
                CpfItem[] items = command.getPacket().getItems();

                if (items.length == 2 &&
                    items[0].getTypeId() == NullAddressItem.TYPE_ID &&
                    items[1].getTypeId() == UnconnectedDataItemResponse.TYPE_ID) {

                    ByteBuf data = ((UnconnectedDataItemResponse) items[1]).getData();

                    future.complete(data);
                } else {
                    future.completeExceptionally(new Exception("received unexpected items"));
                }
            } else {
                future.completeExceptionally(ex);
            }
        });

        channelManager.getChannel().whenComplete((ch, ex) -> {
            if (ch != null) writeCommand(ch, command, future);
            else future.completeExceptionally(ex);
        });*/

    }

    private void encodeReadRequest(PlcRequestContainer<InternalPlcRequest, InternalPlcResponse> msg, List<Object> out) {
        if (!supportsCipEncapsulation) {
            LOGGER.warn("CIP Encapsulation not supported by remote, payload encapsulation must be handled by target and originator");
        }

        InternalPlcReadRequest request = (InternalPlcReadRequest) msg.getRequest();
        for (String fieldName : request.getFieldNames()) {
            PlcField field = request.getField(fieldName);

            // CIP Part
            EtherNetIpField enipField = (EtherNetIpField) field;
            EPath.PaddedEPath path = new EPath.PaddedEPath(new LogicalSegment.ClassId(enipField.getObjectNumber()),
                new LogicalSegment.InstanceId(enipField.getInstanceNumber()),
                new LogicalSegment.AttributeId(enipField.getAttributeNumber()));
            GetAttributeSingleService service = new GetAttributeSingleService(path);

            // ENIP Part
            EnipPacket packet = new EnipPacket(CommandCode.SendRRData, sessionHandle, EnipStatus.EIP_SUCCESS,
                messageId.getAndIncrement(), new SendRRData(new CpfPacket(
                new NullAddressItem(),
                new UnconnectedDataItemRequest(service::encodeRequest)
            )));

            requestsMap.put(packet.getSenderContext(), msg);

            out.add(packet);
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Decoding
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @SuppressWarnings("unchecked")
    @Override
    protected void decode(ChannelHandlerContext ctx, EnipPacket msg, List<Object> out) {
        LOGGER.trace("(-->IN): {}, {}, {}", ctx, msg, out);
        LOGGER.debug("{}: session handle: {}, sender context: {}, EtherNetIPPacket:{}", msg, msg.getSessionHandle(), msg.getSenderContext(), msg);

        EnipPacket packet = null;
        switch (msg.getCommandCode()) {
            case RegisterSession:
                handleRegisterSession(ctx, msg);

                // Now try getting some detailed information about the remote.
                packet = new EnipPacket(CommandCode.ListIdentity, sessionHandle, EnipStatus.EIP_SUCCESS,
                    messageId.getAndIncrement(), new ListIdentity());
                break;

            case UnRegisterSession:
                handleUnregisterSession(ctx, msg);

                // Spec: The receiver shall initiate a close of the underlying
                // TCP/IP connection when it receives this command.
                ctx.channel().disconnect();
                break;

            case ListIdentity:
                handleListIdentity(ctx, msg);

                // Now try listing the services the remote has to offer.
                packet = new EnipPacket(CommandCode.ListServices, sessionHandle, EnipStatus.EIP_SUCCESS,
                    messageId.getAndIncrement(), new ListServices());
                break;

            case ListInterfaces:
                handleListInterfaces(ctx, msg);

                // Here we're done connecting.
                ctx.channel().pipeline().fireUserEventTriggered(new ConnectedEvent());
                break;

            case ListServices:
                handleListServices(ctx, msg);

                // Now try listing the interfaces the remote has to offer.
                packet = new EnipPacket(CommandCode.ListInterfaces, sessionHandle, EnipStatus.EIP_SUCCESS,
                    messageId.getAndIncrement(), new ListInterfaces());
                break;

            case Nop:
                handleNop(ctx, msg);
                break;

            case SendRRData:
                handleSendRRDataResponse(ctx, msg);
                break;

            case SendUnitData:
                // This might be where the connected data is sent (eventually publish/subscribe communication)
                break;
        }

        if (packet != null) {
            ctx.channel().writeAndFlush(packet);
        }
    }

    /**
     * In order to do explicit connected messaging, the client has to register a session with the server.
     * In case of a successful session registration the response will contain the sessionHandle, which is
     * required to be used in all subsequent connected interactions.
     *
     * @param ctx the {@link ChannelHandlerContext} instance.
     * @param msg the packet received from the server.
     */
    private void handleRegisterSession(ChannelHandlerContext ctx, EnipPacket msg) {
        if (msg.getStatus() == EnipStatus.EIP_SUCCESS) {
            sessionHandle = msg.getSessionHandle();

            LOGGER.info("EtherNet/IP session registered session-handle {}", sessionHandle);
        } else {
            ctx.channel().pipeline().fireExceptionCaught(new PlcProtocolException("Got a non-success response."));
        }
    }

    /**
     * As connected operations allocate resources on the server and the client, when receiving a
     * {@link UnRegisterSession} message (request or response) the locally allocated resources have
     * to be released again. As the correct response to a UnRegisterSession is the closing of the
     * connection by the receiving side, this incoming command must be a request sent from the
     * server.
     *
     * @param ctx the {@link ChannelHandlerContext} instance.
     * @param msg the packet received from the server.
     */
    private void handleUnregisterSession(ChannelHandlerContext ctx, EnipPacket msg) {
        if (msg.getStatus() == EnipStatus.EIP_SUCCESS) {
            // Reset all internal variables.
            identityItem = null;
            supportsCipEncapsulation = false;
            supportsClass0Or1UdpConnections = false;
            nonCipInterfaces = null;
            fieldConnectionMap = null;
        } else {
            ctx.channel().pipeline().fireExceptionCaught(new PlcProtocolException("Got a non-success response."));
        }
    }

    /**
     * The response to a {@link ListIdentity} command contains a lot of information about the
     * remote counterpart. In this case we just save this information for further usage.
     *
     * @param ctx the {@link ChannelHandlerContext} instance.
     * @param msg the packet received from the server.
     */
    private void handleListIdentity(ChannelHandlerContext ctx, EnipPacket msg) {
        if (msg.getStatus() == EnipStatus.EIP_SUCCESS) {
            ListIdentity listIdentityResponse = (ListIdentity) msg.getCommand();
            if (listIdentityResponse != null) {
                identityItem = listIdentityResponse.getIdentity().orElse(null);
                if (identityItem != null) {
                    LOGGER.info("Connected to: \n - product name: {} \n - serial number: {} ",
                        identityItem.getProductName().trim(), identityItem.getSerialNumber());
                }
            } else {
                identityItem = null;
            }
        } else {
            ctx.channel().pipeline().fireExceptionCaught(new PlcProtocolException("Got a non-success response."));
        }
    }

    /**
     * Some times EtherNet/IP devices support other devices than the default one.
     * As we are required to ev eventually reference these interfaces, build a map
     * of all the devices the remote supports. This way we can check the validity
     * before actually sending a request.
     *
     * @param ctx the {@link ChannelHandlerContext} instance.
     * @param msg the packet received from the server.
     */
    private void handleListInterfaces(ChannelHandlerContext ctx, EnipPacket msg) {
        if (msg.getStatus() == EnipStatus.EIP_SUCCESS) {
            ListInterfaces listInterfaces = (ListInterfaces) msg.getCommand();
            if (listInterfaces != null) {
                // If the device supports non-CIP interfaces, this array is not empty.
                // In this case build a map so we can access the information when sending
                // data in RR-Requests (Request-Response).
                if (listInterfaces.getInterfaces().length > 0) {
                    nonCipInterfaces = new HashMap<>();
                    for (ListInterfaces.InterfaceInformation interfaceInformation : listInterfaces.getInterfaces()) {
                        String interfaceName = new String(
                            interfaceInformation.getData(), Charset.forName("US-ASCII")).trim();
                        nonCipInterfaces.put(interfaceName, interfaceInformation.hashCode());
                    }
                } else {
                    nonCipInterfaces = null;
                }
            }
        } else {
            ctx.channel().pipeline().fireExceptionCaught(new PlcProtocolException("Got a non-success response."));
        }
    }

    /**
     * Each EtherNet/IP device can support one or more so-called `services`. At least the `Communications`
     * service is required to be supported by every EtherNet/IP compliant device. This is used for default
     * IO operations. Usually vendors support custom services which are adjusted to their particular needs,
     * which might be able to provide better performance than the default. In this case we are ignoring all
     * these as supporting these would require custom adapters on the PLC4X side. However we do inspect the
     * capabilities of the `Communications` service to check if encapsulation of CIP data is supported and
     * if we are able to do connected implicit communication via a parallel UDP channel.
     *
     * @param ctx the {@link ChannelHandlerContext} instance.
     * @param msg the packet received from the server.
     */
    private void handleListServices(ChannelHandlerContext ctx, EnipPacket msg) {
        if (msg.getStatus() == EnipStatus.EIP_SUCCESS) {
            ListServices listServices = (ListServices) msg.getCommand();
            if (listServices != null) {
                for (ListServices.ServiceInformation service : listServices.getServices()) {
                    // Check if the type code matches the communications service and if bit 5 of the
                    // capability flags is set.
                    if (service.getTypeCode() == SERVICE_COMMUNICATIONS_TYPE_CODE) {
                        supportsCipEncapsulation = (service.getCapabilityFlags() & 32) != 0;
                        supportsClass0Or1UdpConnections = (service.getCapabilityFlags() & 256) != 0;
                    }
                }
            } else {
                supportsCipEncapsulation = false;
                supportsClass0Or1UdpConnections = false;
            }
        } else {
            ctx.channel().pipeline().fireExceptionCaught(new PlcProtocolException("Got a non-success response."));
        }
    }

    /**
     * NOP request/responses are simple no-payload messages used to check if a connection is still
     * available. Depending on if it's a request or reply, we simply send back a NOP Reply or not.
     * As no reply is to be generated for an incoming NOP command, this must be a NopRequest.
     *
     * @param ctx the {@link ChannelHandlerContext} instance.
     * @param msg the packet received from the server.
     */
    private void handleNop(ChannelHandlerContext ctx, EnipPacket msg) {
        if (msg.getStatus() == EnipStatus.EIP_SUCCESS) {
            //Nop nop = (Nop) msg.getCommand();
            // TODO: Reset some sort of timer ...
        } else {
            ctx.channel().pipeline().fireExceptionCaught(
                new PlcProtocolException("Got a non-success flagged request."));
        }
    }

    /**
     * As RR Data is Request Response data and the server will not issue a request to
     * the client, we can be pretty sure this is a response to a previously issued request.
     * This contains the actual payload for our requests.
     *
     * @param ctx the {@link ChannelHandlerContext} instance.
     * @param msg the packet received from the server.
     */
    private void handleSendRRDataResponse(ChannelHandlerContext ctx, EnipPacket msg) {
        // This is where the typical request/response stuff is handled.
        long senderContext = msg.getSenderContext();
        PlcRequestContainer<InternalPlcRequest, InternalPlcResponse> plcRequestContainer = requestsMap.get(senderContext);
        if (plcRequestContainer == null) {
            ctx.channel().pipeline().fireExceptionCaught(
                new PlcProtocolException("Unrelated payload received for message " + msg));
            return;
        }

        if (!(plcRequestContainer.getRequest() instanceof PlcReadRequest)) {
            ctx.fireExceptionCaught(new PlcProtocolException("Expecting a PlcReadRequest here."));
            return;
        }
        InternalPlcReadRequest request = (InternalPlcReadRequest) plcRequestContainer.getRequest();
        PlcResponseCode responseCode;
        if (msg.getStatus() != EnipStatus.EIP_SUCCESS) {
            responseCode = PlcResponseCode.NOT_FOUND;
        } else {
            responseCode = PlcResponseCode.OK;
        }

        SendRRData sendRRDataCommand = (SendRRData) msg.getCommand();
        if (sendRRDataCommand == null) {
            ctx.fireExceptionCaught(new PlcProtocolException("Expecting a SendRRData command here."));
            return;
        }
        CpfItem[] items = sendRRDataCommand.getPacket().getItems();
        if (items.length != 2) {
            ctx.fireExceptionCaught(new PlcProtocolException("Expecting 2 items here."));
            return;
        }
        CpfItem payload = items[1];
        if (!(payload instanceof UnconnectedDataItemResponse)) {
            ctx.fireExceptionCaught(new PlcProtocolException("Item[1] should be of type UnconnectedDataItemResponse"));
            return;
        }
        UnconnectedDataItemResponse enipResponse = (UnconnectedDataItemResponse) payload;
        ByteBuf data = enipResponse.getData();
        if (data.readableBytes() > 0) {
            Map<String, Pair<PlcResponseCode, BaseDefaultFieldItem>> values = new HashMap<>();
            // TODO: This is not quite correct as this will probalby not work when requesting more than one item.
            for (String fieldName : request.getFieldNames()) {
                MessageRouterResponse cipResponse = MessageRouterResponse.decode(data);
                short value;
                // TODO: This is not quite correct as we assume everything is an integer.
                if (cipResponse.getData().readableBytes() >= 2) {
                    value = cipResponse.getData().readShort();
                } else {
                    value = -1;
                }
                DefaultLongFieldItem fieldItem = new DefaultLongFieldItem((long) value);
                values.put(fieldName, new ImmutablePair<>(responseCode, fieldItem));
            }
            InternalPlcReadResponse response = new DefaultPlcReadResponse(request, values);
            plcRequestContainer.getResponseFuture().complete(response);
        }
    }

    ////////////////////////////////////////////////////////////////////////////////
    // Encoding helpers.
    ////////////////////////////////////////////////////////////////////////////////

    ////////////////////////////////////////////////////////////////////////////////
    // Decoding helpers.
    ////////////////////////////////////////////////////////////////////////////////

}
