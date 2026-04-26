import Typography from "@mui/material/Typography";
import React from "react";

type MenuSectionTitlePorps = {
    content: string
}
const MenuSectionTitle: React.FC<MenuSectionTitlePorps> = ({content}) => {
    return <Typography variant="subtitle1" component="div"
                       sx={{flexGrow: 1, textAlign: "center", backgroundColor: "lightgrey"}}>
        {content}
    </Typography>

}


export default MenuSectionTitle