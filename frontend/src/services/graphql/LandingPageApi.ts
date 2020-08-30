import { ApolloError } from 'apollo-client';
import gql from 'graphql-tag';
import { AuthOutput, CreateUserInput, LoginInput } from './generated';
import { DocumentNode } from 'graphql';
import { ApiResponse, LoginResponse, CreateUserResponse } from './types';
import BaseApi from './BaseApi';

export default class LandingPageApi extends BaseApi {
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
}
