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
package org.apache.plc4x.java.s7.readwrite.protocol;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.ByteBufUtil;
import io.netty.channel.Channel;
import io.netty.channel.ChannelHandler.Sharable;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.embedded.EmbeddedChannel;

import io.netty.handler.codec.MessageToMessageCodec;
import io.netty.handler.timeout.ReadTimeoutHandler;
import io.netty.util.AttributeKey;
import java.time.LocalTime;
import java.util.List;
import org.apache.plc4x.java.spi.events.ConnectEvent;
import org.apache.plc4x.java.spi.events.DisconnectEvent;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Implementation of a multiplexing channel, from an embedded channel to two 
 * possible TCP connections, primary and secondary.
 * The objective is to allow connections to individual systems 
 * with a two CP (PN CPUs, CP343-1, CP443-1 or similar), or H-type systems 
 * (S7-400H or S7-1500H).
 * 
 * The user App must be in charge of restoring the requests or 
 * subscriptions that it is requesting.
 * 
 * @author cgarcia
 */
@Sharable
public class S7HMuxImpl extends MessageToMessageCodec<ByteBuf, ByteBuf> implements S7HMux {

    private static final Logger logger = LoggerFactory.getLogger(S7HMuxImpl.class);
    
    /*
    * This attribute indicates to the other handlers that the channel is connected
    * or disconnected because a switch is being made between TCP channels or 
    * both TCP channels are disconnected.
    * Default value: false
    */
    final static AttributeKey<Boolean> IS_CONNECTED = AttributeKey.valueOf("IS_CONNECTED");
    
    /*
    * This attribute indicates to the other handlers which channel is being used,
    * this allows the request to be properly prepared.
    * For example, in the case of a CPU with two CPs, you should change 
    * the "slot", in the case of H systems, you should change the "rack",
    * the correct values will be defined in the connection URL.
    * Default value: true
    */
    final static AttributeKey<Boolean> IS_PRIMARY = AttributeKey.valueOf("IS_PRIMARY");
    
    /*
    * This is the maximum waiting time for reading on the TCP channel. 
    * As there is no traffic, it must be assumed that the connection with the 
    * interlocutor was lost and it must be restarted.
    * When the channel is closed, the "fail over" is carried out 
    * in case of having the secondary channel, or it is expected that it 
    * will be restored automatically, which is done every 4 seconds.
    * Default value: 8 sec.
    */
    final static AttributeKey<Integer> READ_TIME_OUT = AttributeKey.valueOf("READ_TIME_OUT");    
        
    /*
    * If your application requires sampling times greater than the 
    * set "watchdog" time, it is important that the PING option is activated, 
    * this will prevent the TCP channel from being closed unnecessarily. 
    * Default value: false
    */
    final static AttributeKey<Boolean> IS_PING_ACTIVE = AttributeKey.valueOf("IS_PIN_ACTIVE");
    
    /*
    * Time value in seconds at which the execution of the PING will be scheduled.
    * Generally set by developer experience, but generally should be the same 
    * as READ_TIME_OUT / 2.
    * Default value: -1    
    */
    final static AttributeKey<Integer> PING_TIME = AttributeKey.valueOf("PING_TIME");      
    
    /*
    * Time for supervision of TCP channels. If the channel is not active,
    * a safe stop of the EventLoop must be performed, to ensure that 
    * no additional tasks are created.
    * Default value: 4    
    */
    final static AttributeKey<Integer> RETRY_TIME = AttributeKey.valueOf("RETRY_TIME");     
    
    ChannelHandlerContext embed_ctx = null;
    protected Channel embeded_channel = null;    
    protected Channel tcp_channel = null;
    protected Channel primary_channel = null;     
    protected Channel secondary_channel = null;    

    /*
    * From S7ProtocolLogic
    * TODO: Evaluate if the "embed_ctx" is really required since we set 
    * the Embeded channel when we created it.
    */
    @Override
    protected void encode(ChannelHandlerContext ctx, ByteBuf outbb, List<Object> list) throws Exception { 
        //System.out.println("Llego aqui: " + ByteBufUtil.hexDump(outbb));
        if ((embed_ctx == null) && (ctx.channel() instanceof EmbeddedChannel)) embed_ctx = ctx;
        if ((tcp_channel != null)  && (embed_ctx == ctx)){  
            tcp_channel.writeAndFlush(outbb.copy());
        } else {
            list.add(outbb.copy());                        
        }
    }

    /*
    * To S7ProtocolLogic
    * The information received here from the channel "tcp_channel" is sent to 
    * the pipeline of the channel "embeded_channel"
    */
    @Override
    protected void decode(ChannelHandlerContext ctx, ByteBuf inbb, List<Object> list) throws Exception {    
        embed_ctx.fireChannelRead(inbb.copy());
    }

    @Override
    public void channelRegistered(ChannelHandlerContext ctx) throws Exception {
        super.channelRegistered(ctx);
        logger.debug("channelRegistered: " + ctx.name());         
    }
    
    @Override
    public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) throws Exception {
        super.exceptionCaught(ctx, cause);
        logger.debug("exceptionCaught: " + ctx.name());         
    }

    @Override
    public void channelWritabilityChanged(ChannelHandlerContext ctx) throws Exception {
        super.channelWritabilityChanged(ctx); 
        logger.debug("channelWritabilityChanged: " + ctx.name() );        
    }

    /*
    * The events detected here flow from the S7ProtocolLogic object.
    * Upon receiving the "ConnectEvent" event, we must safely add the watchdog 
    * to the pipeline of the "tcp_channel" connection.
    * The supervision time can be defined in the connection URL, 
    * the default value being 8 secs, this value being defined experimentally.
    */
    @Override
    public void userEventTriggered(ChannelHandlerContext ctx, Object evt) throws Exception {
        super.userEventTriggered(ctx, evt); 
        logger.info(LocalTime.now().toString() + " userEventTriggered: " + ctx.name() + " Event: "  + evt);         
        if (evt instanceof ConnectEvent) {
            try {
                tcp_channel.pipeline().remove("watchdog");                 
            } catch (Exception ex){
                logger.info(ex.toString());
            }
            try {
                tcp_channel.pipeline().addFirst("watchdog", new ReadTimeoutHandler(30)); 
                if (tcp_channel.isActive()) {
                    embeded_channel.attr(IS_CONNECTED).set(true);
                 } else {
                    embeded_channel.attr(IS_CONNECTED).set(false);
                }               
            } catch (Exception ex){
                logger.info(ex.toString());                
            }            
        }
        
        if (evt instanceof DisconnectEvent) {
            logger.info("DisconnectEvent");
        }        
        
    }

    @Override
    public void channelReadComplete(ChannelHandlerContext ctx) throws Exception {
        super.channelReadComplete(ctx); 
        logger.debug(LocalTime.now().toString() + " channelReadComplete: " + ctx.name() );
    }

    @Override
    public void channelInactive(ChannelHandlerContext ctx) throws Exception {
        super.channelInactive(ctx); 
        logger.debug("channelInactive: " + ctx.name());        
    }

    @Override
    public void channelActive(ChannelHandlerContext ctx) throws Exception {
        super.channelActive(ctx); 
        logger.debug("channelActive: " + ctx.name());          
    }

    @Override
    public void channelUnregistered(ChannelHandlerContext ctx) throws Exception {
        super.channelUnregistered(ctx); 
        logger.debug(LocalTime.now().toString() + " channelUnregistered: " + ctx.name() );        
        String strCanal = (tcp_channel == primary_channel)?"PRIMARY":"SECONDARY";
        logger.debug("Unregistered of channel: " + strCanal);
        //TODO: If embedded channel is closed, we need close all channels
        if (ctx == embed_ctx) return;        
        
        if (tcp_channel == ctx.channel())
        embeded_channel.attr(IS_CONNECTED).set(false);
        
        logger.info(embed_ctx.executor().toString());
                
        if ((tcp_channel == primary_channel) &&
            (primary_channel == ctx.channel()))
        if ((!primary_channel.isActive()) &&
            (secondary_channel != null)) {
        if (secondary_channel.isActive())
            synchronized(tcp_channel) {
                logger.info("Using secondary TCP channel.");
                tcp_channel =  secondary_channel;
                embeded_channel.attr(IS_PRIMARY).set(false);
                embeded_channel.pipeline().fireUserEventTriggered(new ConnectEvent());   
            }
        };


        if ((tcp_channel == secondary_channel) &&
            (secondary_channel == ctx.channel()))
        if ((!secondary_channel.isActive() &&
            (primary_channel.isActive()))) {
            synchronized(tcp_channel) {         
                logger.info("Using primary TCP channel.");               
                tcp_channel = primary_channel;
                embeded_channel.attr(IS_PRIMARY).set(true);                
                embeded_channel.pipeline().fireUserEventTriggered(new ConnectEvent());             
            }
        } 
    }


    @Override
    public void setEmbededhannel(Channel embeded_channel) {
        this.embeded_channel = embeded_channel;
        this.embeded_channel.attr(IS_CONNECTED).set(false);  
        this.embeded_channel.attr(IS_PRIMARY).set(true);         
        this.embeded_channel.attr(READ_TIME_OUT).set(8);
        this.embeded_channel.attr(IS_PING_ACTIVE).set(false); 
        this.embeded_channel.attr(PING_TIME).set(-1); 
        this.embeded_channel.attr(RETRY_TIME).set(8); 
    }        
    
    public void setPrimaryChannel(Channel primary_channel) {       
        if ((this.primary_channel == null) && (tcp_channel == null)){
            if (primary_channel != null){    
                this.primary_channel = primary_channel;
                tcp_channel = primary_channel; 
                embeded_channel.attr(IS_PRIMARY).set(true);                 
            }
        } else if ((!this.primary_channel.isActive()) && (tcp_channel == secondary_channel)){
            this.primary_channel = primary_channel;
        }  else if ((!this.primary_channel.isActive()) && (tcp_channel == this.primary_channel)){
            synchronized(tcp_channel) {
                tcp_channel.close();
                this.primary_channel = primary_channel;
                tcp_channel =  primary_channel;
                embeded_channel.attr(IS_PRIMARY).set(true);                
                if (tcp_channel.isActive()) {
                    embed_ctx.fireUserEventTriggered(new ConnectEvent());                               
                }
            }
        }        
    }

    @Override
    public void setSecondaryChannel(Channel secondary_channel) {
        if ((this.primary_channel == null) && (tcp_channel == null)) {
            if (secondary_channel != null){
                this.secondary_channel = secondary_channel;
                tcp_channel =  secondary_channel;    
                embeded_channel.attr(IS_PRIMARY).set(false);                
            }
        } else if ((this.secondary_channel == null) || (tcp_channel == primary_channel)){
            this.secondary_channel = secondary_channel;
        } else if ((!this.secondary_channel.isActive()) && (tcp_channel == primary_channel)){
            this.secondary_channel = secondary_channel;
        }  else if ((!this.secondary_channel.isActive()) && (tcp_channel == this.secondary_channel)){
            synchronized(tcp_channel) {
                tcp_channel.close();                
                this.secondary_channel = secondary_channel;
                tcp_channel =  secondary_channel;
                embeded_channel.attr(IS_PRIMARY).set(false);                 
            }
            if (tcp_channel.isActive()) {
                embed_ctx.fireUserEventTriggered(new ConnectEvent());                               
            }                 
        }
    }
             
    
}
