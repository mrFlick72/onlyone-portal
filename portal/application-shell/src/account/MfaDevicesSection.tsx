import React, { useEffect, useState } from 'react'
import {
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
import { Delete } from "@mui/icons-material";
import { getMfaDevices } from "./domain/repository/MfaRepository";
import MfaDevice from "./domain/Mfa";
import { MfaDevicesSectionMessageBundle } from "../messages/MessageBundles";

type MfaDevicesSectionProps = {
    messages: MfaDevicesSectionMessageBundle
}

const methodLabelFor = (mfaMethod: MfaDevice["mfaMethod"], messages: MfaDevicesSectionMessageBundle) =>
    mfaMethod === "EMAIL_MFA_METHOD" ? messages.methodEmailLabel : messages.methodSmsLabel

const MfaDevicesSection: React.FC<MfaDevicesSectionProps> = ({ messages }) => {
    const [devices, setDevices] = useState<MfaDevice[]>([])

    useEffect(() => {
        getMfaDevices().then(setDevices)
    }, [])

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
    </>
}

export default MfaDevicesSection
