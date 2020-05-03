const BASE_URL = "http://localhost:9000",
  output = document.getElementById("grid");

const url = (endpoint) => BASE_URL + endpoint;

async function getDocumentaries() {
  const res = await fetch(url("/all")),
    data = await res.json();

  data.forEach(
    (v) =>
      (output.innerHTML += `
      <div class="video">
        <a href="${v.URL}" class="video_link"
        ><img
        src="${v.Thumbnail}"
        srcset="${v.Thumbnail}"
        alt="thumbnail"
        class="video_thumbnail"
        /></a>
        <p class="video_text">${v.Title}</p>
        <p class="video_text">${v.Author}</p>
        <p class="video_text">${v.UploadDate}</p>
        <p class="video_text">${v.Length}</p>
        <p class="video_text">${v.Views}</p>
      </div>`)
  );
}

window.onload = getDocumentaries;
