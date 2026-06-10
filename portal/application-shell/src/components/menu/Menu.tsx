import React from "react"
import {AppBar, Box, Button, Divider, Drawer, IconButton, List, Toolbar, Typography, useTheme} from "@mui/material";
import LogoutIcon from '@mui/icons-material/Logout';
import MenuIcon from '@mui/icons-material/Menu';
import GlobalPageNavigation from "./GlobalPageNavigation";
import MenuSectionTitle from "./MenuSectionTitle";
import { MenuMessageBundle } from "../../messages/MessageBundles";

type MenuItemProps = {
    text: string
}

type MenuProps = {
    messages: MenuMessageBundle,
    children?: React.ReactNode,
    navBarItems: React.ReactNode,
    showGlobalNav?: boolean
}

const Menu: React.FC<MenuProps> = ({messages, children, navBarItems, showGlobalNav = true}) => {
    const links = {
        logOut: "/logout",
        home: "/"
    }
    let theme = useTheme()
    const [openDrawer, setOpenDrawer] = React.useState(false);

    const toggleDrawer = (open: boolean) => (event: React.KeyboardEvent | React.MouseEvent) => {
        // if (event.type === 'keydown' && (event.key === 'Tab' || event.key === 'Shift')) {
        //     return;
        // }
        setOpenDrawer(open)
    };

    return <Box style={{
        marginBottom: "35px"
    }}>
        <AppBar position="fixed">
            <Toolbar variant="dense">
                <IconButton
                    onClick={() => setOpenDrawer(!openDrawer)}
                    size="large"
                    edge="start"
                    color="inherit"
                    aria-label="menu"
                    sx={{mr: 2}}>
                    <MenuIcon/>
                </IconButton>
                <Typography variant="h6" component="div" sx={{flexGrow: 1}}>
                    <a href={links.home} style={{
                        "textDecoration": "none",
                        "color": theme.palette.primary.contrastText
                    }}>{messages.title}</a>
                </Typography>

                <Drawer
                    open={openDrawer}
                    onClose={toggleDrawer(false)}>
                    <Box
                        role="presentation"
                        onClick={toggleDrawer(false)}
                        onKeyDown={toggleDrawer(false)}>
                        <List>
                            {/*{children && <>*/}
                            {/*    <MenuSectionTitle content={"Section Specific Feature"}/>*/}
                            {/*    {children}*/}
                            {/*</>}*/}

                            {children}
                            {children && showGlobalNav && <Divider/>}
                            {showGlobalNav && <>

                                <MenuSectionTitle content={messages.otherFeatureLabel}/>

                                <GlobalPageNavigation messages={{
                                    budget: messages.budgetPageLabel,
                                    revenue: messages.revenuePageLabel,
                                    plans: messages.planPageLabel,
                                    account: messages.accountPageLabel,
                                    analytics: messages.analyticsPageLabel
                                }}/>
                            </>}
                        </List>
                    </Box>
                </Drawer>
                {navBarItems}
                <form action={links.logOut} method="GET">
                    <Button color="inherit" type="submit">
                        {messages.logOutLabel} <LogoutIcon/>
                    </Button>
                </form>
            </Toolbar>
        </AppBar>
    </Box>
}

export default Menu
export type {MenuItemProps}