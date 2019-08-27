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
package org.apache.plc4x.java.abeth.model;

import org.apache.plc4x.java.abeth.model.types.FileType;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.PlcField;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class AbEthField implements PlcField {

    private static final Pattern ADDRESS_PATTERN =
        Pattern.compile("^N(?<fileNumber>\\d{1,7})\\:(?<elementNumber>\\d{1,7})/(?<subElementNumber>\\d{1,7}):(?<dataType>[a-zA-Z_]+)(\\[(?<size>\\d+)])?");

    private static final String FILE_NUMBER = "fileNumber";
    private static final String ELEMENT_NUMBER = "elementNumber";
    private static final String SUB_ELEMENT_NUMBER = "subElementNumber";
    private static final String DATA_TYPE = "dataType";
    private static final String SIZE = "size";

    private final short byteSize;
    private final short fileNumber;
    private final FileType fileType;
    private final short elementNumber;
    private final short subElementNumber;

    public AbEthField(short byteSize, short fileNumber, FileType fileType, short elementNumber, short subElementNumber) {
        this.byteSize = byteSize;
        this.fileNumber = fileNumber;
        this.fileType = fileType;
        this.elementNumber = elementNumber;
        this.subElementNumber = subElementNumber;
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

    public short getElementNumber() {
        return elementNumber;
    }

    public short getSubElementNumber() {
        return subElementNumber;
    }

    public static AbEthField of(String fieldString) {
        Matcher matcher = ADDRESS_PATTERN.matcher(fieldString);
        if(matcher.matches()) {
            short fileNumber = Short.parseShort(matcher.group(FILE_NUMBER));
            short elementNumber = Short.parseShort(matcher.group(ELEMENT_NUMBER));
            short subElementNumber = Short.parseShort(matcher.group(SUB_ELEMENT_NUMBER));
            FileType fileType = FileType.valueOf(Short.parseShort(matcher.group(DATA_TYPE)));
            short byteSize = Short.parseShort(matcher.group(SIZE));
            return new AbEthField(byteSize, fileNumber, fileType,elementNumber, subElementNumber);
        }
        throw new PlcInvalidFieldException("Unable to parse address: " + fieldString);
    }

}
