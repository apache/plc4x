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

package org.apache.plc4x.java.profinet.device;

import org.apache.plc4x.java.api.messages.PlcBrowseItem;
import org.apache.plc4x.java.profinet.gsdml.ProfinetVirtualSubmoduleItem;
import org.apache.plc4x.java.profinet.readwrite.PnIoCm_IoCs;
import org.apache.plc4x.java.profinet.readwrite.PnIoCm_IoDataObject;
import org.apache.plc4x.java.profinet.readwrite.PnIoCm_Submodule;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;

public interface ProfinetModule {

    List<PnIoCm_Submodule> getExpectedSubModuleApiBlocks();
    List<PnIoCm_IoDataObject> getInputIoPsApiBlocks();
    List<PnIoCm_IoCs> getOutputIoCsApiBlocks();
    List<PnIoCm_IoCs> getInputIoCsApiBlocks();
    List<PnIoCm_IoDataObject> getOutputIoPsApiBlocks();
    Integer getIdentNumber();
    Integer getSlotNumber();

    Map<String, List<PlcBrowseItem>> browseTags(Map<String, List<PlcBrowseItem>> browseItems);

}
