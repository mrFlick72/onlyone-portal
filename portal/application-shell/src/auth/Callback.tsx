import React from 'react';
import ComponentInitializer from "../components/ComponentInitializer";

import {useEffect} from "react";
import {authenticate} from "./Authenticator";
import {getAccountData} from "../account/domain/repository/AccountRepository";
import {setCachedLocale} from "../messages/LocalePreference";

const Callback = () => {

    let queryString = location.search
    let params = new URLSearchParams(queryString)

    useEffect(() => {
        let code = params.get("code")!!
        authenticate(code)
            .then(() =>
                getAccountData()
                    .then(data => {
                        if (data.locale) {
                            setCachedLocale(data.locale)
                        }
                    })
                    .catch(() => {
                        // Caching the language is best-effort -- login must not be blocked
                        // by this call failing; the app falls back to browser detection.
                    })
            )
            .then(() => {
                window.location.href = window.sessionStorage.getItem("returnTo")!!
            })
    }, [])
    return <div></div>
}


ComponentInitializer(<Callback/>,)