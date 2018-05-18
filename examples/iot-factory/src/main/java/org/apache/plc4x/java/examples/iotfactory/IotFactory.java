package org.apache.plc4x.java.examples.iotfactory;/*
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

import org.apache.edgent.function.Supplier;
import org.apache.edgent.providers.direct.DirectProvider;
import org.apache.edgent.topology.TStream;
import org.apache.edgent.topology.Topology;
import org.apache.plc4x.edgent.PlcConnectionAdapter;
import org.apache.plc4x.edgent.PlcFunctions;

import java.util.concurrent.TimeUnit;

public class IotFactory {

    public enum ConveyorState {
        STOPPED,
        RUNNING_LEFT,
        RUNNING_RIGHT
    };

    private static int smallBoxes = 0;
    private static int largeBoxes = 0;
    private static ConveyorState conveyorState = ConveyorState.STOPPED;


    public static void main(String[] args) throws Exception {
        // Get a plc connection.
        try (PlcConnectionAdapter plcAdapter = new PlcConnectionAdapter("s7://10.10.64.20/1/1")) {
            // Initialize the Edgent core.
            DirectProvider dp = new DirectProvider();
            Topology top = dp.newTopology();

            // Define the event stream.
            // 1) PLC4X source generating a stream of bytes.
            Supplier<Byte> plcSupplier = PlcFunctions.byteSupplier(plcAdapter, "OUTPUTS/0");
            // 2) Use polling to get an item from the byte-stream in regular intervals.
            TStream<Byte> source = top.poll(plcSupplier, 100, TimeUnit.MILLISECONDS);
            // 3) Output the events in the stream on the console.

            source.sink(value -> {
                boolean runningLeft = (value & 8) != 0;
                boolean runningRight = (value & 16) != 0;

                if(conveyorState == ConveyorState.STOPPED) {
                    if(runningLeft | runningRight) {
                        if (runningLeft) {
                            smallBoxes++;
                            conveyorState = ConveyorState.RUNNING_LEFT;
                            updateOutput();
                        } else {
                            largeBoxes++;
                            conveyorState = ConveyorState.RUNNING_RIGHT;
                            updateOutput();
                        }
                    }
                } else if (!(runningLeft | runningRight)){
                    conveyorState = ConveyorState.STOPPED;
                }
            });

            // Submit the topology and hereby start the event streams.
            dp.submit(top);
        }
    }

    private static void updateOutput() {
        System.out.println(String.format("Small Boxes: %3d, Large Boxes %3d", smallBoxes, largeBoxes));
    }

}
