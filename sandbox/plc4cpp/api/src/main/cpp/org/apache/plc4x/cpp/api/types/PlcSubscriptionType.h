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

#ifndef _PLC_SUBSCRIPTION_TYPE
#define _PLC_SUBSCRIPTION_TYPE

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
					namespace types
					{
						enum PlcSubscriptionType
						{
							/**
							 * A cyclic subscription where a value is sent no matter if it's value changed in a given interval.
							 */
							CYCLIC,

							/**
							 * Only send data, if a value in the PLC changed.
							 */
							CHANGE_OF_STATE,

							/**
							 * Subscribe to events created by the PLC which usually are defined in the PLCs application (Alarms).
							 */
							EVENT
						};
					}
				}
			}
		}
	}
}

#endif