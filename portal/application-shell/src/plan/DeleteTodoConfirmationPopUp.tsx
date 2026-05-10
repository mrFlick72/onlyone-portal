import React from "react"
import ConfirmationPopUp from "../components/layout/ConfirmationPopUp";

interface DeleteTodoConfirmationPopUpProps {
    open: boolean;
    handleClose: () => void;
    saveCallback: () => void;
    modal: {
        title: string;
        message: string;
    };
}

const DeleteTodoConfirmationPopUp: React.FC<DeleteTodoConfirmationPopUpProps> = ({ open, handleClose, saveCallback, modal }) => (
    <ConfirmationPopUp
        open={open}
        handleClose={handleClose}
        confirmationHandler={saveCallback}
        modalTitle={modal.title}
        modalMessageBody={modal.message} />
);

export default DeleteTodoConfirmationPopUp;
