import React from 'react';
import ComponentInitializer from "../components/ComponentInitializer";

import {useEffect} from "react";
import {authenticate} from "./Authenticator";

const Callback = () => {

    let queryString = location.search
    let params = new URLSearchParams(queryString)

    useEffect(() => {
        let code = params.get("code")!!
        authenticate(code)
            .then(_ => {
                window.location.href = window.sessionStorage.getItem("returnTo")!!
            })
    }, [])
    return <div></div>
}


ComponentInitializer(<Callback/>,)