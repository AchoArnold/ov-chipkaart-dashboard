import BaseApi from './BaseApi';
import { DocumentNode } from 'graphql';
import gql from 'graphql-tag';
import { CancelTokenResponse } from './types';
import { ApolloError } from 'apollo-client';

export default class AuthApi extends BaseApi {
    logout(): Promise<CancelTokenResponse> {
        const mutation: DocumentNode = gql`
            mutation {
                cancelToken
            }
        `;

        return new Promise<CancelTokenResponse>((resolve, reject) => {
            this.client
                .mutate({ mutation })
                .then(() => {
                    resolve(new CancelTokenResponse({ data: true }));
                })
                .catch((error: ApolloError) => {
                    reject(
                        new CancelTokenResponse({
                            errorTitle: this.extractMainError(error),
                            validationErrors: this.mapErrorToMessageBag(error),
                            data: false,
                        }),
                    );
                });
        });
    }
}
