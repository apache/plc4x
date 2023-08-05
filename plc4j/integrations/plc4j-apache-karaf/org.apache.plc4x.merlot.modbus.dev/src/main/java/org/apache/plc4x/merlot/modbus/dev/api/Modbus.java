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
package org.apache.plc4x.merlot.modbus.dev.api;


public class Modbus {

    /*********************************************************************
     * MODBUS fuction codes.
     *********************************************************************/

    /**
     * Function 01 (01hex) Read Coils
     * Reads the ON/OFF status of discrete coils in the slave.
     */
    public static final byte FC_READ_COILS  = (byte) 0x01;

    /**
     * Function 02(02hex) Read Discrete Inputs
     * Reads the ON/OFF status of discrete inputs in the slave.
     */
    public static final byte FC_READ_DISCRETE_INPUTS        = (byte) 0x02;

    /**
     * Function 03 (03hex) Read Holding Registers
     * Read the binary contents of holding registers in the slave.
     */
    public static final byte FC_READ_HOLDING_REGISTERS      = (byte) 0x03;

    /**
     * Function 04 (04hex) Read Input Registers
     * Read the binary contents of input registers in the slave.
     */
    public static final byte FC_READ_INPUT_REGISTERS        = (byte) 0x04;

    /**
     * Function 05 (05hex) Write Single Coil
     * Writes a single coil to either ON or OFF.
     */
    public static final byte FC_WRITE_SINGLE_COIL           = (byte) 0x05;

    /**
     * Function 06 (06hex) Write Single Register
     * Writes a value into a single holding register.
     */
    public static final byte FC_WRITE_SINGLE_REGISTER       = (byte) 0x06;

    /**
     * Function code  07 (0x07), Read exception state. (only serial)
     */
    public static final byte FC_READ_EXCEPTION_STATUS       = (byte) 0x07;

    /**
     * Function code  08 (0x08), Diagnostic (only serial)
     */
    public static final byte FC_DIAGNOSTIC                  = (byte) 0x08;

    /**
     * Function code  11 (0x0B), Read the communication event counter (only serial)
     */
    public static final byte FC_GET_COMM_EVENT_COUNTER      = (byte) 0x0B;

    /**
     * Function code  12 (0x0C), Read the event logs (only serial)
     */
    public static final byte FC_GET_COMM_EVENT_LOG          = (byte) 0x0C;

    /**
     * Function 15 (0Fhex) Write Multiple Coils
     * Writes each coil in a sequence of coils to either ON or OFF.
     */
    public static final byte FC_WRITE_MULTIPLE_COILS        = (byte) 0x0F;
    
    /**
     * Function 16 (10hex) Write Multiple Registers
     * Writes values into a sequence of holding registers
     */
    public static final byte FC_WRITE_MULTIPLE_REGISTERS    = (byte) 0x10;

    /**
     * Function code  17 (0x11), Report the slave ID (only serial)
     */
    public static final byte FC_REPORT_SLAVE_ID             = (byte) 0x11;
    
    /**
     * Function code  20 (0x14), Read one record from file
     */
    public static final byte FC_READ_FILE_RECORD            = (byte) 0x14;

    /**
     * Function code  21 (0x15), Write one record to file
     */
    public static final byte FC_WRITE_FILE_RECORD           = (byte) 0x15;

    /**
     * Function code  22 (0x16), enmascara escribir registro
     */
    public static final byte FC_MASK_WRITE_REGISTER         = (byte) 0x16;

    /**
     * Function code  23 (0x17), Read/Write multiple registers
     */
    public static final byte FC_READ_WRITE_MULTIPLE_REGISTERS = (byte) 0x17;

    /**
     * Function code  24 (0x18), Read FIFO QUEUE
     */
    public static final byte FC_READ_FIFO_QUEUE = (byte) 0x18;

    /**
     * Command to put off one coil.
     */
    public static final byte COIL_OFF                 = (byte) 0x00;

    /**
     * Command to put on one coil.
     */
    public static final byte COIL_ON                  = (byte) 0xFF;

    /*********************************************************************
     * MODBUS exception code.
     *********************************************************************/

    /**
     * The function code received in the query is not an
     * allowable action for the server. This may be
     * because the function code is only applicable to
     * newer devices, and was not implemented in the
     * unit selected. It could also indicate that the server
     * is in the wrong state to process a request of this
     * type, for example because it is unconfigure d and
     * is being asked to return register values.
     */
    public static final byte EX_ILLEGAL_FUNCTION            = (byte) 0x01;

    /**
     * The data address received in the query is not an
     * allowable address for the server. More
     * specifically, the combination of reference number
     * and transfer length is invalid. For a controller with
     * 100 registers, the PDU addresses the first register
     * as 0, and the last one as 99. If a request is
     * submitted with a starting register address of 96
     * and a quantity of registers of 4, then this request
     * will successfully operate (address -wise at least)
     * on registers 96, 97, 98, 99. If a request is
     * submitted with a starting register address of 96
     * and a quantity of registers of 5, then this request
     * will fail with Exception Code 0x02 �Illegal Data
     * Address� since it attempts to operate on registers
     * 96, 97, 98, 99 and 100, and there is no register
     * with address 100.
     */
    public static final byte EX_ILLEGAL_DATA_ADDRESS        = (byte) 0x02;

    /**
     * A value contained in the query data field is not an
     * allowable value for server. This indicates a fault in
     * the structure of the remainder of a complex
     * request, such as that th e implied length is
     * incorrect. It specifically does NOT mean that a
     * data item submitted for storage in a register has a
     * value outside the expectation of the application
     * program, since the MODBUS protocol is unaware
     * of the significance of any particular val ue of any
     * particular register.
     */
    public static final byte EX_ILLEGAL_DATA_VALUE          = (byte) 0x03;

    /**
     * An unrecoverable error occurred while the server
     * was attempting to perform the requested action.
     */
    public static final byte EX_SERVER_DEVICE_FAILURE        = (byte) 0x04;

    /**
     * Specialized use in conjunction with programming
     * commands.
     * The server has accepted the request and is
     * processing it, but a long duration of time will be
     * required to do so. This response is returned to
     * prevent a timeout error from occurring in the
     * client. The client can next issue a Poll Program
     * Complete message to determine if processing is
     * completed.
     */
    public static final byte EX_ACKNOWLEDGE                 = (byte) 0x05;

    /**
     * Specialized use in conjunction with programming
     * commands.
     * The server is engaged in processing a long �
     * duration program command. The client should
     * retransmit the message later when the server is
     * free.
     */
    public static final byte EX_SERVER_DEVICE_BUSY           = (byte) 0x06;

    /**
     * Specialized use in conjunction with function codes
     * 20 and 21 and reference type 6, to indicate that
     * the extended file area failed to pass a consistency
     * check.
     * The server attempted to read record file, but
     * detected a parity error in the memory. The client
     * can retry the request, but service may be required
     * on the server device.
     */
    public static final byte EX_MEMORY_PARITY_ERROR          = (byte) 0x08;

    /**
     * Specialized use in conjunction with gateways,
     * indicates that the gateway was unable to allocate
     * an internal communication path from the input port
     * to the output port for processing the request.
     * Usually means that the gateway is misconfigured
     * or overloaded.
     */
    public static final byte EX_GATEWAY_PATH_UNAVAILABLE     = (byte) 0x0A;

    /**
     * Specialized use in conjunction with gateways,
     * indicates that no response was obtained from the
     * target device. Usually means that the device is not
     * present on the network.
     */
    public static final byte EX_GATEWAY_TARGET_DEVICE_FAILED_TO_RESPOND  = (byte) 0x0A;

    /*********************************************************************
     * MODBUS code of undocumented exceptions.
     *********************************************************************/

    /**
     * Exception code, significa que la petición debio haber generado
     * una respuesta larga ilegal.
     */
    public static final byte EX_ILLEGAL_RESPONSE_LENGTH  = (byte) 0x04;

    /**
     * Exception code para el modificador de excepciones.
     * Este valor es agregado al valor de la funcíon para indicar
     * que ocurrio un error.
     */
    public static final byte EX_EXCEPTION_MODIFIER       = (byte) 0x80;

     /*********************************************************************
     * Parameterization of the MODBUS driver.
     *********************************************************************/

    /**
     * Highest permissible address.
     */
    public static final int ADDRESS_MAX            = 65535;

    /**
     * Valor de mensaje mas largo..
     */
    public static final int MAX_MESSAGE_LENGTH = 256;

    /**
     * Máximo valor de un entero sin signo.
     */
    public static int UINT16_MAX = (int) 0xFFFF;

    /**
     * Valor minimo para un entero sin signo.
     */
    public static int UINT16_MIN = (int) 0x0000;

    /**
     * Máximo valor para un entero de 8 bits sin signo.
     */
    public static int UINT8_MAX = (int) 0xFF;

    /**
     * M�nimo valor para un entero de 8 bits sin signo.
     */
    public static int UINT8_MIN = (int) 0x00;

    /**
    * N�mero de registros m�ximos.
    */
    public static int MAX_REGISTERS = (int) 0xFFFF;

    /**
    * N�mero de bytes requeridos para el dispositivo.
    */
    public static long MAX_REGISTER_BYTES = (long) 0x20000;
}
