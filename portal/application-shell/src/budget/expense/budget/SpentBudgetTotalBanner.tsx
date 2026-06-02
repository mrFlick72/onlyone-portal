import React from "react"
import { Box, Paper, Typography } from "@mui/material";

type SpentBudgetTotalBannerProps = {
    total: string,
    totalLabel: string
}

const SpentBudgetTotalBanner: React.FC<SpentBudgetTotalBannerProps> = ({ total, totalLabel }) => {
    return <Box component={Paper} color="primary">
        <Typography sx={{ margin: "10px" }}>{totalLabel} {total || 0.00}</Typography>
    </Box>
}

export default SpentBudgetTotalBanner