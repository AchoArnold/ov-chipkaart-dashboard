import BaseApi from './BaseApi';
import { DocumentNode } from 'graphql';
import gql from 'graphql-tag';
import { ApiResponse } from './types';
import { AnalyzeRequest, StoreAnalyzeRequestInput } from './generated';
import { ApolloError } from 'apollo-client';

export default class DashboardApi extends BaseApi {
    storeRequest(
        input: StoreAnalyzeRequestInput,
    ): Promise<ApiResponse<boolean>> {
        const mutation: DocumentNode = gql`
            mutation StoreAnalyzeRequest($file: Upload, $username : String, $password : String ){
                storeAnalyzeRequest(
                    input:{
                        travelHistoryFile: $file,
                        ovChipkaartUsername: $username,
                        ovChipkaartPassword: $password,
                        ovChipkaartNumber: "${input.ovChipkaartNumber}",
                        startDate: "${input.startDate}",
                        endDate: "${input.endDate}",
                    })
            }
        `;

        return new Promise((resolve, reject) => {
            this.client
                .mutate({
                    mutation,
                    variables: {
                        file: input.travelHistoryFile,
                        username: input.ovChipkaartUsername,
                        password: input.ovChipkaartPassword,
                    },
                })
                .then((data: any) => {
                    resolve(
                        new ApiResponse<boolean>({
                            data: data.data.storeAnalyzeRequest,
                        }),
                    );
                })
                .catch((error: ApolloError) => {
                    reject(
                        new ApiResponse<boolean>({
                            errorTitle: this.extractMainError(error),
                            validationErrors: this.mapErrorToMessageBag(error),
                            data: undefined,
                        }),
                    );
                });
        });
    }

    async getRecentRequests(): Promise<ApiResponse<AnalyzeRequest[]>> {
        const query: DocumentNode = gql`
            query {
                analyzeRequests(
                    take: 10
                    skip: 0
                    orderBy: "created_at"
                    orderDirection: "DESC"
                ) {
                    id
                    ovChipkaartNumber
                    startDate
                    endDate
                    status
                    createdAt
                }
            }
        `;

        return new Promise((resolve, reject) => {
            this.client
                .query({ query })
                .then((data: any) => {
                    resolve(
                        new ApiResponse<AnalyzeRequest[]>({
                            data: data.data.analyzeRequests,
                        }),
                    );
                })
                .catch((error: ApolloError) => {
                    reject(
                        new ApiResponse<AnalyzeRequest[]>({
                            errorTitle: this.extractMainError(error),
                            validationErrors: this.mapErrorToMessageBag(error),
                            data: undefined,
                        }),
                    );
                });
        });
    }
}
