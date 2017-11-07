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

import org.apache.plc4x.java.exceptions.PlcException;
import org.apache.plc4x.java.exceptions.PlcIoException;
import org.apache.plc4x.java.exceptions.PlcProtocolException;
import org.apache.plc4x.java.s7.messages.isotp.types.DeviceGroup;
import org.apache.plc4x.java.s7.messages.isotp.types.ParameterCode;
import org.apache.plc4x.java.s7.messages.isotp.types.TpduSize;

import java.io.DataInputStream;
import java.io.IOException;

public class IsoTpParameterFactory {

    public static IsoTpParameter parse(DataInputStream dis) throws PlcException {
        try {
            ParameterCode parameterCode = ParameterCode.valueOf(dis.readByte());
            if(parameterCode == null) {
                throw new PlcProtocolException("Could not find parameter code");
            }
            byte tmp;
            switch (parameterCode) {
                case TPDU_SIZE:
                    dis.readByte();
                    TpduSize tpduSize = TpduSize.valueOf(dis.readByte());
                    return new PduSizeIsoTpParameter(tpduSize);
                case CALLING_TSAP:
                    dis.readByte();
                    DeviceGroup callingDeviceGroup = DeviceGroup.valueOf(dis.readByte());
                    tmp = dis.readByte();
                    return new CallingTsapIsoTpParameter(callingDeviceGroup,
                        (tmp & 0xF0) >> 4, (tmp & 0x0F));
                case CALLED_TSAP:
                    dis.readByte();
                    DeviceGroup calledDeviceGroup = DeviceGroup.valueOf(dis.readByte());
                    tmp = dis.readByte();
                    return new CalledTsapIsoTpParameter(calledDeviceGroup,
                        (byte) ((tmp & 0xF0) >> 4), (byte) (tmp & 0x0F));
                default:
                    throw new PlcProtocolException("Parameter not implemented yet " +parameterCode.name());
            }
        } catch (IOException e) {
            throw new PlcIoException("Error parsing message", e);
        }
    }

}
