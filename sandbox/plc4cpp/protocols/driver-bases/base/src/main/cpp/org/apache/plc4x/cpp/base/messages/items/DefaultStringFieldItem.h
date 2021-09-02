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

#ifndef _DEFAULT_STRING_FIELD_ITEM
#define _DEFAULT_STRING_FIELD_ITEM

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
							class DefaultStringFieldItem : public BaseDefaultFieldItem<std::string>
							{
								
							public:
								template<typename T>
								DefaultStringFieldItem(std::vector<T> values);

								template<typename T>
								inline T getObject(int index) { return BaseDefaultFieldItem::getValue(index); }
								inline bool isValidString(int index);
								std::string getString(int index);

							private:
								typedef std::string data_type;
							};
						}
					}
				}
			}
		}
	}
}

#endif