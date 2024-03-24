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

package org.apache.plc4x.java.ctrlx.readwrite;

import org.apache.plc4x.java.ctrlx.readwrite.rest.datalayer.ApiClient;
import org.apache.plc4x.java.ctrlx.readwrite.rest.datalayer.api.DataLayerInformationAndSettingsApi;
import org.apache.plc4x.java.ctrlx.readwrite.rest.datalayer.model.BrowseData;
import org.apache.plc4x.java.ctrlx.readwrite.utils.ApiClientFactory;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;

public class BrowseTest {

    @Test
    @Disabled("Only made to run locally")
    public void browseTest() throws Exception {
        ApiClient apiClient = ApiClientFactory.getApiClient("https://192.168.42.100", "boschrexroth", "TimechoEurope2023");
        //
        DataLayerInformationAndSettingsApi api = new DataLayerInformationAndSettingsApi(apiClient);
        BrowseData nodeNames = api.getNodeNames();
        for (String nodeName : nodeNames.getValue()) {
            System.out.println(nodeName);
        }
    }

}
