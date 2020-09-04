import ApolloClient, { ApolloError } from 'apollo-client';
import { NormalizedCacheObject } from 'apollo-cache-inmemory';
import { GraphQLError } from 'graphql';
import { ERROR_MESSAGE_INTERNAL_SERVER_ERROR } from '../../constants/errors';
import { ValidationErrorMessageBag } from '../../domain/ValidationErrorMessageBag';
import MessageBag from '../message-bag/MessageBag';
import { ValidationError } from '../../domain/ValidationError';
import _ from 'lodash';

const VALIDATION_ERROR_CODE = 'VALIDATION_ERROR';

export default class BaseApi {
    protected client: ApolloClient<NormalizedCacheObject>;
    constructor(client: ApolloClient<NormalizedCacheObject>) {
        this.client = client;
    }

    protected extractMainError(error: ApolloError): string {
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

    protected mapErrorToMessageBag(
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
                    let key = element.path[element.path.length - 1].toString();
                    let validationError: ValidationError = {
                        key: key,
                        message: element.message.replace(key, _.startCase(key)),
                    };
                    messageBag.add(validationError.key, validationError);
                }
            });

        return messageBag;
    }
}
