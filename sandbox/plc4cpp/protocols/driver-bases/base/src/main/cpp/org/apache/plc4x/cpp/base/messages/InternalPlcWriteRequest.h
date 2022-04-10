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

#ifndef _INTERNAL_PLC_WRITE_REQUEST
#define _INTERNAL_PLC_WRITE_REQUEST

#include <org/apache/plc4x/cpp/api/messages/PlcWriteRequest.h>
#include "./items/BaseDefaultFieldItem.h"
#include "InternalPlcFieldRequest.h"
#include <vector>
#include <tuple>

using namespace org::apache::plc4x::cpp::api::messages;
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
					namespace messages
					{
						
						class InternalPlcWriteRequest : public PlcWriteRequest, public InternalPlcFieldRequest
						{
						public:
							virtual ValueTypeObject getFieldItem(std::string name) = 0;
							virtual std::vector<ValueTypeObject> getFieldItems(std::string name)= 0;
							virtual std::tuple<std::string, PlcField, ValueTypeObject> getNamedFieldTriples() = 0;
						};
					}
				}
			}
		}
	}
}

#endif