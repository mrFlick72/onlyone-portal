import {AsyncLocalStorage} from "node:async_hooks";

type SecurityContext = {
    token: string;
    userName: string;
    roles: string[];
}

const SecurityContextHolder = new AsyncLocalStorage<SecurityContext>();

export {SecurityContext, SecurityContextHolder}