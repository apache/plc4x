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
package org.apache.plc4x.java.applications.iotree;

import org.apache.edgent.function.Consumer;
import org.apache.edgent.providers.direct.DirectProvider;
import org.apache.edgent.topology.TStream;
import org.apache.edgent.topology.Topology;
import org.apache.plc4x.edgent.PlcConnectionAdapter;
import org.apache.plc4x.edgent.PlcFunctions;

import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;

public class IoTree {

    private final PlcConnectionAdapter plcAdapter;

    private AtomicInteger internalVariables;
    private AtomicInteger digitalInput;
    private AtomicInteger analogInput;
    private AtomicInteger amplification;
    private AtomicInteger maxValue;
    private AtomicInteger digitalOutput;

    private IoTree(PlcConnectionAdapter plcAdapter) {
        this.plcAdapter = plcAdapter;
        internalVariables = new AtomicInteger(0);
        digitalInput = new AtomicInteger(0);
        analogInput = new AtomicInteger(0);
        amplification = new AtomicInteger(0);
        maxValue = new AtomicInteger(0);
        digitalOutput = new AtomicInteger(0);
    }

    private void run() throws Exception {
        DirectProvider dp = new DirectProvider();
        Topology top = dp.newTopology();
        // Automatically update the internal variable values ...
        TStream<Byte> internalVariableStream = top.poll(PlcFunctions.byteSupplier(plcAdapter, "DATA_BLOCKS/0"), 10, TimeUnit.MILLISECONDS);
        internalVariableStream.sink((Consumer<Byte>) inputs -> internalVariables.set(inputs));
        // Automatically update the digital input values ...
        TStream<Byte> digitalInputStream = top.poll(PlcFunctions.byteSupplier(plcAdapter, "INPUTS/0"), 10, TimeUnit.MILLISECONDS);
        digitalInputStream.sink((Consumer<Byte>) inputs -> digitalInput.set(inputs));
        // Automatically update the analog input values ...
        TStream<Short> analogInputStream = top.poll(PlcFunctions.shortSupplier(plcAdapter, "INPUTS/66"), 10, TimeUnit.MILLISECONDS);
        analogInputStream.sink((Consumer<Short>) inputs -> analogInput.set(inputs));
        // Automatically update the amplification ...
        TStream<Float> amplificationStream = top.poll(PlcFunctions.floatSupplier(plcAdapter, "DATA_BLOCKS/1/0"), 10, TimeUnit.MILLISECONDS);
        amplificationStream.sink((Consumer<Float>) inputs -> amplification.set(inputs.intValue()));
        // Automatically update the maxValue ...
        TStream<Integer> maxValueStream = top.poll(PlcFunctions.integerSupplier(plcAdapter, "DATA_BLOCKS/1/4"), 10, TimeUnit.MILLISECONDS);
        maxValueStream.sink((Consumer<Integer>) inputs -> maxValue.set(inputs));
        // Automatically update the digital output values ...
        TStream<Byte> digitalOutputStream = top.poll(PlcFunctions.byteSupplier(plcAdapter, "OUTPUTS/0"), 10, TimeUnit.MILLISECONDS);
        digitalOutputStream.sink((Consumer<Byte>) outputs -> digitalOutput.set(outputs));

        // Start logging ...
        dp.submit(top);

        while (true) {
            Thread.sleep(1000);
            System.out.println(String.format("Internal Variables: %d Digital In: %d Analog In: %d Amplifictaion: %d Max Value: %d Digital Out: %d",
                internalVariables.get(), digitalInput.get(), analogInput.get(), amplification.get(), maxValue.get(), digitalOutput.get()));
        }
    }

    public static void main(String[] args) {
        if((args == null) || (args.length != 1)) {
            System.out.println("Usage: IoTree {connection-string}");
            System.out.println("Example: IoTree s7://192.168.0.1/0/0");
            System.exit(1);
        }

        try (PlcConnectionAdapter plcAdapter = new PlcConnectionAdapter(args[0])) {
            // Initialize the tree itself
            IoTree tree = new IoTree(plcAdapter);
            // Start the tree ...
            tree.run();
            // Yeah ... well prevent the application from exiting ;-)
            while (true) {
                Thread.sleep(1000);
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

}
