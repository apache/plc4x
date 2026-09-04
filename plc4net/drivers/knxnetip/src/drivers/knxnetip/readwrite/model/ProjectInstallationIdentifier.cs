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

// Code generated from the mspec by plc4net-code-gen. DO NOT EDIT.

using System;
using System.Linq;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.generation;

namespace org.apache.plc4net.drivers.knxnetip.readwrite.model
{
    public partial class ProjectInstallationIdentifier : IMessage
    {
        public byte ProjectNumber { get; }
        public byte InstallationNumber { get; }

        public ProjectInstallationIdentifier(byte projectNumber, byte installationNumber)
        {
            ProjectNumber = projectNumber;
            InstallationNumber = installationNumber;
        }

        public static ProjectInstallationIdentifier StaticParse(ReadBuffer readBuffer)
        {
            var projectNumber = readBuffer.ReadByte("projectNumber", 8);
            var installationNumber = readBuffer.ReadByte("installationNumber", 8);
            return new ProjectInstallationIdentifier(projectNumber, installationNumber);
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("projectNumber", 8, ProjectNumber);
            writeBuffer.WriteByte("installationNumber", 8, InstallationNumber);
        }

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            return lengthInBits;
        }

        public int GetLengthInBytes() => (GetLengthInBits() + 7) / 8;
    }
}
