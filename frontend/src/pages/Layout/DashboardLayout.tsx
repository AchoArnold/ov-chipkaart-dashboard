import React from 'react';
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
import ROUTE_NAMES, { ROUTE_LANDING_PAGE } from '../../constants/routes';
import { Link } from 'react-router-dom';
import AssessmentIcon from '@material-ui/icons/Assessment';
import { useTranslation } from 'react-i18next';
import { AuthApi } from '../../serviceProvider';
import { sendToastNotification } from '../../services/notifications';
import { useSnackbar } from 'notistack';
import { VARIANT_ERROR, VARIANT_SUCCESS } from '../../constants/errors';
import { CancelTokenResponse } from '../../services/graphql/types';
import { KEY_TOKEN } from '../../constants/localStorage';

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
    }),
);

export default function DashboardLayout(props: any) {
    const classes = useStyles();
    const theme = useTheme();
    const { t } = useTranslation();
    const { enqueueSnackbar } = useSnackbar();

    const handleLogout = async () => {
        await AuthApi.logout()
            .then(() => {
                sendToastNotification(
                    enqueueSnackbar,
                    'You have successfully logged out',
                    VARIANT_SUCCESS,
                );
                localStorage.removeItem(KEY_TOKEN);
                window.location.href = ROUTE_LANDING_PAGE;
            })
            .catch((response: CancelTokenResponse) => {
                sendToastNotification(
                    enqueueSnackbar,
                    response.getErrorTitle(),
                    VARIANT_ERROR,
                );
            });
    };

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
                    <Tooltip title={t('Logout') as string}>
                        <IconButton
                            color="inherit"
                            aria-label={'Logout'}
                            onClick={handleLogout}
                        >
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
                {props.children}
            </main>
        </div>
    );
}
