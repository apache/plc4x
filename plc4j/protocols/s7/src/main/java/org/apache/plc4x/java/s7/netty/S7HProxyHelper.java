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

import io.netty.channel.Channel;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.embedded.EmbeddedChannel;
import io.netty.handler.codec.MessageToMessageCodec;
import java.util.List;
import org.apache.plc4x.java.base.events.ConnectEvent;
import org.apache.plc4x.java.isotp.protocol.IsoTPProtocol;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 *
 * @author cgarcia
 */
public class S7HProxyHelper  extends MessageToMessageCodec{
    private static final Logger logger = LoggerFactory.getLogger(S7HProxyHelper.class);
    private final EmbeddedChannel proxychannel;
    
    public S7HProxyHelper(Channel proxychannel) {
        this.proxychannel = (EmbeddedChannel) proxychannel;
    }
    
    /*
    *  
    */
    @Override
    public void userEventTriggered(ChannelHandlerContext ctx, Object evt) throws Exception {
        logger.info("S7HProxyHelper event: userEventTriggered: " + evt.getClass().getName());
        proxychannel.pipeline().fireUserEventTriggered(evt);
    }          
    /*
    * TODO: Chequear la propa
    */
    @Override
    protected void encode(ChannelHandlerContext ctx, Object obj, List out) throws Exception {
        logger.info("S7HProxyHelper: Encode en la cola..." + obj.getClass().getName());
        if (obj instanceof ConnectEvent){
            logger.info("002: Recibe objeto y es trasladado al HEAD.");
            ctx.pipeline().fireUserEventTriggered(obj);
        } else {
            out.add(obj);
        }
    }

    /*
    * 
    */
    @Override
    protected void decode(ChannelHandlerContext ctx, Object obj, List out) throws Exception {
        logger.info("S7HProxyHelper: decode " + obj.getClass().getName());
        proxychannel.writeInbound(obj);
    }
    
}
