//-------- DllMain.cpp --------//
#define __dll__
#include "DllExports.h"
#include "S7PlcDriver.h"
#include <windows.h>

using namespace org::apache::plc4x::cpp::s7;
using namespace org::apache::plc4x::cpp::spi;

int WINAPI DllEntryPoint(HINSTANCE hinst, unsigned long reason, void*)
{
    return 1;
}

_declspec(dllexport) PlcDriver* _CreatePlcDriverInstance()
{
    return new S7PlcDriver;
}