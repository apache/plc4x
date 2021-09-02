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

#ifndef _B_LOGGER_H_
#define _B_LOGGER_H_

#include <boost/log/core.hpp>
#include <boost/log/trivial.hpp>
#include <boost/log/expressions.hpp>
#include <boost/log/sources/severity_logger.hpp>
#include <boost/thread/mutex.hpp>
#include "ExLog.h"

#include "DatDmp.h"

#define LOG_TRACE(x)  	org::apache::plc4x::cpp::utils::BLogger::Instance()->Log(x,boost::log::trivial::severity_level::trace)
#define LOG_DEBUG(x)  	org::apache::plc4x::cpp::utils::BLogger::Instance()->Log(x,boost::log::trivial::severity_level::debug)
#define LOG_INFO(x)  	org::apache::plc4x::cpp::utils::BLogger::Instance()->Log(x,boost::log::trivial::severity_level::info)
#define LOG_WARNING(x)  org::apache::plc4x::cpp::utils::BLogger::Instance()->Log(x, boost::log::trivial::severity_level::warning)
#define LOG_ERROR(x)  	org::apache::plc4x::cpp::utils::BLogger::Instance()->Log(x, boost::log::trivial::severity_level::error)
#define LOG_FATAL(x)  	org::apache::plc4x::cpp::utils::BLogger::Instance()->Log(x, boost::log::trivial::severity_level::fatal)
#define LOG_BOOL(x,y)  	org::apache::plc4x::cpp::utils::BLogger::Instance()->LogBool(x,y)

#define IS_LOG(x) 		org::apache::plc4x::cpp::utils::BLogger::Instance()->isLog(x)


namespace org
{
	namespace apache
	{
		namespace plc4x
		{
			namespace cpp
			{
				namespace utils
				{
					namespace src = boost::log::sources;
					namespace expr = boost::log::expressions;
					namespace sinks = boost::log::sinks;
					namespace keywords = boost::log::keywords;
					namespace BoostLogger = boost::log;

					using namespace BoostLogger::trivial;
					using namespace std;

					
					class BLogger
					{

					public:

						static BLogger* Instance();
						void Release();

						~BLogger();

						void Initialize();

						void Log(const char* strMessage, severity_level logLevel = boost::log::trivial::severity_level::info);

						void Log(string strMessage, severity_level logLevel = boost::log::trivial::severity_level::info);

						void LogBool(bool b, severity_level logLevel = boost::log::trivial::severity_level::info);

						void LogBool(string strMessage, bool b, severity_level logLevel = boost::log::trivial::severity_level::info);

						void LogHexDump(unsigned char* pBinaryData, unsigned int nNumBytes, const char* szDumpPrefix = NULL, severity_level logLevel = boost::log::trivial::severity_level::info);

						string getFileExtension() { return _strFileExtension; }
						void setFileExtension(string strExtension);

						string getFileTitle() { return _strFileTitle; }
						void setFileTitle(string strFileTitle);

						int getMaxFileSize() const { return _iMaxFileSize; }
						void setMaxFileSize(int maxFileSize);

						unsigned int getMaxNumberOfFiles() { return _uiMaxNumberOfFiles; }
						void setMaxNumberOfFiles(unsigned int uiMaxNumberOfFiles);

						severity_level getSeverityLogLevel() { return _severityLevel; }
						void setSeverityLogLevel(severity_level severityLevel);

						int getSeverityLogLevelNumber() { return getNumberfromSeverityLevel(_severityLevel); }
						void setSeverityLogLevelNumber(int iSeverityLevelNumber);

						bool isLog(severity_level logLevel);

						string getTrueMessage() const { return _strTrueMessage; }
						void setTrueMessage(string strTrueMessage) { _strTrueMessage = strTrueMessage; }

						string getFalseMessage() const { return _strFalseMessage; }
						void setFalseMessage(string strFalseMessage) { _strFalseMessage = strFalseMessage; }

						void clearLogDir();
						void ForceFlushLog(void);

					private:

						BLogger();  // Private so that it can  not be called
						BLogger(BLogger const&){};             // copy constructor is private
						BLogger& operator=(BLogger const&){ };  // assignment operator is private

						int getNumberfromSeverityLevel(severity_level severityLevel);
						severity_level getSeverityLevelfromNumber(int iSeverityLevelNumber);

						void FlushLog(void);

						static BLogger* m_pInstance;

						int _iMaxFileSize;
						unsigned int _uiMaxNumberOfFiles;

						unsigned int _uiFlushCounter;

						unsigned int _uiMaxFlushEntrys[6];

						severity_level _severityLevel;
						int 		   _iSeverityLevelNumber;

						boost::log::filter _logFilter;

						src::severity_logger_mt<severity_level> _lg;

						boost::mutex         _mtxLock;

						string _strFileTitle;
						string _strFileExtension;

						string _strTrueMessage;
						string _strFalseMessage;

						static const char* _pDefaultFileTitle;
						static const char* _pDefaultFileExtension;
					};
				}
			}
		}
	}
}
#endif