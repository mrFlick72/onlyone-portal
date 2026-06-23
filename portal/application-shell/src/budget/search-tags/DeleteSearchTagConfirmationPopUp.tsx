import React from "react"
import ConfirmationPopUp from "../../components/layout/ConfirmationPopUp";
import { DeleteModalMessageBundle } from "../../messages/MessageBundles";

type DeleteSearchTagConfirmationPopUpProps = {
    open: boolean;
    handleClose: () => void;
    saveCallback: () => void;
    modal: DeleteModalMessageBundle;
}

const DeleteSearchTagConfirmationPopUp: React.FC<DeleteSearchTagConfirmationPopUpProps> = ({ open, handleClose, saveCallback, modal }) => (
    <ConfirmationPopUp
        open={open}
        handleClose={handleClose}
        confirmationHandler={saveCallback}
        modalTitle={modal.title}
        modalMessageBody={modal.message}
        yesLabel={modal.yesLabel}
        noLabel={modal.noLabel} />
);

export default DeleteSearchTagConfirmationPopUp;
