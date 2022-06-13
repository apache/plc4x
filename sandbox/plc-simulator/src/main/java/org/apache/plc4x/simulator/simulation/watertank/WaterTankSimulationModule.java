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
package org.apache.plc4x.simulator.simulation.watertank;

import org.apache.plc4x.simulator.model.Context;
import org.apache.plc4x.simulator.simulation.SimulationModule;

import java.util.concurrent.TimeUnit;

/**
 * This is a little simulation that simulates a Water tank.
 * This tank has a capacity who's "waterLevel" is represented as a Long value.
 * Water can flow into the tank if the input valve is opened and it can flow
 * out of the tank if the output valve is open.
 *
 * The capacity of the output is slightly smaller than that of the input, so
 * opening both valves will result in the tank filling.
 *
 * To prevent the tank from bursting, there's an emergency valve which is opened
 * as soon as the water-level reaches a critical maximum.
 */
public class WaterTankSimulationModule implements SimulationModule {

    private static final long MAX_WATER_LEVEL = 27648L;
    private static final long EMERGENCY_VALVE_WATER_LEVEL = 27500L;

    private static final int NUM_INPUT_VALVE_INPUT = 0;
    private static final int NUM_OUTPUT_VALVE_INPUT = 1;

    private static final int EMERGENCY_VALVE_OUTPUT = 0;

    private static final String PROP_WATER_LEVEL = "waterLevel";

    private final Context context;

    public WaterTankSimulationModule() {
        context = new Context.ContextBuilder()
            // The input valve
            .addDigitalInput(false)
            // The output valve
            .addDigitalInput(true)
            // The emergency valve
            .addDigitalOutput(false)
            // The water level
            .addMemoryVariable(PROP_WATER_LEVEL, 0L).build();
    }

    @Override
    public String getName() {
        return "Water Tank";
    }

    @Override
    public Context getContext() {
        return context;
    }

    @Override
    public void loop() {
        // TODO: Do something sensible ;-)
        // Just a POC for now ... to be replaced by a "real" simulation ...
        try {
            // Get the current value for the water tank level.
            Long value = (Long) context.getMemory().get(PROP_WATER_LEVEL);

            // If the input valve is open, add 10.
            if(context.getDigitalInputs().get(NUM_INPUT_VALVE_INPUT)) {
                value += 10;
                value = Math.min(MAX_WATER_LEVEL, value);
            }

            // If the output valve is open, subtract 8 (It's slightly less throughput than the input)
            if(context.getDigitalInputs().get(NUM_OUTPUT_VALVE_INPUT)) {
                value -= 8;
                value = Math.max(0, value);
            }

            // Calculate if the emergency valve should be open
            boolean emergencyValveOpen = value > EMERGENCY_VALVE_WATER_LEVEL;

            // Update the memory.
            context.getMemory().put(PROP_WATER_LEVEL, value);
            // Update the state of the emergency valve.
            context.getDigitalOutputs().set(EMERGENCY_VALVE_OUTPUT, emergencyValveOpen);

            // Should probably be handled by the simulator
            TimeUnit.MILLISECONDS.sleep(100);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }

}
