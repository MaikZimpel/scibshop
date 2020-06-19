import React, { useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardMedia from '@material-ui/core/CardMedia';
import CardContent from '@material-ui/core/CardContent';
import Avatar from '@material-ui/core/Avatar';
import IconButton from '@material-ui/core/IconButton';
import Checkbox from '@material-ui/core/Checkbox';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import TextField from '@material-ui/core/TextField';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import axios from 'axios'
import Fab from '@material-ui/core/Fab';
import AddIcon from '@material-ui/icons/Add';
import DeleteIcon from '@material-ui/icons/Delete';

const useStyles = makeStyles((theme) => ({
    root: {
        '& .MuiTextField-root': {
            margin: theme.spacing(1),
            width: '25ch',
        },
    }
}));

export default function InventoryCard(props) {

    const [item, setItem] = useState(props.item)


    /* Backend Data Structure
    Id          string             `json:"id"`
	Upc         string             `json:"upc"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Categories  []string           `json:"categories"`
	Brand       string             `json:"brand"`
	Size        string             `json:"size"`
	Color       string             `json:"color"`
	Price       int                `json:"price"`
	Images      []string           `json:"images"`
	Supplier    string             `json:"suppliers"`
	Sku         string             `json:"sku"`
	Cnt         int                `json:"cnt"`
	Stockable   bool               `json:"stockable"`
	Available   bool               `json:"available"`
     */

    const classes = useStyles();

    function addImage(axiosResponse) {
        const imagePath = axiosResponse.data
        console.log(imagePath)
        let imageArray = item.images
        if (imageArray == null){
            imageArray = []
        }
        imageArray.push(imagePath)
        setItem({...item, images: imageArray})
    }

    function removeImage(ix) {
        let imageArray = item.images
        imageArray.splice(ix,1)
        setItem({...item, images: imageArray})
    }

    const onFileChange = event => {
        // get signed url from backend
        const formData = new FormData()
        formData.append("content_type", "image/jpeg")
        const file = event.target.files[0]
        formData.append('originalFile', file)
        axios.post("http://localhost:8082/inventory/"+item.id+"/images", formData)
            .then(addImage)
            .catch((ex) => console.log(ex))
    }

    function deletePicture(imageName, imgIndex) {
        axios.delete("http://localhost:8082/inventory/"+item.id+"/images/"+imageName)
            .then(removeImage(imgIndex))
            .catch((ex) => console.log(ex))
    }

    return (
        <Card className={classes.root} variant={"outlined"} elevation={2}>
            <CardHeader
                avatar={
                    item.images ?
                        <img alt={""} className={"avtr-pic"} src={"http://localhost:8082/inventory/" + item.id + "/images/" + item.images[0]}/>
                        :
                        <Avatar aria-label="inventory-item" className={"_avatar"}/>
                }
                action={
                    <IconButton aria-label="actions">
                        <MoreVertIcon />
                    </IconButton>
                }
                title={item.name}
                subheader={item.description}
            />
            <div className={"horizontal-box"}>
                {
                    item.images ?
                    item.images.map((val, index) => {
                        return (
                            <div key={index} className={"container"}>
                                <img className={"image"} alt={""} src={"http://localhost:8082/inventory/" + item.id + "/images/" + val}/>
                                <div className={"middle"}>
                                    <Fab
                                        color="primary"
                                        size="small"
                                        component="div"
                                        aria-label="add"
                                        variant="extended"
                                    >
                                        <DeleteIcon onClick={event => deletePicture(val, index)}/>
                                    </Fab>
                                </div>
                            </div>


                        );
                    }): <div/>
                }
            </div>
            <CardContent>
                <div className={"horizontal-box"}>
                    <form>
                        <label htmlFor={"upload-picture"}>
                            <input style={{display: 'none'}} id={"upload-picture"} name={"upload-picture"}
                                   type={"file"} onChange={onFileChange} aria-label={"Add picture"}/>
                            <Fab
                                color="primary"
                                size="small"
                                component="div"
                                aria-label="add"
                                variant="extended"
                            >
                                <AddIcon />
                            </Fab>
                        </label>

                    </form>
                </div>
            </CardContent>

            <CardContent>
                <form className={classes.root}>
                    <div>
                        <TextField required id={"name"} value={item.name} label={"Name"} variant={"outlined"}
                                   onChange={e => setItem({...item, name: e.target.value})}/>
                        <TextField id={"upc"} value={item.upc} label={"Upc"} variant={"outlined"}
                                   onChange={e => setItem({...item, upc: e.target.value})}/>
                        <TextField required id={"sku"} value={item.sku} label={"Sku"} variant={"outlined"}
                                   onChange={e => setItem({...item, sku: e.target.value})}/>
                    </div>
                    <div>
                        <TextField id={"brand"} value={item.brand} label={"Brand"} variant={"outlined"}
                                   onChange={e => setItem({...item, brand: e.target.value})}/>
                        <TextField id={"size"} value={item.size} label={"Size"} variant={"outlined"}
                                   onChange={e => setItem({...item, size: e.target.value})}/>
                        <TextField required id={"color"} value={item.color} label={"Color"} variant={"outlined"}
                                   onChange={e => setItem({...item, color: e.target.value})}/>
                    </div>
                    <div>
                        <TextField id={"price"} type={"number"} value={item.price} label={"Price"} variant={"outlined"}
                                   onChange={e => setItem({...item, price: e.target.value})}/>
                        <TextField id={"supplier"} value={item.suppliers} label={"Supplier"} variant={"outlined"}
                                   onChange={e => setItem({...item, suppliers: e.target.value})}/>
                        <TextField type={"number"} required id={"qty"} value={item.cnt} label={"Quantity"} variant={"outlined"}
                                   onChange={e => setItem({...item, cnt: e.target.value})}/>
                    </div>
                    <div>
                        <FormControlLabel control={
                            <Checkbox
                                color="primary"
                                checked={item.stockable}
                                onChange={(e, value) => setItem({...item, stockable: value})}
                            />
                        } label={"Can be stocked"}/>
                    </div>
                    <div>
                        <FormControlLabel control={
                            <Checkbox
                                color="primary"
                                checked={item.available}
                                onChange={(e, value) => setItem({...item, available: value})}
                            />
                        } label={"Can be ordered"}/>
                    </div>
                </form>
            </CardContent>
        </Card>
    );
}
