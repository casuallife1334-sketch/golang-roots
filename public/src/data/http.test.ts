// @vitest-environment jsdom
import { afterEach, expect, it, vi } from "vitest";
import { z } from "zod";
import {
  request,
  setUnauthorizedHandler,
  TOKEN_KEY,
  cancelSessionRequests,
} from "./http";
afterEach(() => {
  vi.unstubAllGlobals();
  vi.restoreAllMocks();
  sessionStorage.clear();
});
it("validates successful HTTP responses", async () => {
  vi.stubGlobal(
    "fetch",
    vi.fn().mockResolvedValue(new Response(JSON.stringify({ id: 123 }))),
  );
  await expect(
    request("/test", z.object({ id: z.string() })),
  ).rejects.toMatchObject({ code: "invalid_response" });
});
it("does not end a session on a server error", async () => {
  const unauthorized = vi.fn(),
    clean = setUnauthorizedHandler(unauthorized);
  vi.stubGlobal(
    "fetch",
    vi.fn().mockResolvedValue(new Response("{}", { status: 500 })),
  );
  await expect(request("/test", z.object({}))).rejects.toMatchObject({
    status: 500,
  });
  expect(unauthorized).not.toHaveBeenCalled();
  clean();
});
it("ignores unauthorized responses from an old token", async () => {
  sessionStorage.setItem(TOKEN_KEY, "old");
  const unauthorized = vi.fn(),
    clean = setUnauthorizedHandler(unauthorized);
  vi.stubGlobal(
    "fetch",
    vi.fn().mockImplementation(async () => {
      sessionStorage.setItem(TOKEN_KEY, "new");
      return new Response("{}", { status: 401 });
    }),
  );
  await expect(request("/test", z.object({}))).rejects.toMatchObject({
    status: 401,
  });
  expect(unauthorized).not.toHaveBeenCalled();
  clean();
});
it("aborts requests when the session changes", async () => {
  vi.stubGlobal(
    "fetch",
    vi
      .fn()
      .mockImplementation(
        (_url, init) =>
          new Promise((_, reject) =>
            init.signal.addEventListener("abort", () =>
              reject(new DOMException("aborted", "AbortError")),
            ),
          ),
      ),
  );
  const promise = request("/test", z.object({}));
  cancelSessionRequests();
  await expect(promise).rejects.toMatchObject({ name: "AbortError" });
});
it("accepts an image returned as application/octet-stream", async () => {
  const png = new Uint8Array([
    0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0, 0, 0, 0,
  ]);
  vi.stubGlobal(
    "fetch",
    vi.fn().mockResolvedValue(
      new Response(png, {
        headers: { "Content-Type": "application/octet-stream" },
      }),
    ),
  );
  const result = await request<Blob>("/photo", "blob");
  expect(result.type).toBe("image/png");
});
it("rejects a binary response that is not an image", async () => {
  vi.stubGlobal(
    "fetch",
    vi.fn().mockResolvedValue(
      new Response("not an image", {
        headers: { "Content-Type": "application/octet-stream" },
      }),
    ),
  );
  await expect(request("/photo", "blob")).rejects.toMatchObject({
    code: "invalid_response",
  });
});
