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
package org.apache.plc4x.java.bacnetip.model;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.PlcField;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class BacNetIpField implements PlcField {

    private static final Pattern ADDRESS_PATTERN =
        Pattern.compile("^N(?<fileNumber>\\d{1,7}):(?<elementNumber>\\d{1,7})(/(?<bitNumber>\\d{1,7}))?:(?<dataType>[a-zA-Z_]+)(\\[(?<size>\\d+)])?");

    public BacNetIpField() {
    }

    public static boolean matches(String fieldString) {
        return ADDRESS_PATTERN.matcher(fieldString).matches();
    }

    public static BacNetIpField of(String fieldString) {
        Matcher matcher = ADDRESS_PATTERN.matcher(fieldString);
        if(matcher.matches()) {
            return new BacNetIpField();
        }
        throw new PlcInvalidFieldException("Unable to parse address: " + fieldString);
    }

}
