import React from "react";
import { Box, Button, Chip, Dialog, DialogActions, DialogContent, DialogTitle, Typography } from "@mui/material";
import { Todo, TodoStatus } from "./domain/Plan";
import { allowedTransitionsFor, colorFor } from "./domain/TodoStatus";
import { ChangeTodoStatusModalMessageBundle, TodoStatusMessageBundle } from "../messages/MessageBundles";

type ChangeTodoStatusPopUpProps = {
    open: boolean;
    todo: Todo | null;
    handleClose: () => void;
    onSelect: (status: TodoStatus) => void;
    messages: ChangeTodoStatusModalMessageBundle;
    statusMessages: TodoStatusMessageBundle;
}

const buttonColorFor = (status: TodoStatus) => {
    const color = colorFor(status);
    return color === "default" ? "inherit" : color;
};

const ChangeTodoStatusPopUp: React.FC<ChangeTodoStatusPopUpProps> = ({ open, todo, handleClose, onSelect, messages, statusMessages }) => {
    const currentStatus = todo?.status ?? "TODO";
    const transitions = todo ? allowedTransitionsFor(todo.status) : [];

    return (
        <Dialog onClose={handleClose} open={open} fullWidth>
            <DialogTitle>{messages.title}</DialogTitle>
            <DialogContent>
                <Box sx={{ mt: 1, display: "flex", flexDirection: "column", gap: 2 }}>
                    <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                        <Typography variant="body2">{messages.currentLabel}</Typography>
                        <Chip label={statusMessages[currentStatus]} color={colorFor(currentStatus)} size="small" />
                    </Box>
                    {transitions.length === 0 ? (
                        <Typography variant="body2" color="text.secondary">
                            {messages.noTransitionsLabel}
                        </Typography>
                    ) : (
                        <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
                            <Typography variant="body2">{messages.moveToLabel}</Typography>
                            <Box sx={{ display: "flex", flexWrap: "wrap", gap: 1 }}>
                                {transitions.map(target => (
                                    <Button
                                        key={target}
                                        variant="outlined"
                                        color={buttonColorFor(target)}
                                        onClick={() => onSelect(target)}>
                                        {statusMessages[target]}
                                    </Button>
                                ))}
                            </Box>
                        </Box>
                    )}
                </Box>
            </DialogContent>
            <DialogActions>
                <Button onClick={handleClose}>{messages.closeButtonLabel}</Button>
            </DialogActions>
        </Dialog>
    );
};

export default ChangeTodoStatusPopUp;
