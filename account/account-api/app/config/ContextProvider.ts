import VAuthenticatorAccountRepository from "../adapter/VAuthenticatorAccountRepository";
import UpdateAccount from "../domain/UpdateAccount";
import dotenv from "dotenv";
dotenv.config();

const accountRepository = new VAuthenticatorAccountRepository(process.env.ISSUER || "");
const updateAccount = new UpdateAccount(accountRepository);


export {accountRepository, updateAccount}
