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
package org.apache.plc4x.simulator;

import org.apache.plc4x.simulator.server.ServerModule;
import org.apache.plc4x.simulator.simulation.SimulationModule;

import java.util.Map;
import java.util.ServiceLoader;
import java.util.TreeMap;
import java.util.concurrent.TimeUnit;

public class PlcSimulator {

    private boolean running;
    private static Map<String, ServerModule> serverModules;
    private static Map<String, SimulationModule> simulationModules;

    public PlcSimulator() {
        this(Thread.currentThread().getContextClassLoader());
    }

    public PlcSimulator(ClassLoader classLoader) {
        // Initialize all the server modules.
        serverModules = new TreeMap<>();
        ServiceLoader<ServerModule> serverModuleLoader = ServiceLoader.load(ServerModule.class, classLoader);
        for (ServerModule serverModule : serverModuleLoader) {
            serverModules.put(serverModule.getName(), serverModule);
        }

        // Initialize all the simulation modules.
        simulationModules = new TreeMap<>();
        ServiceLoader<SimulationModule> simulationModuleLoader = ServiceLoader.load(SimulationModule.class, classLoader);
        for (SimulationModule simulationModule : simulationModuleLoader) {
            simulationModules.put(simulationModule.getName(), simulationModule);
        }

        running = true;
    }

    public void stop() {
        running = false;
    }

    public void run() throws Exception {
        // Start all server modules.
        for (ServerModule serverModule : serverModules.values()) {
            serverModule.start();
        }

        try {
            while (running) {
                // Give all the simulation modules the chance to do something.
                for (SimulationModule simulationModule : simulationModules.values()) {
                    simulationModule.loop();
                }
                // Sleep 100 ms to not run the simulation too eagerly.
                TimeUnit.MILLISECONDS.sleep(100);
            }
        } finally {
            // Start all server modules.
            for (ServerModule serverModule : serverModules.values()) {
                serverModule.stop();
            }
        }
    }

    public static void main(String[] args) throws Exception {
        new PlcSimulator().run();
    }

}
