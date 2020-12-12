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

import io.netty.bootstrap.Bootstrap;
import io.netty.channel.Channel;
import io.netty.channel.ChannelDuplexHandler;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelFutureListener;
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelHandler.Sharable;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInboundHandlerAdapter;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelOption;
import io.netty.channel.ChannelPipeline;
import io.netty.channel.ChannelPromise;
import io.netty.channel.group.ChannelGroup;
import io.netty.channel.group.DefaultChannelGroup;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.nio.NioSocketChannel;
import io.netty.handler.codec.MessageToMessageCodec;
import io.netty.util.AttributeKey;
import io.netty.util.concurrent.Future;
import io.netty.util.concurrent.GenericFutureListener;
import io.netty.util.concurrent.GlobalEventExecutor;
import java.net.InetAddress;
import java.net.InetSocketAddress;
import java.net.SocketAddress;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import org.apache.plc4x.java.base.events.ConnectedEvent;
import org.apache.plc4x.java.isoontcp.protocol.IsoOnTcpProtocol;
import org.apache.plc4x.java.isotp.protocol.IsoTPProtocol;
import org.apache.plc4x.java.isotp.protocol.events.IsoTPConnectedEvent;
import org.apache.plc4x.java.isotp.protocol.model.types.DeviceGroup;
import org.apache.plc4x.java.isotp.protocol.model.types.TpduSize;
import org.apache.plc4x.java.s7.netty.events.S7ConnectedEvent;
import org.apache.plc4x.java.s7.utils.S7TsapIdEncoder;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * S7HProxy
 * It is a proxy towards Siemens high availability systems, which are 
 * generally made up of two (2) CPUs with or without associated CPs.
 * . Here are the individual connections to each of the CP443-1, with 
 * the TCP/IP protocol enabled, to which each of the messages received 
 * from "S7Protocol" are sent.
 * . Duplicate messages (one for each CP) are expected, so they are filtered 
 * in this version at the "S7Protocol" level. 
 * . There is no determinism from which CP will receive the response message.
 * . This object is marked as "Sharable", it allows to share it with the 
 * pipes of chrack0 and chrack1. 
 * . Being "Sharable" is guaranteed to be threadsafe.
 * . For simplicity and maintaining a granular control over the operations, 
 * the use of the "ChannelGroup" is avoided, this allows informing which of the 
 * links is down.
 * . In this first version, packet filtering is not performed, which is left to 
 * the following capable ones, which keep control of the requested packets.
 * 
 *   rack0           rack1 
 * +--+--+-+... ...+--+--+-+
 * |  |  | |       |  |  | |
 * |P |C |C|       |P |C |C|
 * |S |P |P|       |S |P |P|
 * |  |U |1|       |  |U |1|
 * |  |0 | |       |  |1 | |
 * |  |  | |       |  |  | |
 * +--+--+-+... ...+--+--+-+
 * 
 * @author cgarcia
 */

@Sharable
public class S7HProxy  extends MessageToMessageCodec {
    private static final Logger logger = LoggerFactory.getLogger(S7HProxy.class);
    private static final int ISO_ON_TCP_PORT = 102;
    private final AttributeKey<Short> pduSizeKey = AttributeKey.valueOf("pduSizeKey");    
    
    private ChannelHandlerContext ctx = null;
    private final InetAddress[] haddress;
    private final int[] racks;
    private final int[] slots; 
    private Short pduSize = 0;
    
    private Channel chrack0;
    private Channel chrack1; 
    
    private Channel proxychannel = null;
    

    public S7HProxy(InetAddress[] haddress, int [] racks, int[] slots ) {
        this.haddress = haddress;
        this.racks = racks;
        this.slots = slots;
    }
          
    @Override
    public void channelActive(ChannelHandlerContext ctx) throws Exception {
        if (this.ctx != null) return;
        try {
            this.ctx = ctx;
            proxychannel = ctx.channel();
            pduSize = ctx.channel().attr(pduSizeKey).get();
            
            
            final NioEventLoopGroup workerGroup = new NioEventLoopGroup();
            CompletableFuture<Void> sessionSetupCompleteFuture0 = new CompletableFuture<>();
            CompletableFuture<Void> sessionSetupCompleteFuture1 = new CompletableFuture<>();            
            
            Bootstrap bs_rack0 = new Bootstrap();
            bs_rack0.group(workerGroup);
            bs_rack0.channel(NioSocketChannel.class);
            bs_rack0.option(ChannelOption.SO_KEEPALIVE, true);
            bs_rack0.option(ChannelOption.TCP_NODELAY, true);

            Bootstrap bs_rack1 = bs_rack0.clone(workerGroup);
            
            bs_rack0.handler(getChannelHandlerRack0(sessionSetupCompleteFuture0));
            
            bs_rack1.handler(getChannelHandlerRack1(sessionSetupCompleteFuture1));
            
            //Start the client for rack0.
            System.out.println(haddress[0].getHostAddress() +" : " + haddress[1].getHostAddress());
            final ChannelFuture f0 = bs_rack0.connect(haddress[0], ISO_ON_TCP_PORT);
            f0.addListener(new GenericFutureListener<Future<? super Void>>() {
                @Override public void operationComplete(Future<? super Void> future) throws Exception {
                    if (future.isSuccess()) {
                        proxychannel.read();
                    } else {
                        logger.info("Unable to connect Rack0.");
                    } 
                }
            });
            
            try {
                f0.await();
                f0.awaitUninterruptibly(); // jf: unsure if we need that
            } catch (Exception ex){
                ex.printStackTrace();
            }
            
            //Start the client for rack1.            
            final ChannelFuture f1  = bs_rack1.connect(haddress[1], ISO_ON_TCP_PORT);
            f1.addListener(new GenericFutureListener<Future<? super Void>>() {
                @Override public void operationComplete(Future<? super Void> future) throws Exception {
                    if (future.isSuccess()) {
                        proxychannel.read();
                    } else {
                        logger.info("Unable to connect Rack1.");
                    } 
                }
            });

            try {
                f1.await();
                f1.awaitUninterruptibly(); // jf: unsure if we need that  
            } catch (Exception ex){
                ex.printStackTrace();
            }
            
            if (!((f0.isDone() && f0.isSuccess()) || (f1.isDone() && f1.isSuccess()))) {
                logger.info("RACK0 & 1: Shutdown.");
                workerGroup.shutdownGracefully();
                //No connections
                proxychannel.close();
            }
            
            chrack0 = f0.channel();
            chrack1 = f1.channel();

        } catch (Exception ex){
            ex.printStackTrace();
        }
    }   
    
    @Override
    public void channelInactive(ChannelHandlerContext ctx) throws Exception {
        super.channelInactive(ctx); //To change body of generated methods, choose Tools | Templates.
    }  
           
    /*
     * This method intercepts user events such as "setupCommunicationRequest", 
     * which is used to initiate communication with the CPs.
    */
    @Override
    public void userEventTriggered(ChannelHandlerContext ctx, Object evt) throws Exception {
        logger.info("S7HProxy event: " + evt.getClass().getName()); 
        if ((evt instanceof IsoTPConnectedEvent) ||
            (evt instanceof S7ConnectedEvent)    ){
            super.userEventTriggered(ctx, evt);
            return;
        }
        if (chrack0.isActive())chrack0.writeAndFlush(evt).
                addListener(new ChannelFutureListener() {
                @Override
                public void operationComplete(ChannelFuture future) {
                    if (!future.isSuccess()) {
                        logger.info("RACK0: User event send.");
                    }
                }
            });
        
        if (chrack1.isActive())chrack1.writeAndFlush(evt).
                addListener(new ChannelFutureListener() {
                @Override
                public void operationComplete(ChannelFuture future) {
                    if (!future.isSuccess()) {
                        logger.info("RACK1: User event send.");
                    }
                }
            });
    }
    
    /*
    * 
    * Plc4XS7Protocol -> 
    * LongS7MessageProcessorCodec -> 
    * S7Protocol ->
    * S7HProxy -> +--> IsoTProtocol -> IsoOnTcpProtocol -> socket0
    *             +--> IsoTProtocol -> IsoOnTcpProtocol -> socket1
    */
    @Override
    protected void encode(ChannelHandlerContext ctx, Object msg, List out) throws Exception {
        logger.info("ENCODE: " + msg.getClass().getName());
        out.add(msg);
        if (ctx != this.ctx) return;
        if (chrack0.isActive())
        chrack0.writeAndFlush(msg).addListener(new ChannelFutureListener() {
                @Override
                public void operationComplete(ChannelFuture future) {
                    if (!future.isSuccess()) {
                        logger.info("RACK0: Encode error");
                    }
                }
            }); 
        
        if (chrack1.isActive())        
        chrack1.writeAndFlush(msg).addListener(new ChannelFutureListener() {
                @Override
                public void operationComplete(ChannelFuture future) {
                    if (!future.isSuccess()) {
                        logger.info("RACK1: Encode error");
                    }
                }
            });                 
    }

    /*
    * 
    * Plc4XS7Protocol <- 
    * LongS7MessageProcessorCodec <- 
    * S7Protocol <-
    * S7HProxy <- <-+-- IsoTProtocol <- IsoOnTcpProtocol <- socket0
    *             <-+-- IsoTProtocol <- IsoOnTcpProtocol <- socket1
    */
    @Override
    protected void decode(ChannelHandlerContext ctx, Object msg, List out) throws Exception {
        logger.info("DECODE: message decode.");        
        out.add(msg);
    }    
                    
    @Override
    public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) throws Exception {
        super.exceptionCaught(ctx, cause); //To change body of generated methods, choose Tools | Templates.
        logger.info(cause.getMessage());
    }

    private ChannelHandler getChannelHandlerRack0(CompletableFuture<Void> sessionSetupCompleteFuture) {
        short calledTsapId = S7TsapIdEncoder.encodeS7TsapId(DeviceGroup.PG_OR_PC, racks[0], slots[0]);
        short callingTsapId = S7TsapIdEncoder.encodeS7TsapId(DeviceGroup.OS, 0, 0);
        
        return new ChannelInitializer() {
            @Override
            protected void initChannel(Channel channel) {
                // Build the protocol stack for communicating with the s7 protocol.
                ChannelPipeline pipeline = channel.pipeline();
                pipeline.addLast(new ChannelInboundHandlerAdapter() {
                    @Override
                    public void userEventTriggered(ChannelHandlerContext ctx, Object evt) throws Exception {
                        logger.info("Evento en el primero de la tuberia 0 ..." + evt.getClass().getName());
                        if (evt instanceof ConnectedEvent) {
                            sessionSetupCompleteFuture.complete(null);
                        } else {
                            super.userEventTriggered(ctx, evt);
                        }
                    }
                });

                pipeline.addLast(new IsoOnTcpProtocol());
                pipeline.addLast(new IsoTPProtocol(callingTsapId, calledTsapId, TpduSize.valueForGivenSize(pduSize)));
                pipeline.addLast(new S7HProxyHelper(proxychannel)); 
            }
        };
    }    
    
    private ChannelHandler getChannelHandlerRack1(CompletableFuture<Void> sessionSetupCompleteFuture) {
        short calledTsapId = S7TsapIdEncoder.encodeS7TsapId(DeviceGroup.PG_OR_PC, racks[1], slots[1]);
        short callingTsapId = S7TsapIdEncoder.encodeS7TsapId(DeviceGroup.OS, 0, 0);
        
        return new ChannelInitializer() {
            @Override
            protected void initChannel(Channel channel) {
                // Build the protocol stack for communicating with the s7 protocol.
                ChannelPipeline pipeline = channel.pipeline();
                pipeline.addLast(new ChannelInboundHandlerAdapter() {
                    @Override
                    public void userEventTriggered(ChannelHandlerContext ctx, Object evt) throws Exception {
                        logger.info("Evento en el primero de la tuberia 1 ..." + evt.getClass().getName());                        
                        if (evt instanceof ConnectedEvent) {
                            sessionSetupCompleteFuture.complete(null);
                        } else {
                            super.userEventTriggered(ctx, evt);
                        }
                    }
                });

                pipeline.addLast(new IsoOnTcpProtocol());
                pipeline.addLast(new IsoTPProtocol(callingTsapId, calledTsapId, TpduSize.valueForGivenSize(pduSize)));
                pipeline.addLast(new S7HProxyHelper(proxychannel));                
            }
        };
    }       

    
}
