// ---------- AUTH: REGISTER ----------

document.getElementById("register-form").addEventListener("submit", async (e) => {
  e.preventDefault();

  const form = document.getElementById("register-form");

  const data = {
    nickname: form.nickname.value,
    first_name: form.first_name.value,
    last_name: form.last_name.value,
    age: parseInt(form.age.value),
    gender: form.gender.value,
    email: form.email.value,
    password: form.password.value,
  };

  const res = await fetch("/api/register", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });

  if (res.ok) {
    alert("Registration successful. You can now log in.");
    showView("auth-section");
  } else {
    alert("Registration failed");
  }
});

// ---------- AUTH: LOGIN ----------

document.getElementById("login-form").addEventListener("submit", async (e) => {
  e.preventDefault();

  const form = document.getElementById("login-form");

  const data = {
    identifier: form.identifier.value,
    password: form.password.value,
  };

  const res = await fetch("/api/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });

  if (res.ok) {
    const user = await res.json();
    setAuthState(true, user.nickname || user.email || "");
    showView("feed-section");
    loadPosts();
    connectWebSocket();
  } else {
    alert("Invalid login");
  }
});

// ---------- AUTH: LOGOUT ----------

document.getElementById("nav-logout-btn").addEventListener("click", async () => {
  await fetch("/api/logout", { method: "POST" });
  setAuthState(false, "");
  showView("auth-section");
});

// ---------- UI: VIEW SWITCHING ----------

function showView(viewId) {
  const sections = document.querySelectorAll(".view-section");
  sections.forEach((sec) => {
    if (sec.id === viewId) {
      sec.classList.remove("hidden");
    } else {
      sec.classList.add("hidden");
    }
  });
}

// ---------- UI: AUTH STATE ----------

function setAuthState(isLoggedIn, username) {
  const loginBtn = document.getElementById("nav-login-btn");
  const registerBtn = document.getElementById("nav-register-btn");
  const logoutBtn = document.getElementById("nav-logout-btn");
  const homeBtn = document.getElementById("nav-home-btn");
  const newPostBtn = document.getElementById("nav-new-post-btn");
  const testBtn = document.getElementById("testBtn");
  const userLabel = document.getElementById("current-user-label");
  const chatSidebar = document.getElementById("chat-sidebar");
  const chatPanel = document.getElementById("chat-panel");

  if (isLoggedIn) {
    loginBtn.classList.add("hidden");
    registerBtn.classList.add("hidden");
    logoutBtn.classList.remove("hidden");
    homeBtn.classList.remove("hidden");
    newPostBtn.classList.remove("hidden");
    if (testBtn) testBtn.classList.remove("hidden");
    userLabel.textContent = username || "";
    userLabel.classList.remove("hidden");
    chatSidebar.classList.remove("hidden");
    chatPanel.classList.remove("hidden");
  } else {
    loginBtn.classList.remove("hidden");
    registerBtn.classList.remove("hidden");
    logoutBtn.classList.add("hidden");
    homeBtn.classList.add("hidden");
    newPostBtn.classList.add("hidden");
    if (testBtn) testBtn.classList.add("hidden");
    userLabel.classList.add("hidden");
    userLabel.textContent = "";
    chatSidebar.classList.add("hidden");
    chatPanel.classList.add("hidden");
  }
}

// ---------- NAV BUTTONS ----------

document.getElementById("nav-home-btn").addEventListener("click", () => {
  showView("feed-section");
});

document.getElementById("nav-login-btn").addEventListener("click", () => {
  showView("auth-section");
});

document.getElementById("nav-register-btn").addEventListener("click", () => {
  showView("auth-section");
});

// Toggle new post card
document.getElementById("nav-new-post-btn").addEventListener("click", () => {
  const card = document.getElementById("new-post-card");
  if (card) card.classList.toggle("hidden");
});

// ---------- INITIAL /me CHECK ----------

window.addEventListener("DOMContentLoaded", async () => {
  await loadCategories();

  try {
    const res = await fetch("/api/me");
    if (res.ok) {
      const user = await res.json();
      setAuthState(true, user.nickname || user.email || "");
      showView("feed-section");
      loadPosts();
      connectWebSocket();
    } else {
      setAuthState(false, "");
      showView("auth-section");
    }
  } catch (err) {
    console.error("Error checking /api/me:", err);
    setAuthState(false, "");
    showView("auth-section");
  }
});

// ---------- CATEGORIES ----------

async function loadCategories() {
  try {
    const res = await fetch("/api/categories");
    if (!res.ok) return;
    const categories = await res.json();

    // Populate the filter dropdown
    const select = document.getElementById("category-select");
    select.innerHTML = `<option value="">All</option>`;
    categories.forEach((cat) => {
      const opt = document.createElement("option");
      opt.value = cat.Name;
      opt.textContent = cat.Name;
      select.appendChild(opt);
    });

    // Populate the new-post checkbox group
    const group = document.getElementById("new-post-categories");
    group.innerHTML = "";
    categories.forEach((cat) => {
      const label = document.createElement("label");
      label.className = "checkbox-label";
      label.innerHTML = `
        <input type="checkbox" name="category" value="${cat.Name}" />
        <span>${cat.Name}</span>
      `;
      group.appendChild(label);
    });
  } catch (err) {
    console.error("Error loading categories:", err);
  }
}

// Filter posts when category dropdown changes
document.getElementById("category-select").addEventListener("change", () => {
  loadPosts();
});

// ---------- POSTS FEED ----------

async function loadPosts() {
  try {
    const category = document.getElementById("category-select").value;
    const url = category ? `/api/posts?category=${encodeURIComponent(category)}` : "/api/posts";
    const res = await fetch(url);
    if (!res.ok) {
      console.error("Failed to load posts");
      return;
    }
    const posts = await res.json();
    renderPosts(posts);
  } catch (err) {
    console.error("Error loading posts:", err);
  }
}

function renderPosts(posts) {
  const list = document.getElementById("posts-list");
  list.innerHTML = "";

  if (!posts || posts.length === 0) {
    list.innerHTML = `<p class="no-posts">No posts yet. Be the first to create one.</p>`;
    return;
  }

  posts.forEach((post) => {
    const categories = (post.Categories || [])
      .map((c) => `<span class="category-badge">${c}</span>`)
      .join("");

    const date = post.CreatedAt
      ? new Date(post.CreatedAt).toLocaleDateString("en-GB", { day: "numeric", month: "short", year: "numeric" })
      : "";

    const item = document.createElement("article");
    item.className = "post-item";
    item.innerHTML = `
      <h3 class="post-title">${post.Title || post.title}</h3>
      <p class="post-meta">
        By ${post.Username || post.author || "Unknown"} &middot; ${date}
      </p>
      ${categories ? `<div class="post-categories">${categories}</div>` : ""}
      <p class="post-content">${(post.Content || "").slice(0, 200)}${(post.Content || "").length > 200 ? "…" : ""}</p>
      <button class="view-post-btn" data-post-id="${post.ID || post.id}">Read more</button>
    `;
    list.appendChild(item);
  });

  list.querySelectorAll(".view-post-btn").forEach((btn) => {
    btn.addEventListener("click", () => {
      const id = btn.getAttribute("data-post-id");
      loadPostDetail(id);
    });
  });
}

// ---------- NEW POST FORM ----------

document.getElementById("new-post-form").addEventListener("submit", async (e) => {
  e.preventDefault();
  const form = e.target;

  const checkedBoxes = form.querySelectorAll('input[name="category"]:checked');
  const categories = Array.from(checkedBoxes).map((cb) => cb.value);

  const data = {
    title: form.title.value,
    content: form.content.value,
    categories,
  };

  const res = await fetch("/api/posts/create", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });

  if (res.ok) {
    form.reset();
    document.getElementById("new-post-card").classList.add("hidden");
    loadPosts();
  } else {
    document.getElementById("new-post-message").textContent = "Failed to create post.";
  }
});

// ---------- POST DETAIL ----------

async function loadPostDetail(postId) {
  try {
    const res = await fetch(`/api/posts/get?id=${encodeURIComponent(postId)}`);
    if (!res.ok) {
      console.error("Failed to load post detail");
      return;
    }
    const post = await res.json();
    renderPostDetail(post);
    showView("post-detail-section");
  } catch (err) {
    console.error("Error loading post detail:", err);
  }
}

function renderPostDetail(post) {
  const container = document.getElementById("post-detail");

  const categories = (post.Categories || [])
    .map((c) => `<span class="category-badge">${c}</span>`)
    .join("");

  const date = post.CreatedAt
    ? new Date(post.CreatedAt).toLocaleDateString("en-GB", { day: "numeric", month: "long", year: "numeric" })
    : "";

  container.innerHTML = `
    <h2>${post.Title || post.title}</h2>
    <p class="post-meta">
      By ${post.Username || post.author || "Unknown"} &middot; ${date}
    </p>
    ${categories ? `<div class="post-categories" style="margin-bottom:14px">${categories}</div>` : ""}
    <p class="post-content">${post.Content || post.content}</p>
  `;
}

document.getElementById("back-to-feed-btn").addEventListener("click", () => {
  showView("feed-section");
});

// ---------- CHAT / WEBSOCKET ----------

let chatSocket = null;
let currentChatUserId = null;

function connectWebSocket() {
  if (chatSocket && chatSocket.readyState === WebSocket.OPEN) return;

  const protocol = window.location.protocol === "https:" ? "wss" : "ws";
  const wsUrl = `${protocol}://${window.location.host}/ws`;

  chatSocket = new WebSocket(wsUrl);

  chatSocket.addEventListener("open", () => {
    console.log("WebSocket connected");
  });

  chatSocket.addEventListener("message", (event) => {
    try {
      const msg = JSON.parse(event.data);
      handleIncomingMessage(msg);
    } catch (err) {
      console.error("Invalid WS message:", event.data);
    }
  });

  chatSocket.addEventListener("close", () => {
    console.log("WebSocket closed");
  });

  chatSocket.addEventListener("error", (err) => {
    console.error("WebSocket error:", err);
  });
}

function handleIncomingMessage(msg) {
  const container = document.getElementById("chat-messages-container");

  const div = document.createElement("div");
  div.className = msg.fromMe ? "chat-message me" : "chat-message them";
  div.innerHTML = `
    <div class="chat-meta">
      <span class="chat-user">${msg.from}</span>
      <span class="chat-time">${msg.created_at || ""}</span>
    </div>
    <div class="chat-text">${msg.content}</div>
  `;
  container.appendChild(div);
  container.scrollTop = container.scrollHeight;
}

document.getElementById("chat-form").addEventListener("submit", (e) => {
  e.preventDefault();
  const input = document.getElementById("chat-message-input");
  const text = input.value.trim();
  if (!text || !chatSocket || chatSocket.readyState !== WebSocket.OPEN) return;

  const payload = {
    type: "private_message",
    to: currentChatUserId,
    content: text,
  };

  chatSocket.send(JSON.stringify(payload));
  input.value = "";
});

function openChatWith(userId, username) {
  currentChatUserId = userId;
  document.getElementById("chat-with-label").textContent = `Chat with ${username}`;
  document.getElementById("chat-form").classList.remove("hidden");
}
