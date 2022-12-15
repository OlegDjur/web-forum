const likeBtn = document.querySelector(".like_btn");
let likeIcon = document.querySelector("#icon");
let count = document.querySelector("#count");

let clicked = false;

likeBtn.addEventListener("click", () => {
    if (!clicked) {
        clicked = true;
        likeIcon.innerHTML = `<i class="fa-regular fa-thumbs-up"></i>`
        count.textContent++;
    }
    else {
        clicked = false;
        likeIcon.innerHTML = `<i class="fa-solid fa-thumbs-up"></i>`
        count.textContent--;
    }
});

const LikeBtn = document.querySelectorAll("#like_btn");

