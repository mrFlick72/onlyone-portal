import { NextFunction, Request, Response } from 'express';
import { createRemoteJWKSet, jwtVerify } from "jose";
import { SecurityContextHolder } from "./securityContext";

const oAuth2ResourceServerMiddlewareFactory = (issuer: string) => {
    return async (req: Request, res: Response, next: NextFunction) => {
        const JWKS = createRemoteJWKSet(new URL(`${issuer}/oauth2/jwks`))

        const options = {
            issuer: issuer,
        }

        const token = req.headers.authorization?.slice(7, req.headers.authorization?.length) || ""
        await jwtVerify(token, JWKS, options)
            .then(result => {                
                const userName = result.payload["user_name"] as string;
                const roles = result.payload["authorities"] as string[]

                const ALLOWED_ROLE = process.env.ALLOWED_ROLE;
                const hasPermision: boolean = roles.some((role, _) => {
                    console.log("role " + role)
                    console.log("role === ALLOWED_ROLE " + role !== ALLOWED_ROLE)
                    return role === ALLOWED_ROLE
                })

                if (hasPermision) {
                    SecurityContextHolder.run({
                        token: token,
                        userName: userName,
                        roles: roles
                    }, () => next())

                } else {
                    console.error(`The provided roles ${roles}, does not match with the expected role: ${ALLOWED_ROLE}`)
                    res.status(403).end()
                }
            })
            .catch(async (error) => {
                console.error(error)
                SecurityContextHolder.run({
                    token: token,
                    userName: "NOP",
                    roles: ["NOP"] as string[]
                }, () => res.status(401)
                    .end())
            })
    }
}

export default oAuth2ResourceServerMiddlewareFactory