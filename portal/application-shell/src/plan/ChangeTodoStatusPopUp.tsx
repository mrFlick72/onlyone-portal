import React from "react";
import { Box, Button, Chip, Dialog, DialogActions, DialogContent, DialogTitle, Typography } from "@mui/material";
import { Todo, TodoStatus } from "./domain/Plan";
import { allowedTransitionsFor, colorFor, labelFor } from "./domain/TodoStatus";

interface ChangeTodoStatusPopUpProps {
    open: boolean;
    todo: Todo | null;
    handleClose: () => void;
    onSelect: (status: TodoStatus) => void;
}

const buttonColorFor = (status: TodoStatus) => {
    const color = colorFor(status);
    return color === "default" ? "inherit" : color;
};

const ChangeTodoStatusPopUp: React.FC<ChangeTodoStatusPopUpProps> = ({ open, todo, handleClose, onSelect }) => {
    const currentStatus = todo?.status ?? "TODO";
    const transitions = todo ? allowedTransitionsFor(todo.status) : [];

    return (
        <Dialog onClose={handleClose} open={open} fullWidth>
            <DialogTitle>Change Status</DialogTitle>
            <DialogContent>
                <Box sx={{ mt: 1, display: "flex", flexDirection: "column", gap: 2 }}>
                    <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                        <Typography variant="body2">Current:</Typography>
                        <Chip label={labelFor(currentStatus)} color={colorFor(currentStatus)} size="small" />
                    </Box>
                    {transitions.length === 0 ? (
                        <Typography variant="body2" color="text.secondary">
                            No further transitions allowed for this status.
                        </Typography>
                    ) : (
                        <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
                            <Typography variant="body2">Move to:</Typography>
                            <Box sx={{ display: "flex", flexWrap: "wrap", gap: 1 }}>
                                {transitions.map(target => (
                                    <Button
                                        key={target}
                                        variant="outlined"
                                        color={buttonColorFor(target)}
                                        onClick={() => onSelect(target)}>
                                        {labelFor(target)}
                                    </Button>
                                ))}
                            </Box>
                        </Box>
                    )}
                </Box>
            </DialogContent>
            <DialogActions>
                <Button onClick={handleClose}>Close</Button>
            </DialogActions>
        </Dialog>
    );
};

export default ChangeTodoStatusPopUp;
