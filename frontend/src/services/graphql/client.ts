import { createHttpLink } from 'apollo-link-http';
import { InMemoryCache, NormalizedCacheObject } from 'apollo-cache-inmemory';
import { ApolloClient } from 'apollo-client';
import { LOCALE_FALLBACK } from '../../constants/locales';
import { KEY_LOCALE, KEY_TOKEN } from '../../constants/localStorage';

export const link = createHttpLink({
    uri: process.env.REACT_APP_GRAPHQL_SERVER_URL,
    headers: {
        Authorization: localStorage.getItem(KEY_TOKEN) ?? '',
        AcceptLanguage: localStorage.getItem(KEY_LOCALE) ?? LOCALE_FALLBACK,
        Origin: window.location.href,
    },
});

export const client: ApolloClient<NormalizedCacheObject> = new ApolloClient({
    cache: new InMemoryCache(),
    link,
});
