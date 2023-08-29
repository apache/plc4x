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
package org.apache.plc4x.merlot.api.core;

import java.util.HashMap;
import java.util.Map;


public class Merlot {
   
   
    public enum STATE {
        
        STOPPED(0x0000,"TEXTO"),
        RESETTING(0x0002,"TEXTO"),
        IDLE(0x0004,"TEXTO"),
        STARTING(0x0006,"TEXTO"),
        EXECUTE(0x0008,"TEXTO"),
        COMPLETING(0x000A,"TEXTO"),
        COMPLETE(0x000C,"TEXTO"),
        ABORTING(0x000E,"TEXTO"),
        ABORTED(0x0010,"TEXTO"),
        CLEARING(0x0012,"TEXTO"),
        STOPPING(0x0014,"TEXTO"),
        HOLDING(0x0016,"TEXTO"),
        HELD(0x0018,"TEXTO"),
        UNHOLDING(0x001A,"TEXTO"),
        SUSPENDING(0x001C,"TEXTO"),
        SUSPENDED(0x001E,"TEXTO"),
        UNSUSPENDING(0x0020,"TEXTO");
        
        private static final Map<Integer, STATE> map;
        
        static {
            map = new HashMap<>();
            for (STATE state : STATE.values()) {
                map.put(state.code, state);
            }
        }; 
        
        private final String description;
        private final int code;    
        
        STATE(int code, String description){
            this.description = description;
            this.code = code;
        }   
        
        public String getDescription(){
            return description;
        }    

        public int getCode() {
            return code;
        } 
        
        public static STATE  valueOf(int code) {
            return map.get(code);
        }           
        
    }    
    
    
    public enum STATUS {
        M_OK(0x6000,"TEXTO"),
        M_INITIATE_IND(0x6001,"TEXTO"),
        M_CYCL_READ_ABORT_IND(0x6002,"TEXTO"),
        M_CYCL_READ_IND(0x6003,"TEXTO"),

        M_ERR(0x8000,"TEXTO"),
        M_ERR_ADD_VALUE_PDU(0x8001,"TEXTO"),
        M_ERR_BLOCK_ALREADY_EXISTS(0x8002,"TEXTO"),
        M_ERR_BLOCK_NOT_FOUND(0x8003,"TEXTO"),
        M_ERR_BLOCK_PROTECTED(0x8004,"TEXTO"),
        M_ERR_BLOCK_TOO_LARGE(0x8005,"TEXTO"),
        M_ERR_CONN_ABORTED(0x8006,"TEXTO"),
        M_ERR_CONN_ALREADY_FREE(0x8007,"TEXTO"),
        M_ERR_CONN_ALREADY_PASSED(0x8008,"TEXTO"),
        M_ERR_CONN_CNF(0x8009,"TEXTO"),
        M_ERR_CONN_NAME_NOT_FOUND(0x800A,"TEXTO"),
        M_ERR_COORD_RULE(0x800B,"TEXTO"),
        M_ERR_DGN_CONN_NOT_ANNOUNCED(0x800C,"TEXTO"),
        M_ERR_DGN_CONN_TOO_MUCH(0x800D,"TEXTO"),
        M_ERR_DGN_INFO_NOT_AVAIL(0x800E,"TEXTO"),
        M_ERR_FEW_PROT_LEVEL(0x800F,"TEXTO"),
        M_ERR_FILE_NOT_FOUND(0x8010,"TEXTO"),
        M_ERR_FW_ERROR(0x8011,"TEXTO"),
        M_ERR_IEC_DATA_TYPE_MISMATCH(0x8012,"TEXTO"),
        M_ERR_IEC_INVALID_REF(0x8013,"TEXTO"),
        M_ERR_IEC_LOWER_LAYER(0x8014,"TEXTO"),
        M_ERR_IEC_NEG_RESPONSE(0x8015,"TEXTO"),
        M_ERR_IEC_NO(0x8016,"TEXTO"),
        M_ERR_IEC_NO_ACCESS_TO_REM_OBJECT(0x8017,"TEXTO"),
        M_ERR_IEC_PARTNER_IN_WRONG_STATE(0x8018,"TEXTO"),
        M_ERR_IEC_RECEIVER_DISABLED(0x8019,"TEXTO"),
        M_ERR_IEC_RECEIVER_OVERRUN(0x801A,"TEXTO"),
        M_ERR_IEC_RESET_RECEIVED(0x801B,"TEXTO"),
        M_ERR_INSTALL(0x801C,"TEXTO"),
        M_ERR_INTERNAL_ERROR(0x801D,"TEXTO"),
        M_ERR_INVALID_BLOCK(0x801E,"TEXTO"),
        M_ERR_INVALID_BLOCK_NR(0x801F,"TEXTO"),
        M_ERR_INVALID_BLOCK_TYPE(0x8020,"TEXTO"),
        M_ERR_INVALID_CONN_STATE(0x8021,"TEXTO"),
        M_ERR_INVALID_CREF(0x8022,"TEXTO"),
        M_ERR_INVALID_CYCL_READ_STATE(0x8023,"TEXTO"),
        M_ERR_INVALID_DATA_SIZE(0x8024,"TEXTO"),
        M_ERR_INVALID_DATARANGE_OR_TYPE(0x8025,"TEXTO"),
        M_ERR_INVALID_FILENAME(0x8026,"TEXTO"),
        M_ERR_INVALID_FILETYPE(0x8027,"TEXTO"),
        M_ERR_INVALID_ORDERIDM(0x8028,"TEXTO"),
        M_ERR_INVALID_PARAMETER(0x8029,"TEXTO"),
        M_ERR_INVALID_PASSWORD(0x802A,"TEXTO"),
        M_ERR_INVALID_READ_BUFFER(0x802B,"TEXTO"),
        M_ERR_INVALID_SEGMENT(0x802C,"TEXTO"),
        M_ERR_INVALID_WRITE_BUFFER(0x802D,"TEXTO"),
        M_ERR_MAX_REQ(0x802E,"TEXTO"),
        M_ERR_NO_PASSWORD(0x802F,"TEXTO"),
        M_ERR_NO_RCV_BLOCK(0x8030,"TEXTO"),
        M_ERR_NO_RESOURCE(0x8031,"TEXTO"),
        M_ERR_NO_SIN_SERV(0x8032,"TEXTO"),
        M_ERR_OBJ_ACCESS_DENIED(0x8033,"TEXTO"),
        M_ERR_OBJ_ATTR_INCONSISTENT(0x8034,"TEXTO"),
        M_ERR_OBJ_UNDEFINED(0x8035,"TEXTO"),
        M_ERR_ORDERID_USED(0x8036,"TEXTO"),
        M_ERR_R_ID_USED(0x8037,"TEXTO"),
        M_ERR_RECEIVE_BUFFER_FULL(0x8038,"TEXTO"),
        M_ERR_REM_BRCV(0x8039,"TEXTO"),
        M_ERR_REM_BSEND(0x803A,"TEXTO"),
        M_ERR_REM_BSEND_CANCEL(0x803B,"TEXTO"),
        M_ERR_REM_DATABASE_TOO_SMALL(0x803C,"TEXTO"),
        M_ERR_RETRY(0x803D,"TEXTO"),
        M_ERR_SERVICE_CONN_ALREADY_USED(0x803E,"TEXTO"),
        M_ERR_SERVICE_NOT_SUPPORTED(0x803F,"TEXTO"),
        M_ERR_SERVICE_VFD_ALREADY_USED(0x8040,"TEXTO"),
        M_ERR_SYMB_ADDRESS(0x8041,"TEXTO"),
        M_ERR_SYMB_ADDRESS_INCONSISTENT(0x8042,"TEXTO"),
        M_ERR_TOO_LONG_DATA(0x8043,"TEXTO"),
        M_ERR_UNKNOWN_ERROR(0x8044,"TEXTO"),
        M_ERR_WRONG_CP_DESCR(0x8045,"TEXTO"),
        M_ERR_WRONG_IND_CNF(0x8046,"TEXTO");
        
        private static final Map<Integer, STATUS> map;
    
        static {
            map = new HashMap<>();
            for (STATUS  status : STATUS.values()) {
                map.put(status.code, status);
            }
        }    
        
        private final String description;
        private final int code;
        
        STATUS(int code, String description){
            this.description = description;
            this.code = code;
        }   
        
        public String getDescription(){
            return description;
        }    

        public int getCode() {
            return code;
        } 
        
        public static STATUS  valueOf(int code) {
            return map.get(code);
        }           
        
    }    
    
    public enum FUNCTION {
        
        //DRIVER CONFIG
        M_INITIATE_CNF(0x4000,"TEXTO"),
        M_READ_CNF(0x4001,"TEXTO"),
        M_WRITE_CNF(0x4002,"TEXTO"),
        M_MULTIPLE_READ_CNF(0x4003,"TEXTO"),
        M_MULTIPLE_WRITE_CNF(0x4004,"TEXTO"),
        M_CYCL_READ_INIT_CNF(0x4005,"TEXTO"),
        M_CYCL_READ_START_CNF(0x4006,"TEXTO"),
        M_CYCL_READ_STOP_CNF(0x4007,"TEXTO"),
        M_CYCL_READ_DELETE_CNF(0x4008,"TEXTO"),
        M_PASSWORD_CNF(0x4009,"TEXTO"),
        M_PASSWORD_LEN(0x400A,"TEXTO"),
        M_RESET_PASSWORD(0x400B,"TEXTO"),        
        
        //DEVICE
        FC_INITIALIZE_DEVICE(0x0000,"TEXTO"),
        FC_SHUTDOWN_DEVICE(0x0002,"TEXTO"),	

        FC_GET_DEVICE_MODE(0x0004,"TEXTO"),	
        FC_STAR_DEVICE(0x0006,"TEXTO"),
        FC_STOP_DEVICE(0x0008,"TEXTO"),	
        FC_GET_PDU_SIZE(0x0010,"TEXTO"),

        //BITS
        FC_READ_MEMORY_BITS(0x0012,"TEXTO"),
        FC_WRIE_MEMORY_BITS(0x0014,"TEXTO"),

        FC_READ_INPUT_BITS(0x0016,"TEXTO"),

        FC_READ_OUTPUT_BITS(0x0018,"TEXTO"),
        FC_WRITE_OUTPUT_BITS(0x0020,"TEXTO"),

        FC_READ_DATA_BITS(0x0022,"TEXTO"),
        FC_WRITE_DATA_BITS(0x0024,"TEXTO"),

        //BYTES
        FC_READ_MEMORY_BYTES(0x0026,"TEXTO"),
        FC_WRITE_MEMORY_BYTES(0x0028,"TEXTO"),

        FC_READ_INPUT_BYTES(0x0030,"TEXTO"),

        FC_READ_OUTPUT_BYTES(0x0032,"TEXTO"),
        FC_WRITE_OUTPUT_BYTES(0x0034,"TEXTO"),

        FC_READ_DATA_BYTES(0x0036,"TEXTO"),
        FC_WRITE_DATA_BYTES(0x0038,"TEXTO"),

        //WORD
        FC_READ_MEMORY_WORDS(0x0040,"TEXTO"),
        FC_WRITE_MEMORY_WORDS(0x0042,"TEXTO"),

        FC_READ_INPUT_WORDS(0x0044,"TEXTO"),

        FC_READ_OUTPUT_WORDS(0x0046,"TEXTO"),
        FC_WRITE_OUTPUT_WORDS(0x0048,"TEXTO"),

        FC_READ_DATA_WORDS(0x0050,"TEXTO"),
        FC_WRITE_DATA_WORDS(0x0052,"TEXTO"),

        //DWORDS
        FC_READ_MEMORY_DWORDS(0x0054,"TEXTO"),
        FC_WRITE_MEMORY_DWORDS(0x0056,"TEXTO"),

        FC_READ_INPUT_DWORDS(0x0058,"TEXTO"),

        FC_READ_OUTPUT_DWORDS(0x0060,"TEXTO"),
        FC_WRITE_OUTPUT_DWORDS(0x0062,"TEXTO"),

        FC_READ_DATA_DWORDS(0x0064,"TEXTO"),
        FC_WRITE_DATA_DWORDS(0x0066,"TEXTO"),

        //TIMERS
        FC_READ_TIMERS(0x0068,"TEXTO"),
        FC_WRITE_TIMERS(0x0070,"TEXTO"),

        //COUNTER
        FC_READ_COUNTERS(0x0072,"TEXTO"),
        FC_WRITE_COUNTERS(0x0074,"TEXTO");
        
        
        private static final Map<Integer, FUNCTION> map;
    
        static {
            map = new HashMap<>();
            for (FUNCTION  func : FUNCTION.values()) {
                map.put(func.code, func);
            }
        }    
        
        private final String description;
        private final int code;
        
        FUNCTION(int code, String description){
            this.description = description;
            this.code = code;
        }   
        
        public String getDescription(){
            return description;
        }    

        public int getCode() {
            return code;
        } 
        
        public static FUNCTION valueOf(int code) {
            return map.get(code);
        }          
    }
    
    public enum TYPE_IEC {
        M_TYPE_IEC_BOOL(0xC000,"TEXTO"),
        M_TYPE_IEC_SINT(0xC001,"TEXTO"),
        M_TYPE_IEC_INT(0xC002,"TEXTO"),
        M_TYPE_IEC_DINT(0xC003,"TEXTO"),
        M_TYPE_IEC_LINT(0xC004,"TEXTO"),
        M_TYPE_IEC_USINT(0xC005,"TEXTO"),
        M_TYPE_IEC_UINT(0xC006,"TEXTO"),
        M_TYPE_IEC_UDINT(0xC007,"TEXTO"),
        M_TYPE_IEC_ULINT(0xC008,"TEXTO"),
        M_TYPE_IEC_REAL(0xC009,"TEXTO"),
        M_TYPE_IEC_LREAL(0xC00A,"TEXTO"),
        M_TYPE_IEC_TIME(0xC00B,"TEXTO"),
        M_TYPE_IEC_DATE(0xC00C,"TEXTO"),
        M_TYPE_IEC_TIME_OF_DAY(0xC00D,"TEXTO"),
        M_TYPE_IEC_TOD(0xC00E,"TEXTO"),
        M_TYPE_IEC_DATE_AND_TIME(0xC00F,"TEXTO"),
        M_TYPE_IEC_DT(0xC010,"TEXTO"),
        M_TYPE_IEC_STRING(0xC011,"TEXTO"),
        M_TYPE_IEC_BYTE(0xC012,"TEXTO"),
        M_TYPE_IEC_WORD(0xC013,"TEXTO"),
        M_TYPE_IEC_DWORD(0xC014,"TEXTO"),
        M_TYPE_IEC_LWORD(0xC015,"TEXTO"),
        M_TYPE_IEC_WSTRING(0xC016,"TEXTO");
                        
        private static final Map<Integer, TYPE_IEC> map;
    
        static {
            map = new HashMap<>();
            for (TYPE_IEC  func : TYPE_IEC.values()) {
                map.put(func.code, func);
            }
        }    
        
        private final String description;
        private final int code;
        
        TYPE_IEC(int code, String description){
            this.description = description;
            this.code = code;
        }   
        
        public String getDescription(){
            return description;
        }    

        public int getCode() {
            return code;
        } 
        
        public static TYPE_IEC valueOf(int code) {
            return map.get(code);
        }                  
    }
    
    public static final int M_INITIATE_CNF                      = 0x4000;
    public static final int M_READ_CNF                          = 0x4001;
    public static final int M_WRITE_CNF                         = 0x4002;
    public static final int M_MULTIPLE_READ_CNF                 = 0x4003;
    public static final int M_MULTIPLE_WRITE_CNF                = 0x4004;
    public static final int M_CYCL_READ_INIT_CNF                = 0x4005;
    public static final int M_CYCL_READ_START_CNF               = 0x4006;
    public static final int M_CYCL_READ_STOP_CNF                = 0x4007;
    public static final int M_CYCL_READ_DELETE_CNF              = 0x4008;
    public static final int M_PASSWORD_CNF                      = 0x4009;
    public static final int M_PASSWORD_LEN                      = 0x400A;
    public static final int M_RESET_PASSWORD                    = 0x400B;

    public static final int M_OK                                = 0x6000;
    public static final int M_INITIATE_IND                      = 0x6001;
    public static final int M_CYCL_READ_ABORT_IND               = 0x6002;
    public static final int M_CYCL_READ_IND                     = 0x6003;

    public static final int M_ERR                               = 0x8000;
    public static final int M_ERR_ADD_VALUE_PDU                 = 0x8001;
    public static final int M_ERR_BLOCK_ALREADY_EXISTS          = 0x8002;
    public static final int M_ERR_BLOCK_NOT_FOUND               = 0x8003;
    public static final int M_ERR_BLOCK_PROTECTED               = 0x8004;
    public static final int M_ERR_BLOCK_TOO_LARGE               = 0x8005;
    public static final int M_ERR_CONN_ABORTED                  = 0x8006;
    public static final int M_ERR_CONN_ALREADY_FREE             = 0x8007;
    public static final int M_ERR_CONN_ALREADY_PASSED           = 0x8008;
    public static final int M_ERR_CONN_CNF                      = 0x8009;
    public static final int M_ERR_CONN_NAME_NOT_FOUND           = 0x800A;
    public static final int M_ERR_COORD_RULE                    = 0x800B;
    public static final int M_ERR_DGN_CONN_NOT_ANNOUNCED        = 0x800C;
    public static final int M_ERR_DGN_CONN_TOO_MUCH             = 0x800D;
    public static final int M_ERR_DGN_INFO_NOT_AVAIL            = 0x800E;
    public static final int M_ERR_FEW_PROT_LEVEL                = 0x800F;
    public static final int M_ERR_FILE_NOT_FOUND                = 0x8010;
    public static final int M_ERR_FW_ERROR                      = 0x8011;
    public static final int M_ERR_IEC_DATA_TYPE_MISMATCH        = 0x8012;
    public static final int M_ERR_IEC_INVALID_REF               = 0x8013;
    public static final int M_ERR_IEC_LOWER_LAYER               = 0x8014;
    public static final int M_ERR_IEC_NEG_RESPONSE              = 0x8015;
    public static final int M_ERR_IEC_NO                        = 0x8016;
    public static final int M_ERR_IEC_NO_ACCESS_TO_REM_OBJECT   = 0x8017;
    public static final int M_ERR_IEC_PARTNER_IN_WRONG_STATE    = 0x8018;
    public static final int M_ERR_IEC_RECEIVER_DISABLED         = 0x8019;
    public static final int M_ERR_IEC_RECEIVER_OVERRUN          = 0x801A;
    public static final int M_ERR_IEC_RESET_RECEIVED            = 0x801B;
    public static final int M_ERR_INSTALL                       = 0x801C;
    public static final int M_ERR_INTERNAL_ERROR                = 0x801D;
    public static final int M_ERR_INVALID_BLOCK                 = 0x801E;
    public static final int M_ERR_INVALID_BLOCK_NR              = 0x801F;
    public static final int M_ERR_INVALID_BLOCK_TYPE            = 0x8020;
    public static final int M_ERR_INVALID_CONN_STATE            = 0x8021;
    public static final int M_ERR_INVALID_CREF                  = 0x8022;
    public static final int M_ERR_INVALID_CYCL_READ_STATE       = 0x8023;
    public static final int M_ERR_INVALID_DATA_SIZE             = 0x8024;
    public static final int M_ERR_INVALID_DATARANGE_OR_TYPE     = 0x8025;
    public static final int M_ERR_INVALID_FILENAME              = 0x8026;
    public static final int M_ERR_INVALID_FILETYPE              = 0x8027;
    public static final int M_ERR_INVALID_ORDERIDM              = 0x8028;
    public static final int M_ERR_INVALID_PARAMETER             = 0x8029;
    public static final int M_ERR_INVALID_PASSWORD              = 0x802A;
    public static final int M_ERR_INVALID_READ_BUFFER           = 0x802B;
    public static final int M_ERR_INVALID_SEGMENT               = 0x802C;
    public static final int M_ERR_INVALID_WRITE_BUFFER          = 0x802D;
    public static final int M_ERR_MAX_REQ                       = 0x802E;
    public static final int M_ERR_NO_PASSWORD                   = 0x802F;
    public static final int M_ERR_NO_RCV_BLOCK                  = 0x8030;
    public static final int M_ERR_NO_RESOURCE                   = 0x8031;
    public static final int M_ERR_NO_SIN_SERV                   = 0x8032;
    public static final int M_ERR_OBJ_ACCESS_DENIED             = 0x8033;
    public static final int M_ERR_OBJ_ATTR_INCONSISTENT         = 0x8034;
    public static final int M_ERR_OBJ_UNDEFINED                 = 0x8035;
    public static final int M_ERR_ORDERID_USED                  = 0x8036;
    public static final int M_ERR_R_ID_USED                     = 0x8037;
    public static final int M_ERR_RECEIVE_BUFFER_FULL           = 0x8038;
    public static final int M_ERR_REM_BRCV                      = 0x8039;
    public static final int M_ERR_REM_BSEND                     = 0x803A;
    public static final int M_ERR_REM_BSEND_CANCEL              = 0x803B;
    public static final int M_ERR_REM_DATABASE_TOO_SMALL        = 0x803C;
    public static final int M_ERR_RETRY                         = 0x803D;
    public static final int M_ERR_SERVICE_CONN_ALREADY_USED     = 0x803E;
    public static final int M_ERR_SERVICE_NOT_SUPPORTED         = 0x803F;
    public static final int M_ERR_SERVICE_VFD_ALREADY_USED      = 0x8040;
    public static final int M_ERR_SYMB_ADDRESS                  = 0x8041;
    public static final int M_ERR_SYMB_ADDRESS_INCONSISTENT     = 0x8042;
    public static final int M_ERR_TOO_LONG_DATA                 = 0x8043;
    public static final int M_ERR_UNKNOWN_ERROR                 = 0x8044;
    public static final int M_ERR_WRONG_CP_DESCR                = 0x8045;
    public static final int M_ERR_WRONG_IND_CNF                 = 0x8046;

    public static final int M_TYPE_IEC_BOOL                     = 0xC000;
    public static final int M_TYPE_IEC_SINT                     = 0xC001;
    public static final int M_TYPE_IEC_INT                      = 0xC002;
    public static final int M_TYPE_IEC_DINT                     = 0xC003;
    public static final int M_TYPE_IEC_LINT                     = 0xC004;
    public static final int M_TYPE_IEC_USINT                    = 0xC005;
    public static final int M_TYPE_IEC_UINT                     = 0xC006;
    public static final int M_TYPE_IEC_UDINT                    = 0xC007;
    public static final int M_TYPE_IEC_ULINT                    = 0xC008;
    public static final int M_TYPE_IEC_REAL                     = 0xC009;
    public static final int M_TYPE_IEC_LREAL                    = 0xC00A;
    public static final int M_TYPE_IEC_TIME                     = 0xC00B;
    public static final int M_TYPE_IEC_DATE                     = 0xC00C;
    public static final int M_TYPE_IEC_TIME_OF_DAY              = 0xC00D;
    public static final int M_TYPE_IEC_TOD                      = 0xC00E;
    public static final int M_TYPE_IEC_DATE_AND_TIME            = 0xC00F;
    public static final int M_TYPE_IEC_DT                       = 0xC010;
    public static final int M_TYPE_IEC_STRING                   = 0xC011;
    public static final int M_TYPE_IEC_BYTE                     = 0xC012;
    public static final int M_TYPE_IEC_WORD                     = 0xC013;
    public static final int M_TYPE_IEC_DWORD                    = 0xC014;
    public static final int M_TYPE_IEC_LWORD                    = 0xC015;
    public static final int M_TYPE_IEC_WSTRING                  = 0xC016;
        
    //DEVICE
    public final static int FC_INITIALIZE_DEVICE                = 0;
    public final static int FC_SHUTDOWN_DEVICE                  = 2;	

    public final static int FC_GET_DEVICE_MODE                  = 4;	
    public final static int FC_STAR_DEVICE 			= 6;
    public final static int FC_STOP_DEVICE 			= 8;	
    public final static int FC_GET_PDU_SIZE                     = 10;

    //BITS
    public final static int FC_READ_MEMORY_BITS                 = 12;
    public final static int FC_WRIE_MEMORY_BITS                 = 14;

    public final static int FC_READ_INPUT_BITS                  = 16;

    public final static int FC_READ_OUTPUT_BITS                 = 18;
    public final static int FC_WRITE_OUTPUT_BITS                = 20;

    public final static int FC_READ_DATA_BITS                   = 22;
    public final static int FC_WRITE_DATA_BITS                  = 24;

    //BYTES
    public final static int FC_READ_MEMORY_BYTES                = 26;
    public final static int FC_WRITE_MEMORY_BYTES               = 28;

    public final static int FC_READ_INPUT_BYTES                 = 30;

    public final static int FC_READ_OUTPUT_BYTES                = 32;
    public final static int FC_WRITE_OUTPUT_BYTES               = 34;

    public final static int FC_READ_DATA_BYTES                  = 36;
    public final static int FC_WRITE_DATA_BYTES                 = 38;

    //WORD
    public final static int FC_READ_MEMORY_WORDS                = 40;
    public final static int FC_WRITE_MEMORY_WORDS               = 42;

    public final static int FC_READ_INPUT_WORDS                 = 44;

    public final static int FC_READ_OUTPUT_WORDS                = 46;
    public final static int FC_WRITE_OUTPUT_WORDS               = 48;

    public final static int FC_READ_DATA_WORDS                  = 50;
    public final static int FC_WRITE_DATA_WORDS                 = 52;

    //DWORDS
    public final static int FC_READ_MEMORY_DWORDS               = 54;
    public final static int FC_WRITE_MEMORY_DWORDS              = 56;

    public final static int FC_READ_INPUT_DWORDS                = 58;

    public final static int FC_READ_OUTPUT_DWORDS               = 60;
    public final static int FC_WRITE_OUTPUT_DWORDS              = 62;

    public final static int FC_READ_DATA_DWORDS                 = 64;
    public final static int FC_WRITE_DATA_DWORDS                = 66;

    //TIMERS
    public final static int FC_READ_TIMERS 			= 68;
    public final static int FC_WRITE_TIMERS                     = 70;

    //COUNTER
    public final static int FC_READ_COUNTERS                    = 72;
    public final static int FC_WRITE_COUNTERS                   = 74;
    
    /**
 * Return int error description.
 * @param err
 * @return String with the error description.
 */
    public static String getConstantString(int constant){
        String s = "";
        switch(constant){
            
            case M_INITIATE_CNF: s = "M_INITIATE_CNF";break;
            case M_READ_CNF: s = "M_READ_CNF";break;
            case M_WRITE_CNF: s = "M_WRITE_CNF";break;
            case M_MULTIPLE_READ_CNF: s = "M_MULTIPLE_READ_CNF";break;
            case M_MULTIPLE_WRITE_CNF: s = "M_MULTIPLE_WRITE_CNF";break;
            case M_CYCL_READ_INIT_CNF: s = "M_CYCL_READ_INIT_CNF";break;
            case M_CYCL_READ_START_CNF: s = "M_CYCL_READ_START_CNF";break;
            case M_CYCL_READ_STOP_CNF: s = "M_CYCL_READ_STOP_CNF";break;
            case M_CYCL_READ_DELETE_CNF: s = "M_CYCL_READ_DELETE_CNF";break;
            case M_PASSWORD_CNF: s = "M_PASSWORD_CNF";break;
            case M_PASSWORD_LEN: s = "M_PASSWORD_LEN";break;
            case M_RESET_PASSWORD: s = "M_RESET_PASSWORD";break;

            case M_OK: s = "M_OK";break;
            case M_INITIATE_IND: s = "M_INITIATE_IND";break;
            case M_CYCL_READ_ABORT_IND: s = "M_CYCL_READ_ABORT_IND";break;
            case M_CYCL_READ_IND: s = "M_CYCL_READ_IND";break;

            case M_ERR: s = "M_ERR";break;
            case M_ERR_ADD_VALUE_PDU: s = "M_ERR_ADD_VALUE_PDU";break;
            case M_ERR_BLOCK_ALREADY_EXISTS: s = "M_ERR_BLOCK_ALREADY_EXISTS";break;
            case M_ERR_BLOCK_NOT_FOUND: s = "M_ERR_BLOCK_NOT_FOUND";break;
            case M_ERR_BLOCK_PROTECTED: s = "M_ERR_BLOCK_PROTECTED";break;
            case M_ERR_BLOCK_TOO_LARGE: s = "M_ERR_BLOCK_TOO_LARGE";break;
            case M_ERR_CONN_ABORTED: s = "M_ERR_CONN_ABORTED";break;
            case M_ERR_CONN_ALREADY_FREE: s = "M_ERR_CONN_ALREADY_FREE";break;
            case M_ERR_CONN_ALREADY_PASSED: s = "M_ERR_CONN_ALREADY_PASSED";break;
            case M_ERR_CONN_CNF: s = "M_ERR_CONN_CNF";break;
            case M_ERR_CONN_NAME_NOT_FOUND: s = "M_ERR_CONN_NAME_NOT_FOUND";break;
            case M_ERR_COORD_RULE: s = "M_ERR_COORD_RULE";break;
            case M_ERR_DGN_CONN_NOT_ANNOUNCED: s = "M_ERR_DGN_CONN_NOT_ANNOUNCED";break;
            case M_ERR_DGN_CONN_TOO_MUCH: s = "M_ERR_DGN_CONN_TOO_MUCH";break;
            case M_ERR_DGN_INFO_NOT_AVAIL: s = "M_ERR_DGN_INFO_NOT_AVAIL";break;
            case M_ERR_FEW_PROT_LEVEL: s = "M_ERR_FEW_PROT_LEVEL";break;
            case M_ERR_FILE_NOT_FOUND: s = "M_ERR_FILE_NOT_FOUND";break;
            case M_ERR_FW_ERROR: s = "M_ERR_FW_ERROR";break;
            case M_ERR_IEC_DATA_TYPE_MISMATCH: s = "M_ERR_IEC_DATA_TYPE_MISMATCH";break;
            case M_ERR_IEC_INVALID_REF: s = "M_ERR_IEC_INVALID_REF";break;
            case M_ERR_IEC_LOWER_LAYER: s = "M_ERR_IEC_LOWER_LAYER";break;
            case M_ERR_IEC_NEG_RESPONSE: s = "M_ERR_IEC_NEG_RESPONSE";break;
            case M_ERR_IEC_NO: s = "M_ERR_IEC_NO";break;
            case M_ERR_IEC_NO_ACCESS_TO_REM_OBJECT: s = "M_ERR_IEC_NO_ACCESS_TO_REM_OBJECT";break;
            case M_ERR_IEC_PARTNER_IN_WRONG_STATE: s = "M_ERR_IEC_PARTNER_IN_WRONG_STATE";break;
            case M_ERR_IEC_RECEIVER_DISABLED: s = "M_ERR_IEC_RECEIVER_DISABLED";break;
            case M_ERR_IEC_RECEIVER_OVERRUN: s = "M_ERR_IEC_RECEIVER_OVERRUN";break;
            case M_ERR_IEC_RESET_RECEIVED: s = "M_ERR_IEC_RESET_RECEIVED";break;
            case M_ERR_INSTALL: s = "M_ERR_INSTALL";break;
            case M_ERR_INTERNAL_ERROR: s = "M_ERR_INTERNAL_ERROR";break;
            case M_ERR_INVALID_BLOCK: s = "M_ERR_INVALID_BLOCK";break;
            case M_ERR_INVALID_BLOCK_NR: s = "M_ERR_INVALID_BLOCK_NR";break;
            case M_ERR_INVALID_BLOCK_TYPE: s = "M_ERR_INVALID_BLOCK_TYPE";break;
            case M_ERR_INVALID_CONN_STATE: s = "M_ERR_INVALID_CONN_STATE";break;
            case M_ERR_INVALID_CREF: s = "M_ERR_INVALID_CREF";break;
            case M_ERR_INVALID_CYCL_READ_STATE: s = "M_ERR_INVALID_CYCL_READ_STATE";break;
            case M_ERR_INVALID_DATA_SIZE: s = "M_ERR_INVALID_DATA_SIZE";break;
            case M_ERR_INVALID_DATARANGE_OR_TYPE: s = "M_ERR_INVALID_DATARANGE_OR_TYPE";break;
            case M_ERR_INVALID_FILENAME: s = "M_ERR_INVALID_FILENAME";break;
            case M_ERR_INVALID_FILETYPE: s = "M_ERR_INVALID_FILETYPE";break;
            case M_ERR_INVALID_ORDERIDM: s = "M_ERR_INVALID_ORDERIDM";break;
            case M_ERR_INVALID_PARAMETER: s = "M_ERR_INVALID_PARAMETER";break;
            case M_ERR_INVALID_PASSWORD: s = "M_ERR_INVALID_PASSWORD";break;
            case M_ERR_INVALID_READ_BUFFER: s = "M_ERR_INVALID_READ_BUFFER";break;
            case M_ERR_INVALID_SEGMENT: s = "M_ERR_INVALID_SEGMENT";break;
            case M_ERR_INVALID_WRITE_BUFFER: s = "M_ERR_INVALID_WRITE_BUFFER";break;
            case M_ERR_MAX_REQ: s = "M_ERR_MAX_REQ";break;
            case M_ERR_NO_PASSWORD: s = "M_ERR_NO_PASSWORD";break;
            case M_ERR_NO_RCV_BLOCK: s = "M_ERR_NO_RCV_BLOCK";break;
            case M_ERR_NO_RESOURCE: s = "M_ERR_NO_RESOURCE";break;
            case M_ERR_NO_SIN_SERV: s = "M_ERR_NO_SIN_SERV";break;
            case M_ERR_OBJ_ACCESS_DENIED: s = "M_ERR_OBJ_ACCESS_DENIED";break;
            case M_ERR_OBJ_ATTR_INCONSISTENT: s = "M_ERR_OBJ_ATTR_INCONSISTENT";break;
            case M_ERR_OBJ_UNDEFINED: s = "M_ERR_OBJ_UNDEFINED";break;
            case M_ERR_ORDERID_USED: s = "M_ERR_ORDERID_USED";break;
            case M_ERR_R_ID_USED: s = "M_ERR_R_ID_USED";break;
            case M_ERR_RECEIVE_BUFFER_FULL: s = "M_ERR_RECEIVE_BUFFER_FULL";break;
            case M_ERR_REM_BRCV: s = "M_ERR_REM_BRCV";break;
            case M_ERR_REM_BSEND: s = "M_ERR_REM_BSEND";break;
            case M_ERR_REM_BSEND_CANCEL: s = "M_ERR_REM_BSEND_CANCEL";break;
            case M_ERR_REM_DATABASE_TOO_SMALL: s = "M_ERR_REM_DATABASE_TOO_SMALL";break;
            case M_ERR_RETRY: s = "M_ERR_RETRY";break;
            case M_ERR_SERVICE_CONN_ALREADY_USED: s = "M_ERR_SERVICE_CONN_ALREADY_USED";break;
            case M_ERR_SERVICE_NOT_SUPPORTED: s = "M_ERR_SERVICE_NOT_SUPPORTED";break;
            case M_ERR_SERVICE_VFD_ALREADY_USED: s = "M_ERR_SERVICE_VFD_ALREADY_USED";break;
            case M_ERR_SYMB_ADDRESS: s = "M_ERR_SYMB_ADDRESS";break;
            case M_ERR_SYMB_ADDRESS_INCONSISTENT: s = "M_ERR_SYMB_ADDRESS_INCONSISTENT";break;
            case M_ERR_TOO_LONG_DATA: s = "M_ERR_TOO_LONG_DATA";break;
            case M_ERR_UNKNOWN_ERROR: s = "M_ERR_UNKNOWN_ERROR";break;
            case M_ERR_WRONG_CP_DESCR: s = "M_ERR_WRONG_CP_DESCR";break;
            case M_ERR_WRONG_IND_CNF: s = "M_ERR_WRONG_IND_CNF";break;

            case M_TYPE_IEC_BOOL: s = "M_TYPE_IEC_BOOL";break;
            case M_TYPE_IEC_SINT: s = "M_TYPE_IEC_SINT";break;
            case M_TYPE_IEC_INT: s = "M_TYPE_IEC_INT";break;
            case M_TYPE_IEC_DINT: s = "M_TYPE_IEC_DINT";break;
            case M_TYPE_IEC_LINT: s = "M_TYPE_IEC_LINT";break;
            case M_TYPE_IEC_USINT: s = "M_TYPE_IEC_USINT";break;
            case M_TYPE_IEC_UINT: s = "M_TYPE_IEC_UINT";break;
            case M_TYPE_IEC_UDINT: s = "M_TYPE_IEC_UDINT";break;
            case M_TYPE_IEC_ULINT: s = "M_TYPE_IEC_ULINT";break;
            case M_TYPE_IEC_REAL: s = "M_TYPE_IEC_REAL";break;
            case M_TYPE_IEC_LREAL: s = "M_TYPE_IEC_LREAL";break;
            case M_TYPE_IEC_TIME: s = "M_TYPE_IEC_TIME";break;
            case M_TYPE_IEC_DATE: s = "M_TYPE_IEC_DATE";break;
            case M_TYPE_IEC_TIME_OF_DAY: s = "M_TYPE_IEC_TIME_OF_DAY";break;
            case M_TYPE_IEC_TOD: s = "M_TYPE_IEC_TOD";break;
            case M_TYPE_IEC_DATE_AND_TIME: s = "M_TYPE_IEC_DATE_AND_TIME";break;
            case M_TYPE_IEC_DT: s = "M_TYPE_IEC_DT";break;
            case M_TYPE_IEC_STRING: s = "M_TYPE_IEC_STRING";break;
            case M_TYPE_IEC_BYTE: s = "M_TYPE_IEC_BYTE";break;
            case M_TYPE_IEC_WORD: s = "M_TYPE_IEC_WORD";break;
            case M_TYPE_IEC_DWORD: s = "M_TYPE_IEC_DWORD";break;
            case M_TYPE_IEC_LWORD: s = "M_TYPE_IEC_LWORD";break;
            case M_TYPE_IEC_WSTRING: s = "M_TYPE_IEC_WSTRING";break;

            default: s = "M_WITHOUT_DESCRIPTION";break;
      }
        return s;
    };
}
