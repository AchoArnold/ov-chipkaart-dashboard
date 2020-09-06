import React, { useState, MouseEvent, ChangeEvent } from 'react';
import Grid from '@material-ui/core/Grid';
import { makeStyles, createStyles, Theme } from '@material-ui/core/styles';
import Box from '@material-ui/core/Box';
import ROUTE_NAMES, { ROUTE_DASHBOARD } from '../../constants/routes';
import { Link, useHistory } from 'react-router-dom';
import Logo from '../../components/Logo';
import Typography from '@material-ui/core/Typography';
import { useTranslation } from 'react-i18next';
import TransKeys from '../../i18n/keys';
import Button from '@material-ui/core/Button';
import Anchor from '@material-ui/core/Link';
import TextField from '@material-ui/core/TextField';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Checkbox from '@material-ui/core/Checkbox';
import GoogleInvisibleCaptcha from '../../components/GoogleInvisibleCaptcha';
import useTheme from '@material-ui/core/styles/useTheme';
import CircularProgress from '@material-ui/core/CircularProgress';
import { LandingPageAPI } from '../../serviceProvider';
import { ValidationErrorMessageBag } from '../../domain/ValidationErrorMessageBag';
import {
    CreateUserResponse,
    LoginResponse,
} from '../../services/graphql/types';
import { useSnackbar } from 'notistack';
import { VARIANT_ERROR, VARIANT_SUCCESS } from '../../constants/errors';
import { sendToastNotification } from '../../services/notifications';
import { KEY_TOKEN } from '../../constants/localStorage';

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        container: {
            height: '100vh',
            width: '100vw',
            overflow: 'hidden',
        },

        info: {
            height: '100vh',
            padding: theme.spacing(2),
            display: 'flex',
            alignItems: 'center',
        },

        infoContainer: {
            height: '100vh',
            backgroundSize: `cover`,
            backgroundColor: theme.palette.primary.dark,
        },

        card: {
            color: theme.palette.primary.contrastText,
        },

        logo: {
            color: theme.palette.primary.dark,
            position: 'absolute',
            textDecoration: 'none',
        },

        title: {
            fontWeight: 'bold',
        },

        form: {
            '& > *': {
                marginBottom: theme.spacing(2),
            },
        },

        authContainer: {
            position: 'relative',
        },

        authHeader: {
            position: 'absolute',
            width: '100%',
        },

        auth: {
            height: '100vh',
            display: 'flex',
            alignItems: 'center',
            padding: theme.spacing(3),
            backgroundColor: theme.palette.primary.contrastText,
        },

        loadingSpinner: {
            color: theme.palette.secondary.main,
            position: 'absolute',
            left: '50%',
            marginTop: 5,
            marginLeft: -12,
        },
    }),
);

interface LocalState {
    signUpActive: boolean;
    loading: boolean;
    email: string;
    password: string;
    firstName: string;
    lastName: string;
    rememberMe: boolean;
    loginErrors?: ValidationErrorMessageBag;
    signUpErrors?: ValidationErrorMessageBag;
}

export default function LandingPage() {
    const classes = useStyles();
    const { t } = useTranslation();
    const theme = useTheme();
    const history = useHistory();
    const { enqueueSnackbar } = useSnackbar();
    const [state, setState] = useState({
        signUpActive: true,
        loading: false,
        email: '',
        password: '',
        firstName: '',
        lastName: '',
        rememberMe: true,
        loginErrors: undefined,
        signUpErrors: undefined,
    } as LocalState);

    const handleLogin = () => {
        let newState = beforeApiRequest();
        let requiresStateUpdate = true;
        LandingPageAPI.login({
            email: state.email,
            password: state.password,
            rememberMe: state.rememberMe,
            reCaptcha: 'how are you',
        })
            .then((response: LoginResponse) => {
                sendToastNotification(
                    enqueueSnackbar,
                    response.getErrorTitle(),
                    VARIANT_ERROR,
                );
                if (response.hasValidationErrors()) {
                    newState = {
                        ...newState,
                        loginErrors: response.getValidationErrors(),
                    };
                    return;
                }

                if (response.isValid() && response.getData()?.token.value) {
                    localStorage.setItem(
                        KEY_TOKEN,
                        response.getData()?.token.value ?? '',
                    );

                    sendToastNotification(
                        enqueueSnackbar,
                        t('Login successful!'),
                        VARIANT_SUCCESS,
                    );
                    requiresStateUpdate = false;
                    history.push(ROUTE_DASHBOARD);
                }
            })
            .finally(() => {
                if (requiresStateUpdate) {
                    setState({
                        ...newState,
                        loading: false,
                    });
                }
            });
    };

    const handleSignUp = () => {
        let newState = beforeApiRequest();
        LandingPageAPI.signUp({
            email: state.email,
            password: state.password,
            firstName: state.firstName,
            lastName: state.lastName,
            reCaptcha: 'how are you',
        })
            .then((response: CreateUserResponse) => {
                sendToastNotification(
                    enqueueSnackbar,
                    response.getErrorTitle(),
                    VARIANT_ERROR,
                );
                if (response.hasValidationErrors()) {
                    newState = {
                        ...newState,
                        signUpErrors: response.getValidationErrors(),
                    };
                    return;
                }

                if (response.isValid()) {
                    sendToastNotification(
                        enqueueSnackbar,
                        t('Sign up successful!'),
                        VARIANT_SUCCESS,
                    );
                }
            })
            .finally(() => {
                setState({
                    ...newState,
                    loading: false,
                });
            });
    };

    const beforeApiRequest = (): LocalState => {
        let newState = {
            ...state,
            loginErrors: undefined,
            signUpErrors: undefined,
            loading: true,
        };

        setState(newState);

        return newState;
    };

    return (
        <Grid container className={classes.container}>
            <Grid item xs={8}>
                <Box
                    width="100%"
                    className={classes.infoContainer}
                    display="flex"
                    justifyContent="flex-end"
                >
                    <Box width="100%" maxWidth={1000}>
                        <Link
                            className={classes.logo}
                            to={ROUTE_NAMES.LANDING_PAGE}
                        >
                            <Logo
                                backgroundColor={theme.palette.primary.dark}
                            />
                        </Link>
                        <Box className={classes.info}>
                            <Box
                                className={classes.card}
                                width="70%"
                                maxWidth={375}
                                marginTop="300"
                            >
                                <Typography
                                    variant="h2"
                                    className={classes.title}
                                >
                                    {t(TransKeys.LANDING_PAGE.TITLE)}
                                </Typography>
                                <Typography variant="subtitle1" gutterBottom>
                                    {t(TransKeys.LANDING_PAGE.SUB_TITLE)}
                                </Typography>
                            </Box>
                        </Box>
                    </Box>
                </Box>
            </Grid>
            <Grid item xs={4}>
                <Box
                    width="100%"
                    className={classes.authContainer}
                    maxWidth={400}
                >
                    <Box
                        p={3}
                        display="flex"
                        className={classes.authHeader}
                        justifyContent="flex-end"
                    >
                        <Button
                            variant="outlined"
                            onClick={() => {
                                setState({
                                    ...state,
                                    signUpActive: !state.signUpActive,
                                });
                            }}
                            color="primary"
                        >
                            {state.signUpActive ? t('Sign In') : t('Sign Up')}
                        </Button>
                    </Box>
                    <Box className={classes.auth}>
                        {state.signUpActive ? (
                            <Box width="100%">
                                <Typography variant="h5">
                                    {t('Sign Up')}
                                </Typography>
                                <Typography variant="body2">
                                    {t('or')}{' '}
                                    <Anchor
                                        href="#"
                                        onClick={(event: MouseEvent) => {
                                            event.preventDefault();
                                            setState({
                                                ...state,
                                                signUpActive: false,
                                            });
                                        }}
                                    >
                                        {t('sign in to your account')}
                                    </Anchor>
                                </Typography>
                                <form>
                                    <Box
                                        mt={2}
                                        width="100%"
                                        className={classes.form}
                                    >
                                        <TextField
                                            required
                                            fullWidth
                                            error={state.signUpErrors?.has(
                                                'firstName',
                                            )}
                                            helperText={
                                                state.signUpErrors?.first(
                                                    'firstName',
                                                )?.message
                                            }
                                            key="firstName"
                                            name="firstName"
                                            size="small"
                                            label={t('First Name')}
                                            variant="outlined"
                                            autoComplete="given-name"
                                            value={state.firstName}
                                            onChange={(event: any) => {
                                                setState({
                                                    ...state,
                                                    firstName:
                                                        event.target.value,
                                                });
                                            }}
                                        />
                                        <TextField
                                            required
                                            fullWidth
                                            error={state.signUpErrors?.has(
                                                'lastName',
                                            )}
                                            helperText={
                                                state.signUpErrors?.first(
                                                    'lastName',
                                                )?.message
                                            }
                                            size="small"
                                            key="lastName"
                                            name="lastName"
                                            label={t('Last Name')}
                                            variant="outlined"
                                            autoComplete="family-name"
                                            value={state.lastName}
                                            onChange={(event: any) => {
                                                setState({
                                                    ...state,
                                                    lastName:
                                                        event.target.value,
                                                });
                                            }}
                                        />
                                        <TextField
                                            required
                                            fullWidth
                                            error={state.signUpErrors?.has(
                                                'email',
                                            )}
                                            helperText={
                                                state.signUpErrors?.first(
                                                    'email',
                                                )?.message
                                            }
                                            size="small"
                                            name="email"
                                            key="email"
                                            label={t('Email')}
                                            autoComplete="email"
                                            variant="outlined"
                                            value={state.email}
                                            onChange={(event: any) => {
                                                setState({
                                                    ...state,
                                                    email: event.target.value,
                                                });
                                            }}
                                        />
                                        <TextField
                                            required
                                            fullWidth
                                            error={state.signUpErrors?.has(
                                                'password',
                                            )}
                                            helperText={
                                                state.signUpErrors?.first(
                                                    'password',
                                                )?.message
                                            }
                                            size="small"
                                            name="password"
                                            key="password"
                                            label={t('Password')}
                                            type="password"
                                            autoComplete="password"
                                            variant="outlined"
                                            value={state.password}
                                            onChange={(event: any) => {
                                                setState({
                                                    ...state,
                                                    password:
                                                        event.target.value,
                                                });
                                            }}
                                        />

                                        <GoogleInvisibleCaptcha />

                                        <Button
                                            disabled={state.loading}
                                            fullWidth
                                            onClick={(event: MouseEvent) => {
                                                event.preventDefault();
                                                handleSignUp();
                                            }}
                                            color="secondary"
                                            variant="contained"
                                        >
                                            {t('Sign Up')}
                                        </Button>
                                        {state.loading && (
                                            <CircularProgress
                                                size={24}
                                                className={
                                                    classes.loadingSpinner
                                                }
                                            />
                                        )}
                                    </Box>
                                </form>
                            </Box>
                        ) : (
                            <Box width="100%">
                                <Typography variant="h5">
                                    {t('Sign In')}
                                </Typography>
                                <Typography variant="body2">
                                    {t('or')}{' '}
                                    <Anchor
                                        href="#"
                                        onClick={(event: MouseEvent) => {
                                            event.preventDefault();
                                            setState({
                                                ...state,
                                                signUpActive: true,
                                            });
                                        }}
                                    >
                                        {t('create an account')}
                                    </Anchor>
                                </Typography>
                                <form>
                                    <Box
                                        mt={2}
                                        width="100%"
                                        className={classes.form}
                                    >
                                        <TextField
                                            required
                                            fullWidth
                                            error={state.loginErrors?.has(
                                                'email',
                                            )}
                                            helperText={
                                                state.loginErrors?.first(
                                                    'email',
                                                )?.message
                                            }
                                            size="small"
                                            name="email"
                                            key="email"
                                            label={t('Email')}
                                            autoComplete="email"
                                            variant="outlined"
                                            value={state.email}
                                            onChange={(event: any) => {
                                                setState({
                                                    ...state,
                                                    email: event.target.value,
                                                });
                                            }}
                                        />
                                        <TextField
                                            required
                                            fullWidth
                                            error={state.loginErrors?.has(
                                                'password',
                                            )}
                                            helperText={
                                                state.loginErrors?.first(
                                                    'password',
                                                )?.message
                                            }
                                            size="small"
                                            name="password"
                                            key="password"
                                            label={t('Password')}
                                            type="password"
                                            autoComplete="password"
                                            variant="outlined"
                                            value={state.password}
                                            onChange={(event: any) => {
                                                setState({
                                                    ...state,
                                                    password:
                                                        event.target.value,
                                                });
                                            }}
                                        />
                                        <GoogleInvisibleCaptcha />
                                        <FormControlLabel
                                            control={
                                                <Checkbox
                                                    checked={state.rememberMe}
                                                    onChange={(
                                                        event: ChangeEvent<
                                                            HTMLInputElement
                                                        >,
                                                    ) => {
                                                        setState({
                                                            ...state,
                                                            rememberMe:
                                                                event.target
                                                                    .checked,
                                                        });
                                                    }}
                                                    name="remember-me"
                                                    color="primary"
                                                />
                                            }
                                            label={t('Remember Me')}
                                        />
                                        <Button
                                            fullWidth
                                            color="secondary"
                                            variant="contained"
                                            disabled={state.loading}
                                            onClick={(event: MouseEvent) => {
                                                event.preventDefault();
                                                handleLogin();
                                            }}
                                        >
                                            {t('Sign In')}
                                        </Button>
                                        {state.loading && (
                                            <CircularProgress
                                                size={24}
                                                className={
                                                    classes.loadingSpinner
                                                }
                                            />
                                        )}
                                    </Box>
                                </form>
                            </Box>
                        )}
                    </Box>
                </Box>
            </Grid>
        </Grid>
    );
}
