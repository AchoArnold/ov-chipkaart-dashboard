import MessageBag from '../services/message-bag/MessageBag';
import { ValidationError } from './ValidationError';

export type ValidationErrorMessageBag = MessageBag<string, ValidationError>;
