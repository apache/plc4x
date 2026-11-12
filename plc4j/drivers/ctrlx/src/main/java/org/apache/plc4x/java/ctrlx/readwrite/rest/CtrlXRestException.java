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

package org.apache.plc4x.java.ctrlx.readwrite.rest;

/**
 * A ctrlX REST call failed: the request couldn't be sent, the device answered with an error
 * status, or the answer wasn't what the Data Layer API describes.
 */
public class CtrlXRestException extends Exception {

    /** HTTP status of the answer, or 0 if there was none. */
    private final int statusCode;

    public CtrlXRestException(String message) {
        this(message, 0, null);
    }

    public CtrlXRestException(String message, Throwable cause) {
        this(message, 0, cause);
    }

    public CtrlXRestException(String message, int statusCode, Throwable cause) {
        super(message, cause);
        this.statusCode = statusCode;
    }

    public int getStatusCode() {
        return statusCode;
    }

}
