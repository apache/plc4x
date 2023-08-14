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
package org.apache.plc4x.merlot.modbus.svr.core;

import org.apache.plc4x.merlot.modbus.dev.api.Modbus;
import org.apache.plc4x.merlot.modbus.dev.api.ModbusDevice;
import org.apache.plc4x.merlot.modbus.dev.api.ModbusDeviceArray;
import org.apache.plc4x.merlot.modbus.dev.impl.ModbusDeviceArrayImpl;
import org.apache.plc4x.merlot.modbus.svr.api.ModbusADU;
import org.apache.plc4x.merlot.modbus.svr.impl.ModbusADUImpl;
import io.netty.buffer.ByteBuf;
import io.netty.buffer.ByteBufUtil;
import io.netty.buffer.Unpooled;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.SimpleChannelInboundHandler;
import org.osgi.framework.BundleContext;
import org.osgi.framework.ServiceReference;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class ModbusServerADUHandler extends SimpleChannelInboundHandler<ModbusADU>{
    private final Logger LOGGER = LoggerFactory.getLogger(ModbusServerADUHandler.class.getName());

    private BundleContext bc;
    private ModbusDeviceArray myMda = null;
	
    public ModbusServerADUHandler(BundleContext bc){
    	super();
    	this.bc = bc;
    	ServiceReference<?> sr = bc.getServiceReference("org.apache.plc4x.merlot.modbus.dev.api.ModbusDeviceArray");
    	if (sr!=null) {
    		myMda = (ModbusDeviceArray) bc.getService(sr);
    	}
    }

    @Override
    protected void channelRead0(ChannelHandlerContext ctx, ModbusADU rxADU) throws Exception {
        ModbusADU txADU = null;
        if (myMda != null) {
            txADU = doFunctionCode(rxADU);
        } else {
            ServiceReference<ModbusDeviceArrayImpl> sr = bc.getServiceReference(ModbusDeviceArrayImpl.class);
            if (sr != null) {
                myMda = (ModbusDeviceArray) bc.getService(sr);
            }
        }
        if (txADU != null) {
            if (txADU.getFunctionCode() != 0) {
                ctx.writeAndFlush(txADU);
            }
        }
    }

    @Override
    public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) throws Exception {
        //cause.printStackTrace();
        super.exceptionCaught(ctx, cause);
    }
	
    @SuppressWarnings("unused")
    public ModbusADU doFunctionCode(ModbusADU rxADU){

        ModbusADU txADU = null;

        if (myMda.getModbusDevicesArray()[rxADU.getUnitID()] != null) {

            ByteBuf cbuffer = null;
            int startaddress;
            int length;
            int address;
            int value;
            int i;
            byte functioncode ;
            byte bytecount;
            LOGGER.trace(">: " + ByteBufUtil.hexDump(rxADU.getData()));
            switch(rxADU.getFunctionCode()){

                case Modbus.FC_READ_DISCRETE_INPUTS: //Read Discrete Inputs
                        txADU = this.doFC_Read_Discrete_Inputs(rxADU);  					
                        break;

                case Modbus.FC_READ_COILS: //Read Coils
                        txADU = this.doFC_Read_Coils(rxADU);	
                        break;

                case Modbus.FC_WRITE_SINGLE_COIL: //Write Single Coil
                        txADU = this.doFC_Write_Single_Coil(rxADU);						
                        break;	

                case Modbus.FC_WRITE_MULTIPLE_COILS: //Write Multiple Coils
                        txADU = this.doFC_Write_Multiple_Coils(rxADU);					
                        break;	

                case Modbus.FC_READ_INPUT_REGISTERS: //Read Input Register
                        txADU = this.doFC_Read_Input_Registers(rxADU);					
                        break;

                case Modbus.FC_READ_HOLDING_REGISTERS: //Read Holding Registers                        
                        txADU = this.doFC_Read_Holding_Registers(rxADU);					
                        break;	

                case Modbus.FC_WRITE_SINGLE_REGISTER: //Write Single Register
                        txADU = this.doFC_Write_Single_Register(rxADU);					
                        break;	

                case Modbus.FC_WRITE_MULTIPLE_REGISTERS: //Write Multiple Registers
                        txADU = this.doFC_Write_Multiple_Registers(rxADU);					
                        break;

                case Modbus.FC_READ_WRITE_MULTIPLE_REGISTERS: //Read/Write Multiple Registers
                        txADU = this.doFC_Read_Write_Multiple_Registers(rxADU);		
                        break;

                case Modbus.FC_MASK_WRITE_REGISTER: //Mask Write Register
                        txADU = this.doFC_Mask_Write_Register(rxADU);						
                        break;

                case Modbus.FC_READ_FIFO_QUEUE: //Read FIFO queue
                        break;	

                case Modbus.FC_READ_FILE_RECORD: //Read File Record
                        break;	

                case Modbus.FC_WRITE_FILE_RECORD: //Write File Record
                        break;

                case Modbus.FC_READ_EXCEPTION_STATUS: //Read Exception status
                        break;	

                case Modbus.FC_DIAGNOSTIC: //Diagnostic
                        break;	

                case Modbus.FC_GET_COMM_EVENT_COUNTER: //Get Com event counter
                        break;	

                case Modbus.FC_GET_COMM_EVENT_LOG: //Get Com event Log
                        break;

                case Modbus.FC_REPORT_SLAVE_ID: //Report server ID
                        break;	

                case 0x2B: //Read device Identification
                        break;	

                default:
                        break;
            }
            LOGGER.trace("<: " + ByteBufUtil.hexDump(txADU.getData()));
            rxADU.getData().release();
            return txADU;
        }
        rxADU.getData().release();
        return null;
    }

    /* (non-Javadoc)
     * @see io.netty.channel.ChannelInboundHandlerAdapter#channelInactive(io.netty.channel.ChannelHandlerContext)
     */
    @Override
    public void channelInactive(ChannelHandlerContext ctx) throws Exception {
            // TODO Auto-generated method stub
            super.channelInactive(ctx);
    }
	
    /**
     * This function code is used to read from 1 to 2000 contiguous status of discrete inputs in a
     * remote device. The Request PDU specifies the starting address, i.e. the address of the first
     * input specified, and the number of inputs. In the PDU Discrete Inputs are addressed starting
     * at zero. Therefore Discrete inputs numbered 1 -16 are addressed as 0-15.
     * The discrete inputs in the response message are packed as one input per bit of the data field.
     * Status is indicated as 1= ON; 0= OFF. The LSB of the first data byte contains the input
     * addressed in the query. The other inputs follow toward the high order end of this byte, and
     * from low order to high order in subsequent bytes.
     * If the returned input quantity is not a multiple of eight, the remaining bits in the final data byte
     * will be padded with zeros (toward the high order end of the byte). The Byte Count field
     * specifies the quantity of complete bytes of data.
     * 
     * @param rxADU
     * @return txADU
     */
    private ModbusADU doFC_Read_Discrete_Inputs(ModbusADU rxADU){
        ModbusADU txADU = new ModbusADUImpl();

        int Starting_address;
        int Quantity_of_coils;
        int Byte_count;
        ByteBuf DiscreteInputs;

        txADU.setTransactionID(rxADU.getTransactionID());
        txADU.setProtocolID(rxADU.getProtocolID());
        txADU.setUnitID(rxADU.getUnitID());
        txADU.setFunctionCode(Modbus.FC_READ_DISCRETE_INPUTS);

        DiscreteInputs = myMda.getModbusDevicesArray()[rxADU.getUnitID()].getDiscreteInputs();
        Starting_address = rxADU.getData().readShort();
        Quantity_of_coils = rxADU.getData().readShort();

        Byte_count = (byte) ((Quantity_of_coils / 8) + (((Quantity_of_coils % 8) == 0) ? 0 : 1));
        ByteBuf txData = Unpooled.buffer(Byte_count + 1);
        txADU.setData(txData);	

        if ((Quantity_of_coils >= 0x0001) && (Quantity_of_coils <= 0x07D0)) {

            if ((Starting_address + Quantity_of_coils) <= (DiscreteInputs.capacity() * 8)) {
                try {
                    txData.writeByte(Byte_count);
                    DiscreteInputs.getBytes(Starting_address, txADU.getData(), Byte_count);
                    Byte_count += 2;
                    txADU.setLengthField((short) (Byte_count + 1));
                } catch (Exception ex) {
                    txADU.setFunctionCode((byte) (Modbus.FC_READ_DISCRETE_INPUTS + Modbus.EX_EXCEPTION_MODIFIER));
                    txADU.getData().clear();
                    txADU.getData().writeByte(Modbus.EX_SERVER_DEVICE_FAILURE);
                }
            } else {
                txADU.setFunctionCode((byte) (Modbus.FC_READ_DISCRETE_INPUTS + Modbus.EX_EXCEPTION_MODIFIER));
                txADU.getData().writeByte(Modbus.EX_ILLEGAL_DATA_ADDRESS);
            };
        } else {
            txADU.setFunctionCode((byte) (Modbus.FC_READ_DISCRETE_INPUTS + Modbus.EX_EXCEPTION_MODIFIER));
            txADU.getData().writeByte(Modbus.EX_ILLEGAL_DATA_VALUE);
        };

        return txADU;
    }

    /**
     * This function code is used to read from 1 to 2000 contiguous status of coils in a remote
     * device. The Request PDU specifies the starting address, i.e. the address of the first coil
     * specified, and the number of coils. In the PDU Coils are addressed starting at zero. Therefore
     * coils numbered 1 -16 are addressed as 0-15.
     * The coils in the response message are packed as one coil per bit of the data field. Status is
     * indicated as 1= ON and 0= OFF. The LSB of the first data byte contains the output addressed
     * in the query. The other coils follow toward the high order end of this byte, and from low order
     * to high order in subsequent bytes.
     * If the returned output quantity is not a multiple of eight, the remaining bits in the final data
     * byte will be padded with zeros (toward the high order end of the byte). The Byte Count field
     * specifies the quantity of complete bytes of data.
     * 
     * @param rxADU
     * @return txADU 
     */
    private ModbusADU doFC_Read_Coils(ModbusADU rxADU){
        ModbusADU txADU = new ModbusADUImpl();

        int Starting_address;
        int Quantity_of_coils;
        int Byte_count;
        ByteBuf Coils;

        txADU.setTransactionID(rxADU.getTransactionID());
        txADU.setProtocolID(rxADU.getProtocolID());
        txADU.setUnitID(rxADU.getUnitID());
        txADU.setFunctionCode(Modbus.FC_READ_COILS);

        Coils = myMda.getModbusDevicesArray()[rxADU.getUnitID()].getCoils();

        txADU.setFunctionCode(Modbus.FC_READ_COILS);
        Starting_address = rxADU.getData().readShort();
        Quantity_of_coils = rxADU.getData().readShort();

        Byte_count =  ((Quantity_of_coils / 8) + (((Quantity_of_coils % 8) == 0) ? 0 : 1));
        ByteBuf txData = Unpooled.buffer(Byte_count + 1);
        txADU.setData(txData);

        if ((Quantity_of_coils >= 0x0001) && (Quantity_of_coils <= 0x07D0)) {

            if ((Starting_address + Quantity_of_coils) <= (Coils.capacity() * 8)) {
                try {
                    txData.writeByte(Byte_count);
                    Coils.getBytes(Starting_address, txADU.getData(), Byte_count);
                    Byte_count += 2;
                    txADU.setLengthField((short) (Byte_count + 1));
                } catch (Exception ex) {
                    txADU.setFunctionCode((byte) (Modbus.FC_READ_COILS + Modbus.EX_EXCEPTION_MODIFIER));
                    txADU.getData().clear();
                    txADU.getData().writeByte(Modbus.EX_SERVER_DEVICE_FAILURE);
                }
            } else {
                txADU.setFunctionCode((byte) (Modbus.FC_READ_COILS + Modbus.EX_EXCEPTION_MODIFIER));
                txADU.getData().writeByte(Modbus.EX_ILLEGAL_DATA_ADDRESS);
            };
        } else {
            txADU.setFunctionCode((byte) (Modbus.FC_READ_COILS + Modbus.EX_EXCEPTION_MODIFIER));
            txADU.getData().writeByte(Modbus.EX_ILLEGAL_DATA_VALUE);
        };
        return txADU;
    }
	
    /**
     * This function code is used to write a single output to either ON or OFF in a remote device.
     * The requested ON/OFF state is specified by a constant in the request data field. A value of
     * FF00 hex requests the output to be ON. A value of 00 00 requests it to be OFF. All other
     * values are illegal and will not affect the output.
     * The Request PDU specifi es the address of the coil to be forced. Coils are addressed
     * starting at zero. Therefore coil numbered 1 is addressed as 0. The requested ON/OFF state
     * is specified by a constant in the Coil Value field. A value of 0XFF00 requests the coil to
     * be ON. A value of 0X0000 requests the coil to be off. All other values are illegal and will
     * not affect the coil.
     * 
     * @param rxADU
     * @return txADU
     */
    @SuppressWarnings("unused")	
    private ModbusADU doFC_Write_Single_Coil(ModbusADU rxADU){

        ModbusADU txADU = new ModbusADUImpl();

        int Starting_address;

        int Quantity_of_coils;
        int Byte_count;
        short Value;
        ByteBuf Coils;

        txADU.setTransactionID(rxADU.getTransactionID());
        txADU.setProtocolID(rxADU.getProtocolID());
        txADU.setUnitID(rxADU.getUnitID());
        txADU.setFunctionCode(Modbus.FC_WRITE_SINGLE_COIL);

        Starting_address = rxADU.getData().readShort();
        Value = rxADU.getData().readShort();
        ByteBuf txData = Unpooled.buffer(2);
        txADU.setData(txData);
        Coils = myMda.getModbusDevicesArray()[rxADU.getUnitID()].getCoils();
       
        if ((Value == (short) 0x0000) || (Value == (short) 0xFF00)) {
            System.out.println("Es cero o uno...");
            if (Starting_address <= Coils.capacity() * 8) {
                try {
                    System.out.println("N�mero de bobinas aceptada!.");
                    myMda.getModbusDevicesArray()[rxADU.getUnitID()].
                            setCoil(Starting_address, ((Value == (short) 0xFF00) ? false : true));
                    txADU.getData().writeShort(Starting_address);
                    //Unit + Function + Data todo en bytes.
                    txADU.setLengthField((short) 2);

                } catch (Exception ex) {
                    txADU.setFunctionCode((byte) (Modbus.FC_WRITE_SINGLE_COIL + Modbus.EX_EXCEPTION_MODIFIER));
                    txADU.getData().clear();
                    txADU.getData().writeByte(Modbus.EX_SERVER_DEVICE_FAILURE);
                }
            } else {
                txADU.setFunctionCode((byte) (Modbus.FC_WRITE_SINGLE_COIL + Modbus.EX_EXCEPTION_MODIFIER));
                txADU.getData().writeByte(Modbus.EX_ILLEGAL_DATA_ADDRESS);
            }

        } else {
            txADU.setFunctionCode((byte) (Modbus.FC_WRITE_SINGLE_COIL + Modbus.EX_EXCEPTION_MODIFIER));
            txADU.getData().writeByte(Modbus.EX_ILLEGAL_DATA_VALUE);
        }
        return txADU;
    }	
	
    /**
     * This function code is used to force each coil in a sequence of coils to either 
     * ON or OFF in a remote device. The Request PDU specifies the coil references to be
     * forced. Coils are addressed starting at zero. Therefore coil numbered 1 is
     * addressed as 0. The requested ON/OFF states are specified by contents of the
     * request data field. A logical '1' in a bit position of the field requests the 
     * corresponding output to be ON. A logical '0' requests it to be OFF.
     * The normal response returns the function code, starting address, and quantity
     * of coils forced.
     * 
     * @param rxADU
     * @return txADU
     */	
    @SuppressWarnings("unused")
    private ModbusADU doFC_Write_Multiple_Coils(ModbusADU rxADU){

        ModbusADU txADU = new ModbusADUImpl();

        int Starting_address;
        short Quantity_of_Outputs;

        int Byte_count;
        int Value;

        ModbusDevice md;

        txADU.setTransactionID(rxADU.getTransactionID());
        txADU.setProtocolID(rxADU.getProtocolID());
        txADU.setUnitID(rxADU.getUnitID());
        txADU.setFunctionCode(Modbus.FC_WRITE_MULTIPLE_COILS);
        ByteBuf txData = Unpooled.buffer(2);
        txADU.setData(txData);
        Starting_address = rxADU.getData().readShort();
        Quantity_of_Outputs = rxADU.getData().readShort();
        Byte_count = rxADU.getData().readByte();

        ByteBuf Coils = rxADU.getData().discardReadBytes();

        md = myMda.getModbusDevicesArray()[rxADU.getUnitID()];

        if ((Quantity_of_Outputs >= (short) 0x0001) && (Quantity_of_Outputs <= (short) 0x07B0)) {

            if ((Starting_address + Quantity_of_Outputs) <= (md.getCoils().capacity() / 8)) {
                try {

                    for (int i = 0; i < Quantity_of_Outputs; i++) {
                        md.setCoil(Starting_address + i, getDigitalPoint(Coils, i));
                    }
                    txADU.getData().writeShort(Quantity_of_Outputs);
                    
                } catch (Exception ex) {
                    txADU.setFunctionCode((byte) (Modbus.FC_WRITE_SINGLE_COIL + Modbus.EX_EXCEPTION_MODIFIER));
                    txADU.getData().clear();
                    txADU.getData().writeByte(Modbus.EX_SERVER_DEVICE_FAILURE);
                }
            } else {
                txADU.setFunctionCode((byte) (Modbus.FC_WRITE_MULTIPLE_COILS + Modbus.EX_EXCEPTION_MODIFIER));
                txADU.getData().writeByte(Modbus.EX_ILLEGAL_DATA_ADDRESS);
            }
        } else {
            txADU.setFunctionCode((byte) (Modbus.FC_WRITE_MULTIPLE_COILS + Modbus.EX_EXCEPTION_MODIFIER));
            txADU.getData().writeByte(Modbus.EX_ILLEGAL_DATA_VALUE);
        }

        return txADU;
    }
	
    /**
     * This function code is used to read from 1 to 125 contiguous input registers in
     * a remote device. The Request PDU specifies the starting register address and
     * the number of registers. In the PDU Registers are addressed starting at zero.
     * Therefore input registers numbered 1 -16 are addressed as 0-15.
     * The register data in the response message are packed as two bytes per register,
     * with the binary contents right justified within each byte. For each register,
     * the first byte contains the high order bits and the second contains the low 
     * order bits.
     * 
     * @param rxADU
     * @return txADU
     */
    @SuppressWarnings("unused")
    private ModbusADU doFC_Read_Input_Registers(ModbusADU rxADU){

        ModbusADU txADU = new ModbusADUImpl();

        int Starting_address;
        int Quantity_of_Input_Registers;

        int Value;

        txADU.setTransactionID(rxADU.getTransactionID());
        txADU.setProtocolID(rxADU.getProtocolID());
        txADU.setUnitID(rxADU.getUnitID());
        txADU.setFunctionCode(Modbus.FC_READ_INPUT_REGISTERS);

        Starting_address = rxADU.getData().readShort() * 2;
        Quantity_of_Input_Registers = rxADU.getData().readShort() * 2;

        ByteBuf txData = Unpooled.buffer(Quantity_of_Input_Registers + 1);
        txADU.setData(txData);
        ByteBuf InputRegisters = myMda.getModbusDevicesArray()[rxADU.getUnitID()].getInputRegisters();

        if ((Quantity_of_Input_Registers >= 0x0001) && ((Quantity_of_Input_Registers / 2) <= 0x007D)) {

            if ((Starting_address + Quantity_of_Input_Registers) <= InputRegisters.capacity()) {

                try {
                    txADU.getData().writeByte(Quantity_of_Input_Registers);
                    InputRegisters.getBytes(Starting_address, txADU.getData(), Quantity_of_Input_Registers);

                    //Unit + Function + Data todo en bytes.
                    Quantity_of_Input_Registers += 3;
                    txADU.setLengthField((short) Quantity_of_Input_Registers);
                } catch (Exception ex) {
                    txADU.setFunctionCode((byte) (Modbus.FC_READ_INPUT_REGISTERS + Modbus.EX_EXCEPTION_MODIFIER));
                    txADU.getData().clear();
                    txADU.getData().writeByte(Modbus.EX_SERVER_DEVICE_FAILURE);
                }

            } else {
                txADU.setFunctionCode((byte) (Modbus.FC_READ_INPUT_REGISTERS + Modbus.EX_EXCEPTION_MODIFIER));
                txADU.getData().writeByte(Modbus.EX_ILLEGAL_DATA_ADDRESS);
            }
        } else {
            txADU.setFunctionCode((byte) (Modbus.FC_READ_INPUT_REGISTERS + Modbus.EX_EXCEPTION_MODIFIER));
            txADU.getData().writeByte(Modbus.EX_ILLEGAL_DATA_VALUE);
        }

        return txADU;
    }	
	
    /**
     * This function code is used to read the contents of a contiguous block of 
     * holding registers in a remote device. The Request PDU specifies the starting
     * register address and the number of registers. In the PDU Registers are 
     * addressed starting at zero. Therefore registers numbered 1 -16 are addressed 
     * as 0-15.
     * The register data in the response message are packed as two bytes per 
     * register, with the binary contents right justified within each byte. For 
     * each register, the first byte contains the high order bits and the second 
     * contains the low order bits.
     * 
     * @param rxADU
     * @return txADU
     */
    @SuppressWarnings("unused")
    private ModbusADU doFC_Read_Holding_Registers(ModbusADU rxADU){

        int Starting_address;
        int Quantity_of_Holding_Registers;

        int Value;

        ModbusADU txADU = new ModbusADUImpl();
        txADU.setTransactionID(rxADU.getTransactionID());
        txADU.setProtocolID(rxADU.getProtocolID());
        txADU.setUnitID(rxADU.getUnitID());
        txADU.setFunctionCode(Modbus.FC_READ_HOLDING_REGISTERS);

        Starting_address = rxADU.getData().readShort() * 2;
        Quantity_of_Holding_Registers = rxADU.getData().readShort() * 2;

        ByteBuf txData = Unpooled.buffer(Quantity_of_Holding_Registers + 1);
        txADU.setData(txData);
        ByteBuf HoldingRegisters = myMda.getModbusDevicesArray()[rxADU.getUnitID()].getHoldingRegisters();

        if ((Quantity_of_Holding_Registers >= 0x0001) && ((Quantity_of_Holding_Registers / 2) <= 0x007D)) {

            if ((Starting_address + Quantity_of_Holding_Registers) <= HoldingRegisters.capacity()) {

                try {
                    txADU.getData().writeByte(Quantity_of_Holding_Registers);
                    HoldingRegisters.getBytes(Starting_address, txADU.getData(), Quantity_of_Holding_Registers);

                    //Unit + Function + Data todo en bytes.
                    Quantity_of_Holding_Registers += 3;
                    txADU.setLengthField((short) Quantity_of_Holding_Registers);
                } catch (Exception ex) {
                    txADU.setFunctionCode((byte) (Modbus.FC_READ_HOLDING_REGISTERS + Modbus.EX_EXCEPTION_MODIFIER));
                    txADU.getData().clear();
                    txADU.getData().writeByte(Modbus.EX_SERVER_DEVICE_FAILURE);
                }

            } else {
                txADU.setFunctionCode((byte) (Modbus.FC_READ_HOLDING_REGISTERS + Modbus.EX_EXCEPTION_MODIFIER));
                txADU.getData().writeByte(Modbus.EX_ILLEGAL_DATA_ADDRESS);
            }
        } else {
            txADU.setFunctionCode((byte) (Modbus.FC_READ_HOLDING_REGISTERS + Modbus.EX_EXCEPTION_MODIFIER));
            txADU.getData().writeByte(Modbus.EX_ILLEGAL_DATA_VALUE);
        }

        return txADU;
    }	

    /**
     * This function code is used to write a single holding register in a remote
     * device. The Request PDU specifies the address of the register to be written.
     * Registers are addressed starting at zero. Therefore register numbered 1 is
     * addressed as 0.
     * The normal response is an echo of the request, returned after the register 
     * contents have been written.
     * 
     * @param rxADU
     * @return txADU
     */
    private ModbusADU doFC_Write_Single_Register(ModbusADU rxADU){
        ModbusADU txADU = new ModbusADUImpl();

        int Starting_address;
        int Value;

        txADU.setTransactionID(rxADU.getTransactionID());
        txADU.setProtocolID(rxADU.getProtocolID());
        txADU.setUnitID(rxADU.getUnitID());
        txADU.setFunctionCode(Modbus.FC_WRITE_SINGLE_REGISTER);

        Starting_address = rxADU.getData().readShort() * 2;
        Value = rxADU.getData().readShort();

        ByteBuf txData = Unpooled.buffer(4);
        txADU.setData(txData);
        ByteBuf HoldingRegisters = myMda.getModbusDevicesArray()[rxADU.getUnitID()].getHoldingRegisters();

        if ((Value >= 0x000) && (Value <= 0xFFFF)) {

            if (Starting_address <= HoldingRegisters.capacity()) {

                try {
                    HoldingRegisters.setShort(Starting_address, Value);
                    txADU.getData().writeShort((short) (Starting_address / 2));
                    txADU.getData().writeShort((short) Value);
                    txADU.setLengthField(rxADU.getLengthField());
                } catch (Exception ex) {
                    txADU.setFunctionCode((byte) (Modbus.FC_WRITE_SINGLE_REGISTER + Modbus.EX_EXCEPTION_MODIFIER));
                    txADU.getData().clear();
                    txADU.getData().writeByte(Modbus.EX_SERVER_DEVICE_FAILURE);
                }
            } else {
                txADU.setFunctionCode((byte) (Modbus.FC_WRITE_SINGLE_REGISTER + Modbus.EX_EXCEPTION_MODIFIER));
                txADU.getData().writeByte(Modbus.EX_ILLEGAL_DATA_ADDRESS);
            }
        } else {
            txADU.setFunctionCode((byte) (Modbus.FC_WRITE_SINGLE_REGISTER + Modbus.EX_EXCEPTION_MODIFIER));
            txADU.getData().writeByte(Modbus.EX_ILLEGAL_DATA_VALUE);
        }
        return txADU;
    }

    /**
     * This function code is used to write a block of contiguous registers (1 to
     * 123 registers) in a remote device.
     * The requested written values are specified in the request data field. 
     * Data is packed as two bytes per register.
     * The normal response returns the function code, starting address, and 
     * quantity of registers written.
     * 
     * @param rxADU
     * @return txADU
     */
    @SuppressWarnings("unused")
    private ModbusADU doFC_Write_Multiple_Registers(ModbusADU rxADU){

        ModbusADU txADU = new ModbusADUImpl();

        int Starting_address;
        int Quantity_of_Registers;
        int Byte_count;
        int Value;
        int startaddress;

        txADU.setTransactionID(rxADU.getTransactionID());
        txADU.setProtocolID(rxADU.getProtocolID());
        txADU.setUnitID(rxADU.getUnitID());
        txADU.setFunctionCode(Modbus.FC_WRITE_MULTIPLE_REGISTERS);

        Starting_address = rxADU.getData().readShort();
        startaddress = Starting_address * 2;
        Quantity_of_Registers = rxADU.getData().readShort();
        Byte_count = rxADU.getData().readByte();

        ByteBuf txData = Unpooled.buffer(4);
        txADU.setData(txData);
        ByteBuf HoldingRegisters = myMda.getModbusDevicesArray()[rxADU.getUnitID()].getHoldingRegisters();

        if ((Quantity_of_Registers >= 0x0001) && (Quantity_of_Registers <= 0x007B)) {
            if (((Starting_address + Quantity_of_Registers) * 2) <= HoldingRegisters.capacity()) {

                try {
                    for (int i = 0; i < Quantity_of_Registers; i++) {
                        HoldingRegisters.setShort(startaddress, rxADU.getData().readShort());
                        startaddress += 2;
                    }
                    //cbuffer = dynamicBuffer(ByteOrder.BIG_ENDIAN , 4);
                    txADU.getData().writeShort((short) Starting_address);
                    txADU.getData().writeShort((short) Quantity_of_Registers);
                    txADU.setLengthField((short) 6);
                } catch (Exception ex) {
                    txADU.setFunctionCode((byte) (Modbus.FC_READ_HOLDING_REGISTERS + Modbus.EX_EXCEPTION_MODIFIER));
                    txADU.getData().clear();
                    txADU.getData().writeByte(Modbus.EX_SERVER_DEVICE_FAILURE);
                }

            } else {
                txADU.setFunctionCode((byte) (Modbus.FC_WRITE_MULTIPLE_REGISTERS + Modbus.EX_EXCEPTION_MODIFIER));
                txADU.getData().writeByte(Modbus.EX_ILLEGAL_DATA_ADDRESS);
            };

        } else {
            txADU.setFunctionCode((byte) (Modbus.FC_WRITE_MULTIPLE_REGISTERS + Modbus.EX_EXCEPTION_MODIFIER));
            txADU.getData().writeByte(Modbus.EX_ILLEGAL_DATA_VALUE);
        }
        return txADU;
    }

	
    /**
     * This function code is used to modify the contents of a specified holding
     * register using a combination of an AND mask, an OR mask, and the register's
     * current contents. The function can be used to set or clear individual bits
     * in the register.
     * The request specifies the holding register to be written, the data to be used
     * as the AND mask, and the data to be used as the OR mask. Registers are 
     * addressed starting at zero. Therefore registers 1 -16 are addressed as 0-15.
     * The function�s algorithm is:
     * Result = (Current Contents AND And_Mask) OR (Or_Mask AND (NOT And_Mask))
     * 
     * @param rxADU
     * @return
     */
    private ModbusADU doFC_Mask_Write_Register(ModbusADU rxADU){
        ModbusADU txADU = new ModbusADUImpl();

        int Reference_address;
        int AND_Mask;
        int OR_Mask;

        txADU.setTransactionID(rxADU.getTransactionID());
        txADU.setProtocolID(rxADU.getProtocolID());
        txADU.setUnitID(rxADU.getUnitID());
        txADU.setFunctionCode(Modbus.FC_MASK_WRITE_REGISTER);

        Reference_address = rxADU.getData().readShort() * 2;
        AND_Mask = rxADU.getData().readShort();
        OR_Mask = rxADU.getData().readShort();

        ByteBuf txData = Unpooled.buffer(6);
        txADU.setData(txData);
        ModbusDevice md = myMda.getModbusDevicesArray()[rxADU.getUnitID()];
        ByteBuf HoldingRegisters = myMda.getModbusDevicesArray()[rxADU.getUnitID()].getHoldingRegisters();

        if ((Reference_address <= HoldingRegisters.capacity())) {
            //(AND_mask == OK) AND (OR_Mask == OK)
            //How I can check thi?
            try {

                md.setHoldingRegister(Reference_address,
                        (short) ((HoldingRegisters.getShort(Reference_address) & AND_Mask)
                        | (OR_Mask & (~AND_Mask)))
                );
                txADU.getData().writeShort(Reference_address / 2);
                txADU.getData().writeShort(AND_Mask);
                txADU.getData().writeShort(OR_Mask);

            } catch (Exception ex) {
                txADU.setFunctionCode((byte) (Modbus.FC_MASK_WRITE_REGISTER + Modbus.EX_EXCEPTION_MODIFIER));
                txADU.getData().clear();
                txADU.getData().writeByte(Modbus.EX_SERVER_DEVICE_FAILURE);
            }

        } else {
            txADU.setFunctionCode((byte) (Modbus.FC_MASK_WRITE_REGISTER + Modbus.EX_EXCEPTION_MODIFIER));
            txADU.getData().writeByte(Modbus.EX_ILLEGAL_DATA_VALUE);
        }
        return txADU;
    }
	
    /**
     * This function code performs a combination of one read operation and one write
     * operation in a single MODBUS transaction. The write operation is performed
     * before the read.
     * Holding registers are addressed starting at zero. Therefore holding registers
     * 1 -16 are  addressed in the PDU as 0-15.
     * The request specifies the starting address and number of holding registers to
     * be read as well as the starting address, number of holding registers, and the
     * data to be written. The byte count specifies the number of bytes to follow in
     * the write data field.
     * The normal response contains the data from the group of registers that were 
     * read. The byte count field specifies the quantity of bytes to follow in the 
     * read data field.
     * 
     * @param rxADU
     * @return
     */
    private ModbusADU doFC_Read_Write_Multiple_Registers(ModbusADU rxADU){
        ModbusADU txADU = new ModbusADUImpl();

        int Read_Starting_Address;
        int Quantity_to_Read;
        int Write_Starting_Address;
        int Quantity_to_Write;
        byte Write_Byte_Count;

        txADU.setTransactionID(rxADU.getTransactionID());
        txADU.setProtocolID(rxADU.getProtocolID());
        txADU.setUnitID(rxADU.getUnitID());
        txADU.setFunctionCode(Modbus.FC_READ_WRITE_MULTIPLE_REGISTERS);

        Read_Starting_Address = rxADU.getData().readShort() * 2;
        Quantity_to_Read = rxADU.getData().readShort();
        Write_Starting_Address = rxADU.getData().readShort() * 2;
        Quantity_to_Write = rxADU.getData().readShort();
        Write_Byte_Count = rxADU.getData().readByte();

        ByteBuf txData = Unpooled.buffer(Quantity_to_Read + 1);
        txADU.setData(txData);
        ByteBuf HoldingRegisters = myMda.getModbusDevicesArray()[rxADU.getUnitID()].getInputRegisters();

        if (((Quantity_to_Read >= 0x0001) && (Quantity_to_Read <= 0x007D))
                && ((Quantity_to_Write >= 0x0001) && (Quantity_to_Write <= 0x0079))
                && ((Write_Byte_Count == (Quantity_to_Write * 2)))) {

            Quantity_to_Read *= 2;
            Quantity_to_Write *= 2;

            if (((Read_Starting_Address + Quantity_to_Read) <= (HoldingRegisters.capacity() / 2))
                    && (Write_Starting_Address + Quantity_to_Write) <= (HoldingRegisters.capacity() / 2)) {

                try {
                    HoldingRegisters.writeBytes(rxADU.getData(), Write_Starting_Address, Quantity_to_Write);

                    txADU.getData().writeByte(Quantity_to_Read);
                    HoldingRegisters.readBytes(txADU.getData(), 1, Quantity_to_Read);

                } catch (Exception ex) {
                    txADU.setFunctionCode((byte) (Modbus.FC_READ_WRITE_MULTIPLE_REGISTERS + Modbus.EX_EXCEPTION_MODIFIER));
                    txADU.getData().clear();
                    txADU.getData().writeByte(Modbus.EX_SERVER_DEVICE_FAILURE);
                }

            } else {
                txADU.setFunctionCode((byte) (Modbus.FC_READ_WRITE_MULTIPLE_REGISTERS + Modbus.EX_EXCEPTION_MODIFIER));
                txADU.getData().writeByte(Modbus.EX_ILLEGAL_DATA_ADDRESS);
            }

        } else {
            txADU.setFunctionCode((byte) (Modbus.FC_READ_WRITE_MULTIPLE_REGISTERS + Modbus.EX_EXCEPTION_MODIFIER));
            txADU.getData().writeByte(Modbus.EX_ILLEGAL_DATA_VALUE);
        }

        return txADU;
    }	
	
	
    @SuppressWarnings("unused")
    private void setDigitalPoint(ByteBuf Coils,int register, boolean state) {
        int intByte = (register / 8);
        int index = (register % 8);

        int temp = Coils.getByte(intByte);
        if (state) {
            temp = temp | 1 << index; // sets 1 at given index
        } else {
            temp = temp & (~(1 << index)); // sets 0 at given index
        }
        Coils.setByte(intByte, temp);

    };

    private boolean getDigitalPoint(ByteBuf Coils, int register) {
        int intByte = (register / 8);
        int index = (register % 8);
        return (Coils.getByte(intByte) & (1 << index)) != 0;	
    };
	
	
}
