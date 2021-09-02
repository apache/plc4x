/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.examples.helloopm;

import org.apache.plc4x.java.opm.OPMException;
import org.apache.plc4x.java.opm.PlcEntity;
import org.apache.plc4x.java.opm.PlcEntityManager;
import org.apache.plc4x.java.opm.PlcField;
import org.apache.plc4x.java.utils.connectionpool.PooledPlcDriverManager;

/**
 * This Example shows how to use OPM from plc4j via the @{@link PlcEntityManager}.
 * A @{@link PooledPlcDriverManager} is used to optimize the acces and to allow for automatic reconnection.
 *
 * The {@link PlcEntityManager} is similar to JPAs EntityManager.
 * The "connected" Entity (shootCounter) can be kept and passed around and stays connected in the sense that all calls
 * to a getter are forwared to the PLC.
 * Finally, one can disconnect the Entity.
 *
 * This MT works against Tims S7 in Nürtingen.
 * Thus, parameters have to be tuned possibly to get "good" values.
 *
 * @author julian
 * Created by julian on 31.10.18
 */
public class HelloOpm {

    private static final String ADDRESS = "s7://192.168.167.210/0/0";
    private static final String PLC_FIELD_ADDRESS = "%DB225.DBW0:INT";
    private final PlcEntityManager entityManager;

    public static void main(String[] args) throws OPMException {
        HelloOpm helloOpm = new HelloOpm();
        // Do a fetch via connected entity
        helloOpm.readValueFromPlcUsingConnectedEntity();
        // Do a fetch via read
        helloOpm.readValueFromPlcUsingRead();
        // Stop the application
        System.exit(0);
    }

    public HelloOpm() {
        entityManager = new PlcEntityManager(new PooledPlcDriverManager());
    }

    /**
     * The {@link PlcEntityManager#connect(Class, String)} method returns a "connected" Entity, i.e., a proxy Object.
     * Whenever a getter is called on the Proxy object (whose Field Variable is annotated with @{@link PlcEntity}
     * a call to the PLC is made to fetch the value.
     * If another method is called on the Entity all Fields are feched from the Plc first, and then the method is
     * invoked.
     *
     * @throws OPMException
     */
    public void readValueFromPlcUsingConnectedEntity() throws OPMException {
        // Fetch connected Entity
        DistanceSensor distanceSensor = entityManager.connect(DistanceSensor.class, ADDRESS);
        // Read shoot values a hundred times
        long distance = distanceSensor.getDistance();
        System.out.println("Current distance: " + distance);
        // Disconnect the Entity (not necessary)
        entityManager.disconnect(distanceSensor);
    }

    /**
     * The {@link PlcEntityManager#read(Class, String)} method fetches all fields annotated with @{@link PlcField}
     * <b>once</b> and injects them in the new instance. After the constructing this is a regular POJO with no fancy
     * functionality.
     *
     * @throws OPMException
     */
    public void readValueFromPlcUsingRead() throws OPMException {
        // Read Entity from PLC
        DistanceSensor distanceSensor = entityManager.read(DistanceSensor.class, ADDRESS);
        System.out.println("Current distance: " + distanceSensor.getDistance());
    }

    /**
     * Example entity which maps one field on a PLC where a Distance Sensor is connected.
     */
    @PlcEntity
    public static class DistanceSensor {

        @PlcField(PLC_FIELD_ADDRESS)
        private long distance;

        public long getDistance() {
            return distance;
        }

    }
}