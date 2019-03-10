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

#ifndef _PLC_TIMEOUT_EXCEPTION
#define _PLC_TIMEOUT_EXCEPTION

#include <string>
#include <chrono>

#include "PlcRuntimeException.h"

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
					namespace exceptions
					{
						/**
						 * Can be thrown when something times out.
						 */
						class PlcTimeoutException : public PlcRuntimeException
						{
							public:
								/**
								* Indicates something timed out.
								*
								* @param timeout in nanoseconds.
								*/
								explicit PlcTimeoutException(long lTimeout);

								long getTimeout();

							private:
								/**
								* @return the timeout in nanoseconds.
								*/
								long _lTimeout;
						};
					}
				}
			}
		}
	}
}

#endif

