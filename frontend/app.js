const state = {
  route: "auth",
  apiBase: localStorage.getItem("apiBase") || "http://localhost:8080/api/v1",
  accessToken: localStorage.getItem("accessToken") || "",
  refreshToken: localStorage.getItem("refreshToken") || "",
};

const app = document.getElementById("app");
const output = document.getElementById("output");
const nav = document.getElementById("top-nav");

const routes = [
  { key: "auth", label: "Sign In / Sign Up" },
  { key: "main", label: "Main" },
  { key: "admin", label: "Admin" },
];

function setOutput(payload) {
  output.textContent = typeof payload === "string" ? payload : JSON.stringify(payload, null, 2);
}

function persistTokens() {
  localStorage.setItem("accessToken", state.accessToken);
  localStorage.setItem("refreshToken", state.refreshToken);
}

async function callApi(path, method = "GET", body) {
  const headers = { "Content-Type": "application/json" };
  if (state.accessToken) headers.Authorization = `Bearer ${state.accessToken}`;

  const response = await fetch(`${state.apiBase}${path}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  });

  const maybeJson = response.status !== 204 ? await response.json().catch(() => null) : { status: "no-content" };
  if (!response.ok) {
    throw new Error(JSON.stringify({ status: response.status, ...maybeJson }, null, 2));
  }
  return maybeJson;
}

function formDataObject(form) {
  return Object.fromEntries(new FormData(form).entries());
}

function bindForm(id, handler) {
  const form = document.getElementById(id);
  if (!form) return;
  form.addEventListener("submit", async (event) => {
    event.preventDefault();
    try {
      const data = formDataObject(form);
      const result = await handler(data);
      if (result !== undefined) setOutput(result);
    } catch (error) {
      setOutput(error.message || String(error));
    }
  });
}

function renderNav() {
  nav.innerHTML = `
    <label class="api-base">
      API Base
      <input id="api-base-input" value="${state.apiBase}" />
    </label>
  `;
  routes.forEach((route) => {
    const button = document.createElement("button");
    button.textContent = route.label;
    button.onclick = () => {
      state.route = route.key;
      render();
    };
    nav.appendChild(button);
  });
  const apiInput = document.getElementById("api-base-input");
  apiInput.addEventListener("change", () => {
    state.apiBase = apiInput.value.trim();
    localStorage.setItem("apiBase", state.apiBase);
  });
}

function renderTemplate(templateId) {
  const template = document.getElementById(templateId);
  app.innerHTML = "";
  app.appendChild(template.content.cloneNode(true));
}

function bindAuth() {
  bindForm("signin-form", async (data) => {
    const payload = await callApi("/auth/login", "POST", {
      email: data.email,
      password: data.password,
    });
    state.accessToken = payload.accessToken;
    state.refreshToken = payload.refreshToken;
    persistTokens();
    return payload;
  });

  bindForm("signup-form", async (data) => {
    const payload = await callApi("/auth/register", "POST", {
      email: data.email,
      password: data.password,
      role: data.role,
    });
    state.accessToken = payload.accessToken;
    state.refreshToken = payload.refreshToken;
    persistTokens();
    return payload;
  });
}

function bindDashboard() {
  bindForm("refresh-token-form", async () => {
    const payload = await callApi("/auth/refresh", "POST", {
      refreshToken: state.refreshToken,
    });
    state.accessToken = payload.accessToken;
    state.refreshToken = payload.refreshToken;
    persistTokens();
    return payload;
  });

  document.getElementById("signout")?.addEventListener("click", () => {
    state.accessToken = "";
    state.refreshToken = "";
    persistTokens();
    setOutput("Session cleared");
  });

  bindForm("create-workspace-form", (d) =>
    callApi("/workspaces", "POST", {
      workspaceKey: d.workspaceKey,
      displayName: d.displayName,
      ownerUserId: d.ownerUserId,
    }),
  );
  bindForm("list-workspaces-form", (d) => callApi(`/users/${d.userId}/workspaces?page=${d.page}&pageSize=${d.pageSize}`));

  bindForm("add-membership-form", (d) =>
    callApi(`/workspaces/${d.workspaceId}/memberships`, "POST", {
      actorUserId: d.actorUserId,
      targetUserId: d.targetUserId,
      targetRole: d.targetRole,
      invitedByUserId: d.invitedByUserId,
    }),
  );
  bindForm("create-project-form", (d) =>
    callApi(`/workspaces/${d.workspaceId}/projects`, "POST", {
      actorUserId: d.actorUserId,
      projectKey: d.projectKey,
      displayName: d.displayName,
      description: d.description,
    }),
  );
  bindForm("list-projects-form", (d) => callApi(`/workspaces/${d.workspaceId}/projects?page=${d.page}&pageSize=${d.pageSize}`));

  bindForm("create-task-form", (d) =>
    callApi("/tasks", "POST", {
      workspaceId: d.workspaceId,
      projectId: d.projectId,
      title: d.title,
      description: d.description,
      priority: d.priority,
      assigneeUserId: d.assigneeUserId,
      createdByUserId: d.createdByUserId,
    }),
  );
  bindForm("update-task-status-form", (d) => callApi(`/tasks/${d.taskId}/status`, "PATCH", { status: d.status }));
  bindForm("list-tasks-form", (d) => callApi(`/projects/${d.projectId}/tasks?page=${d.page}&pageSize=${d.pageSize}`));

  bindForm("create-comment-form", (d) =>
    callApi("/comments", "POST", {
      workspaceId: d.workspaceId,
      projectId: d.projectId,
      taskId: d.taskId,
      authorUserId: d.authorUserId,
      body: d.body,
    }),
  );
  bindForm("update-comment-form", (d) => callApi(`/comments/${d.commentId}`, "PATCH", { actorUserId: d.actorUserId, body: d.body }));
  bindForm("delete-comment-form", (d) => callApi(`/comments/${d.commentId}`, "DELETE", { actorUserId: d.actorUserId }));
  bindForm("list-comments-form", (d) => callApi(`/tasks/${d.taskId}/comments?page=${d.page}&pageSize=${d.pageSize}`));
  bindForm("list-activities-form", (d) => callApi(`/tasks/${d.taskId}/activities?page=${d.page}&pageSize=${d.pageSize}`));
}

function render() {
  renderNav();
  if (state.route === "auth") {
    renderTemplate("auth-template");
    bindAuth();
    return;
  }
  if (state.route === "main") {
    renderTemplate("dashboard-template");
    bindDashboard();
    return;
  }
  renderTemplate("admin-template");
}

render();
setOutput("Ready. Set API Base if backend runs on a different host.");
