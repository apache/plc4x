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

    public class CEMIDataInd : CEMI
    {
        // Accessors for discriminator values.
        public override byte GetMessageCode() {
            return 0x29;
        }

        // Properties.
        public byte AdditionalInformationLength { get; }
        public CEMIAdditionalInformation[] AdditionalInformation { get; }
        public CEMIDataFrame CemiDataFrame { get; }

        public CEMIDataInd(byte additionalInformationLength, CEMIAdditionalInformation[] additionalInformation, CEMIDataFrame cemiDataFrame)
        {
            AdditionalInformationLength = additionalInformationLength;
            AdditionalInformation = additionalInformation;
            CemiDataFrame = cemiDataFrame;
        }

    }

}
