import { getMessageFor, MessageBundle } from "./MessageRepository";
import {
    AccountPageMessageBundle,
    BudgetExpensePageMessageBundle,
    BudgetRevenuePageMessageBundle,
    CommonMessageBundle,
    HomePageMessageBundle,
    MenuMessageBundle,
    PlanDetailPageMessageBundle,
    PlanPageMessageBundle,
    SearchTagsPageMessageBundle,
} from "./MessageBundles";

/**
 * Labels shown in the global navigation drawer (rendered by `Menu` ->
 * `GlobalPageNavigation` / `MenuSectionTitle`). They are added to every page's
 * `menuMessages` so the drawer is fully localized regardless of which page is
 * currently mounted.
 */
function globalNavLabels(bundle: MessageBundle) {
    return {
        budgetPageLabel: getMessageFor(bundle, "menu.budgetPage.label"),
        revenuePageLabel: getMessageFor(bundle, "menu.revenuePage.label"),
        planPageLabel: getMessageFor(bundle, "menu.planPage.label"),
        accountPageLabel: getMessageFor(bundle, "menu.accountPage.label"),
        otherFeatureLabel: getMessageFor(bundle, "menu.otherFeature.label"),
    };
}

function baseMenuMessages(bundle: MessageBundle): MenuMessageBundle {
    return {
        title: getMessageFor(bundle, "menu.title"),
        userProfileLabel: getMessageFor(bundle, "menu.userProfile.label"),
        logOutLabel: getMessageFor(bundle, "menu.logOut.label"),
        ...globalNavLabels(bundle),
    };
}

export class OnlyonePortalPagesConfigMap {

    /**
     * Shared labels reused across pages: table headers, row action buttons,
     * confirmation buttons and form field labels. Pages pull the group they
     * need and pass it down to the relevant leaf component.
     */
    common(bundle: MessageBundle): CommonMessageBundle {
        return {
            table: {
                date: getMessageFor(bundle, "common.table.date"),
                amount: getMessageFor(bundle, "common.table.amount"),
                note: getMessageFor(bundle, "common.table.note"),
                type: getMessageFor(bundle, "common.table.type"),
                details: getMessageFor(bundle, "common.table.details"),
                options: getMessageFor(bundle, "common.table.options"),
                total: getMessageFor(bundle, "common.table.total"),
                category: getMessageFor(bundle, "common.table.category"),
                title: getMessageFor(bundle, "common.table.title"),
                todos: getMessageFor(bundle, "common.table.todos"),
                content: getMessageFor(bundle, "common.table.content"),
                status: getMessageFor(bundle, "common.table.status"),
                description: getMessageFor(bundle, "common.table.description"),
                operation: getMessageFor(bundle, "common.table.operation"),
            },
            action: {
                edit: getMessageFor(bundle, "common.action.edit"),
                delete: getMessageFor(bundle, "common.action.delete"),
                attach: getMessageFor(bundle, "common.action.attach"),
                open: getMessageFor(bundle, "common.action.open"),
                changeStatus: getMessageFor(bundle, "common.action.changeStatus"),
            },
            confirm: {
                yesLabel: getMessageFor(bundle, "common.confirm.yesLabel"),
                noLabel: getMessageFor(bundle, "common.confirm.noLabel"),
            },
            form: {
                date: getMessageFor(bundle, "common.form.date"),
                amount: getMessageFor(bundle, "common.form.amount"),
                note: getMessageFor(bundle, "common.form.note"),
                searchTags: getMessageFor(bundle, "common.form.searchTags"),
                title: getMessageFor(bundle, "common.form.title"),
                content: getMessageFor(bundle, "common.form.content"),
            },
            totalLabel: getMessageFor(bundle, "common.total.label"),
            nav: {
                budget: getMessageFor(bundle, "common.nav.budget"),
                revenue: getMessageFor(bundle, "common.nav.revenue"),
                plans: getMessageFor(bundle, "common.nav.plans"),
                account: getMessageFor(bundle, "common.nav.account"),
            },
        }
    }

    home(bundle: MessageBundle): HomePageMessageBundle {
        return {
            menuMessages: baseMenuMessages(bundle),
            homePage: {
                homeMenuContent: {
                    account: {
                        titleText: getMessageFor(bundle, "homePage.homeMenuContent.account.titleText"),
                        body: getMessageFor(bundle, "homePage.homeMenuContent.account.body"),
                    },
                    budget: {
                        titleText: getMessageFor(bundle, "homePage.homeMenuContent.budget.titleText"),
                        body: getMessageFor(bundle, "homePage.homeMenuContent.budget.body"),
                    },
                    revenue: {
                        titleText: getMessageFor(bundle, "homePage.homeMenuContent.revenue.titleText"),
                        body: getMessageFor(bundle, "homePage.homeMenuContent.revenue.body"),
                    },
                    plan: {
                        titleText: getMessageFor(bundle, "homePage.homeMenuContent.plan.titleText"),
                        body: getMessageFor(bundle, "homePage.homeMenuContent.plan.body"),
                    }
                }
            }
        }
    }

    budgetExpense(bundle: MessageBundle): BudgetExpensePageMessageBundle {
        return {
            menuMessages: {
                ...baseMenuMessages(bundle),
                insertBudgetModal: getMessageFor(bundle, "budgetExpensePage.menu.newBudgetExpense"),
                searchModal: getMessageFor(bundle, "budgetExpensePage.menu.search"),
                diagrams: getMessageFor(bundle, "budgetExpensePage.menu.diagrams"),
                searchTags: getMessageFor(bundle, "budgetExpensePage.menu.searchTags"),
            },
            cards: {
                dailyDetails: getMessageFor(bundle, "budgetExpensePage.cards.dailyDetails"),
                totalByCategories: getMessageFor(bundle, "budgetExpensePage.cards.totalByCategories")
            },

            tabs: {
                dailyView: getMessageFor(bundle, "budgetExpensePage.tabs.dailyView"),
                byTagsView: getMessageFor(bundle, "budgetExpensePage.tabs.byTagsView")
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
                message: getMessageFor(bundle, "budgetExpensePage.delete.popup.message"),
                yesLabel: getMessageFor(bundle, "common.confirm.yesLabel"),
                noLabel: getMessageFor(bundle, "common.confirm.noLabel")
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

    budgetRevenue(bundle: MessageBundle): BudgetRevenuePageMessageBundle {
        return {
            menuMessages: {
                ...baseMenuMessages(bundle),
                changeYearModal: getMessageFor(bundle, "budgetRevenuePage.menu.changeYear")
            },
            deleteModal: {
                id: "deleteBudgetRevenueModal",
                title: getMessageFor(bundle, "budgetRevenuePage.delete.popup.title"),
                message: getMessageFor(bundle, "budgetRevenuePage.delete.popup.message"),
                yesLabel: getMessageFor(bundle, "common.confirm.yesLabel"),
                noLabel: getMessageFor(bundle, "common.confirm.noLabel")
            },
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

    searchTags(bundle: MessageBundle): SearchTagsPageMessageBundle {
        return {
            menuMessages: baseMenuMessages(bundle),
            form: {
                valueLabel: getMessageFor(bundle, "searchTagsPage.form.value.label"),
                saveLabel: getMessageFor(bundle, "searchTagsPage.form.save.label")
            }
        }
    }

    plan(bundle: MessageBundle): PlanPageMessageBundle {
        return {
            menuMessages: {
                ...baseMenuMessages(bundle),
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
                message: getMessageFor(bundle, "planPage.deletePlan.popup.message"),
                yesLabel: getMessageFor(bundle, "common.confirm.yesLabel"),
                noLabel: getMessageFor(bundle, "common.confirm.noLabel")
            }
        }
    }

    planDetail(bundle: MessageBundle): PlanDetailPageMessageBundle {
        return {
            menuMessages: {
                ...baseMenuMessages(bundle),
                newTodo: getMessageFor(bundle, "planDetailPage.menu.newTodo"),
                backToList: getMessageFor(bundle, "planDetailPage.menu.backToList"),
            },
            statusMessages: {
                TODO: getMessageFor(bundle, "planDetailPage.status.todo"),
                IN_PROGRESS: getMessageFor(bundle, "planDetailPage.status.inProgress"),
                DONE: getMessageFor(bundle, "planDetailPage.status.done"),
                ABORTED: getMessageFor(bundle, "planDetailPage.status.aborted"),
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
                message: getMessageFor(bundle, "planDetailPage.deleteTodo.popup.message"),
                yesLabel: getMessageFor(bundle, "common.confirm.yesLabel"),
                noLabel: getMessageFor(bundle, "common.confirm.noLabel")
            },
            changeStatusModal: {
                id: "changeStatusModal",
                title: getMessageFor(bundle, "planDetailPage.changeStatus.popup.title"),
                currentLabel: getMessageFor(bundle, "planDetailPage.changeStatus.popup.current"),
                moveToLabel: getMessageFor(bundle, "planDetailPage.changeStatus.popup.moveTo"),
                noTransitionsLabel: getMessageFor(bundle, "planDetailPage.changeStatus.popup.noTransitions"),
                closeButtonLabel: getMessageFor(bundle, "common.button.close.label")
            }
        }
    }

    account(bundle: MessageBundle): AccountPageMessageBundle {
        return {
            menuMessages: baseMenuMessages(bundle),
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
