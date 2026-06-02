import React from "react"
import ConfirmationPopUp from "../components/layout/ConfirmationPopUp";
import { DeleteModalMessageBundle } from "../messages/MessageBundles";

interface DeletePlanConfirmationPopUpProps {
    open: boolean;
    handleClose: () => void;
    saveCallback: () => void;
    modal: DeleteModalMessageBundle;
}

const DeletePlanConfirmationPopUp: React.FC<DeletePlanConfirmationPopUpProps> = ({ open, handleClose, saveCallback, modal }) => (
    <ConfirmationPopUp
        open={open}
        handleClose={handleClose}
        confirmationHandler={saveCallback}
        modalTitle={modal.title}
        modalMessageBody={modal.message}
        yesLabel={modal.yesLabel}
        noLabel={modal.noLabel} />
);

export default DeletePlanConfirmationPopUp;
