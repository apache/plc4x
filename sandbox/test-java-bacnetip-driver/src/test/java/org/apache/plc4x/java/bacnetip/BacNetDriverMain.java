/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.java.bacnetip;

import java.util.function.Consumer;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.value.PlcStruct;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionEvent;

public class BacNetDriverMain {

    public static void main(String[] args) throws Exception {
        final PlcConnection connection = new PlcDriverManager().getConnection("bacnet-ip-active:udp://192.168.2.106");
        connection.connect();
        PlcSubscriptionRequest plcSubscriptionRequest = connection.subscriptionRequestBuilder()
            .addEventField("Hurz", "*/*/*")
            .build();

        final PlcSubscriptionResponse subscriptionResponse = plcSubscriptionRequest.execute().get();
        subscriptionResponse.getSubscriptionHandle("Hurz").register(new Consumer<PlcSubscriptionEvent>() {
            @Override
            public void accept(PlcSubscriptionEvent plcSubscriptionEvent) {
                PlcStruct plcStruct = (PlcStruct) ((DefaultPlcSubscriptionEvent) plcSubscriptionEvent).getValues()
                    .get("event").getValue();
                System.out.println(plcStruct);
            }
        });
    }

}
