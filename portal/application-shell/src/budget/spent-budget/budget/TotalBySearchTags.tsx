import React from "react"
import { Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from "@mui/material";

type TotalBySearchTagsProps = {
    totals: {
        searchTagValue: string,
        total: number
    }[]
}

const TotalBySearchTags: React.FC<TotalBySearchTagsProps> = ({ totals }) => {
    return <TableContainer component={Paper}>
        <Table>
            <TableHead>
                <TableRow>
                    <TableCell>Category</TableCell>
                    <TableCell>Total</TableCell>
                </TableRow>
            </TableHead>
            <TableBody>
                {totals.map(total =>
                    <TableRow key={"ST-" + total.searchTagValue} sx={{ '&:last-child td, &:last-child th': { border: 0 } }}>
                        <TableCell>{total.searchTagValue}</TableCell>
                        <TableCell>{total.total}</TableCell>
                    </TableRow>)}
            </TableBody>
        </Table>
    </TableContainer>
}

export default TotalBySearchTags