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

#include "DefaultFloatFieldItem.h"

namespace org
{
	namespace apache
	{
		namespace plc4x
		{
			namespace cpp
			{
				namespace base
				{
					namespace messages
					{
						namespace items
						{
							// ==================================================
							bool DefaultFloatFieldItem::getBoolean(int index)
							{
								if (!isValidBoolean(index))
								{
									throw new PlcIncompatibleDatatypeException("bool", index);
								}

								return getValue(index);
							}

							// ===================================================
							bool DefaultFloatFieldItem::isValidByte(int index)
							{
								float value = getValue(index);
								return value <= -127 && value >= -128;
							}

							// ===================================================
							char DefaultFloatFieldItem::getByte(int index)
							{
								if (!isValidByte(index)) 
								{
									throw new PlcIncompatibleDatatypeException("char", index);
								}
								return (char) getValue(index);
							}

							// ===================================================
							bool DefaultFloatFieldItem::isValidShort(int index)
							{
								float value = getValue(index);
								return value >= -32768 && value <= 32767;
							}

							// ===================================================
							short DefaultFloatFieldItem::getShort(int index)
							{
								if (!isValidShort(index))
								{
									throw new PlcIncompatibleDatatypeException("short", index);
								}
								
								return (short)getValue(index);
							}

							// ===================================================
							bool DefaultFloatFieldItem::isValidInteger(int index)
							{
								float value = getValue(index);
								return value >= 0 && value <= 2147483647;
							}

							// ===================================================
							int DefaultFloatFieldItem::getInteger(int index)
							{
								if (!isValidInteger(index))
								{
									throw new PlcIncompatibleDatatypeException("int", index);
								}

								return (int)getValue(index);
							}

							// ===================================================
							bool DefaultFloatFieldItem::isValidLong(int index)
							{
								float value = getValue(index);
								return value >= -(2 ^ 63) && value <= (2 ^ 63) - 1;
							}

							// ===================================================
							long DefaultFloatFieldItem::getLong(int index)
							{
								if (!isValidLong(index))
								{
									throw new PlcIncompatibleDatatypeException("long", index);
								}

								return (long)getValue(index);
							}

							// ===================================================
							bool DefaultFloatFieldItem::isValidBigInteger(int index)
							{
								return true;
							}

							// ===================================================
							long long DefaultFloatFieldItem::getBigInteger(int index)
							{
								if (!isValidBigInteger(index))
								{
									throw new PlcIncompatibleDatatypeException("long long", index);
								}
								
								// Todo: add conversion from cpp_dec_float_100 TO long long
								return (long)getValue(index);
							}

							// ===================================================
							bool DefaultFloatFieldItem::isValidFloat(int index)
							{
								float value = getValue(index);
								return value >= -3.4e38 && value <= 3.4e38;
							}

							// ===================================================
							float DefaultFloatFieldItem::getFloat(int index)
							{
								if (!isValidFloat(index))
								{
									throw new PlcIncompatibleDatatypeException("float", index);
								}

								return (float)getValue(index);
							}

							// ===================================================
							bool DefaultFloatFieldItem::isValidDouble(int index)
							{
								double value = getValue(index);
								return value >= -1.7e308 && value <= 1.7e308;
							}

							// ===================================================
							double DefaultFloatFieldItem::getDouble(int index)
							{
								if (!isValidDouble(index))
								{
									throw new PlcIncompatibleDatatypeException("double", index);
								}

								return (double)getValue(index);
							}

							// ===================================================
							bool DefaultFloatFieldItem::isValidBigDecimal(int index)
							{
								// same limits but higher precision
								return isValidDouble(index);
							}

							// ===================================================
							cpp_dec_float_100 DefaultFloatFieldItem::getBigDecimal(int index)
							{
								if (!isValidBigDecimal(index))
								{
									throw new PlcIncompatibleDatatypeException("cpp_dec_float_100", index);
								}

								return (cpp_dec_float_100)getValue(index);
							}
						}
						
					}
				}
			}
		}
	}
}
