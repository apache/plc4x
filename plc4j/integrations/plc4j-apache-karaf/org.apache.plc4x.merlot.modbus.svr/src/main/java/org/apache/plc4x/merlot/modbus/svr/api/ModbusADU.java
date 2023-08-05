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
package org.apache.plc4x.merlot.modbus.svr.api;

import io.netty.buffer.ByteBuf;


public interface ModbusADU {

    public void setTransactionID(short TransactionID);
    
    public short getTransactionID();
    
    public void setProtocolID(short ProtocolID);
    
    public short getProtocolID();
    
    public void setLengthField(short LengthField);
    
    public short getLengthField();
    
    public void setUnitID(byte UnitID);
    
    public byte getUnitID();
    
    public void setFunctionCode(byte FunctionCode);
    
    public byte getFunctionCode();
    
    public void setData(ByteBuf  Data);
    
    public ByteBuf getData();
    
    public ByteBuf  getSerialized() ;
    
    @Override
    public String toString();
        
    
}
