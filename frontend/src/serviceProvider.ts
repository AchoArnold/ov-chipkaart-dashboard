import Api from './services/graphql/API';
import { client } from './services/graphql/client';

export const ApiService = new Api(client);
