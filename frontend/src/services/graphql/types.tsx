import { DocumentNode } from 'graphql';
import gql from 'graphql-tag';
import { ValidationErrorMessageBag } from '../../domain/ValidationErrorMessageBag';
import { AuthOutput } from './generated';

export const loginMutation: DocumentNode = gql`
    mutation {
        login(
            input: {
                email: "?email"
                password: "?password"
                reCaptcha: "recaptcha"
                rememberMe: true
            }
        ) {
            user {
                createdAt
                updatedAt
            }
            token {
                value
            }
        }
    }
`;

interface ResponseOptions<T> {
    errorTitle?: string;
    validationErrors?: ValidationErrorMessageBag;
    data?: T;
}

export class ApiResponse<T> {
    errorTitle?: string;
    validationErrors?: ValidationErrorMessageBag;
    data?: T;

    constructor(options: ResponseOptions<T>) {
        this.errorTitle = options.errorTitle;
        this.validationErrors = options.validationErrors;
        this.data = options.data;
    }

    getValidationErrors(): ValidationErrorMessageBag | undefined {
        return this.validationErrors;
    }

    hasValidationErrors(): boolean {
        return this.validationErrors?.isEmpty() === false;
    }

    hasServerError(): boolean {
        return this.errorTitle !== undefined && !this.hasValidationErrors();
    }

    isValid(): boolean {
        return !this.hasValidationErrors() && !this.hasServerError();
    }

    getErrorTitle(): string | undefined {
        return this.errorTitle;
    }
}

export type LoginResponse = ApiResponse<AuthOutput>;
export type CreateUserResponse = ApiResponse<AuthOutput>;
