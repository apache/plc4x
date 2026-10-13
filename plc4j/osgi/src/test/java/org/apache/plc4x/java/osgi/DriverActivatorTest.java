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
package org.apache.plc4x.java.osgi;

import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.spi.transports.api.Transport;
import org.junit.jupiter.api.Test;
import org.osgi.framework.ServiceReference;

import java.util.Collection;
import java.util.Optional;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertInstanceOf;
import static org.junit.jupiter.api.Assertions.assertTrue;

@SuppressWarnings("rawtypes")
class DriverActivatorTest extends FrameworkTestBase {

    @Test
    void startRegistersEachDriverWithItsCodeAndName() throws Exception {
        new DriverActivator().start(context);

        Collection<ServiceReference<PlcDriver>> refs =
            context.getServiceReferences(PlcDriver.class, "(" + DriverActivator.PROTOCOL_CODE + "=osgi-test)");
        assertEquals(1, refs.size());
        assertEquals("OSGi Test Driver", refs.iterator().next().getProperty(DriverActivator.PROTOCOL_NAME));
    }

    @Test
    void driverResolvesTransportsFromTheServiceRegistry() throws Exception {
        new DriverActivator().start(context);
        TestDriver driver = registeredDriver();
        assertInstanceOf(ServiceRegistryTransportManager.class, driver.getTransportManager());
        assertFalse(driver.getTransportManager().getTransport("osgi-test").isPresent());

        // The transport bundle starts after the driver bundle: the driver must still find it.
        TransportActivator transports = new TransportActivator();
        transports.start(context);
        Optional<Transport> transport = driver.getTransportManager().getTransport("osgi-test");
        assertTrue(transport.isPresent());
        assertInstanceOf(TestTransport.class, transport.get());
        assertFalse(driver.getTransportManager().getTransport("tcp").isPresent());

        // And lose it again once the transport bundle stops.
        transports.stop(context);
        assertFalse(driver.getTransportManager().getTransport("osgi-test").isPresent());
    }

    @Test
    void stopUnregistersTheDriversAndStopsTrackingTransports() throws Exception {
        DriverActivator activator = new DriverActivator();
        activator.start(context);
        TestDriver driver = registeredDriver();
        new TransportActivator().start(context);

        activator.stop(context);

        assertTrue(context.getServiceReferences(PlcDriver.class, null).isEmpty());
        assertFalse(driver.getTransportManager().getTransport("osgi-test").isPresent());
    }

    @Test
    void stopWithoutStartIsHarmless() {
        assertDoesNotThrow(() -> new DriverActivator().stop(context));
    }

    private TestDriver registeredDriver() throws Exception {
        ServiceReference<PlcDriver> ref = context.getServiceReferences(PlcDriver.class, null).iterator().next();
        return (TestDriver) context.getService(ref);
    }

}
