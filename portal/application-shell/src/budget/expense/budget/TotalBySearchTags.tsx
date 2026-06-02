import React from "react"
import { Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from "@mui/material";
import { TotalDetail } from "../domain/BudgetExpense";
import { TotalBySearchTagsMessageBundle } from "../../../messages/MessageBundles";

type TotalBySearchTagsProps = {
    totals:  TotalDetail[],
    messages: TotalBySearchTagsMessageBundle
}

const TotalBySearchTags: React.FC<TotalBySearchTagsProps> = ({ totals, messages }) => {
    return <TableContainer component={Paper}>
        <Table>
            <TableHead>
                <TableRow>
                    <TableCell>{messages.category}</TableCell>
                    <TableCell>{messages.total}</TableCell>
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