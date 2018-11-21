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
package org.apache.plc4x.java.ethernetip;

import com.digitalpetri.enip.EtherNetIpClientConfig;
import com.digitalpetri.enip.cip.CipClient;
import com.digitalpetri.enip.cip.epath.EPath;
import com.digitalpetri.enip.cip.epath.LogicalSegment;
import com.digitalpetri.enip.cip.epath.PortSegment;
import com.digitalpetri.enip.cip.services.GetAttributeSingleService;

import java.time.Duration;

public class EnipTest {

    public static void main(String[] args) throws Exception {
        EtherNetIpClientConfig config = EtherNetIpClientConfig.builder("10.10.64.30")
            .setSerialNumber(0x00)
            .setVendorId(0x00)
            .setTimeout(Duration.ofSeconds(2))
            .build();

        // backplane, slot 0
        EPath.PaddedEPath connectionPath = new EPath.PaddedEPath(
            new PortSegment(1, new byte[]{(byte) 0}));

        CipClient client = new CipClient(config, connectionPath);

        client.connect().get();

        GetAttributeSingleService service = new GetAttributeSingleService(
            new EPath.PaddedEPath(new LogicalSegment.ClassId(0x04), new LogicalSegment.InstanceId(0x69), new LogicalSegment.AttributeId(0x03)));

        client.invoke(service).whenComplete((as, ex) -> {
            if (as != null) {
                try {
                    short value = as.readShort();
                    System.out.println(String.format("Value is %s", value));
                } catch (Throwable t) {
                    t.printStackTrace();
                }
            } else {
                ex.printStackTrace();
            }
        });

    }

}
