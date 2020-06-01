import gql from 'graphql-tag';
export type Maybe<T> = T | null;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
};

export type User = {
  __typename?: 'User';
  id: Scalars['ID'];
  firstName: Scalars['String'];
  lastName: Scalars['String'];
  email: Scalars['String'];
  createdAt: Scalars['String'];
  updatedAt: Scalars['String'];
};

export type Token = {
  __typename?: 'Token';
  value: Scalars['String'];
};

export type Query = {
  __typename?: 'Query';
  user: User;
};

export type CreateUserInput = {
  firstName: Scalars['String'];
  lastName: Scalars['String'];
  email: Scalars['String'];
  password: Scalars['String'];
  reCaptcha: Scalars['String'];
};

export type AuthOutput = {
  __typename?: 'AuthOutput';
  user: User;
  token: Token;
};

export type CancelTokenInput = {
  token: Scalars['String'];
};

export type RefreshTokenInput = {
  token: Scalars['String'];
};

export type LoginInput = {
  email: Scalars['String'];
  password: Scalars['String'];
  rememberMe: Scalars['Boolean'];
  reCaptcha: Scalars['String'];
};

export type Mutation = {
  __typename?: 'Mutation';
  createUser: AuthOutput;
  login: AuthOutput;
  cancelToken: Scalars['Boolean'];
  refreshToken: Scalars['String'];
};


export type MutationCreateUserArgs = {
  input: CreateUserInput;
};


export type MutationLoginArgs = {
  input: LoginInput;
};


export type MutationCancelTokenArgs = {
  input: CancelTokenInput;
};


export type MutationRefreshTokenArgs = {
  input: RefreshTokenInput;
};

export type Unnamed_1_MutationVariables = {};


export type Unnamed_1_Mutation = (
  { __typename?: 'Mutation' }
  & { login: (
    { __typename?: 'AuthOutput' }
    & { user: (
      { __typename?: 'User' }
      & Pick<User, 'createdAt' | 'updatedAt'>
    ), token: (
      { __typename?: 'Token' }
      & Pick<Token, 'value'>
    ) }
  ) }
);

