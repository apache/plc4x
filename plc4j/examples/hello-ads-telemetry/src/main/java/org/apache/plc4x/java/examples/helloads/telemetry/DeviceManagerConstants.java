package org.apache.plc4x.java.examples.helloads.telemetry;

// Extracted from https://infosys.beckhoff.com/index.php?content=../content/1031/devicemanager/262982923.html

public enum DeviceManagerConstants {
    NIC(0x00000002),
    Time(0x00000003),
    UserManagement(0x00000004),
    RAS(0x00000005),
    FTP(0x00000006),
    SMB(0x00000007),
    TwinCat(0x00000008),
    Software(0x0000000A),
    COU(0x0000000B),
    Memory(0x0000000C),
    FirewallWinCE(0x0000000E),
    FileSystemObject(0x00000010),
    DisplayDevice(0x00000013),
    EWF(0x00000014),
    FBWF(0x00000015),
    OS(0x00000018),
    RAID(0x00000019),
    Fan(0x0000001B),
    Mainboard(0x0000001C),
    DiskManagement(0x0000001D),
    UPS(0x0000001E),
    PhysicalDrive(0x0000001F),
    MassStorageDrive(0x00000020),
    UnifiedWriteFilter(0x00000021),
    IO(0x00000022),
    Misc(0x00000100);

    final int typeNumber;
    DeviceManagerConstants(int typeNumber) {
        this.typeNumber = typeNumber;
    }
}
