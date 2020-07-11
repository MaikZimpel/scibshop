import React, {useContext, useState} from 'react';
import {makeStyles} from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardContent from '@material-ui/core/CardContent';
import CardMedia from "@material-ui/core/CardMedia";
import CardActions from '@material-ui/core/CardActions';
import Tooltip from '@material-ui/core/Tooltip';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {faCartPlus} from '@fortawesome/free-solid-svg-icons'
import IconButton from '@material-ui/core/IconButton';
import Paper from '@material-ui/core/Paper';
import Avatar from '@material-ui/core/Avatar';
import {CartContext} from "../cart-context/cartContext";
import InfoIcon from '@material-ui/icons/Info';
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';

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
    formControl: {
        margin: theme.spacing(1),
        minWidth: 120,
    },
    selectEmpty: {
        marginTop: theme.spacing(2),
    },
}));

export const ItemCard = (props) => {

    const classes = useStyles();

    const {items, actions} = useContext(CartContext);
    const [item] = useState(items.find(i => i.id === props.itemId));
    const [descriptionDlgOpen, setDescriptionDlgOpen] = useState(false);


    const [selected, setSelected] = useState(0)
    const [selectedImage, setSelectedImage] = useState(0)

    const isSelected = (index) => {
        return selected === index;
    }

    const select = (index) => {
        setSelected(index)
    }

    const selectImage = (index) => {
        setSelectedImage(index)
    }

    const addItemToCart = () => {
        actions.addToCart(item.id, item.variants[selected].sku, item.price, 1);
    }

    function toggleDescriptionDlg() {
        setDescriptionDlgOpen(!descriptionDlgOpen);
    }

    const itemSubheader = () => {
        return (
            <>
                {item.name}<IconButton onClick={toggleDescriptionDlg}><InfoIcon fontSize={"small"}/></IconButton>
                <DescriptionDialog />
            </>
        )
    }

    const DescriptionDialog = () => {
        return (
            <Dialog open={descriptionDlgOpen} onClose={toggleDescriptionDlg}>
                <DialogContent>
                    <DialogContentText>{item.description}</DialogContentText>
                </DialogContent>
            </Dialog>
        )
    }

    return (
        <Card className={classes.root} variant={"outlined"} elevation={2} key={item.id}
              style={{width: '300PX', height: '500px', 'margin-bottom': '10px'}}>
            <CardHeader
                title={item.brand}
                subheader={itemSubheader()}
            />
            <CardContent className={"horizontal-box"}>

                <CardMedia>
                    <img alt={""} className={"product-pic"}
                         src={"http://192.168.178.35:8082/inventory/" + item.id + "/images/" + item.images[selectedImage]}/>
                </CardMedia>

                <div className={"vertical-box"}>
                    {
                        item.images.map((val, index) => {
                            return (
                                <Avatar alt={index}
                                        src={"http://192.168.178.35:8082/inventory/" + item.id + "/images/" + val}
                                        onMouseOver={selectImage.bind(this, index)}/>
                            )
                        })
                    }
                </div>
            </CardContent>
            <CardContent>
                {item.price} EUR
                <Paper elevation={1} variant={"outlined"} square style={{"margin-top": "10px", "padding": "3px"}}>
                    <div>Farben</div>
                    <div className={"horizontal-box"}>
                        {
                            item.variants ?
                                item.variants.map((variant, index) => {
                                    return (
                                        <div onClick={select.bind(this, index)}
                                             className={`${isSelected(index) ? "selected" : ""}`}>
                                            <Tooltip title={variant.color}>
                                                <img key={index} alt={variant.image}
                                                     className={"alt-pic"}
                                                     src={"http://192.168.178.35:8082/inventory/" + item.id + "/images/" + variant.image}/>
                                            </Tooltip>
                                        </div>
                                    )
                                })
                                :
                                <div/>
                        }
                    </div>
                </Paper>
            </CardContent>
            <CardActions>
                <Tooltip title={'In den Einkaufswagen'}>
                    <IconButton aria-label={'addToCart'} onClick={addItemToCart}>
                        <FontAwesomeIcon icon={faCartPlus}/>
                    </IconButton>
                </Tooltip>
            </CardActions>
        </Card>
    )
}