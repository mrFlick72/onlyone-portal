/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly REDIRECT_URI: string;
  readonly CLIENT_APPLICATION_ID: string;
  readonly IDP_BASE_URL: string;
  readonly AUTHENTICATION_CHECK_INTERVAL: string;
  readonly BUDGET_API_BASE_URL: string;
  readonly REVENUE_API_BASE_URL: string;
  readonly ACCOUNT_API_BASE_URL: string;
  readonly TAG_API_BASE_URL: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
