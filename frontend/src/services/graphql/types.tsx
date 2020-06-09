import { DocumentNode } from 'graphql';
import gql from 'graphql-tag';

const loginMutation: DocumentNode = gql`
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

interface defaultContent {
    hasValidationError: boolean;
    hasServerError: boolean;
    errorTitle?: string;
    errors?: Array<ValidationError>;
    data?: object;
}

export interface ValidationError {
    key: string;
    message: string;
}

export interface ValidationErrorContainer {}

export interface LoginResponse extends defaultContent {}
