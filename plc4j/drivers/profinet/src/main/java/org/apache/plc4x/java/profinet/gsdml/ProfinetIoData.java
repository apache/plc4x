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
import java.util.List;

public class ProfinetIoData {

    private int iopsLength;

    private int iocsLength;

    private List<ProfinetIoDataInput> input;

    private List<ProfinetIoDataOutput> output;

    public int getIopsLength() {
        return iopsLength;
    }

    public int getIocsLength() {
        return iocsLength;
    }

    public List<ProfinetIoDataInput> getInput() {
        return input;
    }

    public List<ProfinetIoDataOutput> getOutput() {
        return output;
    }

    static ProfinetIoData fromXml(Element element) throws IOException {
        ProfinetIoData result = new ProfinetIoData();
        for (GsdmlXml.Property property : GsdmlXml.properties(element, false)) {
            switch (property.name()) {
                case "IOPS_Length" -> result.iopsLength = property.intValue();
                case "IOCS_Length" -> result.iocsLength = property.intValue();
                case "Input" -> result.input = property.unwrappedList(ProfinetIoDataInput::fromXml);
                case "Output" -> result.output = property.unwrappedList(ProfinetIoDataOutput::fromXml);
                default -> {
                    // Not used by the driver.
                }
            }
        }
        return result;
    }
}
