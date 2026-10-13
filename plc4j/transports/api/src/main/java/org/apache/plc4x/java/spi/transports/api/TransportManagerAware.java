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
package org.apache.plc4x.java.spi.transports.api;

/**
 * Implemented by drivers whose transport lookup the hosting runtime can replace.
 * <p>
 * By default a driver finds its transports through {@link java.util.ServiceLoader}. In OSGi the transports live in
 * bundles the driver's class loader can't see, so the OSGi activator hands the driver a {@link TransportManager}
 * backed by the service registry instead.
 */
public interface TransportManagerAware {

    void setTransportManager(TransportManager transportManager);

}
