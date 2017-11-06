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
package org.apache.plc4x.java.s7.messages.isotp.params;

import org.apache.plc4x.java.exception.PlcException;
import org.apache.plc4x.java.exception.PlcIoException;
import org.apache.plc4x.java.s7.messages.isotp.types.DeviceGroup;
import org.apache.plc4x.java.s7.messages.isotp.types.ParameterCode;

import java.io.DataOutputStream;
import java.io.IOException;

public class CallingTsapIsoTpParameter implements IsoTpParameter {

    private DeviceGroup deviceGroup;
    private int rackNumber;
    private int slotNumber;

    public CallingTsapIsoTpParameter(DeviceGroup deviceGroup, int rackNumber, int slotNumber) {
        this.deviceGroup = deviceGroup;
        this.rackNumber = rackNumber;
        this.slotNumber = slotNumber;
    }

    public DeviceGroup getDeviceGroup() {
        return deviceGroup;
    }

    public int getRackNumber() {
        return rackNumber;
    }

    public int getSlotNumber() {
        return slotNumber;
    }

    @Override
    public void serialize(DataOutputStream dos) throws PlcException {
        try {
            dos.writeByte(ParameterCode.CALLING_TSAP.getCode());
            dos.writeByte((byte) 0x02);
            dos.writeByte(deviceGroup.getCode());
            dos.writeByte((byte) ((rackNumber << 4) | (slotNumber)));
        } catch (IOException e) {
            throw new PlcIoException("Error serializing message", e);
        }
    }

    @Override
    public short getLength() {
        return 4;
    }

}
