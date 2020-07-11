import axios from 'axios';

export async function loadInventory() {
    async function fetch() {
        const result = await axios('http://192.168.178.35:8082/inventory/?stockable=true&available=true&qty=gt0');
        return result.data
    }

    return await fetch().then(data => {
        return data;
    }).catch(ex => console.log(ex))
}