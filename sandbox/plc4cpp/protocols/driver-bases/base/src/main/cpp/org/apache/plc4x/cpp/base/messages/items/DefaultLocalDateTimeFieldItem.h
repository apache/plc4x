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

#ifndef _DEFAULT_LOCAL_DATE_TIME_FIELD_ITEM
#define _DEFAULT_LOCAL_DATE_TIME_FIELD_ITEM

#include <boost/multiprecision/cpp_dec_float.hpp>
#include <boost/date_time/posix_time/conversion.hpp>
#include <ctime>
#include "BaseDefaultFieldItem.h"

using namespace org::apache::plc4x::cpp::api::exceptions;
using namespace boost::multiprecision;
using namespace boost::posix_time;

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
							class DefaultLocalDateTimeFieldItem : public BaseDefaultFieldItem<time_t>
							{
								typedef time_t data_type;
							public:
								template<typename T>
								DefaultLocalDateTimeFieldItem(std::vector<T> values);

								template<typename T>
								inline T getObject(int index) { return BaseDefaultFieldItem::getValue(index); }
								bool isValidTime(int index);
								time_t getTime(int index);
								bool isValidDate(int index);
								tm getDate(int index);
								bool isValidDateTime(int index);
								time_t getDateTime(int index);
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