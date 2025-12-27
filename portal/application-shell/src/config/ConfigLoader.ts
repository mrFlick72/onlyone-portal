type ApplicationConfig = {
    scope: string,
    redirectUri: string,
    clientApplicationId: string,
    idpBaseUrl: string
    authenticationCheckInterval: number
    budgetApiBaseUrl: string
    revenueApiBaseUrl: string
    accountApiBaseUrl: string,
    tagApiBaseUrl: string

}
export const applicationConfigLoader = async () => {
    const configData = {
        scope: "openid",
        redirectUri: process.env.REDIRECT_URI,
        clientApplicationId: process.env.CLIENT_APPLICATION_ID,
        idpBaseUrl: process.env.IDP_BASE_URL,
        authenticationCheckInterval: Number(process.env.AUTHENTICATION_CHECK_INTERVAL),
        budgetApiBaseUrl: process.env.BUDGET_API_BASE_URL,
        revenueApiBaseUrl: process.env.REVENUE_API_BASE_URL,
        accountApiBaseUrl: process.env.ACCOUNT_API_BASE_URL,
        tagApiBaseUrl: process.env.TAG_API_BASE_URL

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