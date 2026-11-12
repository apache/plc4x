/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.profinet.gsdml;

import org.w3c.dom.Element;

import java.io.IOException;

public class ProfinetProfileBody {

    private ProfinetDeviceIdentity deviceIdentity;

    private ProfinetDeviceFunction deviceFunction;

    private ProfinetApplicationProcess applicationProcess;

    public ProfinetDeviceIdentity getDeviceIdentity() {
        return deviceIdentity;
    }

    public ProfinetDeviceFunction getDeviceFunction() {
        return deviceFunction;
    }

    public ProfinetApplicationProcess getApplicationProcess() {
        return applicationProcess;
    }

    static ProfinetProfileBody fromXml(Element element) throws IOException {
        ProfinetProfileBody result = new ProfinetProfileBody();
        for (GsdmlXml.Property property : GsdmlXml.properties(element, true)) {
            switch (property.name()) {
                case "DeviceIdentity" -> result.deviceIdentity = property.object(ProfinetDeviceIdentity::fromXml);
                case "DeviceFunction" -> result.deviceFunction = property.object(ProfinetDeviceFunction::fromXml);
                case "ApplicationProcess" -> result.applicationProcess = property.object(ProfinetApplicationProcess::fromXml);
                // Strict, as under the former jackson mapping: anything unknown here rejects the whole file.
                default -> throw property.unknownIn(ProfinetProfileBody.class);
            }
        }
        return result;
    }
}
