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

#ifndef _PROXY_CONNECTION
#define _PROXY_CONNECTION

#include <org/apache/plc4x/cpp/api/metadata/PlcConnectionMetadata.h>
#include <org/apache/plc4x/cpp/api/messages/PlcReadRequest.h>
#include <org/apache/plc4x/cpp/api/messages/PlcWriteRequest.h>
#include <org/apache/plc4x/cpp/api/messages/PlcSubscriptionRequest.h>
#include <org/apache/plc4x/cpp/api/messages/PlcUnsubscriptionRequest.h>
#include <org/apache/plc4x/cpp/base/connection/BoostConnection.h>
#include <org/apache/plc4x/cpp/base/connection/ChannelFactory.h>
#include <org/apache/plc4x/cpp/utils/logger/BLogger.h>
#include <org/apache/plc4x/cpp/utils/systemconfig/SystemConfiguration.h>


/**
 * Class implementing the Connection handling for Siemens S7.
 * The adressing of Values in S7 works as follows:
 * <p>
 * For adressing values from Datablocks the following syntax is used:
 * <pre>
 *     DATA_BLOCKS/{blockNumer}/{byteOffset}
 * </pre>
 * <p>
 * For adressing data from other memory segments like I/O, Markers, ...
 * <pre>
 *     {memory area}/{byte offset}
 *     or
 *     {memory area}/{byte offset}/{bit offset}
 * </pre>
 * where the {bit-offset} is optional.
 * All Available Memory Areas for this mode are defined in the {@link MemoryArea} enum.
 */

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
   
                    using namespace org::apache::plc4x::cpp::api::metadata;
                    using namespace org::apache::plc4x::cpp::api::messages;
                    using namespace org::apache::plc4x::cpp::api;
                    using namespace org::apache::plc4x::cpp::base::connection;
                    using namespace org::apache::plc4x::cpp::utils;

                    /**
					* Interface defining the most basic methods a PLC4X connection should support.
					* This generally handles the connection establishment itself and the parsing of
					* field address strings to the platform dependent PlcField instances.
					*/
                    class ProxyConnection : public BoostConnection
                    {
                    public:

                        ProxyConnection();
                        ProxyConnection(std::string strHost, string strParams);
                        
                        ~ProxyConnection();

                        void connect() {};

                        /**
                            * Provides connection metadata.
                            */
                        PlcConnectionMetadata* getMetadata() { return NULL; };

                        /**
                            * Obtain read request builder.
                            * @throws PlcUnsupportedOperationException if the connection does not support reading
                            * Todo: implement pendant for java Builder pattern
                            */
                        PlcReadRequest::Builder* readRequestBuilder() { return NULL; };

                        /**
                            * Obtain write request builder.
                            * @throws PlcUnsupportedOperationException if the connection does not support writing
                            * Todo: implement pendant for java Builder pattern (PlcWriteRequest.Builder)
                            */
                        PlcWriteRequest* writeRequestBuilder() { return NULL; };

                        /**
                            * Obtain subscription request builder.
                            * @throws PlcUnsupportedOperationException if the connection does not support subscription
                            * Todo: implement pendant for java Builder pattern (PlcSubscriptionRequest.Builder)
                            */
                        PlcSubscriptionRequest* subscriptionRequestBuilder() { return NULL; };

                        /**
                            * Obtain unsubscription request builder.
                            * @throws PlcUnsupportedOperationException if the connection does not support subscription
                            * Todo: implement pendant for java Builder pattern (PlcSubscriptionRequest.Builder)

                            */
                        PlcUnsubscriptionRequest* unsubscriptionRequestBuilder() { return NULL; };

                        void ping() {};

                        bool send(unsigned char* pBytesToSend, int iNumBytesToSend) { return true;};

                        void onReceive(const boost::system::error_code& errorCode, std::size_t szRecvBytes) {};
                        
						private:

                            SystemConfiguration CONF;

					};
				}
			}
		}
	}
}

#endif