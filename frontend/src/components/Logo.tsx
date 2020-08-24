import React from 'react';
import Box from '@material-ui/core/Box';
import { makeStyles, createStyles, Theme } from '@material-ui/core/styles';

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        container: {
            ...theme.typography.h4,
            width: 200,
            textAlign: 'center',
            color: theme.palette.primary.contrastText,
            shadow: theme.palette.warning.main,
        },
        suffix: {
            '&::after': {
                content: '""',
                width: '80%',
                marginLeft: 'auto',
                marginRight: 'auto',
                borderBottom: `3px solid ${theme.palette.warning.light}`,
                display: 'block',
                marginTop: -5,
            },
        },

        containerSmall: {
            ...theme.typography.h5,
            width: 140,
            textAlign: 'center',
            color: theme.palette.primary.contrastText,
            shadow: theme.palette.warning.main,
        },

        suffixSmall: {
            '&::after': {
                content: '""',
                width: '80%',
                marginLeft: 'auto',
                marginRight: 'auto',
                borderBottom: `2px solid ${theme.palette.warning.light}`,
                display: 'block',
                marginTop: -5,
            },
        },
    }),
);

type LogoProps = {
    backgroundColor: string;
    variant?: 'small' | 'large';
};

export default function Logo(props: LogoProps) {
    const classes = useStyles();
    return (
        <Box
            className={
                props.variant == 'small'
                    ? classes.containerSmall
                    : classes.container
            }
        >
            <span
                className={
                    props.variant == 'small'
                        ? classes.suffixSmall
                        : classes.suffix
                }
            >
                ov-chi
                <span
                    style={{
                        textShadow: `2px 0 ${props.backgroundColor}, -2px 0px ${props.backgroundColor}`,
                    }}
                >
                    p
                </span>
                kaart
            </span>
            dashboard
        </Box>
    );
}
