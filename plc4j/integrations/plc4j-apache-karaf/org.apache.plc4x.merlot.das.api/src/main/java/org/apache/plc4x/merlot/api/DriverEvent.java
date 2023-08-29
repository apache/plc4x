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
package org.apache.plc4x.merlot.api;

import org.apache.plc4x.merlot.api.core.Merlot.FUNCTION;
import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;


public interface DriverEvent {
	
    public void setSequence(long Sequence);
    
    public long getSequence();	
	
    public void setTransactionID(int TransactionID);
    
    public int getTransactionID();
    
    public void setProtocolID(short ProtocolID);
    
    public short getProtocolID();
    
    public void setField(int Field);
    
    public int getField();
    
    public void setLengthField(int LengthField);
    
    public int getLengthField();
    
    public void setUnitID(byte UnitID);
    
    public byte getUnitID();
    
    public void setFunctionCode(FUNCTION FunctionCode);
    
    public FUNCTION getFunctionCode();
    
    public void setData(ByteBuf  Data);
    
    public ByteBuf getData();
    
    public byte[] getSerialized() ;
    
    public void setCallback(DriverCallback cb);
    
    public DriverCallback getCallback();  
    
    
    public void setPlcReadRequest(PlcReadRequest plcReadRequest);
    public PlcReadRequest getPlcReadRequest();
    
    public void setPlcReadResponse(PlcReadResponse plcReadResponse);
    public PlcReadResponse getPlcReadResponse();  
    
    public void setPlcWriteRequest(PlcWriteRequest plcWriteRequest);
    public PlcWriteRequest getPlcWriteRequest(); 
    
    public void setPlcWriteResponse(PlcWriteResponse plcWriteResponse);
    public PlcWriteResponse getPlcWriteResponse();     
          
    @Override
    public String toString();	

}
