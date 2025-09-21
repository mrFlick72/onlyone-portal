import React, { ReactNode } from "react";

import { ListItem, ListItemButton, ListItemIcon, ListItemText } from "@mui/material";
import { v1 as uuidv1 } from "uuid";


type OpenPopUpMenuItem = {
    icon: ReactNode
    text: string
    openPopupHandler: React.MouseEventHandler
}

const OpenPopUpMenuItem: React.FC<OpenPopUpMenuItem> = ({ icon, text, openPopupHandler }) => {
    return <ListItem key={uuidv1()} disablePadding component="a" onClick={openPopupHandler}>
        <ListItemButton>
            {icon && <ListItemIcon>{icon}</ListItemIcon>}
            <ListItemText primary={text} />
        </ListItemButton>
    </ListItem>
}

export default OpenPopUpMenuItem