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

package org.apache.plc4x.java.issues;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.s7.model.S7Field;
import org.apache.plc4x.java.s7.netty.model.types.MemoryArea;
import org.apache.plc4x.java.s7.netty.model.types.TransportSize;
import org.junit.Assert;
import org.junit.jupiter.api.Test;

import static org.hamcrest.core.IsEqual.equalTo;
import static org.junit.Assert.assertThat;

public class PLC4X56 {

    @Test
    void name() {
        S7Field field = S7Field.of("%DB56.DBB100:SINT[25]");
        assertThat(field.getMemoryArea(), equalTo(MemoryArea.DATA_BLOCKS));
        assertThat(field.getBlockNumber(), equalTo(56));
        assertThat(field.getByteOffset(), equalTo(100));
        assertThat(field.getBitOffset(), equalTo((short) 0));
        assertThat(field.getDataType(), equalTo(TransportSize.SINT));
        assertThat(field.getNumElements(), equalTo(25));
    }

    @Test
    public void invalidBlockLengthThrowsException() throws Exception {
        try {
            S7Field field = S7Field.of("%DB2000.DBB8000000:SINT[25]");
            Assert.fail();
        }
        catch (PlcInvalidFieldException e){
            e.printStackTrace();
            // exception thrown --> Test Ok
        }
    }

    @Test
    public void invalidBlockNumber1ThrowsException() throws Exception {
        try {
            S7Field field = S7Field.of("%DB0.DBB800:SINT[25]");
            Assert.fail();
        }
        catch (PlcInvalidFieldException e){
            e.printStackTrace();
            // exception thrown --> Test Ok
        }
    }

    @Test
    public void invalidBlockNumber2ThrowsException() throws Exception {
        try {
            S7Field field = S7Field.of("%DB80000.DBB800:SINT[25]");
            Assert.fail();
        }
        catch (PlcInvalidFieldException e){
            e.printStackTrace();
            // exception thrown --> Test Ok
        }
    }

}
