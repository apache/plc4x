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

public class ProfinetApplicationProcess {

    private List<ProfinetDeviceAccessPointItem> deviceAccessPointList;

    private List<ProfinetModuleItem> moduleList;

    private List<ProfinetLogBookEntryItem> logBookEntryList;

    private List<ProfinetGraphicItem> graphicsList;

    private ProfinetExternalTextList externalTextList;

    public List<ProfinetDeviceAccessPointItem> getDeviceAccessPointList() {
        return deviceAccessPointList;
    }

    public List<ProfinetModuleItem> getModuleList() {
        return moduleList;
    }

    public List<ProfinetLogBookEntryItem> getLogBookEntryList() {
        return logBookEntryList;
    }

    public List<ProfinetGraphicItem> getGraphicsList() {
        return graphicsList;
    }

    public ProfinetExternalTextList getExternalTextList() {
        return externalTextList;
    }

    static ProfinetApplicationProcess fromXml(Element element) throws IOException {
        ProfinetApplicationProcess result = new ProfinetApplicationProcess();
        for (GsdmlXml.Property property : GsdmlXml.properties(element, false)) {
            switch (property.name()) {
                case "DeviceAccessPointList" -> result.deviceAccessPointList = property.wrappedList(ProfinetDeviceAccessPointItem::fromXml);
                case "ModuleList" -> result.moduleList = property.wrappedList(ProfinetModuleItem::fromXml);
                case "LogBookEntryList" -> result.logBookEntryList = property.wrappedList(ProfinetLogBookEntryItem::fromXml);
                case "GraphicsList" -> result.graphicsList = property.wrappedList(ProfinetGraphicItem::fromXml);
                case "ExternalTextList" -> result.externalTextList = property.object(ProfinetExternalTextList::fromXml);
                default -> {
                    // Not used by the driver.
                }
            }
        }
        return result;
    }
}
