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
package org.apache.plc4x.java.examples.helloopm;

import org.apache.plc4x.java.opm.OPMException;
import org.apache.plc4x.java.opm.PlcEntity;
import org.apache.plc4x.java.opm.PlcEntityManager;
import org.apache.plc4x.java.opm.PlcTag;
import org.apache.plc4x.java.utils.cache.CachedPlcConnectionManager;

/**
 * This Example shows how to use OPM from plc4j via the @{@link PlcEntityManager}.
 * A @{@link CachedPlcConnectionManager} is used to optimize the access and to allow for automatic reconnection.
 *
 * The {@link PlcEntityManager} is similar to JPAs EntityManager.
 * The "connected" Entity (shootCounter) can be kept and passed around and stays connected in the sense that all calls
 * to a getter are forwarded to the PLC.
 * Finally, one can disconnect the Entity.
 *
 * This MT works against Tims S7 in Nürtingen.
 * Thus, parameters have to be tuned possibly to get "good" values.
 */
public class HelloOpm {

    private static final String PLC_ADDRESS = "s7://192.168.23.30";
    private static final String PLC_TAG_ADDRESS = "%DB4:16:SINT";
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
        entityManager = new PlcEntityManager(CachedPlcConnectionManager.getBuilder().build());
    }

    /**
     * The {@link PlcEntityManager#connect(Class, String)} method returns a "connected" Entity, i.e., a proxy Object.
     * Whenever a getter is called on the Proxy object (whose Tag Variable is annotated with @{@link PlcEntity}
     * a call to the PLC is made to fetch the value.
     * If another method is called on the Entity all Tags are fetched from the Plc first, and then the method is
     * invoked.
     *
     * @throws OPMException
     */
    public void readValueFromPlcUsingConnectedEntity() throws OPMException {
        // Fetch connected Entity
        DistanceSensor distanceSensor = entityManager.connect(DistanceSensor.class, PLC_ADDRESS);
        // Read shoot values a hundred times
        long distance = distanceSensor.getDistance();
        System.out.println("Current distance: " + distance);
        // Write the values back ...
        //entityManager.write(DistanceSensor.class, PLC_ADDRESS, distanceSensor);
        // Disconnect the Entity (not necessary)
        entityManager.disconnect(distanceSensor);
    }

    /**
     * The {@link PlcEntityManager#read(Class, String)} method fetches all tags annotated with @{@link PlcTag}
     * <b>once</b> and injects them in the new instance. After the constructing this is a regular POJO with no fancy
     * functionality.
     *
     * @throws OPMException
     */
    public void readValueFromPlcUsingRead() throws OPMException {
        // Read Entity from PLC
        DistanceSensor distanceSensor = entityManager.read(DistanceSensor.class, PLC_ADDRESS);
        System.out.println("Current distance: " + distanceSensor.getDistance());
    }

    /**
     * Example entity which maps one tag on a PLC where a Distance Sensor is connected.
     */
    @PlcEntity
    public static class DistanceSensor {

        @PlcTag("%DB4:0.0:BOOL")
        private boolean bitValue;

        @PlcTag("%DB4:1:BYTE")
        private short byteValue;

        @PlcTag("%DB4:2:WORD")
        private int wordValue;

        @PlcTag("%DB4:4:DWORD")
        private long dwordValue;

        @PlcTag("%DB4:16:SINT")
        private byte sintValue;

        @PlcTag("%DB4:17:USINT")
        private short usintValue;

        @PlcTag("%DB4:18:INT")
        private short intValue;

        @PlcTag("%DB4:20:UINT")
        private int uintValue;

        @PlcTag("%DB4:22:DINT")
        private int dintValue;

        @PlcTag("%DB4:26:UDINT")
        private long udintValue;

        // TODO: For some reason it doesn't work with the simple type "float"
        @PlcTag("%DB4:46:REAL")
        private Float realValue;

        @PlcTag(PLC_TAG_ADDRESS)
        private long distance;

        public long getDistance() {
            return distance;
        }

    }
}