import React from "react";

import {Divider, Grid2} from "@mui/material";
import {commonStyle} from "../../theme/ThemeProvider";

export default function Separator() {
    return <Grid2 style={commonStyle.formRow}>
        <Divider/>
    </Grid2>
}