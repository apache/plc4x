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
package org.apache.plc4x.java.api.connection;

import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.model.Address;

import java.util.Optional;

/**
 * Interface defining the most basic methods a PLC4X connection should support.
 * This generally handles the connection establishment itself and the parsing of
 * address strings to the platform dependent Address instances.
 *
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
     * Closes the connection to the remote PLC.
     *
     * @throws Exception an exception if shutting down the connection failed.
     */
    @Override
    void close() throws Exception;

    /**
     * Parses a PLC/protocol dependent address string into an Address object.
     *
     * @param addressString String representation of an address for the current type of PLC/protocol.
     * @return Address object identifying an address for the current type of PLC/protocol.
     * @throws PlcException an exception if there was a problem parsing the address string.
     */
    Address parseAddress(String addressString) throws PlcException;

    Optional<PlcLister> getLister();

    Optional<PlcReader> getReader();

    Optional<PlcWriter> getWriter();

}
