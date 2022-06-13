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

#ifndef _BASE_DEFAULT_FIELD_ITEM
#define _BASE_DEFAULT_FIELD_ITEM

#include <boost/system/error_code.hpp>
#include <boost/multiprecision/cpp_dec_float.hpp>
#include <boost/functional/hash.hpp>
#include <ctime>
#include <array>

#include <org/apache/plc4x/cpp/api/exceptions/PlcIncompatibleDatatypeException.h>
#include <org/apache/plc4x/cpp/api/exceptions/PlcFieldRangeException.h>
#include <org/apache/plc4x/cpp/api/types/ValueTypeObject.h>

#include <string>

using namespace org::apache::plc4x::cpp::api::exceptions;
using namespace org::apache::plc4x::cpp::api::types;


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
							template<typename T> 
							class BaseDefaultFieldItem : public ValueTypeObject
							{

							public:
								inline int getNumberOfValues() { return _values.length; }
								virtual T getObject(int index) = 0;
								inline bool isValidBoolean(int index) { return false; }
								inline bool getBoolean(int index) { throw new PlcIncompatibleDatatypeException("bool", index); }
								inline bool isValidByte(int index) { return false; }
								inline char getByte(int index) { throw new PlcIncompatibleDatatypeException("byte", index); }
								inline bool isValidShort(int index) { return false; }
								inline short getShort(int index) { throw new PlcIncompatibleDatatypeException("short", index); }
								inline bool isValidInteger(int index) { return false; }
								inline  int getInteger(int index) { throw new PlcIncompatibleDatatypeException("int", index); }
								inline bool isValidLong(int index) { return false; }
								inline long getLong(int index) { throw new PlcIncompatibleDatatypeException("long", index); }
								inline bool isValidBigInteger(int index) { return false; }
								inline long long getBigInteger(int index) { throw new PlcIncompatibleDatatypeException("long long", index); }
								inline bool isValidFloat(int index) { return false; }
								inline float getFloat(int index) { throw new PlcIncompatibleDatatypeException("float", index); }
								inline bool isValidDouble(int index) { return false; }
								inline double getDouble(int index) { throw new PlcIncompatibleDatatypeException("double", index); }
								inline bool isValidBigDecimal(int index) { return false; }
								inline long long getBigDecimal(int index) { throw new PlcIncompatibleDatatypeException("long long", index); }
								inline bool isValidString(int index) { return false; }
								inline std::string getString(int index) { throw new PlcIncompatibleDatatypeException("string", index); }
								inline bool isValidTime(int index) { return false; }
								inline tm* getTime(int index) { throw new PlcIncompatibleDatatypeException("tm*", index); }
								inline bool isValidDate(int index) { return false; }
								inline tm* getDate(int index) { throw new PlcIncompatibleDatatypeException("tm*", index); }
								inline bool isValidDateTime(int index) { return false; }
								inline tm* getDateTime(int index) { throw new PlcIncompatibleDatatypeException("tm*", index); }
								inline bool isValidByteArray(int index) { return false; }
								inline char* getByteArray(int index) { throw new PlcIncompatibleDatatypeException("char", index); }
								inline std::vector<T>* getValues() { return &_values; }
								inline  int hashCode() { return int_hash(_values); } // Todo: @Override
								inline std::string toString() // Todo:  @Override								 
								{									
									return "BaseDefaultFielditem { values=" + std::string(_values) + "}";
								};

								// Todo: add methods for unsigned data types?

								/* Todo:
								@Override
									bool equals(Object o)
								{
									if (this == o)
									{
										return true;
									}
									if (!(o instanceof BaseDefaultFieldItem))
									{
										return false;
									}
									BaseDefaultFieldItem< ? > fieldItem = (BaseDefaultFieldItem< ? >) o;
									return Arrays.equals(values, fieldItem.values);
								}*/

							protected:
								//BaseDefaultFieldItem() { this._values = new array<T>(0); }
								inline BaseDefaultFieldItem(T* values) { _values = values; }
								inline T getValue(unsigned int index)
								{
									if ((index < 0 || (index >= _values.size())))
									{
										throw new PlcFieldRangeException((_values.size() == 0) ? -1 : _values.size() - 1, index);
									}
									return _values[index];
								}

								std::vector<T> _values;
							};
						}
					}
				}
			}
		}
	}
}

#endif

