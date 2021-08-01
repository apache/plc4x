/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using org.apache.plc4net.api;
using org.apache.plc4net.api.authentication;
using org.apache.plc4net.exceptions;

namespace org.apache.plc4net
{
    /// <summary>
    /// Manages connections to PLCs
    /// </summary>
    public class PlcDriverManager
    {
        /// <summary>
        /// Singleton instance of the manager
        /// </summary>
        private static PlcDriverManager _instance;

        /// <summary>
        /// Get the singleton instance
        /// </summary>
        public static PlcDriverManager Instance => _instance ?? (_instance = new PlcDriverManager());

        /// <summary>
        /// Dictionary for the drivers
        /// </summary>
        private readonly Dictionary<string, IPlcDriver> _drivers;

        /// <summary>
        /// Private constructor for the singleton driver manager.
        /// </summary>
        private PlcDriverManager()
        {
            _drivers = new Dictionary<string, IPlcDriver>();

            /*
             * TODO: Implement some mechanism to provide drivers -> MEF?
             */
        }

        /// <summary>
        /// Get the connection to the a PLC identified by the URL
        /// </summary>
        /// <param name="url">URL including the schema to connect to the PLC</param>
        /// <param name="authentication">Authentication to use</param>
        /// <returns>Created PLC connection</returns>
        public async Task<IPlcConnection> GetConnection(string url, IPlcAuthentication authentication)
        {
            var plcDriver = GetDriver(url);
            var connection = await plcDriver.ConnectAsync(url, authentication);
            
            //TODO: Does the driver method already connect or is a separate connect needed?
            //TODO: Should we do it like this?
            if (!connection.IsConnected)
            {
                await connection.ConnectAsync();
            }

            return connection;
        }

        public IPlcDriver GetDriver(string url)
        {
            try
            {
                Uri plcUri = new Uri(url);
                var proto = plcUri.Scheme;

                _drivers.TryGetValue(proto, out var plcDriver);

                if (plcDriver == null)
                {
                    throw new PlcConnectionException($"Unknown driver for protocol '{proto}'");
                }

                return plcDriver;
            }
            catch (UriFormatException invalidUriException)
            {
                throw new PlcConnectionException($"Provided connection string '{url}' is invalid", invalidUriException);
            }
        }
    }
}