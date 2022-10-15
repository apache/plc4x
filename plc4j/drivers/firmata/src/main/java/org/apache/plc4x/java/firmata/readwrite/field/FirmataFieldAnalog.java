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
package org.apache.plc4x.java.firmata.readwrite.field;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class FirmataFieldAnalog extends FirmataField {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("analog:" + FirmataField.ADDRESS_PATTERN);

    public FirmataFieldAnalog(int address, Integer quantity) {
        super(address, quantity);
    }

    @Override
    public String getPlcDataType() {
        return "INT";
    }

    public static FirmataFieldAnalog of(String addressString) throws PlcInvalidFieldException {
        Matcher matcher = ADDRESS_PATTERN.matcher(addressString);
        if (!matcher.matches()) {
            throw new PlcInvalidFieldException(addressString, ADDRESS_PATTERN);
        }
        int address = Integer.parseInt(matcher.group("address"));

        String quantityString = matcher.group("quantity");
        Integer quantity = quantityString != null ? Integer.valueOf(quantityString) : null;
        return new FirmataFieldAnalog(address, quantity);
    }

}
