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
package org.apache.plc4x.java.abeth.field;

import org.apache.plc4x.java.abeth.types.FileType;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcValueType;

import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class AbEthField implements PlcField {

    private static final Pattern ADDRESS_PATTERN =
//        Pattern.compile("^N(?<fileNumber>\\d{1,7}):(?<elementNumber>\\d{1,7})/(?<bitNumber>\\d{1,7}):(?<dataType>[a-zA-Z_]+)(\\[(?<size>\\d+)])?");
        Pattern.compile("^N(?<fileNumber>\\d{1,7}):(?<elementNumber>\\d{1,7})(/(?<bitNumber>\\d{1,7}))?:(?<dataType>[a-zA-Z_]+)(\\[(?<size>\\d+)])?");


    private static final String FILE_NUMBER = "fileNumber";
    private static final String ELEMENT_NUMBER = "elementNumber";
//    private static final String SUB_ELEMENT_NUMBER = "subElementNumber";
    private static final String BIT_NUMBER = "bitNumber";
    private static final String DATA_TYPE = "dataType";
    private static final String SIZE = "size";

    private final short byteSize;
    private final short fileNumber;
    private final FileType fileType;
    private final short elementNumber;
    private final short bitNumber;

    public AbEthField(short byteSize, short fileNumber, FileType fileType, short elementNumber, short bitNumber) {
        this.byteSize = byteSize;
        this.fileNumber = fileNumber;
        this.fileType = fileType;
        this.elementNumber = elementNumber;
        this.bitNumber = bitNumber;
    }

    public short getByteSize() {
        return byteSize;
    }

    public short getFileNumber() {
        return fileNumber;
    }

    public FileType getFileType() {
        return fileType;
    }

    @Override
    public String getAddressString() {
        String address = String.format("N%d:%d", fileNumber, elementNumber);
        if(bitNumber != 0) {
            address += "/" + bitNumber;
        }
        address += ":" + fileType.name();
        if(byteSize != 1) {
            address += "[" + byteSize + "]";
        }
        return address;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return fileType.getPlcValueType();
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        return PlcField.super.getArrayInfo();
    }

    public short getElementNumber() {
        return elementNumber;
    }

    public short getBitNumber() {
        return bitNumber;
    }

    public static boolean matches(String fieldString) {
        return ADDRESS_PATTERN.matcher(fieldString).matches();
    }

    public static AbEthField of(String fieldString) {
        Matcher matcher = ADDRESS_PATTERN.matcher(fieldString);
        if(matcher.matches()) {
            short fileNumber = Short.parseShort(matcher.group(FILE_NUMBER));
            short elementNumber = Short.parseShort(matcher.group(ELEMENT_NUMBER));
            short bitNumber = (matcher.group(BIT_NUMBER) != null) ? Short.parseShort(matcher.group(BIT_NUMBER)) : 0;  //Short.parseShort(matcher.group(BIT_NUMBER));
            FileType fileType = FileType.valueOf(matcher.group(DATA_TYPE).toUpperCase());

            short byteSize;
            switch (fileType) {
                case WORD:
                case SINGLEBIT:
                    byteSize = 2;
                    break;
                case DWORD:
                    byteSize = 4;
                    break;
                default:
                    byteSize = Short.parseShort(matcher.group(SIZE));
            }
            return new AbEthField(byteSize, fileNumber, fileType, elementNumber, bitNumber);
        }
        throw new PlcInvalidFieldException("Unable to parse field address: " + fieldString);
    }

}
