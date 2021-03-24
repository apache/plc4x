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
using System;
using NLog;
using org.apache.plc4net.api;

namespace org.apache.plc4net.examples.helloplc4x
{
    class HelloPlc4x
    {
        private static Logger Log = LogManager.GetCurrentClassLogger();
        /// <summary>
        /// A example program of PLC4x in .NET!
        /// </summary>
        /// <param name="args"></param>
        static void Main(string[] args) 
        {
            using (IPlcConnection PlcConnection = (new PlcDriverManager()).GetConnection("s7://192.168.167.210/1/1"))
            {
                PlcConnection.ReadRequestBuilder.AddItem("001", "%MX1.0:BOOL");
                PlcConnection.ReadRequestBuilder.AddItem("002", "%MB208:CHAR");
                
                PlcConnection.ReadRequestBuilder.Build().ExecuteAsync();
                Log.Info("Request has been sent");
                PlcConnection.CloseAsync();
            }

            using (IPlcConnection PlcConnection = (new PlcDriverManager()).GetConnection("s7://192.168.167.210/1/1"))
            {
                PlcConnection.WriteRequestBuilder.AddItem("001", "%MX1.0:BOOL", true);
                PlcConnection.WriteRequestBuilder.AddItem("002", "%MB3:BYTE", (byte)0x12);

                char?[] characterArray = new char?[] { 'A', 'B', 'C', 'D' };

                PlcConnection.WriteRequestBuilder.AddItem("045", "%DB109.DBB0:CHAR[7]", characterArray);
                PlcConnection.WriteRequestBuilder.Build().ExecuteAsync();
                Log.Info("Writing data");
                PlcConnection.CloseAsync();
            }
        }
    }
}
