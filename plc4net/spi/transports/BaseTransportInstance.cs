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

using System;

namespace org.apache.plc4net.spi.transports
{
    /// <summary>
    /// Shared state for transport implementations: the configuration and the
    /// driver-config carved off the connection URL.
    /// </summary>
    public abstract class BaseTransportInstance : ITransportInstance
    {
        private string _driverConfig = string.Empty;

        protected BaseTransportInstance(ITransportConfiguration configuration)
        {
            Configuration = configuration ?? throw new ArgumentNullException(nameof(configuration));
        }

        public ITransportConfiguration Configuration { get; }

        public string DriverConfig
        {
            get => _driverConfig;
            protected set => _driverConfig = value ?? string.Empty;
        }

        /// <summary>
        /// Records the portion of the address the transport did not consume. Transports
        /// call this immediately after construction.
        /// </summary>
        public void SetDriverConfig(string driverConfig)
        {
            DriverConfig = driverConfig;
        }

        public abstract bool IsOpen { get; }

        public abstract int GetNumBytesAvailable();

        public abstract byte[] PeekReadableBytes(int numBytes);

        public abstract byte[] Read(int numBytes);

        public abstract void Write(byte[] bytes);

        public abstract void Close();

        protected void EnsureOpen()
        {
            if (!IsOpen)
            {
                throw new TransportException("Transport is closed");
            }
        }

        public void Dispose()
        {
            Dispose(true);
            GC.SuppressFinalize(this);
        }

        protected virtual void Dispose(bool disposing)
        {
            if (!disposing)
            {
                return;
            }
            try
            {
                Close();
            }
            catch (TransportException)
            {
                // Dispose must not throw; a failure to close is already terminal.
            }
        }
    }
}
