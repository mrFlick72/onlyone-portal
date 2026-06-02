export type MessageBundle = {
    [name: string]: string
}

export type NestedMessages = {
    [key: string]: string | NestedMessages
}

/**
 * Every message bundle YAML under `./bundle` is parsed into a nested object at
 * build time by `@modyfi/vite-plugin-yaml`. Vite needs a statically analysable
 * glob, so we eagerly load every language here and filter to the requested one
 * in {@link getAllMessageRegistry}. Keys are file paths, e.g.
 * `./bundle/account/message_bundle_en_en.yaml`.
 */
const bundleFiles = import.meta.glob<NestedMessages>("./bundle/**/*.yaml", {
    eager: true,
    import: "default",
});

/**
 * Flattens a nested i18n messages object (as served under the i18n `messages`
 * field / the bundle YAML files) into a flat `MessageBundle` with dot-separated
 * keys, e.g. `{ form: { firstName: { label: "Name" } } }` -> `{ "form.firstName.label": "Name" }`.
 */
export function flattenMessages(nested: NestedMessages, prefix = ""): MessageBundle {
    return Object.entries(nested).reduce<MessageBundle>((flat, [key, value]) => {
        const path = prefix ? `${prefix}.${key}` : key;
        if (value !== null && typeof value === "object") {
            Object.assign(flat, flattenMessages(value, path));
        } else {
            flat[path] = String(value);
        }
        return flat;
    }, {});
}

/**
 * Loads every message bundle for the given language from `./bundle` and returns
 * a single flat `MessageBundle` keyed by dot-separated paths. The bundles are
 * merged together, so callers receive one registry covering every page. The
 * language is fixed to `en_en` for now but can be overridden once more locales
 * are wired up.
 */
export function getAllMessageRegistry(language: string = "it_it"): MessageBundle {
    const suffix = `message_bundle_${language}.yaml`;
    return Object.entries(bundleFiles)
        .filter(([path]) => path.endsWith(suffix))
        .reduce<MessageBundle>(
            (registry, [, nested]) => Object.assign(registry, flattenMessages(nested)),
            {},
        );
}

export function getMessageFor(bundle: MessageBundle, key: string): string {
    return bundle[key] || "";
}
