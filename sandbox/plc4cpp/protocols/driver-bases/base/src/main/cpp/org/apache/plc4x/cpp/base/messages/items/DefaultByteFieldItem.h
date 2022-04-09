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

#ifndef _DEFAULT_BYTE_FIELD_ITEM
#define _DEFAULT_BYTE_FIELD_ITEM

#include <boost/multiprecision/cpp_dec_float.hpp>
#include "BaseDefaultFieldItem.h"

using namespace org::apache::plc4x::cpp::api::exceptions;
using namespace boost::multiprecision;

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
							class DefaultByteFieldItem : public BaseDefaultFieldItem<char>
							{

							public:
								template<typename T>
								DefaultByteFieldItem(std::vector<T> values);

								template<typename T>
								inline T getObject(int index) { return BaseDefaultFieldItem::getValue(index); }
								inline bool isValidBoolean(int index) {	return true;	}
								bool getBoolean(int index);
								bool isValidByte(int index);
								char getByte(int index);
								bool isValidShort(int index);
								short getShort(int index);
								bool isValidInteger(int index);
								int getInteger(int index);
								bool isValidLong(int index);
								long getLong(int index);
								bool isValidBigInteger(int index);
								long long getBigInteger(int index);
								bool isValidFloat(int index);
								float getFloat(int index);
								bool isValidDouble(int index);
								double getDouble(int index);
								bool isValidBigDecimal(int index);
								cpp_dec_float_100 getBigDecimal(int index);

							private:
							};
						}
					}
				}
			}
		}
	}
}

#endif