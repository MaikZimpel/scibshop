import {withStyles} from "@material-ui/core/styles";
import {yellow, blue} from "@material-ui/core/colors";
import Button from "@material-ui/core/Button";

export const CheckoutButton = withStyles((theme) => ({
    root: {
        color: theme.palette.getContrastText(yellow[500]),
        backgroundColor: yellow[500],
        '&:hover': {
            backgroundColor: yellow[700],
        },
    },
}))(Button);

export const CloseDlgButton = withStyles((theme) => ({
    root: {
        color: theme.palette.getContrastText(blue[500]),
        backgroundColor: blue[500],
        '&:hover': {
            backgroundColor: blue[700],
        },
    },
}))(Button);