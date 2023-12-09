import {TabPanel, TabView} from "primereact/tabview";

export default function Editor() {
    return (
        <TabView style={{width: "100%", height: "100%"}}>
            <TabPanel key={"ads://192.168.23.20"} header="ads://192.168.23.20" closable={true}>
                <p>test</p>
            </TabPanel>
        </TabView>
    )
}
