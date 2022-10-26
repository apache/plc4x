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
package org.apache.plc4x.simulator.model;

import org.apache.commons.collections4.list.FixedSizeList;
import org.apache.commons.collections4.map.FixedSizeMap;

import java.util.*;

public class Context {

    private final FixedSizeMap<String, Object> memory;
    private final FixedSizeList<Boolean> digitalInputs;
    private final FixedSizeList<Long> analogInputs;
    private final FixedSizeList<Boolean> digitalOutputs;
    private final FixedSizeList<Long> analogOutputs;

    public Context(FixedSizeMap<String, Object> memory,
                   FixedSizeList<Boolean> digitalInputs, FixedSizeList<Long> analogInputs,
                   FixedSizeList<Boolean> digitalOutputs, FixedSizeList<Long> analogOutputs) {
        this.memory = memory;
        this.digitalInputs = digitalInputs;
        this.analogInputs = analogInputs;
        this.digitalOutputs = digitalOutputs;
        this.analogOutputs = analogOutputs;
    }

    public Map<String, Object> getMemory() {
        return memory;
    }

    public List<Boolean> getDigitalInputs() {
        return digitalInputs;
    }

    public List<Long> getAnalogInputs() {
        return analogInputs;
    }

    public List<Boolean> getDigitalOutputs() {
        return digitalOutputs;
    }

    public List<Long> getAnalogOutputs() {
        return analogOutputs;
    }

    public static class ContextBuilder {
        private final Map<String, Object> memory;
        private final List<Boolean> digitalInputs;
        private final List<Long> analogInputs;
        private final List<Boolean> digitalOutputs;
        private final List<Long> analogOutputs;

        public ContextBuilder() {
            memory = new TreeMap<>();
            digitalInputs = new LinkedList<>();
            analogInputs = new LinkedList<>();
            digitalOutputs = new LinkedList<>();
            analogOutputs = new LinkedList<>();
        }

        public ContextBuilder addMemoryVariable(String name, Object defaultValue) {
            memory.put(name, defaultValue);
            return this;
        }

        public ContextBuilder addDigitalInput(Boolean defaultValue) {
            digitalInputs.add(defaultValue);
            return this;
        }

        public ContextBuilder addAnalogInput(Long defaultValue) {
            analogInputs.add(defaultValue);
            return this;
        }

        public ContextBuilder addDigitalOutput(Boolean defaultValue) {
            digitalOutputs.add(defaultValue);
            return this;
        }

        public ContextBuilder addAnalogOutput(Long defaultValue) {
            analogOutputs.add(defaultValue);
            return this;
        }

        public Context build() {
            return new Context(FixedSizeMap.fixedSizeMap(memory),
                FixedSizeList.fixedSizeList(digitalInputs), FixedSizeList.fixedSizeList(analogInputs),
                FixedSizeList.fixedSizeList(digitalOutputs), FixedSizeList.fixedSizeList(analogOutputs));
        }

    }

}
