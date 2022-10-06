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
package org.apache.plc4x.java.bacnetip;

import org.apache.commons.io.FileUtils;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionEvent;
import org.apache.plc4x.java.spi.values.PlcStruct;

import java.io.File;
import java.net.URL;
import java.nio.file.FileSystems;

public class PassiveBacNetIpDriver2Manual {

    public static void main(String[] args) throws Exception {
        String tempDirectory = FileUtils.getTempDirectoryPath();
        File pcapFile = FileSystems.getDefault().getPath(tempDirectory, "bacnet-stack-services.cap").toFile();
        if (!pcapFile.exists())
            FileUtils.copyURLToFile(new URL("https://wiki.wireshark.org/Protocols/bacnet?action=AttachFile&do=get&target=bacnet-stack-services.cap"), pcapFile);

        final BacNetIpDriver driver = new BacNetIpDriver();
        final PlcConnection connection = driver.getConnection(
            "bacnet-ip:pcap://" + pcapFile.getAbsolutePath());
        connection.connect();
        final PlcSubscriptionResponse subscriptionResponse = connection.subscriptionRequestBuilder().addEventField(
            "Hurz", "*/*/*").build().execute().get();
        subscriptionResponse.getSubscriptionHandle("Hurz").register(plcSubscriptionEvent -> {
            PlcStruct plcStruct = (PlcStruct)
                ((DefaultPlcSubscriptionEvent) plcSubscriptionEvent).getValues().get("event").getValue();
            System.out.println(plcStruct);
        });
    }

}
