const BASE_URL = "https://desolate-fortress-36652.herokuapp.com",
  grid = document.getElementById("grid"),
  counter = document.getElementById("counter"),
  searchInput = document.getElementById("search"),
  searchBtn = document.getElementById("search_btn");

let videos = [];

const url = (endpoint) => BASE_URL + endpoint;

const renderVideos = (video) => (grid.innerHTML += generateMarkup(video));

function debounce(cb, timer, immediate = false) {
  let timeout;
  return function () {
    const later = () => {
      timeout = null;
      if (!immediate) cb.apply(this, arguments);
    };
    const callNow = immediate && !timeout;
    clearTimeout(timeout);
    timeout = setTimeout(later, timer);
    if (callNow) cb.apply(this, arguments);
  };
}

function generateMarkup(video) {
  return `
      <div class="video shadow">
        <a href="${video.URL}" class="video_link"
        ><img
        src="${video.Thumbnail}"
        srcset="${video.Thumbnail}"
        alt="thumbnail"
        class="video_thumbnail"
        /></a>
        <p class="video_text">${video.Title}</p>
        <p class="video_text">${video.Author}</p>
        <p class="video_text">${video.UploadDate}</p>
        <p class="video_text">${video.Length}</p>
        <p class="video_text">${video.Views}</p>
      </div>`;
}

function clearAndRenderVideos(vids) {
  counter.innerText = `${vids.length} videos available`;

  grid.innerHTML = "";

  vids.forEach((v) => debounce(() => renderVideos(v), 100)());
}

function handleSearch({ target, keyCode }) {
  if (target.tagName === "INPUT" && keyCode !== 13) return;
  const searchTerm =
    target.tagName === "INPUT" ? target.value : searchInput.value;
  const filteredVideos = videos.filter((v) =>
    v.Title.toLowerCase().includes(searchTerm.toLowerCase())
  );

  clearAndRenderVideos(filteredVideos);
}

async function getDocumentaries() {
  const res = await fetch(url("/all"));
  videos = await res.json();

  clearAndRenderVideos(videos);
}

// Event handlers
searchInput.onkeydown = handleSearch;
searchBtn.onclick = handleSearch;
window.onload = getDocumentaries;
