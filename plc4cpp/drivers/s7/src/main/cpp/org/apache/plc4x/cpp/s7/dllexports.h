//-------- DllExports.h --------//
#ifndef _DLLEXPORTS_H
#define _DLLEXPORTS_H

#include <org/apache/plc4x/cpp/spi/PlcDriver.h>

/*#ifdef __dll__
#define IMPEXP __declspec(dllexport)
#else
#define IMPEXP __declspec(dllimport)
#endif 	// __dll__*/

extern "C" __declspec(dllexport) org::apache::plc4x::cpp::spi::PlcDriver* _CreatePlcDriverInstance();


#endif	// DLLEXPORTS_H