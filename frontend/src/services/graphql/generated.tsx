import gql from 'graphql-tag';
export type Maybe<T> = T | null;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
  /** The `UploadFile, // b.txt` scalar type represents a multipart file upload. */
  Upload: any;
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

export type AnalyzeRequest = {
  __typename?: 'AnalyzeRequest';
  startDate: Scalars['String'];
  endDate: Scalars['String'];
  ovChipkaartNumber: Scalars['String'];
  id: Scalars['String'];
  status: Scalars['String'];
  createdAt: Scalars['String'];
  updatedAt: Scalars['String'];
};

export type AnalzyeRequestDetails = {
  __typename?: 'AnalzyeRequestDetails';
  analyzeRequestId: Scalars['String'];
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

export type StoreAnalyzeRequestInput = {
  ovChipkaartUsername?: Maybe<Scalars['String']>;
  ovChipkaartPassword?: Maybe<Scalars['String']>;
  travelHistoryFile?: Maybe<Scalars['Upload']>;
  startDate: Scalars['String'];
  endDate: Scalars['String'];
  ovChipkaartNumber: Scalars['String'];
};

/** The `Query` type, represents all of the entry points into our object graph. */
export type Query = {
  __typename?: 'Query';
  user: User;
  analyzeRequests: Array<AnalyzeRequest>;
};


/** The `Query` type, represents all of the entry points into our object graph. */
export type QueryAnalyzeRequestsArgs = {
  skip?: Maybe<Scalars['Int']>;
  take?: Maybe<Scalars['Int']>;
  orderBy?: Maybe<Scalars['String']>;
  orderDirection?: Maybe<Scalars['String']>;
};

/** The `Mutation` type, represents all updates we can make to our data. */
export type Mutation = {
  __typename?: 'Mutation';
  createUser: AuthOutput;
  login: AuthOutput;
  cancelToken: Scalars['Boolean'];
  refreshToken: Scalars['String'];
  storeAnalyzeRequest: Scalars['Boolean'];
};


/** The `Mutation` type, represents all updates we can make to our data. */
export type MutationCreateUserArgs = {
  input: CreateUserInput;
};


/** The `Mutation` type, represents all updates we can make to our data. */
export type MutationLoginArgs = {
  input: LoginInput;
};


/** The `Mutation` type, represents all updates we can make to our data. */
export type MutationCancelTokenArgs = {
  input: CancelTokenInput;
};


/** The `Mutation` type, represents all updates we can make to our data. */
export type MutationRefreshTokenArgs = {
  input: RefreshTokenInput;
};


/** The `Mutation` type, represents all updates we can make to our data. */
export type MutationStoreAnalyzeRequestArgs = {
  input: StoreAnalyzeRequestInput;
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
