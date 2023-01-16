package org.apache.plc4x.java.profinet.gsdml;

import java.util.List;

public interface ProfinetDeviceItem {

    List<ProfinetVirtualSubmoduleItem> getVirtualSubmoduleList();

    ProfinetSystemDefinedSubmoduleList getSystemDefinedSubmoduleList();


}
