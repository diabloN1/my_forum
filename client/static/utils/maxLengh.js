const inputs = [...document.getElementsByTagName('textarea'), ...document.getElementsByTagName('input')]

console.log(inputs)
inputs.forEach((element) => element.addEventListener("input", function () {

    // Enforce max length manually
    if (element.maxLength && element.value.length > element.maxLength) {
        element.value = element.value.substring(0, element.maxLength); // Trim excess
        element.length = element.maxLength;
    }
    console.log(element)
})
)