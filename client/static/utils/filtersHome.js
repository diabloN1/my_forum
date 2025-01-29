const filterContainer = document.getElementById('filterContainer')
const searchInput = document.getElementById('searchInput')
const categoryFilterButtons = document.getElementById('categoryFilterButtons')

const filterButton = document.getElementById('filterButton')
filterButton.addEventListener('click', () => {
    console.log("button clicked")
    triggerVisibility(filterContainer)
})

function triggerVisibility(element) {
    const computedStyle = window.getComputedStyle(element)

    if (computedStyle.display === 'none') {
        element.style.display = "block"
    } else {
        element.style.display = "none"
    }
}

// Read and parse json (takes a string and returns the parsed object)
const postsData = JSON.parse(document.getElementById('postsData').textContent)



///////////////////////////////////// Search Bar Working ///////////////////////////////////////
// console.log(postsData)

// DFS to extract all search suggestions
const searchExemples = new Set(); // Set is an array that only holds unique items
const categories = new Set(); // Set is an array that only holds unique items
const stack = [{ value: postsData, parent: "" }]; // Initialize the stack

while (stack.length > 0) {
    const current = stack.pop(); // Get and remove the last element from the stack
    const { value, parent } = current; // Destructure to get value and parent

    if (parent == "category") {
        searchExemples.add(value + " - " + parent);
        categories.add(value);
    } else if (parent == "user_name" || parent == "title") { //&& (parent != "image")
        searchExemples.add(value + " - " + parent);
    } else if (value instanceof Array) { // We didn't use typeof because it define the array as an object
        // If it's an array, push all its items onto the stack with the current parent name
        value.forEach((item) => {
                stack.push({ value: item, parent: parent }); // Keep the parent name the same for array items      
        });
    } else if (value instanceof Object) {
        // If it's an object, push all its values onto the stack with their keys as parent names
        Object.entries(value).forEach(([key, val]) => {
            stack.push({ value: val, parent: key }); // Use the key as the parent name
        });
    }
}

// searchExemples.forEach((exemple) => {
//     document.getElementById('searchExemples').innerHTML += ("<option value='"+exemple.split(" - ")[0]+"'>" + exemple + "</option></br>")
// })
// console.log(searchExemples)

// // Search Input EventListener
// searchInput.addEventListener('input', (event) => {
//     const searchValue = event.target.value
//     const postsDivs = document.getElementById('postsContainer').children
//     console.log(event.target.value)

//     postsData.forEach((post, index)=>{
//         console.log(post.user_name, post.category, post.title)
//         if (
//             post.user_name.toLowerCase().includes(searchValue) ||
//             post.category.toLowerCase().includes(searchValue) ||
//             post.title.toLowerCase().includes(searchValue)
//             ) {
//             console.log(2)
//             const item = postsDivs[index]
//             item.style.display = "block"
//         } else {
//             console.log(1)
//             const item = postsDivs[index]
//             item.style.display = "none"
//         }
//     })
//     console.log(postsDivs)
// })

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Categories Buttons
console.log(categories)
categories.forEach((item)=>{
    if (item != "") {
        console.log(3141)
        categoryFilterButtons.innerHTML += "<button class='categoryButtonsClass'>"+item+"</button>"
    }
})

const categoryFilterButtonsChildren = document.getElementById('categoryFilterButtons').childNodes
categoryFilterButtonsChildren.forEach((button) => {
    // button.addEventListener('click', () => {
    //     console.log('clkmw', button.innerHTML)
    //     searchInput.value = button.innerHTML
    // });
});