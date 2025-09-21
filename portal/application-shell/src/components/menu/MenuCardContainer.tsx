import React, {ReactNode} from 'react';
import {Grid2} from "@mui/material";


interface MenuCardContainerProps {
    children: ReactNode
}

const MenuCardContainer: React.FC<MenuCardContainerProps> = ({children}) => {
    return <Grid2 container spacing={2}>
        {children}
    </Grid2>
}
export default MenuCardContainer