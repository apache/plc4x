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

#ifndef _PLC_DRIVER
#define _PLC_DRIVER

#include "../api/PlcConnection.h"
#include "../api/authentication/PlcAuthentication.h"
#include <string>

using namespace org::apache::plc4x::cpp::api::authentication;
using namespace org::apache::plc4x::cpp::api;

namespace org
{
	namespace apache
	{
		namespace plc4x
		{
			namespace cpp
			{
				namespace spi
				{
					/**
					* Interface defining the most basic methods a PLC4X connection should support.
					* This generally handles the connection establishment itself and the parsing of
					* field address strings to the platform dependent PlcField instances.
					*/
					class PlcDriver
					{

					public:
						/**
						 * @return code of the implemented protocol. This is usually a lot shorter than the String returned by @see #getProtocolName().
						 */
						virtual std::string getProtocolCode() = 0;

						/**
						 * @return name of the implemented protocol.
						 */
						virtual std::string getProtocolName() = 0;

						/**
						 * Connects to a PLC using the given plc connection string.
						 * @param url plc connection string.
						 * @return PlcConnection object.
						 * @throws PlcConnectionException an exception if the connection attempt failed.
						 */
						virtual PlcConnection* connect(std::string url) = 0;

						/**
						 * Connects to a PLC using the given plc connection string using given authentication credentials.
						 * @param url plc connection string.
						 * @param authentication authentication credentials.
						 * @return PlcConnection object.
						 * @throws PlcConnectionException an exception if the connection attempt failed.
						 */
						virtual PlcConnection* connect(std::string url, PlcAuthentication authentication) = 0;

					private:
					};
				}
			}
		}
	}
}

#endif

