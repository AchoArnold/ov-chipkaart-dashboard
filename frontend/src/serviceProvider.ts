import LandingPageApi from './services/graphql/LandingPageApi';
import { client } from './services/graphql/client';
import DashboardApi from './services/graphql/DashboardApi';
import AuthApiClass from './services/graphql/AuthApi';

export const LandingPageAPI = new LandingPageApi(client);
export const DashboardAPI = new DashboardApi(client);
export const AuthApi = new AuthApiClass(client);
