import React from "react"
import { Dialog, DialogActions, DialogContent, DialogTitle } from "@mui/material";
import { Check } from "@mui/icons-material";
import moment from "moment";
import TodoForm, { TodoFormData } from "./TodoForm";
import YesAndNoButtonGroup from "../components/layout/YesAndNoButtonGroup";

interface SaveTodoPopUpProps {
    open: boolean;
    handleClose: () => void;
    todo: TodoFormData;
    handlers: {
        content: (event: React.ChangeEvent<HTMLTextAreaElement>) => void;
        date: (date: moment.Moment) => void;
    };
    saveCallback: () => void;
    modal: {
        title: string;
        saveButtonLabel: string;
        closeButtonLabel: string;
    };
    formMessages: { date: string; content: string };
}

const SaveTodoPopUp: React.FC<SaveTodoPopUpProps> = ({ open, handleClose, todo, handlers, saveCallback, modal, formMessages }) => (
    <Dialog onClose={handleClose} open={open} fullWidth scroll="paper">
        <DialogTitle>{modal.title}</DialogTitle>
        <DialogContent>
            <TodoForm data={todo} handlers={handlers} messages={formMessages} />
        </DialogContent>
        <DialogActions>
            <YesAndNoButtonGroup
                yesIcon={<Check />}
                yesFun={saveCallback}
                noFun={handleClose}
                buttonMessages={{
                    noLabel: modal.closeButtonLabel,
                    yesLabel: modal.saveButtonLabel,
                }} />
        </DialogActions>
    </Dialog>
);

export default SaveTodoPopUp;
