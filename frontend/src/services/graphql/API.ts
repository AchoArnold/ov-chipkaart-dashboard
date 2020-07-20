import ApolloClient, { ApolloError } from 'apollo-client';
import gql from 'graphql-tag';
import { NormalizedCacheObject } from 'apollo-cache-inmemory';
import { AuthOutput, CreateUserInput, LoginInput } from './generated';
import { DocumentNode, GraphQLError } from 'graphql';
import { ApiResponse, LoginResponse, CreateUserResponse } from './types';
import { ValidationErrorMessageBag } from '../../domain/ValidationErrorMessageBag';
import MessageBag from '../message-bag/MessageBag';
import { ERROR_MESSAGE_INTERNAL_SERVER_ERROR } from '../../constants/errors';
import { ValidationError } from '../../domain/ValidationError';

const VALIDATION_ERROR_CODE = 'VALIDATION_ERROR';

class Api {
    client: ApolloClient<NormalizedCacheObject>;
    constructor(client: ApolloClient<NormalizedCacheObject>) {
        this.client = client;
    }

    async signUp(input: CreateUserInput): Promise<CreateUserResponse> {
        const mutation: DocumentNode = gql`
            mutation{
                createUser(
                    input:{
                        firstName: "${input.firstName}",
                        lastName: "${input.lastName}",
                        email: "${input.email}",
                        password: "${input.password}",
                        reCaptcha: "${input.reCaptcha}",
                    })
                {
                    user {
                        id,
                        firstName,
                        lastName,
                        email,
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
            .then((data: any) => {
                return new ApiResponse<AuthOutput>({
                    data: data.data.login,
                }) as CreateUserResponse;
            })
            .catch((error: ApolloError) => {
                return new ApiResponse<AuthOutput>({
                    errorTitle: this.extractMainError(error),
                    validationErrors: this.mapErrorToMessageBag(error),
                    data: undefined,
                }) as CreateUserResponse;
            });
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
                        id,
                        firstName,
                        lastName,
                        email,
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
            .then((data: any) => {
                return new ApiResponse<AuthOutput>({
                    data: data.data.login,
                }) as LoginResponse;
            })
            .catch((error: ApolloError) => {
                return new ApiResponse<AuthOutput>({
                    errorTitle: this.extractMainError(error),
                    validationErrors: this.mapErrorToMessageBag(error),
                    data: undefined,
                }) as LoginResponse;
            });
    }

    private extractMainError(error: ApolloError): string {
        let mainError: string = error.graphQLErrors
            .filter((element: GraphQLError) => {
                return !(
                    element.extensions &&
                    element.extensions.code === VALIDATION_ERROR_CODE
                );
            })
            .map((element: GraphQLError) => {
                return element.message;
            })[0];

        if (mainError === undefined || mainError === '') {
            return error.message ?? ERROR_MESSAGE_INTERNAL_SERVER_ERROR;
        }

        return mainError;
    }

    private mapErrorToMessageBag(
        error: ApolloError,
    ): ValidationErrorMessageBag {
        let messageBag: ValidationErrorMessageBag = new MessageBag<
            string,
            ValidationError
        >();

        error.graphQLErrors
            .filter((element: GraphQLError) => {
                return (
                    element.extensions &&
                    element.extensions.code === VALIDATION_ERROR_CODE
                );
            })
            .forEach((element: GraphQLError) => {
                if (element.path !== undefined) {
                    let validationError: ValidationError = {
                        key: element.path[element.path.length - 1].toString(),
                        message: element.message,
                    };
                    messageBag.add(validationError.key, validationError);
                }
            });

        return messageBag;
    }
}

interface LoginOutput {}

export default Api;
