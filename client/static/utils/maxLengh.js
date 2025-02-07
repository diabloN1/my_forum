const inputs = [...document.getElementsByTagName('textarea'), ...document.getElementsByTagName('input')]

console.log(inputs)
inputs.forEach((element) => element.addEventListener("input", function () {

    // Enforce max length manually
    if (element.value.length > maxLength) {
        element.value = element.value.substring(0, maxLength); // Trim excess
        value.length = maxLength;
    }
    console.log(element)
})
)