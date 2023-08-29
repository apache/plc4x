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
package org.apache.plc4x.merlot.api.impl;

import org.apache.plc4x.merlot.api.DriverCallback;
import org.apache.plc4x.merlot.api.DriverEvent;
import org.apache.plc4x.merlot.api.core.Merlot.FUNCTION;
import com.lmax.disruptor.EventFactory;
import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;


public class DriverEventImpl implements DriverEvent {

    private int TransactionID = 0;
    private long Sequence = 0;
    private short ProtocolID  = 0;
    private int Field  = 0;
    private int LengthField = 0;
    private byte UnitID = 0;
    private FUNCTION FunctionCode;
    private ByteBuf Data = null;
    private DriverCallback cb = null;
    
    private PlcReadRequest plcReadRequest = null;
    private PlcReadResponse plcReadResponse = null;
    private PlcWriteRequest plcWriteRequest = null;
    private PlcWriteResponse plcWriteResponse = null;
    
    
    public static final EventFactory<DriverEvent> FACTORY = new EventFactory<DriverEvent>() {
        @Override
        public DriverEvent newInstance()
        {
            return new DriverEventImpl();
        }
    };        
	
    public DriverEventImpl(){
        
    }

/*    
    public DriverEventImpl(byte unitID, short protocolID, byte functionCode, short transactionID, ByteBuf data,
                    short lengthField, DriverCallback cb) {
        super();
        UnitID = unitID;
        ProtocolID = protocolID;
        FunctionCode = functionCode;
        TransactionID = transactionID;
        Data = data;
        LengthField = lengthField;
        this.cb = cb;
    }	
*/	
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
    public byte[] getSerialized() {
        return null;
    }

    @Override
    public void setCallback(DriverCallback cb) {
        this.cb = cb;
    }

    @Override
    public DriverCallback getCallback() {
        return cb;
    }
    
    @Override
    public void setSequence(long Sequence) {
        this.Sequence = Sequence;
    }

    @Override
    public long getSequence() {
        return Sequence;
    }

    @Override
    public void setPlcReadRequest(PlcReadRequest plcReadRequest) {
        this.plcReadRequest = plcReadRequest;
    }

    @Override
    public PlcReadRequest getPlcReadRequest() {
        return plcReadRequest;
    }

    @Override
    public void setPlcReadResponse(PlcReadResponse plcReadResponse) {
        this.plcReadResponse = plcReadResponse;
    }

    @Override
    public PlcReadResponse getPlcReadResponse() {
        return plcReadResponse;
    }

    @Override
    public void setPlcWriteRequest(PlcWriteRequest plcWriteRequest) {
        this.plcWriteRequest = plcWriteRequest;
    }

    @Override
    public PlcWriteRequest getPlcWriteRequest() {
        return plcWriteRequest;
    }

    @Override
    public void setPlcWriteResponse(PlcWriteResponse plcWriteResponse) {
        this.plcWriteResponse = plcWriteResponse;
    }

    @Override
    public PlcWriteResponse getPlcWriteResponse() {
        return plcWriteResponse;
    }
             
    @Override
    public String toString() {
        return "DriverEventImpl{" + 
                "TransactionID=" + TransactionID + 
                ", Sequence=" + Sequence + 
                ", ProtocolID=" + ProtocolID + 
                ", Field=" + Field + 
                ", LengthField=" + LengthField + 
                ", UnitID=" + UnitID + 
                ", FunctionCode=" + FunctionCode + 
                ", Data=" + Data + 
                ", cb=" + cb + 
                ", plcReadRequest=" + plcReadRequest + 
                ", plcReadResponse=" + plcReadResponse + 
                ", plcWriteRequest=" + plcWriteRequest + 
                ", plcWriteResponse=" + plcWriteResponse + '}';
    }




}
