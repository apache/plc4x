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

#ifndef _PLC_WRITE_REQUEST
#define _PLC_WRITE_REQUEST

#include "PlcFieldRequest.h"
#include <vector>
#include <boost/multiprecision/cpp_dec_float.hpp>
#include <ctime>

using namespace boost::multiprecision;

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
						class PlcWriteRequest : public PlcFieldRequest
						{
							// Todo: @override CompletableFuture<? extends PlcWriteResponse> execute();
						public:							
							virtual PlcWriteRequest* build() = 0;

							virtual void addItem(std::string name, std::string fieldQuery, std::vector<bool> values) = 0;

							virtual void addItem(std::string name, std::string fieldQuery, std::vector<char> values) = 0;

							virtual void addItem(std::string name, std::string fieldQuery, std::vector<short> values) =  0;

							virtual void addItem(std::string name, std::string fieldQuery, std::vector<int> values)  = 0;

							virtual void  addItem(std::string name, std::string fieldQuery, std::vector<long long> values) = 0;

							virtual void addItem(std::string name, std::string fieldQuery, std::vector<long> values) = 0;

							virtual void addItem(std::string name, std::string fieldQuery, std::vector<float> values) = 0;

							virtual void addItem(std::string name, std::string fieldQuery, std::vector<double> values) = 0;

							virtual void addItem(std::string name, std::string fieldQuery, std::vector<cpp_dec_float_100> values) = 0;

							virtual void addItem(std::string name, std::string fieldQuery, std::vector<std::string> values) = 0;

							virtual void addItem(std::string name, std::string fieldQuery, std::vector<std::tm> values) = 0;

							virtual void addItem(std::string name, std::string fieldQuery, std::vector<std::vector<char> > values) = 0;

							//Todo: <T> PlcWriteRequest.Builder addItem(String name, String fieldQuery, T... values);
							/*template<class T>
							virtual void addItem(std::string name, std::string fieldQuery, std::vector<T> values) = 0; */

						private:
						};
					}
				}
			}
		}
	}
}

#endif

