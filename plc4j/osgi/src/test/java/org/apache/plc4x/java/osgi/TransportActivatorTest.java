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

import org.apache.plc4x.java.spi.transports.api.Transport;
import org.junit.jupiter.api.Test;
import org.osgi.framework.ServiceReference;

import java.util.Collection;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

@SuppressWarnings("rawtypes")
class TransportActivatorTest extends FrameworkTestBase {

    @Test
    void startRegistersEachTransportWithItsCodeAndName() throws Exception {
        new TransportActivator().start(context);

        Collection<ServiceReference<Transport>> refs =
            context.getServiceReferences(Transport.class, "(org.apache.plc4x.transport.code=osgi-test)");
        assertEquals(1, refs.size());
        ServiceReference<Transport> ref = refs.iterator().next();
        assertEquals("OSGi Test Transport", ref.getProperty("org.apache.plc4x.transport.name"));
        assertTrue(context.getService(ref) instanceof TestTransport);
    }

    @Test
    void stopUnregistersTheTransports() throws Exception {
        TransportActivator activator = new TransportActivator();
        activator.start(context);

        activator.stop(context);

        assertTrue(context.getServiceReferences(Transport.class, null).isEmpty());
    }

}
