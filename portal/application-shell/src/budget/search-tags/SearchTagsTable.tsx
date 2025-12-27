import { Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from "@mui/material";
import { ModeEdit } from "@mui/icons-material";
import FormButton from "../../components/form/FormButton";
import { FC } from "react";

interface SearchTagsTableProps {
    searchTagsRegistry: SearchTag[];
    handler: (searchTagKey: string, searchTagValue: string) => void;
}

const SearchTagsTable: FC<SearchTagsTableProps> = ({ searchTagsRegistry, handler }) => {
    return <TableContainer component={Paper}>
        <Table aria-label="simple table">
            <TableHead>
                <TableRow>
                    <TableCell>Description</TableCell>
                    <TableCell>Operation</TableCell>
                </TableRow>
            </TableHead>
            <TableBody>
                {searchTagsRegistry.map((searchTag) => (
                    <TableRow key={searchTag.key}>
                        <TableCell scope="row">
                            {searchTag.value}
                        </TableCell>
                        <TableCell>
                            <FormButton type="button"
                                label=""
                                labelPrefix={<ModeEdit />}
                                onClickHandler={() => handler(searchTag.key, searchTag.value)}
                            />
                        </TableCell>
                    </TableRow>
                ))}
            </TableBody>
        </Table>
    </TableContainer>
}

export default SearchTagsTable;