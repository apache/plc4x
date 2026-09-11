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

namespace org.apache.plc4net
{
    /// <summary>
    /// Extension methods for registering PLC4Net drivers via dependency injection.
    /// Assembly-scanning helpers live in the SPI layer (see
    /// <c>org.apache.plc4net.spi.drivers.Plc4NetScanExtensions</c>).
    /// </summary>
    public static class Plc4NetServiceExtensions
    {
        /// <summary>
        /// Registers a driver with the global <see cref="PlcDriverManager"/>.
        /// </summary>
        public static PlcDriverManager RegisterDriver(this PlcDriverManager manager, api.IPlcDriver driver)
        {
            if (manager == null) throw new ArgumentNullException(nameof(manager));
            if (driver == null) throw new ArgumentNullException(nameof(driver));
            manager.RegisterDriver(driver);
            return manager;
        }
    }
}
