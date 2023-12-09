import {Outlet} from "react-router-dom";
import {AppBar, Box, CssBaseline, Drawer, IconButton, Toolbar, Typography} from "@mui/material";
import MenuIcon from "@mui/icons-material/Menu";
import MainMenu from "../components/MainMenu.tsx";
import React from "react";

const drawerWidth = 480;

export default function MainLayout() {
    const [mobileOpen, setMobileOpen] = React.useState(false);

    const handleDrawerToggle = () => {
        setMobileOpen(!mobileOpen);
    };

    return(
        <Box sx={{display: 'flex'}}>
            <CssBaseline/>
            {/* This is the main application bar */}
            <AppBar
                position="fixed"
                sx={{
                    width: {sm: `calc(100% - ${drawerWidth}px)`},
                    ml: {sm: `${drawerWidth}px`},
                }}>
                <Toolbar>
                    {/* This is the menu button, which is displayed in very narrow screens (Hidden otherwise) */}
                    <IconButton
                        color="inherit"
                        aria-label="open drawer"
                        edge="start"
                        onClick={handleDrawerToggle}
                        sx={{mr: 2, display: {sm: 'none'}}}
                    >
                        <MenuIcon/>
                    </IconButton>
                    {/* This is the title text in the main application bar */}
                    <Typography variant="h6" noWrap component="div">
                        Responsive drawer
                    </Typography>
                </Toolbar>
            </AppBar>
            {/* This is the drawer content (Menu) */}
            <Box
                component="nav"
                sx={{width: {sm: drawerWidth}, flexShrink: {sm: 0}}}
                aria-label="mailbox folders"
            >
                {/* This is the dynamic drawer used on mobile devices. */}
                <Drawer
                    variant="temporary"
                    open={mobileOpen}
                    onClose={handleDrawerToggle}
                    ModalProps={{
                        keepMounted: true, // Better open performance on mobile.
                    }}
                    sx={{
                        display: {xs: 'block', sm: 'none'},
                        '& .MuiDrawer-paper': {boxSizing: 'border-box', width: drawerWidth},
                    }}
                >
                    <MainMenu/>
                </Drawer>
                {/* This is the static drawer used on normal devices. */}
                <Drawer
                    variant="permanent"
                    sx={{
                        display: {xs: 'none', sm: 'block'},
                        '& .MuiDrawer-paper': {boxSizing: 'border-box', width: drawerWidth},
                    }}
                    open
                >
                    <MainMenu/>
                </Drawer>
            </Box>
            {/* This is the main application content */}
            <Box
                component="main"
                sx={{flexGrow: 1, p: 3, width: {sm: `calc(100% - ${drawerWidth}px)`}}}
            >
                {/* This is a dummy toolbar to prevent the content from going behind the real application bar */}
                <Toolbar/>
                <Outlet/>
            </Box>
        </Box>
    )
}