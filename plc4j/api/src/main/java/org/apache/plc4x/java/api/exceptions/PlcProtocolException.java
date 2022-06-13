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
package org.apache.plc4x.java.api.exceptions;

/**
 * Base exception for all protocol related exceptions.
 * This is usually if the messages themselves are malformed, contain incorrect information or
 * any other problems occur which are related to the protocol being implemented.
 */
public class PlcProtocolException extends PlcException {
    private static final long serialVersionUID = 1L;

    public PlcProtocolException(String message) {
        super(message);
    }

    public PlcProtocolException(String message, Throwable cause) {
        super(message, cause);
    }

    public PlcProtocolException(Throwable cause) {
        super(cause);
    }

    public PlcProtocolException(String message, Throwable cause, boolean enableSuppression, boolean writableStackTrace) {
        super(message, cause, enableSuppression, writableStackTrace);
    }

}
