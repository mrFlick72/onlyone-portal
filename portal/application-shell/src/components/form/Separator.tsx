import React from "react";

import {Divider, Grid} from "@mui/material";
import {commonStyle} from "../../theme/ThemeProvider";

export default function Separator() {
    return <Grid style={commonStyle.formRow}>
        <Divider/>
    </Grid>
}