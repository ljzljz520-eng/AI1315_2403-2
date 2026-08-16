const links = [...document.querySelectorAll(".main-nav a")];
const current = location.pathname.split("/").pop() || "index.html";
links.forEach((link) => {
  if (link.getAttribute("href") === current) link.classList.add("active");
});

const form = document.querySelector("[data-booking-form]");
if (form) {
  const status = form.querySelector("[data-status]");
  form.addEventListener("submit", async (event) => {
    event.preventDefault();
    const data = new FormData(form);
    const payload = {
      operator: data.get("operator"),
      field: "预约偏好",
      value: `${data.get("project")} / ${data.get("date")} / ${data.get("slot")}`,
    };
    try {
      const response = await fetch("/api/records/陶艺预约-2026-001/confirm", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      if (!response.ok) throw new Error("booking failed");
      status.textContent = "预约信息已记录，我们会在一个工作日内联系你。";
      status.classList.add("show");
      form.reset();
    } catch {
      status.textContent = "演示服务暂不可用，请直接联系门店。";
      status.classList.add("show");
    }
  });
}
