const load_products = () => new Promise((res) =>{
    fetch("http://localhost:8080/get_products", {
        method: "GET"
    }).then(data => data.json().then(json_data => res(json_data)))
})

const main = async () => {
    let elements = await load_products()
    let main_div = document.getElementById("productList")
    for (let el of elements) {
        let prod = document.createElement("tr")
        let id = document.createElement("td")
        let name = document.createElement("td")
        let price = document.createElement("td")
        let type = document.createElement("td")
        id.innerText = `${el.id}`
        name.innerText = `${el.name}`
        price.innerText = `${el.price}`
        type.innerText =`${el.type}`
        prod.appendChild(id)
        prod.appendChild(name)
        prod.appendChild(price)
        prod.appendChild(type)
        main_div.appendChild(prod)
    }
}

const loop = () => {
    new Promise((resolve) => {
        main()
        setTimeout(() => {
            resolve()
        }, 1000)
    })
}
loop()