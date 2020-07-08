import React, {useState} from 'react';
import Select from "@material-ui/core/Select";
import MenuItem from "@material-ui/core/MenuItem";
import {makeStyles} from '@material-ui/core/styles';

const useStyles = makeStyles((theme) => ({
    formControl: {
        margin: theme.spacing(1),
        minWidth: 120,
    },
    selectEmpty: {
        marginTop: theme.spacing(2),
    },
}));

export const ItemSizes = (props) => {

    const classes = useStyles();
    const [sizes] = useState(props.sizes)

    return (
        <div>
            <div>Größe</div>
            <Select
                displayEmpty
                className={classes.selectEmpty}
                inputProps={{'aria-label': 'Without label'}}
            >{
                sizes.map((s, index) => {
                    return (<MenuItem value={index}>
                        <em>s.size_name</em>
                    </MenuItem>)

                })
            }
            </Select>
        </div>)
    }