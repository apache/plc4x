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
package org.apache.plc4x.merlot.api;

import io.netty.buffer.ByteBuf;


public interface ByteBasedDriver  extends Driver{
	
	//Administrative
	
	public long InitializeDevice(DriverCallback cb);
	public long ShutDownDevice(DriverCallback cb);
	
	//Device Handler
	
	public long GetDeviceMode(DriverCallback cb);
	public long StartDevice(DriverCallback cb);
	public long StopDevice(DriverCallback cb);
	public long GetPDUSize(DriverCallback cb);
	
	//Bits
	
	/*
	 * This function can be used to read bits in a device in packed format.
	 * The different values are stored in the ByteBuf table called. 
	 * The bits which are not read in the last byte that receives the data 
	 * are in an undetermined state. 
	 * @param channel Channel number (0-31)
	 * @param device Device number (0-255)
	 * @param nb Number of bits to be read. The maximum number depends on the protocol and target device.
	 * @param adr Address of the first bit to be read in the device.
	 * @return ByteBuf Table receiving read data.
	 */
	public ByteBuf ReadMemoryBits(DriverCallback cb, int device, int nb, int adr);
	
	/*
	 * This function can be used to write bits in a device in packed format.
	 * The different values are stored in the ByteBuf table called bf. 
	 * The bits which are not read in the last byte that receives the data 
	 * are in an undetermined state. 
	 * @param channel Channel number (0-31)
	 * @param device Device number (0-255)
	 * @param nb Number of bits to be read. The maximum number depends on the protocol and target device.
	 * @param adr Adress of the first bit to be write in the device.
	 * @return int Error status.
	 */	
	public long WriteMemoryBits(DriverCallback cb, int device, int nb, int adr, ByteBuf bf);

	public long ReadInputBits(DriverCallback cb, int device, int nb, int adr);
	
	public long ReadOutputBits(DriverCallback cb, int device, int nb, int adr);
	public long WriteOutputBits(DriverCallback cb, int device, int nb, int adr, ByteBuf bf);
	
	public long ReadDataBits(DriverCallback cb, int device, int nb, int adr);
	public long WriteDataBits(DriverCallback cb, int device, int nb, int adr, ByteBuf bf);	
	
	//Bytes

	public long ReadMemoryBytes(DriverCallback cb, int device, int nb, int adr);
	public long WriteMemoryBytes(DriverCallback cb, int device, int nb, int adr, ByteBuf bf);

	public long ReadInputBytes(DriverCallback cb, int device, int nb, int adr, ByteBuf bf);
	
	public long ReadOutputBytes(DriverCallback cb, int device, int nb, int adr);
	public long WriteOutputBytes(DriverCallback cb, int device, int nb, int adr, ByteBuf bf);

	public long ReadDataBytes(DriverCallback cb, int device, int nb, int adr);
	public long WriteDataBytes(DriverCallback cb, int device, int nb, int adr, ByteBuf bf);
	
	//Words
	
	public long ReadMemoryWords(DriverCallback cb, int device, int nb, int adr);
	public long WriteMemoryWords(DriverCallback cb, int device, int nb, int adr, ByteBuf bf);

	public long ReadInputWords(DriverCallback cb, int device, int nb, int adr);
	
	public long ReadOutputWords(DriverCallback cb, int device, int nb, int adr);
	public long WriteOutputWords(DriverCallback cb, int device, int nb, int adr, ByteBuf bf);

	public long ReadDataWords(DriverCallback cb, int device, int nb, int adr);
	public long WriteDataWords(DriverCallback cb, int device, int nb, int adr, ByteBuf bf);	
	
	//Double Words

	public long ReadMemoryDoubleWords(DriverCallback cb, int device, int nb, int adr);
	public long WriteMemorDoubleyWords(DriverCallback cb, int device, int nb, int adr, ByteBuf bf);

	public long ReadInputDoubleWords(DriverCallback cb, int device, int nb, int adr);
	
	public long ReadOutputDoubleWords(DriverCallback cb, int device, int nb, int adr);
	public long WriteOutputDoubleWords(DriverCallback cb, int device, int nb, int adr, ByteBuf bf);

	public long ReadDataDoubleWords(DriverCallback cb, int device, int nb, int adr);
	public long WriteDataDoubleWords(DriverCallback cb, int device, int nb, int adr, ByteBuf bf);
	
	//Timer and Counter functions
	
	public long ReadTimer(DriverCallback cb, int device, int nb, int adr);
	public long WriteTimer(DriverCallback cb, int device, int nb, int adr, ByteBuf bf);
	
	public long ReadCounter(DriverCallback cb, int device, int nb, int adr);
	public long WriteCounter(DriverCallback cb, int device, int nb, int adr, ByteBuf bf);	
	
}
