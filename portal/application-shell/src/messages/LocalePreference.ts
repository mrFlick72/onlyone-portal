const STORAGE_KEY = "PREFERRED_LOCALE";

export const SUPPORTED_LOCALES = ["it", "en"];
export const DEFAULT_BUNDLE_LOCALE = "it_it";

export function getCachedLocale(): string | null {
    try {
        return window.localStorage.getItem(STORAGE_KEY);
    } catch {
        return null;
    }
}

export function setCachedLocale(locale: string): void {
    try {
        window.localStorage.setItem(STORAGE_KEY, locale);
    } catch {
        // localStorage unavailable (private browsing, blocked site data, ...) -- the app
        // just falls back to browser detection / the default on the next page load.
    }
}

export function toBundleLocale(locale: string): string {
    return `${locale}_${locale}`;
}

function detectBrowserLocale(): string | null {
    const languages = (typeof navigator !== "undefined" && navigator.languages) || [];
    for (const language of languages) {
        const primarySubtag = language.split("-")[0].toLowerCase();
        if (SUPPORTED_LOCALES.includes(primarySubtag)) {
            return primarySubtag;
        }
    }
    return null;
}

/**
 * Resolves the bundle suffix (e.g. "it_it") getAllMessageRegistry should load: the locale
 * cached at login, else a browser-detected guess (never persisted back to vauthenticator),
 * else the hardcoded default.
 */
export function resolveBundleLocale(): string {
    const cached = getCachedLocale();
    if (cached && SUPPORTED_LOCALES.includes(cached)) {
        return toBundleLocale(cached);
    }

    const detected = detectBrowserLocale();
    if (detected) {
        return toBundleLocale(detected);
    }

    return DEFAULT_BUNDLE_LOCALE;
}
