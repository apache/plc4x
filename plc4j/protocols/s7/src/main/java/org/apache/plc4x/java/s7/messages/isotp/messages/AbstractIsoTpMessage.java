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

import org.apache.plc4x.java.exception.PlcException;
import org.apache.plc4x.java.exception.PlcIoException;
import org.apache.plc4x.java.s7.messages.Message;
import org.apache.plc4x.java.s7.messages.isotp.params.IsoTpParameter;
import org.apache.plc4x.java.s7.messages.isotp.types.TpduCode;

import java.io.DataOutputStream;
import java.io.IOException;
import java.util.List;

public abstract class AbstractIsoTpMessage implements Message {

    public static final byte ISO_ON_TCP_MAGIC_NUMBER = 0x03;

    private TpduCode tpduCode;
    private List<IsoTpParameter> isoTpParameters;

    public AbstractIsoTpMessage(TpduCode tpduCode) {
        this.tpduCode = tpduCode;
    }

    public AbstractIsoTpMessage(TpduCode tpduCode, List<IsoTpParameter> isoTpParameters) {
        this.tpduCode = tpduCode;
        this.isoTpParameters = isoTpParameters;
    }

    public TpduCode getTpduCode() {
        return tpduCode;
    }

    public List<IsoTpParameter> getIsoTpParameters() {
        return isoTpParameters;
    }

    @Override
    public void serialize(DataOutputStream dos) throws PlcException {
        try {
            ////////////////////////////////////////////////////
            // RFC 1006 (ISO on TCP)
            dos.writeByte(ISO_ON_TCP_MAGIC_NUMBER);     // Version (is always constant 0x03)
            dos.writeByte((byte) 0x00);                 // Reserved (is always constant 0x00)
            dos.writeShort(getIsoOnTcpPacketLength());  // Packet length (including ISOonTCP header)

            ////////////////////////////////////////////////////
            // RFC 905 (ISO Transport Protocol)
            dos.writeByte((byte) (getIsoTpHeaderLength() - 1)); // Length indicator field (The length byte doesn't count)
            dos.writeByte(tpduCode.getCode());  // TPDU Code (First 4 bits), Initial Credit Allocation (Second 4 bits)
            //                                      (Length of the header excluding the length indicator itself)
            serializeIsoTpHeaderFixedPart(dos);
            serializeIsoTpHeaderParameters(dos);

            serializeIsoTpUserData(dos);
        } catch (IOException e) {
            throw new PlcIoException("Error serializing message", e);
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////
    // ISO on TCP related methods

    private short getIsoOnTcpHeaderLength() {
        return 4;
    }

    protected short getIsoOnTcpPacketLength() {
        return (short) (getIsoOnTcpHeaderLength() + getIsoTpHeaderLength() + getIsoTpUserDataLength());
    }

    ////////////////////////////////////////////////////////////////////////////////////
    // ISO Transport Protocol related methods

    private short getIsoTpHeaderLength() {
        short l = (short) (2 + getIsoTpHeaderFixedPartLength());
        if(isoTpParameters != null) {
            for(IsoTpParameter isoTpParameter : isoTpParameters) {
                l += isoTpParameter.getLength();
            }
        }
        return l;
    }

    protected abstract short getIsoTpHeaderFixedPartLength();

    protected abstract void serializeIsoTpHeaderFixedPart(DataOutputStream dos) throws PlcException;

    private void serializeIsoTpHeaderParameters(DataOutputStream dos) throws PlcException {
        if(isoTpParameters != null) {
            for (IsoTpParameter isoTpParameter : isoTpParameters) {
                isoTpParameter.serialize(dos);
            }
        }
    }

    protected short getIsoTpUserDataLength() {
        return (short) 0x0000;
    }

    protected void serializeIsoTpUserData(DataOutputStream dos) throws PlcException {
        // Ignore ...
    }

}
