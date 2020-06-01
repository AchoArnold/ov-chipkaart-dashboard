import Api from './services/graphql/API';
import { client } from './services/graphql/client';
import i18n from './i18n';

export const ApiService = new Api(client, i18n.t);
