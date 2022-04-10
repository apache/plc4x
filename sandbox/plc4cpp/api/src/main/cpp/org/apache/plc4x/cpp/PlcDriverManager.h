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

#ifndef _PLC_DRIVER_MANAGER
#define _PLC_DRIVER_MANAGER
#include <vector>
#include <map>
#include <string>
#include <iostream>
#include "./api/PlcConnection.h"
#include "./spi/PlcDriver.h"
#include "./api/exceptions/PlcConnectionException.h"

using namespace org::apache::plc4x::cpp::api::exceptions;
using namespace org::apache::plc4x::cpp::spi;

namespace org
{
	namespace apache
	{
		namespace plc4x
		{
			namespace cpp
			{
				class PlcDriverManager
				{				

					public:
						PlcDriverManager();

						~PlcDriverManager();

						PlcConnection* getConnection(std::string url);
						PlcConnection* getConnection(std::string url, PlcAuthentication& authentication);

					protected:
						
					private:
						std::map<std::string, PlcDriver*> _mapDrivers;
						PlcDriver* getDriver(std::string url);

                        void findDrivers();

                        const std::string PLC_DRIVER_TEMPLATE = "plc4cpp-driver-(.*)(.dll|.so)";
                        const std::string PLC_CREATE_DRIVER_INSTANCE = "_CreatePlcDriverInstance";
                        const std::string REGEX_PROTOCOL = "^(?<proto>.*)://";
                        

				};
			}
		}
	}
}

#endif