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

#ifndef _PLC_UNSUBSCRIPTION_REQUEST
#define _PLC_UNSUBSCRIPTION_REQUEST

#include <vector>
#include "PlcRequest.h"
#include "../model/PlcSubscriptionHandle.h"

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
						class PlcUnsubscriptionRequest : public PlcRequest
						{
								
							public:			
								virtual PlcRequest* build() = 0;
								
								// 1st approach to migrate builder pattern: Instead of Builder-object simple methods are used to handle inside the given parameters.
								virtual void addHandles(PlcSubscriptionHandle* plcSubscriptionHandle) = 0;
							
								virtual void addHandles(PlcSubscriptionHandle* plcSubscriptionHandle1, std::vector<PlcSubscriptionHandle*> plcSubscriptionHandles) = 0;

								virtual void addHandles(std::vector<PlcSubscriptionHandle*> plcSubscriptionHandle) = 0;
						
							private:
						};
					}
				}
			}
		}
	}
}
#endif

