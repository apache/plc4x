//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

namespace org.apache.plc4net.drivers.knxnetip.readwrite.model
{

    public class CEMIMPropReadCon : CEMI
    {
        // Accessors for discriminator values.
        public override byte GetMessageCode() {
            return 0xFB;
        }

        // Properties.
        public ushort InterfaceObjectType { get; }
        public byte ObjectInstance { get; }
        public byte PropertyId { get; }
        public byte NumberOfElements { get; }
        public ushort StartIndex { get; }
        public ushort Unknown { get; }

        public CEMIMPropReadCon(ushort interfaceObjectType, byte objectInstance, byte propertyId, byte numberOfElements, ushort startIndex, ushort unknown)
        {
            InterfaceObjectType = interfaceObjectType;
            ObjectInstance = objectInstance;
            PropertyId = propertyId;
            NumberOfElements = numberOfElements;
            StartIndex = startIndex;
            Unknown = unknown;
        }

    }

}
