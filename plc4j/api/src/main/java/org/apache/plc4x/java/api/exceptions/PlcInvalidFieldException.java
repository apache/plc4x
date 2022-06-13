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

import java.util.regex.Pattern;

/**
 * Indicates an invalid field address.
 */
public class PlcInvalidFieldException extends PlcRuntimeException {

    private static final long serialVersionUID = 1L;
    private static final String DOES_NOT_MATCH = " doesn't match ";
    private final String fieldToBeParsed;

    public PlcInvalidFieldException(String fieldToBeParsed) {
        super(fieldToBeParsed + " invalid");
        this.fieldToBeParsed = fieldToBeParsed;
    }

    public PlcInvalidFieldException(String fieldToBeParsed, Throwable cause) {
        super(fieldToBeParsed + " invalid", cause);
        this.fieldToBeParsed = fieldToBeParsed;
    }

    public PlcInvalidFieldException(String fieldToBeParsed, Pattern pattern) {
        super(fieldToBeParsed + DOES_NOT_MATCH + pattern);
        this.fieldToBeParsed = fieldToBeParsed;
    }

    public PlcInvalidFieldException(String fieldToBeParsed, Pattern pattern, Throwable cause) {
        super(fieldToBeParsed + DOES_NOT_MATCH + pattern, cause);
        this.fieldToBeParsed = fieldToBeParsed;
    }

    public PlcInvalidFieldException(String fieldToBeParsed, Pattern pattern, String readablePattern) {
        super(fieldToBeParsed + DOES_NOT_MATCH + readablePattern + '(' + pattern + ')');
        this.fieldToBeParsed = fieldToBeParsed;
    }

    public String getFieldToBeParsed() {
        return fieldToBeParsed;
    }

}
