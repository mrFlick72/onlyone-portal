class UpdateAccount {

    private accountRepository: AccountRepository;

    constructor(accountRepository: AccountRepository) {
        this.accountRepository = accountRepository;
    }

    execute(account: Account): Promise<void> {
        return this.accountRepository.save(account);
    }
}

export default UpdateAccount;