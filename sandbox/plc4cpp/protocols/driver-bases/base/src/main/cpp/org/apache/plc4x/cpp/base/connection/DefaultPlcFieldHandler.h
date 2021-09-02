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

#ifndef _DEFAULT_PLC_FIELD_HANDLER
#define _DEFAULT_PLC_FIELD_HANDLER

#include <boost/system/error_code.hpp>
#include <org/apache/plc4x/cpp/api/PlcConnection.h>
#include <string>
#include <boost/multiprecision/cpp_dec_float.hpp>
#include <boost/functional/hash.hpp>
#include <ctime>
#include <array>
#include <vector>

#include <org/apache/plc4x/cpp/api/model/PlcField.h>
#include "PlcFieldHandler.h"
#include "../messages/items/BaseDefaultFieldItem.h"

using namespace org::apache::plc4x::cpp::api::model;
using namespace org::apache::plc4x::cpp::base::messages::items;

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
					 * Base Implementation of {@link PlcFieldHandler} which throws a {@link PlcRuntimeException} for all
					 * encodeXXX methods.
						 */
						class DefaultPlcFieldhandler : public PlcFieldHandler
						{

						public:
							
							BaseDefaultFieldItem<bool>* encodeBoolean(PlcField plcField, std::vector<void*> pValues) override;
								
							BaseDefaultFieldItem<char>* encodeByte(PlcField plcField, std::vector<void*> pValues) override;
														
							BaseDefaultFieldItem<short>* encodeShort(PlcField plcField, std::vector<void*> pValues) override;
														
							BaseDefaultFieldItem<int>* encodeInteger(PlcField plcField, std::vector<void*> pValues) override;
					
							BaseDefaultFieldItem<long long>* encodeBigInteger(PlcField plcField, std::vector<void*> pValues) override;
								
							BaseDefaultFieldItem<long>* encodeLong(PlcField plcField, std::vector<void*> pValues) override;
													
							BaseDefaultFieldItem<float>* encodeFloat(PlcField plcField, std::vector<void*> pValues) override;
							
							BaseDefaultFieldItem<cpp_dec_float_100>* encodeBigDecimal(PlcField plcField,std::vector<void*> pValues) override;
							
							BaseDefaultFieldItem<double>* encodeDouble(PlcField plcField, std::vector<void*> pValues) override;
							
							BaseDefaultFieldItem<std::string>* encodeString(PlcField plcField, std::vector<void*> pValues) override;
						
							BaseDefaultFieldItem<time_t>* encodeTime(PlcField plcField, std::vector<void*> pValues) override;
						
							BaseDefaultFieldItem<tm>* encodeDate(PlcField plcField, std::vector<void*> pValues) override;
							
							BaseDefaultFieldItem<tm>* encodeDateTime(PlcField plcField, std::vector<void*> pValues) override;
								
							BaseDefaultFieldItem<char*>* encodeByteArray(PlcField plcField, std::vector<void*> pValues) override;
							
						};
					}
				}
			}
		}
	}
}

#endif