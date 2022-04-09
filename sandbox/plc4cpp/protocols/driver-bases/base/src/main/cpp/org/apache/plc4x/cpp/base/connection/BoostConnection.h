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

#ifndef _BOOST_COMMUNICATION
#define _BOOST_COMMUNICATION

#include <string>
#include <boost/asio.hpp>
#include <boost/thread/thread.hpp>
#include <boost/atomic/atomic.hpp>
#include <org/apache/plc4x/cpp/PlcConnection.h>


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
                        using namespace boost::asio;
                        using namespace boost::asio::ip;
                        using namespace org::apache::plc4x::cpp::api;

                        class BoostConnection: public PlcConnection
                        {
                        public:

                            BoostConnection();

                            ~BoostConnection();

                            inline unsigned short getPort() { return _endpoint.port(); }
                            void setPort(unsigned short usPort);

                            inline string getIPAddress() { return _strIPAddress; }
                            void setIPAddress(string strIPAddress);

                            inline int getConnectionTimeout() { return _uiConnectionTimeout; }
                            void setConnectionTimeout(int uiConnectionTimeout) { _uiConnectionTimeout = uiConnectionTimeout; }

                            inline tcp::endpoint getEndpoint() { return _endpoint; }
                            void setEndpoint(tcp::endpoint endpoint);

                            inline bool getIsReceiving() { return _isReceiving; }

                            virtual void connect();
                            virtual void close();
                            virtual bool send(unsigned char* pBytesToSend, int iBytesToSend);

                            void ping();

                        protected:

                            void clearRessources();

                            void resolveEndpoint();

                            virtual void onConnect(const boost::system::error_code& errorCode, tcp::resolver::iterator endpointIterator);
                            virtual void onSend(const boost::system::error_code& ErrorCode, std::size_t bytes_transferred);
                            virtual void onReceive(const boost::system::error_code& ErrorCode, std::size_t bytes_transferred);
                            virtual void dispatchReceive();

                            virtual void continueReceive();

                            string endPointToString(tcp::endpoint ep);

                        protected:
                            io_service _io_service;
                            io_service::work* _pWorker;
                            boost::thread* _pReceiveDispatcher_Thread;

                            tcp::endpoint _endpoint;
                            tcp::resolver::iterator _iterEndPoint;
                            tcp::socket* _pSocket;
                            address _ipAddress;
                            string  _strIPAddress;
                            unsigned short  _usPort;
                            unsigned int _uiConnectionTimeout;
                            long _lCloseDeviceTimeoutMS;
                            
                            boost::mutex _mutexSend;
                            boost::atomic<bool> _isReceiving;

                            // General Properties
                            bool            _bWriteInProgress;
                            boost::thread*  _pClientThread;

                        private:

                        };
                    }
                }
            }
	    }
    }
}

#endif