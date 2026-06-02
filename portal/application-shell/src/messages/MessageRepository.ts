import { getI18nApiBaseUrl } from "../config/ConfigLoader";

export type MessageBundle = {
    [name: string]: string
}

export type NestedMessages = {
    [key: string]: string | NestedMessages
}

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

export function getAllMessageRegistry(): MessageBundle {
    let bundle = {}
    getI18nApiBaseUrl()
        .then(url => {
            fetch(url)
                .then(response => response.json()
                    .then(content => {
                        bundle = flattenMessages(content.messages)
                    }))
        })

    return bundle
}

export function getMessageFor(bundle: MessageBundle, key: string): string {
    return bundle[key] || "";
}
