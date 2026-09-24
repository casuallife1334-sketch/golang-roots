import { test, expect, type Page } from "@playwright/test";
const date = "2026-01-01T00:00:00Z";
const user = {
  id: "user",
  email: "test@example.com",
  created_at: date,
  updated_at: date,
};
const tree = {
  id: "tree",
  owner_id: user.id,
  name: "Тестовое дерево",
  created_at: date,
  updated_at: date,
};
const person = (id: string) => ({
  id,
  first_name: id,
  patronymic: null as string | null,
  last_name: "Тестовый",
  metadata: { city: "Москва" },
  created_at: date,
  photo_url: null as string | null,
});
const initial = [
  "motherA",
  "fatherA",
  "motherB",
  "fatherB",
  "childA",
  "childB",
  "grandchild",
].map(person);
const relation = (id: string, a: string, b: string, type = "parent_child") => ({
  id,
  person1_id: a,
  person2_id: b,
  type,
  direction: type === "spouse" ? null : "parent",
  created_at: date,
});
const relationships = [
  relation("ab", "motherA", "fatherA", "spouse"),
  relation("cd", "motherB", "fatherB", "spouse"),
  relation("xy", "childA", "childB", "spouse"),
  relation("ax", "motherA", "childA"),
  relation("bx", "fatherA", "childA"),
  relation("cy", "motherB", "childB"),
  relation("dy", "fatherB", "childB"),
  relation("xz", "childA", "grandchild"),
  relation("yz", "childB", "grandchild"),
];
async function fixture(
  page: Page,
  meStatus = 200,
  photoFailuresBeforeSuccess = 1,
) {
  const errors: string[] = [];
  page.on("pageerror", (error) => errors.push(error.message));
  let people = structuredClone(initial);
  let creates = 0;
  let photoFailures = photoFailuresBeforeSuccess;
  await page.addInitScript(() =>
    sessionStorage.setItem("roots:access-token", "test-token"),
  );
  await page.route("**/api/v1/**", async (route) => {
    const path = new URL(route.request().url()).pathname,
      method = route.request().method();
    const json = (value: unknown, status = 200) =>
      route.fulfill({
        status,
        contentType: "application/json",
        body: JSON.stringify(value),
      });
    if (path.endsWith("/users/me"))
      return json(meStatus === 200 ? user : {}, meStatus);
    if (path.endsWith("/trees")) return json([tree]);
    if (path.endsWith("/relationships")) return json(relationships);
    if (path.endsWith("/photo")) {
      if (method === "POST") {
        if (photoFailures-- > 0) return json({}, 500);
        const id = path.split("/").at(-2)!;
        const saved = people.find((item) => item.id === id)!;
        saved.photo_url = "/photo";
        return json(saved);
      }
      return route.fulfill({
        status: 200,
        contentType: "application/octet-stream",
        body: Buffer.from(
          "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+jP1sAAAAASUVORK5CYII=",
          "base64",
        ),
      });
    }
    if (path.endsWith("/persons")) {
      if (method === "POST") {
        creates++;
        const saved = {
          ...person(`created${creates}`),
          ...route.request().postDataJSON(),
        };
        people.push(saved);
        return json(saved, 201);
      }
      return json(people);
    }
    if (method === "PATCH" && path.includes("/persons/")) {
      const id = path.split("/").at(-1)!;
      people = people.map((item) =>
        item.id === id ? { ...item, ...route.request().postDataJSON() } : item,
      );
      return json(people.find((item) => item.id === id));
    }
    return json({}, 404);
  });
  return {
    errors,
    get creates() {
      return creates;
    },
  };
}
test("cards retain coordinates on selection, family paths and modal keyboard work", async ({
  page,
}) => {
  const state = await fixture(page);
  await page.goto("/trees/tree");
  const cards = page.locator(".react-flow__node-person");
  await expect(cards).toHaveCount(7);
  const positions = () =>
    cards.evaluateAll((nodes) =>
      nodes.map((node) => [
        node.getAttribute("data-id"),
        (node as HTMLElement).style.transform,
      ]),
    );
  const before = await positions();
  await page.locator('[data-id="childA"]').click();
  await expect(page.locator(".detail-panel h2")).toContainText("childA");
  await page.locator('[data-id="childB"]').click();
  // Opening the panel temporarily virtualizes cards until ResizeObserver fits the viewport.
  await expect.poll(positions).toEqual(before);
  await expect(page.locator(".react-flow__edge-family")).toHaveCount(3);
  const sizes = await cards.evaluateAll((nodes) =>
    nodes.map((node) => [
      (node as HTMLElement).offsetWidth,
      (node as HTMLElement).offsetHeight,
    ]),
  );
  expect(
    sizes.every(([width, height]) => width === 240 && height === 132),
  ).toBe(true);
  await page
    .getByRole("button", { name: "Добавить человека", exact: true })
    .click();
  await expect(page.getByRole("dialog")).toBeVisible();
  await page.keyboard.press("Escape");
  await expect(page.getByRole("dialog")).toHaveCount(0);
  expect(state.errors).toEqual([]);
  await page.screenshot({
    path: "test-results/tree-reviewed.png",
    fullPage: true,
  });
});
test("failed photo upload retries without creating a duplicate", async ({
  page,
}) => {
  const state = await fixture(page);
  await page.goto("/trees/tree");
  await page
    .getByRole("button", { name: "Добавить человека", exact: true })
    .click();
  await page.getByLabel("Имя", { exact: true }).fill("Новый");
  await page.getByLabel("Фамилия", { exact: true }).fill("Человек");
  const png = await page.evaluate(() => {
    const canvas = document.createElement("canvas");
    canvas.width = 20;
    canvas.height = 30;
    canvas.getContext("2d")!.fillRect(0, 0, 20, 30);
    return canvas.toDataURL().split(",")[1];
  });
  await page.getByLabel("Фотография", { exact: true }).setInputFiles({
    name: "photo.png",
    mimeType: "image/png",
    buffer: Buffer.from(png, "base64"),
  });
  await page.getByRole("button", { name: "Применить фото" }).click();
  await expect(
    page.locator(".person-photo-control .photo-control-preview"),
  ).toBeVisible();
  await page
    .getByRole("dialog")
    .getByRole("button", { name: "Добавить человека", exact: true })
    .click();
  await expect(page.getByRole("alert")).toContainText("Человек сохранён");
  await page.getByRole("button", { name: "Сохранить", exact: true }).click();
  await expect(page.getByRole("dialog")).toHaveCount(0);
  expect(state.creates).toBe(1);
  expect(state.errors).toEqual([]);
});
test("uploaded photo appears immediately when the API returns binary content", async ({
  page,
}) => {
  const state = await fixture(page, 200, 0);
  await page.goto("/trees/tree?person=childA");
  const png = await page.evaluate(() => {
    const canvas = document.createElement("canvas");
    canvas.width = 80;
    canvas.height = 120;
    const context = canvas.getContext("2d")!;
    context.fillStyle = "#3478ff";
    context.fillRect(0, 0, 80, 120);
    return canvas.toDataURL().split(",")[1];
  });
  await page.locator('.profile-photo input[type="file"]').setInputFiles({
    name: "portrait.png",
    mimeType: "image/png",
    buffer: Buffer.from(png, "base64"),
  });
  await expect(page.getByRole("dialog")).toBeVisible();
  await expect(page.getByLabel("Область обрезки фотографии")).toBeVisible();
  await expect(page.getByLabel("Масштаб фотографии")).toHaveValue("1");
  await expect(page.getByText("Затемнённые края")).toBeVisible();
  await page.getByLabel("Увеличить масштаб").click();
  await expect(page.getByLabel("Масштаб фотографии")).not.toHaveValue("1");
  await page.screenshot({
    path: "test-results/crop-reviewed.png",
    fullPage: true,
  });
  await page.getByRole("button", { name: "Применить фото" }).click();
  await expect(page.getByRole("dialog")).toHaveCount(0);
  await expect(page.locator(".profile-photo img")).toBeVisible();
  await page
    .getByRole("button", { name: "Редактировать", exact: true })
    .click();
  await expect(page.locator(".person-photo-control")).toBeVisible();
  await expect(
    page.locator(".person-photo-control").getByRole("button", {
      name: "Удалить",
      exact: true,
    }),
  ).toBeVisible();
  await expect(
    page.getByRole("button", { name: "Загрузить фотографию", exact: true }),
  ).toHaveCount(0);
  expect(state.errors).toEqual([]);
});
test("tree selectors keep their controls separated", async ({ page }) => {
  await fixture(page);
  await page.goto("/trees/tree");
  const switcher = page.locator(".tree-switcher");
  await expect(switcher).toBeVisible();
  await expect(switcher.getByLabel("Создать дерево")).toBeVisible();
  await page.goto("/settings");
  const title = page.locator(".tree-settings .card-title");
  const selector = page.locator(".tree-settings > .field").first();
  await expect(title).toBeVisible();
  await expect(selector).toBeVisible();
  const titleBox = await title.boundingBox();
  const selectorBox = await selector.boundingBox();
  expect(selectorBox!.y).toBeGreaterThan(titleBox!.y + titleBox!.height);
  await page.screenshot({
    path: "test-results/settings-reviewed.png",
    fullPage: true,
  });
});
test("network/server failure during session check preserves token", async ({
  page,
}) => {
  await fixture(page, 500);
  await page.goto("/trees/tree");
  await expect(page.getByRole("alert")).toContainText(
    "Не удалось проверить сессию",
  );
  expect(
    await page.evaluate(() => sessionStorage.getItem("roots:access-token")),
  ).toBe("test-token");
  await expect(page).toHaveURL(/\/trees\/tree$/);
});
test("patronymic, comment and cleared metadata survive reload", async ({
  page,
}) => {
  const state = await fixture(page);
  await page.goto("/trees/tree?person=childA");
  await page
    .getByRole("button", { name: "Редактировать", exact: true })
    .click();
  await page.getByLabel("Отчество", { exact: true }).fill("Петрович");
  await page
    .getByLabel("Комментарий", { exact: true })
    .fill("Первая строка\nВторая строка");
  await page.getByLabel("Место / город", { exact: true }).fill("");
  await page.getByRole("button", { name: "Сохранить", exact: true }).click();
  await expect(page.getByRole("dialog")).toHaveCount(0);
  await expect(page.locator(".detail-panel h2")).toContainText(
    "childA Петрович Тестовый",
  );
  await expect(page.locator(".detail-panel")).not.toContainText("Москва");
  await expect(page.locator(".person-comment-view")).toHaveText(
    "Первая строка\nВторая строка",
  );
  await page.reload();
  await expect(page.locator(".detail-panel h2")).toContainText(
    "childA Петрович Тестовый",
  );
  await expect(page.locator(".detail-panel")).not.toContainText("Москва");
  await expect(page.locator(".person-comment-view")).toHaveText(
    "Первая строка\nВторая строка",
  );
  await page
    .getByRole("button", { name: "Редактировать", exact: true })
    .click();
  await page.getByLabel("Отчество", { exact: true }).fill("");
  await page.getByLabel("Комментарий", { exact: true }).fill("");
  await page.getByRole("button", { name: "Сохранить", exact: true }).click();
  await expect(page.locator(".detail-panel h2")).toHaveText("childA Тестовый");
  await expect(page.locator(".person-comment-view")).toHaveCount(0);
  expect(state.errors).toEqual([]);
});
test("comment editor is responsive, limited and scrolls without a visible bar", async ({
  page,
}) => {
  await fixture(page);
  await page.goto("/trees/tree");
  await page
    .getByRole("button", { name: "Добавить человека", exact: true })
    .click();
  const main = page.locator(".person-main-fields");
  const comment = page.locator(".person-comment-section");
  const desktopMain = await main.boundingBox();
  const desktopComment = await comment.boundingBox();
  expect(desktopComment!.x).toBeGreaterThan(
    desktopMain!.x + desktopMain!.width,
  );
  const firstNameBox = await page
    .getByLabel("Имя", { exact: true })
    .boundingBox();
  const birthBox = await page
    .getByLabel("Дата рождения", { exact: true })
    .boundingBox();
  const deathBox = await page
    .getByLabel("Дата смерти", { exact: true })
    .boundingBox();
  expect(birthBox!.y).toBe(deathBox!.y);
  expect(birthBox!.y).toBeGreaterThan(firstNameBox!.y + firstNameBox!.height);
  await expect(
    page.getByRole("button", { name: "Загрузить фотографию", exact: true }),
  ).toBeVisible();

  const textarea = page.getByLabel("Комментарий", { exact: true });
  await textarea.fill("Строка комментария\n".repeat(400));
  await expect(textarea).toHaveValue(/Строка комментария/);
  expect(await textarea.inputValue()).toHaveLength(5000);
  const scroll = await textarea.evaluate((element) => {
    const field = element as HTMLTextAreaElement;
    field.scrollTop = field.scrollHeight;
    return {
      scrollable: field.scrollHeight > field.clientHeight,
      scrolled: field.scrollTop > 0,
      scrollbarWidth: getComputedStyle(field).scrollbarWidth,
    };
  });
  expect(scroll.scrollable).toBe(true);
  expect(scroll.scrolled).toBe(true);
  expect(scroll.scrollbarWidth).toBe("none");
  await page.screenshot({
    path: "test-results/comment-editor-desktop-reviewed.png",
    fullPage: true,
  });

  await page.setViewportSize({ width: 600, height: 820 });
  const mobileMain = await main.boundingBox();
  const mobileComment = await comment.boundingBox();
  expect(mobileComment!.y).toBeGreaterThan(mobileMain!.y + mobileMain!.height);
  await page.screenshot({
    path: "test-results/comment-editor-reviewed.png",
    fullPage: true,
  });
});
test("logout clears cached trees before signing in again", async ({ page }) => {
  await fixture(page);
  await page.goto("/trees/tree");
  await expect(page.locator(".react-flow__node-person")).toHaveCount(7);
  await page.getByRole("button", { name: "Выйти", exact: true }).click();
  await expect(page).toHaveURL(/\/login(?:\?|$)/);
  expect(
    await page.evaluate(() => sessionStorage.getItem("roots:access-token")),
  ).toBeNull();
  await page.route("**/api/v1/auth/login", (route) =>
    route.fulfill({
      json: {
        access_token: "second-session",
        token_type: "Bearer",
        expires_in: 900,
      },
    }),
  );
  await page.route("**/api/v1/trees", (route) => route.fulfill({ json: [] }));
  await page.getByLabel("Email", { exact: true }).fill("test@example.com");
  await page.getByLabel("Пароль", { exact: true }).fill("password123");
  await page.getByRole("button", { name: "Войти", exact: true }).click();
  await expect(
    page.getByRole("heading", { name: "Создайте первое дерево" }),
  ).toBeVisible();
  await expect(page.locator(".react-flow__node-person")).toHaveCount(0);
});
test("a failed people request shows a retry state instead of an empty tree", async ({
  page,
}) => {
  await fixture(page);
  await page.route("**/api/v1/trees/tree/persons", (route) =>
    route.fulfill({ status: 500, json: {} }),
  );
  await page.goto("/trees/tree");
  await expect(page.getByRole("alert")).toContainText(
    "Сервис временно недоступен",
  );
  await expect(
    page.getByRole("heading", { name: "Дерево пока пусто" }),
  ).toHaveCount(0);
  await page.unroute("**/api/v1/trees/tree/persons");
  await page.getByRole("button", { name: "Повторить", exact: true }).click();
  await expect(page.locator(".react-flow__node-person")).toHaveCount(7);
});
