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
package org.apache.plc4x.test.driver.internal.handlers;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcResponse;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBufferXmlBased;
import org.apache.plc4x.java.spi.utils.Serializable;
import org.apache.plc4x.test.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.test.driver.internal.utils.Synchronizer;
import org.apache.plc4x.test.driver.internal.validator.ApiValidator;
import org.dom4j.Element;

import java.util.concurrent.TimeUnit;

public class ApiResponseHandler {

    private final Element payload;

    private final Synchronizer synchronizer;

    public ApiResponseHandler(Element payload, Synchronizer synchronizer) {
        this.payload = payload;
        this.synchronizer = synchronizer;
    }

    public void executeApiResponse() {
        assert synchronizer != null;
        if (synchronizer.responseFuture == null) {
            throw new DriverTestsuiteException("No response expected.");
        }
        PlcResponse plcResponse;
        try {
            plcResponse = synchronizer.responseFuture.get(5000, TimeUnit.MILLISECONDS);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new DriverTestsuiteException("Got no response within 5000ms.", e);
        } catch (Exception e) {
            throw new DriverTestsuiteException("Got no response within 5000ms.", e);
        }

        // Reset the future.
        synchronizer.responseFuture = null;
        final String serializedResponse = serializeToXmlString((Serializable) plcResponse);
        ApiValidator.validateApiMessage(payload, serializedResponse);
    }

    private String serializeToXmlString(Serializable value) {
        WriteBufferXmlBased writeBuffer = new WriteBufferXmlBased();
        try {
            value.serialize(writeBuffer);
        } catch (SerializationException e) {
            throw new PlcRuntimeException(e);
        }
        return writeBuffer.getXmlString();
    }

}
