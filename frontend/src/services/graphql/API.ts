import ApolloClient from 'apollo-client';
import gql from 'graphql-tag';
import { NormalizedCacheObject } from 'apollo-cache-inmemory';
import { LoginInput } from './generated';
import { DocumentNode } from 'graphql';
import { LoginResponse, ValidationError } from './types';
import { TFunction } from 'i18next';

const VALIDATION_ERROR_CODE = 'VALIDATION_ERROR';

class Api {
    client: ApolloClient<NormalizedCacheObject>;
    translate: TFunction;
    constructor(
        client: ApolloClient<NormalizedCacheObject>,
        translationFunction: TFunction,
    ) {
        this.client = client;
        this.translate = translationFunction;
    }

    async login(input: LoginInput): Promise<LoginResponse> {
        const mutation: DocumentNode = gql`
            mutation{
                login(
                    input:{
                        email: "${input.email}",
                        password: "${input.password}",
                        reCaptcha: "${input.reCaptcha}",
                        rememberMe: ${input.rememberMe.toString()}
                })
                {
                    user {
                        createdAt,
                        updatedAt,
                    }
                    token {
                        value
                    }
                }
            }
        `;

        return await this.client
            .mutate({
                mutation,
            })
            .then((data: LoginOutput) => {
                return {
                    hasValidationError: false,
                    hasServerError: false,
                } as LoginResponse;
            })
            .catch((error: any) => {
                let validationErrors = error.graphQLErrors
                    .filter((element: any) => {
                        return (
                            element.extensions &&
                            element.extensions.code === VALIDATION_ERROR_CODE
                        );
                    })
                    .map((element: any) => {
                        return {
                            key: element.path[element.path.length - 1],
                            message: element.message,
                        } as ValidationError;
                    });

                let mainError: string = error.graphQLErrors
                    .filter((element: any) => {
                        return !(
                            element.extensions &&
                            element.extensions.code === VALIDATION_ERROR_CODE
                        );
                    })
                    .map((element: any) => {
                        return element.message;
                    })[0];

                if (mainError === undefined) {
                    mainError = this.translate('internal server error');
                }

                return {
                    hasValidationError: validationErrors.length > 0,
                    hasServerError: validationErrors.length === 0,
                    errorTitle: mainError,
                    errors: validationErrors,
                    data: undefined,
                } as LoginResponse;
            });
    }
}

interface LoginOutput {}

export default Api;
