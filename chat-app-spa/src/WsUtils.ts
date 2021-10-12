export const TYPE_MESSAGE = 'msg';

export const buildMessage = (Type: string, Payload: any) => JSON.stringify({ Type, Payload });
