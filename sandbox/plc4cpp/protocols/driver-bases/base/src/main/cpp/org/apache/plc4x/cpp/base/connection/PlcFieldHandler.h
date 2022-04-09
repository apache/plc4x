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

#ifndef _PLC_FIELD_HANDLER
#define _PLC_FIELD_HANDLER

#include <boost/system/error_code.hpp>
#include <org/apache/plc4x/cpp/api/PlcConnection.h>
#include <string>
#include <vector>
#include <ctime>

#include <org/apache/plc4x/cpp/api/model/PlcField.h>
#include "../messages/items/BaseDefaultFieldItem.h"

using namespace org::apache::plc4x::cpp::api::model;
using namespace org::apache::plc4x::cpp::base::messages;
using namespace org::apache::plc4x::cpp::base::messages::items;
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
					namespace connection
					{
						/**
						 * Base class for implementing connections.
						 * Per default, all operations (read, write, subscribe) are unsupported.
						 * Concrete implementations should override the methods indicating connection capabilities
						 * and for obtaining respective request builders.
						 */
						class PlcFieldHandler
						{

						public:

							virtual PlcField* createField(std::string strFieldQuery) = 0;

							virtual BaseDefaultFieldItem<bool>* encodeBoolean(PlcField plcField, std::vector<void*> pValues) = 0;

							virtual BaseDefaultFieldItem<char>* encodeByte(PlcField plcField, std::vector<void*> pValues)= 0;

							virtual BaseDefaultFieldItem<short>* encodeShort(PlcField plcField, std::vector<void*> pValues) = 0;

							virtual BaseDefaultFieldItem<int>* encodeInteger(PlcField plcField, std::vector<void*> pValues) = 0;

							virtual BaseDefaultFieldItem<long long>* encodeBigInteger(PlcField plcField,std::vector<void*> pValues) = 0;

							virtual BaseDefaultFieldItem<long>* encodeLong(PlcField plcField, std::vector<void*> pValues) = 0;

							virtual BaseDefaultFieldItem<float>* encodeFloat(PlcField plcField, std::vector<void*> pValues) = 0;

							virtual BaseDefaultFieldItem< cpp_dec_float_100>* encodeBigDecimal(PlcField plcField, std::vector<void*> pValues) = 0;

							virtual BaseDefaultFieldItem<double>* encodeDouble(PlcField plcField, std::vector<void*> pValues) = 0;

							virtual BaseDefaultFieldItem<std::string>* encodeString(PlcField plcField, std::vector<void*> pValues) = 0;

							virtual BaseDefaultFieldItem<time_t>* encodeTime(PlcField plcField, std::vector<void*> pValues) = 0;

							virtual BaseDefaultFieldItem<tm>* encodeDate(PlcField plcField, std::vector<void*> pValues) = 0;

							virtual BaseDefaultFieldItem<tm>* encodeDateTime(PlcField plcField, std::vector<void*> pValues)= 0;

							virtual BaseDefaultFieldItem<char*>* encodeByteArray(PlcField fieplcFieldld,std::vector<void*> pValues) = 0;

						};
					}
				}
			}
		}
	}
}

#endif