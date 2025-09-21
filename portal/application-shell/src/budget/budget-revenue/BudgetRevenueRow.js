import React from "react"
import {Button, ButtonGroup, TableCell, TableRow} from "@mui/material";
import {Delete, Edit} from "@mui/icons-material";
import {v1 as uuidv1} from "uuid";


export default ({revenue, openDeletePopUp, openUpdatePopUp}) => {
    return <TableRow key={uuidv1()} sx={{'&:last-child td, &:last-child th': {border: 0}}}>
        <TableCell>{revenue.date}</TableCell>
        <TableCell>{revenue.amount}</TableCell>
        <TableCell>{revenue.note}</TableCell>
        <TableCell>
            <ButtonGroup variant="contained" aria-label="outlined primary button group">
                <Button onClick={openUpdatePopUp}><Edit/> Edit</Button>
                <Button onClick={openDeletePopUp}><Delete/>Delete </Button>
            </ButtonGroup>
        </TableCell>
    </TableRow>

}