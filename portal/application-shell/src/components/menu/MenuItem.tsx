import React, { ReactNode } from "react";
import { ListItem, ListItemButton, ListItemIcon, ListItemText } from "@mui/material";
import { v1 as uuidv1 } from 'uuid';

type MenuItemProps = {
    icon: ReactNode
    text: string
    link: string
}

const MenuItem: React.FC<MenuItemProps> = ({ icon, text, link }) => {
    return <ListItem key={uuidv1()} disablePadding component="a" href={link}>
        <ListItemButton>
            {icon && <ListItemIcon>{icon}</ListItemIcon>}
            <ListItemText primary={text} />
        </ListItemButton>
    </ListItem>
}

export default MenuItem