const deepClone = (obj) => {
    if (obj ===null){
        return null
    }
    const clone = {};
    Object.keys(obj).forEach(item=>{
        if (obj[item] instanceof Object){
            clone[item] = deepClone(obj[item]);
        }else{
            clone[item] = obj[item];
        }
    });
    return clone
}

const a = { foo: "bar", obj: { a: 1, b: 2 } };
const b = deepClone(a);
console.log(b);
