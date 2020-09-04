import { VariantType } from 'notistack';

export const sendToastNotification = (
    enqueueSnackbar: Function,
    message: string | undefined,
    variant: VariantType,
) => {
    if (message !== undefined && message !== '') {
        enqueueSnackbar(message, { variant });
    }
};
