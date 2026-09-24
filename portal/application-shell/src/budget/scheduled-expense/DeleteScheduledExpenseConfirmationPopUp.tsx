import React from "react"
import ConfirmationPopUp from "../../components/layout/ConfirmationPopUp";
import { DeleteModalMessageBundle } from "../../messages/MessageBundles";

type DeleteScheduledExpenseConfirmationPopUpProps = {
    saveCallback: () => void;
    modal: DeleteModalMessageBundle;
    open: boolean;
    handleClose: () => void;
}

const DeleteScheduledExpenseConfirmationPopUp: React.FC<DeleteScheduledExpenseConfirmationPopUpProps> = ({ saveCallback, modal, open, handleClose }) =>
    <ConfirmationPopUp confirmationHandler={saveCallback}
        handleClose={handleClose}
        open={open}
        modalMessageBody={modal.message}
        modalTitle={modal.title}
        yesLabel={modal.yesLabel}
        noLabel={modal.noLabel} />

export default DeleteScheduledExpenseConfirmationPopUp;
