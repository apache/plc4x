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
package org.apache.plc4x.java.umas.manual;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcDriverManager;
import org.apache.plc4x.java.api.messages.PlcBrowseItem;
import org.apache.plc4x.java.api.messages.PlcBrowseResponse;
import org.apache.plc4x.java.api.model.PlcQuery;
import org.apache.plc4x.java.api.types.PlcResponseCode;

import java.util.stream.Collectors;

/**
 * In order to run this using the simulator, we need to forward port 502 on all devices to 127.0.0.1 by running the
 * following command on a shell with admin privileges:
 *   netsh interface portproxy add v4tov4 listenaddress=0.0.0.0 listenport=502 connectaddress=127.0.0.1 connectport=502
 */
public class ManualUmasBrowse {

    public static void main(String[] args) throws Exception {
        long startTime = System.currentTimeMillis();
        try (PlcConnection connection = PlcDriverManager.getDefault().getConnectionManager().getConnection("umas://192.168.42.99")){
            PlcBrowseResponse plcBrowseResponse = connection.browseRequestBuilder()
                .addQuery("all", "**")
                .build().executeWithInterceptor((queryName, query, item) -> {
                    outputItem(queryName, query, item, 0);
                    return true;
                }).get();

            long endTime = System.currentTimeMillis();
            int numNodes = 0;
            for (String queryName : plcBrowseResponse.getQueryNames()) {
                if (plcBrowseResponse.getResponseCode(queryName) != PlcResponseCode.OK) {
                    continue;
                }

                for (PlcBrowseItem value : plcBrowseResponse.getValues(queryName)) {
                    numNodes += countNodes(value);
                }
            }
            System.out.printf("Took %dms returned %d nodes%n", endTime - startTime, numNodes);
        }
    }

    protected static void outputItem(String queryName, PlcQuery query, PlcBrowseItem item, int level) {
        System.out.printf("%s- %s: name: %s address: %s - type: %s%s%s%n",
            "  ".repeat(level),
            queryName,
            item.getName(),
            item.getTag().getAddressString(),
            item.getTag().getPlcValueType(),
            (item.getArrayInformation() != null) && !item.getArrayInformation().isEmpty() ? " " + item.getArrayInformation().stream().map(arrayInfo -> "[" + arrayInfo.getLowerBound() + ".." + arrayInfo.getUpperBound() + "]").collect(Collectors.joining()) : "",
            (item.getOptions() != null) && !item.getOptions().isEmpty() ? " {" + item.getOptions().entrySet().stream().map(stringPlcValueEntry -> stringPlcValueEntry.getKey() + ": \"" + stringPlcValueEntry.getValue().toString() + "\"").collect(Collectors.joining()) + "}" : "");
        if ((item.getChildren() != null) && !item.getChildren().isEmpty()) {
            item.getChildren().forEach((s, plcBrowseItem) -> outputItem(queryName, query, plcBrowseItem, level + 1));
        }
    }

    protected static int countNodes(PlcBrowseItem browseItem) {
        if(browseItem.getChildren().isEmpty()) {
            return 1;
        }

        int nodes = 1;
        for (PlcBrowseItem childBrowseItem : browseItem.getChildren().values()) {
            nodes += countNodes(childBrowseItem);
        }
        return nodes;
    }

}
