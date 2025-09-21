type ApplicationConfig = {
    scope: string,
    redirectUri: string,
    clientApplicationId: string,
    idpBaseUrl: string
    authenticationCheckInterval: number
    budgetApiBaseUrl: string
    accountApiBaseUrl: string

}
export const applicationConfigLoader = async () => {
    const configData = {
        scope: "openid",
        redirectUri: process.env.REDIRECT_URI,
        clientApplicationId: process.env.CLIENT_APPLICATION_ID,
        idpBaseUrl: process.env.IDP_BASE_URL,
        authenticationCheckInterval: Number(process.env.AUTHENTICATION_CHECK_INTERVAL),
        budgetApiBaseUrl: process.env.BUDGET_API_BASE_URL,
        accountApiBaseUrl: process.env.ACCOUNT_API_BASE_URL
    };
    return configData as ApplicationConfig
}

export async function getBudgetApiBaseUrl() {
    const appConfig = await applicationConfigLoader()
    return appConfig.budgetApiBaseUrl;
}

export async function getAccountApiBaseUrl() {
    const appConfig = await applicationConfigLoader()
    return appConfig.accountApiBaseUrl;
}