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

#include <boost/functional/hash.hpp>
#include "PlcUsernamePasswordAuthentication.h"

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
					namespace authentication
					{
						/**---------------------------------------------------------------------
						* ctor
						*---------------------------------------------------------------------*/
						PlcUsernamePasswordAuthentication::PlcUsernamePasswordAuthentication(std::string username, std::string password)
						{																					
							// Check isNull (Java) not required, is every time a valid string
							this->username = username;
							this->password = password;
						}
						PlcUsernamePasswordAuthentication::~PlcUsernamePasswordAuthentication()
						{
							// do nothing
						}
						/**---------------------------------------------------------------------
						* returns username
						*---------------------------------------------------------------------*/
						std::string PlcUsernamePasswordAuthentication::getUsername()
						{
							return username;
						}

						/**---------------------------------------------------------------------
						* returns password
						*---------------------------------------------------------------------*/
						std::string PlcUsernamePasswordAuthentication::getPassword()
						{
							return password;
						}
						/**---------------------------------------------------------------------
						* Check for equality (identical object or identical values)
						*---------------------------------------------------------------------*/						
						bool PlcUsernamePasswordAuthentication::equals(PlcUsernamePasswordAuthentication& auth)
						{
							bool bResult = false;
							
							if ( (this == &auth) || ( (username.compare(auth.getUsername()) == 0) && (password.compare(auth.getPassword()) == 0) ))
							{
								return true;
							}

							return bResult;
						}

						std::string PlcUsernamePasswordAuthentication::toString()
						{
							return "PlcUsernamePasswordAuthentication{ username='" + username + '\'' +
								", password='" + "*****************" + '\'' +
								'}';
						}

						int PlcUsernamePasswordAuthentication::hashCode()
						{
							boost::hash<std::string> string_hash;

							return string_hash(username + password);
						}

					}
				}
			}
		}
	}
}