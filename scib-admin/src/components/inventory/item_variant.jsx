import React, {useContext, useState} from 'react';
import Paper from '@material-ui/core/Paper';
import Fab from "@material-ui/core/Fab";
import DeleteIcon from "@material-ui/icons/Delete";
import TextField from "@material-ui/core/TextField";
import IconButton from '@material-ui/core/IconButton';
import AddIcon from "@material-ui/icons/Add";
import * as api from './inventory_api'
import {InventoryContext} from "./inventoryContext";

export const ItemVariant = (props) => {

    /*

	Sku       string `json:"sku" bson:"sku, omitempty"`
	Color     string `json:"color" bson:"color, omitempty"`
	Image     string `json:"image" bson:"image, omitempty"`
	Size string `json:"size" bson:"size, omitempty"`
	Cnt       int    `json:"cnt" bson:"cnt, omitempty"`
	Stockable bool   `json:"stockable" bson:"stockable, omitempty"`
	Available bool   `json:"available" bson:"available, omitempty"`


     */

    const {items, actions} = useContext(InventoryContext);
    const itemId = props.itemId;
    const index = props.index;
    const [variant, setVariant] = useState(items.find(i => i.id === props.itemId).variants[props.index]);

    const remove = () => {
        if (variant.image) {
            api.deleteImageFile(itemId, variant.image).then(() => actions.removeVariant(itemId, index));
        } else {
            actions.removeVariant(itemId, index);
        }
    }

    function deletePicture(imageName) {
        api.deleteImageFile(itemId, imageName).then(() => updateImage(null)).catch(console.log);
    }

    function upload(event, ind) {
        const formData = new FormData();
        formData.append("content_type", "image/jpeg");
        const file = event.target.files[0];
        formData.append('originalFile', file);
        formData.append('isColorImage', 'true');
        formData.append('variantIndex', ind);
        formData.append('sku', variant.sku);
        formData.append('color', variant.color);
        formData.append('size', variant.size);
        formData.append('cnt', variant.cnt);
        formData.append('stockable', JSON.stringify(variant.stockable));
        formData.append('available', JSON.stringify(variant.available));
        api.upload(itemId, formData).then(r => updateImage(r.data)).catch(console.log);
    }

    function handleUpdate (e)  {
        const name = e.target.name;
        const value = function() {
            if (name === 'cnt') {
                return Number(e.target.value);
            } else {
                return e.target.value;
            }
        }() ;
        setVariant(variant => ({...variant, [name]: value}));
        actions.updateVariant(itemId, index, name, value);
    }

    const updateImage = (imageId) => {
        setVariant(variant => ({...variant, image: imageId}));
        actions.updateVariant(itemId, index, "image", imageId);
    }

    return(
        <Paper>
            <div className={"horizontal-box-middle"}>

                    <div>
                        {
                            variant.image &&
                            <div className={"container"}
                                 style={{height: "25px", width: "25px"}}>
                                <img alt={variant.image}
                                     className={"colorImage"}
                                     src={"http://localhost:8082/inventory/" + itemId + "/images/" + variant.image}/>
                                <div className={"middle"}>
                                    <Fab
                                        color="primary"
                                        size="small"
                                        component="div"
                                        aria-label="remove"
                                        variant="extended"
                                    >
                                        <DeleteIcon onClick={() => deletePicture(variant.image)}/>
                                    </Fab>
                                </div>
                            </div>

                        }
                        {
                            !variant.image &&
                            <div>
                                    <label htmlFor={"upload-color-picture"}>
                                        <input id={"upload-color-picture"} style={{display: 'none'}} type={"file"}
                                               onChange={event => upload(event, index)} aria-label={"Add picture"}/>
                                        <IconButton size="small" component="div" aria-label="add" variant="extended">
                                            <AddIcon/>
                                        </IconButton>
                                    </label>

                            </div>
                        }
                </div>
                <TextField value={variant.color} name={"color"} onChange={handleUpdate} label={"Color"}/>
                <TextField value={variant.size} name={"size"} onChange={handleUpdate} label={"Size"}/>
                <TextField type={"Number"} value={variant.cnt} name={"cnt"} onChange={handleUpdate} label={"Qnt"}/>
                <IconButton onClick={remove}>
                    <DeleteIcon/>
                </IconButton>
            </div>
        </Paper>
    )

}