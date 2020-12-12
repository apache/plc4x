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
package org.apache.plc4x.java.s7.netty.model.payloads.ssls;

public class SslModuleIdentificationDataRecord implements SslDataRecord {

    public static final short INDEX_MODULE = 0x0001;
    public static final short INDEX_BASIC_HARDWARE = 0x0006;
    public static final short INDEX_BASIC_FIRMWARE = 0x0007;

    private short index;
    private String articleNumber;
    private short bgType;
    private short moduleOrOsVersion;
    private short pgDescriptionFileVersion;

    public SslModuleIdentificationDataRecord(short index, String articleNumber, short bgType, short moduleOrOsVersion, short pgDescriptionFileVersion) {
        this.index = index;
        this.articleNumber = articleNumber;
        this.bgType = bgType;
        this.moduleOrOsVersion = moduleOrOsVersion;
        this.pgDescriptionFileVersion = pgDescriptionFileVersion;
    }

    @Override
    public short getLengthInWords() {
        return 14;
    }

    public short getIndex() {
        return index;
    }

    public String getArticleNumber() {
        return articleNumber;
    }

    public short getBgType() {
        return bgType;
    }

    public short getModuleOrOsVersion() {
        return moduleOrOsVersion;
    }

    public short getPgDescriptionFileVersion() {
        return pgDescriptionFileVersion;
    }

}
