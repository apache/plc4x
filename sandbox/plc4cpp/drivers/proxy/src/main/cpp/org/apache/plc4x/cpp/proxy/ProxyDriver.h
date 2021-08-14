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

#ifndef _PROXY_DRIVER
#define _PROXY_DRIVER

#include <org/apache/plc4x/cpp/spi/PlcDriver.h>
#include <org/apache/plc4x/cpp/api/authentication/PlcAuthentication.h>
#include <string>

using namespace org::apache::plc4x::cpp::spi;

namespace org
{
	namespace apache
	{
		namespace plc4x
		{
			namespace cpp
			{
				namespace proxy
				{
					/**
					* Interface defining the most basic methods a PLC4X connection should support.
					* This generally handles the connection establishment itself and the parsing of
					* field address strings to the platform dependent PlcField instances.
					*/
					class ProxyDriver : public PlcDriver
					{

					public:
						/**
						 * @return code of the implemented protocol. This is usually a lot shorter than the String returned by @see #getProtocolName().
						 */
						virtual std::string getProtocolCode();

						/**
						 * @return name of the implemented protocol.
						 */
						virtual std::string getProtocolName();

						/**
						 * Connects to a PLC using the given plc connection string.
						 * @param url plc connection string.
						 * @return PlcConnection object.
						 * @throws PlcConnectionException an exception if the connection attempt failed.
						 */
						virtual PlcConnection* connect(std::string url);

						/**
						 * Connects to a PLC using the given plc connection string using given authentication credentials.
						 * @param url plc connection string.
						 * @param authentication authentication credentials.
						 * @return PlcConnection object.
						 * @throws PlcConnectionException an exception if the connection attempt failed.
						 */
						virtual PlcConnection* connect(std::string url, PlcAuthentication authentication);

					private:

                        const std::string PROXY_URI_PATTERN="proxy://(?<host>.*)(?<params>\\?.*)?";
                  
					};
				}
			}
		}
	}
}

#endif

