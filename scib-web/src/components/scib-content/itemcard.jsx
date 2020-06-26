import React, {useState} from 'react';
import {makeStyles} from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardContent from '@material-ui/core/CardContent';
import CardMedia from "@material-ui/core/CardMedia";
import CardFooter from '@material-ui/core/CardActions';
import ShoppingCartIcon from "@material-ui/icons/ShoppingCart";
import CardActions from '@material-ui/core/CardActions'
import Tooltip from '@material-ui/core/Tooltip';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {faCartPlus} from '@fortawesome/free-solid-svg-icons'
import IconButton from '@material-ui/core/IconButton';

const useStyles = makeStyles((theme) => ({
    root: {
        '& .MuiTextField-root': {
            margin: theme.spacing(1),
            width: '25ch',
        },
        '& .MuiCardHeader-root': {
            'display': 'flex',
            'flex-direction': 'row',
            'align-items': 'flex-start',
        },
        '& .MuiCardHeader-subheader': {}
    },
}));

export const ItemCard = (props) => {

    const classes = useStyles();
    const [item, setItem] = useState(props.item)

    return(
        <Card className={classes.root} variant={"outlined"} elevation={2} key={item.id} style={{width:'300PX', height:'450px', 'margin-bottom':'10px'}}>
            <CardHeader
                title={item.name}
                subheader={item.description}
            />
            <CardContent>
                <CardMedia>
                    <img alt={""} className={"product-pic"}
                         src={"http://localhost:8082/inventory/" + item.id + "/images/" + item.images[0]}/>
                </CardMedia>
            </CardContent>
            <CardContent>
                {item.price}
            </CardContent>
            <CardActions>
                <Tooltip title={'Add to cart'}>
                    <IconButton aria-label={'addToCart'}>
                        <FontAwesomeIcon icon={faCartPlus}/>
                    </IconButton>
                </Tooltip>
            </CardActions>
        </Card>
    )
}