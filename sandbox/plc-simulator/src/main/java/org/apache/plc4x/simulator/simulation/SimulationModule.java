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
package org.apache.plc4x.simulator.simulation;

import java.util.Map;

public interface SimulationModule {

    /**
     * @return the name of the simulation module
     */
    String getName();

    /**
     * Gives access to the internal simulations context.
     * This is an immutable map of named properties that should contain only simple data-types.
     * @return reference to the simulations context
     */
    Map<String, Object> getContext();

    /**
     * Method for doing the actual processing inside the simulation.
     * In this method the simulation can do calculations and update it's context variables.
     */
    void loop();

}
