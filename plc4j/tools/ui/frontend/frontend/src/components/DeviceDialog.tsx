import {Device} from "../generated/plc4j-tools-ui-frontend.ts";
import {Dialog} from "primereact/dialog";
import {Button} from "primereact/button";
import {InputText} from "primereact/inputtext";
import withWidth from "@mui/material/Hidden/withWidth";
import {DataTable} from "primereact/datatable";
import {Column} from "primereact/column";

interface ConnectionDialogProps {
    device: Device
    visible: boolean

    onSave: (device:Device) => void
    onCancel: () => void
}

interface TableEntry {
    key: string
    value: string
}

type Dictionary = { [index: string]: string }

function mapToTableEntry(map : Dictionary):TableEntry[] {
    let tableEntries = [] as TableEntry[]
    for(const key in map) {
        const value = map[key]
        tableEntries = [...tableEntries, {key: key, value: value}]
    }
    return tableEntries
}

export default function DeviceDialog({device, visible, onSave, onCancel}: ConnectionDialogProps) {
    function handleSave() {
        onSave(device)
    }
    function handleCancel() {
        onCancel()
    }

    return (
        <Dialog visible={visible} modal style={{width: '60rem'}} draggable={true} resizable={true} onHide={() => {
        }}>
            <div className="formgrid grid">
                <label htmlFor="name" className="col-12 mb-2 md:col-2 md:mb-0">Device Name</label>
                <div className="col-12 md:col-10">
                    <InputText id="name" value={device.name}/>
                </div>
            </div>
            <div className="formgrid grid">
                <label htmlFor="protocol" className="col-12 mb-2 md:col-2 md:mb-0">Protocol</label>
                <div className="col-12 md:col-10">
                    <InputText id="protocol" value={device.protocolCode}/>
                </div>
            </div>
            <div className="formgrid grid">
                <label htmlFor="transportCode" className="col-12 mb-2 md:col-2 md:mb-0">Transport Type</label>
                <div className="col-12 md:col-10">
                    <InputText id="transportCode" value={device.transportCode}/>
                </div>
            </div>
            <div className="formgrid grid">
                <label htmlFor="transportUrl" className="col-12 mb-2 md:col-2 md:mb-0">Transport URL</label>
                <div className="col-12 md:col-10">
                    <InputText id="transportUrl" value={device.transportUrl}/>
                </div>
            </div>
            <div className="formgrid grid">
                <label htmlFor="options" className="col-12 mb-2 md:col-2 md:mb-0">Options</label>
                <div className="col-12 md:col-10">
                    <DataTable id="options" value={mapToTableEntry(device.options)} tableStyle={{minWidth: '50rem'}}>
                        <Column field="key" header="Name"/>
                        <Column field="value" header="Value"/>
                    </DataTable>
                </div>
            </div>
            <div className="formgrid grid">
                <label htmlFor="attributes" className="col-12 mb-2 md:col-2 md:mb-0">Attributes</label>
                <div className="col-12 md:col-10">
                    <DataTable id="attributes" value={mapToTableEntry(device.attributes)} tableStyle={{minWidth: '50rem'}}>
                        <Column field="key" header="Name"/>
                        <Column field="value" header="Value"/>
                    </DataTable>
                </div>
            </div>
            <div className="formgrid grid">
                <div className="field col">
                    <Button label={"Save"} onClick={handleSave}/>
                    <Button label={"Cancel"} onClick={handleCancel}/>
                </div>
            </div>
        </Dialog>
    )
}