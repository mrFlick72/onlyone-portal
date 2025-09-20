import { SecurityContextHolder } from "../infrastructure/securityContext";

type UserInfo = {
    given_name: string
    family_name: string
    birthdate?: string
    email: string
    phone_number?: string
}

export default class VAuthenticatorAccountRepository implements AccountRepository {
    private readonly baseUrl: string;

    public constructor(baseUrl: string) {
        this.baseUrl = baseUrl;
    }

    async findAnAccount(): Promise<Account> {
        const accessToken = SecurityContextHolder.getStore()?.token

        let userInfoResponse = await fetch(
            this.baseUrl + "/userinfo", {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${accessToken}`
            }
        })
        let jsonResponse = await userInfoResponse.json() as UserInfo
        return {
            firstName: jsonResponse.given_name,
            lastName: jsonResponse.family_name,
            birthDate: jsonResponse.birthdate ?? "",
            email: jsonResponse.email,
            phone: jsonResponse.phone_number ?? ""
        } as Account
    }

    async save(account: Account): Promise<void> {
        const accessToken = SecurityContextHolder.getStore()?.token

        await fetch(this.baseUrl + "/api/accounts", {
            method: "PUT",
            headers: {
                "Authorization": `Bearer ${accessToken}`,
                "Content-Type": "application/json",
            },
            body: JSON.stringify(account)
        })
    }

}
