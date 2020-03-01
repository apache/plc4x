/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

#ifndef _PLC_FIELD_REQUEST
#define _PLC_FIELD_REQUEST

#include "../model/PlcField.h"
#include "PlcFieldResponse.h"

#include <string>
#include <map>

namespace org
{
	namespace apache
	{
		namespace plc4x
		{
			namespace cpp
			{
				namespace api
				{
					namespace messages
					{
						/**
						 * Base type for all messages sent from the plc4x system to a connected plc.
						 */
						class PlcFieldRequest : public PlcRequest
						{
							public:	
								virtual int getNumberOfFields() = 0;
                                virtual std::vector<std::string> getFieldNames();
								virtual PlcField* getField(std::string) = 0;
								virtual std::vector<PlcField*> getFields() = 0;
                                virtual std::map<std::string, PlcField*> getNamedFields() = 0;

							private:
						};
					}
				}
			}
		}
	}
}

#endif

