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
using System.Reflection;
using org.apache.plc4net.api;
using org.apache.plc4net.spi.transports;

namespace org.apache.plc4net.spi.drivers
{
    /// <summary>
    /// Assembly scanning helpers for auto-registering drivers. These live
    /// in the SPI layer (not api) so they have access to <c>ITransportManager</c>
    /// without creating a circular dependency.
    /// </summary>
    public static class Plc4NetScanExtensions
    {
        /// <summary>
        /// Scans <paramref name="assembly"/> for types extending <c>DriverBase</c>,
        /// instantiates each one with the given <paramref name="transportManager"/>,
        /// and registers the resulting driver with the global <c>PlcDriverManager</c>.
        /// </summary>
        public static int ScanAndRegisterDrivers(
            this PlcDriverManager manager,
            Assembly assembly,
            ITransportManager transportManager)
        {
            if (manager == null) throw new ArgumentNullException(nameof(manager));
            if (assembly == null) throw new ArgumentNullException(nameof(assembly));
            if (transportManager == null) throw new ArgumentNullException(nameof(transportManager));

            var count = 0;
            foreach (var type in assembly.GetExportedTypes())
            {
                if (type.IsAbstract || type.IsInterface) continue;

                if (!IsDriverBase(type)) continue;

                var ctors = type.GetConstructors(
                    BindingFlags.Public | BindingFlags.Instance);
                if (ctors.Length == 0) continue;

                foreach (var ctor in ctors)
                {
                    var parameters = ctor.GetParameters();
                    var args = new object[parameters.Length];
                    var matched = true;
                    for (int i = 0; i < parameters.Length; i++)
                    {
                        if (typeof(ITransportManager)
                            .IsAssignableFrom(parameters[i].ParameterType))
                        {
                            args[i] = transportManager;
                        }
                        else if (parameters[i].HasDefaultValue)
                        {
                            args[i] = parameters[i].DefaultValue!;
                        }
                        else if (parameters[i].ParameterType.IsValueType)
                        {
                            args[i] = Activator.CreateInstance(
                                parameters[i].ParameterType)!;
                        }
                        else
                        {
                            matched = false;
                            break;
                        }
                    }
                    if (!matched) continue;

                    try
                    {
                        var driver = (IPlcDriver)ctor.Invoke(args);
                        manager.RegisterDriver(driver);
                        count++;
                        break;
                    }
                    catch
                    {
                        // Try the next constructor.
                    }
                }
            }
            return count;
        }

        private static bool IsDriverBase(Type type)
        {
            var baseType = type.BaseType;
            while (baseType != null)
            {
                if (baseType.FullName ==
                    "org.apache.plc4net.spi.drivers.DriverBase")
                    return true;
                baseType = baseType.BaseType;
            }
            return false;
        }
    }
}
