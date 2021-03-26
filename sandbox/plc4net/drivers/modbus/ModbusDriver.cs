using System;
using System.Collections.Generic;
using System.Net.Sockets;
using org.apache.plc4net.exceptions;

namespace org.apache.plc4net.drivers.modbus
{
    public class ModbusDriver
    {
        public void ModbusTCP()
        {
            TcpClient tcpClient = new TcpClient();
            var connectionResult = tcpClient.BeginConnect("localhost", 502, null, null);
            var isSuccessful = connectionResult.AsyncWaitHandle.WaitOne(1000);
            
            if(!isSuccessful)
            {
                throw new PlcConnectionException();
            }

            tcpClient.EndConnect(connectionResult);
        }
    }
}
