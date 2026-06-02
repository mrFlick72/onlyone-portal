import React from "react"
import ConfirmationPopUp from "../../../components/layout/ConfirmationPopUp";
import { DeleteModalMessageBundle } from "../../../messages/MessageBundles";
type DeleteBudgetExpenseConfirmationPopUpProps = {
    deleteBudgetExpenseAction: () => void,
    modal: DeleteModalMessageBundle,
    open: boolean,
    handleClose: () => void
}

const DeleteBudgetExpenseConfirmationPopUp: React.FC<DeleteBudgetExpenseConfirmationPopUpProps> = ({ deleteBudgetExpenseAction, modal, open, handleClose }) =>
    <ConfirmationPopUp confirmationHandler={deleteBudgetExpenseAction}
        handleClose={handleClose}
        open={open}
        modalMessageBody={modal.message}
        modalTitle={modal.title}
        yesLabel={modal.yesLabel}
        noLabel={modal.noLabel} />

export default DeleteBudgetExpenseConfirmationPopUp