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

#ifndef _PLC_UNSUPPORTED_DATATYPE_EXCEPTION
#define _PLC_UNSUPPORTED_DATATYPE_EXCEPTION

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
						 * Indicate that a data type ({@link Class}) is not supported by Plc4x.
						 */
						class PlcUnsupportedDataTypeException : public PlcRuntimeException
						{
							public:
								explicit PlcUnsupportedDataTypeException(const std::string& strMessage): PlcRuntimeException(strMessage){}
								PlcUnsupportedDataTypeException(const std::string& strMessage, const std::exception &exException): PlcRuntimeException(strMessage, exException){}
								explicit PlcUnsupportedDataTypeException(const std::exception &exException): PlcRuntimeException(exException){}
						};
					}
				}
			}
		}
	}
}

#endif