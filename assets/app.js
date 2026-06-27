// Progressive enhancement only: project filtering + copy buttons.
// The page is fully readable with JavaScript disabled.
(() => {
  "use strict";

  const chips = document.querySelectorAll("[data-filter]");
  const cards = document.querySelectorAll(".project-card");
  const status = document.getElementById("filter-status");

  function applyFilter(filter) {
    let shown = 0;
    cards.forEach((card) => {
      const kinds = (card.dataset.kind || "").split(" ");
      const visible = filter === "all" || kinds.includes(filter);
      card.hidden = !visible;
      if (visible) shown += 1;
    });
    if (status) {
      status.textContent = `Showing ${shown} of ${cards.length} projects`;
    }
  }

  chips.forEach((chip) => {
    chip.addEventListener("click", () => {
      chips.forEach((other) => other.setAttribute("aria-pressed", String(other === chip)));
      applyFilter(chip.dataset.filter);
    });
  });

  document.querySelectorAll("[data-copy]").forEach((button) => {
    button.addEventListener("click", async () => {
      const target = document.querySelector(button.dataset.copy);
      if (!target) return;

      const original = button.textContent;
      try {
        await navigator.clipboard.writeText(target.innerText);
        button.textContent = "Copied";
      } catch {
        // Clipboard API unavailable (e.g. non-secure context): select the text
        // so the user can copy it manually, and say so honestly.
        const range = document.createRange();
        range.selectNodeContents(target);
        const selection = window.getSelection();
        selection.removeAllRanges();
        selection.addRange(range);
        button.textContent = "Selected — press Ctrl/Cmd+C";
      }
      window.setTimeout(() => { button.textContent = original; }, 1600);
    });
  });
})();
