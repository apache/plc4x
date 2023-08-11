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
package org.apache.plc4x.merlot.das.base.core;

import org.apache.plc4x.merlot.das.api.DriverCallback;
import org.apache.plc4x.merlot.das.api.DriverEvent;
import org.apache.plc4x.merlot.das.core.Merlot.FUNCTION;
import com.lmax.disruptor.EventFactory;
import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;


public class BaseDriverEvent implements DriverEvent {

    long Sequence = 0;
    int TransactionID = 0;
    short ProtocolID = 0;
    int Field = 0;
    int LengthField = 0;
    byte UnitID = 0;
    FUNCTION FunctionCode;
    ByteBuf Data = null;	
    DriverCallback cb = null;
    
    
    public static final EventFactory<DriverEvent> FACTORY = new EventFactory<DriverEvent>() {
        @Override
        public DriverEvent newInstance()
        {
            return new BaseDriverEvent();
        }
    };    
    
    @Override
    public byte[] getSerialized() {
            return Data.array();
    }
    
    @Override
    public void setCallback(DriverCallback cb) {
            this.cb = cb;
    }
    
    @Override
    public DriverCallback getCallback() {
            return this.cb;
    }
     
    @Override
    public void setSequence(long Sequence) {
            this.Sequence = Sequence;
    }

    @Override
    public long getSequence() {
            return this.Sequence;
    }
    
    @Override
    public void setTransactionID(int TransactionID) {
            this.TransactionID = TransactionID;
    }
    
    @Override
    public int getTransactionID() {
            return this.TransactionID;
    }
    
    @Override
    public void setProtocolID(short ProtocolID) {
            this.ProtocolID = ProtocolID;
    }
    
    @Override
    public short getProtocolID() {
            return this.ProtocolID;
    }

    @Override
    public void setField(int Field) {
        this.Field = Field;
    }

    @Override
    public int getField() {
        return this.Field;
    }
       
    @Override
    public void setLengthField(int LengthField) {
            this.LengthField = LengthField;
    }
    
    @Override
    public int getLengthField() {
            return this.LengthField;
    }
    
    @Override
    public void setUnitID(byte UnitID) {
            this.UnitID = UnitID;
    }
    
    @Override
    public byte getUnitID() {
            return this.UnitID;
    }
    
    @Override
    public void setFunctionCode(FUNCTION FunctionCode) {
            this.FunctionCode = FunctionCode;
    }
    
    @Override
    public FUNCTION getFunctionCode() {
            return this.FunctionCode;
    }
    
    @Override
    public void setData(ByteBuf Data) {
            this.Data = Data;
    }
    
    @Override
    public ByteBuf getData() {
            return this.Data;
    }

    @Override
    public void setPlcReadRequest(PlcReadRequest plcReadRequest) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public PlcReadRequest getPlcReadRequest() {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public void setPlcReadResponse(PlcReadResponse plcReadResponse) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public PlcReadResponse getPlcReadResponse() {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public void setPlcWriteRequest(PlcWriteRequest plcWriteRequest) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public PlcWriteRequest getPlcWriteRequest() {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public void setPlcWriteResponse(PlcWriteResponse plcWriteResponse) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public PlcWriteResponse getPlcWriteResponse() {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    
    @Override
    public String toString() {
            return "MerlotBasicDriverEvent [Sequence=" + Sequence + ", TransactionID=" + TransactionID + ", ProtocolID="
                            + ProtocolID + ", LengthField=" + LengthField + ", UnitID=" + UnitID + ", FunctionCode=" + FunctionCode
                            + ", Data=" + Data + ", cb=" + cb + "]";
    }        
    
}
