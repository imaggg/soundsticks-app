// Bridge: Go backend calls window.updateState(json), UI calls window.sendCmd(cmd, payload)
// During design/preview, we use mock data.

const MOCK = typeof window.sendCmd === "undefined";

// ── i18n ──────────────────────────────────────────────────────────────────────

const TRANSLATIONS = {
  en: {
    "setup.title": "One-time setup",
    "setup.desc":
      "To connect to your SoundSticks, we need to extract a security certificate from the Harman Kardon One app.",
    "setup.select_apk": "Select <strong>.xapk</strong>",
    "setup.hint": "Harman Kardon One app from APKPure",
    "setup.browse_apk": "Browse file…",
    "setup.done": "Certificate imported. Looking for SoundSticks…",
    "searching.text": "Looking for SoundSticks on your network…",
    "searching.not_found":
      "Could not find SoundSticks. Check WiFi, then use ↺ to retry.",
    "tab.general": "General",
    "tab.lights": "Lights",
    "tab.eq": "EQ",
    "general.volume": "Volume",
    "lights.label": "Lights",
    "lights.brightness": "Brightness",
    "lights.speed": "Speed",
    "lights.pattern": "Pattern",
    "lights.color": "Color",
    "pattern.ocean": "Ocean",
    "pattern.aurora": "Aurora",
    "pattern.blossom": "Blossom",
    "pattern.sunrise": "Sunrise",
    "pattern.fireplace": "Fireplace",
    "pattern.calm": "Calm",
    "pattern.nebula": "Nebula",
    "eq.preset": "Preset",
    "eq.custom_section": "Custom",
    "eq.preset.signature": "Signature",
    "eq.preset.vocal": "Vocal",
    "eq.preset.energetic": "Energetic",
    "eq.preset.chill": "Chill",
    "eq.preset.custom": "Custom",
    "eq.reset": "Reset",
    "eq.save": "Save",
    "tab.moment": "Moment",
    "moment.label": "Moment",
    "moment.playtime": "Play Time",
    "moment.unlimited": "Unlimited",
    "moment.min": "min",
    "moment.mode": "Mode",
    "moment.forest": "Forest",
    "moment.rain": "Rain",
    "moment.ocean": "Ocean",
    "moment.city": "City",
    "moment.forest.0": "Insects",
    "moment.forest.1": "Creek",
    "moment.forest.2": "Birds",
    "moment.rain.0": "Lightning",
    "moment.rain.1": "Frogs",
    "moment.rain.2": "Raindrops",
    "moment.ocean.0": "Seabirds",
    "moment.ocean.1": "Currents",
    "moment.ocean.2": "Waves",
    "moment.city.0": "Explore",
    "moment.city.1": "Traffic",
    "moment.city.2": "Ambience",
  },
  es: {
    "setup.title": "Configuración inicial",
    "setup.desc":
      "Para conectarte a tus SoundSticks, necesitamos extraer un certificado de seguridad de la app Harman Kardon One.",
    "setup.select_apk": "Seleccionar <strong>.xapk</strong>",
    "setup.hint": "App Harman Kardon One o APKPure",
    "setup.browse_apk": "Buscar archivo…",
    "setup.done": "Certificado importado. Buscando SoundSticks…",
    "searching.text": "Buscando SoundSticks en tu red…",
    "searching.not_found":
      "No se encontró SoundSticks. Verifica el WiFi y usa ↺ para reintentar.",
    "tab.general": "General",
    "tab.lights": "Luces",
    "tab.eq": "Ecualizador",
    "general.volume": "Volumen",
    "lights.label": "Luces",
    "lights.brightness": "Brillo",
    "lights.speed": "Velocidad",
    "lights.pattern": "Patrón",
    "lights.color": "Color",
    "pattern.ocean": "Océano",
    "pattern.aurora": "Aurora",
    "pattern.blossom": "Flor",
    "pattern.sunrise": "Amanecer",
    "pattern.fireplace": "Chimenea",
    "pattern.calm": "Calma",
    "pattern.nebula": "Nebulosa",
    "eq.preset": "Ajuste predeterminado",
    "eq.custom_section": "Personalizado",
    "eq.preset.signature": "Firma",
    "eq.preset.vocal": "Vocal",
    "eq.preset.energetic": "Energético",
    "eq.preset.chill": "Tranquilo",
    "eq.preset.custom": "Personalizado",
    "eq.reset": "Restablecer",
    "eq.save": "Guardar",
    "tab.moment": "Momento",
    "moment.label": "Momento",
    "moment.playtime": "Tiempo",
    "moment.unlimited": "Ilimitado",
    "moment.min": "min",
    "moment.mode": "Modo",
    "moment.forest": "Bosque",
    "moment.rain": "Lluvia",
    "moment.ocean": "Océano",
    "moment.city": "Ciudad",
    "moment.forest.0": "Insectos",
    "moment.forest.1": "Arroyo",
    "moment.forest.2": "Pájaros",
    "moment.rain.0": "Relámpago",
    "moment.rain.1": "Ranas",
    "moment.rain.2": "Gotas",
    "moment.ocean.0": "Aves marinas",
    "moment.ocean.1": "Corrientes",
    "moment.ocean.2": "Olas",
    "moment.city.0": "Explorar",
    "moment.city.1": "Tráfico",
    "moment.city.2": "Ambiente",
  },
  fr: {
    "setup.title": "Configuration initiale",
    "setup.desc":
      "Pour vous connecter à vos SoundSticks, nous devons extraire un certificat de sécurité de l'application Harman Kardon One.",
    "setup.select_apk": "Sélectionner <strong>.xapk</strong>",
    "setup.hint": "Application Harman Kardon One depuis APKPure",
    "setup.browse_apk": "Parcourir…",
    "setup.done": "Certificat importé. Recherche de SoundSticks…",
    "searching.text": "Recherche de SoundSticks sur votre réseau…",
    "searching.not_found":
      "SoundSticks introuvable. Vérifiez le WiFi, puis utilisez ↺ pour réessayer.",
    "tab.general": "Général",
    "tab.lights": "Lumières",
    "tab.eq": "Égaliseur",
    "general.volume": "Volume",
    "lights.label": "Lumières",
    "lights.brightness": "Luminosité",
    "lights.speed": "Vitesse",
    "lights.pattern": "Motif",
    "lights.color": "Couleur",
    "pattern.ocean": "Océan",
    "pattern.aurora": "Aurore",
    "pattern.blossom": "Floraison",
    "pattern.sunrise": "Lever du soleil",
    "pattern.fireplace": "Cheminée",
    "pattern.calm": "Calme",
    "pattern.nebula": "Nébuleuse",
    "eq.preset": "Préréglage",
    "eq.custom_section": "Personnalisé",
    "eq.preset.signature": "Signature",
    "eq.preset.vocal": "Vocal",
    "eq.preset.energetic": "Énergique",
    "eq.preset.chill": "Calme",
    "eq.preset.custom": "Personnalisé",
    "eq.reset": "Réinitialiser",
    "eq.save": "Enregistrer",
    "tab.moment": "Moment",
    "moment.label": "Moment",
    "moment.playtime": "Durée",
    "moment.unlimited": "Illimité",
    "moment.min": "min",
    "moment.mode": "Mode",
    "moment.forest": "Forêt",
    "moment.rain": "Pluie",
    "moment.ocean": "Océan",
    "moment.city": "Ville",
    "moment.forest.0": "Insectes",
    "moment.forest.1": "Ruisseau",
    "moment.forest.2": "Oiseaux",
    "moment.rain.0": "Éclair",
    "moment.rain.1": "Grenouilles",
    "moment.rain.2": "Gouttes",
    "moment.ocean.0": "Oiseaux marins",
    "moment.ocean.1": "Courants",
    "moment.ocean.2": "Vagues",
    "moment.city.0": "Explorer",
    "moment.city.1": "Circulation",
    "moment.city.2": "Ambiance",
  },
  uk: {
    "setup.title": "Початкове налаштування",
    "setup.desc":
      "Щоб підключитися до SoundSticks 5, потрібно витягти сертифікат безпеки з додатку Harman Kardon One.",
    "setup.select_apk": "Вибрати <strong>.xapk</strong>",
    "setup.hint": "Додаток Harman Kardon One з APKPure",
    "setup.browse_apk": "Вибрати файл…",
    "setup.done": "Сертифікат імпортовано. Шукаємо SoundSticks…",
    "searching.text": "Пошук SoundSticks у мережі…",
    "searching.not_found":
      "SoundSticks не знайдено. Перевірте WiFi та натисніть ↺ для повтору.",
    "tab.general": "Загальні",
    "tab.lights": "Підсвітка",
    "tab.eq": "Еквалайзер",
    "general.volume": "Гучність",
    "lights.label": "Підсвітка",
    "lights.brightness": "Яскравість",
    "lights.speed": "Швидкість",
    "lights.pattern": "Патерн",
    "lights.color": "Колір",
    "pattern.ocean": "Океан",
    "pattern.aurora": "Аврора",
    "pattern.blossom": "Цвітіння",
    "pattern.sunrise": "Світанок",
    "pattern.fireplace": "Камін",
    "pattern.calm": "Колір",
    "pattern.nebula": "Туманність",
    "eq.preset": "Пресет",
    "eq.custom_section": "Кастом",
    "eq.preset.signature": "Signature",
    "eq.preset.vocal": "Вокал",
    "eq.preset.energetic": "Енергійний",
    "eq.preset.chill": "Релакс",
    "eq.preset.custom": "Кастом",
    "eq.reset": "Скинути",
    "eq.save": "Зберегти",
    "tab.moment": "Момент",
    "moment.label": "Момент",
    "moment.playtime": "Час відтворення",
    "moment.unlimited": "Необмежено",
    "moment.min": "хв",
    "moment.mode": "Режим",
    "moment.forest": "Ліс",
    "moment.rain": "Дощ",
    "moment.ocean": "Океан",
    "moment.city": "Місто",
    "moment.forest.0": "Комахи",
    "moment.forest.1": "Річка",
    "moment.forest.2": "Птахи",
    "moment.rain.0": "Блискавка",
    "moment.rain.1": "Жаби",
    "moment.rain.2": "Краплі дощу",
    "moment.ocean.0": "Морські птахи",
    "moment.ocean.1": "Течія",
    "moment.ocean.2": "Хвилі",
    "moment.city.0": "Навколо",
    "moment.city.1": "Трафік",
    "moment.city.2": "Атмосфера",
  },
};

let currentLang = "en";

// ── Theme ─────────────────────────────────────────────────────────────────────

window.applyTheme = function (theme) {
  document.documentElement.dataset.theme = theme; // "light" | "dark" | "auto"
};

// Apply saved theme immediately (before Go sends state) so there's no flash.
window.applyTheme("auto");

window.t = function (key) {
  return (
    (TRANSLATIONS[currentLang] || TRANSLATIONS.en)[key] ??
    TRANSLATIONS.en[key] ??
    key
  );
};

window.setLanguage = function (lang) {
  if (!TRANSLATIONS[lang]) lang = "en";
  currentLang = lang;
  document.querySelectorAll("[data-i18n]").forEach((el) => {
    el.textContent = window.t(el.dataset.i18n);
  });
  document.querySelectorAll("[data-i18n-html]").forEach((el) => {
    el.innerHTML = window.t(el.dataset.i18nHtml);
  });
  document.querySelectorAll("[data-i18n-placeholder]").forEach((el) => {
    el.placeholder = window.t(el.dataset.i18nPlaceholder);
  });
  document.querySelectorAll(".timer-btn[data-timer]").forEach((btn) => {
    const sec = parseInt(btn.dataset.timer);
    if (sec > 0) btn.textContent = sec / 60 + " " + window.t("moment.min");
  });
};

const state = {
  deviceName: "SoundSticks 5 Wi-Fi",
  language: "en",
  connected: true,
  lights: {
    enable: true,
    brightness: 80,
    speed: 2,
    patternId: 1,
    colorLevel: 50,
  },
  eq: { presetId: 1, custom: [0, 0, 0, 0, 0, 0, 0] },
  moment: { enabled: false, soundscapeId: 1, sleepTimer: 0 },
  player: { vol: 70, mute: false },
};

const EQ_FREQS = ["125", "250", "500", "1k", "2k", "4k", "8k"];

const PATTERN_GRADIENT = {
  1: "linear-gradient(to right, #7cdac2, #5337fa)", // Ocean
  2: "linear-gradient(to right, #8ae474, #528ed7)", // Aurora
  3: "linear-gradient(to right, #c082c2, #7e56e2)", // Blossom
  4: "linear-gradient(to right, #df5b5a, #e4ec75)", // Sunrise
  5: "linear-gradient(to right, #d9b85f, #d54712)", // Fireplace
  6: "linear-gradient(to right, #fff, #ebdb85, #9ee87c, #ace7ef,#608ae9 ,#bb6be0,#f5ddd2)", // Calm
  7: "linear-gradient(to right, #dfb965, #9eed70, #93e1eb , #76a0ea, #6064e8, #c05ee3, #c35a90,#dca867)", // Nebula
};

// ── Volume icon ───────────────────────────────────────────────────────────────

const _S =
  'stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"';
const _BODY =
  '<path d="M1 6H4L8 3V13L4 10H1Z" fill="currentColor" stroke="none"/>';
const VOL_ICONS = {
  mute: `<svg width="24" height="24" viewBox="0 0 16 16" fill="none" ${_S}>${_BODY}<line x1="10.5" y1="5.5" x2="14.5" y2="10.5"/><line x1="10.5" y1="10.5" x2="14.5" y2="5.5"/></svg>`,
  low: `<svg width="24" height="24" viewBox="0 0 16 16" fill="none" ${_S}>${_BODY}<path d="M9 6.3A2 2 0 0 1 9 9.7"/></svg>`,
  mid: `<svg width="24" height="24" viewBox="0 0 16 16" fill="none" ${_S}>${_BODY}<path d="M9 6.3A2 2 0 0 1 9 9.7"/><path d="M9.75 5A3.5 3.5 0 0 1 9.75 11"/></svg>`,
  high: `<svg width="24" height="24" viewBox="0 0 16 16" fill="none" ${_S}>${_BODY}<path d="M9 6.3A2 2 0 0 1 9 9.7"/><path d="M9.75 5A3.5 3.5 0 0 1 9.75 11"/><path d="M10.5 3.5A5.5 5.5 0 0 1 10.5 12.5"/></svg>`,
};

const volIconEl = document.getElementById("vol-icon");

let playerVol = 70;
let playerMute = false;

function updateVolIcon(vol, mute) {
  playerVol = vol;
  playerMute = mute;
  let key, label;
  if (mute) {
    key = "mute";
    label = "Muted";
  } else if (vol < 30) {
    key = "low";
    label = `Volume: ${vol}%`;
  } else if (vol < 75) {
    key = "mid";
    label = `Volume: ${vol}%`;
  } else {
    key = "high";
    label = `Volume: ${vol}%`;
  }
  volIconEl.innerHTML = VOL_ICONS[key];
  volIconEl.title = label;
}

volIconEl.addEventListener("click", () => {
  playerMute = !playerMute;
  updateVolIcon(playerVol, playerMute);
  cmd("toggleMute", { mute: playerMute });
});

// ── Volume popup slider (header) ─────────────────────────────────────────────

const volSliderEl = document.getElementById("vol-slider");
const volValueEl = document.getElementById("vol-value");
const volPopupEl = document.querySelector(".vol-popup");

volSliderEl?.addEventListener("input", () => {
  volValueEl.textContent = volSliderEl.value + "%";
  updateVolIcon(+volSliderEl.value, playerMute);
});
volSliderEl?.addEventListener("change", () => {
  playerVol = +volSliderEl.value;
  cmd("setVol", { vol: playerVol });
});

// Keep popup visible during drag (mouse can leave .header-right while dragging).
volSliderEl?.addEventListener("mousedown", () =>
  volPopupEl?.classList.add("dragging"),
);
document.addEventListener("mouseup", () =>
  volPopupEl?.classList.remove("dragging"),
);

const SOUNDSCAPE_ELEMENTS = {
  1: ["moment.forest.0", "moment.forest.1", "moment.forest.2"],
  2: ["moment.rain.0", "moment.rain.1", "moment.rain.2"],
  3: ["moment.ocean.0", "moment.ocean.1", "moment.ocean.2"],
  4: ["moment.city.0", "moment.city.1", "moment.city.2"],
};

function updateColorStrip(patternId) {
  const grad = PATTERN_GRADIENT[+patternId] ?? PATTERN_GRADIENT[1];
  colorHue.style.background = grad;
}

function cmd(command, payload) {
  if (MOCK) {
    console.log("cmd", command, payload);
    return;
  }
  window.sendCmd(command, JSON.stringify(payload ?? {}));
}

// ── Setup screen ──────────────────────────────────────────────────────────────

const setupStatus = document.getElementById("setup-status");

function setupShowStatus(ok, msg) {
  setupStatus.className = "setup-status " + (ok ? "ok" : "err");
  setupStatus.textContent = msg;
  setupStatus.classList.remove("hidden");
}

function isApkFile(name) {
  return name.endsWith(".apk") || name.endsWith(".xapk");
}

// Browse APK — opens native file dialog via Go (or hidden input in preview)
document.getElementById("btn-browse-apk")?.addEventListener("click", () => {
  if (MOCK) {
    document.getElementById("preview-apk-input").click();
  } else {
    cmd("setupBrowseApk", {});
  }
});

document
  .getElementById("preview-apk-input")
  ?.addEventListener("change", (e) => {
    const file = e.target.files[0];
    if (file && isApkFile(file.name)) {
      setupShowStatus(
        true,
        `Selected: ${file.name} — in production, extraction happens here`,
      );
    }
  });

// Called by Go after successful extraction
window.setupDone = function () {
  setupShowStatus(true, window.t("setup.done"));
};

// ── Tab navigation ────────────────────────────────────────────────────────────

document.querySelectorAll(".tab").forEach((btn) => {
  btn.addEventListener("click", () => {
    document
      .querySelectorAll(".tab")
      .forEach((t) => t.classList.remove("active"));
    document
      .querySelectorAll(".tab-content")
      .forEach((t) => t.classList.remove("active"));
    btn.classList.add("active");
    document.getElementById("tab-" + btn.dataset.tab).classList.add("active");
  });
});

// ── Lights ────────────────────────────────────────────────────────────────────

const lightsToggle = document.getElementById("lights-toggle");
const lightsControls = document.getElementById("lights-controls");
const brightnessEl = document.getElementById("brightness");
const brightnessVal = document.getElementById("brightness-value");
const speedEl = document.getElementById("speed");
const speedVal = document.getElementById("speed-value");
const colorRow = document.getElementById("color-row");
const colorHue = document.getElementById("color-hue");
const speedSection = document.getElementById("speed-section");

// Pattern 6 (Calm) has no animation — lock speed slider when it's active.
function updateSpeedLock(patternId) {
  const locked = patternId === 6;
  speedSection.style.opacity = locked ? "0.4" : "";
  speedSection.style.pointerEvents = locked ? "none" : "";
}

// Per-pattern saved color level (0-100). Restored when switching patterns.
const patternColorLevel = { 1: 50, 2: 50, 3: 50, 4: 50, 5: 50, 6: 50, 7: 50 };
let activePatternId = 1;

lightsToggle.addEventListener("change", () => {
  const on = lightsToggle.checked;
  lightsControls.style.opacity = on ? "1" : "0.4";
  lightsControls.style.pointerEvents = on ? "" : "none";
  cmd("setLightInfo", { enable: on ? "1" : "0" });
});

brightnessEl.addEventListener("input", () => {
  brightnessVal.textContent = brightnessEl.value + "%";
});
brightnessEl.addEventListener("change", () => {
  cmd("setLightInfo", { brightness: brightnessEl.value });
});

speedEl.addEventListener("input", () => {
  speedVal.textContent = speedEl.value;
});
speedEl.addEventListener("change", () => {
  cmd("setLightInfo", { dynamic_level: speedEl.value });
});

document.querySelectorAll(".pattern-btn").forEach((btn) => {
  btn.addEventListener("click", () => {
    document
      .querySelectorAll(".pattern-btn")
      .forEach((b) => b.classList.remove("active"));
    btn.classList.add("active");
    const id = +btn.dataset.id;
    activePatternId = id;
    colorHue.value = patternColorLevel[id] ?? 50;
    updateColorStrip(id);
    updateSpeedLock(id);
    cmd("setLightInfo", {
      active_pattern: { id: btn.dataset.id, level: colorHue.value },
    });
  });
});

colorHue.addEventListener("change", () => {
  patternColorLevel[activePatternId] = +colorHue.value;
  cmd("setLightInfo", {
    active_pattern: { id: String(activePatternId), level: colorHue.value },
  });
});

// ── EQ ────────────────────────────────────────────────────────────────────────

const customEq = document.getElementById("custom-eq");
const eqBands = document.getElementById("eq-bands");
const customNameEl = document.getElementById("custom-eq-name");

const FS = [125, 250, 500, 1000, 2000, 4000, 8000];

const BAND_MIN = [-12, -12, -12, -12, -12, -12, -12];

// activeSlot: null=native preset, -1=official custom (id=0), 0/1/2=user custom slot
let activeSlot = null;

// Official custom gains in UI units (±12). Seeded from device; synced from polls
// whenever a user-custom slot is NOT active (EQID=0 then belongs to us or phone).
let officialCustomGains = [0, 0, 0, 0, 0, 0, 0];

// Per-slot gains and names (populated from Go config on connect)
const customGains = [
  [0, 0, 0, 0, 0, 0, 0],
  [0, 0, 0, 0, 0, 0, 0],
  [0, 0, 0, 0, 0, 0, 0],
];
const customNames = ["Custom 1", "Custom 2", "Custom 3"];

// Build 7-band sliders
const bandSliders = [];
EQ_FREQS.forEach((freq, i) => {
  const band = document.createElement("div");
  band.className = "eq-band";
  band.innerHTML = `
    <span class="band-val" id="band-val-${i}">0</span>
    <input type="range" min="${BAND_MIN[i]}" max="12" value="0" step="1" id="band-${i}">
    <span class="band-label">${freq}</span>
  `;
  eqBands.appendChild(band);
  const slider = document.getElementById(`band-${i}`);
  const valEl = document.getElementById(`band-val-${i}`);
  slider.addEventListener("input", () => {
    valEl.textContent = fmtGain(+slider.value);
  });
  slider.addEventListener("change", () => {
    if (activeSlot === -1) {
      officialCustomGains = bandSliders.map((sl) => +sl.value);
      cmd("setCustomEQ", { fs: FS, gain: officialCustomGains });
    }
  });
  bandSliders.push(slider);
});

function fmtGain(v) {
  return (v > 0 ? "+" : "") + v.toFixed(1);
}

function loadGains(gains) {
  bandSliders.forEach((sl, i) => {
    const v = Math.max(BAND_MIN[i], Math.min(12, gains[i] || 0));
    sl.value = v;
    document.getElementById(`band-val-${i}`).textContent = fmtGain(v);
  });
}

function deactivateAll() {
  document
    .querySelectorAll(".preset-btn")
    .forEach((b) => b.classList.remove("active"));
}

// Preset-grid: native presets (id 1–4) + official Custom (id 0)
document.querySelectorAll("#preset-grid .preset-btn").forEach((btn) => {
  btn.addEventListener("click", () => {
    deactivateAll();
    btn.classList.add("active");
    const id = +btn.dataset.id;
    if (id === 0) {
      // Official custom — restore its gains to device, show sliders
      activeSlot = -1;
      customEq.classList.remove("hidden");
      customEq.classList.add("official-custom");
      loadGains(officialCustomGains);
      cmd("setCustomEQ", { fs: FS, gain: officialCustomGains });
    } else {
      activeSlot = null;
      customEq.classList.add("hidden");
      cmd("setActiveEQ", { eq_id: id });
    }
  });
});

// Custom-preset-grid: user slots 0/1/2
document.querySelectorAll("#custom-preset-grid .preset-btn").forEach((btn) => {
  btn.addEventListener("click", () => {
    const slot = +btn.dataset.slot;
    deactivateAll();
    btn.classList.add("active");
    activeSlot = slot;
    customEq.classList.remove("hidden");
    customEq.classList.remove("official-custom");
    customNameEl.value = customNames[slot];
    loadGains(customGains[slot]);
    cmd("setCustomEQ", { fs: FS, gain: customGains[slot] });
  });
});

document.getElementById("btn-reset-eq")?.addEventListener("click", () => {
  const zeros = [0, 0, 0, 0, 0, 0, 0];
  loadGains(zeros);
  if (activeSlot === -1) {
    officialCustomGains = zeros.slice();
    cmd("setCustomEQ", { fs: FS, gain: officialCustomGains });
  }
});

document.getElementById("btn-apply-eq")?.addEventListener("click", () => {
  const gain = bandSliders.map((sl) => +sl.value);
  if (activeSlot === -1) {
    // Official custom: send to device, remember in memory (not saved to config)
    officialCustomGains = gain;
    cmd("setCustomEQ", { fs: FS, gain });
  } else if (activeSlot !== null) {
    // User custom slot: send to device + persist to config
    customGains[activeSlot] = gain;
    const name = customNameEl.value.trim() || `Custom ${activeSlot + 1}`;
    customNames[activeSlot] = name;
    const btn = document.querySelector(
      `#custom-preset-grid .preset-btn[data-slot="${activeSlot}"]`,
    );
    if (btn) btn.textContent = name;
    cmd("setCustomEQ", { fs: FS, gain });
    cmd("saveCustomPreset", { slot: activeSlot, name, gain });
  }
});

// ── Moment ────────────────────────────────────────────────────────────────────

const momentToggle = document.getElementById("moment-toggle");
const momentControls = document.getElementById("moment-controls");

// Start dimmed (off by default); updateState corrects if device is already playing.
momentControls.style.opacity = "0.4";
momentControls.style.pointerEvents = "none";

let momentActive = false;
let momentSoundscapeId = 1;
let momentSleepTimer = 0;
// Per-mode element volumes [element 0-2], default 50%.
const momentElemVols = {
  1: [50, 50, 50],
  2: [50, 50, 50],
  3: [50, 50, 50],
  4: [50, 50, 50],
};

function renderMomentElements(modeId) {
  const keys = SOUNDSCAPE_ELEMENTS[modeId];
  const vols = momentElemVols[modeId];
  for (let i = 0; i < 3; i++) {
    const labelEl = document.getElementById(`el-label-${i}`);
    const valEl = document.getElementById(`el-val-${i}`);
    const sliderEl = document.getElementById(`el-slider-${i}`);
    if (labelEl) {
      labelEl.dataset.i18n = keys[i];
      labelEl.textContent = window.t(keys[i]);
    }
    if (valEl) valEl.textContent = vols[i] + "%";
    if (sliderEl) sliderEl.value = vols[i];
  }
}

// Initialize labels for default mode (Forest).
renderMomentElements(1);

momentToggle?.addEventListener("change", () => {
  momentActive = momentToggle.checked;
  momentControls.style.opacity = momentActive ? "1" : "0.4";
  momentControls.style.pointerEvents = momentActive ? "" : "none";
  if (momentActive) {
    cmd("momentOn", {
      soundscape_id: momentSoundscapeId,
      sleep_timer: momentSleepTimer,
    });
  } else {
    cmd("momentOff", { soundscape_id: momentSoundscapeId });
  }
});

document.querySelectorAll(".timer-btn").forEach((btn) => {
  btn.addEventListener("click", () => {
    document
      .querySelectorAll(".timer-btn")
      .forEach((b) => b.classList.remove("active"));
    btn.classList.add("active");
    momentSleepTimer = +btn.dataset.timer;
    if (momentActive) {
      cmd("momentConfig", {
        soundscape_id: momentSoundscapeId,
        sleep_timer: momentSleepTimer,
      });
    }
  });
});

document.querySelectorAll(".mode-btn").forEach((btn) => {
  btn.addEventListener("click", () => {
    document
      .querySelectorAll(".mode-btn")
      .forEach((b) => b.classList.remove("active"));
    btn.classList.add("active");
    const newId = +btn.dataset.id;
    const prevId = momentSoundscapeId;
    momentSoundscapeId = newId;
    renderMomentElements(newId);
    if (momentActive) {
      cmd("momentSwitch", {
        prev_id: prevId,
        soundscape_id: newId,
        sleep_timer: momentSleepTimer,
      });
    }
  });
});

for (let i = 0; i < 3; i++) {
  const slider = document.getElementById(`el-slider-${i}`);
  const valEl = document.getElementById(`el-val-${i}`);
  slider?.addEventListener("input", () => {
    valEl.textContent = slider.value + "%";
  });
  slider?.addEventListener("change", () => {
    momentElemVols[momentSoundscapeId][i] = +slider.value;
    cmd("momentSetElement", {
      soundscape_id: momentSoundscapeId,
      element_id: i,
      value: +slider.value,
    });
  });
}

// ── State sync from Go backend ────────────────────────────────────────────────

window.updateState = function (json) {
  const s = typeof json === "string" ? JSON.parse(json) : json;
  if (s.theme) window.applyTheme(s.theme);
  if (s.language) window.setLanguage(s.language);
  if (s.deviceName)
    document.getElementById("device-name").textContent = s.deviceName;
  if (s.connected !== undefined) {
    document.getElementById("status-dot").className =
      "status-dot " + (s.connected ? "connected" : "");
  }
  if (s.lights) {
    const l = s.lights;
    if (l.enable !== undefined) lightsToggle.checked = l.enable;
    if (l.brightness !== undefined) {
      brightnessEl.value = l.brightness;
      brightnessVal.textContent = l.brightness + "%";
    }
    if (l.speed !== undefined) {
      speedEl.value = l.speed;
      speedVal.textContent = l.speed;
    }
    if (l.patternId !== undefined) {
      activePatternId = l.patternId;
      document
        .querySelectorAll(".pattern-btn")
        .forEach((b) =>
          b.classList.toggle("active", +b.dataset.id === l.patternId),
        );
      updateColorStrip(l.patternId);
      updateSpeedLock(l.patternId);
    }
    if (l.colorLevel !== undefined) {
      patternColorLevel[activePatternId] = l.colorLevel;
      colorHue.value = l.colorLevel;
    }
  }
  if (s.eq) {
    const e = s.eq;
    if (e.presetId !== undefined && e.presetId > 0) {
      // Native preset became active (e.g. changed from phone)
      deactivateAll();
      const btn = document.querySelector(
        `#preset-grid .preset-btn[data-id="${e.presetId}"]`,
      );
      if (btn) btn.classList.add("active");
      activeSlot = null;
      customEq.classList.add("hidden");
    }
    if (e.customGain) {
      // Block only when a user-custom slot is active: EQID=0 then holds user-custom
      // data, not official custom. All other states (native preset or official custom
      // active) are safe to sync — phone changes to official custom flow through.
      if (activeSlot === null || activeSlot === -1)
        officialCustomGains = e.customGain;
      if (activeSlot === -1) loadGains(officialCustomGains);
    }
  }
  if (s.player) {
    const p = s.player;
    updateVolIcon(p.vol, p.mute);
    if (volSliderEl && p.vol !== undefined) {
      volSliderEl.value = p.vol;
      volValueEl.textContent = p.vol + "%";
    }
  }
  if (s.momentElements) {
    for (const modeId in s.momentElements) {
      const id = +modeId;
      if (momentElemVols[id]) momentElemVols[id] = s.momentElements[modeId];
    }
    renderMomentElements(momentSoundscapeId);
  }
  if (s.moment) {
    const m = s.moment;
    if (m.enabled !== undefined) {
      momentToggle.checked = m.enabled;
      momentActive = m.enabled;
      momentControls.style.opacity = m.enabled ? "1" : "0.4";
      momentControls.style.pointerEvents = m.enabled ? "" : "none";
    }
    if (m.soundscapeId !== undefined) {
      momentSoundscapeId = m.soundscapeId;
      document
        .querySelectorAll(".mode-btn")
        .forEach((b) =>
          b.classList.toggle("active", +b.dataset.id === m.soundscapeId),
        );
      renderMomentElements(m.soundscapeId);
    }
    if (m.sleepTimer !== undefined) {
      momentSleepTimer = m.sleepTimer;
      document
        .querySelectorAll(".timer-btn")
        .forEach((b) =>
          b.classList.toggle("active", +b.dataset.timer === m.sleepTimer),
        );
    }
  }
  if (s.customPresets) {
    s.customPresets.forEach((p, i) => {
      if (p.name) {
        customNames[i] = p.name;
        const btn = document.querySelector(
          `#custom-preset-grid .preset-btn[data-slot="${i}"]`,
        );
        if (btn) btn.textContent = p.name;
      }
      if (p.gain && p.gain.length === 7) customGains[i] = p.gain;
    });
  }
};

// Preview navigation
window.showScreen = function (id) {
  document
    .querySelectorAll(".screen")
    .forEach((s) => s.classList.add("hidden"));
  document.getElementById(id).classList.remove("hidden");
  document.querySelectorAll("#preview-nav button").forEach((b) => {
    b.classList.toggle("active", b.getAttribute("onclick").includes(id));
  });
};

// Preview with mock state
if (MOCK) {
  window.updateState(state);
  // Hide preview-nav in production (Go injects window.sendCmd before page loads)
} else {
  document.getElementById("preview-nav")?.remove();
  // Auto-resize the native window to match content height.
  const _ro = new ResizeObserver((entries) => {
    cmd("setWindowHeight", {
      height: Math.ceil(entries[0].contentRect.height),
    });
  });
  _ro.observe(document.body);
}
