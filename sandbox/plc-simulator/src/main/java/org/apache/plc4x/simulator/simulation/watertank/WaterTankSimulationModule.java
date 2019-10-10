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
package org.apache.plc4x.simulator.simulation.watertank;

import org.apache.plc4x.simulator.simulation.SimulationModule;

import java.util.Map;
import java.util.TreeMap;
import java.util.concurrent.TimeUnit;

public class WaterTankSimulationModule implements SimulationModule {

    private static final String PROP_WATER_LEVEL = "waterLevel";

    private final Map<String, Object> context;

    public WaterTankSimulationModule() {
        context = new TreeMap<>();
        context.put(PROP_WATER_LEVEL, 0);
    }

    @Override
    public String getName() {
        return "Water Tank";
    }

    @Override
    public Map<String, Object> getContext() {
        return context;
    }

    @Override
    public void loop() {
        // TODO: Do something sensible ;-)
        try {
            // Just increase the level by 1 (Whatever this means ...
            context.put(PROP_WATER_LEVEL, ((Integer) context.get(PROP_WATER_LEVEL)) + 1);
            TimeUnit.MILLISECONDS.sleep(100);
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
    }

}
