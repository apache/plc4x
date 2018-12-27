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

#include "PlcProtocolPayloadTooBigException.h"

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

						PlcProtocolPayloadTooBigException::PlcProtocolPayloadTooBigException(const std::string &protocolName, int maxSize, int actualSize, std::vector<char> payload) :
							PlcProtocolException("Payload for protocol '" + protocolName + "' with size " + std::to_string(actualSize) + " exceeded allowed maximum of " + std::to_string(maxSize))
						{
							_protocolName = protocolName;
							_maxSize = maxSize;
							_actualSize = actualSize;
							_payload = payload;
						}

						std::string PlcProtocolPayloadTooBigException::getProtocolName() 
						{ 
							return _protocolName; 
						}
						
						int PlcProtocolPayloadTooBigException::getMaxSize() 
						{ 
							return _maxSize; 
						}

						int PlcProtocolPayloadTooBigException::getActualSize() 
						{ 
							return _actualSize; 
						}

						std::vector<char> PlcProtocolPayloadTooBigException::getPayload() 
						{ 
							return _payload; 
						}

					}
				}			
			}
		}
	}
}
