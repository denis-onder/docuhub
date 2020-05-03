const BASE_URL = "http://localhost:9000",
  grid = document.getElementById("grid"),
  counter = document.getElementById("counter"),
  searchInput = document.getElementById("search"),
  searchBtn = document.getElementById("search_btn");

let videos = [];

const url = (endpoint) => BASE_URL + endpoint;

function renderVideos(video) {
  grid.innerHTML += `
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

  vids.forEach(renderVideos);
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

searchInput.onkeydown = handleSearch;
searchBtn.onclick = handleSearch;
window.onload = getDocumentaries;
