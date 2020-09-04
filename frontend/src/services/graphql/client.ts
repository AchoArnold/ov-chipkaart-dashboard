import { InMemoryCache, NormalizedCacheObject } from 'apollo-cache-inmemory';
import { ApolloClient } from 'apollo-client';
import { createUploadLink } from 'apollo-upload-client';
import { LOCALE_FALLBACK } from '../../constants/locales';
import { KEY_LOCALE, KEY_TOKEN } from '../../constants/localStorage';

export const link = createUploadLink({
    uri: process.env.REACT_APP_GRAPHQL_SERVER_URL,
    headers: {
        Authorization: localStorage.getItem(KEY_TOKEN) ?? '',
        AcceptLanguage: localStorage.getItem(KEY_LOCALE) ?? LOCALE_FALLBACK,
        Origin: window.location.href,
    },
});

export const client: ApolloClient<NormalizedCacheObject> = new ApolloClient({
    cache: new InMemoryCache(),
    // @ts-ignore
    link,
});
