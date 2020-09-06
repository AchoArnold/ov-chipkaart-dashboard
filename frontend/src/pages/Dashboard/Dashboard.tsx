import React, { MouseEvent, useRef, useState, useEffect } from 'react';
import { createStyles, Theme, makeStyles } from '@material-ui/core/styles';
import Drawer from '@material-ui/core/Drawer';
import AppBar from '@material-ui/core/AppBar';
import CssBaseline from '@material-ui/core/CssBaseline';
import Toolbar from '@material-ui/core/Toolbar';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import LogoutIcon from '@material-ui/icons/ExitToApp';
import IconButton from '@material-ui/core/IconButton';
import Tooltip from '@material-ui/core/Tooltip';
import Logo from '../../components/Logo';
import useTheme from '@material-ui/core/styles/useTheme';
import { Box } from '@material-ui/core';
import ROUTE_NAMES from '../../constants/routes';
import { Link } from 'react-router-dom';
import AssessmentIcon from '@material-ui/icons/Assessment';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import CircularProgress from '@material-ui/core/CircularProgress';
import { useTranslation } from 'react-i18next';
import { ValidationErrorMessageBag } from '../../domain/ValidationErrorMessageBag';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Switch from '@material-ui/core/Switch';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import InputLabel from '@material-ui/core/InputLabel';
import Typography from '@material-ui/core/Typography';
import TableHead from '@material-ui/core/TableHead';
import Table from '@material-ui/core/Table';
import TableCell from '@material-ui/core/TableCell';
import TableBody from '@material-ui/core/TableBody';
import TableRow from '@material-ui/core/TableRow';
import { DashboardAPI } from '../../serviceProvider';
import { sendToastNotification } from '../../services/notifications';
import { useSnackbar } from 'notistack';
import { VARIANT_ERROR, VARIANT_SUCCESS } from '../../constants/errors';
import { ApiResponse } from '../../services/graphql/types';
import { AnalyzeRequest } from '../../services/graphql/generated';
import { Check } from '@material-ui/icons';
import FormHelperText from '@material-ui/core/FormHelperText';
import { localeDate } from '../../services/formatters';

const drawerWidth = 240;

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        root: {
            display: 'flex',
        },
        appBar: {
            backgroundColor: theme.palette.primary.dark,
            zIndex: theme.zIndex.drawer + 1,
        },
        drawer: {
            width: drawerWidth,
            flexShrink: 0,
        },
        drawerPaper: {
            width: drawerWidth,
        },
        drawerContainer: {
            overflow: 'auto',
        },
        content: {
            flexGrow: 1,
            padding: theme.spacing(3),
        },

        title: {
            fontSize: 20,
        },

        appBarLogo: {
            flexGrow: 1,
        },
        logo: {
            textDecoration: 'none',
        },
        form: {
            width: '100%',
            maxWidth: 600,
            margin: '0 auto',
            marginTop: theme.spacing(4),
        },

        recentRequests: {
            marginTop: theme.spacing(4),
            width: '100%',
            maxWidth: 1080,
            margin: '0 auto',
        },

        formContents: {
            '& > *': {
                marginBottom: theme.spacing(2),
            },
        },

        loadingSpinner: {
            color: theme.palette.secondary.main,
            position: 'absolute',
            left: '50%',
            marginTop: 5,
            marginLeft: -12,
        },

        fileInputLabel: {
            '& legend': {
                width: '200px !important',
            },

            '& label': {
                marginBottom: -8,
                fontSize: 12,
                paddingLeft: 12,
            },
        },

        gap1: {
            gap: theme.spacing(1) + 'px',
        },

        buttonContainer: {
            position: 'relative',
        },

        table: {
            marginTop: theme.spacing(2),
            '& thead th': {
                fontWeight: 'bold',
            },
        },

        noItemCell: {
            textAlign: 'center !important' as 'center',
            border: 'none',
        },

        buttonLogo: {
            marginLeft: 0,
            marginRight: theme.spacing(1),
        },

        badgeInProgress: {
            backgroundColor: theme.palette.grey + ' !important',
        },

        badgeDone: {
            color: theme.palette.common.white + '!important',
            backgroundColor: theme.palette.success.light + '!important',
        },

        badgeError: {
            color: theme.palette.common.white + '!important',
            backgroundColor: theme.palette.error.light + '!important',
        },
    }),
);

interface DashboardState {
    OvChipkaartUsername: string;
    OvChipkaartPassword: string;
    OvChipkaartFile?: File;
    OvChipkaartNumber: string;
    uploadTravelHistory: boolean;
    StartDate?: string;
    EndDate?: string;
    AuthorizedRequestsFetched: boolean;
    Loading: boolean;
    Errors?: ValidationErrorMessageBag;
    requestRows: Array<AnalyzeRequest>;
}

export default function Dashboard() {
    const classes = useStyles();
    const theme = useTheme();
    const { t } = useTranslation();
    const refFileInput = useRef(null);
    const { enqueueSnackbar } = useSnackbar();

    const [state, setState] = useState({
        OvChipkaartUsername: '',
        OvChipkaartPassword: '',
        OvChipkaartFile: undefined,
        OvChipkaartNumber: '',
        StartDate: undefined,
        EndDate: undefined,
        uploadTravelHistory: false,
        Loading: false,
        requestRows: [],
        Errors: undefined,
        AuthorizedRequestsFetched: false,
    } as DashboardState);

    const handleNewRequest = () => {
        let newState: DashboardState = {
            ...state,
            Errors: undefined,
            Loading: true,
        };
        setState(newState);

        DashboardAPI.storeRequest({
            ovChipkaartUsername: state.uploadTravelHistory
                ? null
                : state.OvChipkaartUsername,
            ovChipkaartPassword: state.uploadTravelHistory
                ? null
                : state.OvChipkaartPassword,
            travelHistoryFile: state.uploadTravelHistory
                ? state.OvChipkaartFile
                : null,
            ovChipkaartNumber: state.OvChipkaartNumber,
            startDate: state.StartDate ?? '',
            endDate: state.EndDate ?? '',
        })
            .then(async () => {
                sendToastNotification(
                    enqueueSnackbar,
                    'Analyze request added successfully!',
                    VARIANT_SUCCESS,
                );
                await refreshRecentRequests(newState);
                newState = resetForm(newState);
            })
            .catch((response: ApiResponse<boolean>) => {
                sendToastNotification(
                    enqueueSnackbar,
                    response.getErrorTitle(),
                    VARIANT_ERROR,
                );
                newState = {
                    ...newState,
                    Errors: response.getValidationErrors(),
                };
            })
            .finally(() => {
                setState({
                    ...newState,
                    Loading: false,
                });
            });
    };

    const refreshRecentRequests = async (state: DashboardState) => {
        await DashboardAPI.getRecentRequests()
            .then((response: ApiResponse<AnalyzeRequest[]>) => {
                setState({
                    ...state,
                    requestRows: response.getData(),
                });
            })
            .catch((response: ApiResponse<AnalyzeRequest[]>) => {
                sendToastNotification(
                    enqueueSnackbar,
                    response.getErrorTitle(),
                    VARIANT_ERROR,
                );
            });
    };

    useEffect(() => {
        refreshRecentRequests(state);
    }, []);

    const resetForm = (state: DashboardState) => {
        let form: HTMLFormElement = document.getElementById(
            'analyze-request-form',
        ) as HTMLFormElement;

        form.reset();

        return {
            ...state,
            OvChipkaartUsername: '',
            OvChipkaartPassword: '',
            OvChipkaartFile: undefined,
            OvChipkaartNumber: '',
            StartDate: undefined,
            EndDate: undefined,
        };
    };

    // @ts-ignore
    return (
        <div className={classes.root}>
            <CssBaseline />
            <AppBar position="fixed" className={classes.appBar}>
                <Toolbar>
                    <Box className={classes.appBarLogo}>
                        <Link
                            className={classes.logo}
                            to={ROUTE_NAMES.LANDING_PAGE}
                        >
                            <Logo
                                variant="small"
                                backgroundColor={theme.palette.primary.dark}
                            />
                        </Link>
                    </Box>
                    <Tooltip title="Logout">
                        <IconButton color="inherit" aria-label="logout">
                            <LogoutIcon />
                        </IconButton>
                    </Tooltip>
                </Toolbar>
            </AppBar>
            <Drawer
                className={classes.drawer}
                variant="permanent"
                classes={{
                    paper: classes.drawerPaper,
                }}
            >
                <Toolbar />
                <div className={classes.drawerContainer}>
                    <List>
                        <ListItem button key="Analyze">
                            <ListItemIcon>
                                <AssessmentIcon />
                            </ListItemIcon>
                            <ListItemText primary="Analyze" />
                        </ListItem>
                    </List>
                </div>
            </Drawer>
            <main className={classes.content}>
                <Toolbar />
                <Box width="100%">
                    <form autoComplete="off" id="analyze-request-form">
                        <Card className={classes.form} variant="outlined">
                            <CardContent className={classes.formContents}>
                                <FormControlLabel
                                    control={
                                        <Switch
                                            disabled={state.Loading}
                                            checked={state.uploadTravelHistory}
                                            onChange={(
                                                event: React.ChangeEvent<
                                                    HTMLInputElement
                                                >,
                                            ) => {
                                                setState({
                                                    ...state,
                                                    uploadTravelHistory:
                                                        event.target.checked,
                                                });
                                            }}
                                            color="primary"
                                        />
                                    }
                                    label={t('Upload travel history CSV file')}
                                />
                                {state.uploadTravelHistory && (
                                    <div className={classes.fileInputLabel}>
                                        <div>
                                            <InputLabel
                                                error={state.Errors?.has(
                                                    'travelHistoryFile',
                                                )}
                                            >
                                                {t(
                                                    'OV Chipkaart Travel History CSV *',
                                                )}
                                            </InputLabel>
                                        </div>
                                        <OutlinedInput
                                            fullWidth
                                            required
                                            type="file"
                                            ref={refFileInput}
                                            inputProps={{ accept: 'text/csv' }}
                                            key="travel-history-csv"
                                            autoComplete="off"
                                            margin="dense"
                                            error={state.Errors?.has(
                                                'travelHistoryFile',
                                            )}
                                            disabled={state.Loading}
                                            labelWidth={200}
                                            onChange={() => {
                                                setState({
                                                    ...state,
                                                    // @ts-ignore
                                                    OvChipkaartFile: refFileInput.current.getElementsByTagName(
                                                        'input',
                                                    )[0].files[0],
                                                });
                                            }}
                                        />
                                        {state.Errors?.has(
                                            'travelHistoryFile',
                                        ) && (
                                            <FormHelperText
                                                error={state.Errors?.has(
                                                    'travelHistoryFile',
                                                )}
                                            >
                                                {
                                                    state.Errors?.first(
                                                        'travelHistoryFile',
                                                    )?.message
                                                }
                                            </FormHelperText>
                                        )}
                                    </div>
                                )}
                                {!state.uploadTravelHistory && (
                                    <TextField
                                        disabled={state.Loading}
                                        required
                                        fullWidth
                                        error={state.Errors?.has(
                                            'ovChipkaartUsername',
                                        )}
                                        helperText={
                                            state.Errors?.first(
                                                'ovChipkaartUsername',
                                            )?.message
                                        }
                                        size="small"
                                        key="ovChipkaartUsername"
                                        label={t('OV Chipkaart Username')}
                                        autoComplete="off"
                                        variant="outlined"
                                        value={state.OvChipkaartUsername}
                                        onChange={(event: any) => {
                                            setState({
                                                ...state,
                                                OvChipkaartUsername:
                                                    event.target.value,
                                            });
                                        }}
                                    />
                                )}
                                {!state.uploadTravelHistory && (
                                    <TextField
                                        disabled={state.Loading}
                                        required
                                        fullWidth
                                        error={state.Errors?.has(
                                            'ovChipkaartPassword',
                                        )}
                                        helperText={
                                            state.Errors?.first(
                                                'ovChipkaartPassword',
                                            )?.message
                                        }
                                        size="small"
                                        type="password"
                                        key="ovChipkaartPassword"
                                        label={t('OV Chipkaart Password')}
                                        autoComplete="off"
                                        variant="outlined"
                                        value={state.OvChipkaartPassword}
                                        onChange={(event: any) => {
                                            setState({
                                                ...state,
                                                OvChipkaartPassword:
                                                    event.target.value,
                                            });
                                        }}
                                    />
                                )}

                                <TextField
                                    required
                                    fullWidth
                                    error={state.Errors?.has(
                                        'ovChipkaartNumber',
                                    )}
                                    helperText={
                                        state.Errors?.first('ovChipkaartNumber')
                                            ?.message
                                    }
                                    size="small"
                                    key="ovChipkaartNumber"
                                    label={t('OV Chipkaart Number')}
                                    autoComplete="off"
                                    variant="outlined"
                                    type="number"
                                    disabled={state.Loading}
                                    value={state.OvChipkaartNumber}
                                    onChange={(event: any) => {
                                        setState({
                                            ...state,
                                            OvChipkaartNumber:
                                                event.target.value,
                                        });
                                    }}
                                />

                                <Box display="flex" className={classes.gap1}>
                                    <TextField
                                        required
                                        fullWidth
                                        error={state.Errors?.has('startDate')}
                                        helperText={
                                            state.Errors?.first('startDate')
                                                ?.message
                                        }
                                        type="date"
                                        size="small"
                                        key="startDate"
                                        label={t('Start Date')}
                                        autoComplete="off"
                                        variant="outlined"
                                        disabled={state.Loading}
                                        defaultValue={state.StartDate}
                                        onChange={(event: any) => {
                                            setState({
                                                ...state,
                                                StartDate: event.target.value,
                                            });
                                        }}
                                        InputLabelProps={{
                                            shrink: true,
                                        }}
                                    />

                                    <TextField
                                        required
                                        fullWidth
                                        disabled={state.Loading}
                                        error={state.Errors?.has('endDate')}
                                        helperText={
                                            state.Errors?.first('endDate')
                                                ?.message
                                        }
                                        type="date"
                                        size="small"
                                        key="endDate"
                                        label={t('End Date')}
                                        autoComplete="off"
                                        variant="outlined"
                                        defaultValue={state.EndDate}
                                        onChange={(event: any) => {
                                            setState({
                                                ...state,
                                                EndDate: event.target.value,
                                            });
                                        }}
                                        InputLabelProps={{
                                            shrink: true,
                                        }}
                                    />
                                </Box>

                                <div className={classes.buttonContainer}>
                                    <Button
                                        fullWidth
                                        color="secondary"
                                        variant="contained"
                                        disabled={state.Loading}
                                        onClick={(event: MouseEvent) => {
                                            event.preventDefault();
                                            handleNewRequest();
                                        }}
                                    >
                                        {t('Analyze')}
                                    </Button>
                                    {state.Loading && (
                                        <CircularProgress
                                            size={24}
                                            className={classes.loadingSpinner}
                                        />
                                    )}
                                </div>
                            </CardContent>
                        </Card>
                    </form>

                    <Card className={classes.recentRequests}>
                        <CardContent>
                            <Typography variant="h4">
                                Recent Requests
                            </Typography>

                            <Table className={classes.table}>
                                <TableHead>
                                    <TableRow>
                                        <TableCell>{t('Created On')}</TableCell>
                                        <TableCell align="right">
                                            {t('Card Number')}
                                        </TableCell>
                                        <TableCell align="right">
                                            {t('Start Date')}
                                        </TableCell>
                                        <TableCell align="right">
                                            {t('End Date')}
                                        </TableCell>
                                        <TableCell align="right">
                                            {t('Status')}
                                        </TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {state.requestRows.length === 0 && (
                                        <TableRow>
                                            <TableCell
                                                colSpan={5}
                                                align="center"
                                                className={classes.noItemCell}
                                            >
                                                <Typography variant="h6">
                                                    No Request available
                                                </Typography>
                                            </TableCell>
                                        </TableRow>
                                    )}
                                    {state.requestRows.map(
                                        (row: AnalyzeRequest) => (
                                            <TableRow hover={true} key={row.id}>
                                                <TableCell
                                                    component="th"
                                                    scope="row"
                                                >
                                                    {localeDate(row.createdAt)}
                                                </TableCell>
                                                <TableCell align="right">
                                                    {row.ovChipkaartNumber
                                                        .match(/.{1,4}/g)
                                                        ?.join(' ')
                                                        .trim()}
                                                </TableCell>
                                                <TableCell align="right">
                                                    {localeDate(row.startDate)}
                                                </TableCell>
                                                <TableCell align="right">
                                                    {localeDate(row.endDate)}
                                                </TableCell>
                                                <TableCell align="right">
                                                    {row.status ===
                                                        'in-progress' && (
                                                        <Button
                                                            size="small"
                                                            variant="contained"
                                                            className={
                                                                classes.badgeInProgress
                                                            }
                                                            disabled
                                                        >
                                                            <CircularProgress
                                                                className={
                                                                    classes.buttonLogo
                                                                }
                                                                size={16}
                                                            />
                                                            In Progress
                                                        </Button>
                                                    )}

                                                    {row.status === 'done' && (
                                                        <Button
                                                            size="small"
                                                            variant="contained"
                                                            disabled
                                                            className={
                                                                classes.badgeDone
                                                            }
                                                        >
                                                            <Check />
                                                            In Progress
                                                        </Button>
                                                    )}

                                                    {row.status === 'error' && (
                                                        <Button
                                                            size="small"
                                                            variant="contained"
                                                            disabled
                                                            className={
                                                                classes.badgeDone
                                                            }
                                                        >
                                                            <Check />
                                                            In Progress
                                                        </Button>
                                                    )}
                                                </TableCell>
                                            </TableRow>
                                        ),
                                    )}
                                </TableBody>
                            </Table>
                        </CardContent>
                    </Card>
                </Box>
            </main>
        </div>
    );
}
