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
package org.apache.plc4x.java.knxnetip.maual;

import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcBrowseItem;
import org.apache.plc4x.java.api.messages.PlcBrowseResponse;
import org.apache.plc4x.java.api.model.PlcQuery;

import java.util.stream.Collectors;

/**
 * Manual test for browsing KNX group addresses from an ETS project file.
 * Browse is only available when a knxproj file is configured.
 */
public class ManualKnxNetIpBrowse {

    public static void main(String[] args) throws Exception {
        try (PlcConnection connection = new DefaultPlcDriverManager().getConnection(
            "knxnet-ip://192.168.42.28?" +
            "knxproj-file-path=/Users/christoferdutz/Projects/Privat/NLNet/plc4x/plc4j/drivers/knxnetip/Stettiner-Str-13.knxproj&" +
            "knxproj-password=cW171998$")) {

            // Create a browse request for all group addresses
            PlcBrowseResponse plcBrowseResponse = connection.browseRequestBuilder()
                .addQuery("all", "*")
                .build()
                .executeWithInterceptor((queryName, query,item) -> {
                    outputItem(queryName, query, item, 0);
                    return true;
                }).get();
            System.out.println(plcBrowseResponse);
        }
    }

    protected static void outputItem(String queryName, PlcQuery query, PlcBrowseItem item, int level) {
        if (item.getOptions().get("dataPointType") == null) {
            return;
        }
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

}
