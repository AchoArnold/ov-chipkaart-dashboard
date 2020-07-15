import React from 'react';
import ReactDOM from 'react-dom';
import * as serviceWorker from './serviceWorker';
import 'minireset.css';
import { Route, Switch, BrowserRouter as Router } from 'react-router-dom';
import LandingPage from './pages/LandingPage';
import './i18n';
import ROUTE_NAMES from './constants/routes';
import { createMuiTheme, ThemeProvider } from '@material-ui/core/styles';
import { SnackbarProvider } from 'notistack';

const theme = createMuiTheme({});

const routing = (
    <ThemeProvider theme={theme}>
        <SnackbarProvider
            anchorOrigin={{ vertical: 'top', horizontal: 'right' }}
        >
            <Router>
                <Switch>
                    <Route
                        exact
                        path={ROUTE_NAMES.LANDING_PAGE}
                        component={LandingPage}
                    />
                </Switch>
            </Router>
        </SnackbarProvider>
    </ThemeProvider>
);

ReactDOM.render(routing, document.getElementById('root'));

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
