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
package org.apache.plc4x.merlot.modbus.svr.core;

import org.apache.plc4x.merlot.modbus.svr.api.ModbusADU;
import org.apache.plc4x.merlot.modbus.svr.impl.ModbusADUImpl;
import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.ByteToMessageDecoder;
import java.util.List;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class ModbusServerADUDecoder extends ByteToMessageDecoder {
    private final Logger LOGGER = LoggerFactory.getLogger(ModbusServerADUDecoder.class.getName()); 
    
    @Override
    protected void decode(ChannelHandlerContext ctx, ByteBuf in, List<Object> out) throws Exception {
        ModbusADU ADU = new ModbusADUImpl();
        
        try {   		
            ADU.setTransactionID(in.readShort());
            ADU.setProtocolID(in.readShort());
            ADU.setLengthField (in.readShort());
            ADU.setUnitID(in.readByte());
            ADU.setFunctionCode(in.readByte());     
            ADU.setData(in.readBytes(ADU.getLengthField()-2));
        } catch (Exception ex) {  
            LOGGER.info(ex.getMessage());
        }  
        out.add(ADU);
    }

    /* (non-Javadoc)
     * @see io.netty.channel.ChannelInboundHandlerAdapter#exceptionCaught(io.netty.channel.ChannelHandlerContext, java.lang.Throwable)
     *
     */
    @Override
    public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) throws Exception {
        LOGGER.info(cause.getMessage());        	
    }



}
