import { TodoStatus } from "../plan/domain/Plan";

/**
 * Dedicated types for every i18n label object passed around the app.
 *
 * `OnlyonePortalPagesConfigMap` turns the flat {@link MessageBundle} registry
 * (see `MessageRepository.ts`) into the per-page composite bundles at the bottom
 * of this file. Pages slice those composites and hand the leaf bundles to
 * components as `messages` / `modal` / `formMessages` / `actions` props. Both
 * the producer (config map) and the consumers (components) import these same
 * types, so a mismatch between a built label object and the prop that expects it
 * becomes a compile error.
 *
 * NOTE: the flat string registry keeps its own name, `MessageBundle`, in
 * `MessageRepository.ts`; the types here describe the structured, per-component
 * label objects derived from it.
 */

/* ── leaf bundles: consumed directly by a single component ───────────────── */

/** Base global-navigation drawer labels, consumed by `Menu`. */
export type MenuMessageBundle = {
    title: string;
    userProfileLabel: string;
    logOutLabel: string;
    budgetPageLabel: string;
    revenuePageLabel: string;
    planPageLabel: string;
    accountPageLabel: string;
    analyticsPageLabel: string;
    otherFeatureLabel: string;
};

/** Section links `Menu` derives from {@link MenuMessageBundle} for `GlobalPageNavigation`. */
export type GlobalNavMessageBundle = {
    budget: string;
    revenue: string;
    plans: string;
    account: string;
    analytics: string;
};

/** Yes/No confirm buttons (`YesAndNoButtonGroup`). */
export type ConfirmMessageBundle = {
    yesLabel: string;
    noLabel: string;
};

/* modal bundles ----------------------------------------------------------- */

/** Save/search-style dialog: a title plus save + close buttons. */
export type SaveModalMessageBundle = {
    title: string;
    saveButtonLabel: string;
    closeButtonLabel: string;
};

/** Delete-confirmation dialog: a title, a body message and yes/no buttons. */
export type DeleteModalMessageBundle = {
    title: string;
    message: string;
    yesLabel: string;
    noLabel: string;
};

/** Attachment upload dialog (`UploadAttachmentPopUp`). */
export type UploadAttachmentModalMessageBundle = {
    title: string;
    uploadButtonLabel: string;
    closeButtonLabel: string;
    chooseFileLabel: string;
    noFileSelectedLabel: string;
    existingAttachmentsLabel: string;
    noAttachmentsLabel: string;
    downloadAttachmentLabel: string;
    deleteAttachmentLabel: string;
};

/** Todo status-transition dialog (`ChangeTodoStatusPopUp`). */
export type ChangeTodoStatusModalMessageBundle = {
    title: string;
    currentLabel: string;
    moveToLabel: string;
    noTransitionsLabel: string;
    closeButtonLabel: string;
};

/* form field-label bundles ------------------------------------------------- */

export type PlanFormMessageBundle = { title: string; date: string };
export type TodoFormMessageBundle = { date: string; content: string };
export type BudgetRevenueFormMessageBundle = { date: string; amount: string; note: string };
export type SpentBudgetFormMessageBundle = { date: string; amount: string; searchTags: string; note: string };
export type SearchTagFormMessageBundle = { valueLabel: string; saveLabel: string };

/* table / row bundles ------------------------------------------------------ */

/** Localized label per todo status, keyed by the {@link TodoStatus} union. */
export type TodoStatusMessageBundle = Record<TodoStatus, string>;

export type SearchTagsTableMessageBundle = { description: string; operation: string };
export type TotalBySearchTagsMessageBundle = { category: string; total: string };

export type PlanListRowActionsMessageBundle = { open: string; delete: string };
export type TodoRowActionsMessageBundle = { changeStatus: string; edit: string; delete: string };
/** Row actions shared by budget expense and revenue rows. */
export type BudgetActionsMessageBundle = { edit: string; attach: string; delete: string };

export type PlanListContentMessageBundle = {
    headers: { date: string; title: string; todos: string; options: string };
    actions: PlanListRowActionsMessageBundle;
};

export type TodoContentMessageBundle = {
    headers: { date: string; content: string; status: string; options: string };
    actions: TodoRowActionsMessageBundle;
    status: TodoStatusMessageBundle;
};

export type BudgetRevenueTableMessageBundle = {
    headers: { date: string; amount: string; note: string; options: string };
    totalLabel: string;
    actions: BudgetActionsMessageBundle;
};

export type SpentBudgetContentMessageBundle = {
    headers: { date: string; amount: string; note: string; type: string; details: string };
    totalLabel: string;
    actions: BudgetActionsMessageBundle;
};

/* home page tiles ---------------------------------------------------------- */

export type HomeMenuTileMessageBundle = { titleText: string; body: string };
export type HomeMenuContentMessageBundle = {
    account: HomeMenuTileMessageBundle;
    budget: HomeMenuTileMessageBundle;
    revenue: HomeMenuTileMessageBundle;
    plan: HomeMenuTileMessageBundle;
    analytics: HomeMenuTileMessageBundle;
};

/* ── composite per-page bundles: returned by OnlyonePortalPagesConfigMap ──── */

/** Config-map modal entries carry a DOM `id` alongside the labels a component needs. */
type WithId<T> = T & { id: string };

export type CommonMessageBundle = {
    table: {
        date: string;
        amount: string;
        note: string;
        type: string;
        details: string;
        options: string;
        total: string;
        category: string;
        title: string;
        todos: string;
        content: string;
        status: string;
        description: string;
        operation: string;
    };
    action: { edit: string; delete: string; attach: string; open: string; changeStatus: string };
    confirm: ConfirmMessageBundle;
    form: { date: string; amount: string; note: string; searchTags: string; title: string; content: string };
    totalLabel: string;
    nav: { budget: string; revenue: string; plans: string; account: string };
};

export type HomePageMessageBundle = {
    menuMessages: MenuMessageBundle;
    homePage: { homeMenuContent: HomeMenuContentMessageBundle };
};

export type BudgetExpensePageMessageBundle = {
    menuMessages: MenuMessageBundle & {
        insertBudgetModal: string;
        searchModal: string;
        diagrams: string;
        searchTags: string;
    };
    cards: { dailyDetails: string; totalByCategories: string };
    tabs: { dailyView: string; byTagsView: string };
    newBudgetExpenseModal: WithId<SaveModalMessageBundle>;
    searchFilterModal: WithId<SaveModalMessageBundle>;
    deleteModal: WithId<DeleteModalMessageBundle>;
    uploadAttachmentModal: WithId<UploadAttachmentModalMessageBundle>;
};

export type BudgetRevenuePageMessageBundle = {
    menuMessages: MenuMessageBundle & { changeYearModal: string; searchTags: string };
    deleteModal: WithId<DeleteModalMessageBundle>;
    saveBudgetRevenueModal: WithId<SaveModalMessageBundle>;
    changeYearModal: WithId<SaveModalMessageBundle>;
    uploadAttachmentModal: WithId<UploadAttachmentModalMessageBundle>;
};

export type SearchTagsTabsMessageBundle = { expense: string; revenue: string };

export type SearchTagsPageMessageBundle = {
    menuMessages: MenuMessageBundle;
    heading: string;
    tabs: SearchTagsTabsMessageBundle;
    form: SearchTagFormMessageBundle;
};

export type PlanPageMessageBundle = {
    menuMessages: MenuMessageBundle & { newPlan: string };
    savePlanModal: WithId<SaveModalMessageBundle>;
    deletePlanModal: WithId<DeleteModalMessageBundle>;
};

export type PlanDetailPageMessageBundle = {
    menuMessages: MenuMessageBundle & { newTodo: string; backToList: string };
    statusMessages: TodoStatusMessageBundle;
    saveTodoModal: WithId<SaveModalMessageBundle>;
    updateTodoModal: WithId<SaveModalMessageBundle>;
    deleteTodoModal: WithId<DeleteModalMessageBundle>;
    changeStatusModal: WithId<ChangeTodoStatusModalMessageBundle>;
};

export type AnalyticsPageMessageBundle = {
    menuMessages: MenuMessageBundle;
    filters: {
        year: string;
        month: string;
        allMonths: string;
        tags: string;
        fromYear: string;
        toYear: string;
        tag: string;
        allTags: string;
        invalidYearRange: string;
    };
    charts: { totalByTagTitle: string; totalByYearTitle: string };
    states: { empty: string; error: string; retry: string };
    reindex: {
        title: string;
        description: string;
        button: string;
        running: string;
        success: string;
        error: string;
    };
};

export type AccountPageMessageBundle = {
    menuMessages: MenuMessageBundle;
    form: {
        firstNameLabel: string;
        lastNameLabel: string;
        birthDateLabel: string;
        emailLabel: string;
        phoneLabel: string;
        saveButtonLabel: string;
    };
};
