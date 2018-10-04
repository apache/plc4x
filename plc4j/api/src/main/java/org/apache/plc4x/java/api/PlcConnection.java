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
package org.apache.plc4x.java.api;

import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;

import java.util.Optional;

/**
 * Interface defining the most basic methods a PLC4X connection should support.
 * This generally handles the connection establishment itself and the parsing of
 * field address strings to the platform dependent PlcField instances.
 * <p>
 * The individual operations are then defined by other interfaces within this package.
 */
public interface PlcConnection extends AutoCloseable {

    /**
     * Established the connection to the remote PLC.
     *
     * @throws PlcConnectionException an exception if the connection attempt failed.
     */
    void connect() throws PlcConnectionException;

    /**
     * Returns true if the PlcConnection is connected to a remote PLC.
     *
     * @return true, if connected, false, if not.
     */
    boolean isConnected();

    /**
     * Closes the connection to the remote PLC.
     *
     * @throws Exception an exception if shutting down the connection failed.
     */
    @Override
    void close() throws Exception;

    Optional<PlcReadRequest.Builder> readRequestBuilder();

    Optional<PlcWriteRequest.Builder> writeRequestBuilder();

    Optional<PlcSubscriptionRequest.Builder> subscriptionRequestBuilder();

    Optional<PlcUnsubscriptionRequest.Builder> unsubscriptionRequestBuilder();

}
