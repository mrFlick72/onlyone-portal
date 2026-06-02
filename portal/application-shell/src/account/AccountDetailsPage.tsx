import React, { useEffect, useState } from 'react'
import { getAccountData, save } from "./domain/repository/AccountRepository";
import { Container, Paper, ThemeProvider } from "@mui/material";
import { createTheme } from '@mui/material/styles';
import CheckIcon from '@mui/icons-material/Check';
import Separator from "../components/form/Separator";
import FormButton from "../components/form/FormButton";
import { OnlyonePortalPagesConfigMap } from "../messages/OnlyonePortalPagesConfigMap";
import { getAllMessageRegistry, MessageBundle } from "../messages/MessageRepository";
import { isAuthenticated } from "../auth/Authenticator";
import Menu from "../components/menu/Menu";
import FormInputTextField from '../components/form/FormInputTextField';
import FormDatePicker, { FormDateFormatPattern } from '../components/form/FormDatePicker';
import Account from './domain/Account';

const AccountDetailsPage = () => {
    const configMap = new OnlyonePortalPagesConfigMap();
    const [account, setAccount] = useState<Account>(
        {
            email: "",
            firstName: "",
            lastName: "",
            phone: "",
            birthDate: ""
        }
    )

    const [messageRegistry, setMessageRegistry] = useState<MessageBundle>({})


    useEffect(() => {
        getAccountData().then(data => {
            setAccount(
                {
                    email: data.email,
                    firstName: data.firstName,
                    lastName: data.lastName,
                    phone: data.phone,
                    birthDate: data.birthDate
                }
            )
        })
    }, [])

    useEffect(() => {
        isAuthenticated().then()
        setMessageRegistry(getAllMessageRegistry())
    }, [])

    let theme = createTheme({
        palette: {
            primary: {
                main: '#252624',
                contrastText: '#fff',
            }
        },
    });

    return <ThemeProvider theme={theme}>
        <Paper variant="outlined">
            <Menu messages={configMap.searchTags(messageRegistry).menuMessages}
                navBarItems={[]} />

            <Container>
                <FormInputTextField id="firstName"
                    label={configMap.account(messageRegistry).form.firstNameLabel}
                    required={true}
                    handler={(value) => {
                        setAccount((account: Account) => {
                            let copiedAccount = Object.assign({}, account)
                            copiedAccount.firstName = value.target.value
                            return copiedAccount
                        })
                    }}
                    value={account.firstName || ""} />

                <FormInputTextField id="lastName"
                    label={configMap.account(messageRegistry).form.lastNameLabel}
                    required={true}
                    handler={(value) => {
                        setAccount((account: Account) => {
                            let copiedAccount = Object.assign({}, account)
                            copiedAccount.lastName = value.target.value
                            return copiedAccount
                        })
                    }}
                    value={account.lastName || ""} />

                <FormDatePicker
                    pattern={FormDateFormatPattern}
                    value={account.birthDate}
                    onClickHandler={(value) => {
                        let date = "";
                        try {
                            date = value.format(FormDateFormatPattern);
                        } catch (e) {
                        }
                        setAccount((account: Account) => {
                            let copiedAccount = Object.assign({}, account)
                            copiedAccount.birthDate = date
                            return copiedAccount
                        })
                    }}
                    label={configMap.account(messageRegistry).form.birthDateLabel} />

                <FormInputTextField id="phone"
                    label={configMap.account(messageRegistry).form.phoneLabel}
                    required={true}
                    handler={(value) => {
                        setAccount((account: Account) => {
                            let copiedAccount = Object.assign({}, account)
                            copiedAccount.phone = value.target.value
                            return copiedAccount
                        })
                    }}
                    value={account.phone || ""} />

                <FormInputTextField id="email"
                    label={configMap.account(messageRegistry).form.emailLabel}
                    required={true}
                    disabled={true}
                    handler={(value) => {
                        setAccount((account: Account) => {
                            let copiedAccount = Object.assign({}, account)
                            copiedAccount.email = value
                            return copiedAccount
                        })
                    }}
                    value={account.email || ""} />

                <Separator />

                <FormButton type="button"
                    onClickHandler={() => {
                        save({
                            "email": account.email,
                            "firstName": account.firstName,
                            "lastName": account.lastName,
                            "phone": account.phone,
                            "birthDate": account.birthDate
                        })
                    }}
                    labelPrefix={<CheckIcon fontSize="large" />}
                    label={configMap.account(messageRegistry).form.saveButtonLabel} />


            </Container>
        </Paper>
    </ThemeProvider>
}

export default AccountDetailsPage