import React, { useState, MouseEvent, ChangeEvent } from 'react';
import Grid from '@material-ui/core/Grid';
import { makeStyles, createStyles, Theme } from '@material-ui/core/styles';
import Box from '@material-ui/core/Box';
import ROUTE_NAMES from '../../constants/routes';
import { Link } from 'react-router-dom';
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
import { ApiService } from '../../serviceProvider';
import { LoginResponse } from '../../services/graphql/types';

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

export default function LandingPage() {
    const classes = useStyles();
    const { t } = useTranslation();
    const theme = useTheme();
    const [state, setState] = useState({
        signUpActive: true,
        loading: false,
        email: '',
        password: '',
        firstName: '',
        lastName: '',
        rememberMe: true,
    });

    const handleLogin = function () {
        ApiService.login({
            email: state.email,
            password: state.password,
            rememberMe: state.rememberMe,
            reCaptcha: 'how are you',
        })
            .then((response: LoginResponse) => {
                console.log('then', response);
            })
            .catch((response: LoginResponse) => {
                console.log('catch', response);
            });
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
                                                setState({
                                                    ...state,
                                                    loading: true,
                                                });
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
                                                setState({
                                                    ...state,
                                                    loading: true,
                                                });

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
