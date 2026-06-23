import { Button, ButtonGroup, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from "@mui/material";
import { Delete, Edit } from "@mui/icons-material";
import { FC } from "react";
import { SearchTagsTableMessageBundle } from "../../messages/MessageBundles";

type SearchTagsTableProps = {
    searchTagsRegistry: SearchTag[];
    handler: {
        onEdit: (searchTag: SearchTag) => void;
        onDelete: (searchTag: SearchTag) => void;
    };
    messages: SearchTagsTableMessageBundle;
}

const SearchTagsTable: FC<SearchTagsTableProps> = ({ searchTagsRegistry, handler, messages }) => {
    return <TableContainer component={Paper}>
        <Table aria-label="simple table">
            <TableHead>
                <TableRow>
                    <TableCell>{messages.description}</TableCell>
                    <TableCell>{messages.operation}</TableCell>
                </TableRow>
            </TableHead>
            <TableBody>
                {searchTagsRegistry.map((searchTag) => (
                    <TableRow key={searchTag.key}>
                        <TableCell scope="row">
                            {searchTag.value}
                        </TableCell>
                        <TableCell>
                            <ButtonGroup variant="contained" aria-label="outlined primary button group">
                                <Button onClick={() => handler.onEdit(searchTag)}><Edit /> {messages.actions.edit}</Button>
                                <Button onClick={() => handler.onDelete(searchTag)}><Delete /> {messages.actions.delete}</Button>
                            </ButtonGroup>
                        </TableCell>
                    </TableRow>
                ))}
            </TableBody>
        </Table>
    </TableContainer>
}

export default SearchTagsTable;