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
package org.apache.plc4x.java.ethernetip.model;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.PlcField;

import java.util.Objects;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class EtherNetIpField implements PlcField {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("^#(?<class>.*?)#(?<instance>\\d{1,4})(?:#(?<attribute>\\d))?");

    public static boolean matches(String addressString) {
        return ADDRESS_PATTERN.matcher(addressString).matches();
    }

    public static EtherNetIpField of(String addressString) throws PlcInvalidFieldException {
        Matcher matcher = ADDRESS_PATTERN.matcher(addressString);
        if (!matcher.matches()) {
            throw new PlcInvalidFieldException(addressString, ADDRESS_PATTERN);
        }
        int classNumber = Integer.parseInt(matcher.group("class"));
        int instanceNumber = Integer.parseInt(matcher.group("instance"));
        int attributeNumber = Integer.parseInt(matcher.group("attribute"));

        return new EtherNetIpField(classNumber, instanceNumber, attributeNumber);
    }

    private final int objectNumber;
    private final int instanceNumber;
    private final int attributeNumber;

    private int connectionId;

    public EtherNetIpField(int objectNumber, int instanceNumber, int attributeNumber) {
        this.objectNumber = objectNumber;
        this.instanceNumber = instanceNumber;
        this.attributeNumber = attributeNumber;

        this.connectionId = -1;
    }

    public int getObjectNumber() {
        return objectNumber;
    }

    public int getInstanceNumber() {
        return instanceNumber;
    }

    public int getAttributeNumber() {
        return attributeNumber;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof EtherNetIpField)) {
            return false;
        }
        EtherNetIpField that = (EtherNetIpField) o;
        return objectNumber == that.objectNumber &&
            instanceNumber == that.instanceNumber &&
            attributeNumber == that.attributeNumber;
    }

    @Override
    public int hashCode() {

        return Objects.hash(objectNumber, instanceNumber, attributeNumber);
    }

    @Override
    public String toString() {
        return "EtherNetIpField{" +
            "object-number=" + objectNumber +
            ", instance-number=" + instanceNumber +
            ", attribute-number=" + attributeNumber +
            '}';
    }
}
