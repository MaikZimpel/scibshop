import React, {useState} from 'react';
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
import Select from '@material-ui/core/Select';
import MenuItem from '@material-ui/core/MenuItem';
import {ItemSizes} from "./itemSizes";
import Avatar from '@material-ui/core/Avatar';

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
    const [item] = useState(props.item)
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

    return (
        <Card className={classes.root} variant={"outlined"} elevation={2} key={item.id}
              style={{width: '300PX', height: '450px', 'margin-bottom': '10px'}}>
            <CardHeader
                title={item.name}
                subheader={item.description}
            />
            <CardContent className={"horizontal-box"}>
                <CardMedia>
                    <img alt={""} className={"product-pic"}
                         src={"http://localhost:8082/inventory/" + item.id + "/images/" + item.images[selectedImage]}/>
                </CardMedia>
                <div className={"vertical-box"}>
                    {
                        item.images.map((val, index) => {
                            return (
                                <Avatar alt={index} src={"http://localhost:8082/inventory/" + item.id + "/images/" + val}
                                        onMouseOver={selectImage.bind(this, index)}/>
                            )
                        })
                    }
                </div>
            </CardContent>
            <CardContent>
                {item.price} EUR
                <Paper elevation={1} variant={"outlined"} square style={{"margin-top": "10px", "padding": "3px"}}>
                    <div>Erhältlich in</div>
                        <div className={"horizontal-box"}>
                            {
                                item.colors ?
                                    item.colors.map((color, index) => {
                                        return (
                                            <div onClick={select.bind(this, index)} className={`${isSelected(index) ? "selected" : ""}`}>
                                                <Tooltip title={color.color_name}>
                                                    <img key={index} alt={color.color_name}
                                                         className={"alt-pic"}
                                                         src={"http://localhost:8082/inventory/" + item.id + "/images/" + color.image}/>
                                                </Tooltip>
                                            </div>
                                        )
                                    })
                                    :
                                    <div/>
                            }
                        </div>
                    {
                        item.sizes ? <ItemSizes sizes={item.sizes}/> : <div/>
                    }
                </Paper>
            </CardContent>
            <CardActions>
                <Tooltip title={'In den Einkaufswagen'}>
                    <IconButton aria-label={'addToCart'}>
                        <FontAwesomeIcon icon={faCartPlus}/>
                    </IconButton>
                </Tooltip>
            </CardActions>
        </Card>
    )
}