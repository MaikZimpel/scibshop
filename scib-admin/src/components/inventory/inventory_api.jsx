import axios from "axios";

export async function loadInventory() {
    async function fetch() {
        const result = await axios('http://localhost:8082/inventory/');
        return result.data
    }

    return await fetch().then(data => {
        return data;
    }).catch(ex => console.log(ex))
}

export async function removeItem(item) {
    async function del() {
        const result = await axios.delete("http://localhost:8082/inventory/" + item.id);
        return result.status;
    }
    return await del().then(response => {
        return response === 204;
    }).catch(ex => console.log(ex));
}

export async function saveItem(item) {
    async function save() {
        const result = await axios.put("http://localhost:8082/inventory/" + item.id, item);
        return result.status;
    }
    return await save().then(response => {
        return response;
    }).catch(ex => console.log(ex));
}

export async function upload(itemId, formData) {
    async function up() {
        return await axios.post("http://localhost:8082/inventory/" + itemId + "/images", formData)
    }
    return await up().then((response) => {
        return response;
    }).catch((ex) => console.log(ex));
}

export async function deleteImageFile(itemId, imageId) {
    async function del() {
        return await axios.delete("http://localhost:8082/inventory/" + itemId + "/images/" + imageId)
    }
    return await del().then((response) => {
        return response;
    }).catch((ex) => console.log(ex));
}