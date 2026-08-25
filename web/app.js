const taskList = document.querySelector("#taskList");
const taskCount = document.querySelector("#taskCount");
const tasksEmpty = document.querySelector("#tasksEmpty");
const taskForm = document.querySelector("#taskForm");
const taskTitle = document.querySelector("#taskTitle");
const greeting = document.querySelector("#greeting");
const errorBanner = document.querySelector("#errorBanner");

function showError(message) {
  errorBanner.textContent = message;
  errorBanner.hidden = false;

  window.clearTimeout(showError.timeoutId);
  showError.timeoutId = window.setTimeout(() => {
    errorBanner.hidden = true;
  }, 3500);
}

async function request(url, options = {}) {
  const response = await fetch(url, options);

  if (!response.ok) {
    const message = await response.text();
    throw new Error(
      message.trim() || `${response.status} ${response.statusText}`
    );
  }

  if (response.status === 204) {
    return null;
  }

  return response.json();
}

function renderTasks(tasks) {
  taskList.replaceChildren();
  taskCount.textContent = tasks.length;
  tasksEmpty.hidden = tasks.length !== 0;

  for (const task of tasks) {
    const row = document.createElement("div");
    row.className = `task-item${task.completed ? " completed" : ""}`;

    const checkbox = document.createElement("input");
    checkbox.type = "checkbox";
    checkbox.checked = task.completed;
    checkbox.setAttribute("aria-label", `Выполнено: ${task.title}`);

    checkbox.addEventListener("change", async () => {
      checkbox.disabled = true;

      try {
        await request(`/api/tasks/${task.id}`, {
          method: "PATCH",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            completed: checkbox.checked,
          }),
        });

        await loadDashboard();
      } catch (error) {
        checkbox.checked = !checkbox.checked;
        showError(error.message);
      } finally {
        checkbox.disabled = false;
      }
    });

    const title = document.createElement("span");
    title.className = "task-title";
    title.textContent = task.title;

    const deleteButton = document.createElement("button");
    deleteButton.className = "delete-button";
    deleteButton.type = "button";
    deleteButton.textContent = "Удалить";

    deleteButton.addEventListener("click", async () => {
      deleteButton.disabled = true;

      try {
        await request(`/api/tasks/${task.id}`, {
          method: "DELETE",
        });

        await loadDashboard();
      } catch (error) {
        deleteButton.disabled = false;
        showError(error.message);
      }
    });

    row.append(checkbox, title, deleteButton);
    taskList.append(row);
  }
}

async function loadDashboard() {
  try {
    const data = await request("/api/dashboard");

    greeting.textContent = data.greeting;
    renderTasks(data.tasks ?? []);
  } catch (error) {
    showError(`Не удалось загрузить dashboard: ${error.message}`);
  }
}

taskForm.addEventListener("submit", async (event) => {
  event.preventDefault();

  const title = taskTitle.value.trim();

  if (!title) {
    taskTitle.focus();
    return;
  }

  const submitButton = taskForm.querySelector("button");
  submitButton.disabled = true;

  try {
    await request("/api/tasks", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ title }),
    });

    taskTitle.value = "";

    await loadDashboard();

    taskTitle.focus();
  } catch (error) {
    showError(error.message);
  } finally {
    submitButton.disabled = false;
  }
});

document
  .querySelector("#focusTaskInput")
  .addEventListener("click", () => {
    taskTitle.focus();
  });

loadDashboard();