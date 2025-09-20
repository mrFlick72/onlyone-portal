interface AccountRepository {

    findAnAccount(): Promise<Account>;

    save(account: Account): Promise<void>;
}
