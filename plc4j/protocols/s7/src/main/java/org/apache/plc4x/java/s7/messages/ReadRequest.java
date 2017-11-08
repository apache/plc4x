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
package org.apache.plc4x.java.s7.messages;

import org.apache.plc4x.java.s7.messages.s7.messages.S7RequestMessage;
import org.apache.plc4x.java.s7.messages.s7.params.ReadVarParameter;
import org.apache.plc4x.java.s7.messages.s7.types.MemoryArea;
import org.apache.plc4x.java.s7.messages.s7.types.MessageType;
import org.apache.plc4x.java.s7.messages.s7.types.TransportSize;

import java.util.Collections;

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
public class ReadRequest extends S7RequestMessage {

    public ReadRequest(MemoryArea memoryArea, TransportSize transportSize, short numElements,
                       short dataBlockNumber, short byteOffset, byte bitOffset) {
        super(MessageType.JOB, Collections.singletonList(
            new ReadVarParameter(memoryArea, transportSize, numElements, dataBlockNumber, byteOffset, bitOffset)),
            null);
    }

}
