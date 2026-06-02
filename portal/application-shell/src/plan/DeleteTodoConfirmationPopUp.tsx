import React from "react"
import ConfirmationPopUp from "../components/layout/ConfirmationPopUp";
import { DeleteModalMessageBundle } from "../messages/MessageBundles";

interface DeleteTodoConfirmationPopUpProps {
    open: boolean;
    handleClose: () => void;
    saveCallback: () => void;
    modal: DeleteModalMessageBundle;
}

const DeleteTodoConfirmationPopUp: React.FC<DeleteTodoConfirmationPopUpProps> = ({ open, handleClose, saveCallback, modal }) => (
    <ConfirmationPopUp
        open={open}
        handleClose={handleClose}
        confirmationHandler={saveCallback}
        modalTitle={modal.title}
        modalMessageBody={modal.message}
        yesLabel={modal.yesLabel}
        noLabel={modal.noLabel} />
);

export default DeleteTodoConfirmationPopUp;
