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
package org.apache.plc4x.java.s7.messages.isotp.messages;

import org.apache.plc4x.java.exceptions.PlcException;
import org.apache.plc4x.java.exceptions.PlcIoException;
import org.apache.plc4x.java.s7.messages.isotp.params.CalledTsapIsoTpParameter;
import org.apache.plc4x.java.s7.messages.isotp.params.CallingTsapIsoTpParameter;
import org.apache.plc4x.java.s7.messages.isotp.params.PduSizeIsoTpParameter;
import org.apache.plc4x.java.s7.messages.isotp.types.DeviceGroup;
import org.apache.plc4x.java.s7.messages.isotp.types.ProtocolClass;
import org.apache.plc4x.java.s7.messages.isotp.types.TpduCode;
import org.apache.plc4x.java.s7.messages.isotp.types.TpduSize;

import java.io.DataOutputStream;
import java.io.IOException;
import java.util.Arrays;

/**
 * Related Links:
 * <ul>
 *  <li>@see <a href="http://gmiru.com/article/s7comm/">S7 Protocol</a></li>
 *  <li>@see <a href="https://tools.ietf.org/html/rfc905">ISO Transport Protocol (Class 0)</a></li>
 *  <li>@see <a href="https://tools.ietf.org/html/rfc1006">ISO on TCP</a></li>
 *  <li>@see <a href="https://www.tanindustrie.de/fr/Help/ConfigClient/tsap_s7.htm">Reference to calculating the TSAP ids</a>
 *  <li>@see <a href="https://support.industry.siemens.com/tf/ww/en/posts/classic-style-any-pounter-to-variant-type/126024/?page=0&pageSize=10">Structure and some constants of a variable read/write request</a></li>
 * </ul>
 */
public abstract class AbstractConnectionMessage extends AbstractIsoTpMessage {

    private short destinationReference;
    private short sourceReference;

    public AbstractConnectionMessage(TpduCode tpduCode) {
        super(tpduCode);
    }

    public AbstractConnectionMessage(TpduCode tpduCode, short connectionId, byte rackNumber, byte slotNumber) {
        super(tpduCode,
            Arrays.asList(
                new PduSizeIsoTpParameter(TpduSize.SIZE_1024),
                new CallingTsapIsoTpParameter(DeviceGroup.PG_OR_PC, 0, 0),
                new CalledTsapIsoTpParameter(DeviceGroup.OTHERS, rackNumber, slotNumber)));

        this.destinationReference = 0x0000;
        this.sourceReference = connectionId;
    }

    public short getDestinationReference() {
        return destinationReference;
    }

    public short getSourceReference() {
        return sourceReference;
    }

    @Override
    protected short getIsoTpHeaderFixedPartLength() {
        return 5;
    }

    @Override
    protected void serializeIsoTpHeaderFixedPart(DataOutputStream dos) throws PlcException {
        try {
            dos.writeShort(destinationReference); // Destination Reference (constantly set to 0)
            dos.writeShort(sourceReference);      // Source Reference (Connection ID set by the client to identify the connection)
            dos.writeByte(ProtocolClass.CLASS_0.getCode());
        } catch (IOException e) {
            throw new PlcIoException("Error serializing message", e);
        }
    }

}
