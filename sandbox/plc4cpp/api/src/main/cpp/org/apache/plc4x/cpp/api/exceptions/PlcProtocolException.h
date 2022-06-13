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

#ifndef _PLC_PROTOCOL_EXCEPTION
#define _PLC_PROTOCOL_EXCEPTION

#include "PlcException.h"

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
						 * Base exception for all protocol related exceptions.
						 * This is usually if the messages themselves are malformed, contain incorrect information or
						 * any other problems occur which are related to the protocol being implemented.
						 */
						class PlcProtocolException : public PlcException
						{
							public:
								explicit PlcProtocolException(const std::string& strMessage): PlcException(strMessage){}
								PlcProtocolException(const std::string& strMessage, const std::exception& exException): PlcException(strMessage, exException){};
								explicit PlcProtocolException(const std::exception& exException): PlcException(exException){}
						};
					}
				}
			}
		}
	}
}

#endif

