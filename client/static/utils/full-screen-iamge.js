const expandBtn = document.querySelector(".expand-image-btn");
const fullscreenImage = document.getElementById("fullscreen-image");
const fullscreenImageContent = fullscreenImage.querySelector(
  ".fullscreen-image-content"
);
const closeBtn = document.querySelector(".close-fullscreen-btn");
const postImage = document.querySelector(".post-image");

expandBtn.addEventListener("click", () => {
  fullscreenImageContent.src = postImage.src;
  fullscreenImage.classList.add("active");
});

closeBtn.addEventListener("click", () => {
  fullscreenImage.classList.remove("active");
});

fullscreenImage.addEventListener("click", (e) => {
  if (e.target === fullscreenImage) {
    fullscreenImage.classList.remove("active");
  }
});
