using System;
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
            var formatName = "F32";
            var input = StrToByteArray("0041b00000");
            IPlcValue expected = new PlcREAL(22.0f);
            
            var actual = new KnxDatapoint().Parse(new ReadBuffer(input), formatName);
            
            Assert.Equal(expected, actual);
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