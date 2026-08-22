import React, { useCallback, useEffect, useState } from 'react'
import {
    Button,
    Chip,
    IconButton,
    Paper,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Tooltip,
    Typography
} from "@mui/material";
import { Add, Delete } from "@mui/icons-material";
import { getMfaDevices } from "./domain/repository/MfaRepository";
import MfaDevice from "./domain/Mfa";
import { AccountPageMessageBundle } from "../messages/MessageBundles";
import EnrollMfaDeviceDialog from "./EnrollMfaDeviceDialog";

type MfaDevicesSectionProps = {
    email: string;
    phone: string;
    messages: AccountPageMessageBundle["mfaDevices"];
    enrollDialogMessages: AccountPageMessageBundle["enrollDialog"];
}

const methodLabelFor = (mfaMethod: MfaDevice["mfaMethod"], messages: MfaDevicesSectionProps["messages"]) =>
    mfaMethod === "EMAIL_MFA_METHOD" ? messages.methodEmailLabel : messages.methodSmsLabel

const MfaDevicesSection: React.FC<MfaDevicesSectionProps> = ({ email, phone, messages, enrollDialogMessages }) => {
    const [devices, setDevices] = useState<MfaDevice[]>([])
    const [dialogOpen, setDialogOpen] = useState(false)

    const loadDevices = useCallback(() => {
        getMfaDevices().then(setDevices)
    }, [])

    useEffect(() => { loadDevices() }, [loadDevices])

    return <>
        <Typography variant="h6">{messages.title}</Typography>

        {devices.length === 0
            ? <Typography>{messages.emptyStateLabel}</Typography>
            : <TableContainer component={Paper}>
                <Table aria-label="mfa devices table">
                    <TableHead>
                        <TableRow>
                            <TableCell>{messages.methodColumnLabel}</TableCell>
                            <TableCell>{messages.channelColumnLabel}</TableCell>
                            <TableCell>{messages.defaultColumnLabel}</TableCell>
                            <TableCell>{messages.actionsColumnLabel}</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {devices.map(device => (
                            <TableRow key={device.mfaDeviceId}>
                                <TableCell scope="row">{methodLabelFor(device.mfaMethod, messages)}</TableCell>
                                <TableCell>{device.mfaChannel}</TableCell>
                                <TableCell>
                                    {device.default && <Chip label={messages.defaultChipLabel} color="primary" size="small" />}
                                </TableCell>
                                <TableCell>
                                    <Tooltip title={messages.deleteDisabledTooltip}>
                                        <span>
                                            <IconButton disabled aria-label={messages.deleteActionLabel}>
                                                <Delete />
                                            </IconButton>
                                        </span>
                                    </Tooltip>
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>
        }

        <Button variant="contained" startIcon={<Add />} onClick={() => setDialogOpen(true)} sx={{ marginTop: 2 }}>
            {messages.addDeviceLabel}
        </Button>

        <EnrollMfaDeviceDialog
            open={dialogOpen}
            onClose={() => setDialogOpen(false)}
            onEnrolled={() => {
                setDialogOpen(false)
                loadDevices()
            }}
            email={email}
            phone={phone}
            existingDevices={devices}
            messages={enrollDialogMessages} />
    </>
}

export default MfaDevicesSection
