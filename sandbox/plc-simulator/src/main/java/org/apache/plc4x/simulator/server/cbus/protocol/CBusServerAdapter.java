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
package org.apache.plc4x.simulator.server.cbus.protocol;

import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInboundHandlerAdapter;
import org.apache.plc4x.java.cbus.readwrite.*;
import org.apache.plc4x.simulator.model.Context;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class CBusServerAdapter extends ChannelInboundHandlerAdapter {

    private static final Logger LOGGER = LoggerFactory.getLogger(CBusServerAdapter.class);

    private Context context;

    private static final RequestContext requestContext = new RequestContext(false, false, false);
    private static final CBusOptions cBusOptions = new CBusOptions(false, false, false, false, false, false, false, false, false);

    public CBusServerAdapter(Context context) {
        this.context = context;
    }

    @Override
    public void channelRead(ChannelHandlerContext ctx, Object msg) throws Exception {
        if (!(msg instanceof CBusMessage)) {
            return;
        }
        CBusMessage packet = (CBusMessage) msg;
        if (packet instanceof CBusMessageToClient) {
            LOGGER.info("Message to client not supported\n{}", packet);
            return;
        }
        CBusMessageToServer cBusMessageToServer = (CBusMessageToServer) packet;
        Request request = cBusMessageToServer.getRequest();
        if (request instanceof RequestEmpty || request instanceof RequestNull) {
            LOGGER.debug("Ignoring\n{}", request);
            return;
        }
        if (request instanceof RequestDirectCommandAccess) {
            RequestDirectCommandAccess requestDirectCommandAccess = (RequestDirectCommandAccess) request;
            // TODO: handle this
            return;
        }
        if (request instanceof RequestCommand) {
            RequestCommand requestCommand = (RequestCommand) request;
            // TODO: handle this
            Alpha alpha = requestCommand.getAlpha();
            if (alpha != null) {
                Confirmation confirmation = new Confirmation(alpha, null, ConfirmationType.CONFIRMATION_SUCCESSFUL);
                ReplyOrConfirmationConfirmation replyOrConfirmationConfirmation = new ReplyOrConfirmationConfirmation(alpha.getCharacter(), confirmation, null, cBusOptions, requestContext);
                CBusMessage response = new CBusMessageToClient(replyOrConfirmationConfirmation, requestContext, cBusOptions);
                ctx.writeAndFlush(response);
            }
            return;
        }
        if (request instanceof RequestObsolete) {
            RequestObsolete requestObsolete = (RequestObsolete) request;
            // TODO: handle this
            return;
        }
        if (request instanceof RequestReset) {
            // TODO: handle this
            return;
        }
        if (request instanceof RequestSmartConnectShortcut) {
            // TODO: handle this
            return;
        }
    }

}
