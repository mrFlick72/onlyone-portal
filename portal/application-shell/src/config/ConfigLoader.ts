type ApplicationConfig = {
    scope: string,
    redirectUri: string,
    clientApplicationId: string,
    idpBaseUrl: string
    authenticationCheckInterval: number
    budgetApiBaseUrl: string
    revenueApiBaseUrl: string
    accountApiBaseUrl: string,
    tagApiBaseUrl: string,
    planApiBaseUrl: string,

}
export const applicationConfigLoader = async () => {
    const configData = {
        scope: "openid",
        redirectUri: import.meta.env.REDIRECT_URI,
        clientApplicationId: import.meta.env.CLIENT_APPLICATION_ID,
        idpBaseUrl: import.meta.env.IDP_BASE_URL,
        authenticationCheckInterval: Number(import.meta.env.AUTHENTICATION_CHECK_INTERVAL),
        budgetApiBaseUrl: import.meta.env.BUDGET_API_BASE_URL,
        revenueApiBaseUrl: import.meta.env.REVENUE_API_BASE_URL,
        accountApiBaseUrl: import.meta.env.ACCOUNT_API_BASE_URL,
        tagApiBaseUrl: import.meta.env.TAG_API_BASE_URL,
        planApiBaseUrl: import.meta.env.PLAN_API_BASE_URL,

    };
    return configData as ApplicationConfig
}

export async function getBudgetApiBaseUrl() {
    const appConfig = await applicationConfigLoader()
    return appConfig.budgetApiBaseUrl;
}

export async function getRevenueApiBaseUrl() {
    const appConfig = await applicationConfigLoader()
    return appConfig.revenueApiBaseUrl;
}

export async function getAccountApiBaseUrl() {
    const appConfig = await applicationConfigLoader()
    return appConfig.accountApiBaseUrl;
}

export async function getTagApiBaseUrl() {
    const appConfig = await applicationConfigLoader()
    return appConfig.tagApiBaseUrl;
}

export async function getPlanApiBaseUrl() {
    const appConfig = await applicationConfigLoader()
    return appConfig.planApiBaseUrl;
}