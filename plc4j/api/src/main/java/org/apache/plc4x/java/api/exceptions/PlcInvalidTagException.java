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
 * Indicates an invalid tag address.
 */
public class PlcInvalidTagException extends PlcRuntimeException {

    private static final long serialVersionUID = 1L;
    private static final String DOES_NOT_MATCH = " doesn't match ";
    private final String tagToBeParsed;

    public PlcInvalidTagException(String tagToBeParsed) {
        super(tagToBeParsed + " invalid");
        this.tagToBeParsed = tagToBeParsed;
    }

    public PlcInvalidTagException(String tagToBeParsed, Throwable cause) {
        super(tagToBeParsed + " invalid", cause);
        this.tagToBeParsed = tagToBeParsed;
    }

    public PlcInvalidTagException(String tagToBeParsed, Pattern pattern) {
        super(tagToBeParsed + DOES_NOT_MATCH + pattern);
        this.tagToBeParsed = tagToBeParsed;
    }

    public PlcInvalidTagException(String tagToBeParsed, Pattern pattern, Throwable cause) {
        super(tagToBeParsed + DOES_NOT_MATCH + pattern, cause);
        this.tagToBeParsed = tagToBeParsed;
    }

    public PlcInvalidTagException(String tagToBeParsed, Pattern pattern, String readablePattern) {
        super(tagToBeParsed + DOES_NOT_MATCH + readablePattern + '(' + pattern + ')');
        this.tagToBeParsed = tagToBeParsed;
    }

    public String getTagToBeParsed() {
        return tagToBeParsed;
    }

}
