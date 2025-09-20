import { Express, Request, Response } from "express";
import moment from "moment";
import { accountRepository, updateAccount } from "../config/ContextProvider";
import { SecurityContextHolder } from "../infrastructure/securityContext";

const ENDPOINT_PREFIX = "/api/account/user-account";

export function registerAccountEndPointFor(
    app: Express,
): void {
    app.get(ENDPOINT_PREFIX, async (req: Request, res: Response) => {
        let account = await accountRepository.findAnAccount();
        const formattedDate = getFormattedDate(account);

        res.status(200)
            .json(
                {
                    firstName: account.firstName,
                    lastName: account.lastName,
                    birthDate: formattedDate,
                    email: account.email,
                    phone: account.phone
                }
            )
            .end()
    });

    app.put(ENDPOINT_PREFIX, (req: Request, res: Response) => {
        const account = req.body as Account
        const username = SecurityContextHolder.getStore()?.userName!!
        
        const parsedBirthDate = moment(account.birthDate, "DD/MM/YYYY") || "";
        const accountToBeStored = {
            firstName: account.firstName,
            lastName: account.lastName,
            birthDate: moment(parsedBirthDate).format("YYYY-MM-DD") || "",
            email: username,
            phone: account.phone
        };
        
        updateAccount.execute(accountToBeStored)
            .then(_ => {
                res.status(204).end()
            }).catch(err => {
                console.error(err)
                res.status(500).end()
            });
    });


    const getFormattedDate = (account: Account): string => {
        let formattedDate: string = "";
        try {
            formattedDate = moment(account.birthDate, "YYYY-MM-DD").format("DD/MM/YYYY")
        } catch (e) {
        }
        return formattedDate;
    }
}