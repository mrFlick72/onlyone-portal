import { getMessageFor, MessageBundle } from "./MessageRepository";

export class OnlyonePortalPagesConfigMap {

    home(bundle: MessageBundle) {
        return {
            menuMessages: {
                title: getMessageFor(bundle, "menu.title"),
                userProfileLabel: getMessageFor(bundle, "menu.userProfile.label"),
                logOutLabel: getMessageFor(bundle, "menu.logOut.label")
            }
        }
    }

    budgetExpense(bundle: MessageBundle) {
        return {
            menuMessages: {
                insertBudgetModal: getMessageFor(bundle, "budgetExpensePage.menu.newBudgetExpense"),
                searchModal: getMessageFor(bundle, "budgetExpensePage.menu.search"),
                diagrams: getMessageFor(bundle, "budgetExpensePage.menu.diagrams"),
                searchTags: getMessageFor(bundle, "budgetExpensePage.menu.searchTags"),
                title: getMessageFor(bundle, "menu.title"),
                userProfileLabel: getMessageFor(bundle, "menu.userProfile.label"),
                budgetPageLabel: getMessageFor(bundle, "menu.budgetPage.label"),
                revenuePageLabel: getMessageFor(bundle, "menu.revenuePage.label"),
                planPageLabel: getMessageFor(bundle, "menu.planPage.label"),
                logOutLabel: getMessageFor(bundle, "menu.logOut.label")
            },
            cards: {
                dailyDetails: getMessageFor(bundle, "budgetExpensePage.cards.dailyDetails"),
                totalByCategories: getMessageFor(bundle, "budgetExpensePage.cards.totalByCategories")
            },

            newBudgetExpenseModal: {
                id: "insertBudgetModal",
                title: getMessageFor(bundle, "budgetExpensePage.newBudgetExpense.popup.title"),
                closeButtonLabel: getMessageFor(bundle, "common.button.close.label"),
                saveButtonLabel: getMessageFor(bundle, "common.button.save.label")
            },

            searchFilterModal: {
                id: "searchByTagsModal",
                title: getMessageFor(bundle, "budgetExpensePage.search.popup.title"),
                closeButtonLabel: getMessageFor(bundle, "common.button.close.label"),
                saveButtonLabel: getMessageFor(bundle, "common.button.save.label")
            },

            deleteModal: {
                id: "deleteModal",
                title: getMessageFor(bundle, "budgetExpensePage.delete.popup.title"),
                message: getMessageFor(bundle, "budgetExpensePage.delete.popup.message")
            },

            uploadAttachmentModal: {
                id: "uploadAttachmentModal",
                title: getMessageFor(bundle, "budgetExpensePage.attachment.popup.title"),
                uploadButtonLabel: getMessageFor(bundle, "common.button.upload.label"),
                closeButtonLabel: getMessageFor(bundle, "common.button.close.label"),
                chooseFileLabel: getMessageFor(bundle, "attachment.popup.chooseFile.label"),
                noFileSelectedLabel: getMessageFor(bundle, "attachment.popup.noFileSelected.label"),
                existingAttachmentsLabel: getMessageFor(bundle, "attachment.popup.existingAttachments.label"),
                noAttachmentsLabel: getMessageFor(bundle, "attachment.popup.noAttachments.label"),
                downloadAttachmentLabel: getMessageFor(bundle, "attachment.popup.download.label"),
                deleteAttachmentLabel: getMessageFor(bundle, "attachment.popup.delete.label")
            }
        }
    }

    budgetRevenue(bundle: MessageBundle) {
        return {
            menuMessages: {
                title: getMessageFor(bundle, "menu.title"),
                userProfileLabel: getMessageFor(bundle, "menu.userProfile.label"),
                logOutLabel: getMessageFor(bundle, "menu.logOut.label"),
                changeYearModal: getMessageFor(bundle, "budgetRevenuePage.menu.changeYear")
            },
            deleteModal: { id: "deleteBudgetRevenueModal" },
            saveBudgetRevenueModal: {
                id: "saveBudgetRevenueModal",
                title: getMessageFor(bundle, "budgetRevenuePage.newBudgetRevenueModal.title"),
                closeButtonLabel: getMessageFor(bundle, "common.button.close.label"),
                saveButtonLabel: getMessageFor(bundle, "common.button.save.label")
            },
            changeYearModal: {
                id: "changeRevenueYearModal",
                title: getMessageFor(bundle, "budgetRevenuePage.changeYear.popup.title"),
                closeButtonLabel: getMessageFor(bundle, "common.button.close.label"),
                saveButtonLabel: getMessageFor(bundle, "common.button.save.label")
            },
            uploadAttachmentModal: {
                id: "uploadAttachmentModal",
                title: getMessageFor(bundle, "budgetRevenuePage.attachment.popup.title"),
                uploadButtonLabel: getMessageFor(bundle, "common.button.upload.label"),
                closeButtonLabel: getMessageFor(bundle, "common.button.close.label"),
                chooseFileLabel: getMessageFor(bundle, "attachment.popup.chooseFile.label"),
                noFileSelectedLabel: getMessageFor(bundle, "attachment.popup.noFileSelected.label"),
                existingAttachmentsLabel: getMessageFor(bundle, "attachment.popup.existingAttachments.label"),
                noAttachmentsLabel: getMessageFor(bundle, "attachment.popup.noAttachments.label"),
                downloadAttachmentLabel: getMessageFor(bundle, "attachment.popup.download.label"),
                deleteAttachmentLabel: getMessageFor(bundle, "attachment.popup.delete.label")
            }
        }
    }

    searchTags(bundle: MessageBundle) {
        return {
            menuMessages: {
                title: getMessageFor(bundle, "menu.title"),
                userProfileLabel: getMessageFor(bundle, "menu.userProfile.label"),
                logOutLabel: getMessageFor(bundle, "menu.logOut.label")
            }
        }
    }

    plan(bundle: MessageBundle) {
        return {
            menuMessages: {
                title: getMessageFor(bundle, "menu.title"),
                userProfileLabel: getMessageFor(bundle, "menu.userProfile.label"),
                logOutLabel: getMessageFor(bundle, "menu.logOut.label"),
                newPlan: getMessageFor(bundle, "planPage.menu.newPlan"),
            },
            savePlanModal: {
                id: "savePlanModal",
                title: getMessageFor(bundle, "planPage.savePlan.popup.title"),
                closeButtonLabel: getMessageFor(bundle, "common.button.close.label"),
                saveButtonLabel: getMessageFor(bundle, "common.button.save.label")
            },
            deletePlanModal: {
                id: "deletePlanModal",
                title: getMessageFor(bundle, "planPage.deletePlan.popup.title"),
                message: getMessageFor(bundle, "planPage.deletePlan.popup.message")
            }
        }
    }

    planDetail(bundle: MessageBundle) {
        return {
            menuMessages: {
                title: getMessageFor(bundle, "menu.title"),
                userProfileLabel: getMessageFor(bundle, "menu.userProfile.label"),
                logOutLabel: getMessageFor(bundle, "menu.logOut.label"),
                newTodo: getMessageFor(bundle, "planDetailPage.menu.newTodo"),
                backToList: getMessageFor(bundle, "planDetailPage.menu.backToList"),
            },
            saveTodoModal: {
                id: "saveTodoModal",
                title: getMessageFor(bundle, "planDetailPage.saveTodo.popup.title"),
                closeButtonLabel: getMessageFor(bundle, "common.button.close.label"),
                saveButtonLabel: getMessageFor(bundle, "common.button.save.label")
            },
            updateTodoModal: {
                id: "updateTodoModal",
                title: getMessageFor(bundle, "planDetailPage.updateTodo.popup.title"),
                closeButtonLabel: getMessageFor(bundle, "common.button.close.label"),
                saveButtonLabel: getMessageFor(bundle, "common.button.save.label")
            },
            deleteTodoModal: {
                id: "deleteTodoModal",
                title: getMessageFor(bundle, "planDetailPage.deleteTodo.popup.title"),
                message: getMessageFor(bundle, "planDetailPage.deleteTodo.popup.message")
            }
        }
    }

    account(bundle: MessageBundle) {
        return {
            menuMessages: {
                title: getMessageFor(bundle, "menu.title"),
                userProfileLabel: getMessageFor(bundle, "menu.userProfile.label"),
                logOutLabel: getMessageFor(bundle, "menu.logOut.label")
            },
            form: {
                firstNameLabel: getMessageFor(bundle, "form.firstName.label"),
                lastNameLabel: getMessageFor(bundle, "form.lastName.label"),
                birthDateLabel: getMessageFor(bundle, "form.birthDate.label"),
                emailLabel: getMessageFor(bundle, "form.email.value"),
                phoneLabel: getMessageFor(bundle, "form.phone.value"),
                saveButtonLabel: getMessageFor(bundle, "form.save.value"),

            }
        }
    }
}