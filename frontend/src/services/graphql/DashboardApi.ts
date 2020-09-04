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
            mutation StoreAnalyzeRequest($file: Upload ){
                storeAnalyzeRequest(
                    input:{
                        travelHistoryFile: $file,
                        ovChipkaartUsername: "${input.ovChipkaartUsername}",
                        ovChipkaartPassword: "${input.ovChipkaartPassword}",
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
                    },
                })
                .then((data: any) => {
                    resolve(
                        new ApiResponse<boolean>({ data }),
                    );
                })
                .catch((error: ApolloError) => {
                    console.log('error');
                    console.log(error);
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

    async getRecentRequests() {
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
                }
            }
        `;

        return await this.client
            .query({ query })
            .then((data: any) => {
                console.log(data);
                return new ApiResponse<AnalyzeRequest[]>({ data });
            })
            .catch((error: ApolloError) => {
                return new ApiResponse<boolean>({
                    errorTitle: this.extractMainError(error),
                    validationErrors: this.mapErrorToMessageBag(error),
                    data: undefined,
                });
            });
    }
}
