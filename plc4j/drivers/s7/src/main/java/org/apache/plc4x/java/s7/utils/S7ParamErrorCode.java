/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.s7.utils;

import java.util.HashMap;
import java.util.Map;

public enum S7ParamErrorCode {

    ERROR_0x0000((short) 0x0000, "No error"),
    ERROR_0x0110((short) 0x0110, "Invalid block number"),
    ERROR_0x0111((short) 0x0111, "Invalid request length"),
    ERROR_0x0112((short) 0x0112, "Invalid parameter"),
    ERROR_0x0113((short) 0x0113, "Invalid block type"),
    ERROR_0x0114((short) 0x0114, "Block not found"),
    ERROR_0x0115((short) 0x0115, "Block already exists"),
    ERROR_0x0116((short) 0x0116, "Block is write-protected"),
    ERROR_0x0117((short) 0x0117, "The block/operating system update is too large"),
    ERROR_0x0118((short) 0x0118, "Invalid block number"),
    ERROR_0x0119((short) 0x0119, "Incorrect password entered"),
    ERROR_0x011A((short) 0x011A, "PG resource error"),
    ERROR_0x011B((short) 0x011B, "PLC resource error"),
    ERROR_0x011C((short) 0x011C, "Protocol error"),
    ERROR_0x011D((short) 0x011D, "Too many blocks (module-related restriction)"),
    ERROR_0x011E((short) 0x011E, "There is no longer a connection to the database, or S7DOS handle is invalid"),
    ERROR_0x011F((short) 0x011F, "Result buffer too small"),
    ERROR_0x0120((short) 0x0120, "End of block list"),
    ERROR_0x0140((short) 0x0140, "Insufficient memory available"),
    ERROR_0x0141((short) 0x0141, "Job cannot be processed because of a lack of resources"),
    ERROR_0x8001((short) 0x8001, "The requested service cannot be performed while the block is in the current status"),
    ERROR_0x8003((short) 0x8003, "S7 protocol error: Error occurred while transferring the block"),
    ERROR_0x8100((short) 0x8100, "Application, general error: Service unknown to remote module"),
    ERROR_0x8104((short) 0x8104, "This service is not implemented on the module or a frame error was reported"),
    ERROR_0x8204((short) 0x8204, "The type specification for the object is inconsistent"),
    ERROR_0x8205((short) 0x8205, "A copied block already exists and is not linked"),
    ERROR_0x8301((short) 0x8301, "Insufficient memory space or work memory on the module, or specified storage medium not accessible"),
    ERROR_0x8302((short) 0x8302, "Too few resources available or the processor resources are not available"),
    ERROR_0x8304((short) 0x8304, "No further parallel upload possible. There is a resource bottleneck"),
    ERROR_0x8305((short) 0x8305, "Function not available"),
    ERROR_0x8306((short) 0x8306, "Insufficient work memory (for copying, linking, loading AWP)"),
    ERROR_0x8307((short) 0x8307, "Not enough retentive work memory (for copying, linking, loading AWP)"),
    ERROR_0x8401((short) 0x8401, "S7 protocol error: Invalid service sequence (for example, loading or uploading a block)"),
    ERROR_0x8402((short) 0x8402, "Service cannot execute owing to status of the addressed object"),
    ERROR_0x8404((short) 0x8404, "S7 protocol: The function cannot be performed"),
    ERROR_0x8405((short) 0x8405, "Remote block is in DISABLE state (CFB). The function cannot be performed"),
    ERROR_0x8500((short) 0x8500, "S7 protocol error: Wrong frames"),
    ERROR_0x8503((short) 0x8503, "Alarm from the module: Service canceled prematurely"),
    ERROR_0x8701((short) 0x8701, "Error addressing the object on the communications partner (for example, area length error)"),
    ERROR_0x8702((short) 0x8702, "The requested service is not supported by the module"),
    ERROR_0x8703((short) 0x8703, "Access to object refused"),
    ERROR_0x8704((short) 0x8704, "Access error: Object damaged"),
    ERROR_0xD001((short) 0xD001, "Protocol error: Illegal job number"),
    ERROR_0xD002((short) 0xD002, "Parameter error: Illegal job variant"),
    ERROR_0xD003((short) 0xD003, "Parameter error: Debugging function not supported by module"),
    ERROR_0xD004((short) 0xD004, "Parameter error: Illegal job status"),
    ERROR_0xD005((short) 0xD005, "Parameter error: Illegal job termination"),
    ERROR_0xD006((short) 0xD006, "Parameter error: Illegal link disconnection ID"),
    ERROR_0xD007((short) 0xD007, "Parameter error: Illegal number of buffer elements"),
    ERROR_0xD008((short) 0xD008, "Parameter error: Illegal scan rate"),
    ERROR_0xD009((short) 0xD009, "Parameter error: Illegal number of executions"),
    ERROR_0xD00A((short) 0xD00A, "Parameter error: Illegal trigger event"),
    ERROR_0xD00B((short) 0xD00B, "Parameter error: Illegal trigger condition"),
    ERROR_0xD011((short) 0xD011, "Parameter error in path of the call environment: Block does not exist"),
    ERROR_0xD012((short) 0xD012, "Parameter error: Wrong address in block"),
    ERROR_0xD014((short) 0xD014, "Parameter error: Block being deleted/overwritten"),
    ERROR_0xD015((short) 0xD015, "Parameter error: Illegal tag address"),
    ERROR_0xD016((short) 0xD016, "Parameter error: Test jobs not possible, because of errors in user program"),
    ERROR_0xD017((short) 0xD017, "Parameter error: Illegal trigger number"),
    ERROR_0xD025((short) 0xD025, "Parameter error: Invalid path"),
    ERROR_0xD026((short) 0xD026, "Parameter error: Illegal access type"),
    ERROR_0xD027((short) 0xD027, "Parameter error: This number of data blocks is not permitted"),
    ERROR_0xD031((short) 0xD031, "Internal protocol error"),
    ERROR_0xD032((short) 0xD032, "Parameter error: Wrong result buffer length"),
    ERROR_0xD033((short) 0xD033, "Protocol error: Wrong job length"),
    ERROR_0xD03F((short) 0xD03F, "Coding error: Error in parameter section (for example, reserve bytes not equal to 0)"),
    ERROR_0xD041((short) 0xD041, "Data error: Illegal status list ID"),
    ERROR_0xD042((short) 0xD042, "Data error: Illegal tag address"),
    ERROR_0xD043((short) 0xD043, "Data error: Referenced job not found, check job data"),
    ERROR_0xD044((short) 0xD044, "Data error: Illegal tag value, check job data"),
    ERROR_0xD045((short) 0xD045, "Data error: Exiting the ODIS control is not allowed in HOLD"),
    ERROR_0xD046((short) 0xD046, "Data error: Illegal measuring stage during run-time measurement"),
    ERROR_0xD047((short) 0xD047, "Data error: Illegal hierarchy in 'Read job list'"),
    ERROR_0xD048((short) 0xD048, "Data error: Illegal deletion ID in 'Delete job'"),
    ERROR_0xD049((short) 0xD049, "Invalid substitute ID in 'Replace job'"),
    ERROR_0xD04A((short) 0xD04A, "Error executing 'program status'"),
    ERROR_0xD05F((short) 0xD05F, "Coding error: Error in data section (for example, reserve bytes not equal to 0, ...)"),
    ERROR_0xD061((short) 0xD061, "Resource error: No memory space for job"),
    ERROR_0xD062((short) 0xD062, "Resource error: Job list full"),
    ERROR_0xD063((short) 0xD063, "Resource error: Trigger event occupied"),
    ERROR_0xD064((short) 0xD064, "Resource error: Not enough memory space for one result buffer element"),
    ERROR_0xD065((short) 0xD065, "Resource error: Not enough memory space for several  result buffer elements"),
    ERROR_0xD066((short) 0xD066, "Resource error: The timer available for run-time measurement is occupied by another job"),
    ERROR_0xD067((short) 0xD067, "Resource error: Too many 'modify tag' jobs active (in particular multi-processor operation)"),
    ERROR_0xD081((short) 0xD081, "Function not permitted in current mode"),
    ERROR_0xD082((short) 0xD082, "Mode error: Cannot exit HOLD mode"),
    ERROR_0xD0A1((short) 0xD0A1, "Function not permitted in current protection level"),
    ERROR_0xD0A2((short) 0xD0A2, "Function not possible at present, because a function is running that modifies memory"),
    ERROR_0xD0A3((short) 0xD0A3, "Too many 'modify tag' jobs active on the I/O (in particular multi-processor operation)"),
    ERROR_0xD0A4((short) 0xD0A4, "Forcing' has already been established"),
    ERROR_0xD0A5((short) 0xD0A5, "Referenced job not found"),
    ERROR_0xD0A6((short) 0xD0A6, "Job cannot be disabled/enabled"),
    ERROR_0xD0A7((short) 0xD0A7, "Job cannot be deleted, for example because it is currently being read"),
    ERROR_0xD0A8((short) 0xD0A8, "Job cannot be replaced, for example because it is currently being read or deleted"),
    ERROR_0xD0A9((short) 0xD0A9, "Job cannot be read, for example because it is currently being deleted"),
    ERROR_0xD0AA((short) 0xD0AA, "Time limit exceeded in processing operation"),
    ERROR_0xD0AB((short) 0xD0AB, "Invalid job parameters in process operation"),
    ERROR_0xD0AC((short) 0xD0AC, "Invalid job data in process operation"),
    ERROR_0xD0AD((short) 0xD0AD, "Operating mode already set"),
    ERROR_0xD0AE((short) 0xD0AE, "The job was set up over a different connection and can only be handled over this connection"),
    ERROR_0xD0C1((short) 0xD0C1, "At least one error has been detected while accessing the tag(s)"),
    ERROR_0xD0C2((short) 0xD0C2, "Change to STOP/HOLD mode"),
    ERROR_0xD0C3((short) 0xD0C3, "At least one error was detected while accessing the tag(s). Mode change to STOP/HOLD"),
    ERROR_0xD0C4((short) 0xD0C4, "Timeout during run-time measurement"),
    ERROR_0xD0C5((short) 0xD0C5, "Display of block stack inconsistent, because blocks were deleted/reloaded"),
    ERROR_0xD0C6((short) 0xD0C6, "Job was automatically deleted as the jobs it referenced have been deleted"),
    ERROR_0xD0C7((short) 0xD0C7, "The job was automatically deleted because STOP mode was exited"),
    ERROR_0xD0C8((short) 0xD0C8, "Block status' aborted because of inconsistencies between test job and running program"),
    ERROR_0xD0C9((short) 0xD0C9, "Exit the status area by resetting OB90"),
    ERROR_0xD0CA((short) 0xD0CA, "Exiting the status range by resetting OB90 and access error reading tags before exiting"),
    ERROR_0xD0CB((short) 0xD0CB, "The output disable for the peripheral outputs has been activated again"),
    ERROR_0xD0CC((short) 0xD0CC, "The amount of data for the debugging functions is restricted by the time limit"),
    ERROR_0xD201((short) 0xD201, "Syntax error in block name"),
    ERROR_0xD202((short) 0xD202, "Syntax error in function parameters"),
    ERROR_0xD205((short) 0xD205, "Linked block already exists in RAM: Conditional copying is not possible"),
    ERROR_0xD206((short) 0xD206, "Linked block already exists in EPROM: Conditional copying is not possible"),
    ERROR_0xD208((short) 0xD208, "Maximum number of copied (not linked) blocks on module exceeded"),
    ERROR_0xD209((short) 0xD209, "(At least) one of the given blocks not found on the module"),
    ERROR_0xD20A((short) 0xD20A, "The maximum number of blocks that can be linked with one job was exceeded"),
    ERROR_0xD20B((short) 0xD20B, "The maximum number of blocks that can be deleted with one job was exceeded"),
    ERROR_0xD20C((short) 0xD20C, "OB cannot be copied because the associated priority class does not exist"),
    ERROR_0xD20D((short) 0xD20D, "SDB cannot be interpreted (for example, unknown number)"),
    ERROR_0xD20E((short) 0xD20E, "No (further) block available"),
    ERROR_0xD20F((short) 0xD20F, "Module-specific maximum block size exceeded"),
    ERROR_0xD210((short) 0xD210, "Invalid block number"),
    ERROR_0xD212((short) 0xD212, "Incorrect header attribute (run-time relevant)"),
    ERROR_0xD213((short) 0xD213, "Too many SDBs. Note the restrictions on the module being used"),
    ERROR_0xD216((short) 0xD216, "Invalid user program - reset module"),
    ERROR_0xD217((short) 0xD217, "Protection level specified in module properties not permitted"),
    ERROR_0xD218((short) 0xD218, "Incorrect attribute (active/passive)"),
    ERROR_0xD219((short) 0xD219, "Incorrect block lengths (for example, incorrect length of first section or of the whole block)"),
    ERROR_0xD21A((short) 0xD21A, "Incorrect local data length or write-protection code faulty"),
    ERROR_0xD21B((short) 0xD21B, "Module cannot compress or compression was interrupted early"),
    ERROR_0xD21D((short) 0xD21D, "The volume of dynamic project data transferred is illegal"),
    ERROR_0xD21E((short) 0xD21E, "Unable to assign parameters to a module (such as FM, CP). The system data could not be linked"),
    ERROR_0xD220((short) 0xD220, "Invalid programming language. Note the restrictions on the module being used"),
    ERROR_0xD221((short) 0xD221, "The system data for connections or routing are not valid"),
    ERROR_0xD222((short) 0xD222, "The system data of the global data definition contain invalid parameters"),
    ERROR_0xD223((short) 0xD223, "Error in instance data block for communication function block or maximum number of instance DBs exceeded"),
    ERROR_0xD224((short) 0xD224, "The SCAN system data block contains invalid parameters"),
    ERROR_0xD225((short) 0xD225, "The DP system data block contains invalid parameters"),
    ERROR_0xD226((short) 0xD226, "A structural error occurred in a block"),
    ERROR_0xD230((short) 0xD230, "A structural error occurred in a block"),
    ERROR_0xD231((short) 0xD231, "At least one loaded OB cannot be copied because the associated priority class does not exist"),
    ERROR_0xD232((short) 0xD232, "At least one block number of a loaded block is illegal"),
    ERROR_0xD234((short) 0xD234, "Block exists twice in the specified memory medium or in the job"),
    ERROR_0xD235((short) 0xD235, "The block contains an incorrect checksum"),
    ERROR_0xD236((short) 0xD236, "The block does not contain a checksum"),
    ERROR_0xD237((short) 0xD237, "You are about to load the block twice, i.e. a block with the same time stamp already exists on the CPU"),
    ERROR_0xD238((short) 0xD238, "At least one of the blocks specified is not a DB"),
    ERROR_0xD239((short) 0xD239, "At least one of the DBs specified is not available as a linked variant in the load memory"),
    ERROR_0xD23A((short) 0xD23A, "At least one of the specified DBs is considerably different from the copied and linked variant"),
    ERROR_0xD240((short) 0xD240, "Coordination rules violated"),
    ERROR_0xD241((short) 0xD241, "The function is not permitted in the current protection level"),
    ERROR_0xD242((short) 0xD242, "Protection violation while processing F blocks"),
    ERROR_0xD250((short) 0xD250, "Update and module ID or version do not match"),
    ERROR_0xD251((short) 0xD251, "Incorrect sequence of operating system components"),
    ERROR_0xD252((short) 0xD252, "Checksum error"),
    ERROR_0xD253((short) 0xD253, "No executable loader available; update only possible using a memory card"),
    ERROR_0xD254((short) 0xD254, "Storage error in operating system"),
    ERROR_0xD280((short) 0xD280, "Error compiling block in S7-300 CPU"),
    ERROR_0xD2A1((short) 0xD2A1, "Another block function or a trigger on a block is active"),
    ERROR_0xD2A2((short) 0xD2A2, "A trigger is active on a block. Complete the debugging function first"),
    ERROR_0xD2A3((short) 0xD2A3, "The block is not active (linked), the block is occupied or the block is currently marked for deletion"),
    ERROR_0xD2A4((short) 0xD2A4, "The block is already being processed by another block function"),
    ERROR_0xD2A6((short) 0xD2A6, "It is not possible to save and change the user program simultaneously"),
    ERROR_0xD2A7((short) 0xD2A7, "The block has the attribute 'unlinked' or is not processed"),
    ERROR_0xD2A8((short) 0xD2A8, "An active debugging function is preventing parameters from being assigned to the CPU"),
    ERROR_0xD2A9((short) 0xD2A9, "New parameters are being assigned to the CPU"),
    ERROR_0xD2AA((short) 0xD2AA, "New parameters are currently being assigned to the modules"),
    ERROR_0xD2AB((short) 0xD2AB, "The dynamic configuration limits are currently being changed"),
    ERROR_0xD2AC((short) 0xD2AC, "A running active or deactivate assignment (SFC 12) is temporarily preventing R-KiR process"),
    ERROR_0xD2B0((short) 0xD2B0, "An error occurred while configuring in RUN (CiR)"),
    ERROR_0xD2C0((short) 0xD2C0, "The maximum number of technological objects has been exceeded"),
    ERROR_0xD2C1((short) 0xD2C1, "The same technology data block already exists on the module"),
    ERROR_0xD2C2((short) 0xD2C2, "Downloading the user program or downloading the hardware configuration is not possible"),
    ERROR_0xD401((short) 0xD401, "Information function unavailable"),
    ERROR_0xD402((short) 0xD402, "Information function unavailable"),
    ERROR_0xD403((short) 0xD403, "Service has already been logged on/off (Diagnostics/PMC)"),
    ERROR_0xD404((short) 0xD404, "Maximum number of nodes reached. No more logons possible for diagnostics/PMC"),
    ERROR_0xD405((short) 0xD405, "Service not supported or syntax error in function parameters"),
    ERROR_0xD406((short) 0xD406, "Required information currently unavailable"),
    ERROR_0xD407((short) 0xD407, "Diagnostics error occurred"),
    ERROR_0xD408((short) 0xD408, "Update aborted"),
    ERROR_0xD409((short) 0xD409, "Error on DP bus"),
    ERROR_0xD601((short) 0xD601, "Syntax error in function parameter"),
    ERROR_0xD602((short) 0xD602, "Incorrect password entered"),
    ERROR_0xD603((short) 0xD603, "The connection has already been legitimized"),
    ERROR_0xD604((short) 0xD604, "The connection has already been enabled"),
    ERROR_0xD605((short) 0xD605, "Legitimization not possible because password does not exist"),
    ERROR_0xD801((short) 0xD801, "At least one tag address is invalid"),
    ERROR_0xD802((short) 0xD802, "Specified job does not exist"),
    ERROR_0xD803((short) 0xD803, "Illegal job status"),
    ERROR_0xD804((short) 0xD804, "Illegal cycle time (illegal time base or multiple)"),
    ERROR_0xD805((short) 0xD805, "No more cyclic read jobs can be set up"),
    ERROR_0xD806((short) 0xD806, "The referenced job is in a state in which the requested function cannot be performed"),
    ERROR_0xD807((short) 0xD807, "Function aborted due to overload, meaning executing the read cycle takes longer than the set scan cycle time"),
    ERROR_0xDC01((short) 0xDC01, "Date and/or time invalid"),
    ERROR_0xE201((short) 0xE201, "CPU is already the master"),
    ERROR_0xE202((short) 0xE202, "Connect and update not possible due to different user program in flash module"),
    ERROR_0xE203((short) 0xE203, "Connect and update not possible due to different firmware"),
    ERROR_0xE204((short) 0xE204, "Connect and update not possible due to different memory configuration"),
    ERROR_0xE205((short) 0xE205, "Connect/update aborted due to synchronization error"),
    ERROR_0xE206((short) 0xE206, "Connect/update denied due to coordination violation"),
    ERROR_0xEF01((short) 0xEF01, "S7 protocol error: Error at ID2; only 00H permitted in job"),
    ERROR_0xEF02((short) 0xEF02, "S7 protocol error: Error at ID2; set of resources does not exist"),;
    
    private static final Map<Short, S7ParamErrorCode> map;
    
    static {
        map = new HashMap<>();
        for (S7ParamErrorCode  subevent : S7ParamErrorCode .values()) {
            map.put(subevent.code, subevent);
        }
    }    
    
    private final String event;
    private final short code;
    
    S7ParamErrorCode(short code, String event){
        this.event = event;
        this.code = code;
    }
    
    public String getEvent(){
        return event;
    }    
    
    public short getCode() {
        return code;
    }    
    
    public static S7ParamErrorCode  valueOfEvent(String event) {
        for (S7ParamErrorCode  value : S7ParamErrorCode .values()) {
            if(value.getEvent().equals(event)) {
                return value;
            }
        }
        return null;
    }
    public static S7ParamErrorCode  valueOf(short code) {
        return map.get(code);
    }    
}