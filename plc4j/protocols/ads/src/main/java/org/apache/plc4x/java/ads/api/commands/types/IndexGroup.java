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
package org.apache.plc4x.java.ads.api.commands.types;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.util.UnsignedIntLEByteValue;

@SuppressWarnings("unused") // Due to predefined IndexGroups
public class IndexGroup extends UnsignedIntLEByteValue {

    public static final int NUM_BYTES = UnsignedIntLEByteValue.UNSIGNED_INT_LE_NUM_BYTES;

    private IndexGroup(byte... values) {
        super(values);
    }

    private IndexGroup(long value) {
        super(value);
    }

    private IndexGroup(String value) {
        super(value);
    }

    private IndexGroup(ByteBuf byteBuf) {
        super(byteBuf);
    }

    public static IndexGroup of(byte... values) {
        return new IndexGroup(values);
    }

    public static IndexGroup of(long value) {
        return new IndexGroup(value);
    }

    public static IndexGroup of(String value) {
        return new IndexGroup(value);
    }

    public static IndexGroup of(ByteBuf byteBuf) {
        return new IndexGroup(byteBuf);
    }

    public static final class ReservedGroups {
        public static final IndexGroup ADSIGRP_SYMTAB = IndexGroup.of(0xF000);
        public static final IndexGroup ADSIGRP_SYMNAME = IndexGroup.of(0xF001);
        public static final IndexGroup ADSIGRP_SYMVAL = IndexGroup.of(0xF002);
        public static final IndexGroup ADSIGRP_SYM_HNDBYNAME = IndexGroup.of(0xF003);
        public static final IndexGroup ADSIGRP_SYM_VALBYNAME = IndexGroup.of(0xF004);
        public static final IndexGroup ADSIGRP_SYM_VALBYHND = IndexGroup.of(0xF005);
        public static final IndexGroup ADSIGRP_SYM_RELEASEHND = IndexGroup.of(0xF006);
        public static final IndexGroup ADSIGRP_SYM_INFOBYNAME = IndexGroup.of(0xF007);
        public static final IndexGroup ADSIGRP_SYM_VERSION = IndexGroup.of(0xF008);
        public static final IndexGroup ADSIGRP_SYM_INFOBYNAMEEX = IndexGroup.of(0xF009);
        public static final IndexGroup ADSIGRP_SYM_DOWNLOAD = IndexGroup.of(0xF00A);
        public static final IndexGroup ADSIGRP_SYM_UPLOAD = IndexGroup.of(0xF00B);
        public static final IndexGroup ADSIGRP_SYM_UPLOADINFO = IndexGroup.of(0xF00C);
        public static final IndexGroup ADSIGRP_SYMNOTE = IndexGroup.of(0xF010);
        public static final IndexGroup ADSIGRP_IOIMAGE_RWIB = IndexGroup.of(0xF020);
        public static final IndexGroup ADSIGRP_IOIMAGE_RWIX = IndexGroup.of(0xF021);
        public static final IndexGroup ADSIGRP_IOIMAGE_RISIZE = IndexGroup.of(0xF025);
        public static final IndexGroup ADSIGRP_IOIMAGE_RWOB = IndexGroup.of(0xF030);
        public static final IndexGroup ADSIGRP_IOIMAGE_RWOX = IndexGroup.of(0xF031);
        public static final IndexGroup ADSIGRP_IOIMAGE_RWOSIZE = IndexGroup.of(0xF035);
        public static final IndexGroup ADSIGRP_IOIMAGE_CLEARI = IndexGroup.of(0xF040);
        public static final IndexGroup ADSIGRP_IOIMAGE_CLEARO = IndexGroup.of(0xF050);
        public static final IndexGroup ADSIGRP_IOIMAGE_RWIOB = IndexGroup.of(0xF060);
        public static final IndexGroup ADSIGRP_DEVICE_DATA = IndexGroup.of(0xF100);
        public static final IndexGroup ADSIOFFS_DEVDATA_ADSSTATE = IndexGroup.of(0x0000);
        public static final IndexGroup ADSIOFFS_DEVDATA_DEVSTATE = IndexGroup.of(0x0002);

        private ReservedGroups() {
            // Container class
        }
    }

    public static final class SystemServiceGroups {
        public static final IndexGroup SYSTEMSERVICE_OPENCREATE = IndexGroup.of(100);
        public static final IndexGroup SYSTEMSERVICE_OPENREAD = IndexGroup.of(101);
        public static final IndexGroup SYSTEMSERVICE_OPENWRITE = IndexGroup.of(102);
        public static final IndexGroup SYSTEMSERVICE_CREATEFILE = IndexGroup.of(110);
        public static final IndexGroup SYSTEMSERVICE_CLOSEHANDLE = IndexGroup.of(111);
        public static final IndexGroup SYSTEMSERVICE_FOPEN = IndexGroup.of(120);
        public static final IndexGroup SYSTEMSERVICE_FCLOSE = IndexGroup.of(121);
        public static final IndexGroup SYSTEMSERVICE_FREAD = IndexGroup.of(122);
        public static final IndexGroup SYSTEMSERVICE_FWRITE = IndexGroup.of(123);
        public static final IndexGroup SYSTEMSERVICE_FSEEK = IndexGroup.of(124);
        public static final IndexGroup SYSTEMSERVICE_FTELL = IndexGroup.of(125);
        public static final IndexGroup SYSTEMSERVICE_FGETS = IndexGroup.of(126);
        public static final IndexGroup SYSTEMSERVICE_FPUTS = IndexGroup.of(127);
        public static final IndexGroup SYSTEMSERVICE_FSCANF = IndexGroup.of(128);
        public static final IndexGroup SYSTEMSERVICE_FPRINTF = IndexGroup.of(129);
        public static final IndexGroup SYSTEMSERVICE_FEOF = IndexGroup.of(130);
        public static final IndexGroup SYSTEMSERVICE_FDELETE = IndexGroup.of(131);
        public static final IndexGroup SYSTEMSERVICE_FRENAME = IndexGroup.of(132);
        public static final IndexGroup SYSTEMSERVICE_REG_HKEYLOCALMACHINE = IndexGroup.of(200);
        public static final IndexGroup SYSTEMSERVICE_SENDEMAIL = IndexGroup.of(300);
        public static final IndexGroup SYSTEMSERVICE_TIMESERVICES = IndexGroup.of(400);
        public static final IndexGroup SYSTEMSERVICE_STARTPROCESS = IndexGroup.of(500);
        public static final IndexGroup SYSTEMSERVICE_CHANGENETID = IndexGroup.of(600);

        private SystemServiceGroups() {
            // Container class
        }
    }
}
