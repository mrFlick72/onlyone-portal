import React from "react"
import {Box, Paper, Typography} from "@mui/material";

export default ({total}) => {
    return <Box component={Paper} color="primary">
        <Typography sx={{margin: "10px"}}>Total: {total || 0.00}</Typography>
    </Box>
}