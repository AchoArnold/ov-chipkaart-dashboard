import LandingPageApi from './services/graphql/LandingPageApi';
import { client } from './services/graphql/client';

export const ApiService = new LandingPageApi(client);
