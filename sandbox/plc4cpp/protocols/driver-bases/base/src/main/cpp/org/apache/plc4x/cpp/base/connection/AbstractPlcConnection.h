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

#ifndef _ABSTRCT_PLC_HANDLER
#define _ABSTRCT_PLC_HANDLER

#include <org/apache/plc4x/cpp/api/PlcConnection.h>
#include <boost/system/error_code.hpp>
#include <string>

using namespace org::apache::plc4x::cpp::api;

namespace org
{
	namespace apache
	{
		namespace plc4x
		{
			namespace cpp
			{
				namespace base
				{
					namespace connection
					{
						/**
						 * Base class for implementing connections.
						 * Per default, all operations (read, write, subscribe) are unsupported.
						 * Concrete implementations should override the methods indicating connection capabilities
						 * and for obtaining respective request builders.
						 */
						class AbstractPlcConnection : public PlcConnection, PlcConnectionMetadata
						{

						public:
							PlcConnectionMetadata* getMetadata();
							bool canRead();
							bool canWrite();
							bool canSubscribe();

							virtual void open() = 0;
							virtual void close() = 0;
							virtual bool send(unsigned char* pBytesToSend, int iNumBytesToSend) = 0;

							PlcReadRequest::Builder* readRequestBuilder();
							// PlcWriteRequest.Builder writeRequestBuilder();
							// PlcSubscriptionRequest.Builder subscriptionRequestBuilder();
							// PlcUnsubscriptionRequest.Builder nsubscriptionRequestBuilder();

						protected:
							virtual void onReceive(const boost::system::error_code& errorCode, std::size_t bytes_transferred) = 0;
						private:

						};
					}
				}
			}
		}
	}
}

#endif

