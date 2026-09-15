//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

using System.Collections.Generic;
using org.apache.plc4net.api.value;
using org.apache.plc4net.drivers.knxnetip.readwrite.model;
using org.apache.plc4net.spi.generation;
using org.apache.plc4net.spi.model.values;
using Xunit;

namespace org.apache.plc4net.test.knxnetip.readwrite.model
{
    public class KnxDatapointTests
    {
        [Fact]
        public void DecodeKnxData1()
        {
            // DPT 14.019 is a reserved byte followed by an IEEE-754 float32, so the
            // frame is five bytes. The test previously passed only the four float
            // bytes; it never showed up because the suite was not being executed
            // (the test project had no Microsoft.NET.Test.Sdk reference).
            var input = StrToByteArray("0041b00000");
            IPlcValue expected = new PlcREAL(22.0f);

            var actual = KnxDatapoint.StaticParse(new ReadBuffer(input), KnxDatapointType.DPT_Value_Electric_Current);

            Assert.Equal(expected, actual);
            Assert.Equal(22.0f, actual.GetFloat());
        }
        
        private static byte[] StrToByteArray(string str)
        {
            var hexIndex = new Dictionary<string, byte>();
            for (var i = 0; i <= 255; i++) 
            {
                hexIndex.Add(i.ToString("X2"), (byte) i);
            }

            var hexRes = new List<byte>();
            for (var i = 0; i < str.Length; i += 2)
            {
                hexRes.Add(hexIndex[str.Substring(i, 2).ToUpper()]);
            }
            return hexRes.ToArray();
        }
    }

}