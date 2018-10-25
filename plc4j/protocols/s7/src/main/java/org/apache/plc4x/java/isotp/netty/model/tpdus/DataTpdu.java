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
package org.apache.plc4x.java.isotp.netty.model.tpdus;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.base.messages.PlcProtocolMessage;
import org.apache.plc4x.java.isotp.netty.model.params.Parameter;
import org.apache.plc4x.java.isotp.netty.model.types.TpduCode;

import java.util.List;

public class DataTpdu extends Tpdu {

    private final boolean eot;
    private final byte tpduRef;

    public DataTpdu(boolean eot, byte tpduRef, List<Parameter> parameters, ByteBuf userData) {
        this(eot, tpduRef, parameters, userData, null);
    }

    public DataTpdu(boolean eot, byte tpduRef, List<Parameter> parameters, ByteBuf userData, PlcProtocolMessage parent) {
        super(TpduCode.DATA, parameters, userData, parent);
        this.eot = eot;
        this.tpduRef = tpduRef;
    }

    public boolean isEot() {
        return eot;
    }

    public byte getTpduRef() {
        return tpduRef;
    }

}
