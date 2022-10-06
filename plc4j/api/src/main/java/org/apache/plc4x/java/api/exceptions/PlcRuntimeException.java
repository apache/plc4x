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
 * Most generic type of plc4x runtime exception. All plc4x exceptions are derived from
 * PlcException.
 */
public class PlcRuntimeException extends RuntimeException {
    private static final long serialVersionUID = 1L;

    public PlcRuntimeException(String message) {
        super(message);
    }

    public PlcRuntimeException(String message, Throwable cause) {
        super(message, cause);
    }

    public PlcRuntimeException(Throwable cause) {
        super(cause);
    }

    public PlcRuntimeException(String message, Throwable cause, boolean enableSuppression, boolean writableStackTrace) {
        super(message, cause, enableSuppression, writableStackTrace);
    }

}
